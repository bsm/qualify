package qualify

import (
	"github.com/bsm/qualify/internal/intsets"
)

type multiOneOf struct {
	reqs map[int]*intsets.Dense
	size int
}

func newMultiOneOf() *multiOneOf {
	return &multiOneOf{
		reqs: make(map[int]*intsets.Dense, 1),
	}
}

func (m *multiOneOf) OneOf(vals ...int) {
	if len(vals) == 0 {
		return
	}

	for _, v := range vals {
		s, ok := m.reqs[v]
		if !ok {
			s = new(intsets.Dense)
			m.reqs[v] = s
		}
		s.Insert(m.size)
	}
	m.size++
}

func (m *multiOneOf) Match(tmp *intsets.Dense, vals ...int) bool {
	if len(vals) == 0 {
		return false
	}

	if tmp == nil {
		tmp = new(intsets.Dense)
	} else {
		tmp.Clear()
	}

	for _, v := range vals {
		if s, ok := m.reqs[v]; ok {
			if tmp.UnionWith(s) && tmp.Len() == m.size {
				return true
			}
		}
	}
	return false
}
