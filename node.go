package pathrouter

import (
	"slices"
	"strings"
)

type node[v any] struct {
	param    bool
	wildcard bool
	path     string
	value    *v
	lut      []byte
	children []*node[v]
}

func newNode[v any]() *node[v] {
	node := node[v]{
		path:     "",
		value:    nil,
		children: nil,
	}

	return &node
}

func (n *node[v]) SetPath(path string) {
	if string(path[0]) == ":" {
		n.param = true
		n.path = path[1:]
		return
	}

	if string(path[0]) == "*" {
		n.wildcard = true
	}
	n.path = path
}

func (n *node[v]) SetValue(value v) {
	n.value = &value
}

// adds a child node to the current node
func (n *node[v]) AddNode(key string, child *node[v]) {
	if n.children == nil {
		n.lut = append(n.lut, key[0])
		n.children = make([]*node[v], 0)
		n.children = append(n.children, child)
		return
	}

	i, _ := slices.BinarySearchFunc(n.children, child, compareNode[v])
	n.lut = slices.Insert(n.lut, i, key[0])
	n.children = slices.Insert(n.children, i, child)
}

// checks if the current node matches the url segment
func (n *node[v]) matchPath(path *string, ps *Params) bool {
	// if the node is a parameter segment, add the value to the params and remove the segment
	// from the url
	// return true if matches
	if n.param {
		val := *path
		idx := -1
		for i := 0; i < len(*path); i++ {
			if (*path)[i] == '/' {
				idx = i
				break
			}
		}

		if idx == -1 {
			val = *path
			*path = ""
		} else {
			val = (*path)[:idx]
			*path = (*path)[idx:]
		}

		ps.Push(n.path, val)
		return true
	}

	if n.wildcard {
		ps.Push(n.path, *path)
		*path = ""
		return true
	}

	// checks if the url segment is a complete match to the node path
	pathLen := len(*path)
	endOfSegmentIndex := -1
	for i := 0; i < len(*path); i++ {
		if (*path)[i] == '/' {
			endOfSegmentIndex = i
			break
		}
	}
	if endOfSegmentIndex == -1 {
		endOfSegmentIndex = pathLen
	}

	nLen := len(n.path)
	if len((*path)) >= nLen && (*path)[:nLen] == n.path && endOfSegmentIndex == nLen {
		idx := -1
		for i := 0; i < len(*path); i++ {
			if (*path)[i] == '/' {
				idx = i
				break
			}
		}

		if idx == -1 {
			*path = ""
		}

		return true
	}

	return false
}

func removeSegment(s *string) string {
	idx := -1
	for i := 0; i < len(*s); i++ {
		if (*s)[i] == '/' {
			idx = i
			break
		}
	}

	if idx == -1 {
		val := *s
		*s = ""
		return val
	}

	val := (*s)[:idx]
	*s = (*s)[idx:]

	return val
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

	if a.path == "" || b.path == "" {
		return strings.Compare(a.path, b.path)
	}

	return strings.Compare(a.path, b.path)
}
