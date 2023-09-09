package pathrouter

type Params struct {
	entries paramMapInterface
}

type param struct {
	name  string
	value string
}

func newParams() *Params {
	ps := &Params{
		entries: newArrayMap(),
	}
	return ps
}

func (ps *Params) Get(name string) string {
	return ps.entries.get(name)
}

func (ps *Params) Push(name, value string) {
	switch (ps.entries).(type) {
	case *arraymap:
		if ps.entries.full() {
			pMap := newParamMap()
			arr := ps.entries.(*arraymap)
			for i := 0; i < int(arr.count); i++ {
				pMap.push(arr.entries[i].name, arr.entries[i].value)
				ps.entries = pMap
			}
		} else {
			ps.entries.push(name, value)
		}
	case *paramMap:
		ps.entries.push(name, value)
	}
}

func (ps *Params) clear() {
	ps.entries.clear()
}

type paramMapInterface interface {
	push(name, value string)
	get(name string) string
	clear()
	full() bool
}

type arraymap struct {
	count   int8
	entries [8]param
}

func newArrayMap() *arraymap {
	a := &arraymap{}
	return a
}

func (a *arraymap) get(name string) string {
	for i := 0; i < int(a.count); i++ {
		if a.entries[i].name == name {
			return a.entries[i].value
		}
	}
	return ""
}

func (a *arraymap) push(name, value string) {
	a.entries[a.count] = param{
		name:  name,
		value: value,
	}
	a.count++
}

func (a *arraymap) clear() {
	a.count = 0
}

func (a *arraymap) full() bool {
	return a.count == 8
}

type paramMap struct {
	entries map[string]string
}

func newParamMap() *paramMap {
	m := &paramMap{
		entries: make(map[string]string),
	}
	return m
}

func (m *paramMap) get(name string) string {
	return m.entries[name]
}

func (m *paramMap) push(name, value string) {
	m.entries[name] = value
}

func (m *paramMap) clear() {
	for k := range m.entries {
		delete(m.entries, k)
	}
}

func (m *paramMap) full() bool {
	return false
}
