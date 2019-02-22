package noticeSystem

type NoticeMsg struct {
	Phone []string `json:"phoneNumbers"`
	Args  []string `json:"args"`
}

// sms msg data
type noticeSdkMsg struct {
	AppId      string `json:"appid"`
	AppKey     string `json:"appkey"`
	TemplateId string `json:"templateId"`
	SmsSign    string `json:"smsSign"`
	NoticeMsg
}

type NoticeResp struct {
	ReturnString string
	Err          error
}

// wx mp msg
type WxMpMsg struct {
	Touser      string                 `json:"touser"`
	Template_Id string                 `json:"template_id"`
	Url         string                 `json:"url"`
	Data        map[string]WxMpMsgData `json:"data"`
}

type WxMpMsgData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type WxMpMsgResp struct {
	Error  string `json:"error"`
	ErrMsg string `json:"errmsg"`
	MsgId  int    `json:"msgid"`
}

// phone
type PhoneMsgReq struct {
	CalledNumber   string
	Id             string `json:"TtsCode"`
	CallShowNumber string
	// ttsparam是一个json字符串
	TtsParam string
}

type PhoneMsgResp struct {
	RequestId string
	Code      string
	Message   string
	CallId    string
}
