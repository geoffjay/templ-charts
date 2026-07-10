package swarmplot

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// Defaults mirrors @nivo/swarmplot defaultProps. Fields left zero in a
// SwarmPlotProps fall back to these via applyDefaults.
var Defaults = SwarmPlotProps{
	ValueScale:           scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
	Size:                 6,
	Spacing:              2,
	Layout:               "vertical",
	Gap:                  0,
	ForceStrength:        1,
	SimulationIterations: 120,
	Colors:               colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	ColorBy:              "group",
	BorderWidth:          0,
	BorderColor:          "rgba(0, 0, 0, 0)",
	EnableGridX:          core.BoolPtr(true),
	EnableGridY:          core.BoolPtr(true),
	Layers:               DefaultLayers,
	Role:                 "img",
}
