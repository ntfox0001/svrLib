package userSystem

import (
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/userSystem/userDefine"
	"github.com/ntfox0001/svrLib/util"
	"reflect"
	"time"

	"github.com/ntfox0001/svrLib/log"
)

const (
	id_quit          = 0
	id_heartbeat     = 1
	id_svrTimer      = 2
	id_commonHandler = 3
	id_sendJsonMsg   = 4
	id_sendMsg       = 5
	id_listen        = 6
	id_listenJson    = 7

	id_max = id_listenJson + 1

	heartbeatTimeoutCheck = 5 // 心跳检测间隔
)

type RunInUserReq struct {
	UserId int
	F      func(usr *User)
}

// 每一个user都在自己的协程里运行，所以在user中访问其他资源都要考虑多线程问题
// user不能有update，理论上为了减少服务器压力，
// 服务器不提供时间轮询的功能，类似功能都通过客户端实现
type User struct {
	acceptConn        networkInterface.IMsgHandler // ac
	attachACId        []uint64
	heartbeatMsg      map[string]interface{}
	heartbeat         util.TimeLimit
	timeoutCheck      time.Duration // 心跳检查
	sendJsonMsgChan   chan interface{}
	sendMsgChan       chan networkInterface.IMsgData
	acceptJsonMsgChan chan map[string]interface{}
	acceptMsgChan     chan *networkInterface.RawMsgData
	userInfo          *UserInfo
	selectLoop        *selectCase.SelectLoop
	callback          IUserCallback
}

// 从内存创建user，数据库已经创建完成
func NewUser(usrData userDefine.UserData, usrcb IUserCallback) *User {
	usr := newUser(usrcb)
	usr.userInfo = newUserInfoForUserData(&usrData)
	return usr
}

