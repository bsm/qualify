package intsets

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive", func() {
	var subject *Naive

	var seed = func(seed int64, n, max int) *Naive {
		set := &Naive{}
		rnd := rand.New(rand.NewSource(seed))
		for i := 0; i < n; i++ {
			set.Insert(rnd.Intn(max))
		}
		return set
	}

	s1 := seed(11, 24, 1080)
	s2 := seed(12, 12, 2060)
	s3 := seed(13, 60, 110)

	BeforeEach(func() {
		subject = &Naive{}
		subject.Copy(s1)
	})

	It("should insert", func() {
		Expect(s1.vals).To(HaveLen(24))
		Expect(s1.vals).To(HaveCap(32))
	})

	It("should copy", func() {
		Expect(subject.vals).To(HaveLen(24))
		Expect(subject.vals).To(HaveCap(24))

		subject.Insert(500)
		Expect(subject.Len()).To(Equal(25))
		Expect(s1.Len()).To(Equal(24))
	})

	It("should append", func() {
		Expect(s1.AppendTo(nil)).To(Equal([]int{
			25, 49, 109, 133, 156, 196, 205, 254,
			389, 393, 413, 430, 551, 558, 563, 629,
			664, 840, 886, 950, 959, 995, 1014, 1060,
		}))

		Expect(s2.AppendTo(nil)).To(Equal([]int{
			178, 535, 1009, 1123, 1126, 1174,
			1383, 1510, 1652, 1843, 1868, 1893,
		}))
	})

	It("should union", func() {
		Expect(subject.Len()).To(Equal(24))
		Expect(subject.UnionWith(subject)).To(BeFalse())
		Expect(subject.Len()).To(Equal(24))

		Expect(subject.UnionWith(s2)).To(BeTrue())
		Expect(subject.Len()).To(Equal(36))
		Expect(s2.Len()).To(Equal(12))

		Expect(subject.AppendTo(nil)).To(Equal([]int{
			25, 49, 109, 133, 156, 178, 196, 205, 254,
			389, 393, 413, 430, 535, 551, 558, 563, 629,
			664, 840, 886, 950, 959, 995, 1009, 1014, 1060,
			1123, 1126, 1174, 1383, 1510, 1652, 1843, 1868, 1893,
		}))

		Expect(subject.UnionWith(s3)).To(BeTrue())
		Expect(subject.Len()).To(Equal(79))
		Expect(s3.Len()).To(Equal(45))
	})

	It("should intersect", func() {
		Expect(subject.Len()).To(Equal(24))
		subject.IntersectionWith(subject)
		Expect(subject.Len()).To(Equal(24))

		subject.IntersectionWith(s3)
		Expect(subject.Len()).To(Equal(2))
		Expect(s3.Len()).To(Equal(45))

		Expect(subject.AppendTo(nil)).To(Equal([]int{25, 109}))
	})

	It("should difference", func() {
		Expect(subject.Len()).To(Equal(24))
		subject.DifferenceWith(s3)
		Expect(subject.Len()).To(Equal(22))
		Expect(s3.Len()).To(Equal(45))
		Expect(subject.AppendTo(nil)).To(Equal([]int{
			49, 133, 156, 196, 205, 254,
			389, 393, 413, 430, 551, 558, 563, 629,
			664, 840, 886, 950, 959, 995, 1014, 1060,
		}))
		subject.DifferenceWith(subject)
		Expect(subject.Len()).To(Equal(0))
	})

	It("should have len", func() {
		Expect(s1.Len()).To(Equal(24))
		Expect(s2.Len()).To(Equal(12))
		Expect(s3.Len()).To(Equal(45))
		Expect(subject.IsEmpty()).To(BeFalse())
	})

	It("should clear", func() {
		subject.Clear()
		Expect(subject.Len()).To(Equal(0))
		Expect(subject.IsEmpty()).To(BeTrue())
	})

})
