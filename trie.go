package pathrouter

import (
	"bytes"
	"strings"
)

func newTrie[v any]() *node[v] {
	return newNode[v]()
}

func (t *node[v]) Get(path []byte, ps *Params) *v {
	n := t

	index := 0

start:
	if len(path) == 1 && path[0] == '/' {
		return n.value
	}

	for i := 0; i < len(n.lut); i++ {
		v := n.lut[i]
		if v == path[1] {
			for j := i; j < len(n.children); j++ {
				v := n.children[j]

				if len(path) < len(v.path) {
					continue
				}

				for i := 0; i < len(v.path); i++ {
					if v.path[i] != path[:len(v.path)][i] {
						goto end
					}
				}

				n = v
				path = path[len(v.path)-1:]
				index += len(v.path) - 1
				goto start
			end:
			}
			break
		}
	}

	if n.lut[len(n.lut)-1] == ':' {
		n = n.children[len(n.lut)-1]
		idx := -1
		for i := 1; i < len(path); i++ {
			if path[i] == '/' {
				idx = i
				break
			}
		}

		if idx > -1 {
			path = path[idx:]
			ps.push(n.path, index+1)
			index = index + idx
		} else {
			ps.push(n.path, index)
			return n.value
		}

		goto start
	}

	if n.lut[len(n.lut)-1] == '*' {
		n = n.children[len(n.lut)-1]
		ps.push(n.path, index)
		return n.value
	}

	return nil
}

func (t *node[v]) Insert(path string, value v) {
	if path == "/" {
		t.value = &value
		return
	}

	xs := strings.Split(path, "/")
	xs = filter[string](xs, func(s string) bool {
		return s != ""
	})

	n := t

start:
	if n.children == nil || len(n.children) == 0 {
		goto insertAll
	}

	for _, v := range n.children {
		if v.param {
			if xs[0][0] == ':' {
				v.setPath(xs[0])
				xs = xs[1:]
				n = v
				goto start
			}
		}
		if v.wildcard {
			if xs[0][0] == '*' {
				v.path = []byte("*")
				xs = xs[1:]
				n = v
				goto start
			}
		}
		if bytes.Equal(v.path, []byte("/"+xs[0]+"/")) {
			xs = xs[1:]
			n = v
			goto start
		}
	}

insertAll:
	for _, p := range xs {
		child := newNode[v]()
		child.setPath(p)
		n.addNode(p, child)
		n = child
	}
	n.setValue(value)
}

func filter[t comparable](slice []t, check func(v t) bool) []t {
	var result []t

	for _, v := range slice {
		if check(v) {
			result = append(result, v)
		}
	}

	return result
}
