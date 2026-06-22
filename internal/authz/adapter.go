package authz

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxAdapter is a Casbin persist.Adapter backed by pgxpool.
type pgxAdapter struct {
	db *pgxpool.Pool
}

// NewPgxAdapter returns a Casbin adapter using the provided pgxpool.
func NewPgxAdapter(db *pgxpool.Pool) persist.Adapter {
	return &pgxAdapter{db: db}
}

const insertRuleSQL = `
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT DO NOTHING`

func (a *pgxAdapter) LoadPolicy(m model.Model) error {
	ctx := context.Background()
	rows, err := a.db.Query(ctx, `SELECT ptype, v0, v1, v2, v3, v4, v5 FROM casbin_rule`)
	if err != nil {
		return fmt.Errorf("casbin load: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ptype, v0, v1, v2, v3, v4, v5 string
		if err := rows.Scan(&ptype, &v0, &v1, &v2, &v3, &v4, &v5); err != nil {
			return fmt.Errorf("casbin load scan: %w", err)
		}
		persist.LoadPolicyLine(buildPolicyLine(ptype, v0, v1, v2, v3, v4, v5), m)
	}
	return rows.Err()
}

func (a *pgxAdapter) SavePolicy(m model.Model) error {
	ctx := context.Background()
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("casbin save begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `DELETE FROM casbin_rule`); err != nil {
		return fmt.Errorf("casbin save truncate: %w", err)
	}

	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			if err := execInsert(ctx, tx, ptype, rule); err != nil {
				return err
			}
		}
	}
	for ptype, ast := range m["g"] {
		for _, rule := range ast.Policy {
			if err := execInsert(ctx, tx, ptype, rule); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func (a *pgxAdapter) AddPolicy(_ string, ptype string, rule []string) error {
	ctx := context.Background()
	_, err := a.db.Exec(ctx, insertRuleSQL,
		ptype, strAt(rule, 0), strAt(rule, 1), strAt(rule, 2),
		strAt(rule, 3), strAt(rule, 4), strAt(rule, 5))
	return err
}

func (a *pgxAdapter) RemovePolicy(_ string, ptype string, rule []string) error {
	ctx := context.Background()
	_, err := a.db.Exec(ctx,
		`DELETE FROM casbin_rule
         WHERE ptype=$1 AND v0=$2 AND v1=$3 AND v2=$4 AND v3=$5 AND v4=$6 AND v5=$7`,
		ptype, strAt(rule, 0), strAt(rule, 1), strAt(rule, 2),
		strAt(rule, 3), strAt(rule, 4), strAt(rule, 5))
	return err
}

func (a *pgxAdapter) RemoveFilteredPolicy(_ string, ptype string, fieldIndex int, fieldValues ...string) error {
	ctx := context.Background()
	cols := []string{"v0", "v1", "v2", "v3", "v4", "v5"}
	query := `DELETE FROM casbin_rule WHERE ptype = $1`
	args := []any{ptype}
	for i, val := range fieldValues {
		if val != "" {
			query += fmt.Sprintf(" AND %s = $%d", cols[fieldIndex+i], len(args)+1)
			args = append(args, val)
		}
	}
	_, err := a.db.Exec(ctx, query, args...)
	return err
}

func execInsert(ctx context.Context, tx pgx.Tx, ptype string, rule []string) error {
	_, err := tx.Exec(ctx, insertRuleSQL,
		ptype, strAt(rule, 0), strAt(rule, 1), strAt(rule, 2),
		strAt(rule, 3), strAt(rule, 4), strAt(rule, 5))
	return err
}

func strAt(s []string, i int) string {
	if i < len(s) {
		return s[i]
	}
	return ""
}

func buildPolicyLine(parts ...string) string {
	line := ""
	for i, p := range parts {
		if i == 0 {
			line = p
		} else if p != "" {
			line += ", " + p
		}
	}
	return line
}
