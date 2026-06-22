// Package hooks implements the hook chain execution engine.
package hooks

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/hooks"
)

// Chain executes a slice of HookFn in declaration order.
// The first non-nil error aborts the chain and is returned (wrapped with position context).
type Chain struct {
	EntityName string
	Event      hooks.HookEvent
	fns        []hooks.HookFn
}

// NewChain constructs a Chain for the given entity and event.
func NewChain(entityName string, event hooks.HookEvent, fns []hooks.HookFn) *Chain {
	return &Chain{EntityName: entityName, Event: event, fns: fns}
}

// Execute calls each hook in order. Stops on first error.
func (c *Chain) Execute(ctx context.Context, hctx *hooks.HookContext) error {
	for i, fn := range c.fns {
		if err := fn(ctx, hctx); err != nil {
			return fmt.Errorf("hook[%d] %s/%s: %w", i, c.EntityName, c.Event, err)
		}
	}
	return nil
}

// Len returns the number of hooks in the chain.
func (c *Chain) Len() int { return len(c.fns) }
