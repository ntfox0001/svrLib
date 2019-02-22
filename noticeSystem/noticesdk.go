package noticeSystem

import (
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/util"

	"github.com/ntfox0001/svrLib/log"
	jsoniter "github.com/json-iterator/go"
)

type smsSdk struct {
	SmsSdkCfg
}

func newSmsSdk(sdk SmsSdkCfg) *smsSdk {
	return &smsSdk{
		SmsSdkCfg: sdk,
	}
}

func (n *smsSdk) send(data map[string]string) (rt string, rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = err.(error)
			rt = ""
			log.Error("noticeSdk send error", "Error", rtErr.Error())
			return
		}
	}()

	sdkTemplate, err := n.getTemplateFromType(data["{type}"])
	if err != nil {
		return "", err
	}
	sdkmsg := noticeSdkMsg{
		AppId:      n.AppId,
		AppKey:     n.AppKey,
		SmsSign:    n.Sign,
		TemplateId: sdkTemplate.Id,
		NoticeMsg: NoticeMsg{
			Phone: []string{data["{phone}"]},
			Args:  make([]string, 0, 10),
		},
	}

	// 根据sdk模板，填写args
	for _, v := range sdkTemplate.Args {
		//sdkmsg.Args = append(sdkmsg.Args, msg[v])
		s := util.StringReplace(v, data)
		sdkmsg.Args = append(sdkmsg.Args, s)
	}

	if s, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(sdkmsg); err != nil {
		log.Warn("NoticeSdk", "send err", err.Error())
		return "", err
	} else {
		postStr := "msg=" + string(s)
		//log.Debug("sdkmsg req", "msg", postStr)
		if rt, err := network.SyncHttpPost(n.Url, postStr, network.ContentTypeFrom); err != nil {
			log.Warn("NoticeSdk", "post err", err.Error())
			return "", err
		} else {
			return rt, nil
		}
	}
}

func (n *smsSdk) getTemplateFromType(noticeType string) (*SmsMsgTemplate, error) {
	for _, v := range n.Templates {
		if v.Type == noticeType {
			return &v, nil
		}
	}
	log.Error("noticeSdk", "TemplateType does not exist", noticeType)
	return nil, commonError.NewStringErr("TemplateType does not exist.")
}
