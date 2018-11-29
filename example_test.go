package qualify_test

import (
	"fmt"

	"github.com/bsm/qualify"
)

func ExampleQualifier() {
	// We can use dictionary encoding to translate
	// string values into numerics
	dict := qualify.NewStrDict()

	// Init a new builder popute with rules
	// for each outcome.
	builder := qualify.NewBuilder()

	// Outcome #34 requires:
	//  * Country to be GBR or FRA, and
	//  * Attrs to contain 101 or 102, and
	//  * Attrs to contain 202 or 203
	builder.Require(34, FieldCountry,
		qualify.OneOf(dict.Add("GBR"), dict.Add("FRA")))
	builder.Require(34, FieldAttrs,
		qualify.OneOf(101, 102))
	builder.Require(34, FieldAttrs,
		qualify.OneOf(202, 203))

	// Outcome #35 requires:
	//  * Country NOT to be NLD nor GER, and
	//  * Browser NOT to be Safari, and
	//  * OS to be either Android or iOS
	builder.Require(35, FieldCountry,
		qualify.NoneOf(dict.Add("NLD"), dict.Add("GER")))
	builder.Require(35, FieldBrowser,
		qualify.NoneOf(dict.Add("Safari")))
	builder.Require(35, FieldOS,
		qualify.OneOf(dict.Add("Android"), dict.Add("iOS")))

	// Setup the qualifier
	qfy := builder.Compile()

	// Init result set
	var res []int

	// Matches outcome #34
	res = qfy.Qualify(res[:0], &factReader{Dict: dict, Fact: Fact{Country: "GBR", Attrs: []int{101, 202}}})
	fmt.Println(res)

	// Matches outcome #35
	res = qfy.Qualify(res[:0], &factReader{Dict: dict, Fact: Fact{Country: "IRE", OS: "iOS"}})
	fmt.Println(res)

	// Matches nothing
	res = qfy.Qualify(res[:0], &factReader{Dict: dict, Fact: Fact{Country: "NLD"}})
	fmt.Println(res)

	// Output:
	// [34]
	// [35]
	// []
}

// --------------------------------------------------------------------

// Fact is an example fact
type Fact struct {
	Country string
	Browser string
	OS      string
	Attrs   []int
}

// Enumeration of our fact features
const (
	FieldCountry qualify.Field = iota
	FieldBrowser
	FieldOS
	FieldAttrs
)

// factReader is a wrapper around facts to
// make them comply with qualify.Fact
type factReader struct {
	Dict qualify.StrDict
	Fact
}

func (r *factReader) AppendFieldValues(x []int, f qualify.Field) []int {
	switch f {
	case FieldCountry:
		return r.Dict.Lookup(x, r.Country)
	case FieldBrowser:
		return r.Dict.Lookup(x, r.Browser)
	case FieldOS:
		return r.Dict.Lookup(x, r.OS)
	case FieldAttrs:
		return append(x, r.Attrs...)
	}
	return x
}
