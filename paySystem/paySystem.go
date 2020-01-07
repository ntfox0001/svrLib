package paySystem

import (
	"fmt"

	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
	"github.com/ntfox0001/svrLib/util"

	"github.com/ntfox0001/svrLib/log"
)

type PaySystem struct {
	wxNotifyUrl string
	server      *network.Server
	wxMpPayMap  map[string]*WxMpPay
	goPool      *util.GoroutinePool
	callback    *selectCaseInterface.CallbackHandler // 支付成功时，调用
}

const (
	WxNotifyPath = "wxPayNotify"
)

var _self *PaySystem

func Instance() *PaySystem {
	if _self == nil {
		_self = &PaySystem{
			wxMpPayMap: make(map[string]*WxMpPay),
		}
	}
	return _self
}

// 支付系统初始化，goPoolSize在支付量很大时，要适当提高
func (*PaySystem) Initial(listenip, port string, goPoolSize, execSize int) error {

	_self.server = network.NewServer(listenip, port)
	_self.goPool = util.NewGoPool("PaySystem", goPoolSize, execSize)

	// 接受微信支付服务器通知的地址
	_self.wxNotifyUrl = fmt.Sprintf("http://%s:%s/%s", listenip, port, WxNotifyPath)
	return nil
}
func (*PaySystem) Release() {
	_self.goPool.Release()
}
func (*PaySystem) Run() {

	if _self.callback == nil {
		log.Warn("PaySystem wxCallback is nil.")
	}
	_self.server.Start()
}

// 设置一个微信支付回调，回调的参数是一个 PaySystemNotify，当支付成功时调用
func (*PaySystem) SetCallbackFunc(callback *selectCaseInterface.CallbackHandler) {
	_self.callback = callback

}

// 添加wx支付数据，wxCallback当微信服务器返回支付成功时，PaySystem会先验证消息，更新数据库，然后调用这个函数
func (*PaySystem) AddWxPay(appId string, mchId string, mckKey string) error {
	if _, ok := _self.wxMpPayMap[appId]; ok {
		return commonError.NewStringErr("There is AppId already:" + appId)
	}

	pay := NewWxMpPay(appId, mchId, mckKey, _self.wxNotifyUrl)
	_self.server.RegisterRouter(WxNotifyPath, network.RouterHandler{ProcessHttpFunc: wxNotifyReq})
	_self.wxMpPayMap[appId] = pay

	return nil
}

// 发起一笔微信支付,通过cb，返回一个prePayId string，用于客户端拉起微信
func (*PaySystem) WxPay(pd payDataStruct.WxPayReqData, cb *selectCaseInterface.CallbackHandler) error {

	if pay, ok := _self.wxMpPayMap[pd.AppId]; !ok {
		return commonError.NewStringErr("appid does not exist:" + pd.AppId)
	} else {
		_self.goPool.Go(func() {
			resp, err := pay.BeginPay(pd)
			if err != nil {
				log.Warn("wx pay failed.", "err", err.Error())
			} else {
				cb.SendReturnMsgNoReturn(resp)
			}
		}, nil)

	}
	return nil
}

// 开始验证一个客户端发过来的收据是否正确
func (*PaySystem) ApplePay(userId int, receipt string, productId string) {
	appleItem := newApplePayItem(userId, receipt, productId)
	_self.goPool.Go(func() {
		appleItem.run()
		appleItem.waitForClose()
	}, nil)
}
