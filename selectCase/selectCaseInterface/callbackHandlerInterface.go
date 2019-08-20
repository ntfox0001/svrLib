package selectCaseInterface

type ICallbackHandler interface {
	GetUserData() interface{}
	SendReturnMsg(msg EventChanMsg)
	SendReturnMsgNoReturn(data interface{})
}
