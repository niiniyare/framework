package store

import (
	"fmt"
	"strings"

	"awo.so/framework/pkg/filter"
)

// SQLResult holds a parameterised SQL fragment and its arguments.
type SQLResult struct {
	SQL  string
	Args []any
}

// FilterToSQL translates a filter.Filter into a SQL WHERE clause and argument list.
// argOffset is the starting parameter index (use 1 for the first query, or pass the
// number of already-bound args + 1 to append to an existing query).
//
// Returns an empty SQLResult (no WHERE clause) for a zero-value Filter.
func FilterToSQL(f filter.Filter, argOffset int) (SQLResult, error) {
	parts := make([]string, 0, len(f.Conditions)+len(f.Groups))
	args := make([]any, 0)

	for _, c := range f.Conditions {
		frag, a, err := conditionToSQL(c, argOffset+len(args))
		if err != nil {
			return SQLResult{}, err
		}
		parts = append(parts, frag)
		args = append(args, a...)
	}

	for _, g := range f.Groups {
		frag, a, err := groupToSQL(g, argOffset+len(args))
		if err != nil {
			return SQLResult{}, err
		}
		parts = append(parts, frag)
		args = append(args, a...)
	}

	if len(parts) == 0 {
		return SQLResult{}, nil
	}

	return SQLResult{
		SQL:  strings.Join(parts, " AND "),
		Args: args,
	}, nil
}

func conditionToSQL(c filter.Condition, argN int) (string, []any, error) {
	col := columnExpr(c)

	switch c.Op {
	case filter.Eq:
		return fmt.Sprintf("%s = $%d", col, argN+1), []any{c.Value}, nil
	case filter.Neq:
		return fmt.Sprintf("%s != $%d", col, argN+1), []any{c.Value}, nil
	case filter.Gt:
		return fmt.Sprintf("%s > $%d", col, argN+1), []any{c.Value}, nil
	case filter.Gte:
		return fmt.Sprintf("%s >= $%d", col, argN+1), []any{c.Value}, nil
	case filter.Lt:
		return fmt.Sprintf("%s < $%d", col, argN+1), []any{c.Value}, nil
	case filter.Lte:
		return fmt.Sprintf("%s <= $%d", col, argN+1), []any{c.Value}, nil
	case filter.In:
		vals, ok := c.Value.([]any)
		if !ok {
			return "", nil, fmt.Errorf("filter: IN value must be []any, got %T", c.Value)
		}
		placeholders := make([]string, len(vals))
		for i := range vals {
			placeholders[i] = fmt.Sprintf("$%d", argN+i+1)
		}
		return fmt.Sprintf("%s = ANY(ARRAY[%s])", col, strings.Join(placeholders, ",")), vals, nil
	case filter.NotIn:
		vals, ok := c.Value.([]any)
		if !ok {
			return "", nil, fmt.Errorf("filter: NOT IN value must be []any, got %T", c.Value)
		}
		placeholders := make([]string, len(vals))
		for i := range vals {
			placeholders[i] = fmt.Sprintf("$%d", argN+i+1)
		}
		return fmt.Sprintf("%s != ALL(ARRAY[%s])", col, strings.Join(placeholders, ",")), vals, nil
	case filter.Contains:
		return fmt.Sprintf("%s LIKE $%d", col, argN+1), []any{fmt.Sprintf("%%%v%%", c.Value)}, nil
	case filter.StartsWith:
		return fmt.Sprintf("%s LIKE $%d", col, argN+1), []any{fmt.Sprintf("%v%%", c.Value)}, nil
	case filter.EndsWith:
		return fmt.Sprintf("%s LIKE $%d", col, argN+1), []any{fmt.Sprintf("%%%v", c.Value)}, nil
	case filter.ILike:
		return fmt.Sprintf("%s ILIKE $%d", col, argN+1), []any{c.Value}, nil
	case filter.IsNull:
		return fmt.Sprintf("%s IS NULL", col), nil, nil
	case filter.IsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", col), nil, nil
	default:
		return "", nil, fmt.Errorf("filter: unknown operator %q", c.Op)
	}
}

// columnExpr returns the SQL column expression for a condition.
// JSONB path conditions use the ->> operator to extract as text.
func columnExpr(c filter.Condition) string {
	if c.JSONPath != "" {
		// e.g. data->>'some_field'
		return fmt.Sprintf("data->>'%s'", strings.ReplaceAll(c.JSONPath, "'", "''"))
	}
	// Quoted identifier to prevent injection
	return fmt.Sprintf("%q", c.Field)
}

func groupToSQL(g filter.Group, argN int) (string, []any, error) {
	parts := make([]string, 0, len(g.Conditions)+len(g.Groups))
	args := make([]any, 0)

	for _, c := range g.Conditions {
		frag, a, err := conditionToSQL(c, argN+len(args))
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, frag)
		args = append(args, a...)
	}
	for _, sg := range g.Groups {
		frag, a, err := groupToSQL(sg, argN+len(args))
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, frag)
		args = append(args, a...)
	}

	if len(parts) == 0 {
		return "TRUE", nil, nil
	}

	sep := " AND "
	if g.Op == filter.LogicalOr {
		sep = " OR "
	}

	joined := strings.Join(parts, sep)
	if g.Op == filter.LogicalNot {
		return fmt.Sprintf("NOT (%s)", joined), args, nil
	}
	return fmt.Sprintf("(%s)", joined), args, nil
}

// OrderBySQL builds an ORDER BY clause from the filter's OrderBy slice.
func OrderBySQL(f filter.Filter) string {
	if len(f.OrderBy) == 0 {
		return "ORDER BY created_at DESC"
	}
	parts := make([]string, len(f.OrderBy))
	for i, o := range f.OrderBy {
		dir := "ASC"
		if o.Direction == filter.Desc {
			dir = "DESC"
		}
		nulls := ""
		if o.NullsLast {
			nulls = " NULLS LAST"
		}
		parts[i] = fmt.Sprintf("%q %s%s", o.Field, dir, nulls)
	}
	return "ORDER BY " + strings.Join(parts, ", ")
}
