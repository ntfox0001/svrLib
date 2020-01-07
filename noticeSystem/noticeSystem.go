package noticeSystem

import (
	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/noticeSystem/wxAccessRefMsg"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/util"
	"sync"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var _self *NoticeSystem

type NoticeSystem struct {

	// 短信用
	sdks    []*smsSdk
	lastSdk int32
	goPool  *util.GoroutineFixedPool

	// 微信用
	accessTokenUrl    string // 用于获取微信刷新token
	wxMpMsgSender     *noticeWxMpMsgSender
	wxAccessToken     string
	wxAccessTokenLock sync.RWMutex
	quitChan          chan interface{}

	// phone用
	phoneSender *phone

	// robot
	wxRobotSender *wxRobot
}

type NoticeSystemParams struct {
	GoPoolSize     int    `json:"goPoolSize"`
	ExecSize       int    `json:"execSize"`
	AccessTokenUrl string `json:"accessTokenUrl"`
}

type wxMpRefreshAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func Instance() *NoticeSystem {
	if _self == nil {
		_self = &NoticeSystem{
			sdks:     make([]*smsSdk, 0, 10),
			lastSdk:  0,
			goPool:   nil,
			quitChan: make(chan interface{}, 1),
		}
	}
	return _self
}

func (*NoticeSystem) Initial(params NoticeSystemParams) {
	_self.goPool = util.NewGoFixedPool("NoticeSystem sdk", params.GoPoolSize, params.ExecSize)
	_self.accessTokenUrl = params.AccessTokenUrl
	_self.refreshAccessToken()

	go _self.run()
}

func (*NoticeSystem) run() {
	// 每10分钟，取一次accesstoken
	ticker := time.NewTicker(time.Second * 600)
runable:
	for {
		select {
		case <-_self.quitChan:
			break runable
		case <-ticker.C:
			_self.refreshAccessToken()
		}

	}
}
func (*NoticeSystem) refreshAccessToken() {
	req := wxAccessRefMsg.WxAccessTokenReq{AppId: NoticeTemplate.WxAppId}
	if reqJs, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(req); err != nil {
		log.Error("NoticeSystem req marshal error", "err", err.Error())
	} else {
		if result, err := network.SyncHttpPost(_self.accessTokenUrl, string(reqJs), network.ContentTypeJson); err == nil {
			var resp wxAccessRefMsg.WxAccessTokeyResp
			if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(result), &resp); err == nil {

				//log.Debug("WxMp access refresh success", "old", _self.wxAccessToken, "new", resp.Token)
				log.Debug("WxMp access refresh success")
				_self.wxAccessTokenLock.Lock()
				_self.wxAccessToken = resp.Token
				_self.wxAccessTokenLock.Unlock()
			} else {
				log.Error("Failed to RefreshAccessToken", "result", result, "err", err.Error())
			}
		}
	}
}
func (*NoticeSystem) Release() {
	_self.wxRobotSender.close()
	_self.quitChan <- struct{}{}
	_self.goPool.Release()

	log.Debug("NoticeSystem release")
}

// 添加sms短信sdk
func (*NoticeSystem) AddMsgSdk(sdk SmsSdkCfg) {
	item := newSmsSdk(sdk)

	_self.sdks = append(_self.sdks, item)
}

// 添加微信模板
func (*NoticeSystem) AddWxMpMsgTemplate(templates []WxMpMsgTemplateCfg) {
	_self.wxMpMsgSender = newNoticeWxMpMsgSender(templates)
}

// 添加phone模板
func (*NoticeSystem) AddPhoneTemplate(templates []PhoneTemplateCfg) {
	_self.phoneSender = newPhone(templates)
}

func (*NoticeSystem) AddWxRobotTemplate(ip, port string, refreshTime uint, templates []WxRobotTemplateCfg) error {
	var err error
	_self.wxRobotSender, err = newWxRobot(ip, port, refreshTime, templates)
	if err != nil {
		log.Error("AddWxRobotTemplate error", "err", err.Error())
		return err
	}
	return nil
}

// 发送短信通知
func (*NoticeSystem) SendSmsNotice(cbHandler *selectCaseInterface.CallbackHandler, data map[string]string) {
	if len(_self.sdks) == 0 {
		return
	}
	// 获得下一个sdk
	newid := int(atomic.AddInt32(&_self.lastSdk, 1))

	sdk := _self.sdks[newid%len(_self.sdks)]

	_self.goPool.Go(func() {

		resp := NoticeResp{}
		resp.ReturnString, resp.Err = sdk.send(data)
		//log.Debug("SendSmsNotice resp", "resp", resp)

		if cbHandler != nil {
			cbHandler.SendReturnMsgNoReturn(resp)
		}

	}, nil)

	return
}
func (*NoticeSystem) getAccessToken() (rt string) {
	_self.wxAccessTokenLock.RLock()
	rt = _self.wxAccessToken
	_self.wxAccessTokenLock.RUnlock()
	return
}

// 发送微信通知
func (*NoticeSystem) SendWxMpNotice(cbHandler *selectCaseInterface.CallbackHandler, data map[string]string) {
	if _self.wxMpMsgSender == nil {
		return
	}
	data["{accessToken}"] = _self.getAccessToken()

	_self.goPool.Go(func() {

		resp, err := _self.wxMpMsgSender.send(data)
		if err != nil {
			log.Error("SendWxMpNotice", "err", err.Error())
			return
		}
		if resp.ErrMsg != "ok" {
			log.Error("SendWxMpNotice resp err", "resp", resp.ErrMsg)
		}

		if cbHandler != nil {
			cbHandler.SendReturnMsgNoReturn(resp)
		}

	}, nil)
}

func (*NoticeSystem) SendPhoneNotice(cbHandler *selectCaseInterface.CallbackHandler, data map[string]string) {
	if _self.phoneSender == nil {
		return
	}
	_self.goPool.Go(func() {
		resp, err := _self.phoneSender.send(data)
		if err != nil {
			log.Error("SendPhoneNotice", "err", err.Error())
		}

		if cbHandler != nil {
			cbHandler.SendReturnMsgNoReturn(resp)
		}
	}, nil)
}

func (*NoticeSystem) SendWxRobotNotice(data map[string]string) {
	if _self.wxRobotSender == nil {
		return
	}
	_self.goPool.Go(func() {
		err := _self.wxRobotSender.roomSend(data)
		if err != nil {
			log.Error("SendPhoneNotice", "err", err.Error())
		}
	}, nil)
}

func (*NoticeSystem) GetWxRobotSender() *wxRobot {
	return _self.wxRobotSender
}
