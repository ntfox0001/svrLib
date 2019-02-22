package noticeSystem_test

import (
	"fmt"
	"github.com/ntfox0001/svrLib/network"
	"github.com/ntfox0001/svrLib/noticeSystem"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestNoticeSystem(t *testing.T) {
	r, _ := network.SyncHttpPost("http://192.168.1.117:8888/wxrobot/getList", "{}", network.ContentTypeJson)
	var resp noticeSystem.NoticeWxRobotGetListResp
	jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(r, &resp)

	for _, i := range resp.List {
		for l, k := range i {
			if l == 0 {
				ShowQun(k.UserName)
				fmt.Println("--------------------------------------------------------------------------------")
			}
			//fmt.Println(k)
		}
		fmt.Println("----------------------------------------------------------------------------------")
	}
	fmt.Println(resp.ErrorId)
}
func ShowQun(id string) {
	s := fmt.Sprintf(`{"id":"%s"}`, id)
	r, _ := network.SyncHttpPost("http://192.168.1.117:8888/wxrobot/getRoomList", s, network.ContentTypeJson)
	fmt.Println(r)
	var resp noticeSystem.NoticeWxRobotGetRoomListResp
	jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(r, &resp)

	for _, i := range resp.List {
		fmt.Println(i)
	}

}
func TestNoticeSystem2(t *testing.T) {
	fmt.Println(network.SyncHttpPost("http://192.168.1.117:8888/wxrobot/sendmsg",
		`{"id":"@c57b2f2a48038ec5627487c600f57def", "target":["@cb533e04d518712e4edd5c5583ce85b4"],"msg":"有空我请你吃饭"}`,
		network.ContentTypeJson))
}

func TestNoticeSystem3(t *testing.T) {
	fmt.Println(network.SyncHttpPost("http://192.168.1.117:8888/wxrobot/getmsg",
		`{"id":"@617d86fa0fbbd524acd8e84316490df2"}`,
		network.ContentTypeJson))
}
