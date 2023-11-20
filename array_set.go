package pathrouter

type arraySet struct {
	data []string
}

func newArraySet() arraySet {
	return arraySet{}
}

func (set *arraySet) insert(value string) {
	if len(set.data) == 0 {
		set.data = append(set.data, value)
		return
	}

	for i := 0; i < len(set.data); i++ {
		if set.data[i] == value {
			return
		}
	}

	set.data = append(set.data, value)
}
