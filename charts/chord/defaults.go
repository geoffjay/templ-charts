package chord

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/chord svgDefaultProps. Fields left zero in a
// ChordProps fall back to these via applyDefaults.
var Defaults = ChordProps{
	PadAngle:          0,
	InnerRadiusRatio:  0.9,
	InnerRadiusOffset: 0,

	Colors: colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},

	ArcOpacity:         1,
	ActiveArcOpacity:   1,
	InactiveArcOpacity: 0.15,
	ArcBorderWidth:     1,
	ArcBorderColor:     colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 0.4}}),

	RibbonOpacity:         0.5,
	ActiveRibbonOpacity:   0.85,
	InactiveRibbonOpacity: 0.15,
	RibbonBorderWidth:     1,
	RibbonBorderColor:     colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 0.4}}),
	RibbonBlendMode:       "normal",

	EnableLabel:    boolPtr(true),
	Label:          "id",
	LabelOffset:    12,
	LabelRotation:  0,
	LabelTextColor: colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 1}}),

	Layers:      DefaultLayers,
	Role:        "img",
	MotionProps: core.MotionProps{Animate: true, MotionConfig: "gentle"},
}

func boolPtr(b bool) *bool { return &b }
