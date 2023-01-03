package qualify

import (
	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
)

var _ = Describe("StrDict", func() {
	var subject StrDict

	BeforeEach(func() {
		subject = make(StrDict)
	})

	It("should add", func() {
		Expect(subject.Add("a")).To(Equal(0))
		Expect(subject.Add("b")).To(Equal(1))
		Expect(subject.Add("a")).To(Equal(0))
		Expect(subject.Add("c")).To(Equal(2))
		Expect(subject.Add("c")).To(Equal(2))
		Expect(subject.Add("b")).To(Equal(1))
		Expect(subject.Add("d")).To(Equal(3))
	})

	It("should lookup", func() {
		Expect(subject.Add("a")).To(Equal(0))
		Expect(subject.Add("b")).To(Equal(1))
		Expect(subject.Add("c")).To(Equal(2))
		Expect(subject.Lookup(nil, "a", "b")).To(Equal([]int{0, 1}))
		Expect(subject.Lookup(nil, "c", "b", "d", "a")).To(Equal([]int{2, 1, 0}))
	})

})
