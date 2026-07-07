package colors

import (
	"sort"
	"testing"
)

func TestPaletteKindString(t *testing.T) {
	cases := []struct {
		kind PaletteKind
		want string
	}{
		{KindCategorical, "categorical"},
		{KindSequential, "sequential"},
		{KindDiverging, "diverging"},
		{PaletteKind(99), "categorical"}, // unknown kinds default to categorical
	}
	for _, c := range cases {
		if got := c.kind.String(); got != c.want {
			t.Errorf("PaletteKind(%d).String() = %q, want %q", c.kind, got, c.want)
		}
	}
}

func TestPaletteConfigConstructors(t *testing.T) {
	p, ok := LookupPalette(PaletteViridis)
	if !ok {
		t.Fatal("viridis not in catalog")
	}
	if got := p.Ordinal(); got.Type != OrdinalTypeScheme || got.Scheme != "viridis" {
		t.Errorf("Ordinal() = %+v", got)
	}
	if got := p.Sequential(); got.Type != "sequential" || got.Scheme != "viridis" {
		t.Errorf("Sequential() = %+v", got)
	}
	d, _ := LookupPalette(PaletteSpectral)
	if got := d.Diverging(); got.Type != "diverging" || got.Scheme != "spectral" {
		t.Errorf("Diverging() = %+v", got)
	}
}

func TestPalettesByKind(t *testing.T) {
	total := 0
	for _, kind := range []PaletteKind{KindCategorical, KindSequential, KindDiverging} {
		ps := PalettesByKind(kind)
		if len(ps) == 0 {
			t.Errorf("no palettes of kind %v", kind)
		}
		for _, p := range ps {
			if p.Kind != kind {
				t.Errorf("palette %q kind = %v, want %v", p.ID, p.Kind, kind)
			}
		}
		total += len(ps)
	}
	if got := len(Palettes()); got != total {
		t.Errorf("kinds sum to %d, catalog has %d", total, got)
	}
}

func TestPaletteIDs(t *testing.T) {
	ids := PaletteIDs()
	if len(ids) != len(Palettes()) {
		t.Fatalf("PaletteIDs len = %d, want %d", len(ids), len(Palettes()))
	}
	if !sort.StringsAreSorted(ids) {
		t.Errorf("PaletteIDs not sorted: %v", ids)
	}
	found := false
	for _, id := range ids {
		if id == "viridis" {
			found = true
		}
	}
	if !found {
		t.Errorf("PaletteIDs missing viridis")
	}
}

func TestLookupPaletteNotFound(t *testing.T) {
	if _, ok := LookupPalette("not-a-palette"); ok {
		t.Error("unknown palette should not be found")
	}
}

func TestColorblindSafePalettes(t *testing.T) {
	safe := ColorblindSafePalettes()
	if len(safe) == 0 {
		t.Fatal("expected some colorblind-safe palettes")
	}
	for _, p := range safe {
		if !p.ColorblindSafe {
			t.Errorf("palette %q not flagged colorblind-safe", p.ID)
		}
	}
}

func TestPaletteSwatchEdgeCases(t *testing.T) {
	cat, _ := LookupPalette(PaletteNivo)
	// n <= 0 yields nil.
	if got := cat.Swatch(0); got != nil {
		t.Errorf("Swatch(0) = %v, want nil", got)
	}
	// Categorical: base colors, cycled past the set size.
	base := CategoricalColorSchemes["nivo"]
	got := cat.Swatch(len(base) + 2)
	if len(got) != len(base)+2 {
		t.Fatalf("Swatch len = %d", len(got))
	}
	if got[0] != base[0] || got[len(base)] != base[0] {
		t.Errorf("categorical swatch should cycle: %v", got)
	}
	// Unknown categorical id yields nil.
	if got := (Palette{ID: "nope", Kind: KindCategorical}).Swatch(3); got != nil {
		t.Errorf("unknown categorical Swatch = %v, want nil", got)
	}
	// Sequential: sampled gradient endpoints.
	seq, _ := LookupPalette(PaletteViridis)
	sw := seq.Swatch(5)
	if len(sw) != 5 || sw[0] != "#440154" || sw[4] != "#ece51b" {
		t.Errorf("viridis Swatch(5) = %v", sw)
	}
	// n == 1 samples the middle.
	one := seq.Swatch(1)
	if len(one) != 1 || one[0] != ColorInterpolators["viridis"](0.5) {
		t.Errorf("Swatch(1) = %v, want mid sample", one)
	}
	// Unknown gradient id yields nil.
	if got := (Palette{ID: "nope", Kind: KindSequential}).Swatch(3); got != nil {
		t.Errorf("unknown sequential Swatch = %v, want nil", got)
	}
}

