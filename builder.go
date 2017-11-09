package qualify

// Builder instances are used to map sets of rules to
// outcomes.
type Builder struct {
	outcomes map[int]int
	fields   fieldSet // fields with rules
	explicit fieldMap // explicitly targeted fields

	oneOf  fieldValueMap // one-of targeting
	oneOfX fieldMultiMap // one-of extra conditions
	noneOf fieldValueMap // none-of targeting

	readOnly bool
}

// NewBuilder inits a new builder
func NewBuilder() *Builder {
	return &Builder{
		outcomes: make(map[int]int),
		fields:   make(fieldSet),
		explicit: make(fieldMap),

		oneOf:  make(fieldValueMap),
		oneOfX: make(fieldMultiMap),
		noneOf: make(fieldValueMap),
	}
}

// Require adds a requirement condition for a particular outcome. Returns true if
// the requirement was acknowledged.
func (b *Builder) Require(outcome int, field Field, cond Condition) bool {
	if b.readOnly {
		return false
	}

	switch c := cond.(type) {
	case oneOf:
		b.requireOneOf(outcome, field, c)
	case noneOf:
		b.requireNoneOf(outcome, field, c)
	}
	return true
}

// Compile compiles a qualifier. The builder will become read-only as soon as this method is called.
func (b *Builder) Compile() *Qualifier {
	b.readOnly = true

	size := len(b.outcomes)

	// set outcomes
	outcomes := make([]int, size)
	for outcome, pos := range b.outcomes {
		outcomes[pos] = outcome
	}

	// set fields
	fields := make([]Field, 0, len(b.fields))
	for field := range b.fields {
		fields = append(fields, field)
	}

	// prune
	b.oneOfX.Prune()

	return &Qualifier{
		outcomes: outcomes,
		fields:   fields,
		implicit: b.explicit.InvertTo(size),

		oneOf:  b.oneOf,
		oneOfX: b.oneOfX,
		noneOf: b.noneOf,
	}
}

func (b *Builder) requireOneOf(outcome int, field Field, c oneOf) {
	if len(c) == 0 {
		return
	}

	pos := b.findPos(outcome)
	b.fields.Add(field)
	b.explicit.Fetch(field).Insert(pos)

	for _, val := range c {
		b.oneOf.Fetch(field, val).Insert(pos)
	}
	b.oneOfX.Fetch(outcome, field).OneOf(c...)
}

func (b *Builder) requireNoneOf(outcome int, field Field, c noneOf) {
	if len(c) == 0 {
		return
	}

	pos := b.findPos(outcome)
	b.fields.Add(field)
	b.explicit.Fetch(field)

	for _, val := range c {
		b.noneOf.Fetch(field, val).Insert(pos)
	}
}

func (b *Builder) findPos(outcome int) int {
	n, ok := b.outcomes[outcome]
	if !ok {
		n = len(b.outcomes)
		b.outcomes[outcome] = n
	}
	return n
}
