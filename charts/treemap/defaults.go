package treemap

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/treemap svgDefaultProps. Fields left zero in a
// TreemapProps fall back to these via applyDefaults.
var Defaults = TreemapProps{
	Tile:               TileSquarify,
	InnerPadding:       0,
	OuterPadding:       0,
	Colors:             colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	NodeOpacity:        0.33,
	BorderWidth:        1,
	EnableLabel:        core.BoolPtr(true),
	LabelSkipSize:      0,
	EnableParentLabel:  core.BoolPtr(true),
	ParentLabelSize:    20,
	ParentLabelPadding: 6,
	Role:               "img",
}
