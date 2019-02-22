package config

import (
	"io/ioutil"
	"github.com/ntfox0001/svrLib/commonError"
	"strconv"

	jsoniter "github.com/json-iterator/go"

	"github.com/ntfox0001/svrLib/log"
)

type Config struct {
	jsonObj map[string]interface{}
}

func NewConfig(filename string) (*Config, error) {
	cf := Config{
		jsonObj: make(map[string]interface{}),
	}
	var err error
	if cf.jsonObj, err = cf.readFile(filename); err != nil {
		return nil, err
	}
	return &cf, nil
}

func NewConfigByString(json string) (*Config, error) {
	cf := Config{
		jsonObj: make(map[string]interface{}),
	}
	var err error
	if cf.jsonObj, err = cf.readString(json); err != nil {
		return nil, err
	}
	return &cf, nil
}
func NewConfigByInterface(obj interface{}) (*Config, error) {
	if js, ok := obj.(map[string]interface{}); ok {
		cf := Config{
			jsonObj: js,
		}
		return &cf, nil
	}
	return nil, commonError.NewStringErr("invalid json object.")
}

func (c *Config) readFile(filename string) (j map[string]interface{}, e error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("config", "ReadFile: ", err.Error())
		return nil, err
	} else {
		j, e = c.readbytes(bytes)
	}

	return j, nil
}
func (c *Config) readbytes(bytes []byte) (j map[string]interface{}, e error) {
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, &j); err != nil {
		log.Error("config", "Unmarshal: ", err.Error())
		return nil, err
	} else {
		return j, err
	}
}
func (c *Config) readString(json string) (j map[string]interface{}, e error) {
	j, e = c.readbytes([]byte(json))
	return
}

func (c *Config) GetStringValue(key string, defVal string) string {
	if v, ok := c.jsonObj[key]; ok {
		return v.(string)
	}
	return defVal
}

func (c *Config) GetValue(key string, defVal interface{}) interface{} {
	if v, ok := c.jsonObj[key]; ok {
		//fmt.Println(reflect.TypeOf(v))
		return v
	}
	return defVal
}

func (c *Config) GetIntValue(key string, defVal int) int {
	if v, ok := c.jsonObj[key]; ok {
		switch v.(type) {
		case int:
			return v.(int)
		case float64:
			return int(v.(float64))
		case string:
			if b, err := strconv.Atoi(v.(string)); err == nil {
				return b
			} else {
				return defVal
			}
		default:
			return defVal
		}
	}
	return defVal
}

func (c *Config) GetGroup(key string) (*Config, error) {
	v := c.GetValue(key, struct{}{})
	if group, ok := v.(map[string]interface{}); ok {
		return &Config{
			jsonObj: group,
		}, nil
	}
	return nil, commonError.NewStringErr("group " + key + " does not exist.")
}
