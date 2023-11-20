package pathrouter

import (
	"context"
	"net/http"
	"strings"
)

type pathRouter struct {
	getHandler     *node[HandlerFunc]
	postHandler    *node[HandlerFunc]
	putHandler     *node[HandlerFunc]
	patchHandler   *node[HandlerFunc]
	deleteHandler  *node[HandlerFunc]
	connectHandler *node[HandlerFunc]
	optionsHandler *node[HandlerFunc]
	optionsTable   *node[arraySet]
	errorHandler   map[int]HandlerFunc
	middleware     []MiddlewareFunc
}

var _ http.Handler = (*pathRouter)(nil)

func NewRouter() IRouter {
	router := pathRouter{
		getHandler:     newNode[HandlerFunc](),
		postHandler:    newNode[HandlerFunc](),
		putHandler:     newNode[HandlerFunc](),
		patchHandler:   newNode[HandlerFunc](),
		deleteHandler:  newNode[HandlerFunc](),
		connectHandler: newNode[HandlerFunc](),
		optionsHandler: newNode[HandlerFunc](),
		optionsTable:   newNode[arraySet](),
		errorHandler:   make(map[int]HandlerFunc, 0),
		middleware:     nil,
	}

	router.optionsHandler.Insert("*", defaultOptionsHandler(&router))

	return &router
}

func (r *pathRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	url := formatURL(req.URL.Path)

	handler := subTrie.Get(url, ps)
	if handler == nil {
		r.useErrorHandler(http.StatusNotFound, w, req)
		return
	}

	(*handler)(w, req, ps)
}

func (r *pathRouter) Scope(prefix string) IRoutes {
	g := newScope(r, prefix)
	return g
}

func (r *pathRouter) Use(middleware MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware)
	r.optionsHandler.Insert("*", applyMiddleware(defaultOptionsHandler(r), r.middleware))
}

func (r *pathRouter) Get(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.getHandler.Insert(path, handler)
	addMethod(r, path, http.MethodGet)
}

func (r *pathRouter) Post(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.postHandler.Insert(path, handler)
	addMethod(r, path, http.MethodPost)
}

func (r *pathRouter) Put(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.putHandler.Insert(path, handler)
	addMethod(r, path, http.MethodPut)
}

func (r *pathRouter) Patch(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.patchHandler.Insert(path, handler)
	addMethod(r, path, http.MethodPatch)
}

func (r *pathRouter) Delete(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.deleteHandler.Insert(path, handler)
	addMethod(r, path, http.MethodDelete)
}

func (r *pathRouter) Connect(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.connectHandler.Insert(path, handler)
	addMethod(r, path, http.MethodConnect)
}

func (r *pathRouter) Options(path string, handler HandlerFunc) {
	handler = applyMiddleware(handler, r.middleware)
	r.optionsHandler.Insert(path, handler)
}

func (r *pathRouter) Handle(method, path string, handler http.Handler) {
	h := func(w http.ResponseWriter, r *http.Request, ps *Params) {
		ctx := context.WithValue(r.Context(), ParamsKey, ps)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	}

	r.getMethodHandler(method).Insert(path, h)
}

func (r *pathRouter) HandleErr(errorCode int, handler HandlerFunc) {
	r.errorHandler[errorCode] = applyMiddleware(handler, r.middleware)
}

func (r *pathRouter) getMethodHandler(method string) *node[HandlerFunc] {
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

func (r *pathRouter) useErrorHandler(code int, w http.ResponseWriter, req *http.Request) {
	errHandler := r.errorHandler[http.StatusNotFound]
	if errHandler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	errHandler(w, req, nil)
}

func defaultOptionsHandler(router *pathRouter) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *Params) {
		url := formatURL(r.URL.Path)

		set := router.optionsTable.Get(url, ps)
		if set == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(nil)
			return
		}

		w.Header().Set("Allow", strings.Join(set.data, ", "))
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	}
}

func addMethod(r *pathRouter, path, method string) {
	ps := newParams(path)

	set := r.optionsTable.Get(formatURL(path), ps)

	if set == nil {
		set := newArraySet()
		set.insert(method)

		r.optionsTable.Insert(path, set)
		return
	}

	(*set).insert(method)
}
