package noticeSystem

import (
	"io/ioutil"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"
)

// 短信
type SmsMsgTemplate struct {
	Type string   `json:"type"`
	Id   string   `json:"id"`
	Args []string `json:"args"`
}
type SmsSdkCfg struct {
	AppId     string           `json:"appId"`
	AppKey    string           `json:"appKey"`
	Sign      string           `json:"sign"`
	Url       string           `json:"url"`
	Templates []SmsMsgTemplate `json:"templates"`
}

// 微信
type WxMpMsgTemplateDataCfg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Color string `json:"color"`
}
type WxMpMsgTemplateCfg struct {
	Id                  string                   `json:"id"`
	Type                string                   `json:"type"`
	Url                 string                   `json:"url"`
	UrlArgs             []string                 `json:"urlArgs"`
	WxMpMsgTemplateData []WxMpMsgTemplateDataCfg `json:"wxMpMsgTemplateData"`
}

// phone
type PhoneTemplateCfg struct {
	Id             string            `json:"id"`
	Type           string            `json:"type"`
	CallShowNumber string            `json:"callNum"`
	Args           map[string]string `json:"args"`
}

type WxRobotTemplateCfg struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

type NoticeTemplateCfg struct {
	MsgSdk           []SmsSdkCfg          `json:"msgSdk"`
	WxMpMsgTemplates []WxMpMsgTemplateCfg `json:"wxMpMsgTemplates"`
	WxAppId          string               `json:"wxAppId"`
	PhoneUrl         string               `json:"phoneUrl"`
	PhoneTemplates   []PhoneTemplateCfg   `json:"phoneTemplates"`
	WxRobotTemplates []WxRobotTemplateCfg `json:"wxRobotTemplates"`
}

var NoticeTemplate NoticeTemplateCfg

func InitNoticeTemplate(configFilename string) error {
	bytes, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Error("NoticeTeplate", "ReadFile", err.Error())
		return err
	}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(bytes, &NoticeTemplate); err != nil {
		return err
	}
	return nil
}
