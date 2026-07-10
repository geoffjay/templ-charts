package heatmap

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/heatmap commonDefaultProps + svgDefaultProps. Fields
// left zero in a HeatMapProps fall back to these via applyDefaults.
var Defaults = HeatMapProps{
	ForceSquare:     false,
	Colors:          HeatMapColorConfig{Type: "sequential", Scheme: "brown_blueGreen"},
	EmptyColor:      "#000000",
	EnableGridX:     core.BoolPtr(false),
	EnableGridY:     core.BoolPtr(false),
	Opacity:         1,
	ActiveOpacity:   1,
	InactiveOpacity: 0.15,
	BorderWidth:     0,
	BorderColor:     colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 0.8}}),
	BorderRadius:    0,
	EnableLabels:    core.BoolPtr(true),
	LabelTextColor:  colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 2.0}}),
	Layers:          DefaultLayers,
	Role:            "img",
}
