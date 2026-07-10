package boxplot

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/boxplot svgDefaultProps. Fields left zero in a
// BoxPlotProps fall back to these via applyDefaults.
var Defaults = BoxPlotProps{
	Quantiles:       []float64{0.1, 0.25, 0.5, 0.75, 0.9},
	Layout:          BoxPlotLayoutVertical,
	Padding:         0.1,
	InnerPadding:    6,
	Opacity:         1,
	ActiveOpacity:   1,
	InactiveOpacity: 0.25,
	EnableGridX:     core.BoolPtr(false),
	EnableGridY:     core.BoolPtr(true),
	ColorBy:         "subGroup",
	Colors:          colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderRadius:    0,
	BorderWidth:     0,
	BorderColor:     colors.NewFromContextColor("color", nil),
	MedianWidth:     2,
	MedianColor:     colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 2.0}}),
	WhiskerWidth:    2,
	WhiskerColor:    colors.NewFromContextColor("color", nil),
	WhiskerEndSize:  0.6,
	Layers:          DefaultLayers,
	Role:            "img",
}
