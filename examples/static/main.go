package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pathrouter "github.com/hehaowen00/path-router"
)

func main() {
	r := pathrouter.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request, ps *pathrouter.Params) {
		path := ps.Get("*")

		if strings.Contains(path, "..") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if len(path) == 0 {
			path = "index.html"
		}

		bytes, err := os.ReadFile("./www/" + path)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		w.Header().Add("Content-Type", contentType)
		w.Write(bytes)
	})

	addr := ":8000"
	log.Printf("started server at %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
