package intsets

const bitsPerWord = 64

// Dense is an intsets representation, optimised for 64-bits
type Dense struct {
	words []uint64
}

// Copy sets s to the value of x
func (s *Dense) Copy(x *Dense) {
	sz := len(x.words)
	s.ensure(sz)
	s.words = s.words[:sz]
	copy(s.words, x.words)
}

// AppendTo appends the entries to dst and returns the response
func (s *Dense) AppendTo(dst []int) []int {
	for pos, word := range s.words {
		if word == 0 {
			continue
		}

		v := pos * bitsPerWord
		for i := 0; word != 0 && i < bitsPerWord; i++ {
			if word&1 != 0 {
				dst = append(dst, v)
			}
			v++
			word >>= 1
		}
	}
	return dst
}

// IsEmpty reports whether the set s is empty.
func (s *Dense) IsEmpty() bool { return s.Len() == 0 }

// Len returns the number of elements in the set s.
func (s *Dense) Len() int { return popcntSlice(s.words) }

// Clear clears the set
func (s *Dense) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
	s.words = s.words[:0]
}

// Insert adds x to the set s
func (s *Dense) Insert(x int) {
	pos, i := x/bitsPerWord, uint(x%bitsPerWord)

	s.ensure(pos + 1)

	if mask := uint64(1) << i; s.words[pos]&mask == 0 {
		s.words[pos] |= mask
	}
}

// UnionWith sets s to the union s ∪ x, and reports whether s grew.
func (s *Dense) UnionWith(x *Dense) bool {
	changed := false

	s.ensure(len(x.words))
	for pos, wx := range x.words {
		if ws := s.words[pos]; wx != ws {
			s.words[pos] = ws | wx
			changed = true
		}
	}
	return changed
}

// IntersectionWith sets s to the intersection s ∩ x.
func (s *Dense) IntersectionWith(x *Dense) {
	sz, sx := len(s.words), len(x.words)

	if sx < sz {
		for pos := sx; pos < sz; pos++ {
			s.words[pos] = 0
		}
		s.words = s.words[:sx]
	}

	for pos, wx := range x.words {
		if pos >= sz {
			break
		}
		s.words[pos] &= wx
	}
}

// DifferenceWith sets s to the difference s ∖ x.
func (s *Dense) DifferenceWith(x *Dense) {
	if s == x {
		s.Clear()
		return
	}

	sz := len(s.words)
	for pos, wx := range x.words {
		if pos >= sz {
			break
		}
		s.words[pos] &= ^wx
	}
}

func (s *Dense) ensure(sz int) {
	cp := cap(s.words)
	if delta := sz - cp; delta > 0 {
		s.words = append(s.words[:cp], make([]uint64, delta)...)
	}

	if len(s.words) < sz {
		s.words = s.words[:sz]
	}
}
