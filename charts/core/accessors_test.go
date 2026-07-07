package core_test

import (
	"testing"
	"time"

	"github.com/geoffjay/templ-charts/charts/core"
)

type accessorDatum struct {
	ID    string
	Value float64
	Data  map[string]any
}

func TestGetPropertyAccessor_NilAndEmpty(t *testing.T) {
	d := accessorDatum{ID: "a", Value: 3}
	if got := core.GetPropertyAccessor[accessorDatum, float64](nil)(d); got != 0 {
		t.Errorf("nil accessor = %v, want 0", got)
	}
	if got := core.GetPropertyAccessor[accessorDatum, string]("")(d); got != "" {
		t.Errorf("empty accessor = %q, want empty", got)
	}
	if got := core.GetPropertyAccessor[accessorDatum, string]("   ")(d); got != "" {
		t.Errorf("whitespace accessor = %q, want empty", got)
	}
}

func TestGetPropertyAccessor_StringPaths(t *testing.T) {
	d := accessorDatum{
		ID:    "series-1",
		Value: 42.5,
		Data:  map[string]any{"color": "#f00", "nested": map[string]any{"n": 7}},
	}
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"struct field", core.GetPropertyAccessor[accessorDatum, string]("ID")(d), "series-1"},
		{"case-insensitive field", core.GetPropertyAccessor[accessorDatum, string]("id")(d), "series-1"},
		{"float field", core.GetPropertyAccessor[accessorDatum, float64]("Value")(d), 42.5},
		{"map key", core.GetPropertyAccessor[accessorDatum, string]("Data.color")(d), "#f00"},
		{"nested map", core.GetPropertyAccessor[accessorDatum, float64]("Data.nested.n")(d), 7.0},
		{"missing field", core.GetPropertyAccessor[accessorDatum, string]("Nope")(d), ""},
		{"missing map key", core.GetPropertyAccessor[accessorDatum, float64]("Data.missing")(d), 0.0},
		{"path through scalar", core.GetPropertyAccessor[accessorDatum, string]("Value.x")(d), ""},
	}
	for _, tc := range tests {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
}

func TestGetPropertyAccessor_PointerDatum(t *testing.T) {
	d := &accessorDatum{ID: "p"}
	if got := core.GetPropertyAccessor[*accessorDatum, string]("ID")(d); got != "p" {
		t.Errorf("pointer datum = %q, want %q", got, "p")
	}
}

func TestGetPropertyAccessor_MapDatum(t *testing.T) {
	m := map[string]any{"a": map[string]any{"b": "deep"}}
	if got := core.GetPropertyAccessor[map[string]any, string]("a.b")(m); got != "deep" {
		t.Errorf("map datum = %q, want %q", got, "deep")
	}
	// Non-string map keys cannot be addressed by a string path.
	im := map[int]any{1: "x"}
	if got := core.GetPropertyAccessor[map[int]any, string]("1")(im); got != "" {
		t.Errorf("int-keyed map = %q, want empty", got)
	}
}

func TestGetPropertyAccessor_Func(t *testing.T) {
	d := accessorDatum{Value: 5}
	f := core.GetPropertyAccessor[accessorDatum, float64](func(d accessorDatum) float64 { return d.Value * 2 })
	if got := f(d); got != 10 {
		t.Errorf("func accessor = %v, want 10", got)
	}
}

func TestGetPropertyAccessor_FuncLikeViaReflection(t *testing.T) {
	d := accessorDatum{Value: 2.5}
	// A func with a differing signature is invoked via reflection.
	f := core.GetPropertyAccessor[accessorDatum, float64](func(d accessorDatum) any { return d.Value })
	if got := f(d); got != 2.5 {
		t.Errorf("func-like accessor = %v, want 2.5", got)
	}
	// Wrong return type yields the zero value.
	bad := core.GetPropertyAccessor[accessorDatum, float64](func(d accessorDatum) any { return "nope" })
	if got := bad(d); got != 0 {
		t.Errorf("wrong-typed func = %v, want 0", got)
	}
	// A func with no return values yields the zero value.
	void := core.GetPropertyAccessor[accessorDatum, float64](func(d accessorDatum) {})
	if got := void(d); got != 0 {
		t.Errorf("void func = %v, want 0", got)
	}
}

func TestGetPropertyAccessor_UnsupportedType(t *testing.T) {
	f := core.GetPropertyAccessor[accessorDatum, string](42)
	if got := f(accessorDatum{ID: "x"}); got != "" {
		t.Errorf("unsupported accessor = %q, want empty", got)
	}
}

func TestGetPropertyAccessor_Coercion(t *testing.T) {
	d := accessorDatum{
		Value: 2.9,
		Data: map[string]any{
			"num":  "3.5",
			"flag": true,
			"when": time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	if got := core.GetPropertyAccessor[accessorDatum, string]("Value")(d); got != "2.9" {
		t.Errorf("float→string = %q, want %q", got, "2.9")
	}
	if got := core.GetPropertyAccessor[accessorDatum, int]("Value")(d); got != 2 {
		t.Errorf("float→int = %v, want 2", got)
	}
	if got := core.GetPropertyAccessor[accessorDatum, float64]("Data.num")(d); got != 3.5 {
		t.Errorf("string→float = %v, want 3.5", got)
	}
	if got := core.GetPropertyAccessor[accessorDatum, bool]("Data.flag")(d); !got {
		t.Errorf("bool = %v, want true", got)
	}
	want := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	if got := core.GetPropertyAccessor[accessorDatum, time.Time]("Data.when")(d); !got.Equal(want) {
		t.Errorf("time = %v, want %v", got, want)
	}
	// Un-coercible values yield the zero V.
	if got := core.GetPropertyAccessor[accessorDatum, time.Time]("Data.num")(d); !got.IsZero() {
		t.Errorf("string→time = %v, want zero", got)
	}
}

func TestGetLabelGenerator(t *testing.T) {
	d := accessorDatum{ID: "lab", Value: 9}
	if got := core.GetLabelGenerator[accessorDatum, string](func(d accessorDatum) string { return "fn:" + d.ID }, "")(d); got != "fn:lab" {
		t.Errorf("func label = %q", got)
	}
	if got := core.GetLabelGenerator[accessorDatum, string]("ID", "")(d); got != "lab" {
		t.Errorf("string label = %q, want %q", got, "lab")
	}
	// Blank label string falls through to labelFrom.
	if got := core.GetLabelGenerator[accessorDatum, string]("  ", "ID")(d); got != "lab" {
		t.Errorf("blank label + labelFrom = %q, want %q", got, "lab")
	}
	if got := core.GetLabelGenerator[accessorDatum, string](nil, "ID")(d); got != "lab" {
		t.Errorf("nil label + labelFrom = %q, want %q", got, "lab")
	}
	// A label func of the wrong type is ignored in favor of labelFrom.
	if got := core.GetLabelGenerator[accessorDatum, string](func(d accessorDatum) int { return 1 }, "ID")(d); got != "lab" {
		t.Errorf("mistyped label + labelFrom = %q, want %q", got, "lab")
	}
	if got := core.GetLabelGenerator[accessorDatum, string](nil, "")(d); got != "" {
		t.Errorf("no label = %q, want empty", got)
	}
}
