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

	// Focused rescale (d3 zoomable-sunburst): remap the angular span so the
	// focused node's [x0,x1] fills the full 2π, and the radial (r²) span so the
	// focused node's depth becomes the inner edge and its subtree fills to the
	// outer radius. focus == nil ⇒ identity (byte-identical to pre-zoom output).
	focus := findNode(root, props.FocusID)
	if focus != nil && focus.Depth > 0 {
		fx0, fx1 := focus.X0, focus.X1
		xk := 2 * math.Pi
		if fx1 > fx0 {
			xk = 2 * math.Pi / (fx1 - fx0)
		}
		fy0 := focus.Y0
		maxY := radius * radius
		yk := 1.0
		if maxY > fy0 {
			yk = maxY / (maxY - fy0)
		}
		root.Each(func(n *d3hierarchy.Node) {
			n.X0 = (n.X0 - fx0) * xk
			n.X1 = (n.X1 - fx0) * xk
			n.Y0 = math.Max(0, n.Y0-fy0) * yk
			n.Y1 = math.Max(0, n.Y1-fy0) * yk
		})
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	visible := visibleSet(focus)

	out := make([]ComputedArc, 0)
	root.Each(func(n *d3hierarchy.Node) {
		if n.Depth == 0 {
			return // root is the hollow center
		}
		if visible != nil && !visible[n] {
			return
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

	return SunburstResult{Center: center, Arcs: out, Breadcrumb: breadcrumbOf(focus, func(n *d3hierarchy.Node) string {
		return n.Data.(SunburstNode).ID
	})}
}

// findNode returns the node whose SunburstNode.ID matches id, or nil if id is
// empty/unmatched (the root is never a zoom target).
func findNode(root *d3hierarchy.Node, id string) *d3hierarchy.Node {
	if id == "" {
		return nil
	}
	var found *d3hierarchy.Node
	root.Each(func(n *d3hierarchy.Node) {
		if found == nil && n.Data.(SunburstNode).ID == id {
			found = n
		}
	})
	return found
}

// visibleSet returns the focus node plus its descendants; nil when unfocused.
func visibleSet(focus *d3hierarchy.Node) map[*d3hierarchy.Node]bool {
	if focus == nil || focus.Depth == 0 {
		return nil
	}
	set := map[*d3hierarchy.Node]bool{}
	for _, d := range focus.Descendants() {
		set[d] = true
	}
	return set
}

// breadcrumbOf returns the root→focus ancestor path as crumbs; nil when
// unfocused.
func breadcrumbOf(focus *d3hierarchy.Node, id func(*d3hierarchy.Node) string) []Crumb {
	if focus == nil || focus.Depth == 0 {
		return nil
	}
	anc := focus.Ancestors()
	crumbs := make([]Crumb, 0, len(anc))
	for i := len(anc) - 1; i >= 0; i-- {
		crumbs = append(crumbs, Crumb{ID: id(anc[i]), Label: id(anc[i])})
	}
	return crumbs
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
