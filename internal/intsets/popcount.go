package intsets

// bit population count, take from
// https://code.google.com/p/go/issues/detail?id=4988#c11
// credit: https://code.google.com/u/arnehormann/
func popcntSliceGo(s []uint64) (n int) {
	for _, x := range s {
		x -= (x >> 1) & 0x5555555555555555
		x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
		x += x >> 4
		x &= 0x0f0f0f0f0f0f0f0f
		x *= 0x0101010101010101
		n += int(x >> 56)
	}
	return n
}
