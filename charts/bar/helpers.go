package bar

import (
	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/rects"
)

// roundedRectPath builds the SVG path-data for a rounded rect at (0,0,w,h)
// with a uniform corner radius. Delegates to rects.BuildRoundedRectPath.
func roundedRectPath(w, h, radius float64) string {
	return rects.BuildRoundedRectPath(0, 0, w, h, radius, radius, radius, radius)
}

// maxF returns the larger of a or b.
func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// bindBarAnnotations mirrors @nivo/bar BarAnnotations' useAnnotations: binds
// annotation specs to bars (position = bar center; dimensions = bar w/h) and
// computes the rendering instructions for each bound annotation.
func bindBarAnnotations(bars []ComputedBarDatum, specs []annotations.AnnotationSpec[ComputedBarDatum]) []annotations.AnnotationInstructions {
	bound := annotations.BindAnnotations(
		bars, specs,
		func(b ComputedBarDatum) (float64, float64) {
			return b.X + b.Width/2, b.Y + b.Height/2
		},
		func(b ComputedBarDatum) (float64, float64) {
			return b.Width, b.Height
		},
	)
	out := make([]annotations.AnnotationInstructions, len(bound))
	for i, b := range bound {
		out[i] = annotations.ComputeAnnotation(b)
	}
	return out
}

// legendItems converts a slice of bar.LegendData to legends.Datum (the type
// legends.BoxLegendSvg expects).
func legendItems(data []LegendData) []legends.Datum {
	out := make([]legends.Datum, len(data))
	for i, d := range data {
		out[i] = legends.Datum{
			ID:     d.ID,
			Label:  d.Label,
			Color:  d.Color,
			Hidden: d.Hidden,
		}
	}
	return out
}
