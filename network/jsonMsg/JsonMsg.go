package jsonMsg

import (
	"time"

	"github.com/ntfox0001/svrLib/litjson"
)

type JsonMsg struct {
	Data     *litjson.JsonData
	keepRoot *litjson.JsonData
	parent   *litjson.JsonData
}

// 创建jm
func NewJsonMsg(name string) *JsonMsg {
	jm := &JsonMsg{
		Data:     litjson.NewJsonDataByType(litjson.Type_Map),
		keepRoot: nil,
		parent:   nil,
	}
	jm.Data.SetKey(JsonMsg_MsgIdName, name)
	jm.SetKeepRoot(litjson.NewJsonDataByType(litjson.Type_Map))
	jm.SetTimeout(JsonMsg_DefaultTimeout)
	return jm
}

// 使用外部json对象创建jm
func NewJsonMsgWithData(data *litjson.JsonData) *JsonMsg {
	jm := &JsonMsg{
		Data:     data,
		keepRoot: nil,
		parent:   nil,
	}

	if keepData := data.Get(JsonMsg_KeepRootName); keepData != nil {
		jm.keepRoot = keepData
	} else {
		jm.SetKeepRoot(litjson.NewJsonDataByType(litjson.Type_Map))
	}

	jm.SetTimeout(JsonMsg_DefaultTimeout)
	return jm
}

func (jm *JsonMsg) SetParent(parent *litjson.JsonData) {
	jm.parent = parent
	jm.Data.SetKey(JsonMsg_ParentName, parent)
}

func (jm *JsonMsg) SetKeepRoot(kr *litjson.JsonData) {
	jm.Data.SetKey(JsonMsg_KeepRootName, kr)
	jm.keepRoot = kr
}

// 设置超时时间
func (jm *JsonMsg) SetTimeout(timeout int64) {
	jm.SetKeepData(JsonMsg_TimeoutName, timeout)
	jm.SetKeepData(JsonMsg_BuildTimeName, time.Now().Unix())
}
func (jm *JsonMsg) SetTimeoutWhitBuildTime(timeout, buildTime int64) {
	jm.SetKeepData(JsonMsg_TimeoutName, timeout)
	jm.SetKeepData(JsonMsg_BuildTimeName, buildTime)
}

func (jm *JsonMsg) GetTimeout() int64 {
	return jm.keepRoot.Get(JsonMsg_TimeoutName).GetInt64()
}
func (jm *JsonMsg) GetBuildTime() int64 {
	return jm.keepRoot.Get(JsonMsg_BuildTimeName).GetInt64()
}

// 设置一个keep数据
func (jm *JsonMsg) SetKeepData(name string, value interface{}) {
	if jm.keepRoot == nil {
		jm.keepRoot = litjson.NewJsonDataByType(litjson.Type_Map)
		jm.Data.SetKey(JsonMsg_KeepRootName, jm.keepRoot)
	}
	jm.keepRoot.SetKey(name, value)
}

// 获取一个keep数据
func (jm *JsonMsg) GetKeepData(name string) *litjson.JsonData {
	return jm.keepRoot.Get(name)
}
func (jm *JsonMsg) ClearAllExtentData() {

	jm.Data.RemoveKey(JsonMsg_KeepRootName)
	jm.Data.RemoveKey(JsonMsg_ParentName)

	jm.keepRoot = nil
	jm.parent = nil
}

// 创建新的msg，老数据压入stack, full指定是否将所有数据都压入stack,false将只压入keep和stack数据
func (jm *JsonMsg) NewSubJsonMsg(name string) *JsonMsg {
	newjm := NewJsonMsg(name)

	newjm.SetParent(jm.Data)
	newjm.SetTimeoutWhitBuildTime(jm.GetTimeout(), jm.GetBuildTime())
	return newjm
}

func (jm *JsonMsg) GetParentJsonMsg() *JsonMsg {
	return NewJsonMsgWithData(jm.parent)
}
