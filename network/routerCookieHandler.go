package network

import (
	"net/http"

	"github.com/ntfox0001/svrLib/network/networkInterface"

	"github.com/gorilla/securecookie"
)

type RouterCookieHandler struct {
	Cookie          *securecookie.SecureCookie
	ProcessHttpFunc func(w http.ResponseWriter, r *http.Request)

	w http.ResponseWriter
	r *http.Request
}

func NewRouterCookieHandler(f func(w http.ResponseWriter, r *http.Request)) networkInterface.IRouterHandler {
	// 初始化cookie
	// Hash keys should be at least 32 bytes long
	var hashKey = securecookie.GenerateRandomKey(32)
	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	var blockKey = securecookie.GenerateRandomKey(16)

	router := &RouterCookieHandler{
		Cookie:          securecookie.New(hashKey, blockKey),
		ProcessHttpFunc: f,
	}

	return router
}

func (h *RouterCookieHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.w = w
	h.r = r
	h.ProcessHttpFunc(w, r)
}

func (h *RouterCookieHandler) SetCookieHandler(value map[string]string) {

	if encoded, err := h.Cookie.Encode("cookie-name", value); err == nil {
		cookie := &http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(h.w, cookie)
	}
}

func (h *RouterCookieHandler) ReadCookieHandler() map[string]string {
	if cookie, err := h.r.Cookie("cookie-name"); err == nil {
		value := make(map[string]string)
		if err = h.Cookie.Decode("cookie-name", cookie.Value, &value); err == nil {

			return value
		}
	}
	return nil
}
