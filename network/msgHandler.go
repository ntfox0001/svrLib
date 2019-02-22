package network

import (
	"fmt"
	"net/http"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network/msgData"
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"runtime/debug"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
	jsoniter "github.com/json-iterator/go"
)

type WsMsgHandler struct {
	conn        *websocket.Conn
	msgMap      map[string]func(msg *networkInterface.RawMsgData)
	jsonMsgMap  map[string]func(msg map[string]interface{})
	request     *http.Request
	jsonMsgChan chan map[string]interface{}
	msgChan     chan networkInterface.IMsgData
	useExternal bool
	headMsg     msgData.MsgHead

	jsonMsgHandler func(msg map[string]interface{})
	msgHandler     func(f *networkInterface.RawMsgData)
}

func NewMsgHander(conn *websocket.Conn, r *http.Request) networkInterface.IMsgHandler {
	return &WsMsgHandler{
		conn:        conn,
		request:     r,
		msgMap:      make(map[string]func(*networkInterface.RawMsgData)),
		jsonMsgMap:  make(map[string]func(map[string]interface{})),
		jsonMsgChan: make(chan map[string]interface{}),
		msgChan:     make(chan networkInterface.IMsgData),
		useExternal: true,
		headMsg:     msgData.MsgHead{},
	}
}

func (h *WsMsgHandler) GetRequest() *http.Request {
	return h.request
}
func (h *WsMsgHandler) Disconnect() {
	h.conn.Close()
	h.headMsg.Reset()
}

func (h *WsMsgHandler) SendMsg(msg networkInterface.IMsgData) error {
	msgName := proto.MessageName(msg)
	head := msgData.MsgHead{
		MsgName: msgName,
	}

	if headBuf, err := head.Marshal(); err == nil {
		err := h.conn.WriteMessage(websocket.BinaryMessage, headBuf)
		if err != nil {
			log.Error("network", "writeMsg", err.Error())
			return err
		} else {
			if msgBuf, err := msg.Marshal(); err == nil {
				err := h.conn.WriteMessage(websocket.BinaryMessage, msgBuf)
				if err != nil {
					log.Error("network", "writeMsg", err.Error())
					return err
				}
			} else {
				log.Error("network", "Marshal", err.Error())
				return err
			}
		}
	} else {
		log.Error("network", "Marshal", err.Error())
		return err
	}

	return nil
}
func (h *WsMsgHandler) SendJsonMsg(msg interface{}) error {
	return h.conn.WriteJSON(msg)
}

func (h *WsMsgHandler) DispatchJsonMsg(msg map[string]interface{}) error {
	if msgId, ok := msg["msgId"]; !ok {
		log.Error("network", "jsonMsg", "msgId does not exist .")
		return commonError.NewCommErr("jsonMsg: msgId does not exist.", NetErrorUnknowMsg)
	} else {
		if h.jsonMsgHandler != nil {
			if err := callJsonFunc(h.jsonMsgHandler, msg); err != nil {
				return err
			}
		} else {
			if handler, ok := h.jsonMsgMap[msgId.(string)]; ok {
				// handler应该有缓存处理
				if err := callJsonFunc(handler, msg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func callJsonFunc(handler func(map[string]interface{}), msg map[string]interface{}) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("\n%s\n", string(debug.Stack()))
			log.Error("network", "logic error:", err.(error).Error(), "msg", msg, "\nstack", es)
			rterr = commonError.NewCommErr("logic error:"+err.(error).Error(), NetErrorLogic)
		}
	}()
	handler(msg)
	return nil
}
func (h *WsMsgHandler) RegisterJsonMsg(msgId string, handler func(map[string]interface{})) error {
	if _, ok := h.jsonMsgMap[msgId]; ok {
		return commonError.NewCommErr("msgId has exist:"+msgId, NetErrorExistMsg)
	} else {
		h.jsonMsgMap[msgId] = handler
	}
	return nil
}
func (h *WsMsgHandler) DispatchMsg(msg *networkInterface.RawMsgData) error {

	if msg.Name() == "" {
		log.Error("network", "msg", "msgId is nil.")
		return commonError.NewCommErr("msg: msgId is nil.", NetErrorUnknowMsg)
	}
	if h.msgHandler != nil {
		if err := callFunc(h.msgHandler, msg); err != nil {
			return err
		}
	} else {
		if handler, ok := h.msgMap[msg.Name()]; ok {
			// handler应该有缓存处理
			if err := callFunc(handler, msg); err != nil {
				return err
			}
		} else {
			log.Warn("Unknow msgId", "msgId", msg.Name())
		}
	}
	return nil
}
func callFunc(handler func(*networkInterface.RawMsgData), msg *networkInterface.RawMsgData) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			es := fmt.Sprintf("\n%s\n", string(debug.Stack()))
			log.Error("network", "logic error:", err.(error).Error(), "msg", msg, "\nstack", es)

			rterr = commonError.NewCommErr("logic error:"+err.(error).Error(), NetErrorLogic)
		}
	}()
	handler(msg)
	return nil
}
func (h *WsMsgHandler) RegisterMsg(msgId string, handler func(*networkInterface.RawMsgData)) error {
	if _, ok := h.msgMap[msgId]; ok {
		return commonError.NewCommErr("msgId has exist:"+msgId, NetErrorExistMsg)
	} else {
		h.msgMap[msgId] = handler
	}
	return nil
}

func (h *WsMsgHandler) SetDispatchJsonMsgHandler(f func(map[string]interface{})) {
	h.jsonMsgHandler = f
}
func (h *WsMsgHandler) SetDispatchMsgHandler(f func(*networkInterface.RawMsgData)) {
	h.msgHandler = f
}

func (h *WsMsgHandler) ProcessMsg() (rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = commonError.NewCommErr(err.(error).Error(), NetErrorProcessMsg)
		}
	}()

	mt, msg, err := h.conn.ReadMessage()
	if err != nil {
		// 读取错误，直接断开
		log.Warn("network", "readMessageErr:", err.Error())
		rtErr = commonError.NewCommErr(err.Error(), NetErrorReadMsg)
		return
	}
	if mt == websocket.TextMessage {
		// is json msg
		var jsonMsg interface{}
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(msg, &jsonMsg)
		if err == nil {
			if jm, ok := jsonMsg.(map[string]interface{}); ok {

				if err := h.DispatchJsonMsg(jm); err != nil {
					//逻辑错误
					if err.(commonError.CommError).GetType() == NetErrorUnknowMsg {
						return err
					}
				}
			} else {
				log.Error("invalid to jsonMsg format.", "msg", string(msg))
				return commonError.NewCommErr("invalid to jsonMsg format.", NetErrorJsonFormat)
			}
		} else {
			// 解析错误直接断开
			log.Error("network", "Unmarshal", err.Error(), "msg", string(msg))
			return commonError.NewCommErr(err.Error(), NetErrorUnmarshal)
		}
	} else if mt == websocket.BinaryMessage {
		// is protobuf msg
		if h.headMsg.MsgName == "" {
			if err := h.headMsg.Unmarshal(msg); err != nil {
				log.Error("network", "Unmarshal", err.Error())
				h.headMsg.Reset()
			}
		} else {
			rawMsg := networkInterface.NewRawMsgData(h.headMsg.MsgName, msg)

			h.DispatchMsg(rawMsg)

			h.headMsg.Reset()
		}

	}
	rtErr = nil
	return
}
