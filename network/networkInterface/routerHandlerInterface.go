package networkInterface

import (
	"net/http"
)

type IRouterHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
