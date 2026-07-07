package colors

import (
	"fmt"
	"testing"
)

// --- ParseOrdinalColorScaleConfig -------------------------------------------

func TestParseOrdinalColorScaleConfig_Nil(t *testing.T) {
	cfg, err := ParseOrdinalColorScaleConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeStatic || cfg.Static != "" {
		t.Fatalf("nil = %+v, want empty static", cfg)
	}
}

func TestParseOrdinalColorScaleConfig_Func(t *testing.T) {
	cfg, err := ParseOrdinalColorScaleConfig(func(any) string { return "#0f0" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeFunc || cfg.Func(nil) != "#0f0" {
		t.Fatalf("func config = %+v", cfg)
	}
}

func TestParseOrdinalColorScaleConfig_StringSlice(t *testing.T) {
	cfg, err := ParseOrdinalColorScaleConfig([]string{"#111", "#222"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeColors || len(cfg.Colors) != 2 {
		t.Fatalf("[]string config = %+v", cfg)
	}
}

func TestParseOrdinalColorScaleConfig_AnySlice(t *testing.T) {
	// Non-string entries are dropped.
	cfg, err := ParseOrdinalColorScaleConfig([]any{"#111", 3, "#222"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeColors || len(cfg.Colors) != 2 || cfg.Colors[1] != "#222" {
		t.Fatalf("[]any config = %+v", cfg)
	}
}

func TestParseOrdinalColorScaleConfig_MapForms(t *testing.T) {
	// datum path.
	cfg, err := ParseOrdinalColorScaleConfig(map[string]any{"datum": "data.color"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeDatum || cfg.DatumPath != "data.color" {
		t.Fatalf("datum config = %+v", cfg)
	}
	// scheme with each numeric size form.
	for _, size := range []any{5, int64(5), 5.0} {
		cfg, err = ParseOrdinalColorScaleConfig(map[string]any{"scheme": "blues", "size": size})
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Type != OrdinalTypeScheme || cfg.Scheme != "blues" || cfg.SchemeSize != 5 {
			t.Fatalf("scheme config with size %T = %+v", size, cfg)
		}
	}
	// scheme with unusable size type keeps size 0.
	cfg, err = ParseOrdinalColorScaleConfig(map[string]any{"scheme": "blues", "size": "5"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SchemeSize != 0 {
		t.Fatalf("string size should be ignored, got %+v", cfg)
	}
	// colors list.
	cfg, err = ParseOrdinalColorScaleConfig(map[string]any{"colors": []any{"#111", 7, "#222"}})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Type != OrdinalTypeColors || len(cfg.Colors) != 2 {
		t.Fatalf("colors config = %+v", cfg)
	}
	// invalid map.
	if _, err = ParseOrdinalColorScaleConfig(map[string]any{"bogus": true}); err == nil {
		t.Fatal("invalid map should error")
	}
}

func TestParseOrdinalColorScaleConfig_InvalidType(t *testing.T) {
	if _, err := ParseOrdinalColorScaleConfig(42); err == nil {
		t.Fatal("unsupported type should error")
	}
}

// --- GetOrdinalColorScale ----------------------------------------------------

func TestGetOrdinalColorScale_Static(t *testing.T) {
	s := GetOrdinalColorScale[string](OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: "#123"}, nil)
	if got := s("anything"); got != "#123" {
		t.Fatalf("static scale = %q, want #123", got)
	}
}

func TestGetOrdinalColorScale_Func(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeFunc, Func: func(d any) string { return "c-" + d.(string) }}
	s := GetOrdinalColorScale[string](cfg, nil)
	if got := s("x"); got != "c-x" {
		t.Fatalf("func scale = %q, want c-x", got)
	}
}

func TestGetOrdinalColorScale_Datum(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeDatum, DatumPath: "color"}
	s := GetOrdinalColorScale[map[string]string](cfg, nil)
	if got := s(map[string]string{"color": "#fedcba"}); got != "#fedcba" {
		t.Fatalf("datum scale = %q, want #fedcba", got)
	}
}

func TestGetOrdinalColorScale_ColorsCycleAndCache(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: []string{"#111", "#222"}}
	s := GetOrdinalColorScale[string](cfg, "")
	a1 := s("a")
	b1 := s("b")
	c1 := s("c") // cycles back to first color
	if a1 != "#111" || b1 != "#222" || c1 != "#111" {
		t.Fatalf("cycle = %s %s %s, want #111 #222 #111", a1, b1, c1)
	}
	if got := s("a"); got != a1 {
		t.Fatalf("repeat key = %q, want cached %q", got, a1)
	}
}

func TestGetOrdinalColorScale_SchemeWithSize(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: "brown_blueGreen", SchemeSize: 3}
	s := GetOrdinalColorScale[string](cfg, "")
	want := DivergingColorSchemes["brown_blueGreen"][3]
	if got := s("first"); got != want[0] {
		t.Fatalf("scheme scale first = %q, want %q", got, want[0])
	}
}

func TestGetOrdinalColorScale_ZeroConfig(t *testing.T) {
	// Zero value is OrdinalTypeStatic with empty color.
	s := GetOrdinalColorScale[string](OrdinalColorScaleConfig{}, nil)
	if got := s("x"); got != "" {
		t.Fatalf("zero config = %q, want empty", got)
	}
}

// --- resolveIdentity ---------------------------------------------------------

func TestResolveIdentityForms(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: []string{"#111", "#222"}}

	// func(D) string identity.
	sf := GetOrdinalColorScale[int](cfg, func(d int) string { return fmt.Sprintf("%d", d%2) })
	if sf(1) != sf(3) {
		t.Error("identity func should map 1 and 3 to the same color")
	}
	// func(any) string identity.
	sa := GetOrdinalColorScale[int](cfg, func(d any) string { return fmt.Sprintf("%v", d) })
	if sa(1) == sa(2) {
		t.Error("distinct keys should get distinct colors for a 2-color palette")
	}
	// string path identity over a map datum.
	sp := GetOrdinalColorScale[map[string]string](cfg, "id")
	if sp(map[string]string{"id": "a"}) != sp(map[string]string{"id": "a"}) {
		t.Error("same path key should map to the same color")
	}
	// empty string identity → datum itself.
	se := GetOrdinalColorScale[string](cfg, "")
	if se("k") != se("k") {
		t.Error("empty identity should key on the datum itself")
	}
	// nil identity → datum itself.
	sn := GetOrdinalColorScale[string](cfg, nil)
	if sn("k") != sn("k") {
		t.Error("nil identity should key on the datum itself")
	}
	// unsupported identity type → default formatting.
	su := GetOrdinalColorScale[string](cfg, 42)
	if su("k") != su("k") {
		t.Error("fallback identity should key on the datum itself")
	}
}

// --- resolveSchemeColors ------------------------------------------------------

func TestResolveSchemeColors(t *testing.T) {
	// Categorical: size ignored.
	if got := resolveSchemeColors("nivo", 3); len(got) != 6 {
		t.Errorf("nivo len = %d, want 6", len(got))
	}
	// Diverging: honored size and clamped defaults.
	if got := resolveSchemeColors("brown_blueGreen", 5); len(got) != 5 {
		t.Errorf("diverging size 5 len = %d", len(got))
	}
	if got := resolveSchemeColors("brown_blueGreen", 0); len(got) != 11 {
		t.Errorf("diverging size 0 len = %d, want 11", len(got))
	}
	if got := resolveSchemeColors("brown_blueGreen", 99); len(got) != 11 {
		t.Errorf("diverging size 99 len = %d, want 11", len(got))
	}
	// Sequential: honored size and clamped default.
	if got := resolveSchemeColors("blues", 4); len(got) != 4 {
		t.Errorf("sequential size 4 len = %d", len(got))
	}
	if got := resolveSchemeColors("blues", 0); len(got) != 9 {
		t.Errorf("sequential size 0 len = %d, want 9", len(got))
	}
	// Interpolator-only scheme: sampled at size stops (default 9).
	got := resolveSchemeColors("viridis", 0)
	if len(got) != 9 {
		t.Fatalf("viridis sampled len = %d, want 9", len(got))
	}
	if got[0] != "#440154" || got[8] != "#ece51b" {
		t.Errorf("viridis endpoints = %s..%s, want #440154..#ece51b", got[0], got[8])
	}
	if got := resolveSchemeColors("viridis", 3); len(got) != 3 {
		t.Errorf("viridis size 3 len = %d", len(got))
	}
	// Unknown scheme falls back to nivo.
	if got := resolveSchemeColors("not-a-scheme", 0); len(got) != len(CategoricalColorSchemes["nivo"]) {
		t.Errorf("unknown scheme fallback len = %d", len(got))
	}
}

// --- sequential / diverging / quantize scales ---------------------------------

func TestGetSequentialColorScale_TwoColors(t *testing.T) {
	cfg := SequentialColorScaleConfig{Colors: [2]string{"#000000", "#ffffff"}}
	s := GetSequentialColorScale(cfg, SequentialColorScaleValues{Min: 0, Max: 10})
	if got := s(0); got != "#000000" {
		t.Errorf("s(min) = %s, want #000000", got)
	}
	if got := s(10); got != "#ffffff" {
		t.Errorf("s(max) = %s, want #ffffff", got)
	}
	// Out-of-domain values clamp.
	if got := s(-5); got != "#000000" {
		t.Errorf("s(<min) = %s, want clamp to #000000", got)
	}
	if got := s(50); got != "#ffffff" {
		t.Errorf("s(>max) = %s, want clamp to #ffffff", got)
	}
	if got := s(5); !hexRe.MatchString(got) || got == "#000000" || got == "#ffffff" {
		t.Errorf("s(mid) = %s, want an intermediate hex", got)
	}
}

func TestGetSequentialColorScale_MinMaxOverride(t *testing.T) {
	min, max := 100.0, 200.0
	cfg := SequentialColorScaleConfig{Colors: [2]string{"#000000", "#ffffff"}, MinValue: &min, MaxValue: &max}
	s := GetSequentialColorScale(cfg, SequentialColorScaleValues{Min: 0, Max: 10})
	if got := s(100); got != "#000000" {
		t.Errorf("s(overridden min) = %s, want #000000", got)
	}
	if got := s(200); got != "#ffffff" {
		t.Errorf("s(overridden max) = %s, want #ffffff", got)
	}
}

func TestGetSequentialColorScale_CustomInterpolator(t *testing.T) {
	var seen []float64
	cfg := SequentialColorScaleConfig{Interpolator: func(tt float64) string {
		seen = append(seen, tt)
		return "#abcabc"
	}}
	s := GetSequentialColorScale(cfg, SequentialColorScaleValues{Min: 0, Max: 4})
	if got := s(1); got != "#abcabc" {
		t.Fatalf("custom interpolator output = %s", got)
	}
	if len(seen) != 1 || seen[0] != 0.25 {
		t.Errorf("interpolator saw t=%v, want [0.25]", seen)
	}
}

func TestGetSequentialColorScale_DefaultScheme(t *testing.T) {
	// Empty config uses the turbo default.
	s := GetSequentialColorScale(SequentialColorScaleConfig{}, SequentialColorScaleValues{Min: 0, Max: 1})
	if got, want := s(0), ColorInterpolators["turbo"](0); got != want {
		t.Errorf("default scheme s(0) = %s, want turbo(0) %s", got, want)
	}
}

func TestGetSequentialColorScale_ZeroSpan(t *testing.T) {
	// min == max must not divide by zero.
	cfg := SequentialColorScaleConfig{Colors: [2]string{"#000000", "#ffffff"}}
	s := GetSequentialColorScale(cfg, SequentialColorScaleValues{Min: 5, Max: 5})
	if got := s(5); got != "#000000" {
		t.Errorf("zero-span s(5) = %s, want #000000 (t=0)", got)
	}
}

func TestGetDivergingColorScale_ThreeColors(t *testing.T) {
	cfg := DivergingColorScaleConfig{Colors: [3]string{"#ff0000", "#ffffff", "#0000ff"}}
	s := GetDivergingColorScale(cfg, SequentialColorScaleValues{Min: -10, Max: 10})
	if got := s(-10); got != "#ff0000" {
		t.Errorf("s(min) = %s, want #ff0000", got)
	}
	if got := s(10); got != "#0000ff" {
		t.Errorf("s(max) = %s, want #0000ff", got)
	}
	if got := s(-50); got != "#ff0000" {
		t.Errorf("s(<min) = %s, want clamp", got)
	}
	if got := s(50); got != "#0000ff" {
		t.Errorf("s(>max) = %s, want clamp", got)
	}
}

func TestGetDivergingColorScale_DivergeAtAndOverrides(t *testing.T) {
	var seen []float64
	rec := func(tt float64) string {
		seen = append(seen, tt)
		return "#000001"
	}
	divergeAt := 0.25
	min, max := 0.0, 8.0
	cfg := DivergingColorScaleConfig{Interpolator: rec, DivergeAt: &divergeAt, MinValue: &min, MaxValue: &max}
	s := GetDivergingColorScale(cfg, SequentialColorScaleValues{Min: -1, Max: 1})
	s(0) // t=0, offset 0.5-0.25=0.25 → interp(0.25)
	s(8) // t=1 → interp(1.25)
	s(2) // t=0.25 → interp(0.5)
	want := []float64{0.25, 1.25, 0.5}
	if len(seen) != len(want) {
		t.Fatalf("interp saw %v, want %v", seen, want)
	}
	for i := range want {
		if !approxEqualF(seen[i], want[i]) {
			t.Errorf("interp t[%d] = %v, want %v", i, seen[i], want[i])
		}
	}
}

func TestGetDivergingColorScale_DefaultSchemeAndZeroSpan(t *testing.T) {
	s := GetDivergingColorScale(DivergingColorScaleConfig{}, SequentialColorScaleValues{Min: 3, Max: 3})
	// span==0 forces span=1: t=(3-3)/1=0, default divergeAt 0.5 → offset 0.
	want := ColorInterpolators[DivergingColorScaleDefaults.Scheme](0)
	if got := s(3); got != want {
		t.Errorf("default diverging s(3) = %s, want %s", got, want)
	}
}

func TestGetQuantizeColorScale_ExplicitColors(t *testing.T) {
	cfg := QuantizeColorScaleConfig{Colors: []string{"#a", "#b", "#c", "#d"}, Domain: [2]float64{0, 100}}
	s := GetQuantizeColorScale(cfg, SequentialColorScaleValues{})
	cases := []struct {
		v    float64
		want string
	}{
		{-5, "#a"}, {0, "#a"}, {10, "#a"}, {30, "#b"}, {50, "#c"}, {80, "#d"}, {100, "#d"}, {150, "#d"},
	}
	for _, c := range cases {
		if got := s(c.v); got != c.want {
			t.Errorf("quantize(%v) = %s, want %s", c.v, got, c.want)
		}
	}
}

func TestGetQuantizeColorScale_DefaultSchemeSteps(t *testing.T) {
	// No colors: 7 turbo steps over the values' domain.
	s := GetQuantizeColorScale(QuantizeColorScaleConfig{}, SequentialColorScaleValues{Min: 0, Max: 1})
	if got, want := s(-1), ColorInterpolators["turbo"](0); got != want {
		t.Errorf("quantize first bucket = %s, want turbo(0) %s", got, want)
	}
	if got, want := s(2), ColorInterpolators["turbo"](1); got != want {
		t.Errorf("quantize last bucket = %s, want turbo(1) %s", got, want)
	}
}

func TestGetQuantizeColorScale_CustomStepsAndScheme(t *testing.T) {
	cfg := QuantizeColorScaleConfig{Scheme: "viridis", Steps: 3}
	s := GetQuantizeColorScale(cfg, SequentialColorScaleValues{Min: 0, Max: 3})
	if got, want := s(-1), ColorInterpolators["viridis"](0); got != want {
		t.Errorf("first = %s, want viridis(0) %s", got, want)
	}
	if got, want := s(9), ColorInterpolators["viridis"](1); got != want {
		t.Errorf("last = %s, want viridis(1) %s", got, want)
	}
}

func TestNiceDomain(t *testing.T) {
	if got := niceDomain(0, 100); got != [2]float64{0, 100} {
		t.Errorf("niceDomain(0,100) = %v", got)
	}
	if got := niceDomain(5, 5); got != [2]float64{5, 5} {
		t.Errorf("niceDomain(5,5) = %v, want unchanged", got)
	}
}

func TestGetContinuousColorScale(t *testing.T) {
	vals := SequentialColorScaleValues{Min: 0, Max: 1}
	seq := &SequentialColorScaleConfig{Colors: [2]string{"#000000", "#ffffff"}}
	if got := GetContinuousColorScale(ContinuousColorScaleConfig{Type: "sequential", Sequential: seq}, vals)(0); got != "#000000" {
		t.Errorf("continuous sequential = %s", got)
	}
	div := &DivergingColorScaleConfig{Colors: [3]string{"#ff0000", "#ffffff", "#0000ff"}}
	if got := GetContinuousColorScale(ContinuousColorScaleConfig{Type: "diverging", Diverging: div}, vals)(0); got != "#ff0000" {
		t.Errorf("continuous diverging = %s", got)
	}
	qz := &QuantizeColorScaleConfig{Colors: []string{"#111", "#222"}, Domain: [2]float64{0, 1}}
	if got := GetContinuousColorScale(ContinuousColorScaleConfig{Type: "quantize", Quantize: qz}, vals)(0); got != "#111" {
		t.Errorf("continuous quantize = %s", got)
	}
	// Missing sub-config or unknown type falls back to black.
	if got := GetContinuousColorScale(ContinuousColorScaleConfig{Type: "sequential"}, vals)(0.5); got != "#000" {
		t.Errorf("missing sub-config = %s, want #000", got)
	}
	if got := GetContinuousColorScale(ContinuousColorScaleConfig{Type: "bogus"}, vals)(0.5); got != "#000" {
		t.Errorf("unknown type = %s, want #000", got)
	}
}

func TestIsContinuousColorScale(t *testing.T) {
	for _, typ := range []string{"sequential", "diverging", "quantize"} {
		if !IsContinuousColorScale(map[string]any{"type": typ}) {
			t.Errorf("type %q should be continuous", typ)
		}
	}
	if IsContinuousColorScale(map[string]any{"type": "ordinal"}) {
		t.Error("ordinal should not be continuous")
	}
	if IsContinuousColorScale("sequential") {
		t.Error("non-map config should not be continuous")
	}
}

func TestToIntHelper(t *testing.T) {
	if v, ok := toInt(7); !ok || v != 7 {
		t.Errorf("toInt(int) = %v %v", v, ok)
	}
	if v, ok := toInt(int64(7)); !ok || v != 7 {
		t.Errorf("toInt(int64) = %v %v", v, ok)
	}
	if v, ok := toInt(7.9); !ok || v != 7 {
		t.Errorf("toInt(float64) = %v %v", v, ok)
	}
	if _, ok := toInt("7"); ok {
		t.Error("toInt(string) should fail")
	}
}

func TestToFloatHelper(t *testing.T) {
	if v, ok := toFloat(1.5); !ok || v != 1.5 {
		t.Errorf("toFloat(float64) = %v %v", v, ok)
	}
	if v, ok := toFloat(2); !ok || v != 2 {
		t.Errorf("toFloat(int) = %v %v", v, ok)
	}
	if v, ok := toFloat(int64(3)); !ok || v != 3 {
		t.Errorf("toFloat(int64) = %v %v", v, ok)
	}
	if _, ok := toFloat("1.5"); ok {
		t.Error("toFloat(string) should fail")
	}
}
