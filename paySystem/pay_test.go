package paySystem_test

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"github.com/ntfox0001/svrLib/database"
	"github.com/ntfox0001/svrLib/paySystem/payDataStruct"
	"github.com/ntfox0001/svrLib/util"
	"testing"
)

type XmlStru struct {
	I int
	S string
}

func TestXml1(t *testing.T) {
	s, _ := xml.Marshal(XmlStru{I: 1, S: "33"})
	fmt.Println(string(s))
}
func TestUrl1(t *testing.T) {
	u, _ := url.Parse("https://pay.weixin.qq.com/wxpay/pay.action")
	fmt.Println(u.EscapedPath())
}

func TestLoadDb(t *testing.T) {
	database.Instance().Initial("192.168.1.157", "3306", "root", "123456", "bit007", 10, 10)

	op := database.Instance().NewOperation("call WxPayBill_Query(?)", "03abcb6a46d225e5989bae33557b2ec4")
	if rt, err := database.Instance().SyncExecOperation(op); err != nil {
		t.Fail()
		return
	} else {
		wxPayBillDS := rt.FirstSet()
		if len(wxPayBillDS) != 1 {
			t.Fail()
			return
		}
		// 解析数据
		var wxpaybill payDataStruct.WxPayBill
		for _, v := range wxPayBillDS {
			if err := util.I2Stru(v, &wxpaybill); err != nil {
				fmt.Println(err)
				t.Fail()
				return
			}
			break
		}
		fmt.Println(wxpaybill)
	}
	database.Instance().Release()
}
