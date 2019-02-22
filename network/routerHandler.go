package network

import (
	"net/http"
)

type RouterHandler struct {
	ProcessHttpFunc func(w http.ResponseWriter, r *http.Request)
}

func (h RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ProcessHttpFunc(w, r)
}
