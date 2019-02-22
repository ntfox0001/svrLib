package selectCase

import (
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"reflect"
)

type SelectLoopHelper struct {
	selectLoop *SelectLoop
}

func (s *SelectLoopHelper) SendMsgToMe(data selectCaseInterface.EventChanMsg) {
	s.selectLoop.handler.Touch(data)
}
func (s *SelectLoopHelper) RegisterEvent(event string, f func(interface{}) bool) uint64 {
	return s.selectLoop.handler.RegisterEvent(event, f)
}
func (s *SelectLoopHelper) UnregisterEvent(id uint64) {
	s.selectLoop.handler.UnregisterEvent(id)
}

func (s *SelectLoopHelper) AddSelectCaseFront(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	return s.selectLoop.AddSelectCaseFront(ch, cb)
}
func (s *SelectLoopHelper) AddSelectCase(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	return s.selectLoop.AddSelectCase(ch, cb)
}

func (s *SelectLoopHelper) RemoveSelectCase(id uint64) {
	s.selectLoop.RemoveSelectCase(id)
}

func (s *SelectLoopHelper) NewCallbackHandler(returnMsg string, userData interface{}) *selectCaseInterface.CallbackHandler {
	return selectCaseInterface.NewCallbackHandler(returnMsg, s, userData)
}
