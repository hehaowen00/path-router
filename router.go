package pathrouter

import (
	"context"
	"net/http"
)

type Router struct {
	getHandler     *node[HandlerFunc]
	postHandler    *node[HandlerFunc]
	putHandler     *node[HandlerFunc]
	patchHandler   *node[HandlerFunc]
	deleteHandler  *node[HandlerFunc]
	connectHandler *node[HandlerFunc]
	optionsHandler *node[HandlerFunc]
	errorHandler   map[int]HandlerFunc
	middleware     []MiddlewareFunc
}

var _ http.Handler = (*Router)(nil)

func NewRouter() *Router {
	router := Router{
		getHandler:     newNode[HandlerFunc](),
		postHandler:    newNode[HandlerFunc](),
		putHandler:     newNode[HandlerFunc](),
		patchHandler:   newNode[HandlerFunc](),
		deleteHandler:  newNode[HandlerFunc](),
		connectHandler: newNode[HandlerFunc](),
		optionsHandler: newNode[HandlerFunc](),
		errorHandler:   make(map[int]HandlerFunc, 0),
		middleware:     nil,
	}

	return &router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ps := newParams(req.URL.Path)

	var subTrie *node[HandlerFunc]

	if req.Method == http.MethodGet {
		subTrie = r.getHandler
	} else if req.Method == http.MethodPost {
		subTrie = r.postHandler
	} else if req.Method == http.MethodPut {
		subTrie = r.putHandler
	} else if req.Method == http.MethodPatch {
		subTrie = r.patchHandler
	} else if req.Method == http.MethodDelete {
		subTrie = r.deleteHandler
	} else if req.Method == http.MethodConnect {
		subTrie = r.connectHandler
	} else if req.Method == http.MethodOptions {
		subTrie = r.optionsHandler
	}

	if subTrie == nil {
		r.useErrorHandler(http.StatusMethodNotAllowed, w, req)
		return
	}

	url := req.URL.Path
	if url != "/" {
		url = url + "/"
	}

	handler := subTrie.Get(url, ps)
	if handler == nil {
		r.useErrorHandler(http.StatusNotFound, w, req)
		return
	}

	(*handler)(w, req, ps)
}

func (r *Router) Group(prefix string, callback func(*Group)) {
	g := newGroup(prefix)
	g.call(r, callback)
}

func (r *Router) Use(middleware ...MiddlewareFunc) {
	if r.middleware != nil {
		panic("router.Use can only be called once")
	}

	r.middleware = append(r.middleware, middleware...)
}

func (r *Router) Get(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.getHandler.Insert(path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.postHandler.Insert(path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.putHandler.Insert(path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.patchHandler.Insert(path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.deleteHandler.Insert(path, handler)
}

func (r *Router) Connect(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.connectHandler.Insert(path, handler)
}

func (r *Router) Handle(method, path string, handler http.Handler) {
	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		ctx := context.WithValue(r.Context(), ParamsKey, ps)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}

	r.getMethodHandler(method).Insert(path, h)
}

func (r *Router) HandleErr(errorCode int, handler HandlerFunc) {
	r.errorHandler[errorCode] = applyMiddleware(handler, r.middleware)
}

func (r *Router) getMethodHandler(method string) *node[HandlerFunc] {
	if method == http.MethodGet {
		return r.getHandler
	} else if method == http.MethodPost {
		return r.postHandler
	} else if method == http.MethodPut {
		return r.putHandler
	} else if method == http.MethodPatch {
		return r.patchHandler
	} else if method == http.MethodDelete {
		return r.deleteHandler
	} else if method == http.MethodConnect {
		return r.connectHandler
	} else if method == http.MethodOptions {
		return r.optionsHandler
	}
	return nil
}

func (r *Router) useErrorHandler(code int, w http.ResponseWriter, req *http.Request) {
	errHandler := r.errorHandler[http.StatusNotFound]
	if errHandler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	errHandler(w, req, nil)
}
