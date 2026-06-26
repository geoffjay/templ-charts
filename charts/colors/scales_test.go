package colors

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/theming"
)

func TestCategoricalSchemes(t *testing.T) {
	if len(CategoricalColorSchemes["nivo"]) != 6 {
		t.Fatalf("nivo scheme len = %d, want 6", len(CategoricalColorSchemes["nivo"]))
	}
	if len(CategoricalColorSchemes["category10"]) != 10 {
		t.Fatalf("category10 len = %d, want 10", len(CategoricalColorSchemes["category10"]))
	}
}

func TestDivergingSchemes(t *testing.T) {
	brbg := DivergingColorSchemes["brown_blueGreen"]
	if len(brbg) != 9 {
		t.Fatalf("BrBG sizes = %d, want 9 (3..11)", len(brbg))
	}
	if len(brbg[11]) != 11 {
		t.Fatalf("BrBG[11] len = %d, want 11", len(brbg[11]))
	}
}

func TestSequentialSchemes(t *testing.T) {
	blues := SequentialColorSchemes["blues"]
	if len(blues[9]) != 9 {
		t.Fatalf("blues[9] len = %d, want 9", len(blues[9]))
	}
}

func TestIsCategoricalColorScheme(t *testing.T) {
	if !IsCategoricalColorScheme("nivo") {
		t.Fatal("nivo should be categorical")
	}
	if IsCategoricalColorScheme("blues") {
		t.Fatal("blues should not be categorical")
	}
}

func TestOrdinalColorScale_Scheme(t *testing.T) {
	cfg, err := ParseOrdinalColorScaleConfig(map[string]any{"scheme": "nivo"})
	if err != nil {
		t.Fatal(err)
	}
	scale := GetOrdinalColorScale[map[string]string](cfg, "id")
	d1 := map[string]string{"id": "a"}
	d2 := map[string]string{"id": "b"}
	if scale(d1) != scale(d1) {
		t.Fatal("same key should return same color")
	}
	if scale(d1) == scale(d2) {
		// different keys *could* collide on a 1-color scheme, but nivo has 6.
		t.Fatal("different keys should get different colors on the nivo scheme")
	}
}

func TestOrdinalColorScale_Static(t *testing.T) {
	cfg := OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: "#ff0000"}
	scale := GetOrdinalColorScale[any](cfg, nil)
	if scale(nil) != "#ff0000" {
		t.Fatalf("got %q, want #ff0000", scale(nil))
	}
}

func TestInheritedColorGenerator_Static(t *testing.T) {
	c := NewStaticColor("#abc")
	g := GetInheritedColorGenerator(c, nil)
	if g(nil) != "#abc" {
		t.Fatalf("got %q, want #abc", g(nil))
	}
}

func TestInheritedColorGenerator_Theme(t *testing.T) {
	c := NewThemeColor("text.fill")
	g := GetInheritedColorGenerator(c, &theming.DefaultTheme)
	if got := g(nil); got != theming.DefaultTheme.Text.Fill {
		t.Fatalf("got %q, want %q", got, theming.DefaultTheme.Text.Fill)
	}
}

func TestInheritedColorGenerator_FromContext(t *testing.T) {
	c := NewFromContextColor("color", []ColorModifier{{"darker", 1.0}})
	g := GetInheritedColorGenerator(c, &theming.DefaultTheme)
	d := map[string]any{"color": "#ffffff"}
	got := g(d)
	if got == "#ffffff" {
		t.Fatal("darker modifier should have changed the color")
	}
}

func TestSequentialColorScale(t *testing.T) {
	cfg := SequentialColorScaleConfig{Type: "sequential", Scheme: "blues"}
	scale := GetSequentialColorScale(cfg, SequentialColorScaleValues{Min: 0, Max: 100})
	if scale(0) != scale(0) {
		t.Fatal("same value should return same color")
	}
	if scale(0) == scale(100) {
		t.Fatal("min and max should differ")
	}
}

func TestInterpolateTurbo(t *testing.T) {
	if interpolateTurbo(0) == interpolateTurbo(1) {
		t.Fatal("turbo endpoints should differ")
	}
}

func TestApplyColorModifiers(t *testing.T) {
	got := ApplyColorModifiers("#ffffff", [][2]any{{"darker", 1.0}})
	if got == "#ffffff" {
		t.Fatal("darker should change white")
	}
}
