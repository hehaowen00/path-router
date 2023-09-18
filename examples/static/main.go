package main

import (
	"log"
	"net/http"

	pathrouter "github.com/hehaowen00/path-router"
)

func main() {
	r := pathrouter.NewRouter()
	r.Handle(http.MethodGet, "/*", http.FileServer(http.Dir("www")))

	addr := ":8000"
	log.Printf("started server at %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
