package sunburst

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3hierarchy "github.com/geoffjay/templ-charts/internal/d3/hierarchy"
)

// UseSunburst mirrors @nivo/sunburst's useSunburst: it builds a hierarchy, sums
// leaf values, runs the partition layout over [2π, r²], and maps each node to
// an arc. Colors are grouped by the depth-1 ancestor (inheritColorFromParent).
// The root node (depth 0) is the center and is not drawn.
func UseSunburst(props SunburstProps) SunburstResult {
	center := [2]float64{props.Width / 2, props.Height / 2}
	radius := math.Min(props.Width, props.Height) / 2

	root := d3hierarchy.Hierarchy(props.Data, sunburstChildren).
		Sum(func(d any) float64 { return d.(SunburstNode).Value })

	d3hierarchy.NewPartition().Size(2*math.Pi, radius*radius).Layout(root)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	out := make([]ComputedArc, 0)
	root.Each(func(n *d3hierarchy.Node) {
		if n.Depth == 0 {
			return // root is the hollow center
		}
		out = append(out, ComputedArc{
			ID:             n.Data.(SunburstNode).ID,
			Depth:          n.Depth,
			Value:          n.Value,
			FormattedValue: format(n.Value),
			Color:          getColor(colorGroup(n)),
			Arc: arcs.Arc{
				StartAngle:  n.X0,
				EndAngle:    n.X1,
				InnerRadius: math.Sqrt(n.Y0),
				OuterRadius: math.Sqrt(n.Y1),
			},
		})
	})

	return SunburstResult{Center: center, Arcs: out}
}

func sunburstChildren(d any) []any {
	n := d.(SunburstNode)
	if len(n.Children) == 0 {
		return nil
	}
	out := make([]any, len(n.Children))
	for i, c := range n.Children {
		out[i] = c
	}
	return out
}

// colorGroup returns the id of the node's depth-1 ancestor (colorBy 'id' with
// inheritColorFromParent), so a subtree shares its top-level ancestor's color.
func colorGroup(n *d3hierarchy.Node) string {
	cur := n
	for cur.Depth > 1 {
		cur = cur.Parent
	}
	return cur.Data.(SunburstNode).ID
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
