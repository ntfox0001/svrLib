package selectCase

import (
	"reflect"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
)

const (
	RunInReqEventName = "__RunInReqEvent"
)

type SelectLoopHelper struct {
	selectLoop *SelectLoop
}

func NewSelectLoopHelper(s *SelectLoop) *SelectLoopHelper {
	helper := &SelectLoopHelper{
		selectLoop: s,
	}
	helper.RegisterEvent(RunInReqEventName, helper.onRunInReq)
	return helper
}
func (s *SelectLoopHelper) SendMsgToMe(data selectCaseInterface.EventChanMsg) {
	s.selectLoop.handler.Touch(data)
}
func (s *SelectLoopHelper) RegisterEvent(event string, f func(selectCaseInterface.EventChanMsg)) uint64 {
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

// 在协程中，运行指定的函数
func (s *SelectLoopHelper) RunIn(f func()) {
	s.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(RunInReqEventName, nil, f))
}

// 在 logicsystem 里运行指定函数的协程中，运行制定的函数，阻塞返回
func (s *SelectLoopHelper) SyncRunIn(f func() interface{}) interface{} {
	rtchan := make(chan interface{})
	rtfunc := func() {
		rtchan <- f()
	}
	s.selectLoop.GetHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(RunInReqEventName, nil, rtfunc))
	return <-rtchan
}

// 在协程里运行指定函数
func (s *SelectLoopHelper) onRunInReq(msg selectCaseInterface.EventChanMsg) {
	f := msg.Content.(func())

	f()
}
