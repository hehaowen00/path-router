package pathrouter

import (
	"unsafe"
)

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

func joinURL(lhs, rhs string) string {
	if lhs == "" {
		return rhs
	}
	if rhs == "" {
		return lhs
	}

	final := lhs

	if final[len(final)-1:] != "/" {
		final = final + "/"
	}

	if rhs[:1] == "/" {
		final = final + rhs[1:]
	}

	return final
}
