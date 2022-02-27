package ext

import (
	"sort"
	"strings"
)

type Strings []string

func ExtractStrings(args ...interface{}) (Strings, []interface{}) {
	a := Strings{}
	b := []interface{}{}
	for _, i := range args {
		switch v := i.(type) {
		case string:
			a = append(a, v)
		case []string:
			a = append(a, v...)
		case Strings:
			a = append(a, v...)
		default:
			b = append(b, v)
		}
	}
	return a, b
}

func (s Strings) Set() StringSet {
	return NewStringSet(s...)
}

func (s Strings) Sort() Strings {
	sort.Strings(s)
	return s
}

func (s Strings) Equal(other Strings) bool {
	if len(s) != len(other) {
		return false
	}
	for i, v := range s {
		if v != other[i] {
			return false
		}
	}
	return true
}

func (s Strings) Include(v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

func (s Strings) Select(f func(string) bool) Strings {
	r := Strings{}
	for _, i := range s {
		if f(i) {
			r = append(r, i)
		}
	}
	return r
}

func (s Strings) Interfaces() []interface{} {
	r := []interface{}{}
	for _, i := range s {
		r = append(r, i)
	}
	return r
}

func (s Strings) Exclude(vals ...string) Strings {
	values := Strings(vals)
	return s.Select(func(i string) bool {
		return !values.Include(i)
	})
}

func (s Strings) Join(delimiter string) string {
	return strings.Join(s, delimiter)
}
