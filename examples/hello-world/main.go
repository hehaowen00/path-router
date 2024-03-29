package main

import (
	"fmt"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func main() {
	r := pathrouter.NewRouter()

	r.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, r.URL)
			next(w, r)
		}
	})

	r.Use(pathrouter.GzipMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})

	r.Get("/hello/:user", func(w http.ResponseWriter, r *http.Request) {
		value := r.PathValue("user")
		fmt.Fprintf(w, "Hello, %s!\n", value)
	})

	r.HandleErr(http.StatusNotFound, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Page Not Found: %s\n", r.URL.Path)
	})

	addr := "localhost:8000"
	log.Printf("started server at http://%s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
