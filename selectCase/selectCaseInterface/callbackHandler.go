package selectCaseInterface

import "github.com/ntfox0001/svrLib/log"

type CallbackHandler struct {
	callbackMsgName string
	sender          ISelectLoopHelper
	userData        interface{}
}

func NewCallbackHandler(returnMsg string, sender ISelectLoopHelper, userData interface{}) *CallbackHandler {
	return &CallbackHandler{
		callbackMsgName: returnMsg,
		sender:          sender,
		userData:        userData,
	}
}

func (c *CallbackHandler) GetUserData() interface{} {
	return c.userData
}

// 向调用者发送一个消息，userdata回自动返回
func (c *CallbackHandler) SendReturnMsg(msg EventChanMsg) {
	if msg.UserData != nil {
		log.Error("CallbackHandler UserData must be nil", "msg", msg.MsgId, "content", msg.Content)
		return
	}
	msg.UserData = c.userData
	c.sender.SendMsgToMe(msg)
}

// 向调用者发送一个消息，不需要对方回答时使用，userdata回自动返回
func (c *CallbackHandler) SendReturnMsgNoReturn(data interface{}) {
	msg := NewEventChanMsg(c.callbackMsgName, nil, data)
	msg.UserData = c.userData
	c.sender.SendMsgToMe(msg)
}
