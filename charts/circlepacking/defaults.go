package circlepacking

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/circle-packing svgDefaultProps.
var Defaults = CirclePackingProps{
	Padding:          0,
	Colors:           colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	BorderWidth:      0,
	EnableLabels:     BoolPtr(false),
	LabelsSkipRadius: 8,
	Role:             "img",
	MotionProps:      core.MotionProps{Animate: true, MotionConfig: "default"},
}
