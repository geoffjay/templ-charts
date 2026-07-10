package polarbar

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/polar-bar defaultProps. Fields left zero in a
// PolarBarProps fall back to these via applyDefaults.
var Defaults = PolarBarProps{
	Keys:                  []string{"value"},
	IndexBy:               "id",
	StartAngle:            0,
	EndAngle:              360,
	InnerRadius:           0,
	CornerRadius:          0,
	EnableRadialGrid:      core.BoolPtr(true),
	EnableCircularGrid:    core.BoolPtr(true),
	Colors:                colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	EnableArcLabels:       core.BoolPtr(false),
	ArcLabel:              "formattedValue",
	ArcLabelsRadiusOffset: 0.5,
	Layers:                DefaultLayers,
	Role:                  "img",
}
