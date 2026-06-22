package fieldtype

// FieldOption configures optional attributes on a field declaration.
type FieldOption func(*FieldOptions)

// FieldOptions holds the complete set of options for a field.
// Zero value is a valid, unconstrained field.
type FieldOptions struct {
	Required     bool
	Unique       bool
	Immutable    bool   // set on create, rejected on update
	Sensitive    bool   // excluded from logs and API responses by default
	Translatable bool   // value stored with locale key
	ReadOnly     bool   // never accepted in write payloads

	Default any    // static value; use DefaultFn for computed defaults
	MaxLen  int    // for Data, SmallText
	Min     *int64 // for Int, Float, Currency
	Max     *int64 // for Int, Float, Currency

	// Select / MultiSelect
	Choices []string

	// Link / Table
	LinkTarget string // name of the target EntityDefinition

	// Naming series (format string, e.g. "INV-{YYYY}-{SEQ:5}")
	NamingSeries string
}

// Required marks the field as non-nullable.
func Required() FieldOption {
	return func(o *FieldOptions) { o.Required = true }
}

// Unique adds a unique index on the field.
func Unique() FieldOption {
	return func(o *FieldOptions) { o.Unique = true }
}

// Immutable prevents the field from being updated after creation.
func Immutable() FieldOption {
	return func(o *FieldOptions) { o.Immutable = true }
}

// Sensitive excludes the field from logs and responses unless explicitly requested.
func Sensitive() FieldOption {
	return func(o *FieldOptions) { o.Sensitive = true }
}

// Translatable marks the field value as locale-aware.
func Translatable() FieldOption {
	return func(o *FieldOptions) { o.Translatable = true }
}

// ReadOnly prevents the field from being set via the API.
func ReadOnly() FieldOption {
	return func(o *FieldOptions) { o.ReadOnly = true }
}

// Default sets a static default value.
func Default(v any) FieldOption {
	return func(o *FieldOptions) { o.Default = v }
}

// MaxLen sets the maximum character length for string fields.
func MaxLen(n int) FieldOption {
	return func(o *FieldOptions) { o.MaxLen = n }
}

// MinVal sets a minimum numeric value.
func MinVal(n int64) FieldOption {
	return func(o *FieldOptions) { o.Min = &n }
}

// MaxVal sets a maximum numeric value.
func MaxVal(n int64) FieldOption {
	return func(o *FieldOptions) { o.Max = &n }
}

// Choices declares the allowed values for Select and MultiSelect fields.
func Choices(values ...string) FieldOption {
	return func(o *FieldOptions) { o.Choices = values }
}

// LinkTo sets the target EntityDefinition name for Link and Table fields.
func LinkTo(entityName string) FieldOption {
	return func(o *FieldOptions) { o.LinkTarget = entityName }
}

// WithNamingSeries sets the naming series format string.
// Tokens: {PREFIX}, {YYYY}, {MM}, {DD}, {SEQ}, {SEQ:N} (N = zero-padded width), {TENANT}.
func WithNamingSeries(format string) FieldOption {
	return func(o *FieldOptions) { o.NamingSeries = format }
}

// Apply applies a list of options to a FieldOptions value.
func Apply(opts []FieldOption) FieldOptions {
	var o FieldOptions
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
