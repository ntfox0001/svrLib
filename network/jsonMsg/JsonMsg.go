package jsonMsg

import (
	"time"

	"github.com/ntfox0001/svrLib/litjson"
)

type JsonMsg struct {
	Data     *litjson.JsonData
	keepRoot *litjson.JsonData
	stack    *litjson.JsonData
}

// 创建jm
func NewJsonMsg(name string) *JsonMsg {
	return NewJsonMsgWithTimeout(name, JsonMsg_DefaultTimeout)
}

// 使用外部json对象创建jm
func NewJsonMsgWithData(data *litjson.JsonData) *JsonMsg {
	jm := &JsonMsg{
		Data:     data,
		keepRoot: nil,
		stack:    nil,
	}

	if keepData := data.Get(JsonMsg_KeepRootName); keepData != nil {
		jm.SetKeepRoot(keepData)
	} else {
		jm.SetKeepRoot(litjson.NewJsonDataByType(litjson.Type_Map))
	}

	return jm
}

// 创建jm，带有超时信息
func NewJsonMsgWithTimeout(name string, timeout int) *JsonMsg {
	jm := _newEmptyJsonMsg(name)

	jm.SetKeepRoot(litjson.NewJsonDataByType(litjson.Type_Map))
	jm.SetTimeout(timeout)
	return jm
}

// 创建一个空jm，没有keep和stack，外部程序不应该调用
func _newEmptyJsonMsg(name string) *JsonMsg {
	jm := &JsonMsg{
		Data:     litjson.NewJsonDataByType(litjson.Type_Map),
		keepRoot: nil,
		stack:    nil,
	}
	jm.Data.SetKey(JsonMsg_MsgIdName, name)
	return jm
}

func (jm *JsonMsg) SetStack(data *litjson.JsonData) {
	jm.stack = litjson.NewJsonDataByType(litjson.Type_Map)
	jm.Data.SetKey(JsonMsg_StackName, jm.stack)
}

func (jm *JsonMsg) SetKeepRoot(data *litjson.JsonData) {
	jm.Data.SetKey(JsonMsg_KeepRootName, data)
	jm.keepRoot = data
}

// 设置超时时间
func (jm *JsonMsg) SetTimeout(timeout int) {
	jm.SetKeepData(JsonMsg_TimeoutName, timeout)
	jm.SetKeepData(JsonMsg_BuildTimeName, time.Now().Unix())
}

func (jm *JsonMsg) GetTimeout() int {
	return jm.Data.Get(JsonMsg_TimeoutName).GetInt()
}
func (jm *JsonMsg) GetBuildTime() int64 {
	return jm.Data.Get(JsonMsg_BuildTimeName).GetInt64()
}

// 设置一个keep数据
func (jm *JsonMsg) SetKeepData(name string, value interface{}) {
	jm.keepRoot.SetKey(name, value)
}

// 获取一个keep数据
func (jm *JsonMsg) GetKeepData(name string) *litjson.JsonData {
	return jm.keepRoot.Get(name)
}

// 创建新的msg，并且迁移keep和stack数据
func (jm *JsonMsg) NewJsonMsg(name string) *JsonMsg {
	newjm := _newEmptyJsonMsg(name)

	newjm.SetTimeout(jm.GetTimeout())
	newjm.SetStack(jm.stack)
	newjm.SetKeepRoot(jm.keepRoot)
	return newjm
}

// 创建新的msg，老数据压入stack, full指定是否将所有数据都压入stack,false将只压入keep和stack数据
func (jm *JsonMsg) NewSubJsonMsg(name string, full bool) *JsonMsg {
	newjm := NewJsonMsgWithTimeout(name, jm.GetTimeout())

	if full {
		newjm.SetStack(jm.Data)
	} else {
		data := litjson.NewJsonDataByType(litjson.Type_Map)
		// 复制所有特殊字段数据
		data.SetKey(JsonMsg_TimeoutName, jm.GetTimeout())
		data.SetKey(JsonMsg_BuildTimeName, jm.GetBuildTime())
		data.SetKey(JsonMsg_KeepRootName, jm.keepRoot)
		data.SetKey(JsonMsg_StackName, jm.stack)
		newjm.SetStack(data)
	}
	return newjm
}

// 使用上一个stack的数据创建jm
func (jm *JsonMsg) NewParentJsonMsg(name string) *JsonMsg {
	newjm := NewJsonMsgWithData(jm.stack)
	newjm.Data.SetKey(JsonMsg_MsgIdName, name)

	newjm.keepRoot = newjm.Data.Get(JsonMsg_KeepRootName)
	newjm.stack = newjm.Data.Get(JsonMsg_StackName)

	return newjm
}
