package selectCase

import (
	"reflect"
	"sync/atomic"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/log"
)

type selectAddAndDelMsg struct {
	add bool

	// add
	front bool
	ch    reflect.Value
	cb    func(data interface{}) bool

	// add,del
	id uint64
}

type SelectLoop struct {
	name        string
	callbacks   *CallbackList
	selectCases *SelectCaseList

	actChan chan selectAddAndDelMsg

	unqiueId uint64 // 用于唯一id

	handler  selectCaseInterface.IEventChanHandler
	slHelper *SelectLoopHelper

	loopCount    uint64 // 用于保存循环次数
	preLoopCount uint64 // 上一次读取数据的次数
	// 用于退出loop
	quitChan chan interface{}
}

// actSize：注册通道的大小，防止在循环中调用注册函数导致卡死
// preRegSize：消息handler在循环没有开始时，可缓存的注册消息数量
func NewSelectLoop(name string, actSize int, preRegSize int) *SelectLoop {
	sl := &SelectLoop{
		name:         name,
		callbacks:    NewCallbackList(),
		selectCases:  NewSelectCaseList(),
		actChan:      make(chan selectAddAndDelMsg, actSize),
		handler:      NewEventChanHandler(name, actSize, preRegSize),
		loopCount:    0,
		preLoopCount: 0,
		quitChan:     make(chan interface{}, 1),
	}
	sl.slHelper = NewSelectLoopHelper(sl)

	sl.initSelectCase(reflect.ValueOf(sl.quitChan), sl.quit)
	sl.initSelectCase(reflect.ValueOf(sl.actChan), sl.processMsg)

	sl.handler.Initial(sl.slHelper)

	return sl
}
func NewSelectLoop2(name string, handler selectCaseInterface.IEventChanHandler, actSize int) *SelectLoop {
	sl := &SelectLoop{
		name:         name,
		callbacks:    NewCallbackList(),
		selectCases:  NewSelectCaseList(),
		actChan:      make(chan selectAddAndDelMsg, actSize),
		handler:      handler,
		loopCount:    0,
		preLoopCount: 0,
		quitChan:     make(chan interface{}, 1),
	}
	sl.slHelper = NewSelectLoopHelper(sl)

	sl.initSelectCase(reflect.ValueOf(sl.quitChan), sl.quit)
	sl.initSelectCase(reflect.ValueOf(sl.actChan), sl.processMsg)

	sl.handler.Initial(sl.slHelper)

	return sl
}
func (s *SelectLoop) GetEventChanHandler() selectCaseInterface.IEventChanHandler {
	return s.handler
}
func (s *SelectLoop) initSelectCase(ch reflect.Value, cb func(data interface{}) bool) {
	id := s.NextUnqiueId()
	se := s.selectCases.PushBack(reflect.SelectCase{Dir: reflect.SelectRecv, Chan: ch})
	se.UnqiueId = id
	ce := s.callbacks.PushBack(cb)
	ce.UnqiueId = id
}

// 运行循环
func (s *SelectLoop) Run() {
runable:
	for {
		chosen, recv, recvOk := reflect.Select(s.selectCases.ToSlice())
		if recvOk {
			if !s.callbacks.ToSlice()[chosen](recv.Interface()) {
				break runable
			}
			s.loopCount++
		}
	}
	log.Debug("- selectLoop end", "name", s.name)
}
func (h *SelectLoop) GetHelper() selectCaseInterface.ISelectLoopHelper {
	return h.slHelper
}

func (s *SelectLoop) NextUnqiueId() uint64 {
	return atomic.AddUint64(&s.unqiueId, 1)
}
func (s *SelectLoop) quit(interface{}) bool {
	return false
}
func (s *SelectLoop) processMsg(data interface{}) bool {
	msg := data.(selectAddAndDelMsg)
	if msg.add {
		if msg.front {
			se := s.selectCases.PushFront(reflect.SelectCase{Dir: reflect.SelectRecv, Chan: msg.ch})
			se.UnqiueId = msg.id
			ce := s.callbacks.PushFront(msg.cb)
			ce.UnqiueId = msg.id
		} else {
			se := s.selectCases.PushBack(reflect.SelectCase{Dir: reflect.SelectRecv, Chan: msg.ch})
			se.UnqiueId = msg.id
			ce := s.callbacks.PushBack(msg.cb)
			ce.UnqiueId = msg.id
		}
	} else {
		s.selectCases.RemoveForId(msg.id)
		s.callbacks.RemoveForId(msg.id)
	}

	return true
}

// ch: <-chan XXX类型的chan，必须用reflect.ValueOf转换过的chan，返回一个id，用于remove 线程安全
func (s *SelectLoop) AddSelectCaseFront(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	msg := selectAddAndDelMsg{
		add:   true,
		front: true,
		ch:    ch,
		cb:    cb,
		id:    s.NextUnqiueId(),
	}

	s.actChan <- msg

	return msg.id
}

// ch: <-chan XXX类型的chan，必须用reflect.ValueOf转换过的chan，返回一个id，用于remove 线程安全
func (s *SelectLoop) AddSelectCase(ch reflect.Value, cb func(data interface{}) bool) uint64 {
	msg := selectAddAndDelMsg{
		add:   true,
		front: false,
		ch:    ch,
		cb:    cb,
		id:    s.NextUnqiueId(),
	}

	s.actChan <- msg

	return msg.id
}

// 线程安全
func (s *SelectLoop) RemoveSelectCase(id uint64) {
	msg := selectAddAndDelMsg{
		add: false,

		id: id,
	}

	s.actChan <- msg
}
func (s *SelectLoop) Close() {
	s.quitChan <- struct{}{}
	close(s.actChan)
	s.handler.Close()
}

func (s *SelectLoop) Handler() selectCaseInterface.IEventChanHandler {
	return s.handler
}

func (s *SelectLoop) GetLoopCount() uint64 {
	return s.loopCount
}

// 读取LoopCount并重置
func (s *SelectLoop) PopLoopCount() uint64 {
	lc := s.loopCount - s.preLoopCount
	s.preLoopCount = s.loopCount
	return lc
}
