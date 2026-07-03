package marimekko

import (
	"github.com/geoffjay/templ-charts/charts/colors"
)

// Defaults mirrors @nivo/marimekko defaultProps. Fields left zero in a
// MarimekkoProps fall back to these via applyDefaults.
var Defaults = MarimekkoProps{
	Layout:       MarimekkoLayoutVertical,
	Offset:       OffsetNone,
	OuterPadding: 0,
	InnerPadding: 3,
	Colors:       colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	EnableGridX:  false,
	EnableGridY:  true,
	Layers:       DefaultLayers,
	Role:         "img",
}
