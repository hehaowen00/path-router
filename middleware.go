package pathrouter

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func applyMiddleware(handler HandlerFunc, middleware []MiddlewareFunc) HandlerFunc {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}

	return func(w http.ResponseWriter, r *http.Request, ps *Params) {
		handler(w, r, ps)
	}
}

type GZipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (rw GZipResponseWriter) Write(data []byte) (int, error) {
	return rw.writer.Write(data)
}

func GzipMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps *Params) {
		value := r.Header.Get("Accept-Encoding")
		if strings.Contains(value, "gzip") {
			w.Header().Add("Content-Encoding", "gzip")

			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			rw := GZipResponseWriter{
				ResponseWriter: w,
				writer:         gzipWriter,
			}

			next(rw, r, ps)
		} else {
			next(w, r, ps)
		}
	}
}
