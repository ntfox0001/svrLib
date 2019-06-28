package network

import (
	"errors"
	"reflect"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

const (
	Msg_Info_Conn_Id = "Msg_Info_Conn_Id"
)

type registerConnPair struct {
	id   uint64
	conn *websocket.Conn
}

// 每个客户端连入，都会创建一个Conn，然后在自己的循环中向selectLoop发送消息
// 也可以主动链接目标url
type JsonMsgHandler struct {
	connMap            map[uint64]*websocket.Conn
	selectLoopHelper   selectCaseInterface.ISelectLoopHelper
	unqiueId           uint64
	registerConnChan   chan registerConnPair
	unregisterConnChan chan uint64
}

// 处理web socket长连接，
func NewJsonMsgHandler(selectLoopHelper selectCaseInterface.ISelectLoopHelper) JsonMsgHandler {
	router := JsonMsgHandler{
		connMap:            make(map[uint64]*websocket.Conn),
		selectLoopHelper:   selectLoopHelper,
		unqiueId:           0,
		registerConnChan:   make(chan registerConnPair, 1),
		unregisterConnChan: make(chan uint64, 1),
	}

	router.selectLoopHelper.AddSelectCaseFront(reflect.ValueOf(router.registerConnChan), router.onRegisterConnChan)
	// 反注册要放在最前面
	router.selectLoopHelper.AddSelectCaseFront(reflect.ValueOf(router.unregisterConnChan), router.onUnregisterConnChan)

	return router
}

func (h *JsonMsgHandler) GetSelectLoopHelper() selectCaseInterface.ISelectLoopHelper {
	return h.selectLoopHelper
}

// 向conn发送消息，msg中必须带有KeepData：ConnId，非线程安全
func (h *JsonMsgHandler) SendMsg(msg *JsonMsg) error {

	jdConnId := msg.GetKeepData(Msg_Info_Conn_Id)
	if jdConnId == nil {
		log.Error("SendMsg", "err", "NotFoundConnId")
		return errors.New("NotFoundConnId")
	}

	connId := jdConnId.GetUInt64()
	if conn, ok := h.connMap[connId]; ok {
		data := []byte(msg.Data.ToJson())
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return err
		}
		defer func() {
			w.Close()
		}()
		if _, err := w.Write(data); err != nil {
			return err
		}
	} else {
		log.Error("SendMsg", "err", "NotFoundConn", "id", connId)
		return errors.New("NotFoundConn")
	}
	return nil
}

// 向conn发送消息，msg中必须带有KeepData：ConnId，线程安全
func (h *JsonMsgHandler) SendMsgSafe(msg *JsonMsg) {
	h.GetSelectLoopHelper().RunIn(func() {
		h.SendMsg(msg)
	})
}

func (h *JsonMsgHandler) Close() {
	for _, v := range h.connMap {
		v.Close()
	}
}
