package pathrouter

import "net/http"

type IRouter interface {
	IRoutes

	HandleErr(errorCode int, handler http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type IRoutes interface {
	Scope(prefix string) IRoutes
	Use(MiddlewareFunc)

	Handle(method, path string, handler http.Handler)

	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Patch(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
	Connect(paath string, handler http.HandlerFunc)
	Options(paath string, handler http.HandlerFunc)
}
