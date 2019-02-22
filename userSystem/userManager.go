package userSystem

import (
	"github.com/ntfox0001/svrLib/userSystem/userDefine"

	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"github.com/ntfox0001/svrLib/selectCase"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"net/http"

	"reflect"

	"github.com/ntfox0001/svrLib/database"

	"github.com/ntfox0001/svrLib/util"

	"github.com/ntfox0001/svrLib/log"

	"github.com/gorilla/websocket"
)

type UserManager struct {
	wsUserMap     map[networkInterface.IMsgHandler]*User
	userToken     map[string]userDefine.UserData // token和wxunionid的对应，用户登录后就删除
	userMap       map[string]*User               // key是wxunionid
	userUserIdMap map[int]string                 // userId：unionId
	svrip         string
	svrport       string
	selectLoop    *selectCase.SelectLoop
	goPool        *util.GoroutinePool
	mgrCallback   IUserManagerCallback
	usrCallback   IUserCallback
}

/*
	user登陆流程：
	需要实现一个php微信登陆，然后构造一个UserData传给服务器


*/

func NewUserManager(ip string, port string, usrMgrcb IUserManagerCallback, usrcb IUserCallback) *UserManager {

	usrMgr := &UserManager{
		wsUserMap:     make(map[networkInterface.IMsgHandler]*User),
		userToken:     make(map[string]userDefine.UserData),
		userMap:       make(map[string]*User),
		userUserIdMap: make(map[int]string),
		svrip:         ip,
		svrport:       port,
		selectLoop:    selectCase.NewSelectLoop("user manager", 20, 30),
		goPool:        util.NewGoPool("user manager", 20, 20),
		mgrCallback:   usrMgrcb,
		usrCallback:   usrcb,
	}
	// 注册user manager事件
	usrMgr.registerEvent()

	// 加载所有用户数据
	if err := usrMgr.loadAllUser(); err != nil {
		return nil
	}

	if err := usrMgr.mgrCallback.OnInitial(usrMgr.GetSelectLoopHelper()); err != nil {
		log.Error("callback initial error", "err", err.Error())
		return nil
	}

	// 运行
	go usrMgr.selectLoop.Run()

	return usrMgr
}

func (m *UserManager) loadAllUser() error {

	// 去数据库加载用户数据
	op := database.Instance().NewOperation("call playerLoadAll()")
	if rt, err := database.Instance().SyncExecOperation(op); err != nil {
		return err
	} else {
		// 1, load player talbe
		usrDS := rt.FirstSet()
		if usrDS != nil {
			log.Info("userManager", "loadUser", len(usrDS))
			for _, v := range usrDS {
				var usrData userDefine.UserData
				if err := util.I2Stru(v, &usrData); err == nil {

					// 创建新user
					usr := NewUser(usrData, m.usrCallback)
					// 调用user回调
					m.mgrCallback.OnInitUser(usr)

					m.userMap[usrData.UnionId] = usr
					m.userUserIdMap[usrData.UserId] = usrData.UnionId
				} else {
					log.Error("userManager", "loadAllUserInfo error", v)
					return err
				}
			}
		} else {
			return nil
		}

		// 运行user
		for _, u := range m.userMap {
			u.Run()
		}
	}

	return nil
}

func (m *UserManager) registerEvent() {
	m.GetSelectLoopHelper().RegisterEvent("AttachUser", m.attachUser)
	m.GetSelectLoopHelper().RegisterEvent("DeattachUser", m.deattachUser)

	// 发送一个函数到user，并在user协程内执行
	m.GetSelectLoopHelper().RegisterEvent("RunInUserReq", m.runInUserReq)
	m.GetSelectLoopHelper().RegisterEvent("RunInUserManagerReq", m.runInUserManagerReq)

	// user加载
	m.GetSelectLoopHelper().RegisterEvent("FindTokenReq", m.findTokenReq)
	m.GetSelectLoopHelper().RegisterEvent("GenerateTokenReq", m.generateTokenReq)
	m.GetSelectLoopHelper().RegisterEvent("NewUserInfoResp", m.newUserInfoResp)
}

