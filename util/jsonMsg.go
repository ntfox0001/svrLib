package util

// 创建一个json消息
func JsMsg(msgId string) map[string]interface{} {
	v := make(map[string]interface{})
	v["msgId"] = msgId
	v["errorId"] = "0"
	return v
}
