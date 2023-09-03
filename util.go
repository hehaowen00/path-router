package pathrouter

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, ps *Params)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

var _ http.Handler = (*Router)(nil)
