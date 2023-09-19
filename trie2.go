package pathrouter

import (
	"bytes"
	"fmt"
	"strings"
)

/*

radix trie notes

/ -> current node

/a/b/c and /a/b/d should match /a/b/

/a/b/:d and /a/b/* should match all

/a/b and /a/c should remove /b from current and make a new child node and take
the value in current and trim the prefix from path and add child ndoe

to match node paths

split both urls by / and compare each segment individually
*/

// assume both paths start with /
func comparePaths(a, b string) {
	fmt.Println("start", a, b)
	for len(a) > 0 && len(b) > 0 {
		a = a[1:]
		b = b[1:]

		aI := strings.IndexByte(a, '/')
		bI := strings.IndexByte(b, '/')
		if aI == -1 || bI == -1 {
			break
		}

		if a[0] == ':' || a[0] == '*' {
			return
		}

		if b[0] == ':' || b[0] == '*' {
			return
		}

		if a[:aI] == b[:bI] {
			a = a[aI:]
			b = b[bI:]
			fmt.Println(a, b)
			continue
		}
		break
	}

	fmt.Println("final", a, b)
}

func longest(path string, index *int) string {
	// fmt.Println(path)
	if len(path) <= 1 {
		return path
	}
	if path[0] == '/' && (path[1] == ':' || path[1] == '*') {
		idx := strings.IndexByte(path[1:], '/')
		if idx == -1 {
			*index += len(path)
			return path[1:]
		}
		*index += idx + 1
		return path[1 : idx+1]
	}

	for i := 0; i < len(path); i++ {
		if path[i] == ':' || path[i] == '*' {
			*index += i - 1
			return path[:i-1]
		}
	}
	*index += len(path)
	return path[1:]
}

func (cursor *node[v]) Insert2(path string, value v) {
	// fmt.Println("path", path)
	path = formatPath(path)
start:
	if len(path) == 0 || path == "/" {
		cursor.value = &value
		return
	}

	if len(cursor.children) == 0 {
		goto insertAll
	}

	for i := 0; i < len(cursor.children); i++ {
		if cursor.lut[i] == path[1] {
			n := cursor.children[i]
			fmt.Println("check node", formatPath(string(n.path)), formatPath(path))
			least := min(len(n.path)-1, len(formatPath(path))-1)
			arg1 := formatPath(string(n.path))
			arg2 := formatPath(path)
			if arg1[least] != '/' || arg2[least] != '/' {
				continue
			}
			fmt.Println("args", string(arg1), string(arg2))
			if arg1[:least] == arg2[:least] {
				nEnd := n.path[least:]
				nStart := n.path[:least]
				fmt.Println("twin", string(nEnd), string(nStart))
				if len(nEnd) != 0 && string(nEnd) != "/" {
					c := newNode[v]()
					c.value = n.value
					n.value = nil
					c.lut = n.lut
					n.lut = nil
					c.children = n.children
					n.children = nil
					n.setPath(string(nStart))
					c.setPath(string(nEnd))
					n.addNode(formatPath(string(nEnd)), c)
					// fmt.Println("split", string(c.path), string(n.path))
				}
				cursor = n
				path = path[least:]
				goto start
			}
			// comparePaths(string(cursor.children[i].path), path[index:])
			continue
		}
		break
	}

insertAll:
	index := 0
	for index < len(path) {
		cs := longest(path[index:], &index)
		if cs == "" || cs == "/" {
			break
		}
		child := newNode[v]()
		child.setPath(cs)
		cursor.addNode(formatPath(cs), child)
		cursor = child
	}
	cursor.setValue(value)
}

func (cursor *node[v]) Get2(url string, ps *Params) *v {
	path := unsafeStringToBytes(url)
	fmt.Println("get2", string(path))
	index := 0

start:
	if len(path) == 0 || bytes.Equal(path, []byte("/")) {
		if cursor.value == nil {
			if cursor.lut[len(cursor.lut)-1] == '*' {
				cursor = cursor.children[len(cursor.lut)-1]
				ps.push(cursor.path, index)
				return cursor.value
			}
		}

		return cursor.value
	}

	if len(cursor.lut) == 0 {
		return nil
	}

	for i := 0; i < len(cursor.lut); i++ {
		fmt.Println("find", string(path[0]))
		if cursor.lut[i] == path[1] {
			for j := i; j < len(cursor.children); j++ {
				child := cursor.children[j]
				fmt.Println("compare static", string(child.path), string(path))

				if len(path) < len(child.path) {
					continue
				}

				fmt.Println("check")

				for k := 0; k < len(child.path); k++ {
					if child.path[k] != path[:len(child.path)][k] {
						goto end
					}
				}

				fmt.Println("eq")

				cursor = child
				path = path[len(child.path)-1:]
				index += len(child.path) - 1
				fmt.Println("end", string(path))
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

func formatPath(p string) string {
	if p[0] != '/' {
		p = "/" + p
	}
	if len(p) > 1 {
		if p[len(p)-1] != '/' {
			p = p + "/"
		}
	}
	return p
}
