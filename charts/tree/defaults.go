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
	NodeColor:     colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	LinkThickness: 1,
	LinkOpacity:   0.4,
	EnableLabel:   BoolPtr(true),
	LabelOffset:   6,
	UseMesh:       true,
	Role:          "img",
	MotionProps:   core.MotionProps{Animate: true, MotionConfig: "default"},
}
