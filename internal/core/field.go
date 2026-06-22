package core

import "awo.so/framework/pkg/fieldtype"

// Field is the internal representation of a single field on an EntityDefinition.
type Field struct {
	Name  string
	Label string
	Type  fieldtype.FieldType

	// Constraints
	Required     bool
	Unique       bool
	Immutable    bool
	Sensitive    bool
	Translatable bool
	ReadOnly     bool

	// Value constraints
	Default any
	MaxLen  int
	Min     *int64
	Max     *int64

	// Select / MultiSelect
	Choices []string

	// Link / DynamicLink / Table
	LinkTarget string // target EntityDefinition name

	// Naming series format string
	NamingSeries string

	// UI hints (used by SDUI page builder)
	Placeholder string
	Description string
	Hidden      bool
	ColSpan     int // 1–12 grid columns in form layout
}
