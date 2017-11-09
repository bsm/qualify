package qualify

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = DescribeTable("Qualifier",
	func(fact *mockFact, m types.GomegaMatcher) {
		q := seedBuilder().Compile()
		Expect(q.Qualify(nil, fact)).To(m)
	},

	Entry("match: simple", &mockFact{
		Country: "gbr",
		Attrs:   []int{33, 55},
	}, ConsistOf(7)),

	Entry("match: combo", &mockFact{
		Country: "gbr",
		OS:      "ios",
	}, ConsistOf(8)),

	Entry("match: multi", &mockFact{
		Country: "gbr",
		Attrs:   []int{33, 34, 55},
		OS:      "ios",
	}, ConsistOf(7, 8)),

	Entry("no match: one-of missing", &mockFact{
		Country: "usa",
	}, BeEmpty()),

	Entry("no match: one-of partially missing", &mockFact{
		Country: "gbr",
		Attrs:   []int{33, 34},
	}, BeEmpty()),

	Entry("no match: none-of violation", &mockFact{
		Country: "fra",
		Attrs:   []int{33, 55},
		OS:      "android",
	}, BeEmpty()),
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	BeforeEach(func() {
		strDict = make(StrDict)
	})
	RunSpecs(t, "qualify")
}

// --------------------------------------------------------------------

const (
	FieldCountry Field = iota
	FieldBrowser
	FieldOS
	FieldAttrs
)

var strDict StrDict

type mockFact struct {
	Country string
	Browser string
	OS      string
	Attrs   []int
}

func (m *mockFact) AppendFieldValues(x []int, f Field) []int {
	switch f {
	case FieldCountry:
		return strDict.Lookup(x, m.Country)
	case FieldBrowser:
		return strDict.Lookup(x, m.Browser)
	case FieldOS:
		return strDict.Lookup(x, m.OS)
	case FieldAttrs:
		return append(x, m.Attrs...)
	}
	return x
}

func seedBuilder() *Builder {
	b := NewBuilder()

	b.Require(7, FieldCountry, OneOf(strDict.Add("gbr"), strDict.Add("irl")))
	b.Require(7, FieldBrowser, NoneOf(strDict.Add("safari")))
	b.Require(7, FieldAttrs, OneOf(33, 34))
	b.Require(7, FieldAttrs, OneOf(55, 56))

	b.Require(8, FieldCountry, NoneOf(strDict.Add("ger"), strDict.Add("fra")))
	b.Require(8, FieldOS, OneOf(strDict.Add("android"), strDict.Add("ios")))

	return b
}
