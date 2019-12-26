package selectCase

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/ntfox0001/svrLib/debug"
	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
)

type EventChanHandler struct {
	eventChan    chan selectCaseInterface.EventChanMsg
	regChan      chan selectCaseInterface.EventRegMsg
	eventMap     map[string]map[uint64]func(selectCaseInterface.EventChanMsg)
	eventIdMap   map[uint64]string
	count        uint64
	name         string
	dispatchFunc func(data interface{}) bool
}

// 在初始化时，还没开始消息循环，此时注册将暂时储存在chan中
// preRegSize是设置一个足够的大小
// eventChanSize在消息处理函数中给自己发送消息的数量
func NewEventChanHandler(name string, eventChanSize, preRegSize int) *EventChanHandler {
	return &EventChanHandler{
		eventChan:    make(chan selectCaseInterface.EventChanMsg, eventChanSize),
		regChan:      make(chan selectCaseInterface.EventRegMsg, preRegSize),
		eventMap:     make(map[string]map[uint64]func(selectCaseInterface.EventChanMsg)),
		eventIdMap:   make(map[uint64]string),
		count:        1,
		name:         name,
		dispatchFunc: nil,
	}
}
func (h *EventChanHandler) Close() {
	//强制关闭，导致对面panic
	close(h.eventChan)
	h.eventMap = nil
}
func (h *EventChanHandler) Initial(helper selectCaseInterface.ISelectLoopHelper) error {
	helper.AddSelectCase(reflect.ValueOf(h.regChan), h.ProcessRegMsg)
	helper.AddSelectCase(reflect.ValueOf(h.eventChan), h.DispatchEvent)

	return nil
}

// 触发事件（向目标对象发送消息），线程安全
func (h *EventChanHandler) Touch(msg selectCaseInterface.EventChanMsg) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("EventChanHandler", "name", h.name, "err", err.(error).Error())

			return
		}
	}()
	select {
	case h.eventChan <- msg:
	default:
		log.Error("EventChanHandlerTouch", "err", "ChannelFull", "msg", msg.MsgId)
	}

}

// 处理注册消息
func (h *EventChanHandler) ProcessRegMsg(data interface{}) bool {
	msg := data.(selectCaseInterface.EventRegMsg)
	if msg.Reg {
		h.registerEvent(msg)
	} else {
		h.unregisterEvent(msg.EventFuncId)
	}
	return true
}

// 设置一个分派消息函数，如果返回真，跳过后面的处理逻辑
func (h *EventChanHandler) SetDispatchEvent(f func(data interface{}) bool) {
	h.dispatchFunc = f
}

// 分派事件
func (h *EventChanHandler) DispatchEvent(data interface{}) (rt bool) {
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("\n%s\n", string(debug.RuntimeStacks()))
			log.Error("event", "name", h.name, "err", err.(error).Error(), "\nstack", es)
			rt = true
			return
		}
	}()

	if h.dispatchFunc != nil {
		if h.dispatchFunc(data) {
			// 如果返回假，那么继续处理消息
			return true
		}
	}

	msg := data.(selectCaseInterface.EventChanMsg)
	//log.Error("DispatchEvent", "msg", msg.MsgId)
	if fs, ok := h.eventMap[msg.MsgId]; ok {
		for _, v := range fs {
			h.execEvent(v, msg)
		}
	} else {
		log.Warn("EventHandler, Unknow msgId", "name", h.name, "msgId", msg.MsgId)
	}
	return true
}

func (h *EventChanHandler) execEvent(f func(data selectCaseInterface.EventChanMsg), msg selectCaseInterface.EventChanMsg) {
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("\n%s\n", string(debug.RuntimeStacks()))
			log.Error("execEvent error", "err", err.(error).Error(), "msg", msg, "\nstack", es)
			return
		}
	}()

	f(msg)
}

func (h *EventChanHandler) GetNextId() uint64 {
	// 返回一个唯一值
	return atomic.AddUint64(&h.count, 1)
}

func (h *EventChanHandler) RegisterEvent(event string, f func(selectCaseInterface.EventChanMsg)) uint64 {
	msg := selectCaseInterface.EventRegMsg{
		Reg:         true,
		EventId:     event,
		EventFuncId: h.GetNextId(),
		F:           f,
	}
	h.regChan <- msg

	return msg.EventFuncId
}

func (h *EventChanHandler) UnregisterEvent(id uint64) {
	msg := selectCaseInterface.EventRegMsg{
		Reg:         false,
		EventId:     "",
		EventFuncId: id,
		F:           nil,
	}
	h.regChan <- msg

}

func (h *EventChanHandler) registerEvent(msg selectCaseInterface.EventRegMsg) {
	if fs, ok := h.eventMap[msg.EventId]; ok {
		fs[msg.EventFuncId] = msg.F
	} else {
		fs := make(map[uint64]func(selectCaseInterface.EventChanMsg))
		fs[msg.EventFuncId] = msg.F
		h.eventMap[msg.EventId] = fs
	}
	h.eventIdMap[msg.EventFuncId] = msg.EventId
	return
}

func (h *EventChanHandler) unregisterEvent(id uint64) {
	if v, ok := h.eventIdMap[id]; ok {
		delete(h.eventMap[v], id)
		delete(h.eventIdMap, id)
	}

}
