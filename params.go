package pathrouter

import (
	"bytes"
)

type Params struct {
	indices [32]int
	keys    [32][]byte
	len     int
	url     string
}

func newParams(url string) *Params {
	return &Params{
		url: url,
	}
}

func (ps *Params) Get(name string) string {
	key := unsafeStringToBytes(name)

	for i := 0; i < ps.len; i++ {
		if bytes.Equal(ps.keys[i], key) {
			v := ps.url[ps.indices[i]:]
			if v[0] == '/' {
				v = v[1:]
			}

			if name == "*" {
				return v
			}

			idx := len(v)
			for j := 0; j < len(v); j++ {
				if v[j] == '/' {
					idx = j
					break
				}
			}

			return v[:idx]
		}
	}

	return ""
}

func (ps *Params) push(name []byte, pos int) {
	ps.indices[ps.len] = pos
	ps.keys[ps.len] = name
	ps.len++
}

func (ps *Params) clear() {
	ps.len = 0
}
