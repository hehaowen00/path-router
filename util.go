package pathrouter

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, *Params)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

var _ http.Handler = (*Router)(nil)
