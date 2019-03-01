package util

import (
	"github.com/ntfox0001/svrLib/litjson"
)

const MsgKeepDataName = "_MsgKeepDataName_"

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
	if from.HasKey(MsgKeepDataName) {
		jd.SetKey(MsgKeepDataName, from.Get(MsgKeepDataName))
	}

	return jd
}

func JdMsgSetKeep(jd *litjson.JsonData, key string, value interface{}) {
	if !jd.HasKey(MsgKeepDataName) {
		jd.SetKey(MsgKeepDataName, litjson.NewJsonData())
	}
	jd.Get(MsgKeepDataName).SetKey(key, value)
}

func JdMsgGetKeep(jd *litjson.JsonData, key string) *litjson.JsonData {
	keepjd := jd.Get(MsgKeepDataName)
	if keepjd != nil {
		return keepjd.Get(key)
	}

	return nil
}
