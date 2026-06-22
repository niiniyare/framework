package hooks

import (
	"context"
	"fmt"

	"awo.so/framework/internal/core"
	"awo.so/framework/pkg/hooks"
)

// Executor runs lifecycle hooks for entities. It holds no state — callers pass
// the EntityDefinition whose hooks should be invoked. Thread-safe.
type Executor struct{}

// New returns an Executor.
func New() *Executor { return &Executor{} }

// Run looks up the hook chain for (entity, event) and executes it.
// Returns nil if no hooks are registered for that event.
func (e *Executor) Run(
	ctx context.Context,
	def *core.EntityDefinition,
	event hooks.HookEvent,
	record *core.EntityRecord,
	prevRecord *core.EntityRecord,
	data map[string]any,
	principal any,
	repo any,
) error {
	fns := def.HooksFor(event)
	if len(fns) == 0 {
		return nil
	}

	hctx := &hooks.HookContext{
		Event:      event,
		EntityName: def.Name,
		Record:     record,
		PrevRecord: prevRecord,
		Data:       data,
		Principal:  principal,
		Repo:       repo,
	}

	chain := NewChain(def.Name, event, fns)
	if err := chain.Execute(ctx, hctx); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
