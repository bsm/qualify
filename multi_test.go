package qualify

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("multiOneOf", func() {
	var subject *multiOneOf

	BeforeEach(func() {
		subject = newMultiOneOf()
		subject.OneOf(1, 2, 3)
		subject.OneOf(3, 4, 5)
	})

	It("should match", func() {
		Expect(subject.Match(nil, 3, 5)).To(BeTrue())
		Expect(subject.Match(nil, 1, 4)).To(BeTrue())
		Expect(subject.Match(nil, 1, 2)).To(BeFalse())
		Expect(subject.Match(nil, 6, 4)).To(BeFalse())
		Expect(subject.Match(nil, 7, 3)).To(BeTrue())

		Expect(subject.Match(nil, 3)).To(BeTrue())
		Expect(subject.Match(nil, 1, 2, 6, 7, 8, 9)).To(BeFalse())
		Expect(subject.Match(nil, 1, 2, 6, 7, 8, 5)).To(BeTrue())
	})

})
