package qualify

import "github.com/bsm/qualify/internal/intsets"

// StrDict can be used to convert string values to integers
// usingdictionary encoding.
type StrDict map[string]int

// NewStrDict inits a new dictionary
func NewStrDict() StrDict { return make(StrDict) }

// Add adds a string to the dictionary
// returning the numeric value/index. Returns the
// previous index value if the string has been added
// before.
func (m StrDict) Add(s string) int {
	v, ok := m[s]
	if !ok {
		v = len(m)
		m[s] = v
	}
	return v
}

// Lookup performs a lookup of multiple strings and appends
// the results to dst returning it
func (m StrDict) Lookup(dst []int, strs ...string) []int {
	for _, s := range strs {
		if v, ok := m[s]; ok {
			dst = append(dst, v)
		}
	}
	return dst
}

// --------------------------------------------------------------------

type none struct{}

// --------------------------------------------------------------------

type fieldSet map[Field]none

func (s fieldSet) Add(field Field) {
	if _, ok := s[field]; !ok {
		s[field] = none{}
	}
}

// --------------------------------------------------------------------

type fieldMap map[Field]*intsets.Dense

func (m fieldMap) Human() map[Field][]int {
	x := make(map[Field][]int, len(m))
	for f, s := range m {
		x[f] = s.AppendTo(nil)
	}
	return x
}

// InvertTo creates a copy of the field-map with the values inverted
func (m fieldMap) InvertTo(sz int) fieldMap {
	full := new(intsets.Dense)
	for i := 0; i < sz; i++ {
		full.Insert(i)
	}

	o := make(fieldMap, len(m))
	for k, v := range m {
		s := new(intsets.Dense)
		s.Copy(full)
		s.DifferenceWith(v)
		o[k] = s
	}
	return o
}

func (m fieldMap) Fetch(field Field) *intsets.Dense {
	v, ok := m[field]
	if !ok {
		v = new(intsets.Dense)
		m[field] = v
	}
	return v
}

// --------------------------------------------------------------------

type fieldValue struct {
	Field Field
	Value int
}

type fieldValueMap map[fieldValue]*intsets.Dense

func (m fieldValueMap) Fetch(field Field, value int) *intsets.Dense {
	key := fieldValue{Field: field, Value: value}
	v, ok := m[key]
	if !ok {
		v = new(intsets.Dense)
		m[key] = v
	}
	return v
}

func (m fieldValueMap) Copy() fieldValueMap {
	o := make(fieldValueMap, len(m))
	for k, v := range m {
		s := new(intsets.Dense)
		s.Copy(v)
		o[k] = s
	}
	return o
}

func (m fieldValueMap) Human() map[fieldValue][]int {
	x := make(map[fieldValue][]int, len(m))
	for f, s := range m {
		x[f] = s.AppendTo(nil)
	}
	return x
}

// --------------------------------------------------------------------

// outcome -> field -> multiOneOf
type fieldMultiMap map[int]map[Field]*multiOneOf

func (m fieldMultiMap) Human() map[int]map[Field]int {
	x := make(map[int]map[Field]int, len(m))
	for o, s := range m {
		if _, ok := x[o]; !ok {
			x[o] = make(map[Field]int)
		}
		for f, m := range s {
			x[o][f] = m.size
		}
	}
	return x
}

// Prune removes irrelevant checks
func (m fieldMultiMap) Prune() {
	for outcome, sub := range m {
		for field, check := range sub {
			if check.size < 2 {
				delete(sub, field)
				m[outcome] = sub
			}
		}
		if len(m[outcome]) == 0 {
			delete(m, outcome)
		}
	}
}

func (m fieldMultiMap) Fetch(outcome int, field Field) *multiOneOf {
	sub, ok := m[outcome]
	if !ok {
		sub = make(map[Field]*multiOneOf, 1)
		m[outcome] = sub
	}

	set, ok := sub[field]
	if !ok {
		set = newMultiOneOf()
		sub[field] = set
	}

	return set
}
