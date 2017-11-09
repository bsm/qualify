package qualify_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bsm/qualify"
)

func Benchmark_Qualify(b *testing.B) {
	qfy, err := benchCompile()
	if err != nil {
		b.Fatal(err)
	}

	facts, err := benchLoadFacts()
	if err != nil {
		b.Fatal(err)
	}

	var tmp []int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fct := facts[i%len(facts)]
		tmp = qfy.Qualify(tmp[:0], &fct)
	}
}

func benchCompile() (*qualify.Qualifier, error) {
	f, err := os.Open("testdata/outcomes.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	builder := qualify.NewBuilder()
	for dec := json.NewDecoder(f); dec.More(); {
		var o benchOutcome

		if err := dec.Decode(&o); err != nil {
			return nil, err
		}
		for _, rule := range o.Rules {
			builder.Require(o.Outcome, rule.Field(), rule.Condition())
		}
	}

	return builder.Compile(), nil
}

func benchLoadFacts() ([]benchFact, error) {
	f, err := os.Open("testdata/facts.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var res []benchFact
	for dec := json.NewDecoder(f); dec.More(); {
		var fct benchFact

		if err := dec.Decode(&fct); err != nil {
			return nil, err
		}
		res = append(res, fct)
	}
	return res, nil
}

// --------------------------------------------------------------------

const (
	fieldUnknown qualify.Field = iota
	fieldA
	fieldB
	fieldC
	fieldD
	fieldE
	fieldF
	fieldG
	fieldH
	fieldI
	fieldJ
	fieldK
	fieldL
	fieldM
	fieldN
	fieldO
	fieldP
)

type benchOutcome struct {
	Outcome int
	Rules   []benchRule
}

type benchRule struct {
	Key  string `json:"k"`
	Type int    `json:"t"`
	Vals []int  `json:"v"`
}

func (r benchRule) Condition() qualify.Condition {
	if r.Type == 2 {
		return qualify.NoneOf(r.Vals...)
	}
	return qualify.OneOf(r.Vals...)
}

func (r benchRule) Field() qualify.Field {
	switch r.Key {
	case "a":
		return fieldA
	case "b":
		return fieldB
	case "c":
		return fieldC
	case "d":
		return fieldD
	case "e":
		return fieldE
	case "f":
		return fieldF
	case "g":
		return fieldG
	case "h":
		return fieldH
	case "i":
		return fieldI
	case "j":
		return fieldJ
	case "k":
		return fieldK
	case "l":
		return fieldL
	case "m":
		return fieldM
	case "n":
		return fieldN
	case "o":
		return fieldO
	case "p":
		return fieldP
	}
	return fieldUnknown
}

type benchFact struct {
	A int
	B []int
	C int
	D []int
	E int
	F int
	G int
	H []int
	I int
	J int
	K int
	L int
	M int
	N int
	O int
	P int
}

func (f *benchFact) AppendFieldValues(dst []int, field qualify.Field) []int {
	switch field {
	case fieldA:
		return append(dst, f.A)
	case fieldB:
		return append(dst, f.B...)
	case fieldC:
		return append(dst, f.C)
	case fieldD:
		return append(dst, f.D...)
	case fieldE:
		return append(dst, f.E)
	case fieldF:
		return append(dst, f.F)
	case fieldG:
		return append(dst, f.G)
	case fieldH:
		return append(dst, f.H...)
	case fieldI:
		return append(dst, f.I)
	case fieldJ:
		return append(dst, f.J)
	case fieldK:
		return append(dst, f.K)
	case fieldL:
		return append(dst, f.L)
	case fieldM:
		return append(dst, f.M)
	case fieldN:
		return append(dst, f.N)
	case fieldO:
		return append(dst, f.O)
	case fieldP:
		return append(dst, f.P)
	}
	return dst
}
