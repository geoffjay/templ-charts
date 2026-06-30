package waffle

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
)

// Defaults mirrors @nivo/waffle commonDefaultProps + svgDefaultProps. Fields
// left zero in a WaffleProps fall back to these via applyDefaults.
var Defaults = WaffleProps{
	FillDirection: grid.GridFillTop,
	Padding:       1,
	Colors:        colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	EmptyColor:    "#cccccc",
	EmptyOpacity:  1,
	BorderRadius:  0,
	BorderWidth:   0,
	BorderColor:   colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 1.0}}),
	Layers:        DefaultLayers,
	Role:          "img",
	MotionProps:   core.MotionProps{Animate: true, MotionConfig: core.DefaultMotionConfig},
}
