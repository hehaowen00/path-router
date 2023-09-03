package pathrouter

import (
	"context"
	"net/http"
)

type Key string

const ParamsKey = Key("__ROUTER_PARAMS__")

func (r *Router) Handle(method, path string, handler http.Handler) {
	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		ctx := context.WithValue(r.Context(), ParamsKey, ps)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}

	r.getMethodHandler(method).Insert(path, h)
}
