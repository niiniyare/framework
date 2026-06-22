package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const flagTTL = 5 * time.Minute

// flagKey returns the Redis key for a feature flag value.
// Format: flags:{tenantID}:{flagName}
func flagKey(tenantID, flagName string) string {
	return fmt.Sprintf("flags:%s:%s", tenantID, flagName)
}

// GetFlag reads a boolean feature flag for a tenant.
// Returns (false, ErrKeyNotFound) if not cached — callers should fall back to DB.
func (c *Client) GetFlag(ctx context.Context, tenantID, flagName string) (bool, error) {
	raw, err := c.Get(ctx, flagKey(tenantID, flagName))
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(raw)
}

// SetFlag caches a boolean feature flag with the standard TTL.
func (c *Client) SetFlag(ctx context.Context, tenantID, flagName string, value bool) error {
	return c.SetEX(ctx, flagKey(tenantID, flagName), strconv.FormatBool(value), flagTTL)
}

// InvalidateFlags removes all cached flags for a tenant.
func (c *Client) InvalidateFlags(ctx context.Context, tenantID string) error {
	pattern := fmt.Sprintf("flags:%s:*", tenantID)
	cmd := c.B().Scan().Cursor(0).Match(pattern).Count(200).Build()
	_, keys, err := c.Do(ctx, cmd).AsScanEntry()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.Del(ctx, keys...)
}
