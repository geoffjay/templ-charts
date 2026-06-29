package demos

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/htmx"
)

// PaletteDemo describes one palette tile on the /palettes page: the palette
// metadata, a swatch preview, and a bar chart rendered with that palette.
type PaletteDemo struct {
	Palette colors.Palette
	Swatch  []string
	Kind    htmx.ChartKind
	Props   any
}

// paletteKeys is the set of bar keys used by the palette gallery — enough
// distinct series to show several palette colors at once.
var paletteKeys = []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"}

// PaletteDemos returns one tile per catalog palette for the /palettes page.
// Categorical palettes drive the bar colors via their ordinal scheme;
// sequential/diverging palettes sample their gradient into the same number of
// discrete colors so they're demonstrable on a categorical chart.
func PaletteDemos() []PaletteDemo {
	palettes := colors.Palettes()
	out := make([]PaletteDemo, 0, len(palettes))
	for _, p := range palettes {
		var colorsCfg colors.OrdinalColorScaleConfig
		var swatch []string
		switch p.Kind {
		case colors.KindCategorical:
			colorsCfg = p.Ordinal()
			swatch = p.Swatch(swatchCount(p))
		default:
			// Gradients aren't categorical: sample them into discrete colors
			// for both the chart and the swatch.
			colorsCfg = colors.PaletteColors(p.Swatch(len(paletteKeys))...)
			swatch = p.Swatch(12)
		}
		out = append(out, PaletteDemo{
			Palette: p,
			Swatch:  swatch,
			Kind:    htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				IndexBy: "country",
				Keys:    paletteKeys,
				Data:    paletteData(),
				Margin:  defaultMargin(),
				Colors:  colorsCfg,
			},
		})
	}
	return out
}

// swatchCount returns how many swatches to show for a categorical palette: its
// full color set (so large palettes like tableau20/material show their range).
func swatchCount(p colors.Palette) int {
	if n := len(colors.CategoricalColorSchemes[string(p.ID)]); n > 0 {
		return n
	}
	return len(paletteKeys)
}

// paletteData is a small deterministic dataset (5 indices × paletteKeys) for
// the palette gallery bars.
func paletteData() []bar.BarDatum {
	indices := []string{"A", "B", "C", "D", "E"}
	data := make([]bar.BarDatum, len(indices))
	for i, idx := range indices {
		row := map[string]any{"country": idx}
		for j, k := range paletteKeys {
			row[k] = float64((i*5+j*11)%60 + 15)
		}
		data[i] = row
	}
	return data
}
