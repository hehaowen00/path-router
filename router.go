package pathrouter

import (
	"net/http"
)

type Router struct {
	getHandler     *trie[HandlerFunc]
	postHandler    *trie[HandlerFunc]
	putHandler     *trie[HandlerFunc]
	patchHandler   *trie[HandlerFunc]
	deleteHandler  *trie[HandlerFunc]
	connectHandler *trie[HandlerFunc]
	errorHandler   map[int]HandlerFunc
	middleware     []MiddlewareFunc
}

func NewRouter() *Router {
	router := Router{
		getHandler:     newTrie[HandlerFunc](),
		postHandler:    newTrie[HandlerFunc](),
		putHandler:     newTrie[HandlerFunc](),
		patchHandler:   newTrie[HandlerFunc](),
		deleteHandler:  newTrie[HandlerFunc](),
		connectHandler: newTrie[HandlerFunc](),
		errorHandler:   make(map[int]HandlerFunc, 0),
		middleware:     nil,
	}
	return &router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ps := newParams()

	var subTrie *trie[HandlerFunc]

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

	// applyMiddleware(*handler, r.middleware)(w, req, ps)
	(*handler)(w, req, ps)
}

func (r *Router) Group(prefix string, callback func(*Group)) {
	g := newGroup(prefix)
	g.call(r, callback)
}

func (r *Router) Use(middleware ...MiddlewareFunc) {
	if r.middleware == nil {
		r.middleware = append(r.middleware, middleware...)
	}
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

func (r *Router) getMethodHandler(method string) *trie[HandlerFunc] {
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
