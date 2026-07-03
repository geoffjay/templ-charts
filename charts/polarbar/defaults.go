package polarbar

import (
	"github.com/geoffjay/templ-charts/charts/colors"
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
	EnableRadialGrid:      BoolPtr(true),
	EnableCircularGrid:    BoolPtr(true),
	Colors:                colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	EnableArcLabels:       BoolPtr(false),
	ArcLabel:              "formattedValue",
	ArcLabelsRadiusOffset: 0.5,
	Layers:                DefaultLayers,
	Role:                  "img",
}
