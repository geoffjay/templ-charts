package pie

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/pie defaultProps. Fields with zero values in a
// PieProps fall back to these via applyDefaults.
var Defaults = PieProps{
	ID:                          "id",
	Value:                       "value",
	SortByValue:                 false,
	InnerRadius:                 0,
	PadAngle:                    0,
	CornerRadius:                0,
	Layers:                      DefaultLayers,
	StartAngle:                  0,
	EndAngle:                    360,
	Fit:                         true,
	ActiveInnerRadiusOffset:     0,
	ActiveOuterRadiusOffset:     0,
	BorderWidth:                 0,
	BorderColor:                 colors.NewFromContextColor("color", []colors.ColorModifier{{"darker", float64(1)}}),
	EnableArcLabels:             true,
	ArcLabel:                    "formattedValue",
	ArcLabelsSkipAngle:          0,
	ArcLabelsSkipRadius:         0,
	ArcLabelsRadiusOffset:       0.5,
	ArcLabelsTextColor:          colors.NewThemeColor("labels.text.fill"),
	EnableArcLinkLabels:         true,
	ArcLinkLabel:                "id",
	ArcLinkLabelsSkipAngle:      0,
	ArcLinkLabelsOffset:         0,
	ArcLinkLabelsDiagonalLength: 16,
	ArcLinkLabelsStraightLength: 24,
	ArcLinkLabelsThickness:      1,
	ArcLinkLabelsTextOffset:     6,
	ArcLinkLabelsTextColor:      colors.NewThemeColor("labels.text.fill"),
	ArcLinkLabelsColor:          colors.NewThemeColor("axis.ticks.line.stroke"),
	Colors:                      colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
	Defs:                        []core.Def{},
	Fill:                        []core.DefRule{},
	IsInteractive:               true,
	MotionProps:                 core.MotionProps{Animate: true, MotionConfig: core.DefaultMotionConfig},
	TransitionMode:              TransitionModeInnerRadius,
	Legends:                     nil,
	Role:                        "img",
}
