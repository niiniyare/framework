// Package redis provides rueidis client construction and namespace helpers.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/rueidis"
)

// Config holds Redis connection parameters.
type Config struct {
	Addr     string // host:port, e.g. "localhost:6379"
	Password string
	DB       int
}

// Client is a thin wrapper around rueidis.Client with convenience helpers.
type Client struct {
	rueidis.Client
}

// New creates and returns a connected Redis client.
// The caller is responsible for calling Close() when done.
func New(cfg Config) (*Client, error) {
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{cfg.Addr},
		Password:    cfg.Password,
		SelectDB:    cfg.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("redis: connect to %s: %w", cfg.Addr, err)
	}
	return &Client{c}, nil
}

// Ping verifies the connection is alive. Used in health checks.
func (c *Client) Ping(ctx context.Context) error {
	cmd := c.B().Ping().Build()
	return c.Do(ctx, cmd).Error()
}

// SetEX sets key to value with a TTL.
func (c *Client) SetEX(ctx context.Context, key string, value string, ttl time.Duration) error {
	cmd := c.B().Set().Key(key).Value(value).Ex(ttl).Build()
	return c.Do(ctx, cmd).Error()
}

// Get retrieves a string value by key.
// Returns ("", ErrKeyNotFound) if the key does not exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	cmd := c.B().Get().Key(key).Build()
	val, err := c.Do(ctx, cmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", ErrKeyNotFound
		}
		return "", err
	}
	return val, nil
}

// Del deletes one or more keys.
func (c *Client) Del(ctx context.Context, keys ...string) error {
	cmd := c.B().Del().Key(keys...).Build()
	return c.Do(ctx, cmd).Error()
}

// ErrKeyNotFound is returned when a Redis key does not exist.
var ErrKeyNotFound = fmt.Errorf("redis: key not found")
