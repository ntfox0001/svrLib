package jsonMsg

const (
	// jm默认超时，每个节点应该判断当前时间，如果超时，那么应该做相应处理（可选功能）
	JsonMsg_DefaultTimeout = 999
	// 超时时间
	JsonMsg_TimeoutName = "jmTimeout"
	// 消息创建时间
	JsonMsg_BuildTimeName = "jmBuildTime"
	// keep根节点名字
	JsonMsg_KeepRootName = "jmKeepRoot"
	// 消息压栈，当从一个jsonMsg生成新的jm时，旧消息压入新消息中，可以选择只压入keep数据
	JsonMsg_ParentName = "jmParent"
	// 消息id
	JsonMsg_MsgIdName = "msgId"
)
