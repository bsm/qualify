// +build !amd64 appengine

package intsets

func popcntSlice(s []uint64) int {
	return popcntSliceGo(s)
}
