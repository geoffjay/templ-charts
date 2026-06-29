package colors

import "sort"

// PaletteID is a typed color-palette identifier. Constants below give callers
// a discoverable, autocomplete-friendly, compile-checked set of palette names
// (vs. raw scheme strings). It is a string under the hood, so it remains
// interchangeable with the existing Scheme string fields.
type PaletteID string

// Categorical palette ids — discrete colors cycled across series/slices.
const (
	// nivo + d3 categorical schemes (from schemes_data.go).
	PaletteNivo       PaletteID = "nivo"
	PaletteCategory10 PaletteID = "category10"
	PaletteAccent     PaletteID = "accent"
	PaletteDark2      PaletteID = "dark2"
	PalettePaired     PaletteID = "paired"
	PalettePastel1    PaletteID = "pastel1"
	PalettePastel2    PaletteID = "pastel2"
	PaletteSet1       PaletteID = "set1"
	PaletteSet2       PaletteID = "set2"
	PaletteSet3       PaletteID = "set3"
	PaletteTableau10  PaletteID = "tableau10"

	// Curated categorical palettes (from palettes_data.go).
	PaletteObservable10 PaletteID = "observable10"
	PaletteTableau20    PaletteID = "tableau20"
	PaletteMaterial     PaletteID = "material"
	PaletteOkabeIto     PaletteID = "okabe_ito"
	PaletteTolVibrant   PaletteID = "tol_vibrant"
	PaletteTolMuted     PaletteID = "tol_muted"
)

// Sequential palette ids — single-direction value gradients.
const (
	PaletteViridis PaletteID = "viridis"
	PaletteMagma   PaletteID = "magma"
	PaletteInferno PaletteID = "inferno"
	PalettePlasma  PaletteID = "plasma"
	PaletteTurbo   PaletteID = "turbo"
	PaletteCividis PaletteID = "cividis"
	PaletteWarm    PaletteID = "warm"
	PaletteCool    PaletteID = "cool"

	// Curated sequential gradients (from palettes_data.go).
	PaletteSunset  PaletteID = "sunset"
	PaletteTeal    PaletteID = "teal"
	PaletteEmerald PaletteID = "emerald"
)

// Diverging palette ids — two-direction gradients around a midpoint.
const (
	PaletteRedBlue        PaletteID = "red_blue"
	PaletteRedYellowBlue  PaletteID = "red_yellow_blue"
	PaletteBrownBlueGreen PaletteID = "brown_blueGreen"

	// Curated diverging gradients (from palettes_data.go).
	PaletteSpectral PaletteID = "spectral"
	PaletteTemps    PaletteID = "temps"
)

// PaletteKind classifies how a palette assigns color.
type PaletteKind int

const (
	// KindCategorical: discrete colors cycled across series/slices.
	KindCategorical PaletteKind = iota
	// KindSequential: a single-direction value gradient.
	KindSequential
	// KindDiverging: a gradient that diverges around a midpoint.
	KindDiverging
)

// String returns the lowercase kind label (categorical/sequential/diverging).
func (k PaletteKind) String() string {
	switch k {
	case KindSequential:
		return "sequential"
	case KindDiverging:
		return "diverging"
	default:
		return "categorical"
	}
}

// Palette is a named, discoverable color palette with metadata for galleries,
// docs, and ergonomic chart configuration. It wraps the underlying scheme /
// interpolator id; use Ordinal/Sequential/Diverging to get a chart-ready
// config, or Swatch to preview it.
type Palette struct {
	ID             PaletteID
	Name           string // human-friendly display name
	Kind           PaletteKind
	Group          string // provenance, e.g. "nivo", "d3", "curated"
	ColorblindSafe bool
}

// Ordinal returns an OrdinalColorScaleConfig for categorical use (the form
// bar/line/pie Colors fields accept). For sequential/diverging palettes the
// gradient is sampled into discrete steps via the ordinal scheme machinery.
func (p Palette) Ordinal() OrdinalColorScaleConfig {
	return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: string(p.ID)}
}

// Sequential returns a SequentialColorScaleConfig bound to this palette's
// gradient. Meaningful for KindSequential palettes (and usable for any id
// registered as an interpolator).
func (p Palette) Sequential() SequentialColorScaleConfig {
	return SequentialColorScaleConfig{Type: "sequential", Scheme: string(p.ID)}
}

// Diverging returns a DivergingColorScaleConfig bound to this palette's
// gradient. Meaningful for KindDiverging palettes.
func (p Palette) Diverging() DivergingColorScaleConfig {
	return DivergingColorScaleConfig{Type: "diverging", Scheme: string(p.ID)}
}

// Swatch returns up to n representative colors for previewing the palette.
// Categorical palettes return their colors (cycled if n exceeds the set);
// sequential/diverging palettes sample their gradient at n evenly-spaced
// stops. n <= 0 yields an empty slice.
func (p Palette) Swatch(n int) []string {
	if n <= 0 {
		return nil
	}
	switch p.Kind {
	case KindCategorical:
		base := CategoricalColorSchemes[string(p.ID)]
		if len(base) == 0 {
			return nil
		}
		out := make([]string, n)
		for i := range out {
			out[i] = base[i%len(base)]
		}
		return out
	default:
		interp := ColorInterpolators[string(p.ID)]
		if interp == nil {
			return nil
		}
		out := make([]string, n)
		if n == 1 {
			out[0] = interp(0.5)
			return out
		}
		for i := range out {
			out[i] = interp(float64(i) / float64(n-1))
		}
		return out
	}
}

