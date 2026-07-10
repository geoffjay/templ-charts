package bar

import (
	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// Defaults mirrors @nivo/bar commonDefaultProps + svgDefaultProps. Fields
// with zero values in a BarProps fall back to these via UseBar.
var Defaults = BarProps{
	IndexBy:      "id",
	Keys:         []string{"value"},
	GroupMode:    GroupModeStacked,
	Layout:       LayoutVertical,
	Padding:      0.1,
	InnerPadding: 0,
	ValueScale: scales.ScaleLinearSpec{
		Min: scales.FloatVal(0), Max: scales.AutoFloat(),
		Nice: true, Round: false,
	},
	IndexScale:       scales.ScaleBandSpec{Round: false},
	EnableGridX:      core.BoolPtr(false),
	EnableGridY:      core.BoolPtr(true),
	EnableLabel:      core.BoolPtr(true),
	Label:            "formattedValue",
	LabelPosition:    LabelPositionMiddle,
	LabelOffset:      0,
	LabelSkipWidth:   0,
	LabelSkipHeight:  0,
	LabelTextColor:   colors.NewThemeColor("labels.text.fill"),
	ColorBy:          ColorByID,
	Colors:           colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderRadius:     0,
	BorderWidth:      0,
	BorderColor:      colors.NewFromContextColor("color", nil),
	Interactive:      true,
	Legends:          []BarLegendProps{},
	InitialHiddenIDs: []string{},
	Annotations:      []annotations.AnnotationSpec[ComputedBarDatum]{},
	EnableTotals:     core.BoolPtr(false),
	TotalsOffset:     10,
	Layers:           DefaultLayers,
	MotionProps:      core.MotionProps{Animate: true, MotionConfig: core.DefaultMotionConfig},
	Role:             "img",
	IsFocusable:      false,
}
