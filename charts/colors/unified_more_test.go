package colors

import "testing"

// Tests for the typed Set* constructors and the remaining resolver branches of
// ColorSetting (Sequential/Diverging/Static and passthrough sources).

func TestSetSchemeConstructor(t *testing.T) {
	got := SetScheme(PaletteViridis).Ordinal()
	if got.Type != OrdinalTypeScheme || got.Scheme != "viridis" {
		t.Fatalf("SetScheme(viridis).Ordinal() = %+v, want scheme viridis", got)
	}
}

func TestSetColorListConstructor(t *testing.T) {
	got := SetColorList("#111111", "#222222").Ordinal()
	if got.Type != OrdinalTypeColors || len(got.Colors) != 2 || got.Colors[0] != "#111111" {
		t.Fatalf("SetColorList(...).Ordinal() = %+v, want colors list", got)
	}
}

func TestSetStaticConstructor(t *testing.T) {
	s := SetStatic("#abcdef")
	if got := s.Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "#abcdef" {
		t.Errorf("SetStatic.Ordinal() = %+v, want static #abcdef", got)
	}
	if got := s.Static(); got != "#abcdef" {
		t.Errorf("SetStatic.Static() = %q, want #abcdef", got)
	}
}

func TestSetFuncConstructor(t *testing.T) {
	fn := func(d any) string { return "#00ff00" }
	cfg := SetFunc(fn).Ordinal()
	if cfg.Type != OrdinalTypeFunc || cfg.Func == nil {
		t.Fatalf("SetFunc.Ordinal() = %+v, want func config", cfg)
	}
	if got := cfg.Func(nil); got != "#00ff00" {
		t.Errorf("func config returned %q, want #00ff00", got)
	}
}

func TestSetNilAndUnknown(t *testing.T) {
	if got := Set(nil).Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "" {
		t.Errorf("Set(nil).Ordinal() = %+v, want empty static", got)
	}
	if got := Set(42).Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "" {
		t.Errorf("Set(42).Ordinal() = %+v, want empty static", got)
	}
}

func TestSetPassthroughColorSetting(t *testing.T) {
	orig := SetStatic("#111111")
	if got := Set(orig).Static(); got != "#111111" {
		t.Errorf("Set(ColorSetting).Static() = %q, want #111111", got)
	}
}

func TestSetPalette(t *testing.T) {
	p, ok := LookupPalette(PaletteViridis)
	if !ok {
		t.Fatal("viridis palette not in catalog")
	}
	if got := Set(p).Ordinal(); got.Type != OrdinalTypeScheme || got.Scheme != "viridis" {
		t.Errorf("Set(Palette).Ordinal() = %+v, want scheme viridis", got)
	}
}

func TestSetPrebuiltConfigs(t *testing.T) {
	// Ordinal passthrough.
	oc := OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: []string{"#123456"}}
	if got := Set(oc).Ordinal(); got.Type != OrdinalTypeColors || got.Colors[0] != "#123456" {
		t.Errorf("Set(ordinal).Ordinal() = %+v, want passthrough", got)
	}
	// Sequential with scheme resolves to an ordinal scheme.
	sc := SequentialColorScaleConfig{Scheme: "blues"}
	if got := Set(sc).Ordinal(); got.Type != OrdinalTypeScheme || got.Scheme != "blues" {
		t.Errorf("Set(sequential).Ordinal() = %+v, want scheme blues", got)
	}
	// Sequential without scheme falls back to empty static.
	if got := Set(SequentialColorScaleConfig{}).Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "" {
		t.Errorf("Set(sequential no scheme).Ordinal() = %+v, want empty static", got)
	}
	// Diverging with scheme.
	dc := DivergingColorScaleConfig{Scheme: "spectral"}
	if got := Set(dc).Ordinal(); got.Type != OrdinalTypeScheme || got.Scheme != "spectral" {
		t.Errorf("Set(diverging).Ordinal() = %+v, want scheme spectral", got)
	}
	if got := Set(DivergingColorScaleConfig{}).Ordinal(); got.Type != OrdinalTypeStatic {
		t.Errorf("Set(diverging no scheme).Ordinal() = %+v, want static fallback", got)
	}
	// Inherited static resolves to a static ordinal.
	ic := NewStaticColor("#0f0f0f")
	if got := Set(ic).Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "#0f0f0f" {
		t.Errorf("Set(inherited static).Ordinal() = %+v, want static #0f0f0f", got)
	}
	// Inherited non-static falls back to empty static.
	if got := Set(NewThemeColor("Text.Fill")).Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "" {
		t.Errorf("Set(inherited theme).Ordinal() = %+v, want empty static", got)
	}
}

func TestSequentialResolution(t *testing.T) {
	// Scheme intent.
	got := SetScheme(PaletteViridis).Sequential()
	if got.Type != "sequential" || got.Scheme != "viridis" {
		t.Errorf("SetScheme.Sequential() = %+v, want scheme viridis", got)
	}
	// Color list uses first + last.
	got = SetColorList("#111111", "#222222", "#333333").Sequential()
	if got.Colors != [2]string{"#111111", "#333333"} {
		t.Errorf("SetColorList.Sequential() colors = %v, want [#111111 #333333]", got.Colors)
	}
	// Empty color list yields empty pair.
	if got := SetColorList().Sequential(); got.Colors != [2]string{"", ""} {
		t.Errorf("empty SetColorList.Sequential() colors = %v, want empty", got.Colors)
	}
	// Pre-built config passes through, space applied.
	pre := SequentialColorScaleConfig{Scheme: "blues"}
	got = Set(pre).InSpace(SpaceLab).Sequential()
	if got.Scheme != "blues" || got.Space != SpaceLab || got.Type != "sequential" {
		t.Errorf("Set(sequential).InSpace(Lab).Sequential() = %+v", got)
	}
	// Non-continuous intent falls back to bare sequential config.
	got = SetStatic("#ff0000").Sequential()
	if got.Type != "sequential" || got.Scheme != "" || got.Colors != [2]string{"", ""} {
		t.Errorf("SetStatic.Sequential() = %+v, want bare config", got)
	}
}

