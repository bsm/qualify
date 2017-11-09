// +build amd64,!appengine

package intsets

// go:noescape
func hasAsm() bool

// useAsm is a flag used to select the GO or ASM implementation of the popcnt function
var useAsm = hasAsm()

//go:noescape
func popcntSliceAsm(s []uint64) int

func popcntSlice(s []uint64) int {
	if useAsm {
		return popcntSliceAsm(s)
	}
	return popcntSliceGo(s)
}
