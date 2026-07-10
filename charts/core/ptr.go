package core

// Ptr returns a pointer to v. It is a convenience for populating the optional
// pointer-typed fields on chart Props — most notably the Enable* toggles, where
// a nil pointer means "use the chart's default" and a non-nil pointer is an
// explicit override.
func Ptr[T any](v T) *T { return &v }

// BoolPtr returns a pointer to b. Equivalent to Ptr(b), kept as a named helper
// for readable call sites that set an optional bool prop, e.g.
// core.BoolPtr(true).
func BoolPtr(b bool) *bool { return &b }

// FloatPtr returns a pointer to f, for optional float64 props.
func FloatPtr(f float64) *float64 { return &f }
