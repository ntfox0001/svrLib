package util

import (
	"github.com/ntfox0001/svrLib/commonError"
)

// 获得path路径的值
func JsonGetValue(json interface{}, path []string) (interface{}, error) {
	j := json
	for _, p := range path {
		if !IsObj(j) {
			return nil, commonError.NewStringErr("json format error")
		}
		if v, ok := j.(map[string]interface{})[p]; !ok {
			return nil, commonError.NewStringErr("json format error")
		} else {
			j = v
		}
	}
	return j, nil
}

func IsObj(json interface{}) bool {
	_, ok := json.(map[string]interface{})
	return ok
}

func HasSKey(m interface{}, k string) bool {
	_, ok := m.(map[string]interface{})
	if ok {
		_, ok1 := m.(map[string]interface{})[k]
		if ok1 {
			return true
		}
	}

	return false
}
