package network

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/litjson"
	"github.com/ntfox0001/svrLib/log"
	"github.com/ntfox0001/svrLib/network/jsonMsg"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"
)

// 获得一个本次启动唯一id,线程安全
func (h *JsonMsgHandler) nextUnqiueId() uint64 {
	return atomic.AddUint64(&h.unqiueId, 1)
}
func (h *JsonMsgHandler) registerConn(conn *websocket.Conn) uint64 {
	connPair := registerConnPair{id: h.nextUnqiueId(), conn: conn}
	h.registerConnChan <- connPair
	return connPair.id
}
func (h *JsonMsgHandler) unregisterConn(id uint64) {
	h.unregisterConnChan <- id
}

func (h *JsonMsgHandler) processConn(conn *websocket.Conn) {
	id := h.registerConn(conn)
	h.processMsg(conn, id)
	h.unregisterConn(id)
}

func (h *JsonMsgHandler) processMsg(conn *websocket.Conn, id uint64) {
	for {
		if err := h._processMsg(conn, id); err != nil {
			break
		}
	}
}

// 处理消息
func (h *JsonMsgHandler) _processMsg(conn *websocket.Conn, id uint64) (rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = commonError.NewCommErr(err.(error).Error(), NetErrorProcessMsg)
		}
	}()

	mt, msg, err := conn.ReadMessage()
	if err != nil {
		// 读取错误，直接断开
		log.Warn("network", "readMessageErr", err.Error())

		return commonError.NewCommErr(err.Error(), NetErrorReadMsg)
	}
	// 未知的原因导致msg为空，那么跳过
	if mt == websocket.TextMessage && len(msg) > 0 {
		// is json msg
		msgjd := litjson.NewJsonDataFromJson(string(msg))
		if msgjd != nil {
			h.dispatchJsonMsg(msgjd, id)
		} else {
			// 解析错误直接断开
			log.Error("network", "Unmarshal", err.Error(), "msg", string(msg))
			return commonError.NewCommErr(err.Error(), NetErrorUnmarshal)
		}
	} else {
		return commonError.NewCommErr("NetErrorInvalidMsgType", NetErrorInvalidMsgType)
	}
	return nil
}

func (h *JsonMsgHandler) dispatchJsonMsg(msg *litjson.JsonData, id uint64) {
	msgIdjd := msg.Get("msgId")
	if msgIdjd == nil {
		log.Error("dispatchJsonMsg", "err", "NotFoundMsgId")
		return
	}
	sMsgId := msg.Get("msgId").GetString()

	//jm := NewJsonMsgWithData(msg)
	// 设置ConnId
	jsonMsg.JdMsgSetKeep(msg, Msg_Info_Conn_Id, id)
	//jm.SetKeepData(Msg_Info_Conn_Id, id)

	if sMsgId != "" {
		h.GetSelectLoopHelper().SendMsgToMe(selectCaseInterface.NewEventChanMsg(sMsgId, nil, msg))
	} else {
		log.Warn("dispatchJsonMsg", "err", "InvalidMsgFormat", "msg", msg)
	}
}

func (h *JsonMsgHandler) onRegisterConnChan(data interface{}) bool {
	regData := data.(registerConnPair)
	if _, ok := h.connMap[regData.id]; !ok {
		h.connMap[regData.id] = regData.conn
	} else {
		log.Error("RegisterConn", "err", "IdHasExist")
	}

	return true
}
func (h *JsonMsgHandler) onUnregisterConnChan(data interface{}) bool {
	id := data.(uint64)

	if _, ok := h.connMap[id]; ok {
		delete(h.connMap, id)
	} else {
		log.Error("UnregisterConn", "err", "IdNotExist")
	}

	return true
}
