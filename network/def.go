package network

import (
	"net/http"
	"net/url"
)

func CheckSameOrigin(r *http.Request, ip string, port string) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	lhost := ip + ":" + port

	return u.Host == lhost
}
