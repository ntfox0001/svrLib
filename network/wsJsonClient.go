package network

import (
	"net/http"
	"time"

	"github.com/ntfox0001/svrLib/network/jsonMsg"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

type WsJsonClient struct {
	jsonMsgHandler JsonMsgHandler
	connId         uint64
	conn           *websocket.Conn
}

// 创建一个客户端，参数是一个处理消息的selectLoop
func NewWsJsonClient(slHelper selectCaseInterface.ISelectLoopHelper) *WsJsonClient {
	client := &WsJsonClient{
		jsonMsgHandler: NewJsonMsgHandler(slHelper),
		connId:         0,
	}
	return client
}
func (c *WsJsonClient) Dial(url string) <-chan uint64 {
	return c.DialWithHeader(url, http.Header{})
}

// 发起一个连接
func (c *WsJsonClient) DialWithHeader(url string, header http.Header) <-chan uint64 {
	rtid := make(chan uint64, 1)
	go func() {
		for {

			conn, _, err := websocket.DefaultDialer.Dial(url, header)
			if err == nil {
				c.conn = conn
				id := c.jsonMsgHandler.registerConn(c.conn)
				c.connId = id
				rtid <- id
				c.jsonMsgHandler.processMsg(c.conn, id)
				c.jsonMsgHandler.unregisterConn(id)

				c.conn.Close()
				c.conn = nil
				log.Warn("Dial", "err", "connection broken. reconnect in 1s.", "url", url)
			} else {
				log.Warn("Dial", "err", err.Error(), "url", url, "desc", "can't connect url. reconnect in 1s.")
			}

			timer := time.NewTimer(time.Second)
			<-timer.C
		}
	}()
	return rtid
}

func (c *WsJsonClient) SendMsg(msg *jsonMsg.JsonMsg) error {
	msg.SetKeepData(Msg_Info_Conn_Id, c.connId)
	return c.jsonMsgHandler.SendMsg(msg)
}
func (c *WsJsonClient) SendMsgSafe(msg *jsonMsg.JsonMsg) {
	msg.SetKeepData(Msg_Info_Conn_Id, c.connId)
	c.jsonMsgHandler.SendMsgSafe(msg)
}
