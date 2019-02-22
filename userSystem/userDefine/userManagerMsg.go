package userDefine

import (
	"github.com/ntfox0001/svrLib/network/networkInterface"
)

type UserPair struct {
	Ac      networkInterface.IMsgHandler
	UnionId string
}
type FindTokenReq struct {
	Token      string
	WaitWxChan chan string
}

func NewFindTokenReq(token string, waitWxChan chan string) FindTokenReq {
	return FindTokenReq{
		Token:      token,
		WaitWxChan: waitWxChan,
	}
}

type WxMpLoginReq struct {
	UserData
}

type GenerateTokenReq struct {
	WaitTokenChan chan GenerateTokenResp
	UserData
}

type GenerateTokenResp struct {
	Token string
	UserData
}

type WxMpLoginResp struct {
	MsgId   string `json:"msgId"`
	Token   string `json:"token"`
	UserId  int    `json:"userId,string"`
	ErrorId string `json:"errorId"`
}

type NewUserInfoReq struct {
	UserData
	WaitTokenChan chan GenerateTokenResp
}

// wx code login
type WxMpCodeLoginReq struct {
	AppId string `json:"appId"`
	Code  string `json:"code"`
}

type WxMpCodeLoginResp struct {
}
