package noticeSystem

// 获取好友id列表
type NoticeWxRobotGetListReq struct {
}

type NoticeWxRobotUserInfo struct {
	NickName   string
	HeadImgUrl string
	Sex        string
	Province   string
	Signature  string
	UserName   string
}
type NoticeWxRobotGetListResp struct {
	List    [][]NoticeWxRobotUserInfo `json:"list"`
	ErrorId string                    `json:"errorId"`
}

// 微信机器人发送消息
type NoticeWxRobotSendMsgReq struct {
	UserName string   `json:"userName"`
	Target   []string `json:"target"`
	Msg      string   `json:"msg"`
}
type NoticeWxRobotSendMsgResp struct {
	ErrorId string `json:"errorId"`
}

// 获取好友消息
type NoticeWxRobotGetMsgReq struct {
	UserName string `json:"userName"`
}
type NoticeWxRobotMsgItem struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

type NoticeWxRobotGetMsgResp struct {
	ErrorId string                            `json:"errorId"`
	Robot   string                            `json:"robot"`
	Msg     map[string][]NoticeWxRobotMsgItem `json:"msg"`
}

// 群 --------------------------------------------------------------------------------
// 获取群列表
type NoticeWxRobotGetRoomListReq struct {
	UserName string `json:"userName"`
}
type NoticeWxRobotRoomInfo struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
}
type NoticeWxRobotGetRoomListResp struct {
	List    []NoticeWxRobotRoomInfo `json:"list"`
	ErrorId string                  `json:"errorId"`
}

// 请求指定群的消息
type NoticeWxRobotGetRoomMsgReq struct {
	RoomId string `json:"roomId"`
}

type NoticeWxRobotRoomMsgItem struct {
	UserName string `json:"userName"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}

type NoticeWxRobotGetRoomMsgResp struct {
	List    []NoticeWxRobotRoomMsgItem `json:"list"`
	ErrorId string                     `json:"errorId"`
}
