// Package workflow wraps the Temporal SDK client and provides workflow dispatch
// helpers used by the after_save and on_submit hook pipeline.
package workflow

import (
	"fmt"

	"go.temporal.io/sdk/client"
)

// Config holds Temporal connection parameters.
type Config struct {
	HostPort  string // e.g. "localhost:7233"
	Namespace string // e.g. "awo"
}

// DefaultConfig returns sensible local-dev defaults.
func DefaultConfig() Config {
	return Config{
		HostPort:  "localhost:7233",
		Namespace: "default",
	}
}

// Client wraps the Temporal client with framework helpers.
type Client struct {
	client.Client
	cfg Config
}

// New creates a connected Temporal client.
// The caller must call Close() when the process shuts down.
func New(cfg Config) (*Client, error) {
	c, err := client.Dial(client.Options{
		HostPort:  cfg.HostPort,
		Namespace: cfg.Namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("temporal: connect to %s: %w", cfg.HostPort, err)
	}
	return &Client{Client: c, cfg: cfg}, nil
}
