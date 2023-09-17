package pathrouter

import (
	"bytes"
	"strings"
)

func (cursor *node[v]) Get(url string, ps *Params) *v {
	path := unsafeStringToBytes(url)
	index := 0

start:
	if len(path) == 1 && path[0] == '/' {
		return cursor.value
	}

	for i := 0; i < len(cursor.lut); i++ {
		if cursor.lut[i] == path[1] {
			for j := i; j < len(cursor.children); j++ {
				child := cursor.children[j]

				if len(path) < len(child.path) {
					continue
				}

				for k := 0; k < len(child.path); k++ {
					if child.path[k] != path[:len(child.path)][k] {
						goto end
					}
				}

				cursor = child
				path = path[len(child.path)-1:]
				index += len(child.path) - 1
				goto start
			end:
			}
			break
		}
	}

	if cursor.lut[len(cursor.lut)-1] == ':' {
		cursor = cursor.children[len(cursor.lut)-1]
		end := -1

		for i := 1; i < len(path); i++ {
			if path[i] == '/' {
				end = i
				break
			}
		}

		if end == -1 {
			ps.push(cursor.path, index)
			return cursor.value
		}

		path = path[end:]
		ps.push(cursor.path, index+1)
		index = index + end

		goto start
	}

	if cursor.lut[len(cursor.lut)-1] == '*' {
		cursor = cursor.children[len(cursor.lut)-1]
		ps.push(cursor.path, index)
		return cursor.value
	}

	return nil
}

func (cursor *node[v]) Insert(path string, value v) {
	if path == "/" {
		cursor.value = &value
		return
	}

	xs := strings.Split(path, "/")
	xs = filter[string](xs, func(s string) bool {
		return s != ""
	})

start:
	if cursor.children == nil || len(cursor.children) == 0 {
		goto insertAll
	}

	for i, child := range cursor.children {
		if child.param || child.wildcard {
			if xs[0][0] == ':' {
				child.setPath(xs[0])
				cursor.lut[i] = ':'
				cursor = child
				xs = xs[1:]
				goto start
			}
			if xs[0][0] == '*' {
				child.setPath("*")
				cursor.lut[i] = '*'
				cursor = child
				xs = xs[1:]
				goto start
			}
		}

		if bytes.Equal(child.path, []byte("/"+xs[0]+"/")) {
			xs = xs[1:]
			cursor = child
			goto start
		}
	}

insertAll:
	for _, p := range xs {
		child := newNode[v]()
		child.setPath(p)
		cursor.addNode(p, child)
		cursor = child
	}
	cursor.setValue(value)
}
