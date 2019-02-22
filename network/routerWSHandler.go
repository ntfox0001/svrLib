package network

import (
	"net/http"

	"github.com/ntfox0001/svrLib/network/networkInterface"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/svrLib/log"
)

type RouterWSHandler struct {
	upgrader           websocket.Upgrader
	processor          networkInterface.IConnProcessor
	disableCheckOrigin bool // 禁止同源检查
}

func NewRouterWSHandler(processor networkInterface.IConnProcessor) *RouterWSHandler {

	return &RouterWSHandler{
		upgrader:           websocket.Upgrader{},
		processor:          processor,
		disableCheckOrigin: false,
	}
}

// 设置同源检查，由于禁止这项会造成安全问题
func (h *RouterWSHandler) DisableCheckOrigin(s bool) {
	h.disableCheckOrigin = s
}

func (h *RouterWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.processor.CheckConn(w, r) == false {
		return
	}
	if h.disableCheckOrigin {
		h.upgrader.CheckOrigin = func(r *http.Request) bool {
			// 同源测试再checkconn中进行
			return true
		}
	}
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("network", "upgradeError:", err.Error())
		return
	}

	mh := h.processor.NewMsgHandler(c, r)

	defer func() {
		log.Debug("- connect closed.")
		c.Close()
		h.processor.Close(mh)
	}()

	if !h.processor.Fetch(mh) {
		return
	}

	for {
		if err := mh.ProcessMsg(); err != nil {
			log.Debug("- connect exit msg loop")
			break
		}
	}

}
