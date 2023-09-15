package pathrouter

import (
	"strings"
)

type trie[v any] struct {
	root *node[v]
}

func newTrie[v any]() *trie[v] {
	trie := trie[v]{
		root: newNode[v](),
	}
	return &trie
}

func get[v any](n *node[v], path string, ps *Params) *v {
	if len(path) == 0 || path == "/" {
		return n.value
	}

	for i := 0; i < len(n.lut); i++ {
		if n.lut[i] == path[1] {
			for j := i; j < len(n.children); j++ {
				child := n.children[j]
				if len(path) >= len(child.path) && child.path == path[:len(child.path)] {
					return get[v](child, path[len(child.path)-1:], ps)
				}
			}
			break
		}
	}

	if n.lut[len(n.lut)-1] == ':' {
		child := n.children[len(n.lut)-1]
		idx := -1
		for i := 1; i < len(path); i++ {
			if (path)[i] == '/' {
				idx = i
				break
			}
		}

		if idx > -1 {
			val := (path)[1:idx]
			// path = (path)[idx:]
			ps.Push(child.path, val)
			// fmt.Println(ps)
			return get[v](child, path[idx:], ps)
		} else {
			ps.Push(n.path, path)
			return n.value
		}
	}

	if n.lut[len(n.lut)-1] == '*' {
		n = n.children[len(n.lut)-1]
		ps.Push(n.path, path)
		return n.value
	}

	return nil
}

func (t *trie[v]) Get(path string, ps *Params) *v {
	n := t.root
	return get[v](n, path, ps)

	// 	if len(path) == 0 {
	// 		return n.value
	// 	}

	// start:
	// 	if path == "/" {
	// 		return n.value
	// 	}

	// 	for i := 0; i < len(n.lut); i++ {
	// 		v := n.lut[i]
	// 		if v == (path)[1] {
	// 			for j := i; j < len(n.children); j++ {
	// 				v := n.children[j]
	// 				if len(path) >= len(v.path) && v.path == (path)[:len(v.path)] {
	// 					n = v
	// 					path = (path)[len(v.path)-1:]
	// 					goto start
	// 				}
	// 			}
	// 			break
	// 		}
	// 	}

	// 	if n.lut[len(n.lut)-1] == ':' {
	// 		n = n.children[len(n.lut)-1]
	// 		idx := -1
	// 		for i := 1; i < len(path); i++ {
	// 			if (path)[i] == '/' {
	// 				idx = i
	// 				break
	// 			}
	// 		}

	// 		if idx > -1 {
	// 			val := (path)[1:idx]
	// 			path = (path)[idx:]
	// 			ps.Push(n.path, val)
	// 		} else {
	// 			ps.Push(n.path, path)
	// 			return n.value
	// 		}

	// 		goto start
	// 	}

	// 	if n.lut[len(n.lut)-1] == '*' {
	// 		n = n.children[len(n.lut)-1]
	// 		ps.Push(n.path, path)
	// 		return n.value
	// 	}

	// return nil
}

func (t *trie[v]) Insert(path string, value v) {
	if path == "/" {
		t.root.value = &value
		return
	}

	xs := strings.Split(path, "/")
	xs = filter[string](xs, func(s string) bool {
		return s != ""
	})

	n := t.root

start:
	if n.children == nil || len(n.children) == 0 {
		goto insertAll
	}

	for _, v := range n.children {
		if v.param {
			if xs[0][0] == byte(':') {
				v.setPath(xs[0])
				xs = xs[1:]
				n = v
				goto start
			}
		}
		if v.wildcard {
			if xs[0][0] == byte('*') {
				v.path = "*"
				xs = xs[1:]
				n = v
				goto start
			}
		}
		if v.path == "/"+xs[0]+"/" {
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
	n.value = &value
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
