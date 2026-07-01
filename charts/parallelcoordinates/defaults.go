package parallelcoordinates

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/parallel-coordinates defaultProps. Fields left zero in
// a PCProps fall back to these via applyDefaults.
var Defaults = PCProps{
	Layout:            PCLayoutHorizontal,
	Curve:             core.CurveLinear,
	LineWidth:         2,
	LineOpacity:       0.5,
	AxesTicksPosition: "after",
	Colors:            colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "category10"},
	Layers:            DefaultLayers,
	Role:              "img",
	MotionProps:       core.MotionProps{Animate: true, MotionConfig: "default"},
}
