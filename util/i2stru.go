package util

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/ntfox0001/svrLib/log"
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

func Json2Stru(js string, struPtr interface{}) error {
	err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(js, struPtr)
	return err
}

func ToJson(i interface{}) ([]byte, error) {
	return jsoniter.Marshal(i)
}

func LoadConfigFile(filename string, confgPtr interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("config", "ReadFile", err.Error())
		return err
	}

	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, confgPtr); err != nil {
		log.Error("config", "Unmarshal: ", err.Error())
		return err
	}

	if err != nil {
		log.Error("Read Json Error")
	}

	return err
}
