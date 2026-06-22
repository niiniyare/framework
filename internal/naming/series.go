// Package naming implements naming series generation for auto-numbered fields.
//
// Pattern syntax: "SE-{YYYY}-{MM}-{SEQ:5}"
// Tokens:
//   {YYYY}    — 4-digit year (posting date or today)
//   {YY}      — 2-digit year
//   {MM}      — 2-digit month
//   {DD}      — 2-digit day
//   {SEQ:N}   — zero-padded sequence, N digits wide (default 5)
//   {SEQ}     — same as {SEQ:5}
//
// Sequence counters are stored in naming_series_counters and incremented
// atomically via SELECT ... FOR UPDATE.
package naming

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var tokenRE = regexp.MustCompile(`\{([^}]+)\}`)

// Generator generates naming series values.
type Generator struct {
	db *pgxpool.Pool
}

// New creates a Generator.
func New(db *pgxpool.Pool) *Generator {
	return &Generator{db: db}
}

// Next generates the next value for the given series pattern and optional date.
// If date is zero, time.Now() is used.
func (g *Generator) Next(ctx context.Context, pattern string, date time.Time) (string, error) {
	if date.IsZero() {
		date = time.Now()
	}

	// Find the SEQ token to determine the counter key
	seqToken, seqWidth := extractSeqToken(pattern)
	if seqToken == "" {
		// No sequence — just expand date tokens
		return expandDate(pattern, date, 0), nil
	}

	// Build the counter key: pattern with date tokens expanded, SEQ replaced by literal
	counterKey := expandDate(strings.ReplaceAll(pattern, seqToken, "{SEQ}"), date, 0)

	// Atomically increment the counter
	seq, err := g.nextSeq(ctx, counterKey)
	if err != nil {
		return "", fmt.Errorf("naming series %q: %w", pattern, err)
	}

	result := expandDate(pattern, date, 0)
	result = strings.ReplaceAll(result, seqToken, fmt.Sprintf("%0*d", seqWidth, seq))
	return result, nil
}

// nextSeq increments the counter for key and returns the new value.
func (g *Generator) nextSeq(ctx context.Context, key string) (int64, error) {
	tx, err := g.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var seq int64
	err = tx.QueryRow(ctx, `
		INSERT INTO naming_series_counters (series_key, current_value)
		VALUES ($1, 1)
		ON CONFLICT (series_key) DO UPDATE
		    SET current_value = naming_series_counters.current_value + 1
		RETURNING current_value`, key).Scan(&seq)
	if err != nil {
		return 0, fmt.Errorf("counter upsert: %w", err)
	}

	return seq, tx.Commit(ctx)
}

// expandDate replaces date tokens in s with values from t.
// seqVal is only used for non-SEQ tokens (pass 0 when SEQ not yet known).
func expandDate(s string, t time.Time, _ int) string {
	return tokenRE.ReplaceAllStringFunc(s, func(m string) string {
		token := m[1 : len(m)-1] // strip { }
		switch token {
		case "YYYY":
			return fmt.Sprintf("%04d", t.Year())
		case "YY":
			return fmt.Sprintf("%02d", t.Year()%100)
		case "MM":
			return fmt.Sprintf("%02d", int(t.Month()))
		case "DD":
			return fmt.Sprintf("%02d", t.Day())
		default:
			// SEQ tokens — leave as-is for second pass
			return m
		}
	})
}

// extractSeqToken returns the full SEQ token (e.g. "{SEQ:5}") and its width.
func extractSeqToken(pattern string) (token string, width int) {
	matches := tokenRE.FindAllString(pattern, -1)
	for _, m := range matches {
		inner := m[1 : len(m)-1]
		if inner == "SEQ" {
			return m, 5
		}
		if strings.HasPrefix(inner, "SEQ:") {
			w, err := strconv.Atoi(strings.TrimPrefix(inner, "SEQ:"))
			if err != nil || w <= 0 {
				w = 5
			}
			return m, w
		}
	}
	return "", 0
}
