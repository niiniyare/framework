package filter

// Builder constructs a Filter via a fluent API.
type Builder struct {
	f Filter
}

// New returns a new Builder with defaults applied.
func New() *Builder {
	return &Builder{f: Filter{
		Mode:  PaginationCursor,
		Limit: 20,
	}}
}

// Where adds a simple equality condition (shorthand for Eq).
func (b *Builder) Where(field string, value any) *Builder {
	b.f.Conditions = append(b.f.Conditions, Condition{Field: field, Op: Eq, Value: value})
	return b
}

// Condition adds an arbitrary condition.
func (b *Builder) Condition(field string, op Operator, value any) *Builder {
	b.f.Conditions = append(b.f.Conditions, Condition{Field: field, Op: op, Value: value})
	return b
}

// JSONBCondition adds a condition on a JSONB path (for custom entity fields).
func (b *Builder) JSONBCondition(jsonPath string, op Operator, value any) *Builder {
	b.f.Conditions = append(b.f.Conditions, Condition{
		Field:    jsonPath,
		JSONPath: jsonPath,
		Op:       op,
		Value:    value,
	})
	return b
}

// And wraps a set of conditions in an AND group.
func (b *Builder) And(conditions ...Condition) *Builder {
	b.f.Groups = append(b.f.Groups, Group{Op: LogicalAnd, Conditions: conditions})
	return b
}

// Or wraps a set of conditions in an OR group.
func (b *Builder) Or(conditions ...Condition) *Builder {
	b.f.Groups = append(b.f.Groups, Group{Op: LogicalOr, Conditions: conditions})
	return b
}

// Search sets the full-text search query.
func (b *Builder) Search(q string) *Builder {
	b.f.Search = q
	return b
}

// OrderBy appends a sort clause.
func (b *Builder) OrderBy(field string, dir SortDirection) *Builder {
	b.f.OrderBy = append(b.f.OrderBy, OrderBy{Field: field, Direction: dir, NullsLast: true})
	return b
}

// Limit sets the page size.
func (b *Builder) Limit(n int) *Builder {
	b.f.Limit = n
	return b
}

// Cursor sets the opaque pagination cursor (cursor mode).
func (b *Builder) Cursor(cursor string) *Builder {
	b.f.Mode = PaginationCursor
	b.f.Cursor = cursor
	return b
}

// Offset switches to offset pagination and sets the offset.
func (b *Builder) Offset(n int) *Builder {
	b.f.Mode = PaginationOffset
	b.f.Offset = n
	return b
}

// Build returns the constructed Filter.
func (b *Builder) Build() Filter {
	return b.f
}

// C is a convenience constructor for a single Condition.
func C(field string, op Operator, value any) Condition {
	return Condition{Field: field, Op: op, Value: value}
}
