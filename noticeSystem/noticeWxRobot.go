package noticeSystem

import (
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/util"
	"sync"
	"time"

	"github.com/ntfox0001/svrLib/log"

	jsoniter "github.com/json-iterator/go"
)

type wxRobot struct {
	ip   string
	port string

	myselfId   string
	friendList []NoticeWxRobotUserInfo
	roomList   []NoticeWxRobotRoomInfo
	wxInfoLock sync.RWMutex
	quitChan   chan interface{}

	templates []WxRobotTemplateCfg

	// url commend
	getListUrl     string
	sendMsgUrl     string
	getMsgUrl      string
	getPicUrl      string
	sendPicUrl     string
	getRoomListUrl string
	sendRoomMsgUrl string
	getRoomMsgUrl  string
}

func newWxRobot(ip, port string, refreshTime uint, templates []WxRobotTemplateCfg) (*wxRobot, error) {
	robot := &wxRobot{
		ip:         ip,
		port:       port,
		templates:  templates,
		friendList: nil,
		quitChan:   make(chan interface{}),
	}
	robot.getListUrl = fmt.Sprintf("http://%s:%s/wxrobot/getList", ip, port)
	robot.sendMsgUrl = fmt.Sprintf("http://%s:%s/wxrobot/sendMsg", ip, port)
	robot.getMsgUrl = fmt.Sprintf("http://%s:%s/wxrobot/getMsg", ip, port)
	robot.getPicUrl = fmt.Sprintf("http://%s:%s/wxrobot/getPic", ip, port)
	robot.sendPicUrl = fmt.Sprintf("http://%s:%s/wxrobot/sendPic", ip, port)

	robot.getRoomListUrl = fmt.Sprintf("http://%s:%s/wxrobot/getRoomList", ip, port)
	robot.sendRoomMsgUrl = fmt.Sprintf("http://%s:%s/wxrobot/sendRoomMsg", ip, port)
	robot.getRoomMsgUrl = fmt.Sprintf("http://%s:%s/wxrobot/getRoomData", ip, port)

	// // 获取好友列表
	// if myselfId, friendList, err := robot.getList(); err != nil {
	// 	return nil, err
	// } else {
	// 	robot.myselfId = myselfId
	// 	robot.friendList = friendList
	// }

	// // 获取群列表
	// if roomlist, err := robot.getQun(robot.myselfId); err != nil {
	// 	return nil, err
	// } else {
	// 	robot.roomList = roomlist
	// }
	robot.refreshWxInfo()

	t := time.NewTicker(time.Second * time.Duration(refreshTime))

	go func() {
	runable:
		for {
			select {
			case <-robot.quitChan:
				break runable
			case <-t.C:
				robot.refreshWxInfo()
			}
		}
	}()

	return robot, nil
}

func (r *wxRobot) refreshWxInfo() {

	// 获取好友列表
	myselfId, friendList, err := r.getList()
	if err != nil {
		log.Error("wxrobot refresh error", "err", err.Error())
		return
	}

	// 获取群列表
	roomlist, err := r.getQun(myselfId)
	if err != nil {
		log.Error("wxrobot refresh error", "err", err.Error())
		return
	}

	defer r.wxInfoLock.Unlock()
	r.wxInfoLock.Lock()
	r.myselfId = myselfId
	r.friendList = friendList
	r.roomList = roomlist
}

func (r *wxRobot) close() {
	r.quitChan <- struct{}{}
}

func (r *wxRobot) getList() (string, []NoticeWxRobotUserInfo, error) {
	rt, err := network.SyncHttpPost(r.getListUrl, "{}", network.ContentTypeJson)
	if err != nil {
		log.Error("getlisturl error ", "err", err.Error())
		return "", nil, err
	}

	var resp NoticeWxRobotGetListResp
	jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(rt, &resp)

	if resp.ErrorId != "0" {
		return "", nil, commonError.NewStringErr(resp.ErrorId)
	}

	if !(len(resp.List) > 0 && len(resp.List[0]) > 0) {
		return "", nil, commonError.NewStringErr("There is not robot connect the server.")
	}
	// 目前只取第一个机器人使用
	return resp.List[0][0].UserName, resp.List[0], nil
}

func (r *wxRobot) getQun(id string) ([]NoticeWxRobotRoomInfo, error) {
	req := NoticeWxRobotGetRoomListReq{
		UserName: id,
	}
	js, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)

	rt, _ := network.SyncHttpPost(r.getRoomListUrl, js, network.ContentTypeJson)

	var resp NoticeWxRobotGetRoomListResp
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(rt, &resp); err != nil {
		log.Error("wxRobot error: Invalid to format of getRoomList", "js", js)
		return nil, commonError.NewStringErr("wxRobot error: Invalid to format of getRoomList")
	}

	return resp.List, nil
}

func (r *wxRobot) GetMyselfId() string {
	r.wxInfoLock.RLock()
	myselfId := r.myselfId
	r.wxInfoLock.RUnlock()
	return myselfId
}

func (r *wxRobot) GetRoomList() (myself string, roomList []string) {
	r.wxInfoLock.RLock()
	roomList = make([]string, 0, len(r.roomList))
	for _, i := range r.roomList {
		roomList = append(roomList, i.UserName)
	}
	myself = r.myselfId
	r.wxInfoLock.RUnlock()
	return
}

