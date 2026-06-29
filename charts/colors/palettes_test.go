package colors

import (
	"regexp"
	"testing"
)

var hexRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// TestCuratedCategoricalSchemesValid checks every curated categorical palette
// is non-empty and contains only well-formed hex colors.
func TestCuratedCategoricalSchemesValid(t *testing.T) {
	for id, cols := range curatedCategoricalSchemes {
		if len(cols) == 0 {
			t.Errorf("curated categorical %q is empty", id)
		}
		for i, c := range cols {
			if !hexRe.MatchString(c) {
				t.Errorf("curated categorical %q[%d] = %q is not a #rrggbb hex", id, i, c)
			}
		}
		// Registered into the shared categorical map at init().
		if !IsCategoricalColorScheme(id) {
			t.Errorf("curated categorical %q not registered as a categorical scheme", id)
		}
	}
}

// TestCuratedContinuousRegistered checks curated sequential/diverging gradients
// resolve to interpolators producing valid colors across [0,1].
func TestCuratedContinuousRegistered(t *testing.T) {
	ids := append(append([]string{}, curatedSequentialIds...), curatedDivergingIds...)
	for _, id := range ids {
		interp := ColorInterpolators[id]
		if interp == nil {
			t.Errorf("curated continuous %q has no registered interpolator", id)
			continue
		}
		for _, tv := range []float64{0, 0.25, 0.5, 0.75, 1} {
			if got := interp(tv); got == "" {
				t.Errorf("interpolator %q returned empty color at t=%v", id, tv)
			}
		}
	}
}

// TestPaletteCatalogResolves checks every catalog palette resolves to a usable
// color source for its kind and that lookups are consistent.
func TestPaletteCatalogResolves(t *testing.T) {
	for _, p := range Palettes() {
		got, ok := LookupPalette(p.ID)
		if !ok || got.ID != p.ID {
			t.Errorf("LookupPalette(%q) failed", p.ID)
		}
		switch p.Kind {
		case KindCategorical:
			if !IsCategoricalColorScheme(string(p.ID)) {
				t.Errorf("categorical palette %q does not resolve to a categorical scheme", p.ID)
			}
		case KindSequential, KindDiverging:
			if ColorInterpolators[string(p.ID)] == nil {
				t.Errorf("%s palette %q has no interpolator", p.Kind, p.ID)
			}
		}
	}
}

// TestPaletteSwatch checks Swatch returns exactly n colors for each kind, and
// cycles categorical colors when n exceeds the base set.
func TestPaletteSwatch(t *testing.T) {
	cat, _ := LookupPalette(PaletteTableau20)
	if got := cat.Swatch(0); got != nil {
		t.Errorf("Swatch(0) = %v, want nil", got)
	}
	base := CategoricalColorSchemes[string(PaletteTableau20)]
	big := cat.Swatch(len(base) + 3)
	if len(big) != len(base)+3 {
		t.Fatalf("Swatch(len+3) returned %d colors, want %d", len(big), len(base)+3)
	}
	if big[0] != base[0] || big[len(base)] != base[0] {
		t.Errorf("categorical Swatch did not cycle: big[0]=%q big[len]=%q base[0]=%q", big[0], big[len(base)], base[0])
	}

	seq, _ := LookupPalette(PaletteSunset)
	sw := seq.Swatch(5)
	if len(sw) != 5 {
		t.Fatalf("sequential Swatch(5) returned %d colors", len(sw))
	}
	for i, c := range sw {
		if c == "" {
			t.Errorf("sequential swatch[%d] empty", i)
		}
	}
}

// TestSchemeConstructor checks the ergonomic constructors build the expected
// configs.
func TestSchemeConstructor(t *testing.T) {
	cfg := Scheme(PaletteObservable10)
	if cfg.Type != OrdinalTypeScheme || cfg.Scheme != "observable10" {
		t.Errorf("Scheme(observable10) = %+v", cfg)
	}
	cc := PaletteColors("#f00", "#0f0", "#00f")
	if cc.Type != OrdinalTypeColors || len(cc.Colors) != 3 {
		t.Errorf("PaletteColors(...) = %+v", cc)
	}
}

// TestColorblindSafeFlagged checks the known colorblind-safe palettes are
// flagged and surfaced by the helper.
func TestColorblindSafeFlagged(t *testing.T) {
	want := map[PaletteID]bool{
		PaletteOkabeIto:   true,
		PaletteTolVibrant: true,
		PaletteTolMuted:   true,
		PaletteViridis:    true,
	}
	got := map[PaletteID]bool{}
	for _, p := range ColorblindSafePalettes() {
		got[p.ID] = true
	}
	for id := range want {
		if !got[id] {
			t.Errorf("palette %q should be flagged colorblind-safe", id)
		}
	}
}

// TestOrdinalScaleUsesCuratedPalette is an end-to-end check that a curated
// palette drives an ordinal color scale (the path charts use).
func TestOrdinalScaleUsesCuratedPalette(t *testing.T) {
	scale := GetOrdinalColorScale[string](Scheme(PaletteOkabeIto), func(s string) string { return s })
	first := scale("a")
	want := curatedCategoricalSchemes["okabe_ito"][0]
	if first != want {
		t.Errorf("ordinal scale first color = %q, want %q", first, want)
	}
	if scale("a") != first {
		t.Errorf("ordinal scale not stable for same key")
	}
	if scale("b") == first {
		t.Errorf("ordinal scale should assign a different color to a new key")
	}
}
