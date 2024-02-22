package pathrouter

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

func applyMiddleware(handler http.HandlerFunc, middleware []MiddlewareFunc) http.HandlerFunc {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}

	return handler
}

type CorsHandler struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

func (cors *CorsHandler) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(cors.AllowedOrigins) == 0 {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		} else {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(cors.AllowedOrigins, ","))
		}

		if cors.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		next(w, r)
	}
}

type GZipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (rw GZipResponseWriter) Write(data []byte) (int, error) {
	return rw.writer.Write(data)
}

func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get("Accept-Encoding")
		if strings.Contains(value, "gzip") {
			w.Header().Add("Content-Encoding", "gzip")

			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			rw := GZipResponseWriter{
				ResponseWriter: w,
				writer:         gzipWriter,
			}

			next(rw, r)
		} else {
			next(w, r)
		}
	}
}
