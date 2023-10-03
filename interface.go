package pathrouter

import "net/http"

type IRouter interface {
	IRoutes

	HandleErr(errorCode int, handler HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type IRoutes interface {
	Scope(prefix string) IRoutes
	Use(MiddlewareFunc)

	Handle(method, path string, handler http.Handler)

	Get(path string, handler HandlerFunc)
	Post(path string, handler HandlerFunc)
	Put(path string, handler HandlerFunc)
	Patch(path string, handler HandlerFunc)
	Delete(path string, handler HandlerFunc)
	Connect(paath string, handler HandlerFunc)
	Options(paath string, handler HandlerFunc)
}
