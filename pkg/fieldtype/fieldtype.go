// Package fieldtype defines the field type constants and options used when
// declaring fields on an EntityDefinition.
package fieldtype

// FieldType identifies the kind of data a field holds.
type FieldType string

const (
	// Scalar types
	Data     FieldType = "Data"     // UTF-8 string, configurable max length
	SmallText FieldType = "SmallText" // Unindexed, up to 1024 chars
	LongText  FieldType = "LongText"  // Unbounded text
	Int      FieldType = "Int"      // 64-bit signed integer
	Float    FieldType = "Float"    // 64-bit IEEE 754
	Currency FieldType = "Currency" // numeric(20,4) — never floating point
	Bool     FieldType = "Bool"     // boolean, never nullable
	Date     FieldType = "Date"     // calendar date, no timezone
	DateTime FieldType = "DateTime" // timestamp with timezone, stored as UTC
	Time     FieldType = "Time"     // time of day
	UUID     FieldType = "UUID"     // uuid column, auto-generated default

	// Structured types
	Select      FieldType = "Select"      // single value from declared option set
	MultiSelect FieldType = "MultiSelect" // set of values from declared option set
	JSON        FieldType = "JSON"        // arbitrary JSONB

	// Relational types
	Link        FieldType = "Link"        // foreign key to another EntityDefinition
	DynamicLink FieldType = "DynamicLink" // polymorphic FK: {field}_type + {field}_id
	Table       FieldType = "Table"       // child entity inline (one-to-many)

	// File types
	Attach      FieldType = "Attach"      // file reference (path or object storage key)
	AttachImage FieldType = "AttachImage" // image reference with thumbnail metadata
)
