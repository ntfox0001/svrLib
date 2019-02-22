package util

import (
	jsoniter "github.com/json-iterator/go"
)

// 将一个map[string]interface{}结构类型转换成struct, struPtr目标对象的指针
func I2Stru(i interface{}, struPtr interface{}) error {
	if js, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(i); err == nil {
		//fmt.Println(string(js))
		if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(js, struPtr); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func ToJson(i interface{}) ([]byte, error) {
	return jsoniter.Marshal(i)
}
