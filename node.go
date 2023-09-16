package pathrouter

import (
	"bytes"
	"slices"
)

type node[v any] struct {
	param    bool
	wildcard bool
	value    *v
	path     []byte
	lut      []byte
	children []*node[v]
}

func newNode[v any]() *node[v] {
	node := node[v]{
		path:     nil,
		value:    nil,
		children: nil,
	}

	return &node
}

func (n *node[v]) setPath(path string) {
	if path[0] == ':' {
		n.param = true
		n.path = []byte(path[1:])
		return
	}

	if path[0] == '*' {
		n.wildcard = true
		n.path = []byte("*")
		return
	}

	n.path = []byte("/" + path + "/")
}

func (n *node[v]) setValue(value v) {
	n.value = &value
}

func (n *node[v]) addNode(key string, child *node[v]) {
	val := key[0]

	if child.param {
		val = ':'
	} else if child.wildcard {
		val = '*'
	}

	if n.children == nil {
		n.lut = append(n.lut, val)
		n.children = make([]*node[v], 0)
		n.children = append(n.children, child)
		return
	}

	i, _ := slices.BinarySearchFunc(n.children, child, compareNode[v])
	n.lut = slices.Insert(n.lut, i, val)
	n.children = slices.Insert(n.children, i, child)
}

func compareNode[v any](a, b *node[v]) int {
	if a.param && b.param {
		return 0
	}

	if a.wildcard && b.wildcard {
		return 0
	}

	if a.wildcard && b.param {
		return 1
	}

	if a.param && b.wildcard {
		return -1
	}

	if a.param || a.wildcard {
		return 1
	}

	if b.param || b.wildcard {
		return -1
	}

	return bytes.Compare(a.path, b.path)
}
