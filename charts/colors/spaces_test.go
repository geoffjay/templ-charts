package colors

import "testing"

// TestInterpolateInSpaceRGBIsLegacy is the byte-stability guard: SpaceRGB must
// route through the exact legacy interpolateRgbBasis, so no existing scale or
// palette golden can move when the Space selector is added.
func TestInterpolateInSpaceRGBIsLegacy(t *testing.T) {
	cols := []string{"#3070b0", "#f0a030", "#00ff00"}
	legacy := interpolateRgbBasis(cols)
	got := interpolateInSpace(cols, SpaceRGB)
	for i := 0; i <= 20; i++ {
		tt := float64(i) / 20
		if got(tt) != legacy(tt) {
			t.Fatalf("SpaceRGB drifted from interpolateRgbBasis at t=%.2f: %s vs %s", tt, got(tt), legacy(tt))
		}
	}
}

// TestSequentialScaleSpaceDefaultRGB confirms a zero-value Space and an
// explicit SpaceRGB yield identical sequential scales.
func TestSequentialScaleSpaceDefaultRGB(t *testing.T) {
	base := SequentialColorScaleConfig{Type: "sequential", Scheme: "viridis"}
	explicit := SequentialColorScaleConfig{Type: "sequential", Scheme: "viridis", Space: SpaceRGB}
	vals := SequentialColorScaleValues{Min: 0, Max: 100}
	a := GetSequentialColorScale(base, vals)
	b := GetSequentialColorScale(explicit, vals)
	for i := 0; i <= 10; i++ {
		v := float64(i) * 10
		if a(v) != b(v) {
			t.Fatalf("default Space != SpaceRGB at v=%.0f: %s vs %s", v, a(v), b(v))
		}
	}
}

// TestSequentialScaleLab checks a two-color Lab sequential scale interpolates
// perceptually — matching d3's interpolateLab exactly at sampled points, and
// differing from the RGB path.
func TestSequentialScaleLab(t *testing.T) {
	vals := SequentialColorScaleValues{Min: 0, Max: 1}
	lab := GetSequentialColorScale(SequentialColorScaleConfig{
		Type: "sequential", Colors: [2]string{"#ff0000", "#0000ff"}, Space: SpaceLab,
	}, vals)
	if got := lab(0.5); got != "#c10088" { // d3 interpolateLab(red,blue)(0.5)
		t.Errorf("Lab seq red->blue @0.5 = %s want #c10088", got)
	}

	lab2 := GetSequentialColorScale(SequentialColorScaleConfig{
		Type: "sequential", Colors: [2]string{"#3070b0", "#f0a030"}, Space: SpaceLab,
	}, vals)
	if got := lab2(0.3); got != "#837e8f" { // d3 interpolateLab @0.3 = rgb(131,126,143)
		t.Errorf("Lab seq @0.3 = %s want #837e8f", got)
	}

	rgb := GetSequentialColorScale(SequentialColorScaleConfig{
		Type: "sequential", Colors: [2]string{"#ff0000", "#0000ff"},
	}, vals)
	if rgb(0.5) == lab(0.5) {
		t.Errorf("RGB and Lab midpoints should differ, both = %s", lab(0.5))
	}
}

// TestModifiersInSpace pins RGB (legacy) and Lab/Lch (perceptual) lightness
// modifiers against d3-color.
func TestModifiersInSpace(t *testing.T) {
	mods := [][2]any{{"darker", 1.0}}
	if got := ApplyColorModifiers("#ff0000", mods); got != "#b30000" {
		t.Errorf("RGB darker(1) red = %s want #b30000", got)
	}
	if got := ApplyColorModifiersInSpace("#ff0000", mods, SpaceRGB); got != "#b30000" {
		t.Errorf("SpaceRGB darker(1) red = %s want #b30000 (must equal legacy)", got)
	}
	if got := ApplyColorModifiersInSpace("#ff0000", mods, SpaceLab); got != "#c30000" {
		t.Errorf("Lab darker(1) red = %s want #c30000", got)
	}
	if got := ApplyColorModifiersInSpace("#ff0000", mods, SpaceLch); got != "#c30000" {
		t.Errorf("Lch darker(1) red = %s want #c30000", got)
	}
	bright := [][2]any{{"brighter", 0.5}}
	if got := ApplyColorModifiersInSpace("#123456", bright, SpaceLab); got != "#2a486c" {
		t.Errorf("Lab brighter(0.5) #123456 = %s want #2a486c", got)
	}
}

// TestInheritedColorModifierSpace confirms the ModifierSpace field flows into
// the generated inherited-color func, and defaults (zero value) to RGB.
func TestInheritedColorModifierSpace(t *testing.T) {
	mods := []ColorModifier{{"darker", 1.0}}
	datum := map[string]any{"color": "#ff0000"}

	rgbGen := GetInheritedColorGenerator(NewFromContextColor("color", mods), nil)
	if got := rgbGen(datum); got != "#b30000" {
		t.Errorf("default (RGB) inherited darker(1) = %s want #b30000", got)
	}
	labGen := GetInheritedColorGenerator(NewFromContextColorInSpace("color", mods, SpaceLab), nil)
	if got := labGen(datum); got != "#c30000" {
		t.Errorf("Lab inherited darker(1) = %s want #c30000", got)
	}
}

// TestSwatchInSpace confirms perceptual swatches keep their endpoints, differ
// from the RGB swatch, and that categorical/SpaceRGB fall back to Swatch.
func TestSwatchInSpace(t *testing.T) {
	p, ok := LookupPalette(PaletteViridis)
	if !ok {
		t.Fatal("viridis palette missing")
	}
	rgbSw := p.Swatch(7)
	labSw := p.SwatchIn(7, SpaceLab)
	if len(labSw) != 7 {
		t.Fatalf("SwatchIn len = %d want 7", len(labSw))
	}
	// Endpoints are the scheme's extremes and independent of interpolation space.
	if labSw[0] != rgbSw[0] || labSw[6] != rgbSw[6] {
		t.Errorf("Lab swatch endpoints drifted: %v vs %v", []string{labSw[0], labSw[6]}, []string{rgbSw[0], rgbSw[6]})
	}
	// Interior stops should differ between RGB and Lab sampling.
	same := true
	for i := 1; i < 6; i++ {
		if labSw[i] != rgbSw[i] {
			same = false
		}
	}
	if same {
		t.Errorf("Lab swatch identical to RGB swatch: %v", labSw)
	}
	// SpaceRGB and categorical fall back to Swatch.
	for i, c := range p.SwatchIn(7, SpaceRGB) {
		if c != rgbSw[i] {
			t.Errorf("SwatchIn(SpaceRGB)[%d] = %s want %s", i, c, rgbSw[i])
		}
	}
	cat, _ := LookupPalette(PaletteCategory10)
	if got, want := cat.SwatchIn(5, SpaceLab), cat.Swatch(5); !equalStrs(got, want) {
		t.Errorf("categorical SwatchIn(Lab) = %v want %v (== Swatch)", got, want)
	}
}

func equalStrs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