// Scheme is the ergonomic constructor for an ordinal (categorical) color
// config from a palette id: colors.Scheme(colors.PaletteTableau10). It is a
// concise alternative to spelling out the OrdinalColorScaleConfig struct.
func Scheme(id PaletteID) OrdinalColorScaleConfig {
	return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: string(id)}
}

// PaletteColors builds an ordinal color config from an explicit color list,
// for one-off custom palettes: colors.PaletteColors("#f00", "#0f0", "#00f").
func PaletteColors(cols ...string) OrdinalColorScaleConfig {
	return OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: cols}
}

// paletteRegistry is the canonical, ordered catalog of named palettes. Built
// at init() from the curated metadata below plus the registered scheme ids.
var paletteRegistry []Palette

// paletteByID indexes paletteRegistry by id for O(1) lookup.
var paletteByID map[PaletteID]Palette

// paletteMeta carries the display metadata for catalog entries. Order here
// defines catalog order within each kind.
type paletteMeta struct {
	id             PaletteID
	name           string
	kind           PaletteKind
	group          string
	colorblindSafe bool
}

var paletteCatalog = []paletteMeta{
	// Categorical — nivo/d3.
	{PaletteNivo, "Nivo", KindCategorical, "nivo", false},
	{PaletteCategory10, "Category 10", KindCategorical, "d3", false},
	{PaletteAccent, "Accent", KindCategorical, "d3", false},
	{PaletteDark2, "Dark 2", KindCategorical, "d3", false},
	{PalettePaired, "Paired", KindCategorical, "d3", false},
	{PalettePastel1, "Pastel 1", KindCategorical, "d3", false},
	{PalettePastel2, "Pastel 2", KindCategorical, "d3", false},
	{PaletteSet1, "Set 1", KindCategorical, "d3", false},
	{PaletteSet2, "Set 2", KindCategorical, "d3", false},
	{PaletteSet3, "Set 3", KindCategorical, "d3", false},
	{PaletteTableau10, "Tableau 10", KindCategorical, "d3", false},
	// Categorical — curated.
	{PaletteObservable10, "Observable 10", KindCategorical, "curated", false},
	{PaletteTableau20, "Tableau 20", KindCategorical, "curated", false},
	{PaletteMaterial, "Material", KindCategorical, "curated", false},
	{PaletteOkabeIto, "Okabe–Ito", KindCategorical, "curated", true},
	{PaletteTolVibrant, "Tol Vibrant", KindCategorical, "curated", true},
	{PaletteTolMuted, "Tol Muted", KindCategorical, "curated", true},
	// Sequential.
	{PaletteViridis, "Viridis", KindSequential, "d3", true},
	{PaletteMagma, "Magma", KindSequential, "d3", true},
	{PaletteInferno, "Inferno", KindSequential, "d3", true},
	{PalettePlasma, "Plasma", KindSequential, "d3", true},
	{PaletteTurbo, "Turbo", KindSequential, "d3", false},
	{PaletteCividis, "Cividis", KindSequential, "d3", true},
	{PaletteWarm, "Warm", KindSequential, "d3", false},
	{PaletteCool, "Cool", KindSequential, "d3", false},
	{PaletteSunset, "Sunset", KindSequential, "curated", false},
	{PaletteTeal, "Teal", KindSequential, "curated", false},
	{PaletteEmerald, "Emerald", KindSequential, "curated", false},
	// Diverging.
	{PaletteRedBlue, "Red–Blue", KindDiverging, "d3", false},
	{PaletteRedYellowBlue, "Red–Yellow–Blue", KindDiverging, "d3", false},
	{PaletteBrownBlueGreen, "Brown–BlueGreen", KindDiverging, "d3", false},
	{PaletteSpectral, "Spectral", KindDiverging, "d3", false},
	{PaletteTemps, "Temps", KindDiverging, "curated", false},
}

func init() {
	paletteRegistry = make([]Palette, 0, len(paletteCatalog))
	paletteByID = make(map[PaletteID]Palette, len(paletteCatalog))
	for _, m := range paletteCatalog {
		p := Palette{ID: m.id, Name: m.name, Kind: m.kind, Group: m.group, ColorblindSafe: m.colorblindSafe}
		paletteRegistry = append(paletteRegistry, p)
		paletteByID[m.id] = p
	}
}

// Palettes returns the full ordered catalog of named palettes.
func Palettes() []Palette {
	out := make([]Palette, len(paletteRegistry))
	copy(out, paletteRegistry)
	return out
}

// PalettesByKind returns the catalog filtered to one kind, in catalog order.
func PalettesByKind(kind PaletteKind) []Palette {
	var out []Palette
	for _, p := range paletteRegistry {
		if p.Kind == kind {
			out = append(out, p)
		}
	}
	return out
}

// LookupPalette returns the catalog entry for id and whether it was found.
func LookupPalette(id PaletteID) (Palette, bool) {
	p, ok := paletteByID[id]
	return p, ok
}

// ColorblindSafePalettes returns the catalog entries flagged colorblind-safe,
// in catalog order.
func ColorblindSafePalettes() []Palette {
	var out []Palette
	for _, p := range paletteRegistry {
		if p.ColorblindSafe {
			out = append(out, p)
		}
	}
	return out
}

// PaletteIDs returns the sorted ids of all catalog palettes (handy for tests
// and validation).
func PaletteIDs() []string {
	out := make([]string, 0, len(paletteRegistry))
	for _, p := range paletteRegistry {
		out = append(out, string(p.ID))
	}
	sort.Strings(out)
	return out
}
