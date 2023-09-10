package pathrouter

import (
	"strings"
)

type trie[v any] struct {
	root  *node[v]
	arena []node[v]
}

func newTrie[v any]() *trie[v] {
	trie := trie[v]{
		root: nil,
	}
	trie.root = trie.newNode()
	return &trie
}

func (t *trie[v]) newNode() *node[v] {
	node := node[v]{}
	t.arena = append(t.arena, node)
	return &t.arena[len(t.arena)-1]
}

func (t *trie[v]) Get(path *string, ps *Params) *v {
	n := t.root

start:
	if len(*path) > 0 && (*path)[0] == '/' {
		*path = (*path)[1:]
	}

	if len(*path) == 0 || *path == "/" {
		return n.value
	}

	ch := (*path)[0]
	for i := 0; i < len(n.lut); i++ {
		v := n.lut[i]
		if v == ch {
			for j := i; j < len(n.children); j++ {
				v := n.children[j]
				if v.matchPath(path, ps) {
					n = v
					goto start
				}
			}
			break
		}
	}

	if n.lut[len(n.lut)-1] == ':' {
		n = n.children[len(n.lut)-1]
		idx := -1
		for i := 0; i < len(*path); i++ {
			if (*path)[i] == '/' {
				idx = i
				break
			}
		}

		if idx > -1 {
			val := (*path)[:idx]
			*path = (*path)[idx:]
			ps.Push(n.path, val)
		} else {
			ps.Push(n.path, *path)
			return n.value
		}

		goto start
	}

	if n.lut[len(n.lut)-1] == '*' {
		n = n.children[len(n.lut)-1]
		ps.Push(n.path, *path)
		return n.value
	}

	return nil
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
		for _, p := range xs {
			child := t.newNode()
			child.SetPath(p)
			n.AddNode(p, child)
			n = child
		}
		n.SetValue(value)
		return
	} else {
		for _, v := range n.children {
			if v.param {
				if xs[0][0] == byte(':') {
					v.SetPath(xs[0])
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
			if v.path == xs[0] {
				xs = xs[1:]
				n = v
				goto start
			}
		}
	}

	for _, p := range xs {
		child := t.newNode()
		child.SetPath(p)
		n.AddNode(p, child)
		n = child
	}
	n.SetValue(value)
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
