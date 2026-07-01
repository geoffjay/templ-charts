package circlepacking

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3hierarchy "github.com/geoffjay/templ-charts/internal/d3/hierarchy"
)

// UseCirclePacking mirrors @nivo/circle-packing's useCirclePacking: it builds a
// hierarchy, sums leaf values, runs the pack layout, and produces positioned,
// colored circles. Colors are keyed by depth (nivo's colorBy:'depth').
func UseCirclePacking(props CirclePackingProps) CirclePackingResult {
	root := d3hierarchy.Hierarchy(props.Data, cpChildren).
		Sum(func(d any) float64 { return d.(CirclePackingNode).Value })

	d3hierarchy.NewPack().Size(props.Width, props.Height).Padding(props.Padding).Layout(root)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(k string) string { return k })
	format := valueFormatter(props.ValueFormat)

	circles := make([]ComputedCircle, 0)
	root.Each(func(n *d3hierarchy.Node) {
		circles = append(circles, ComputedCircle{
			ID:             n.Data.(CirclePackingNode).ID,
			Depth:          n.Depth,
			Value:          n.Value,
			FormattedValue: format(n.Value),
			X:              n.X, Y: n.Y, R: n.R,
			Color:  getColor(strconv.Itoa(n.Depth)),
			IsLeaf: len(n.Children) == 0,
		})
	})
	return CirclePackingResult{Circles: circles}
}

func cpChildren(d any) []any {
	n := d.(CirclePackingNode)
	if len(n.Children) == 0 {
		return nil
	}
	out := make([]any, len(n.Children))
	for i, c := range n.Children {
		out[i] = c
	}
	return out
}

func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
