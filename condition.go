package qualify

// Condition represents a rule requirement
type Condition interface {
	_isQualifyCondition()
}

// --------------------------------------------------------------------

type oneOf []int

// OneOf returns a condition that
// requires any of the values to be present
func OneOf(values ...int) Condition {
	return oneOf(values)
}

func (oneOf) _isQualifyCondition() {}

// --------------------------------------------------------------------

type noneOf []int

// NoneOf returns a condition that
// requires none of the values to be present
func NoneOf(values ...int) Condition {
	return noneOf(values)
}

func (noneOf) _isQualifyCondition() {}
