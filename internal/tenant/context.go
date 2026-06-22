package tenant

import "context"

type tenantContextKey struct{}

// WithTenant attaches a Tenant to the context.
func WithTenant(ctx context.Context, t *Tenant) context.Context {
	return context.WithValue(ctx, tenantContextKey{}, t)
}

// FromCtx extracts the Tenant from the context. Returns nil if not set.
func FromCtx(ctx context.Context) *Tenant {
	t, _ := ctx.Value(tenantContextKey{}).(*Tenant)
	return t
}

// MustFromCtx extracts the Tenant and panics if absent.
// Use only in request handlers where the tenant middleware is guaranteed to run first.
func MustFromCtx(ctx context.Context) *Tenant {
	t := FromCtx(ctx)
	if t == nil {
		panic("tenant: context has no tenant — middleware not applied?")
	}
	return t
}