func (m *UserManager) CheckConn(w http.ResponseWriter, r *http.Request) bool {
	return true //network.CheckSameOrigin(r, m.svrip, m.svrport)
}
func (m *UserManager) Fetch(ac networkInterface.IMsgHandler) bool {

	waitWxChan := make(chan string, 1)
	// 判断链接合法性, 使用函数调用方式
	ac.RegisterJsonMsg("UserLoginReq", func(msg map[string]interface{}) {
		var req userDefine.UserLoginReq
		if err := util.I2Stru(msg, req); err == nil {
			m.GetSelectLoopHelper().SendMsgToMe(
				selectCaseInterface.NewEventChanMsg("FindTokenReq", nil, userDefine.NewFindTokenReq(req.Token, waitWxChan)))
		}
	})
	// 接收第一个消息
	ac.ProcessMsg()

	// 等待结果
	rt := util.WaitChanWithTimeout(reflect.ValueOf(waitWxChan), 10)
	if rt == nil {
		// 关闭这个通道，避免发送端锁死
		close(waitWxChan)
		return true
	}
	wxUnionId := rt.(string)
	if wxUnionId == "" {
		// 没有找到token
		log.Warn("invalid to Connect of Arrived")
		return false
	}

	m.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("AttachUser", nil, userDefine.UserPair{Ac: ac, UnionId: wxUnionId}))

	log.Debug("+ user arrived", "request", ac.(*network.WsMsgHandler).GetRequest().RemoteAddr)
	return true
}
func (m *UserManager) Close(ac networkInterface.IMsgHandler) {
	m.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("DeattachUser", nil, userDefine.UserPair{Ac: ac, UnionId: ""}))
	log.Debug("- user left", "request", ac.(*network.WsMsgHandler).GetRequest().RemoteAddr)
}

func (m *UserManager) NewMsgHandler(c *websocket.Conn, r *http.Request) networkInterface.IMsgHandler {
	return network.NewMsgHander(c, r)
}

func (m *UserManager) Release() {
	m.mgrCallback.OnRelease()
	for _, v := range m.userMap {
		log.Debug("release user", "userId", v.UserInfo().GetUserData().UserId)
		v.Release()
	}
	m.goPool.Release()
	m.selectLoop.Close()
	log.Debug("UserManager release.")
}
func (m *UserManager) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return m.selectLoop.GetHelper()
}

func (m *UserManager) attachUser(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	up := msg.Content.(userDefine.UserPair)

	if usr, ok := m.userMap[up.UnionId]; ok {
		m.wsUserMap[up.Ac] = usr

	}

	return true
}
func (m *UserManager) deattachUser(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	up := msg.Content.(userDefine.UserPair)
	m.wsUserMap[up.Ac].DeattachWSAcceptConn()
	delete(m.wsUserMap, up.Ac)
	return true
}

func (m *UserManager) findTokenReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	req := msg.Content.(userDefine.FindTokenReq)

	if wxid, ok := m.userToken[req.Token]; ok {
		req.WaitWxChan <- wxid.UnionId
		delete(m.userToken, req.Token)
	} else {
		req.WaitWxChan <- ""
	}

	return true
}

func (m *UserManager) generateTokenReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	req := msg.Content.(userDefine.GenerateTokenReq)

	var resp userDefine.GenerateTokenResp
	//加载userinfo
	if usr, ok := m.userMap[req.UnionId]; ok {
		token := util.NewToken(req.UnionId)
		m.userToken[token] = req.UserData

		// 如果user已经加载，那么要刷新一下wx信息
		usr.UpdateWxInfo(req.UserData)

		resp.UserData = *usr.UserInfo().GetUserData()
		resp.Token = token
		req.WaitTokenChan <- resp
	} else {
		// 新创建一个userInfo
		NewUIReq := userDefine.NewUserInfoReq{
			UserData:      req.UserData,
			WaitTokenChan: req.WaitTokenChan,
		}
		//向数据库插入新用户

		cb := m.GetSelectLoopHelper().NewCallbackHandler("NewUserInfoResp", NewUIReq)
		usr := AsyncNewUser(cb, req.UserData, m.usrCallback)
		m.userMap[req.UnionId] = usr

	}

	return true
}

