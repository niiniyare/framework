// Package filter provides the query DSL used across the persistence layer and HTTP API.
package filter

// Operator is a comparison or membership operator.
type Operator string

const (
	// Comparison
	Eq    Operator = "eq"
	Neq   Operator = "neq"
	Gt    Operator = "gt"
	Gte   Operator = "gte"
	Lt    Operator = "lt"
	Lte   Operator = "lte"
	In    Operator = "in"
	NotIn Operator = "not_in"

	// String
	Contains   Operator = "contains"
	StartsWith Operator = "starts_with"
	EndsWith   Operator = "ends_with"
	ILike      Operator = "ilike"

	// Null
	IsNull    Operator = "is_null"
	IsNotNull Operator = "is_not_null"
)

// Condition is a single predicate: field op value.
type Condition struct {
	Field    string
	Op       Operator
	Value    any      // nil for IsNull / IsNotNull
	JSONPath string   // non-empty for JSONB field access (custom entity fields)
}

// LogicalOp combines conditions.
type LogicalOp string

const (
	LogicalAnd LogicalOp = "AND"
	LogicalOr  LogicalOp = "OR"
	LogicalNot LogicalOp = "NOT"
)

// Group is a set of Conditions combined with a logical operator.
type Group struct {
	Op         LogicalOp
	Conditions []Condition
	Groups     []Group // nested groups
}

// SortDirection is ascending or descending.
type SortDirection string

const (
	Asc  SortDirection = "ASC"
	Desc SortDirection = "DESC"
)

// OrderBy specifies a sort field and direction.
type OrderBy struct {
	Field     string
	Direction SortDirection
	NullsLast bool // default true in the framework
}

// PaginationMode selects the pagination strategy.
type PaginationMode string

const (
	PaginationCursor PaginationMode = "cursor" // default; stable under inserts
	PaginationOffset PaginationMode = "offset" // reports only
)

// Filter is the complete query specification passed to EntityRepository.FindMany.
type Filter struct {
	// Root-level conditions are ANDed together unless wrapped in a Group.
	Conditions []Condition
	Groups     []Group

	// Sorting
	OrderBy []OrderBy

	// Pagination
	Mode   PaginationMode
	Limit  int
	Offset int    // offset mode only
	Cursor string // cursor mode only; opaque to callers

	// Full-text search query (pg_trgm backed)
	Search string
}
