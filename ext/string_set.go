package ext

type StringSet map[string]bool

func ExtractStringSets(args ...interface{}) ([]StringSet, []interface{}) {
	a := []StringSet{}
	b := []interface{}{}
	for _, i := range args {
		switch v := i.(type) {
		case StringSet:
			a = append(a, v)
		default:
			b = append(b, v)
		}
	}
	return a, b
}

func NewStringSet(s ...string) StringSet {
	r := StringSet{}
	for _, i := range s {
		r.Add(i)
	}
	return r
}

func (m StringSet) Add(s string) {
	m[s] = true
}

func (m StringSet) Remove(s string) {
	delete(m, s)
}

func (m StringSet) Include(s string) bool {
	_, ok := m[s]
	return ok
}

func (m StringSet) Slice() Strings {
	r := make(Strings, len(m))
	i := 0
	for k, _ := range m {
		r[i] = k
		i++
	}
	return r
}

func (m StringSet) Equal(other StringSet) bool {
	if len(m) != len(other) {
		return false
	}
	for k, _ := range m {
		if !other.Include(k) {
			return false
		}
	}
	return true
}
