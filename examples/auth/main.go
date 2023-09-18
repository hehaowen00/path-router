package main

import (
	"fmt"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func basicAuth(next pathrouter.HandlerFunc) pathrouter.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		if user != "admin" && pass != "password" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		w.WriteHeader(http.StatusOK)
		next(w, r, ps)
	}
}

func main() {
	r := pathrouter.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		fmt.Fprintf(w, "Hello, World!\n")
	})

	r.Group("/", func(g *pathrouter.Group) {
		g.Use(basicAuth)
		g.Get(
			"/protected",
			func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
				fmt.Fprintf(w, "Protected\n")
			})
	})

	addr := ":8000"
	log.Printf("started server at %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