// 公众号新用户登陆创建新用户
func (m *UserManager) newUserInfoResp(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	rt := msg.Content.(*database.DataResult)
	uiReq := rt.Opt.UserData.(userDefine.NewUserInfoReq)

	resp := userDefine.GenerateTokenResp{}

	if rt.Err != nil {
		resp.Token = ""
		uiReq.WaitTokenChan <- resp
	}
	// 读取数据
	userDS := rt.FirstSet()
	if userDS != nil && len(userDS) == 1 {
		v := userDS[0]
		var usrData userDefine.UserData
		if err := util.I2Stru(v, &usrData); err == nil {
			if usr, ok := m.userMap[usrData.UnionId]; ok {
				// 启动user
				usr.Run()
			} else {
				delete(m.userMap, usrData.UnionId)
				log.Error("userMap does not exist unionId", "unionId", usrData.UnionId, "nickName", usrData.Nickname)
			}

			// new token
			token := util.NewToken(usrData.UnionId)
			m.userToken[token] = uiReq.UserData

			m.userUserIdMap[usrData.UserId] = usrData.UnionId
			resp.Token = token
			resp.UserData = usrData
			uiReq.WaitTokenChan <- resp
		}

	} else {
		log.Debug("newUserInfoResp returnRowError", "wxName", uiReq.Nickname, "unionId", uiReq.UnionId)
		resp.Token = ""
		uiReq.WaitTokenChan <- resp
	}

	return true
}

func (m *UserManager) runInUserReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	req := msg.Content.(RunInUserReq)

	if unionId, ok := m.userUserIdMap[req.UserId]; ok {
		if usr, ok := m.userMap[unionId]; ok {
			usr.SelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("RunInUser", nil, req.F))
			return true
		}
	}
	req.F(nil)
	return true
}
func (m *UserManager) runInUserManagerReq(data interface{}) bool {
	msg := data.(selectCaseInterface.EventChanMsg)
	f := msg.Content.(func(usrMgr *UserManager))

	f(m)

	return true
}

func (m *UserManager) HasUser(userId int) bool {
	pChan := make(chan bool)
	m.RunInUser(userId, func(usr *User) {
		pChan <- usr != nil
	})
	return <-pChan
}

func (m *UserManager) GetOpenIdByUserId(userId int) string {
	if m.HasUser(userId) {
		openIdChan := make(chan string)
		m.RunInUser(userId, func(usr *User) {
			openIdChan <- usr.UserInfo().GetUserData().OpenId
		})
		return <-openIdChan
	}
	return ""
}

// 在user的协程中运行给定函数
func (m *UserManager) RunInUser(userId int, f func(usr *User)) {
	m.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("RunInUserReq", nil, RunInUserReq{UserId: userId, F: f}))
}

// 在userManager的协程中，运行制定的函数
func (m *UserManager) RunInUserManager(f func(usrMgr *UserManager)) {
	m.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("RunInUserManagerReq", nil, f))
}

// 在userManager的协程中，运行制定的函数，阻塞返回
func (m *UserManager) SyncRunInUserManager(f func(usrMgr *UserManager) interface{}) interface{} {
	rtchan := make(chan interface{})
	rtfunc := func(usrMgr *UserManager) {
		rtchan <- f(usrMgr)
	}
	m.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("RunInUserManagerReq", nil, rtfunc))
	return <-rtchan
}

// 在user的协程中运行给定函数，阻塞返回
func (m *UserManager) SyncRunInUser(userId int, f func(usr *User) interface{}) interface{} {
	rtchan := make(chan interface{})
	rtfunc := func(usr *User) {
		rtchan <- f(usr)
	}
	m.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg("RunInUserReq", nil, RunInUserReq{UserId: userId, F: rtfunc}))
	return <-rtchan
}
