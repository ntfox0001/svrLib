package wxAccessRefMsg

type WxAccessTokenReq struct {
	AppId string `json:"appId"`
}
type WxAccessTokeyResp struct {
	Token string `json:"token"`
}

type WxTicketReq struct {
	AppId string `json:"appId"`
}
type WxTicketResp struct {
	Ticket string `json:"ticket"`
}

// 微信回复格式
type WxMpRefreshAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxMpRefreshTickedResp struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}
