package network_test

import (
	"fmt"
	"testing"

	"github.com/ntfox0001/svrLib/selectCase"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/network"
)

func Test_Server(t *testing.T) {

	StartSvr()

}

func StartSvr() {
	sel := selectCase.NewSelectLoop("tests", 20, 20)
	sel.GetHelper().RegisterEvent("TestA", onTestA)

	svr := network.NewServer("0.0.0.0", "7777")

	handler := network.NewRouterWsJsonHandler(sel.GetHelper())
	svr.RegisterRouter("/test1", handler)

	svr.Start()
	fmt.Println("server closed.")
}

func onTestA(msg selectCaseInterface.EventChanMsg) {
	jm := msg.Content.(network.JsonMsg)
	fmt.Println(jm.Data.ToJson())
}

func StartClient() {
	sel := selectCase.NewSelectLoop("testc", 20, 20)

	client := network.NewWsJsonClient(sel.GetHelper())
	client.Dial("http://127.0.0.1:7777/test1")
	jm := network.NewJsonMsg("TestA")
	client.SendMsgSafe(jm)
}
