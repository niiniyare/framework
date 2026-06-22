package core

import "context"

type contextKey int

const (
	contextKeyTenant    contextKey = iota
	contextKeyPrincipal            // *Principal — set by auth middleware
	contextKeyRequestID            // string
	contextKeyLogger               // *slog.Logger
	contextKeyTx                   // active store transaction marker
)

// WithRequestID attaches a request ID to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, id)
}

// RequestIDFromCtx extracts the request ID from the context.
func RequestIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(contextKeyRequestID).(string)
	return v
}
