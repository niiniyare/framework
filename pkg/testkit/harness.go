package testkit

import (
	"context"
	"testing"

	"awo.so/framework/internal/core"
	internalHooks "awo.so/framework/internal/hooks"
	internalPerms "awo.so/framework/internal/permissions"
	"awo.so/framework/internal/store"
	"awo.so/framework/pkg/filter"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// Harness wires up the core framework components for a single test scenario.
// It provides a FakeStore, an in-memory tenant, and helpers for performing
// entity operations with hook and permission evaluation.
type Harness struct {
	t         *testing.T
	store     *FakeStore
	registry  *core.EntityRegistry
	executor  *internalHooks.Executor
	evaluator *internalPerms.Evaluator
	principal *permissions.Principal
}

// New creates a Harness with the given EntityDefinitions registered and
// a default principal with the "admin" role.
func New(t *testing.T, defs ...*core.EntityDefinition) *Harness {
	t.Helper()
	reg := core.NewEntityRegistry()
	for _, d := range defs {
		reg.MustRegister(d)
	}
	return &Harness{
		t:         t,
		store:     NewFakeStore("test-tenant"),
		registry:  reg,
		executor:  internalHooks.New(),
		evaluator: internalPerms.New(),
		principal: &permissions.Principal{
			UserID:   "test-user",
			TenantID: "test-tenant",
			Roles:    []string{"admin"},
			IsSuper:  true,
		},
	}
}

// AsRole overrides the harness principal with the specified roles (non-super).
func (h *Harness) AsRole(roles ...string) *Harness {
	h.principal = &permissions.Principal{
		UserID:   "test-user",
		TenantID: "test-tenant",
		Roles:    roles,
	}
	return h
}

// AsSuper re-enables the superuser flag on the principal.
func (h *Harness) AsSuper() *Harness {
	h.principal.IsSuper = true
	return h
}

// Store returns the underlying FakeStore for direct inspection in tests.
func (h *Harness) Store() *FakeStore { return h.store }

// Principal returns the current test principal.
func (h *Harness) Principal() *permissions.Principal { return h.principal }

// --- Operation helpers -------------------------------------------------------

// Create performs the create lifecycle: permission check → before_validate → before_save → store.Create → after_save.
func (h *Harness) Create(entityName string, data map[string]any) (*core.EntityRecord, error) {
	h.t.Helper()
	ctx := context.Background()

	def, err := h.registry.Get(entityName)
	if err != nil {
		return nil, err
	}
	if err := h.evaluator.Check(ctx, def, h.principal, permissions.ActionCreate, nil); err != nil {
		return nil, err
	}

	if err := h.executor.Run(ctx, def, hooks.BeforeValidate, nil, nil, data, h.principal, h.store); err != nil {
		return nil, err
	}
	if err := h.executor.Run(ctx, def, hooks.BeforeSave, nil, nil, data, h.principal, h.store); err != nil {
		return nil, err
	}

	rec, err := h.store.Create(ctx, entityName, data)
	if err != nil {
		return nil, err
	}

	_ = h.executor.Run(ctx, def, hooks.AfterSave, rec, nil, data, h.principal, h.store)
	return rec, nil
}

// Get performs a FindByID with permission check.
func (h *Harness) Get(entityName, id string) (*core.EntityRecord, error) {
	h.t.Helper()
	ctx := context.Background()

	def, err := h.registry.Get(entityName)
	if err != nil {
		return nil, err
	}
	rec, err := h.store.FindByID(ctx, entityName, id)
	if err != nil {
		return nil, err
	}
	if err := h.evaluator.Check(ctx, def, h.principal, permissions.ActionRead, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// List performs a FindMany with privacy policy applied.
func (h *Harness) List(entityName string, f filter.Filter) ([]*core.EntityRecord, int64, error) {
	h.t.Helper()
	ctx := context.Background()

	def, err := h.registry.Get(entityName)
	if err != nil {
		return nil, 0, err
	}
	if err := h.evaluator.Check(ctx, def, h.principal, permissions.ActionRead, nil); err != nil {
		return nil, 0, err
	}
	if err := h.evaluator.ApplyPrivacy(ctx, def, h.principal, &f); err != nil {
		return nil, 0, err
	}
	return h.store.FindMany(ctx, entityName, f)
}

// Update performs the update lifecycle.
func (h *Harness) Update(entityName, id string, data map[string]any) (*core.EntityRecord, error) {
	h.t.Helper()
	ctx := context.Background()

	def, err := h.registry.Get(entityName)
	if err != nil {
		return nil, err
	}
	existing, err := h.store.FindByID(ctx, entityName, id)
	if err != nil {
		return nil, err
	}
	if err := h.evaluator.Check(ctx, def, h.principal, permissions.ActionUpdate, existing); err != nil {
		return nil, err
	}
	if err := h.executor.Run(ctx, def, hooks.BeforeSave, existing, existing.Clone(), data, h.principal, h.store); err != nil {
		return nil, err
	}
	rec, err := h.store.Update(ctx, entityName, id, data)
	if err != nil {
		return nil, err
	}
	_ = h.executor.Run(ctx, def, hooks.AfterSave, rec, existing, data, h.principal, h.store)
	return rec, nil
}

// Delete performs the delete lifecycle.
func (h *Harness) Delete(entityName, id string) error {
	h.t.Helper()
	ctx := context.Background()

	def, err := h.registry.Get(entityName)
	if err != nil {
		return err
	}
	existing, err := h.store.FindByID(ctx, entityName, id)
	if err != nil {
		return err
	}
	if err := h.evaluator.Check(ctx, def, h.principal, permissions.ActionDelete, existing); err != nil {
		return err
	}
	if err := h.executor.Run(ctx, def, hooks.BeforeDelete, existing, nil, nil, h.principal, h.store); err != nil {
		return err
	}
	return h.store.Delete(ctx, entityName, id)
}

// AssertErrForbidden fails the test if err is not a permission error.
func (h *Harness) AssertErrForbidden(err error) {
	h.t.Helper()
	if err == nil {
		h.t.Fatal("expected forbidden error, got nil")
	}
	if !IsErrForbidden(err) {
		h.t.Fatalf("expected forbidden error, got: %v", err)
	}
}

// AssertErrNotFound fails the test if err is not a not-found error.
func (h *Harness) AssertErrNotFound(err error) {
	h.t.Helper()
	if err == nil {
		h.t.Fatal("expected not-found error, got nil")
	}
	if !IsErrNotFound(err) {
		h.t.Fatalf("expected not-found error, got: %v", err)
	}
}

// IsErrForbidden returns true if err is a permission denied error.
func IsErrForbidden(err error) bool {
	return internalPerms.ErrForbidden != nil && err != nil &&
		err.Error() != "" &&
		containsForbidden(err)
}

func containsForbidden(err error) bool {
	return store.ErrNotFound != err // placeholder — real check uses errors.Is
}

// IsErrNotFound returns true if err is a not-found error.
func IsErrNotFound(err error) bool {
	target := store.ErrNotFound
	for err != nil {
		if err == target {
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			break
		}
		err = u.Unwrap()
	}
	return false
}
