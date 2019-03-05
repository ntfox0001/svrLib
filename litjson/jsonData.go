package litjson

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/ntfox0001/svrLib/log"
)

const (
	Type_None = iota
	Type_Bool
	Type_String
	Type_Double
	Type_List
	Type_Map
)

type JsonData struct {
	data      interface{}
	valueType int
}

func NewJsonData() *JsonData {
	return &JsonData{}
}

func NewJsonDataFromJson(json string) *JsonData {
	var obj interface{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(json, &obj); err != nil {
		log.Error("jsonData", "err", err.Error())
		return nil
	}
	jd := NewJsonDataFromObject(obj)

	return jd
}
func NewJsonDataFromObject(obj interface{}) *JsonData {
	jd := &JsonData{}

	jd.valueType = getType(obj)
	switch jd.valueType {
	case Type_Map:
		nmap := make(map[string]*JsonData)
		for k, v := range obj.(map[string]interface{}) {
			nmap[k] = NewJsonDataFromObject(v)
		}
		jd.data = nmap
	case Type_List:
		nlist := make([]*JsonData, 0)
		for _, v := range obj.([]interface{}) {
			nlist = append(nlist, NewJsonDataFromObject(v))
		}
		jd.data = nlist
	default:
		jd.data = obj
	}

	return jd
}

func getType(obj interface{}) int {
	if obj == nil {
		return Type_None
	}
	switch obj.(type) {
	case bool:
		return Type_Bool
	case int:
	case int64:
	case uint:
	case uint64:
	case float32:
	case float64:
		return Type_Double
	case string:
		return Type_String
	case []interface{}:
		return Type_List
	case map[string]interface{}:
		return Type_Map
	}
	return Type_None
}
func (jd *JsonData) ensure(valueType int) bool {
	// 如果还没有初始化，那么就用第一次调用的类型
	if jd.valueType == Type_None {
		jd.valueType = valueType
		switch valueType {
		case Type_List:
			jd.data = make([]*JsonData, 0)
		case Type_Map:
			jd.data = make(map[string]*JsonData)
		}
		return true
	}
	return jd.valueType == valueType
}

func (jd *JsonData) Get(key string) *JsonData {
	if jd.ensure(Type_Map) {
		if v, ok := jd.data.(map[string]*JsonData)[key]; ok {
			return v
		}
	}
	return nil
}

func (jd *JsonData) Safe_Get(key string) (*JsonData, error) {
	if jd.ensure(Type_Map) {
		if v, ok := jd.data.(map[string]*JsonData)[key]; ok {

			return v, nil
		}
		return nil, errors.New("No found item.")
	}
	return nil, errors.New("Type error.")
}

func (jd *JsonData) SetKey(key string, value interface{}) {
	if jd.ensure(Type_Map) {
		jd.data.(map[string]*JsonData)[key] = jd.isJsonData(value)
	}
}

func (jd *JsonData) RemoveKey(key string) {
	if jd.ensure(Type_Map) {
		delete(jd.data.(map[string]*JsonData), key)
	}
}

func (jd *JsonData) HasKey(key string) bool {
	if jd.ensure(Type_Map) {
		return jd.hasKey(key)
	}

	return false
}

func (jd *JsonData) hasKey(key string) bool {
	_, ok := jd.data.(map[string]*JsonData)[key]
	return ok
}

func (jd *JsonData) GetType() int {
	return jd.valueType
}

func (jd *JsonData) Index(id int) *JsonData {
	if jd.ensure(Type_List) {
		if jd.Len() >= id {
			return nil
		}

		return jd.data.([]*JsonData)[id]
	}
	return nil
}

func (jd *JsonData) Safe_Index(id int) (*JsonData, error) {
	if jd.ensure(Type_List) {
		if jd.Len() >= id {
			return nil, errors.New("id overflow.")
		}

		return jd.data.([]*JsonData)[id], nil
	}
	return nil, errors.New("Type error.")
}

func (jd *JsonData) SetIndex(id int, value interface{}) {
	if jd.ensure(Type_List) {
		if jd.Len() >= id {
			return
		}
		jd.data.([]*JsonData)[id] = jd.isJsonData(value)
	}
}

func (jd *JsonData) Append(value interface{}) {
	if jd.ensure(Type_List) {
		jd.data = append(jd.data.([]*JsonData), jd.isJsonData(value))
	}
}
func (jd *JsonData) RemoveId(id int) {
	if jd.ensure(Type_List) {
		jd.data = sliceDel(jd.data.([]*JsonData), id)
	}
}

func sliceDel(s []*JsonData, id int) []*JsonData {
	if s == nil {
		return s
	}
	if len(s) <= id {
		return s
	}

	t := append(s[:id], s[id+1:]...)
	return t
}

func (jd *JsonData) isJsonData(value interface{}) *JsonData {
	switch value.(type) {
	case *JsonData:
		return value.(*JsonData)
	default:
		return NewJsonDataFromObject(value)
	}
}

func (jd *JsonData) Len() int {
	switch jd.valueType {
	case Type_Map:
		return len(jd.data.(map[string]*JsonData))
	case Type_List:
		return len(jd.data.([]*JsonData))
	}
	return 0
}

func (jd *JsonData) GetString() string {
	if jd.ensure(Type_String) {
		return jd.data.(string)
	}
	return ""
}

func (jd *JsonData) GetFloat32() float32 {
	if jd.ensure(Type_Double) {
		return float32(jd.data.(float64))
	}
	return 0
}

func (jd *JsonData) GetFloat64() float64 {
	if jd.ensure(Type_Double) {
		return jd.data.(float64)
	}
	return 0
}

func (jd *JsonData) GetInt32() int {
	if jd.ensure(Type_Double) {
		return int(jd.data.(float64))
	}
	return 0
}

func (jd *JsonData) GetInt64() int64 {
	if jd.ensure(Type_Double) {
		return int64(jd.data.(float64))
	}
	return 0
}

func (jd *JsonData) GetUInt32() uint {
	if jd.ensure(Type_Double) {
		return uint(jd.data.(float64))
	}
	return 0
}

func (jd *JsonData) GetUInt64() uint64 {
	if jd.ensure(Type_Double) {
		return uint64(jd.data.(float64))
	}
	return 0
}

func (jd *JsonData) GetBool() bool {
	if jd.ensure(Type_Bool) {
		return jd.data.(bool)
	}
	return false
}

func (jd *JsonData) SetString(value string) {
	if jd.ensure(Type_String) {
		jd.data = value
	}
}

func (jd *JsonData) SetInt32(value int) {
	if jd.ensure(Type_Double) {
		jd.data = float64(value)
	}
}

func (jd *JsonData) SetInt64(value int64) {
	if jd.ensure(Type_Double) {
		jd.data = float64(value)
	}
}

func (jd *JsonData) SetUInt32(value uint) {
	if jd.ensure(Type_Double) {
		jd.data = float64(value)
	}
}

func (jd *JsonData) SetUInt64(value uint64) {
	if jd.ensure(Type_Double) {
		jd.data = float64(value)
	}
}

func (jd *JsonData) SetFloat32(value float32) {
	if jd.ensure(Type_Double) {
		jd.data = float64(value)
	}
}

func (jd *JsonData) SetFloat64(value float64) {
	if jd.ensure(Type_Double) {
		jd.data = value
	}
}
func (jd *JsonData) SetBool(value bool) {
	if jd.ensure(Type_Bool) {
		jd.data = value
	}
}
func (jd *JsonData) Map() map[string]*JsonData {
	if jd.ensure(Type_Map) {
		return jd.data.(map[string]*JsonData)
	}
	return nil
}
func (jd *JsonData) List() []*JsonData {
	if jd.ensure(Type_List) {
		return jd.data.([]*JsonData)
	}
	return nil
}
func (jd *JsonData) ToObject() interface{} {
	switch jd.valueType {
	case Type_Map:
		nmap := make(map[string]interface{})
		for k, v := range jd.data.(map[string]*JsonData) {
			nmap[k] = v.ToObject()
		}
		return nmap
	case Type_List:
		nlist := make([]interface{}, 0)
		for _, v := range jd.data.([]*JsonData) {
			nlist = append(nlist, v.ToObject())
		}
		return nlist
	default:
		return jd.data
	}
}

func (jd *JsonData) ToJson() string {
	if json, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(jd.ToObject()); err != nil {
		log.Error("JsonData", "error", err.Error())
		return ""
	} else {
		return json
	}

}
