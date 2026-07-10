package tree

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/tree svgDefaultProps.
var Defaults = TreeProps{
	Mode:          ModeDendogram,
	Layout:        LayoutTopToBottom,
	NodeSize:      12,
	Colors:        colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	LinkThickness: 1,
	LinkOpacity:   0.4,
	EnableLabel:   core.BoolPtr(true),
	LabelOffset:   6,
	UseMesh:       true,
	Role:          "img",
}
