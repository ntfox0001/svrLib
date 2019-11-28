package network

import (
	"net/http"

	"github.com/ntfox0001/svrLib/network/networkInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

type WsClient struct {
	WsMsgHandler
	resp   *http.Response
	url    string
	header http.Header
}

func NewWsClient(url string) (*WsClient, error) {
	header := http.Header{}
	return NewWsClientWithHeader(url, header)

}

func NewWsClientWithHeader(url string, header http.Header) (*WsClient, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}
	client := &WsClient{
		url:    url,
		header: header,
		resp:   resp,
		WsMsgHandler: WsMsgHandler{

			conn:       conn,
			msgMap:     make(map[string]func(*networkInterface.RawMsgData, interface{})),
			jsonMsgMap: make(map[string]func(map[string]interface{}, interface{})),
			// jsonMsgChan: make(chan map[string]interface{}),
			// msgChan:     make(chan networkInterface.IMsgData),
		},
	}

	return client, nil
}

func (w *WsClient) Start() {
	defer func() {
		log.Debug("- connect closed.")
		w.conn.Close()
	}()

	for {
		if err := w.ProcessMsg(); err != nil {
			break
		}
	}
}

func (w *WsClient) Reconnect() error {
	conn, resp, err := websocket.DefaultDialer.Dial(w.url, w.header)
	if err != nil {
		return err
	}

	w.conn = conn
	w.resp = resp
	return nil
}
