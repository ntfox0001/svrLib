package paySystem

import (
	"encoding/xml"
	"fmt"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"github.com/ntfox0001/svrLib/util"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ntfox0001/svrLib/log"
)

const (
	WxPreBillUrl     = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WxNotifyFailResp = `<xml>
	<return_code><![CDATA[FAIL]]></return_code>
	<return_msg><![CDATA[INVALIDFORMAT]]></return_msg>
  </xml>`
	WxNotifySuccessResp = `<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
</xml>`
)

type WxMpPay struct {
	appId     string
	mchId     string
	mckKey    string
	notifyUrl string
	count     int64 // 用于产生bill
	goPool    *util.GoroutinePool
}

// 微信公众号支付
// 微信公众号处理流程
// 客户端通知服务器发起支付-》服务器调用“统一下单”获得支付preid-》客户端收到preid，呼叫微信支付sdk
// -》用户输入密码确认-》微信后台收到确认-》通知服务器结果，服务器处理结果-》发给客户端显示结果
// 这里解决的问题是服务器和微信服务器和数据库订单状态管理
// 监听端口用来告诉微信服务器发送数据到哪里
func NewWxMpPay(appId string, mchId string, mckKey string, notifyUrl string) *WxMpPay {
	wp := &WxMpPay{
		appId:     appId,
		mchId:     mchId,
		mckKey:    mckKey,
		count:     10232,
		notifyUrl: notifyUrl,
	}

	return wp
}

//发起一笔用户微信支付， 返回prePayId
func (w *WxMpPay) BeginPay(pd payDataStruct.WxPayReqData) (string, error) {
	req := payDataStruct.WxUnifiedorderReq{
		WxPayBase:  pd.WxPayBase,
		AppId:      w.appId,
		MchId:      w.mchId,
		Key:        w.mckKey,
		NonceStr:   "wxpay_1122h78",
		OutTradeNo: util.GetUniqueId(),
		TradeType:  "JSAPI",
		NotifyUrl:  w.notifyUrl,
	}

	// 创建签名数据
	if xmlStr, err := util.MakeWxSign(req); err != nil {
		log.Error("MakeWxSign error", "err", err.Error())
		return "", err
	} else {
		//向微信服务器发送“统一下单”请求
		wxRespStr, err := network.SyncHttpPost(WxPreBillUrl, xmlStr, network.ContentTypeText)
		fmt.Println(wxRespStr)
		if err != nil {
			log.Error("SyncHttpPost error", "err", err.Error())
			return "", err
		}
		// 解析微信返回值
		var resp payDataStruct.WxUnifiedorderResp
		if err := xml.Unmarshal([]byte(wxRespStr), &resp); err != nil {
			// 如果格式解析失败，那是一个严重错误
			log.Error("Failed to Unmarshal of WxMpPay resp.", "resp", wxRespStr)
			return "", err
		}

		// 检查微信返回值, 通信成功标识
		if resp.ReturnCode != "SUCCESS" {
			log.Warn("wx pay failed: ReturnCode", "resp", resp)
			return "", err
		}

		// 检查微信返回值，订单成功标识
		if resp.ResultCode != "SUCCESS" {
			log.Warn("wx pay failed: ResultCode", "resp", resp)
			return "", err
		}

		// 都成功了
		// 更新数据库
		extentInfo := w.appId // 额外信息保存appid
		if err := _self.PayRecord_NewBill(pd.UserId, req.OutTradeNo, pd.ProductId, pd.TotalFee, "WX", extentInfo, true); err != nil {
			log.Warn("wx pay PayRecord_NewBill failed.", "err", err.Error())
			return "", err
		}
		// 返回,等待用户付款，快输入密码，快快快~
		return resp.PrePayId, err

	}

}

// 微信支付通知
func wxNotifyReq(w http.ResponseWriter, r *http.Request) {
	s, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(s))
	req := payDataStruct.WxPayNotifyReq{}
	if err := xml.Unmarshal(s, &req); err != nil {
		log.Error("wxNotifyReq", "Marshal error", err.Error())
		io.WriteString(w, WxNotifyFailResp)
		return
	}

	if pd, err := _self.PayRecord_Query(req.OutTradeNo); err != nil {

	} else {
		// 订单状态必须是等待用户支付
		if pd.Status != payDataStruct.PayStatusWaitForUserPay {
			if pd.Status == payDataStruct.PayStatusSuccess {
				// 多余的补单，直接忽略
				return
			}
			log.Error("Invalid to WxPayBill status", "billId", pd.BillId, "status", pd.Status)
			io.WriteString(w, WxNotifyFailResp)
			return
		}

		// 设置订单完结
		_self.PayRecord_SetPayStatusSuccess(pd.BillId, req.TransactionId)

		// 调用通知函数
		if _self.callback == nil {
			log.Error("WxCallback is nil")
			io.WriteString(w, WxNotifyFailResp)
			return
		}

		notify := payDataStruct.PaySystemNotify{
			ExtentData: req,
			ProductId:  pd.ProductId,
			UserId:     pd.UserId,
			PayType:    payDataStruct.PaySystemNotify_PayType_Wx,
		}

		// 发送回调
		_self.callback.SendReturnMsgNoReturn(notify)

		// 返回成功
		io.WriteString(w, WxNotifySuccessResp)
	}

}
