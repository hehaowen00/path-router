package pathrouter

import (
	"net/http"
	"unsafe"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, ps *Params)

type paramsKey struct{}

var ParamsKey = paramsKey{}

func filter[t comparable](slice []t, check func(v t) bool) []t {
	var result []t

	for _, v := range slice {
		if check(v) {
			result = append(result, v)
		}
	}

	return result
}

func unsafeStringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func formatURL(url string) string {
	if url != "/" {
		url = url + "/"
	}
	return url
}
