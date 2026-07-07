package colors

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/theming"
)

func TestNewFuncColor(t *testing.T) {
	cfg := NewFuncColor(func(any) string { return "#0f0f0f" })
	if cfg.Type != InheritedColorTypeFunc {
		t.Fatalf("type = %v, want func", cfg.Type)
	}
	if got := cfg.Func(nil); got != "#0f0f0f" {
		t.Fatalf("func = %q, want #0f0f0f", got)
	}
}

func TestParseInheritedColorConfig(t *testing.T) {
	// nil → empty static.
	cfg, err := ParseInheritedColorConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeStatic || cfg.Static != "" {
		t.Errorf("nil = %+v", cfg)
	}
	// string → static.
	cfg, err = ParseInheritedColorConfig("#abc")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeStatic || cfg.Static != "#abc" {
		t.Errorf("string = %+v", cfg)
	}
	// func(any) string.
	cfg, err = ParseInheritedColorConfig(func(any) string { return "#111" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeFunc || cfg.Func(nil) != "#111" {
		t.Errorf("func = %+v", cfg)
	}
	// func(any) any is adapted to string.
	cfg, err = ParseInheritedColorConfig(func(any) any { return "#222" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeFunc || cfg.Func(nil) != "#222" {
		t.Errorf("func any = %+v", cfg)
	}
	// map theme.
	cfg, err = ParseInheritedColorConfig(map[string]any{"theme": "labels.text.fill"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeTheme || cfg.ThemePath != "labels.text.fill" {
		t.Errorf("theme = %+v", cfg)
	}
	// map from + modifiers (short entries skipped).
	cfg, err = ParseInheritedColorConfig(map[string]any{
		"from":      "color",
		"modifiers": []any{[]any{"darker", 1.0}, []any{"tooshort"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != InheritedColorTypeFromContext || cfg.FromPath != "color" {
		t.Errorf("from = %+v", cfg)
	}
	if len(cfg.Modifiers) != 1 || cfg.Modifiers[0][0] != "darker" {
		t.Errorf("modifiers = %+v, want single darker", cfg.Modifiers)
	}
	// invalid map.
	if _, err = ParseInheritedColorConfig(map[string]any{"bogus": 1}); err == nil {
		t.Error("invalid map should error")
	}
	// unsupported type.
	if _, err = ParseInheritedColorConfig(42); err == nil {
		t.Error("unsupported type should error")
	}
}

func TestGetInheritedColorGenerator_Static(t *testing.T) {
	gen := GetInheritedColorGenerator(NewStaticColor("#333333"), nil)
	if got := gen(nil); got != "#333333" {
		t.Fatalf("static gen = %q", got)
	}
}

func TestGetInheritedColorGenerator_Func(t *testing.T) {
	gen := GetInheritedColorGenerator(NewFuncColor(func(d any) string { return d.(string) }), nil)
	if got := gen("#444444"); got != "#444444" {
		t.Fatalf("func gen = %q", got)
	}
}

func TestGetInheritedColorGenerator_Theme(t *testing.T) {
	theme := theming.Theme{Text: theming.TextStyle{Fill: "#987654"}}
	gen := GetInheritedColorGenerator(NewThemeColor("Text.Fill"), &theme)
	if got := gen(nil); got != "#987654" {
		t.Errorf("theme gen = %q, want #987654", got)
	}
	// Path lookup is case-insensitive on struct fields.
	gen = GetInheritedColorGenerator(NewThemeColor("text.fill"), &theme)
	if got := gen(nil); got != "#987654" {
		t.Errorf("case-insensitive theme gen = %q, want #987654", got)
	}
	// Unknown path resolves to empty.
	gen = GetInheritedColorGenerator(NewThemeColor("Not.A.Path"), &theme)
	if got := gen(nil); got != "" {
		t.Errorf("unknown path = %q, want empty", got)
	}
	// Nil theme resolves to empty.
	gen = GetInheritedColorGenerator(NewThemeColor("Text.Fill"), nil)
	if got := gen(nil); got != "" {
		t.Errorf("nil theme = %q, want empty", got)
	}
}

func TestGetInheritedColorGenerator_FromContext(t *testing.T) {
	// No modifiers: raw path value.
	gen := GetInheritedColorGenerator(NewFromContextColor("color", nil), nil)
	if got := gen(map[string]string{"color": "#ff0000"}); got != "#ff0000" {
		t.Fatalf("from-context gen = %q", got)
	}
	// With a darker modifier (RGB space, pinned against d3-color).
	cfg := NewFromContextColor("color", []ColorModifier{{"darker", 1.0}})
	gen = GetInheritedColorGenerator(cfg, nil)
	if got := gen(map[string]string{"color": "#ff0000"}); got != "#b30000" {
		t.Fatalf("darker(1) red = %q, want #b30000", got)
	}
}

func TestApplyColorModifiers_OpacityAndChaining(t *testing.T) {
	if got := ApplyColorModifiers("#ff0000", [][2]any{{"opacity", 0.5}}); got != "rgba(255, 0, 0, 0.5)" {
		t.Errorf("opacity = %q, want rgba(255, 0, 0, 0.5)", got)
	}
	// Empty modifier list returns the color unchanged.
	if got := ApplyColorModifiers("#123456", nil); got != "#123456" {
		t.Errorf("no modifiers = %q, want passthrough", got)
	}
	// Non-string kind and integer amounts are tolerated.
	got := ApplyColorModifiers("#808080", [][2]any{{42, 1}, {"brighter", 1}})
	if got != "#b7b7b7" {
		t.Errorf("brighter(1) gray = %q, want #b7b7b7", got)
	}
}

func TestGetPathValue(t *testing.T) {
	type inner struct{ Color string }
	type outer struct {
		Data  inner
		Ptr   *inner
		Other int
	}
	d := outer{Data: inner{Color: "#101010"}, Ptr: &inner{Color: "#202020"}}
	if got := getPathValue(d, "Data.Color"); got != "#101010" {
		t.Errorf("struct path = %q", got)
	}
	if got := getPathValue(&d, "Ptr.Color"); got != "#202020" {
		t.Errorf("pointer path = %q", got)
	}
	if got := getPathValue(d, "data.color"); got != "#101010" {
		t.Errorf("case-insensitive path = %q", got)
	}
	if got := getPathValue(d, "Nope"); got != "" {
		t.Errorf("missing field = %q, want empty", got)
	}
	if got := getPathValue(d, "Other.Deeper"); got != "" {
		t.Errorf("path through scalar = %q, want empty", got)
	}
	// Map lookups.
	m := map[string]any{"a": map[string]any{"b": "#303030"}}
	if got := getPathValue(m, "a.b"); got != "#303030" {
		t.Errorf("map path = %q", got)
	}
	if got := getPathValue(m, "a.missing"); got != "" {
		t.Errorf("missing map key = %q, want empty", got)
	}
	// Map with non-string keys cannot be walked.
	im := map[int]string{1: "x"}
	if got := getPathValue(im, "1"); got != "" {
		t.Errorf("int-keyed map = %q, want empty", got)
	}
}
