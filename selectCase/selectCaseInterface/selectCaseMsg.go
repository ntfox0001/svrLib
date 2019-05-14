package selectCaseInterface

type EventChanMsg struct {
	MsgId    string
	Sender   *CallbackHandler
	Content  interface{}
	UserData interface{} // 用来传递额外数据,一般用来Resp时回调回来
}

// sender 跨selectLoop调用时，对方收到消息后，用sender回复消息
// 当多个相同item向一个管理类注册时，管理类在消息处理时，也必须通过sender分辨是哪个item的消息
// 一个selectLoop向自己注册消息，触发这类消息时，不需要填写sender，置nil即可
func NewEventChanMsg(msgId string, sender *CallbackHandler, content interface{}) EventChanMsg {
	return EventChanMsg{
		MsgId:    msgId,
		Sender:   sender,
		Content:  content,
		UserData: nil,
	}
}

type EventRegMsg struct {
	Reg         bool
	EventId     string
	EventFuncId uint64
	F           func(EventChanMsg)
}
