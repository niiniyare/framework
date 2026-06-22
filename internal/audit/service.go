// Package audit provides append-only audit event recording.
// Writes to audit_events (system schema) and domain_events outbox (async delivery).
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// EventType constants for well-known IAM events.
const (
	EventUserCreated     = "USER_CREATED"
	EventUserSuspended   = "USER_SUSPENDED"
	EventUserDeactivated = "USER_DEACTIVATED"
	EventSessionStarted  = "SESSION_STARTED"
	EventSessionEnded    = "SESSION_ENDED"
	EventLoginFailed     = "LOGIN_FAILED"
	EventMFAEnabled      = "MFA_ENABLED"
	EventMFADisabled     = "MFA_DISABLED"
	EventPasswordChanged = "PASSWORD_CHANGED"
	EventRoleAssigned    = "ROLE_ASSIGNED"
	EventRoleRevoked     = "ROLE_REVOKED"
	EventPolicyAdded     = "POLICY_ADDED"
	EventPolicyRemoved   = "POLICY_REMOVED"
	EventFlagChanged     = "FLAG_CHANGED"
	EventSettingChanged  = "SETTING_CHANGED"
	EventAccessDenied    = "ACCESS_DENIED"
	EventAccessGranted   = "ACCESS_GRANTED"

	// Tenant lifecycle
	EventTenantProvisioned  = "TENANT_PROVISIONED"
	EventTenantSuspended    = "TENANT_SUSPENDED"
	EventTenantReactivated  = "TENANT_REACTIVATED"
)

// Event represents a single auditable occurrence.
type Event struct {
	ID            uuid.UUID
	TenantID      *uuid.UUID
	EventType     string
	ActorSubject  string
	ActorDomain   string
	TargetSubject string
	Resource      string
	Action        string
	Outcome       string // "allowed" | "denied"
	Metadata      map[string]any
	CreatedAt     time.Time
}

// DomainEvent is an outbox event for async delivery.
type DomainEvent struct {
	EventType  string
	TenantID   *uuid.UUID
	Payload    map[string]any
	OccurredAt time.Time
}

// Filter for querying audit events.
type Filter struct {
	TenantID  *uuid.UUID
	EventType string
	Since     *time.Time
	Until     *time.Time
	Limit     int
}

// Service writes and queries audit events.
type Service struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

// New creates an audit Service.
func New(db *pgxpool.Pool, log *slog.Logger) *Service {
	return &Service{db: db, log: log}
}

// Log records an audit event. Non-blocking — errors are logged, not returned,
// so audit failures never block the main request path.
func (s *Service) Log(ctx context.Context, e Event) {
	if e.ID == (uuid.UUID{}) {
		e.ID = uuid.New()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	meta, _ := json.Marshal(e.Metadata)

	_, err := s.db.Exec(ctx, `
		INSERT INTO audit_events
		    (id, tenant_id, event_type, actor_subject, actor_domain,
		     target_subject, resource, action, outcome, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		e.ID, e.TenantID, e.EventType, e.ActorSubject, e.ActorDomain,
		e.TargetSubject, e.Resource, e.Action, e.Outcome, meta, e.CreatedAt)
	if err != nil {
		s.log.Error("audit: write failed",
			slog.String("event_type", e.EventType),
			slog.Any("err", err))
	}
}

// Emit writes a domain event to the outbox for async delivery.
// Must be called within the same DB transaction as the triggering operation
// for at-least-once delivery guarantees.
func (s *Service) Emit(ctx context.Context, e DomainEvent) error {
	if e.OccurredAt.IsZero() {
		e.OccurredAt = time.Now()
	}
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return fmt.Errorf("audit emit: marshal payload: %w", err)
	}
	_, err = s.db.Exec(ctx, `
		INSERT INTO domain_events (tenant_id, event_type, payload, occurred_at)
		VALUES ($1, $2, $3, $4)`,
		e.TenantID, e.EventType, payload, e.OccurredAt)
	return err
}

// Query retrieves audit events matching the filter.
func (s *Service) Query(ctx context.Context, f Filter) ([]*Event, error) {
	query := `
		SELECT id, tenant_id, event_type, actor_subject, actor_domain,
		       COALESCE(target_subject,''), COALESCE(resource,''), COALESCE(action,''),
		       COALESCE(outcome,''), COALESCE(metadata,'{}'), created_at
		FROM audit_events WHERE 1=1`
	args := []any{}

	if f.TenantID != nil {
		args = append(args, *f.TenantID)
		query += fmt.Sprintf(" AND tenant_id = $%d", len(args))
	}
	if f.EventType != "" {
		args = append(args, f.EventType)
		query += fmt.Sprintf(" AND event_type = $%d", len(args))
	}
	if f.Since != nil {
		args = append(args, *f.Since)
		query += fmt.Sprintf(" AND created_at >= $%d", len(args))
	}
	if f.Until != nil {
		args = append(args, *f.Until)
		query += fmt.Sprintf(" AND created_at <= $%d", len(args))
	}
	query += " ORDER BY created_at DESC"
	if f.Limit > 0 {
		args = append(args, f.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("audit query: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		var meta []byte
		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.EventType, &e.ActorSubject, &e.ActorDomain,
			&e.TargetSubject, &e.Resource, &e.Action, &e.Outcome, &meta, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("audit query scan: %w", err)
		}
		_ = json.Unmarshal(meta, &e.Metadata)
		events = append(events, &e)
	}
	return events, rows.Err()
}