func TestDivergingResolution(t *testing.T) {
	got := SetScheme(PaletteSpectral).Diverging()
	if got.Type != "diverging" || got.Scheme != "spectral" {
		t.Errorf("SetScheme.Diverging() = %+v, want scheme spectral", got)
	}
	// Three colors map positionally.
	got = SetColorList("#111111", "#222222", "#333333").Diverging()
	if got.Colors != [3]string{"#111111", "#222222", "#333333"} {
		t.Errorf("SetColorList(3).Diverging() colors = %v", got.Colors)
	}
	// Two colors leave the third empty.
	got = SetColorList("#111111", "#222222").Diverging()
	if got.Colors != [3]string{"#111111", "#222222", ""} {
		t.Errorf("SetColorList(2).Diverging() colors = %v", got.Colors)
	}
	// Pre-built config passes through with space.
	pre := DivergingColorScaleConfig{Scheme: "spectral"}
	got = Set(pre).InSpace(SpaceLch).Diverging()
	if got.Scheme != "spectral" || got.Space != SpaceLch {
		t.Errorf("Set(diverging).InSpace(Lch).Diverging() = %+v", got)
	}
	// Fallback.
	got = SetStatic("#ff0000").Diverging()
	if got.Type != "diverging" || got.Scheme != "" {
		t.Errorf("SetStatic.Diverging() = %+v, want bare config", got)
	}
}

func TestStaticResolution(t *testing.T) {
	if got := SetStatic("#123456").Static(); got != "#123456" {
		t.Errorf("static = %q, want #123456", got)
	}
	if got := SetColorList("#111111", "#222222").Static(); got != "#111111" {
		t.Errorf("color-list static = %q, want first color", got)
	}
	if got := SetColorList().Static(); got != "" {
		t.Errorf("empty color-list static = %q, want empty", got)
	}
	if got := Set(NewStaticColor("#aabbcc")).Static(); got != "#aabbcc" {
		t.Errorf("inherited static = %q, want #aabbcc", got)
	}
	if got := Set(NewThemeColor("Text.Fill")).Static(); got != "" {
		t.Errorf("inherited theme static = %q, want empty fallback", got)
	}
	if got := SetScheme(PaletteNivo).Static(); got != "" {
		t.Errorf("scheme static = %q, want empty fallback", got)
	}
}

func TestInheritedResolution(t *testing.T) {
	// Static.
	if got := SetStatic("#ff0000").Inherited(); got.Type != InheritedColorTypeStatic || got.Static != "#ff0000" {
		t.Errorf("static inherited = %+v", got)
	}
	// Func.
	got := SetFunc(func(any) string { return "#00ff00" }).Inherited()
	if got.Type != InheritedColorTypeFunc || got.Func(nil) != "#00ff00" {
		t.Errorf("func inherited = %+v", got)
	}
	// Theme path.
	if got := SetFromTheme("Labels.Text.Fill").Inherited(); got.Type != InheritedColorTypeTheme || got.ThemePath != "Labels.Text.Fill" {
		t.Errorf("theme inherited = %+v", got)
	}
	// Datum path with modifiers + space.
	got = SetFromDatum("color").WithModifiers(ColorModifier{"darker", 1.0}).InSpace(SpaceLab).Inherited()
	if got.Type != InheritedColorTypeFromContext || got.FromPath != "color" {
		t.Errorf("datum inherited = %+v", got)
	}
	if len(got.Modifiers) != 1 || got.ModifierSpace != SpaceLab {
		t.Errorf("datum inherited modifiers/space = %+v", got)
	}
	// Pre-built inherited passthrough, modifiers/space overridden.
	pre := NewFromContextColor("data.color", nil)
	got = Set(pre).WithModifiers(ColorModifier{"brighter", 0.5}).InSpace(SpaceLch).Inherited()
	if got.FromPath != "data.color" || len(got.Modifiers) != 1 || got.ModifierSpace != SpaceLch {
		t.Errorf("prebuilt inherited with overrides = %+v", got)
	}
	// Color list falls back to first color as static.
	if got := SetColorList("#0000ff", "#00ff00").Inherited(); got.Type != InheritedColorTypeStatic || got.Static != "#0000ff" {
		t.Errorf("color-list inherited = %+v", got)
	}
	// Empty color list falls back to empty static.
	if got := SetColorList().Inherited(); got.Type != InheritedColorTypeStatic || got.Static != "" {
		t.Errorf("empty color-list inherited = %+v", got)
	}
	// Scheme (no inherited meaning) falls back to empty static.
	if got := SetScheme(PaletteNivo).Inherited(); got.Type != InheritedColorTypeStatic || got.Static != "" {
		t.Errorf("scheme inherited = %+v", got)
	}
}

func TestOrdinalFromDatumAndFunc(t *testing.T) {
	if got := SetFromDatum("data.color").Ordinal(); got.Type != OrdinalTypeDatum || got.DatumPath != "data.color" {
		t.Errorf("SetFromDatum.Ordinal() = %+v", got)
	}
	got := SetFunc(func(any) string { return "#fff" }).Ordinal()
	if got.Type != OrdinalTypeFunc || got.Func == nil {
		t.Errorf("SetFunc.Ordinal() = %+v", got)
	}
}
