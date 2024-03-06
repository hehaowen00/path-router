package main

import (
	"fmt"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if user != "admin" || pass != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		next(w, r)
	}
}

func main() {
	r := pathrouter.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!\n")
	})

	protected := r.Scope("/")
	protected.Use(basicAuth)
	protected.Get(
		"/protected",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Protected\n")
		})

	addr := "localhost:8000"
	log.Printf("started server at http://%s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
