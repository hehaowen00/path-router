package main

import (
	"fmt"
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func main() {
	r := pathrouter.NewRouter()
	r.Handle(http.MethodGet, "/*", http.FileServer(http.Dir("www")))
	r.HandleErr(
		http.StatusNotFound,
		func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
			fmt.Fprintf(w, "Page Not Found: %s\n", r.URL.Path)
		})

	addr := ":8000"
	log.Printf("started server at %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
