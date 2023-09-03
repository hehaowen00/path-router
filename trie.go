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

// matches a url to a node in the trie
func (t *trie[v]) Get(path *string, ps *Params) *v {
	n := t.root

start:
	if len(*path) > 0 && []rune(*path)[0] == rune('/') {
		trimSlash(path)
	}
	if len(*path) == 0 || *path == "/" {
		return n.value
	}

	for i, v := range n.lut {
		// fmt.Println("lut", string(v), string((*path)[0]))
		if v == (*path)[0] {
			for _, v := range n.children[i:] {
				// fmt.Println("visit", v.path, v.param, v.wildcard, *path)
				if v.matchPath(path, ps) {
					n = v
					removeSegment(path)
					goto start
				}
			}
		}
		if v == byte(':') {
			n = n.children[i]
			n.matchPath(path, ps)
			goto start
		}
		if v == byte('*') {
			n = n.children[i]
			removeSegment(path)
			goto start
		}
	}

	// for _, v := range n.children {
	// 	// fmt.Println("visit", v.path, v.param, v.wildcard, *path)
	// 	if v.matchPath(path, ps) {
	// 		n = v
	// 		removeSegment(path)
	// 		goto start
	// 	}
	// }

	return nil
}

// inserts a url into the trie
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

	// if the root node has no children, add the path to the root node
start:
	if n.children == nil || len(n.children) == 0 {
		for _, p := range xs {
			child := newNode[v]()
			child.SetPath(p)
			n.AddNode(p, child)
			n = child
		}
		n.SetValue(value)
		return
	} else {
		// if the root node has children, check if the first path segment matches any of the children
		for _, v := range n.children {
			if v.param {
				if xs[0][0] == byte(':') {
					xs = xs[1:]
					v.path = xs[0][1:]
					n = v
					goto start
				}
			}
			if v.wildcard {
				if xs[0][0] == byte('*') {
					xs = xs[1:]
					v.path = "*"
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
	// if the first path segment does not match any of the children, add the path to the root node
	for _, p := range xs {
		child := newNode[v]()
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

func trimSlash(s *string) {
	if *s != "" {
		*s = string([]rune(*s)[1:])
	}
}