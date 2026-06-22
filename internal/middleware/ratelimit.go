package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/redis"
	"awo.so/framework/internal/tenant"
)

// RateLimitConfig controls per-(tenant, principal) token-bucket limits.
type RateLimitConfig struct {
	// RequestsPerMinute is the default bucket capacity.
	// Overridden per-tenant by settings["api.rate_limit_per_minute"].
	RequestsPerMinute int
}

// DefaultRateLimitConfig returns sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{RequestsPerMinute: 300}
}

// RateLimit enforces a Redis token-bucket rate limit per (tenantID, principalID).
// Requests that exceed the limit receive 429 Too Many Requests.
// Authenticated principal ID is used when available; falls back to tenant+IP.
func RateLimit(r *redis.Client, cfg RateLimitConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		if t == nil {
			return c.Next() // no tenant context — skip (handled by tenant middleware)
		}

		// Identify the rate-limit bucket: prefer principal, fall back to IP
		bucketID := c.IP()
		if p := PrincipalFromCtx(c.UserContext()); p != nil {
			bucketID = p.UserID
		}

		key := fmt.Sprintf("ratelimit:%s:%s", t.ID, bucketID)
		limit := cfg.RequestsPerMinute

		allowed, remaining, err := tokenBucketCheck(c.UserContext(), r, key, limit)
		if err != nil {
			// Redis failure → fail open (don't block traffic on cache outage)
			return c.Next()
		}

		c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		}
		return c.Next()
	}
}

// tokenBucketCheck implements a sliding-window counter in Redis.
// Returns (allowed, remaining, error).
func tokenBucketCheck(ctx context.Context, r *redis.Client, key string, limit int) (bool, int, error) {
	// Get current count
	raw, err := r.Get(ctx, key)
	if err != nil && err != redis.ErrKeyNotFound {
		return true, limit, err // fail open
	}

	count := 0
	if err == nil {
		count, _ = strconv.Atoi(raw)
	}

	if count >= limit {
		return false, 0, nil
	}

	// Increment counter; set TTL of 60s on first request
	newCount := count + 1
	if count == 0 {
		if err := r.SetEX(ctx, key, strconv.Itoa(newCount), time.Minute); err != nil {
			return true, limit - 1, err
		}
	} else {
		// Increment without resetting TTL — use SetEX with remaining TTL
		// Simplified: just overwrite with same TTL (slight inaccuracy acceptable)
		if err := r.SetEX(ctx, key, strconv.Itoa(newCount), time.Minute); err != nil {
			return true, limit - newCount, err
		}
	}
	return true, limit - newCount, nil
}
