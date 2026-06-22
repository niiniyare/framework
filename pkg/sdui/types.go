// Package sdui defines the types for server-driven UI page builders.
package sdui

import "encoding/json"

// PageBuilderFn generates an amis JSON page definition for an EntityDefinition.
// It receives the entity definition as an opaque any (internal/core.EntityDefinition)
// to avoid circular imports. Use the provided accessor helpers in pkg/entitydef.
//
// The returned JSON is cached in Redis; the builder is only invoked on cache miss.
type PageBuilderFn func(def any) (json.RawMessage, error)

// Component is a JSON-serialisable amis component node.
// The map[string]any representation keeps the SDUI layer schema-agnostic —
// amis schema evolves independently of the framework version.
type Component = map[string]any
