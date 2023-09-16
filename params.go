package pathrouter

import (
	"bytes"
	"strings"
)

type Params struct {
	locations [32]int
	names     [32][]byte
	len       int
}

func newParams() *Params {
	return &Params{}
}

func (ps *Params) Get(path string, name string) string {
	nameBytes := []byte(name)
	for i := 0; i < ps.len; i++ {
		if bytes.Equal(ps.names[i], nameBytes) {
			if name == "*" {
				return path[ps.locations[i]:]
			}
			v := strings.TrimPrefix(path[ps.locations[i]:], "/")
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

func (ps *Params) Push(name []byte, pos int) {
	ps.locations[ps.len] = pos
	ps.names[ps.len] = name
	ps.len++
}

func (ps *Params) Clear() {
	ps.len = 0
}
