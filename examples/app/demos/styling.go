package demos

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
)

// StylingDemos returns the /styling page tiles: gradients, patterns, match
// rules, and blend modes via the Defs + Fill props (mirroring nivo's
// "Patterns & Gradients" guide). Each def carries an "inherit" color stop
// where per-series variants should be generated from the node's own color.
func StylingDemos() []Demo {
	return []Demo{
		{
			ID:          "styling-line-gradient-area",
			Title:       "Gradient area",
			Description: "EnableArea + a linearGradient def with \"inherit\" stops: each series' area fades from its own color to transparent. Fill: [{ID, Match: \"*\"}] binds the def to every series.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: defaultMargin(),
				Curve:  core.CurveMonotoneX,
				Data:   stylingLineData(),

				EnableArea:  core.BoolPtr(true),
				AreaOpacity: 1,
				LineWidth:   3,

				Defs: []core.Def{
					core.LinearGradientDef("lineAreaGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 0.55},
						{Offset: 100, Color: "inherit", Opacity: 0.02},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "lineAreaGrad", Match: "*"}},

				UseMesh:         true,
				EnableCrosshair: core.BoolPtr(true),
			},
		},
		{
			ID:          "styling-bar-gradient",
			Title:       "Gradient bars",
			Description: "Grouped bars with a vertical gradient generated from each key's color (\"inherit\" stops, full opacity at the top fading toward the baseline) plus a border radius.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:     commonChartWidth,
				Height:    commonChartHeight,
				Margin:    defaultMargin(),
				IndexBy:   "country",
				Keys:      []string{"hot dogs", "burgers", "kebab"},
				Data:      barData(),
				GroupMode: bar.GroupModeGrouped,

				BorderRadius: 4,
				Defs: []core.Def{
					core.LinearGradientDef("barGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 1},
						{Offset: 100, Color: "inherit", Opacity: 0.35},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "barGrad", Match: "*"}},
			},
		},
		{
			ID:          "styling-bar-duotone",
			Title:       "Duotone gradient",
			Description: "One fixed two-color gradient (indigo → cyan) applied to every bar via Match: \"*\" — a branded, monochrome look independent of the color scale.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				Margin:  defaultMargin(),
				IndexBy: "country",
				Keys:    []string{"burgers"},
				Data:    barData(),

				BorderRadius: 4,
				Defs: []core.Def{
					core.LinearGradientDef("duotone", []core.GradientStop{
						{Offset: 0, Color: "#6366f1", Opacity: 1},
						{Offset: 100, Color: "#22d3ee", Opacity: 1},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "duotone", Match: "*"}},
			},
		},
		{
			ID:          "styling-bar-patterns",
			Title:       "Pattern fills + match rules",
			Description: "Conditional fills: dots on the \"fries\" key and diagonal lines on \"donut\" (patterns keep each key's color via Background: \"inherit\"); the other keys stay solid. Match rules compare against the bar's datum (dataKey \"data\").",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				Margin:  defaultMargin(),
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
				Data:    barData(),

				Defs: []core.Def{
					{
						ID: "dots", Type: core.DefTypePatternDots,
						Background: "inherit", Color: "rgba(255,255,255,0.55)",
						Size: 4, Padding: 3, Stagger: true,
					},
					{
						ID: "lines", Type: core.DefTypePatternLines,
						Background: "inherit", Color: "rgba(255,255,255,0.45)",
						Spacing: 7, Rotation: -45, LineWidth: 3,
					},
				},
				Fill: []core.DefRule{
					{ID: "dots", Match: map[string]any{"id": "fries"}},
					{ID: "lines", Match: map[string]any{"id": "donut"}},
				},
			},
		},
		{
			ID:          "styling-pie-gradient",
			Title:       "Gradient donut",
			Description: "Each arc fades radially* from its own color via \"inherit\" gradient stops. (*SVG linearGradient, top-to-bottom — the fade direction is shared by all arcs.)",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:        commonChartWidth,
				Height:       commonChartHeight,
				Margin:       defaultMargin(),
				InnerRadius:  0.55,
				PadAngle:     0.7,
				CornerRadius: 4,
				Data:         pieData(),

				ActiveOuterRadiusOffset: 8,
				Defs: []core.Def{
					core.LinearGradientDef("pieGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 1},
						{Offset: 100, Color: "inherit", Opacity: 0.45},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "pieGrad", Match: "*"}},
			},
		},
		{
			ID:          "styling-line-blend",
			Title:       "Area blend mode",
			Description: "AreaBlendMode: \"multiply\" — overlapping areas mix like ink instead of occluding, so every series stays readable without transparency tricks.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: defaultMargin(),
				Curve:  core.CurveMonotoneX,
				Data:   stylingLineData(),

				EnableArea:    core.BoolPtr(true),
				AreaOpacity:   0.55,
				AreaBlendMode: core.MixBlendMultiply,
				LineWidth:     2,
				Colors:        colors.Scheme(colors.PaletteSet2),
			},
		},
	}
}

// stylingLineData returns three deliberately curvy, offset series so gradient
// areas and blend modes read clearly (the shared samples.Line series are
// near-parallel diagonals, which muddies area-fill demos).
func stylingLineData() []line.LineSeries {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	values := map[string][]float64{
		"mobile":  {42, 58, 51, 74, 68, 92, 105, 98, 84, 96, 78, 88},
		"desktop": {68, 62, 79, 65, 88, 76, 70, 84, 97, 82, 94, 76},
		"tablet":  {25, 32, 28, 41, 36, 48, 42, 55, 46, 38, 50, 44},
	}
	var out []line.LineSeries
	for _, id := range []string{"mobile", "desktop", "tablet"} {
		s := line.LineSeries{ID: id}
		for i, m := range months {
			s.Data = append(s.Data, line.LinePointData{X: m, Y: values[id][i]})
		}
		out = append(out, s)
	}
	return out
}
