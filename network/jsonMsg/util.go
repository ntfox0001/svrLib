package jsonMsg

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

func JdMsgFrom(msgId string, from *litjson.JsonData) *litjson.JsonData {
	jd := JdMsg(msgId)
	if from.HasKey(JsonMsg_KeepRootName) {
		jd.SetKey(JsonMsg_KeepRootName, from.Get(JsonMsg_KeepRootName))
	}

	return jd
}

func JdMsgSetKeep(jd *litjson.JsonData, key string, value interface{}) {
	if !jd.HasKey(JsonMsg_KeepRootName) {
		jd.SetKey(JsonMsg_KeepRootName, litjson.NewJsonData())
	}
	jd.Get(JsonMsg_KeepRootName).SetKey(key, value)
}

func JdMsgGetKeep(jd *litjson.JsonData, key string) *litjson.JsonData {
	keepjd := jd.Get(JsonMsg_KeepRootName)
	if keepjd != nil {
		return keepjd.Get(key)
	}

	return nil
}

func JdMsgGetKeepRoot(jd *litjson.JsonData) *litjson.JsonData {
	return jd.Get(JsonMsg_KeepRootName)
}

// 只是迁移Keep data数据
func JdMsgMigrate(from *litjson.JsonData, to *litjson.JsonData) {
	to.SetKey(JsonMsg_KeepRootName, JdMsgGetKeepRoot(from))
}