func TestPaletteSwatchIn(t *testing.T) {
	seq, _ := LookupPalette(PaletteViridis)
	// SpaceRGB falls back to Swatch exactly.
	rgb := seq.SwatchIn(4, SpaceRGB)
	plain := seq.Swatch(4)
	for i := range plain {
		if rgb[i] != plain[i] {
			t.Errorf("SwatchIn(RGB)[%d] = %s, want Swatch %s", i, rgb[i], plain[i])
		}
	}
	// Categorical palettes ignore the space.
	cat, _ := LookupPalette(PaletteNivo)
	catLab := cat.SwatchIn(3, SpaceLab)
	catPlain := cat.Swatch(3)
	for i := range catPlain {
		if catLab[i] != catPlain[i] {
			t.Errorf("categorical SwatchIn[%d] = %s, want %s", i, catLab[i], catPlain[i])
		}
	}
	// Lab swatch: right length, valid hexes, n<=0 nil, n==1 mid.
	if got := seq.SwatchIn(0, SpaceLab); got != nil {
		t.Errorf("SwatchIn(0, Lab) = %v, want nil", got)
	}
	lab := seq.SwatchIn(5, SpaceLab)
	if len(lab) != 5 {
		t.Fatalf("SwatchIn(5, Lab) len = %d", len(lab))
	}
	for i, c := range lab {
		if !hexRe.MatchString(c) {
			t.Errorf("SwatchIn(Lab)[%d] = %q, not a hex color", i, c)
		}
	}
	if got := seq.SwatchIn(1, SpaceLab); len(got) != 1 || !hexRe.MatchString(got[0]) {
		t.Errorf("SwatchIn(1, Lab) = %v", got)
	}
	// Diverging palettes use their scheme stop arrays.
	div, _ := LookupPalette(PaletteSpectral)
	dl := div.SwatchIn(3, SpaceLch)
	if len(dl) != 3 {
		t.Fatalf("diverging SwatchIn len = %d", len(dl))
	}
	stops := DivergingColorSchemes["spectral"][11]
	if dl[0] != stops[0] || dl[2] != stops[len(stops)-1] {
		t.Errorf("diverging SwatchIn endpoints = %s..%s, want %s..%s", dl[0], dl[2], stops[0], stops[len(stops)-1])
	}
}

func TestSchemeStops(t *testing.T) {
	if got := schemeStops("spectral"); len(got) != 11 {
		t.Errorf("spectral stops len = %d, want 11", len(got))
	}
	if got := schemeStops("blues"); len(got) != 9 {
		t.Errorf("blues stops len = %d, want 9", len(got))
	}
	if got := schemeStops("viridis"); got != nil {
		t.Errorf("function-based scheme stops = %v, want nil", got)
	}
}

func TestSchemeInterpolatorInSpaceFallbacks(t *testing.T) {
	// Unknown scheme in RGB falls back to turbo.
	f := schemeInterpolatorInSpace("not-a-scheme", SpaceRGB)
	if got, want := f(0), ColorInterpolators["turbo"](0); got != want {
		t.Errorf("unknown scheme RGB = %s, want turbo %s", got, want)
	}
	// Unknown scheme in Lab also falls back to turbo.
	f = schemeInterpolatorInSpace("not-a-scheme", SpaceLab)
	if got, want := f(0), ColorInterpolators["turbo"](0); got != want {
		t.Errorf("unknown scheme Lab = %s, want turbo %s", got, want)
	}
	// Function-based scheme in Lab is densely sampled: endpoints preserved.
	f = schemeInterpolatorInSpace("viridis", SpaceLab)
	if got := f(0); got != "#440154" {
		t.Errorf("viridis Lab f(0) = %s, want #440154", got)
	}
	if got := f(1); got != "#ece51b" {
		t.Errorf("viridis Lab f(1) = %s, want #ece51b", got)
	}
}

func TestInterpolateInSpaceEdgeCases(t *testing.T) {
	// Empty stop list → black.
	f := interpolateInSpace(nil, SpaceLab)
	if got := f(0.5); got != "#000000" {
		t.Errorf("empty stops = %s, want #000000", got)
	}
	// Single stop → constant.
	f = interpolateInSpace([]string{"#123456"}, SpaceLab)
	if got := f(0.7); got != "#123456" {
		t.Errorf("single stop = %s, want #123456", got)
	}
	// Multi-stop endpoints round-trip exactly and clamp.
	f = interpolateInSpace([]string{"#ff0000", "#00ff00", "#0000ff"}, SpaceLab)
	if got := f(-1); got != "#ff0000" {
		t.Errorf("f(-1) = %s, want first stop", got)
	}
	if got := f(2); got != "#0000ff" {
		t.Errorf("f(2) = %s, want last stop", got)
	}
	if got := f(0.5); got != "#00ff00" {
		t.Errorf("f(0.5) = %s, want middle stop", got)
	}
}
