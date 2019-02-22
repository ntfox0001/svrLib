package util_test

import (
	"encoding/xml"
	"fmt"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/util"
	"testing"
)

type UnifiedorderData struct {
	AppId          string `json:"appid"`            //公众账号ID
	MchId          string `json:"mch_id"`           //商户号
	DeviceInfo     string `json:"device_info"`      //设备号
	NonceStr       string `json:"nonce_str"`        //随机字符串
	Body           string `json:"body"`             //商品描述
	Attach         string `json:"attach"`           //附加数据
	OutTradeNo     string `json:"out_trade_no"`     //商户订单号
	TotalFee       string `json:"total_fee"`        //标价金额
	SpBillCreateIp string `json:"spbill_create_ip"` //终端IP
	NotifyUrl      string `json:"notify_url"`       //通知地址
	TradeType      string `json:"trade_type"`       //交易类型
	ProductId      string `json:"product_id"`       //商品ID
	OpenId         string `json:"openid"`           //用户标识
}

func TestWxSign(t *testing.T) {
	m := make(map[string]string)
	m["appid"] = "wx7a922f55b320fdf4"
	m["mch_id"] = "1508964331"
	m["device_info"] = ""
	m["body"] = "100点卡"
	m["nonce_str"] = "4213421"
	m["key"] = "gamefoxlynx1234abcxindongmajiang"
	m["notify_url"] = "http://111.198.0.142:23003/wxPayNotify"
	m["openid"] = "o5LKu0bTv7DrVX75e9rI6ZHYaSAI"
	m["total_fee"] = "1"
	m["spbill_create_ip"] = "218.60.120.69"
	m["trade_type"] = "JSAPI"
	m["product_id"] = "1"
	m["out_trade_no"] = "32142121k4j21k4j2k13jk213jk21"
	s, _ := util.MakeWxSign(m)

	fmt.Println(s)
	fmt.Println(util.ValidateWxSign(s, "gamefoxlynx1234abcxindongmajiang"))

	wxRespStr, _ := network.SyncHttpPost("https://api.mch.weixin.qq.com/pay/unifiedorder", s, network.ContentTypeText)

	fmt.Println(wxRespStr)
}
func TestWxSign2(t *testing.T) {
	x := `<Xml><nonce_str>1TUgaGSJTubp8qPuYQyNa7Z85BiH9ahP</nonce_str><trade_type>JSAPI</trade_type><attach></attach><body>100点卡</body><total_fee>1</total_fee><openid>o5LKu0bTv7DrVX75e9rI6ZHYaSAI</openid><notify_url>http://111.198.0.142:23003/wxPayNotify</notify_url><device_info></device_info><appid>wx7a922f55b320fdf4</appid><mch_id>1508964331</mch_id><out_trade_no>2e34d0125c4452c13ddfb1192b61dcec</out_trade_no><spbill_create_ip>218.60.120.74</spbill_create_ip><product_id>1</product_id><sign>C6B8A15362FC28BE02D2A2AEFB64529D</sign></Xml>`
	wxRespStr, _ := network.SyncHttpPost("https://api.mch.weixin.qq.com/pay/unifiedorder", x, network.ContentTypeText)

	fmt.Println(wxRespStr)

}

type StringResources struct {
	XMLName xml.Name `xml:"resources"`
}
type WxData struct {
	XMLName        xml.Name `xml:"xml"`
	util.StringMap `xml:"ss"`
}

func TestXml(t *testing.T) {
	w := WxData{StringMap: make(util.StringMap)}
	w.StringMap["sss"] = "fff"
	s, _ := xml.Marshal(w)
	fmt.Println(string(s))
}
func TestRandStr(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(len(util.RandString("fff")))
	}
}
