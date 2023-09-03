package pathrouter

import (
	"sync"
)

var paramsArena = sync.Pool{
	New: func() any {
		return &Params{}
	},
}

type Params struct {
	entries []*param
	arena   *sync.Pool
}

type param struct {
	name  string
	value string
}

func newParams() *Params {
	ps := &Params{}
	ps.arena = &paramsArena
	return ps
}

func (ps *Params) Get(name string) string {
	for _, e := range ps.entries {
		if e.name == name {
			return e.value
		}
	}

	return ""
}

func (ps *Params) Push(name, value string) {
	entry := param{
		name:  name,
		value: value,
	}

	ps.entries = append(ps.entries, &entry)
}

func (ps *Params) clear() {
	ps.entries = []*param{}
}

func (ps *Params) release() {
	ps.clear()
	ps.arena.Put(ps)
}
