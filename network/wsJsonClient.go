package network

import (
	"net/http"
	"time"

	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

type WsJsonClient struct {
	jsonMsgHandler JsonMsgHandler
}

// 创建一个客户端，参数是一个处理消息的selectLoop
func NewWsJsonClient(slHelper selectCaseInterface.ISelectLoopHelper) *WsJsonClient {
	client := &WsJsonClient{
		jsonMsgHandler: NewJsonMsgHandler(slHelper),
	}
	return client
}
func (c *WsJsonClient) Dial(url string) {
	c.DialWithHeader(url, http.Header{})
}

// 发起一个连接
func (c *WsJsonClient) DialWithHeader(url string, header http.Header) {
	go func() {
		for {

			conn, _, err := websocket.DefaultDialer.Dial(url, header)
			if err == nil {

				c.jsonMsgHandler.processConn(conn)

				conn.Close()
				conn = nil
				log.Warn("Dial", "err", "connection broken. reconnect in 1s.", "url", url)
			} else {
				log.Warn("Dial", "err", "can't connect url. reconnect in 1s.", "url", url)
			}

			timer := time.NewTimer(time.Second)
			<-timer.C
		}
	}()
}

func (c *WsJsonClient) SendMsg(msg *JsonMsg) error {
	return c.jsonMsgHandler.SendMsg(msg)
}
func (c *WsJsonClient) SendMsgSafe(msg *JsonMsg) {
	c.jsonMsgHandler.SendMsgSafe(msg)
}
