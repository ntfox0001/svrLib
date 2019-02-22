package networkInterface

type IMsgHandler interface {
	// 发送数据
	SendJsonMsg(msg interface{}) error
	SendMsg(msg IMsgData) error

	// 注册json消息
	RegisterJsonMsg(msgId string, handler func(map[string]interface{})) error
	// 设置json消息处理函数，替代默认函数
	SetDispatchJsonMsgHandler(f func(msg map[string]interface{}))
	DispatchJsonMsg(msg map[string]interface{}) error

	// 注册消息
	RegisterMsg(msgId string, handler func(*RawMsgData)) error
	// 设置消息处理函数，替代默认函数
	SetDispatchMsgHandler(f func(*RawMsgData))
	DispatchMsg(msg *RawMsgData) error

	// 断开连接
	Disconnect()

	// 处理接收到的消息，需要每帧调用
	ProcessMsg() error
}
