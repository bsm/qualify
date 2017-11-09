package intsets

import (
	"sort"
)

// Naive is a naive intsets representation
type Naive struct {
	vals []int
}

// Copy sets s to the value of x
func (s *Naive) Copy(x *Naive) {
	n := len(x.vals)
	s.ensure(n)
	s.vals = s.vals[:n]
	copy(s.vals, x.vals)
}

// AppendTo appends the entries to dst and returns the response
func (s *Naive) AppendTo(dst []int) []int {
	return append(dst, s.vals...)
}

// IsEmpty reports whether the set s is empty.
func (s *Naive) IsEmpty() bool { return s.Len() == 0 }

// Len returns the number of elements in the set s.
func (s *Naive) Len() int { return len(s.vals) }

// Clear clears the set
func (s *Naive) Clear() {
	s.vals = s.vals[:0]
}

// Insert adds x to the set s
func (s *Naive) Insert(x int) bool {
	if pos := sort.SearchInts(s.vals, x); pos < len(s.vals) {
		if s.vals[pos] == x {
			return false
		}

		s.vals = append(s.vals, 0)
		copy(s.vals[pos+1:], s.vals[pos:])
		s.vals[pos] = x
	} else {
		s.vals = append(s.vals, x)
	}
	return true
}

// UnionWith sets s to the union s ∪ x, and reports whether s grew.
func (s *Naive) UnionWith(x *Naive) bool {
	if s == x {
		return false
	}

	on := len(s.vals)
	s.vals = append(s.vals, x.vals...)
	sn := len(s.vals)
	if on == sn {
		return false
	}

	sort.Ints(s.vals)

	p := 0
	for i := 1; i < sn; i++ {
		if s.vals[p] >= s.vals[i] {
			continue
		}
		p++
		if p < i {
			s.vals[p], s.vals[i] = s.vals[i], s.vals[p]
		}
	}
	s.vals = s.vals[:p+1]
	return p+1 > on
}

// IntersectionWith sets s to the intersection s ∩ x.
func (s *Naive) IntersectionWith(x *Naive) {
	if s == x {
		return
	}

	p, i, j := 0, 0, 0
	sn, xn := len(s.vals), len(x.vals)

	for i < sn && j < xn {
		switch {
		case s.vals[i] < x.vals[j]:
			i++
		case s.vals[i] > x.vals[j]:
			j++
		case p < i:
			s.vals[p], s.vals[i] = s.vals[i], s.vals[p]
			fallthrough
		default:
			p, i, j = p+1, i+1, j+1
		}
	}
	s.vals = s.vals[:p]
}

// DifferenceWith sets s to the difference s ∖ x.
func (s *Naive) DifferenceWith(x *Naive) {
	if s == x {
		s.Clear()
		return
	}

	p, i, j := 0, 0, 0
	sn, xn := len(s.vals), len(x.vals)

	for i < sn && j < xn {
		switch {
		case s.vals[i] < x.vals[j]:
			if p < i {
				s.vals[p], s.vals[i] = s.vals[i], s.vals[p]
			}
			p, i = p+1, i+1
		case s.vals[i] > x.vals[j]:
			j++
		default:
			i, j = i+1, j+1
		}
	}

	for p < sn && i < sn {
		s.vals[p], s.vals[i] = s.vals[i], s.vals[p]
		p, i = p+1, i+1
	}
	s.vals = s.vals[:p]
}

func (s *Naive) ensure(n int) {
	cp := cap(s.vals)
	if delta := n - cp; delta > 0 {
		s.vals = append(s.vals[:cp], make([]int, delta)...)
	}

	if len(s.vals) < n {
		s.vals = s.vals[:n]
	}
}
