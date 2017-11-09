package qualify

import (
	"sync"

	"github.com/bsm/qualify/internal/intsets"
)

// Field identifies a particular fact field
type Field int

// --------------------------------------------------------------------

// Fact instances must implement a simple interface
type Fact interface {

	// AppendFieldValues must append of values for a given field to dst
	AppendFieldValues(dst []int, field Field) []int
}

// Qualifier instances can be used to determine qualified outcomes for
// any given fact.
type Qualifier struct {
	outcomes []int
	fields   []Field
	implicit fieldMap

	oneOf  fieldValueMap
	oneOfX fieldMultiMap
	noneOf fieldValueMap

	setPool, islPool sync.Pool
}

// Qualify will find outcomes matching the given fact and append
// them to dst
func (q *Qualifier) Qualify(dst []int, fact Fact) []int {
	if len(q.fields) == 0 {
		return dst
	}

	var oc *intsets.Dense // outcome candidates

	fc := q.recycleSet() // per-field candidates (from pool)
	defer q.setPool.Put(fc)

	vv := q.recycleSlice() // init scratch slice (from pool)
	defer q.islPool.Put(vv)

	for _, field := range q.fields {
		fc.Clear()

		// merge any implicit outcomes for this
		// field
		if set, ok := q.implicit[field]; ok {
			fc.UnionWith(set)
		}

		// retrieve fact values
		vv = fact.AppendFieldValues(vv[:0], field)

		// merge all explicit oneOf outcomes
		for _, val := range vv {
			fv := fieldValue{Field: field, Value: val}

			if set, ok := q.oneOf[fv]; ok {
				fc.UnionWith(set)
			}
		}

		// assign candidates
		if oc == nil {
			oc = q.recycleSet()
			defer q.setPool.Put(oc)
			oc.Copy(fc)
		} else {
			oc.IntersectionWith(fc)
		}

		// now, exclude all explicit noneOf outcomes
		for _, val := range vv {
			fv := fieldValue{Field: field, Value: val}

			if set, ok := q.noneOf[fv]; ok {
				oc.DifferenceWith(set)
			}
		}

		// finally, check if it's worth proceeding
		if oc.IsEmpty() {
			return dst
		}
	}

	// extract candidate positions (from pool)
	poss := oc.AppendTo(q.recycleSlice())
	defer q.islPool.Put(poss)

	// check each candidate against oneOfX restrictions
	for _, pos := range poss {
		outcome := q.outcomes[pos]

		skip := false
		for field, multi := range q.oneOfX[outcome] {
			vv = fact.AppendFieldValues(vv[:0], field)
			if !multi.Match(fc, vv...) {
				skip = true
				break
			}
		}

		if !skip {
			dst = append(dst, outcome)
		}
	}

	return dst
}

func (q *Qualifier) recycleSet() *intsets.Dense {
	if v := q.setPool.Get(); v != nil {
		s := v.(*intsets.Dense)
		return s
	}
	return new(intsets.Dense)
}

func (q *Qualifier) recycleSlice() []int {
	if v := q.islPool.Get(); v != nil {
		return v.([]int)[:0]
	}
	return nil
}
