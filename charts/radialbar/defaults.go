package radialbar

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/radial-bar commonDefaultProps + svgDefaultProps.
// Fields left zero in a RadialBarProps fall back to these via applyDefaults.
var Defaults = RadialBarProps{
	Layers:                DefaultLayers,
	StartAngle:            0,
	EndAngle:              270,
	InnerRadius:           0.3,
	Padding:               0.2,
	PadAngle:              0,
	CornerRadius:          0,
	EnableTracks:          core.BoolPtr(true),
	TracksColor:           "rgba(0, 0, 0, .15)",
	EnableRadialGrid:      core.BoolPtr(true),
	EnableCircularGrid:    core.BoolPtr(true),
	ShowRadialAxisStart:   core.BoolPtr(true),
	ShowRadialAxisEnd:     core.BoolPtr(false),
	ShowCircularAxisInner: core.BoolPtr(false),
	ShowCircularAxisOuter: core.BoolPtr(true),
	Colors:                colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderWidth:           0,
	BorderColor:           colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 1.0}}),
	EnableLabels:          core.BoolPtr(false),
	Label:                 "formattedValue",
	LabelsSkipAngle:       10,
	LabelsRadiusOffset:    0.5,
	LabelsTextColor:       colors.NewThemeColor("labels.text.fill"),
	Role:                  "img",
}