func (r *wxRobot) GetFriendList() (myself string, friendList []string) {
	r.wxInfoLock.RLock()
	friendList = make([]string, 0, len(r.friendList))
	for _, i := range r.friendList {
		friendList = append(friendList, i.UserName)
	}
	myself = r.myselfId
	r.wxInfoLock.RUnlock()
	return
}

// 群发好友
func (r *wxRobot) send(data map[string]string) (rtErr error) {
	myself, target := r.GetFriendList()

	templ, err := r.GetTemplateFromType(data["{type}"])
	if err != nil {
		log.Error("wxrobot room send error: Template does not exist.", "type", data["{type}"])
		return err
	}

	return r._send(myself, target, util.StringReplace(templ.Msg, data))
}

// 群发群
func (r *wxRobot) roomSend(data map[string]string) (rtErr error) {
	myself, target := r.GetRoomList()

	templ, err := r.GetTemplateFromType(data["{type}"])
	if err != nil {
		log.Error("wxrobot room send error: Template does not exist.", "type", data["{type}"])
		return err
	}

	return r._send(myself, target, util.StringReplace(templ.Msg, data))
}

func (r *wxRobot) _send(myself string, target []string, msg string) (rtErr error) {
	if len(target) == 0 {
		log.Warn("There is nothing what robot's friend list.")
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			rtErr = err.(error)
			log.Error("wxrobot room send error", "Error", rtErr.Error())
			return
		}
	}()

	req := NoticeWxRobotSendMsgReq{
		Target: make([]string, 0, 10),
	}

	req.UserName = myself
	req.Msg = msg
	req.Target = target

	js, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Error("wxrobot room send error: Invalid to formate of NoticeWxRobotSendMsgReq")
		return err
	}
	rs, err := network.SyncHttpPost(r.sendMsgUrl, js, network.ContentTypeJson)
	if err != nil {
		log.Error("wxrobot room send error: Post error", "err", err.Error())
		return err
	}

	var resp NoticeWxRobotSendMsgResp
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(rs, &resp); err != nil {
		log.Error("wxrobot room send error: Invalid to formate of NoticeWxRobotSendMsgResp", "js", rs)
		return err
	}

	if resp.ErrorId != "0" {
		log.Error("wxrobot room send error", "errorId", resp.ErrorId)
		return commonError.NewStringErr(resp.ErrorId)
	}
	return nil
}

func (r *wxRobot) GetTemplateFromType(noticeType string) (*WxRobotTemplateCfg, error) {
	for _, v := range r.templates {
		if v.Type == noticeType {
			return &v, nil
		}
	}
	log.Error("wxrobot", "TemplateType does not exist", noticeType)
	return nil, commonError.NewStringErr("TemplateType does not exist.")
}

func (r *wxRobot) Send(myselfId string, target []string, msg string) error {
	return r._send(myselfId, target, msg)
}

func (r *wxRobot) GetRoomMsg(target string) ([]NoticeWxRobotRoomMsgItem, error) {
	req := NoticeWxRobotGetRoomMsgReq{
		RoomId: target,
	}
	js, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Error("wxrobot room GetRoomMsg error: Invalid to formate of NoticeWxRobotGetRoomMsgReq")
		return nil, err
	}
	rs, err := network.SyncHttpPost(r.getRoomMsgUrl, js, network.ContentTypeJson)
	if err != nil {
		log.Error("wxrobot room GetRoomMsg error: Post error", "err", err.Error())
		return nil, err
	}

	var resp NoticeWxRobotGetRoomMsgResp
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(rs, &resp); err != nil {
		log.Error("wxrobot room GetRoomMsg error: Invalid to formate of NoticeWxRobotGetRoomMsgResp", "js", rs)
		return nil, err
	}

	if resp.ErrorId != "recordNone" && resp.ErrorId != "0" {
		log.Error("wxrobot room GetRoomMsg error", "errorId", resp.ErrorId)
		return nil, commonError.NewStringErr(resp.ErrorId)
	}

	return resp.List, nil
}

func (r *wxRobot) GetMsg(myselfId string) (map[string][]NoticeWxRobotMsgItem, error) {
	req := NoticeWxRobotGetMsgReq{
		UserName: myselfId,
	}
	js, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(req)
	if err != nil {
		log.Error("wxrobot GetMsg error: Invalid to formate of NoticeWxRobotGetMsgReq")
		return nil, err
	}
	rs, err := network.SyncHttpPost(r.getMsgUrl, js, network.ContentTypeJson)
	if err != nil {
		log.Error("wxrobot GetMsg error: Post error", "err", err.Error())
		return nil, err
	}

	var resp NoticeWxRobotGetMsgResp
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(rs, &resp); err != nil {
		log.Error("wxrobot GetMsg error: Invalid to formate of NoticeWxRobotGetMsgResp", "js", rs)
		return nil, err
	}

	if resp.ErrorId != "recordNone" && resp.ErrorId != "0" {
		log.Error("wxrobot GetMsg error", "errorId", resp.ErrorId)
		return nil, commonError.NewStringErr(resp.ErrorId)
	}

	return resp.Msg, nil
}
