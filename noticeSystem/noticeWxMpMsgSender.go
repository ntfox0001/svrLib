package noticeSystem

import (
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/util"
	"net/url"

	"github.com/ntfox0001/svrLib/log"

	jsoniter "github.com/json-iterator/go"
)

type noticeWxMpMsgSender struct {
	wxMpMsgTemplates []WxMpMsgTemplateCfg
}

func newNoticeWxMpMsgSender(templates []WxMpMsgTemplateCfg) *noticeWxMpMsgSender {
	return &noticeWxMpMsgSender{
		wxMpMsgTemplates: templates,
	}
}

func (w *noticeWxMpMsgSender) send(data map[string]string) (wxResp *WxMpMsgResp, rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = err.(error)
			wxResp = nil
			log.Error("noticeWxMpMsgSender send error", "Error", rtErr.Error())
			return
		}
	}()

	// 查找模板
	wxTemplate, err := w.getTemplateFromType(data["{type}"])
	if err != nil {
		return nil, err
	}

	wxurl, err := w.makeUrl(wxTemplate, data)
	if err != nil {
		return nil, err
	}

	wxmsg := WxMpMsg{
		Touser:      data["{openId}"],
		Template_Id: wxTemplate.Id,
		Url:         wxurl,
		Data:        make(map[string]WxMpMsgData),
	}

	for _, d := range wxTemplate.WxMpMsgTemplateData {
		v := util.StringReplace(d.Value, data)
		wxd := WxMpMsgData{
			Value: v,
			Color: d.Color,
		}
		wxmsg.Data[d.Name] = wxd
	}

	accessToken := data["{accessToken}"]
	content, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(wxmsg)
	//log.Debug("wxmp msg", "msg", string(content))
	wxResult, err := network.SyncHttpPost("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token="+accessToken, string(content), network.ContentTypeJson)
	if err != nil {
		return nil, err
	}
	var wxresp WxMpMsgResp
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(wxResult), &wxresp); err != nil {
		log.Error("invalid format of WX Resp", "resp", wxResult)
		return &wxresp, err
	}

	return &wxresp, nil
}

func (w *noticeWxMpMsgSender) getTemplateFromType(noticeType string) (*WxMpMsgTemplateCfg, error) {
	for _, v := range w.wxMpMsgTemplates {
		if v.Type == noticeType {
			return &v, nil
		}
	}
	log.Error("noticeWxMpMsgSender", "TemplateType does not exist", noticeType)
	return nil, commonError.NewStringErr("TemplateType does not exist.")
}

func (w *noticeWxMpMsgSender) makeUrl(template *WxMpMsgTemplateCfg, data map[string]string) (string, error) {
	if template.Url == "" {
		return "", nil
	}
	//if wxurl, err := url.Parse(template.Url); err != nil {
	//	return "", err
	//} else {
	v := url.Values{}
	for _, urlArgName := range template.UrlArgs {
		if urlArg, ok := data[urlArgName]; ok {
			v.Add(urlArgName, urlArg)
		} else {
			log.Warn("noticeWxMpMsgSender key does not found", "key", urlArgName)
			return "", commonError.NewStringErr("key does not found:" + urlArgName)
		}
	}
	//wxurl.RawQuery = v.Encode()
	return fmt.Sprint(template.Url, "?", v.Encode()), nil
	//}

}
