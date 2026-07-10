package line

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/tooltip"
)

// Defaults mirrors @nivo/line commonDefaultProps + svgDefaultProps. Fields
// with zero values in a LineProps fall back to these via UseLine. AxisBottom/
// AxisLeft are left nil here; the render path substitutes axes.DefaultAxisProps
// when nil (matching nivo's `axisBottom: defaultAxisProps`).
var Defaults = LineProps{
	XScale:               scales.ScalePointSpec{},
	YScale:               scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
	Curve:                core.CurveLinear,
	Colors:               colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	LineWidth:            2,
	Layers:               DefaultLayers,
	EnablePoints:         core.BoolPtr(true),
	PointSize:            6,
	PointColor:           colors.NewFromContextColor("series.color", nil),
	PointBorderWidth:     0,
	PointBorderColor:     colors.NewThemeColor("background"),
	EnableArea:           core.BoolPtr(false),
	AreaBaselineValue:    0,
	AreaOpacity:          0.2,
	AreaBlendMode:        core.MixBlendNormal,
	EnableGridX:          core.BoolPtr(true),
	EnableGridY:          core.BoolPtr(true),
	Legends:              []legends.LegendProps{},
	Interactive:          true,
	DebugMesh:            false,
	EnablePointLabel:     core.BoolPtr(false),
	PointLabel:           "data.yFormatted",
	AxisTop:              nil,
	AxisRight:            nil,
	AxisBottom:           nil,
	AxisLeft:             nil,
	UseMesh:              false,
	EnableSlices:         EnableSlicesFalse,
	DebugSlices:          false,
	EnableCrosshair:      core.BoolPtr(true),
	CrosshairType:        tooltip.CrosshairTypeBottomLeft,
	EnableTouchCrosshair: core.BoolPtr(false),
	InitialHiddenIDs:     []string{},
	MotionProps:          core.MotionProps{Animate: true, MotionConfig: core.DefaultMotionConfig},
	Role:                 "img",
	IsFocusable:          false,
}

// defaultAxisBottomProps / defaultAxisLeftProps are the concrete default axis
// configs used when the caller leaves AxisBottom/AxisLeft nil. They mirror
// @nivo/axes defaultAxisProps (tickSize=5, tickPadding=5, textAlign=center,
// textBaseline=middle, legendPosition=end).
var defaultAxisBottomProps = axes.AxisProps{
	TickSize:       5,
	TickPadding:    5,
	TickRotation:   0,
	TextAlign:      "center",
	TextBaseline:   "middle",
	LegendPosition: axes.AxisLegendEnd,
	LegendOffset:   0,
}

var defaultAxisLeftProps = defaultAxisBottomProps
