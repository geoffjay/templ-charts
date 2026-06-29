package colors

// Curated palettes extend the d3/nivo scheme set with a broader, more visually
// appealing selection — including several colorblind-safe palettes. They are
// registered into the existing scheme/interpolator maps at init() so the rest
// of the color machinery (resolveSchemeColors, ColorInterpolators, the chart
// Colors configs) picks them up with no further wiring. The generated
// schemes_data.go is left untouched.

// curatedCategoricalSchemes maps id → flat color array for the curated
// categorical palettes. These are merged into CategoricalColorSchemes.
var curatedCategoricalSchemes = map[string][]string{
	// Observable Plot's default categorical scheme — bright, modern, balanced.
	"observable10": {"#4269d0", "#efb118", "#ff725c", "#6cc5b0", "#3ca951", "#ff8ab7", "#a463f2", "#97bbf5", "#9c6b4e", "#9498a0"},
	// Tableau's 20-color palette — pairs of light/dark hues for large series counts.
	"tableau20": {
		"#4e79a7", "#a0cbe8", "#f28e2b", "#ffbe7d", "#59a14f", "#8cd17d", "#b6992d", "#f1ce63",
		"#499894", "#86bcb6", "#e15759", "#ff9d9a", "#79706e", "#bab0ab", "#d37295", "#fabfd2",
		"#b07aa1", "#d4a6c8", "#9d7660", "#d7b5a6",
	},
	// Material Design 500-level hues — saturated and cheerful.
	"material": {
		"#f44336", "#e91e63", "#9c27b0", "#673ab7", "#3f51b5", "#2196f3", "#03a9f4", "#00bcd4",
		"#009688", "#4caf50", "#8bc34a", "#cddc39", "#ffeb3b", "#ffc107", "#ff9800", "#ff5722",
		"#795548", "#607d8b",
	},
	// Okabe–Ito colorblind-safe palette (grey substituted for black so the
	// first series isn't pure black). Widely recommended for accessibility.
	"okabe_ito": {"#e69f00", "#56b4e9", "#009e73", "#f0e442", "#0072b2", "#d55e00", "#cc79a7", "#999999"},
	// Paul Tol "vibrant" — colorblind-safe, high-contrast.
	"tol_vibrant": {"#0077bb", "#33bbee", "#009988", "#ee7733", "#cc3311", "#ee3377", "#bbbbbb"},
	// Paul Tol "muted" — colorblind-safe, softer tones for dense charts.
	"tol_muted": {"#332288", "#88ccee", "#44aa99", "#117733", "#999933", "#ddcc77", "#cc6677", "#882255", "#aa4499", "#dddddd"},
}

// curatedCategoricalIds preserves a stable display order for the curated
// categorical palettes.
var curatedCategoricalIds = []string{
	"observable10",
	"tableau20",
	"material",
	"okabe_ito",
	"tol_vibrant",
	"tol_muted",
}

// curatedSequentialStops maps id → ordered color stops for curated sequential
// gradients. Registered as continuous interpolators (interpolateRgbBasis).
var curatedSequentialStops = map[string][]string{
	// CARTO "Sunset" — warm low-to-high gradient.
	"sunset": {"#f3e79b", "#fac484", "#f8a07e", "#eb7f86", "#ce6693", "#a059a0", "#5c53a5"},
	// CARTO "Teal" — calm cyan-to-deep-teal.
	"teal": {"#d1eeea", "#a8dbd9", "#85c4c9", "#68abb8", "#4f90a6", "#3b738f", "#2a5674"},
	// CARTO "Emrld" — fresh green ramp.
	"emerald": {"#d3f2a3", "#97e196", "#6cc08b", "#4c9b82", "#217a79", "#105965", "#074050"},
}

var curatedSequentialIds = []string{"sunset", "teal", "emerald"}

// curatedDivergingStops maps id → ordered color stops for curated diverging
// gradients. Registered as continuous interpolators (interpolateRgbBasis).
var curatedDivergingStops = map[string][]string{
	// CARTO "Temps" — teal-to-red diverging, good for temperature-style data.
	// ("spectral" is already provided by the generated d3 diverging schemes.)
	"temps": {"#009392", "#39b185", "#9ccb86", "#e9e29c", "#eeb479", "#e88471", "#cf597e"},
}

var curatedDivergingIds = []string{"temps"}

func init() {
	// Merge curated categorical palettes into the categorical scheme registry
	// and id list so IsCategoricalColorScheme / resolveSchemeColors resolve them.
	for id, cols := range curatedCategoricalSchemes {
		CategoricalColorSchemes[id] = cols
	}
	categoricalColorSchemeIds = append(categoricalColorSchemeIds, curatedCategoricalIds...)
	// ColorSchemeIds is built by a package-var initializer that ran before this
	// init(), so append the curated categorical ids to keep enumeration complete.
	ColorSchemeIds = append(ColorSchemeIds, curatedCategoricalIds...)

	// Register curated continuous gradients as interpolators (the form
	// GetSequentialColorScale / GetDivergingColorScale consume). They are not
	// added to the ordinal scheme maps since gradients aren't categorical.
	for _, id := range curatedSequentialIds {
		if _, exists := ColorInterpolators[id]; !exists {
			ColorInterpolators[id] = interpolateRgbBasis(curatedSequentialStops[id])
			ColorInterpolatorIds = append(ColorInterpolatorIds, id)
		}
	}
	for _, id := range curatedDivergingIds {
		if _, exists := ColorInterpolators[id]; !exists {
			ColorInterpolators[id] = interpolateRgbBasis(curatedDivergingStops[id])
			ColorInterpolatorIds = append(ColorInterpolatorIds, id)
		}
	}
}
