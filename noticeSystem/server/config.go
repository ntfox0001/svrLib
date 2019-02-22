package main

import (
	"io/ioutil"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"
)

type WxMpCfg struct {
	AppId  string `json:"appId"`
	Secret string `json:"secret"`
}

type ConfigCfg struct {
	WxMp        []WxMpCfg `json:"wxMp"`
	RefreshTime int64     `json:"refreshTime"`
	ListenIp    string    `json:"listenIp"`
	Port        string    `json:"port"`
	WhiteIp     []string  `json:"whiteIp"`
}

var Config ConfigCfg

func InitApplicationConfig(configFilename string) error {
	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Error("config", "ReadFile", err.Error())
		return err
	}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, &Config); err != nil {
		return err
	}
	return nil
}
