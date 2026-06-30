package scatterplot

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// Defaults mirrors @nivo/scatterplot svgDefaultProps. Fields left zero in a
// ScatterPlotProps fall back to these via applyDefaults.
var Defaults = ScatterPlotProps{
	XScale:      scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
	YScale:      scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
	EnableGridX: true,
	EnableGridY: true,
	NodeSize:    9,
	Colors:      colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	Layers:      DefaultLayers,
	Role:        "img",
	MotionProps: core.MotionProps{Animate: true, MotionConfig: "default"},
}
