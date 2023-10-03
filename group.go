package pathrouter

import (
	"net/http"
)

type Scope struct {
	prefix     string
	middleware []MiddlewareFunc
	routes     IRoutes
}

func newScope(router IRoutes, prefix string) *Scope {
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	group := Scope{
		prefix:     prefix,
		middleware: nil,
		routes:     router,
	}

	return &group
}

func (s *Scope) Handle(method, path string, handler http.Handler) {
	s.routes.Handle(method, path, handler)
}

func (s *Scope) Use(middleware MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware)
}

func (s *Scope) Scope(prefix string) IRoutes {
	s2 := newScope(s, prefix)
	return s2
}

func (s *Scope) Get(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Get(p, h)
}

func (s *Scope) Post(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Post(p, h)
}

func (s *Scope) Put(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Put(p, h)
}

func (s *Scope) Patch(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Patch(p, h)
}

func (s *Scope) Delete(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Delete(p, h)
}

func (s *Scope) Connect(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Connect(p, h)
}

func (s *Scope) Options(path string, handler HandlerFunc) {
	h := applyMiddleware(handler, s.middleware)
	p := joinURL(s.prefix, path)
	s.routes.Options(p, h)
}
