package network

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/selectCase/selectCaseInterface"

	"github.com/ntfox0001/svrLib/log"
)

type RouterWsJsonHandler struct {
	JsonMsgHandler
	upgrader websocket.Upgrader
}

func NewRouterWsJsonHandler(slHelper selectCaseInterface.ISelectLoopHelper) *RouterWsJsonHandler {
	handler := &RouterWsJsonHandler{
		upgrader:       websocket.Upgrader{},
		JsonMsgHandler: NewJsonMsgHandler(slHelper),
	}

	return handler
}

// 通过注册路由可以成为服务器，接受来自客户端的链接
func (h *RouterWsJsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if h.disableCheckOrigin {
	// 	h.upgrader.CheckOrigin = func(r *http.Request) bool {
	// 		// 同源测试再checkconn中进行
	// 		return true
	// 	}
	// }

	log.Debug("ServeHTTP", "info", "connect arrived", "ip", r.RemoteAddr)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("ServeHTTP", "err", err.Error())
		return
	}

	h.processConn(conn)

	conn.Close()
	log.Debug("ServeHTTP", "info", "connect closed", "ip", r.RemoteAddr)
}
