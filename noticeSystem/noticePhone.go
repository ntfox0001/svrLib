package noticeSystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/util"
	"strings"

	"github.com/ntfox0001/svrLib/log"

	jsoniter "github.com/json-iterator/go"
)

type phone struct {
	phoneTemplates []PhoneTemplateCfg
}

func newPhone(templates []PhoneTemplateCfg) *phone {
	return &phone{
		phoneTemplates: templates,
	}
}

func (p *phone) send(data map[string]string) (rt *PhoneMsgResp, rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = err.(error)
			rt = nil
			log.Error("phone send error", "Error", rtErr.Error())
			return
		}
	}()

	templ, err := p.getTemplateFromType(data["{type}"])
	if err != nil {
		return nil, err
	}

	ttsParam, err := p.getTtsParamString(templ, data)

	if err != nil {
		return nil, err
	}

	req := PhoneMsgReq{
		CalledNumber:   data["{phone}"],
		Id:             templ.Id,
		CallShowNumber: templ.CallShowNumber,
		TtsParam:       ttsParam,
	}

	if s, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(req); err != nil {
		log.Warn("phone", "send err", err.Error())
		return nil, err
	} else {
		postStr := "msg=" + string(s)
		if result, err := network.SyncHttpPost(NoticeTemplate.PhoneUrl, postStr, network.ContentTypeFrom); err != nil {
			log.Warn("phone", "post err", err.Error())
			return nil, err
		} else {
			// var resp PhoneMsgResp
			// if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(result), &resp); err != nil {
			// 	log.Error("invalid format of PhoneMsgResp", "resp", result)
			// 	return nil, err
			// }
			// 阿里云返回的是xml接口，简单处理一下
			if strings.Index(result, "OK") == -1 {
				return nil, commonError.NewStringErr(result)
			}

			return nil, nil
		}
	}

}

func (p *phone) getTemplateFromType(noticeType string) (*PhoneTemplateCfg, error) {
	for _, v := range p.phoneTemplates {
		if v.Type == noticeType {
			return &v, nil
		}
	}
	log.Error("phone", "TemplateType does not exist", noticeType)
	return nil, commonError.NewStringErr("TemplateType does not exist.")
}

func (p *phone) getTtsParamString(templ *PhoneTemplateCfg, data map[string]string) (string, error) {
	jsobj := make(map[string]string)
	for k, v := range templ.Args {
		s := util.StringReplace(v, data)
		jsobj[k] = s
	}

	if js, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(jsobj); err != nil {
		return "", err
	} else {
		return string(js), nil
	}
}