// 异步创建
func AsyncNewUser(cb *selectCaseInterface.CallbackHandler, usrData userDefine.UserData, usrcb IUserCallback) *User {
	usr := newUser(usrcb)
	usr.userInfo = asyncNewUserInfo(cb, &usrData)

	return usr
}
func newUser(usrcb IUserCallback) *User {
	hbmsg := make(map[string]interface{})
	hbmsg["msgId"] = "HeartbeatNotify"

	usr := &User{

		acceptConn:        nil,
		attachACId:        make([]uint64, 6, 6),
		heartbeatMsg:      hbmsg,
		heartbeat:         util.NewTimeLimit(6, int64(time.Second*40), 8, 2),
		timeoutCheck:      heartbeatTimeoutCheck,
		sendJsonMsgChan:   make(chan interface{}, 20),
		sendMsgChan:       make(chan networkInterface.IMsgData, 20),
		acceptJsonMsgChan: make(chan map[string]interface{}, 20),
		acceptMsgChan:     make(chan *networkInterface.RawMsgData, 20),
		userInfo:          nil,
		selectLoop:        selectCase.NewSelectLoop("user", 10, 10),
		callback:          usrcb,
	}

	usr.SelectLoopHelper().RegisterEvent("RunInUser", usr.runInUser)

	return usr
}
func (u *User) AttachWSAcceptConn(acceptConn networkInterface.IMsgHandler) {
	u.acceptConn = acceptConn

	// 替换消息处理函数
	u.acceptConn.SetDispatchMsgHandler(func(msg *networkInterface.RawMsgData) {
		u.acceptMsgChan <- msg
	})
	u.acceptConn.SetDispatchJsonMsgHandler(func(msg map[string]interface{}) {
		u.acceptJsonMsgChan <- msg
	})
	// 注册心跳
	u.SelectLoopHelper().RegisterEvent("HeartbeatNotify", func(msg interface{}) bool {
		//log.Debug("HeartbeatNotify", "msg", msg)
		u.heartbeat.Hit()
		u.SendJsonMsg(u.heartbeatMsg)
		return true
	})
	// 初始化第一个心跳
	u.heartbeat.Hit()

	// 加入心跳
	u.attachACId[0] = u.selectLoop.AddSelectCase(reflect.ValueOf(u.heartbeat.TouchChan()), func(data interface{}) bool {
		return false
	})

	// 服务器timer
	timerSecond := time.NewTimer(time.Second * u.timeoutCheck)
	u.attachACId[1] = u.selectLoop.AddSelectCase(reflect.ValueOf(timerSecond.C), func(data interface{}) bool {
		u.heartbeat.CheckIntervalMax()
		timerSecond.Reset(time.Second * u.timeoutCheck)
		return true
	})

	// 用户网络消息
	u.attachACId[2] = u.selectLoop.AddSelectCase(reflect.ValueOf(u.sendJsonMsgChan), func(data interface{}) bool {
		u.Msghandler().SendJsonMsg(data)
		return true
	})
	u.attachACId[3] = u.selectLoop.AddSelectCase(reflect.ValueOf(u.sendMsgChan), func(data interface{}) bool {
		u.Msghandler().SendMsg(data.(networkInterface.IMsgData))
		return true
	})
	u.attachACId[4] = u.selectLoop.AddSelectCase(reflect.ValueOf(u.acceptMsgChan), func(data interface{}) bool {
		//u.Msghandler().DispatchMsg(data.(networkInterface.MsgData))
		msg := data.(*networkInterface.RawMsgData)
		u.SelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(msg.Name(), nil, msg))
		return true
	})
	u.attachACId[6] = u.selectLoop.AddSelectCase(reflect.ValueOf(u.acceptJsonMsgChan), func(data interface{}) bool {
		//u.Msghandler().DispatchJsonMsg(data.(map[string]interface{}))
		if u.userInfo == nil || u.userInfo.usrData.UserId == 0 {
			log.Warn("User need initial before accept msg.", "msg", data)
			return true
		}

		if msg, ok := data.(map[string]interface{}); ok {
			// 强制设置userId
			msg["UserId"] = u.userInfo.usrData.UserId
			if msgId, ok := msg["msgId"]; ok {
				if sMsgId, ok := msgId.(string); ok {
					u.SelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(sMsgId, nil, data))
					return true
				}
			}
		}
		log.Warn("UserClientMsg format error", "msg", data)
		return true
	})
}
func (u *User) DeattachWSAcceptConn() {
	for _, v := range u.attachACId {
		u.selectLoop.RemoveSelectCase(v)
	}
}
func (u *User) Run() {
	// 更新realtimesystem
	go u.run()
}
func (u *User) run() {
	u.selectLoop.Run()
	log.Debug("- user circle quit.")
}

// 用于自己关闭自己
func (u *User) close() {
	// 关闭链接，最终会调用Release
	u.Msghandler().Disconnect()
}

// 只能用于usermanger关闭
func (u *User) Release() {
	u.selectLoop.Close() //关闭chan防止阻塞协程
	//u.goPool.Release()
	log.Debug("- user release.")
}

func (u *User) SelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return u.selectLoop.GetHelper()
}

// 非线程安全
func (u *User) Msghandler() networkInterface.IMsgHandler {
	return u.acceptConn
}

func (u *User) SendJsonMsg(msg interface{}) {
	//u.goPool.Go(func(data interface{}) {
	u.sendJsonMsgChan <- msg
	//}, nil)
}
func (u *User) SendMsg(msg networkInterface.IMsgData) {
	//u.goPool.Go(func(data interface{}) {
	u.sendMsgChan <- msg
	//}, nil)
}

func (u *User) UserInfo() *UserInfo {
	return u.userInfo
}

// 外部传入函数在user中执行
func (u *User) runInUser(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	f := msg.Content.(func(*User))
	f(u)
	return true
}

// 当用户再次登陆时，更新信息
func (u *User) UpdateWxInfo(usrData userDefine.UserData) {

}
