package util_test

import (
	"encoding/json"
	"fmt"
	"github.com/ntfox0001/svrLib/util"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

type RealtimeNoticeUnit struct {
	ExchangeName  string
	PairsCoinName string
	DataName      string
	NoticeType    string
	NoticeValue   float64
}
type RealtimeNotify struct {
	MsgId      string
	NoticeList []RealtimeNoticeUnit
}

func TestJson12(t *testing.T) {
	js := `{"ExchangeName":"11111"}`
	var rtnu RealtimeNoticeUnit
	fmt.Println(jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(js), &rtnu))
	fmt.Println(rtnu)
}
func Test_i2stru(t *testing.T) {
	f := make(map[string]interface{})
	f["msgId"] = "SetRealtimeNoticeReq"
	f["noticeList"] = make([]interface{}, 1)
	n := make(map[string]interface{})
	n["exchangeName"] = "binance"
	n["pairsCoinName"] = "ETH/BTC"
	n["dataName"] = "last"
	n["noticeType"] = "NoticeRise"
	n["noticeValue"] = 0.1
	f["noticeList"].([]interface{})[0] = n

	var d RealtimeNotify
	err := util.I2Stru(f, &d)
	if err != nil {
		t.Log(err.Error())
	}

	if d.MsgId != "SetRealtimeNoticeReq" {
		t.Fail()
	}
}

func Test_i2stru2(t *testing.T) {

	n := make(map[string]interface{})
	n["ExchangeName"] = "binance"
	n["PairsCoinName"] = "ETH/BTC"
	n["DataName"] = "last"
	n["NoticeType"] = "NoticeRise"
	n["NoticeValue"] = 0.1

	var d RealtimeNoticeUnit
	util.I2Stru(n, &d)

	if d.ExchangeName == "binance" {
		t.Log("success")
	} else {
		t.Fail()
	}
}
func Test_i2stru3(t *testing.T) {
	s := "{\"msgId\":\"SetRealtimeNoticeReq\",\"noticeList\":[{\"DataName\":\"last\",\"ExchangeName\":\"binance\",\"NoticeType\":\"NoticeRise\",\"NoticeValue\":0.1,\"PairsCoinName\":\"ETH/BTC\"},{\"DataName\":\"last\",\"ExchangeName\":\"binance\",\"NoticeType\":\"NoticeRise\",\"NoticeValue\":0.1,\"PairsCoinName\":\"ETH/USDT\"}]}"

	var v interface{}
	jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(s, &v)

	fmt.Println(v)
	var d RealtimeNoticeUnit
	util.I2Stru(v, &d)

	if d.ExchangeName == "binance" {
		t.Log("success")
	} else {
		t.Fail()
	}
}

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

func TestOther(t *testing.T) {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)

	fmt.Println(s)
}

func TestJson1(t *testing.T) {
	str := `{"msgId":"SyncRealtimeData", "timestamp":15357546477762398, "data":{"huobipro_ETH/USDT":{"open":435.38, "last":450.47, "high":463.15, "low":405.44}}}`
	var v interface{}
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(str), &v)

	fmt.Println(v.(map[string]interface{})["timestamp"].(float64))
	fmt.Println(uint64(v.(map[string]interface{})["timestamp"].(float64)))
	fmt.Println(v)
}
