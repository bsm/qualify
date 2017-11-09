# Qualify!

[![Build Status](https://travis-ci.org/bsm/qualify.png?branch=master)](https://travis-ci.org/bsm/qualify)
[![GoDoc](https://godoc.org/github.com/bsm/qualify?status.png)](http://godoc.org/github.com/bsm/qualify)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsm/qualify)](https://goreportcard.com/report/github.com/bsm/qualify)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Library for fast rules evaluation for Go. Qualify is able to quickly match a fact against large number of pre-defined rules.

### Example:

```go
import (
	"fmt"

	"github.com/bsm/qualify"
)

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

func main() {{ "Example_Qualifier" | code }}
```

