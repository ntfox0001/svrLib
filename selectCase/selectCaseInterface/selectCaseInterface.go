package selectCaseInterface

import (
	"reflect"
)

type ISelectLoopHelper interface {
	// 向loop发送消息
	SendMsgToMe(data EventChanMsg)
	// 向loop的common handler注册消息
	RegisterEvent(event string, f func(EventChanMsg)) uint64
	UnregisterEvent(id uint64)

	// 向loop中注册通道和处理函数
	AddSelectCaseFront(ch reflect.Value, cb func(data interface{}) bool) uint64
	AddSelectCase(ch reflect.Value, cb func(data interface{}) bool) uint64
	// 移除loop中的通道
	RemoveSelectCase(id uint64)

	// 创建一个回调Handler,用来向另一个selectloop表示发送者的对象，每一个消息只能使用自己的CallbackHandler
	NewCallbackHandler(returnMsg string, userData interface{}) *CallbackHandler

	RunIn(f func())
	SyncRunIn(f func() interface{}) interface{}
}
