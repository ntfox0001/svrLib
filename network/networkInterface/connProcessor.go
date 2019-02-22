package networkInterface

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type IConnProcessor interface {
	NewMsgHandler(c *websocket.Conn, r *http.Request) IMsgHandler
	// 首先调用，返回值控制是否继续
	CheckConn(w http.ResponseWriter, r *http.Request) bool
	Fetch(mh IMsgHandler) bool
	Close(mh IMsgHandler)
}
