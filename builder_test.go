package qualify

import (
	. "github.com/bsm/ginkgo/v2"
	. "github.com/bsm/gomega"
)

var _ = Describe("Builder", func() {
	var subject *Builder

	BeforeEach(func() {
		subject = seedBuilder()
	})

	It("should add rules", func() {
		Expect(subject.outcomes).To(Equal(map[int]int{7: 0, 8: 1}))

		Expect(subject.fields).To(Equal(fieldSet{
			FieldCountry: {},
			FieldBrowser: {},
			FieldOS:      {},
			FieldAttrs:   {},
		}))
		Expect(subject.explicit.Human()).To(Equal(map[Field][]int{
			FieldCountry: {0},
			FieldBrowser: nil,
			FieldOS:      {1},
			FieldAttrs:   {0},
		}))
		Expect(subject.oneOf.Human()).To(Equal(map[fieldValue][]int{
			{Field: FieldCountry, Value: strDict.Add("gbr")}: {0},
			{Field: FieldCountry, Value: strDict.Add("irl")}: {0},
			{Field: FieldAttrs, Value: 33}:                   {0},
			{Field: FieldAttrs, Value: 34}:                   {0},
			{Field: FieldAttrs, Value: 55}:                   {0},
			{Field: FieldAttrs, Value: 56}:                   {0},
			{Field: FieldOS, Value: strDict.Add("android")}:  {1},
			{Field: FieldOS, Value: strDict.Add("ios")}:      {1},
		}))
		Expect(subject.oneOfX.Human()).To(Equal(map[int]map[Field]int{
			7: {FieldCountry: 1, FieldAttrs: 2},
			8: {FieldOS: 1},
		}))

		Expect(subject.noneOf.Human()).To(Equal(map[fieldValue][]int{
			{Field: FieldBrowser, Value: strDict.Add("safari")}: {0},
			{Field: FieldCountry, Value: strDict.Add("ger")}:    {1},
			{Field: FieldCountry, Value: strDict.Add("fra")}:    {1},
		}))
	})

	It("should compile qualifier", func() {
		qfy := subject.Compile()
		Expect(qfy.outcomes).To(Equal([]int{7, 8}))
		Expect(qfy.implicit.Human()).To(Equal(map[Field][]int{
			FieldCountry: {1},
			FieldBrowser: {0, 1},
			FieldOS:      {0},
			FieldAttrs:   {1},
		}))
		Expect(qfy.oneOf.Human()).To(Equal(map[fieldValue][]int{
			{Field: FieldCountry, Value: strDict.Add("gbr")}: {0},
			{Field: FieldCountry, Value: strDict.Add("irl")}: {0},
			{Field: FieldAttrs, Value: 33}:                   {0},
			{Field: FieldAttrs, Value: 34}:                   {0},
			{Field: FieldAttrs, Value: 55}:                   {0},
			{Field: FieldAttrs, Value: 56}:                   {0},
			{Field: FieldOS, Value: strDict.Add("android")}:  {1},
			{Field: FieldOS, Value: strDict.Add("ios")}:      {1},
		}))
		Expect(qfy.oneOfX.Human()).To(Equal(map[int]map[Field]int{
			7: {FieldAttrs: 2},
		}))
		Expect(qfy.noneOf.Human()).To(Equal(map[fieldValue][]int{
			{Field: FieldBrowser, Value: strDict.Add("safari")}: {0},
			{Field: FieldCountry, Value: strDict.Add("ger")}:    {1},
			{Field: FieldCountry, Value: strDict.Add("fra")}:    {1},
		}))
	})

	It("should snapshot on compile", func() {
		Expect(subject.Compile()).To(BeAssignableToTypeOf(&Qualifier{}))
		Expect(subject.Require(9, FieldAttrs, OneOf(77, 78))).To(BeFalse())
	})

})
