package pathrouter

import (
	"net/http"
	"strings"
)

type Group struct {
	prefix     string
	middleware []MiddlewareFunc
	routes     []route
}

type route struct {
	method  string
	path    string
	handler HandlerFunc
}

func newGroup(prefix string) *Group {
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	group := Group{
		prefix:     prefix,
		middleware: make([]MiddlewareFunc, 0),
	}

	return &group
}

func (g *Group) addRoute(method, path string, handler HandlerFunc) {
	g.prefix = strings.TrimSuffix(g.prefix, "/")
	g.routes = append(g.routes, route{
		method:  method,
		path:    g.prefix + "/" + strings.TrimPrefix(path, "/"),
		handler: handler,
	})
}

func (g *Group) call(r *Router, callback func(g *Group)) {
	callback(g)
	for _, route := range g.routes {
		handler := addMiddleware(route.handler, g.middleware)
		r.getMethodHandler(route.method).Insert(route.path, handler)
	}
}

func (g *Group) Get(path string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, path, handler)
}

func (g *Group) Post(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, path, handler)
}

func (g *Group) Put(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPut, path, handler)
}

func (g *Group) Patch(path string, handler HandlerFunc) {
	g.addRoute(http.MethodPatch, path, handler)
}

func (g *Group) Delete(path string, handler HandlerFunc) {
	g.addRoute(http.MethodDelete, path, handler)
}

func (g *Group) Connect(path string, handler HandlerFunc) {
	g.addRoute(http.MethodConnect, path, handler)
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func addMiddleware(handler HandlerFunc, middleware []MiddlewareFunc) HandlerFunc {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}

	return func(w http.ResponseWriter, r *http.Request, ps *Params) {
		handler(w, r, ps)
	}
}