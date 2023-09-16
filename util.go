package pathrouter

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, ps *Params)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type paramsKey struct{}

var ParamsKey = paramsKey{}

var _ http.Handler = (*Router)(nil)

func applyMiddleware(handler HandlerFunc, middleware []MiddlewareFunc) HandlerFunc {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}

	return func(w http.ResponseWriter, r *http.Request, ps *Params) {
		handler(w, r, ps)
	}
}
