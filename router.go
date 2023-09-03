package pathrouter

import (
	"net/http"
	"strings"
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
		middleware:     make([]MiddlewareFunc, 0),
	}
	return &router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ps := newParams()

	subTrie := r.getMethodHandler(req.Method)
	if subTrie == nil {
		r.useErrorHandler(http.StatusMethodNotAllowed, w, req)
		return
	}

	cloned := strings.Clone(req.URL.Path)
	handler := subTrie.Get(&cloned, ps)
	if handler == nil {
		r.useErrorHandler(http.StatusNotFound, w, req)
		return
	}

	addMiddleware(*handler, r.middleware)(w, req, ps)
	defer ps.release()
}

func (r *Router) Group(prefix string, callback func(*Group)) {
	g := newGroup(prefix)
	g.call(r, callback)
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.getHandler.Insert(path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.postHandler.Insert(path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.putHandler.Insert(path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.patchHandler.Insert(path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.deleteHandler.Insert(path, handler)
}

func (r *Router) Connect(path string, handler HandlerFunc) {
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
