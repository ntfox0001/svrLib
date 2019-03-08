package selectCaseInterface

type IEventChanHandler interface {
	// 初始化
	Initial(helper ISelectLoopHelper) error
	// 处理注册消息
	ProcessRegMsg(data interface{}) bool
	// 分派事件
	DispatchEvent(data interface{}) (rt bool)
	// 触发事件（向目标对象发送消息），线程安全
	Touch(msg EventChanMsg)

	// 注册事件，返回唯一id用于反注册
	RegisterEvent(event string, f func(EventChanMsg)) uint64
	UnregisterEvent(id uint64)

	Close()
}
