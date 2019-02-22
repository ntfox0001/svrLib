package selectCase

import (
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"reflect"
)

// 当线程可阻塞等待消息完成时，可以向另一个selectLoop发送SelectChannel，然后使用eventChan等待结果
type SelectChannel struct {
	eventChan chan interface{}
}

func NewSelectChannel() *SelectChannel {
	return &SelectChannel{
		eventChan: make(chan interface{}),
	}
}

// 消息返回之前，会阻塞
func (s *SelectChannel) GetReturn() interface{} {
	return <-s.eventChan
}
func (s *SelectChannel) GetReturnChan() chan interface{} {
	return s.eventChan
}
func (s *SelectChannel) SendMsgToMe(data selectCaseInterface.EventChanMsg) {
	s.eventChan <- data.Content
}

func (s *SelectChannel) RegisterEvent(event string, f func(interface{}) bool) uint64 {
	return 0
}
func (s *SelectChannel) UnregisterEvent(id uint64) {

}

func (s *SelectChannel) AddSelectCaseFront(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	return 0
}
func (s *SelectChannel) AddSelectCase(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	return 0
}

func (s *SelectChannel) RemoveSelectCase(id uint64) {

}
func (s *SelectChannel) NewCallbackHandler(returnMsg string, userData interface{}) *selectCaseInterface.CallbackHandler {
	return selectCaseInterface.NewCallbackHandler(returnMsg, s, userData)
}
