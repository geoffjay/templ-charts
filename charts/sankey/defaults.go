package sankey

import (
	"github.com/geoffjay/templ-charts/charts/colors"
)

// Defaults mirrors @nivo/sankey svgDefaultProps. Fields left zero in a
// SankeyProps fall back to these via applyDefaults.
var Defaults = SankeyProps{
	Layout: SankeyLayoutHorizontal,
	Align:  SankeyAlignCenter,
	Sort:   SankeySortAuto,

	Colors: colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},

	NodeOpacity:      0.75,
	NodeThickness:    12,
	NodeSpacing:      12,
	NodeInnerPadding: 0,
	NodeBorderWidth:  1,
	NodeBorderColor:  colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 0.5}}),
	NodeBorderRadius: 0,

	LinkOpacity:   0.25,
	LinkContract:  0,
	LinkBlendMode: "multiply",

	NodeHoverOpacity:       1,
	NodeHoverOthersOpacity: 0.35,
	LinkHoverOpacity:       0.6,
	LinkHoverOthersOpacity: 0.1,

	EnableLabels:     boolPtr(true),
	Label:            "id",
	LabelPosition:    SankeyLabelInside,
	LabelPadding:     9,
	LabelOrientation: SankeyLabelHorizontal,
	LabelTextColor:   colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", 0.8}}),

	Layers: DefaultLayers,
	Role:   "img",
}

func boolPtr(b bool) *bool { return &b }
