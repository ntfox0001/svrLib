package jsonMsg_test

import (
	"testing"

	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/network/jsonMsg"
)

func Test1(t *testing.T) {
	jm := jsonMsg.JdMsg("testMsg")
	log.Info(jm.ToJson())
	jsonMsg.JdMsgSetKeep(jm, "keepdata1", "data1")
	log.Info(jm.ToJson())
	jsonMsg.JdMsgClearKeepRoot(jm)
	log.Info(jm.ToJson())
}

func Test2(t *testing.T) {
	jm := jsonMsg.NewJsonMsg("testMsg")
	log.Info(jm.Data.ToJson())
	jm.SetKeepData("keepdata1", "data1")
	log.Info(jm.Data.ToJson())
	jm = jm.NewSubJsonMsg("testMsg2")
	log.Info(jm.Data.ToJson())
}
