package sunburst

import (
	"github.com/geoffjay/templ-charts/charts/colors"
)

// Defaults mirrors @nivo/sunburst svgDefaultProps.
var Defaults = SunburstProps{
	CornerRadius:          0,
	Colors:                colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderWidth:           1,
	BorderColor:           "white",
	EnableArcLabels:       BoolPtr(false),
	ArcLabelsRadiusOffset: 0.5,
	Role:                  "img",
}
