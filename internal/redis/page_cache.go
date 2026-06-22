package redis

import (
	"context"
	"fmt"
	"time"
)

const pageTTL = time.Hour

// pageKey returns the Redis key for a cached SDUI page definition.
// Format: sdui:{tenantID}:{entityName}:{role}
func pageKey(tenantID, entityName, role string) string {
	return fmt.Sprintf("sdui:%s:%s:%s", tenantID, entityName, role)
}

// GetPage retrieves a cached amis page JSON. Returns ("", ErrKeyNotFound) on miss.
func (c *Client) GetPage(ctx context.Context, tenantID, entityName, role string) (string, error) {
	return c.Get(ctx, pageKey(tenantID, entityName, role))
}

// SetPage stores a generated amis page JSON with the standard TTL.
func (c *Client) SetPage(ctx context.Context, tenantID, entityName, role, json string) error {
	return c.SetEX(ctx, pageKey(tenantID, entityName, role), json, pageTTL)
}

// InvalidatePage removes the cached page for a specific role.
func (c *Client) InvalidatePage(ctx context.Context, tenantID, entityName, role string) error {
	return c.Del(ctx, pageKey(tenantID, entityName, role))
}

// InvalidateEntityPages removes all cached pages for an entity across all roles.
func (c *Client) InvalidateEntityPages(ctx context.Context, tenantID, entityName string) error {
	pattern := fmt.Sprintf("sdui:%s:%s:*", tenantID, entityName)
	cmd := c.B().Scan().Cursor(0).Match(pattern).Count(100).Build()
	_, keys, err := c.Do(ctx, cmd).AsScanEntry()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.Del(ctx, keys...)
}
