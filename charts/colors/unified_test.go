package colors

import "testing"

// TestSetResolvesToUnderlyingShapes is the plan's key guard: the unified Set
// entry point must resolve to configs that behave identically to building the
// underlying config shapes by hand.

func TestSetOrdinalMatchesScheme(t *testing.T) {
	// Set(PaletteID).Ordinal() == Scheme(PaletteID).
	got := Set(PaletteTableau10).Ordinal()
	want := Scheme(PaletteTableau10)
	if got.Type != want.Type || got.Scheme != want.Scheme {
		t.Fatalf("Set(palette).Ordinal() = %+v want %+v", got, want)
	}
	// And they produce the same colors through the scale.
	gs := GetOrdinalColorScale[string](got, "")
	ws := GetOrdinalColorScale[string](want, "")
	for _, k := range []string{"a", "b", "c", "d"} {
		if gs(k) != ws(k) {
			t.Errorf("scale color for %q differs: %s vs %s", k, gs(k), ws(k))
		}
	}
}

func TestSetOrdinalFromStringSchemeVsColor(t *testing.T) {
	// A recognized scheme id resolves to a scheme config.
	if got := Set("nivo").Ordinal(); got.Type != OrdinalTypeScheme || got.Scheme != "nivo" {
		t.Errorf("Set(\"nivo\") = %+v want scheme nivo", got)
	}
	// An unrecognized string resolves to a static color.
	if got := Set("#ff0000").Ordinal(); got.Type != OrdinalTypeStatic || got.Static != "#ff0000" {
		t.Errorf("Set(\"#ff0000\") = %+v want static #ff0000", got)
	}
}

func TestSetColorListMatchesPaletteColors(t *testing.T) {
	got := Set([]string{"#f00", "#0f0", "#00f"}).Ordinal()
	want := PaletteColors("#f00", "#0f0", "#00f")
	if got.Type != want.Type || len(got.Colors) != len(want.Colors) {
		t.Fatalf("Set(list).Ordinal() = %+v want %+v", got, want)
	}
	for i := range got.Colors {
		if got.Colors[i] != want.Colors[i] {
			t.Errorf("color[%d] = %s want %s", i, got.Colors[i], want.Colors[i])
		}
	}
}

func TestSetInheritedMatchesConstructors(t *testing.T) {
	// Static color → NewStaticColor.
	if got, want := Set("#123456").Inherited(), NewStaticColor("#123456"); got.Type != want.Type || got.Static != want.Static {
		t.Errorf("Set(static).Inherited() = %+v want %+v", got, want)
	}
	// From-theme → NewThemeColor.
	if got, want := SetFromTheme("labels.text.fill").Inherited(), NewThemeColor("labels.text.fill"); got.Type != want.Type || got.ThemePath != want.ThemePath {
		t.Errorf("SetFromTheme().Inherited() = %+v want %+v", got, want)
	}
	// From-datum + modifiers → NewFromContextColorInSpace, producing the same
	// resolved color as the hand-built config.
	mods := []ColorModifier{{"darker", 1.0}}
	got := Set(SetFromDatum("color").WithModifiers(mods...).InSpace(SpaceLab)).Inherited()
	want := NewFromContextColorInSpace("color", mods, SpaceLab)
	datum := map[string]any{"color": "#ff0000"}
	gg := GetInheritedColorGenerator(got, nil)
	wg := GetInheritedColorGenerator(want, nil)
	if gg(datum) != wg(datum) {
		t.Errorf("inherited from-datum color differs: %s vs %s", gg(datum), wg(datum))
	}
	if gg(datum) != "#c30000" { // Lab darker(1) on red, from the color-space port
		t.Errorf("resolved from-datum Lab darker = %s want #c30000", gg(datum))
	}
}

func TestSetSequentialAndDiverging(t *testing.T) {
	// Scheme + InSpace flows into the scale Space and matches a hand-built config.
	got := Set(PaletteViridis).InSpace(SpaceLab).Sequential()
	want := SequentialColorScaleConfig{Type: "sequential", Scheme: "viridis", Space: SpaceLab}
	if got.Type != want.Type || got.Scheme != want.Scheme || got.Space != want.Space {
		t.Fatalf("Set(viridis).InSpace(Lab).Sequential() = %+v want %+v", got, want)
	}
	vals := SequentialColorScaleValues{Min: 0, Max: 100}
	gs := GetSequentialColorScale(got, vals)
	ws := GetSequentialColorScale(want, vals)
	for i := 0; i <= 10; i++ {
		v := float64(i) * 10
		if gs(v) != ws(v) {
			t.Errorf("sequential color at %.0f differs: %s vs %s", v, gs(v), ws(v))
		}
	}

	// Two-color diverging.
	d := Set([]string{"#ff0000", "#ffffff", "#0000ff"}).Diverging()
	if d.Type != "diverging" || d.Colors != [3]string{"#ff0000", "#ffffff", "#0000ff"} {
		t.Errorf("Set(3 colors).Diverging() = %+v", d)
	}
}

func TestSetPassthroughConfigs(t *testing.T) {
	// A pre-built config passes through its matching resolver verbatim.
	oc := OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: "set2"}
	if got := Set(oc).Ordinal(); got.Type != oc.Type || got.Scheme != oc.Scheme {
		t.Errorf("Set(ordinalCfg).Ordinal() = %+v want %+v", got, oc)
	}
	ic := NewStaticColor("#abcdef")
	if got := Set(ic).Inherited(); got.Type != ic.Type || got.Static != ic.Static {
		t.Errorf("Set(inheritedCfg).Inherited() = %+v want %+v", got, ic)
	}
}
