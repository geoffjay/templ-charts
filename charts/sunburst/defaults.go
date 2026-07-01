package sunburst

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
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
	MotionProps:           core.MotionProps{Animate: true, MotionConfig: "default"},
}
