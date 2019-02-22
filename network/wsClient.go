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
	return NewWsClient2(url, header)

}

func NewWsClient2(url string, header http.Header) (*WsClient, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	client := &WsClient{
		url:    url,
		header: header,
		resp:   resp,
		WsMsgHandler: WsMsgHandler{

			conn:        conn,
			msgMap:      make(map[string]func(*networkInterface.RawMsgData)),
			jsonMsgMap:  make(map[string]func(map[string]interface{})),
			jsonMsgChan: make(chan map[string]interface{}),
			msgChan:     make(chan networkInterface.IMsgData),
		},
	}

	if err == nil {
		return client, nil
	} else {
		return nil, err
	}
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
