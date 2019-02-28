package util

import (
	"github.com/ntfox0001/svrLib/litjson"
)

// 创建一个json消息
func JsMsg(msgId string) map[string]interface{} {
	v := make(map[string]interface{})
	v["msgId"] = msgId
	v["errorId"] = "0"
	return v
}

func JdMsg(msgId string) *litjson.JsonData {
	jd := litjson.NewJsonData()
	jd.SetKey("msgId", msgId)
	jd.SetKey("errorId", "0")
	return jd
}
