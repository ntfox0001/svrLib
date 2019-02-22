package userDefine

// 客户端通知服务器发起支付
type WxPayReq struct {
	UserId int `json:"userId"`
	PayId  int `json:"payId"`
}

// 服务器申请preid，返回客户端
type WxPayResp struct {
	ErrorId  string `json:"errorId"`
	PrePayId string `json:"prePayId"`
	Price    int    `json:"price"`
}

// 服务器收到微信支付成功通知后，处理完相关逻辑，发送消息到客户端
type WxPayNotify struct {
	ErrorId string `json:"errorId"`
}
