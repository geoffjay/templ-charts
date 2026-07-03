package stream

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/stream svgDefaultProps. Fields left zero in a
// StreamProps fall back to these via applyDefaults.
var Defaults = StreamProps{
	OffsetType:  core.StackOffsetWiggle,
	Order:       core.StackOrderNone,
	Curve:       core.CurveCatmullRom,
	EnableGridX: false,
	EnableGridY: true,
	Colors:      colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	FillOpacity: 1,
	BorderWidth: 0,
	BorderColor: colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 1.0}}),
	Layers:      DefaultLayers,
	Role:        "img",
}
