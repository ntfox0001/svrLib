package userDefine

const (
	UserNotExist = "100404"
)

type ReqBase struct {
	MsgId string `json:"msgId"`
}
type RespBase struct {
	MsgId   string `json:"msgId"`
	ErrorId string `json:"errorId"`
}

// websocket--------------------------------------------------------------------------------
type UserLoginReq struct {
	ReqBase
	Token string `json:"token"`
}

type UserLoginResp struct {
	RespBase
	MsgId string `json:"msgId"`
}

// http ---------------------------------------------------------------------------------------
type HttpMsgReqBase struct {
	UserId int `json:"userId,string"`
}
type HttpMsgRespBase struct {
	ErrorId string `json:"errorId"`
}
