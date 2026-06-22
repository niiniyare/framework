package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const sessionTTL = 24 * time.Hour

// sessionKey returns the Redis key for a session.
// Format: session:{tenantID}:{sessionID}
func sessionKey(tenantID, sessionID string) string {
	return fmt.Sprintf("session:%s:%s", tenantID, sessionID)
}

// SessionData is the payload stored per session.
type SessionData struct {
	UserID    string   `json:"user_id"`
	TenantID  string   `json:"tenant_id"`
	Roles     []string `json:"roles"`
	IsSuper   bool     `json:"is_super"`
	CreatedAt int64    `json:"created_at"` // unix seconds
}

// SetSession stores a session with the default TTL.
func (c *Client) SetSession(ctx context.Context, tenantID, sessionID string, data SessionData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("session marshal: %w", err)
	}
	return c.SetEX(ctx, sessionKey(tenantID, sessionID), string(b), sessionTTL)
}

// GetSession retrieves a session. Returns ErrKeyNotFound if absent or expired.
func (c *Client) GetSession(ctx context.Context, tenantID, sessionID string) (SessionData, error) {
	raw, err := c.Get(ctx, sessionKey(tenantID, sessionID))
	if err != nil {
		return SessionData{}, err
	}
	var d SessionData
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		return SessionData{}, fmt.Errorf("session unmarshal: %w", err)
	}
	return d, nil
}

// DeleteSession removes a session (logout).
func (c *Client) DeleteSession(ctx context.Context, tenantID, sessionID string) error {
	return c.Del(ctx, sessionKey(tenantID, sessionID))
}

// DeleteAllSessions removes all sessions for a user across all session IDs.
// This is a best-effort operation — it scans for matching keys.
// For high-volume tenants, consider storing session IDs in a separate set.
func (c *Client) DeleteAllSessions(ctx context.Context, tenantID, userID string) error {
	pattern := fmt.Sprintf("session:%s:*", tenantID)
	cmd := c.B().Scan().Cursor(0).Match(pattern).Count(100).Build()
	entry, err := c.Do(ctx, cmd).AsScanEntry()
	if err != nil {
		return fmt.Errorf("session scan: %w", err)
	}
	keys := entry.Elements

	toDelete := make([]string, 0, len(keys))
	for _, k := range keys {
		// Fetch and check user_id before deleting
		raw, err := c.Get(ctx, k)
		if err != nil {
			continue
		}
		var d SessionData
		if json.Unmarshal([]byte(raw), &d) == nil && d.UserID == userID {
			toDelete = append(toDelete, k)
		}
	}
	if len(toDelete) == 0 {
		return nil
	}
	return c.Del(ctx, toDelete...)
}
