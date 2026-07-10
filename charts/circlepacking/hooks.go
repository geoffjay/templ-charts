package circlepacking

import (
	"math"
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

	// Focused zoom: the d3 zoomable-pack transform — translate/scale so the
	// focused node's circle fills the viewport (k = size / (2*F.r)), applied to
	// every node. focus == nil ⇒ identity (byte-identical to pre-zoom output).
	focus := findNode(root, props.FocusID)
	if focus != nil && focus.Depth > 0 && focus.R > 0 {
		size := math.Min(props.Width, props.Height)
		k := size / (2 * focus.R)
		cx, cy := props.Width/2, props.Height/2
		fx, fy := focus.X, focus.Y
		root.Each(func(n *d3hierarchy.Node) {
			n.X = (n.X-fx)*k + cx
			n.Y = (n.Y-fy)*k + cy
			n.R *= k
		})
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(k string) string { return k })
	format := valueFormatter(props.ValueFormat)

	visible := visibleSet(focus)

	circles := make([]ComputedCircle, 0)
	root.Each(func(n *d3hierarchy.Node) {
		if visible != nil && !visible[n] {
			return
		}
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
	return CirclePackingResult{Circles: circles, Breadcrumb: breadcrumbOf(focus, func(n *d3hierarchy.Node) string {
		return n.Data.(CirclePackingNode).ID
	})}
}

// findNode returns the node whose CirclePackingNode.ID matches id, or nil if id
// is empty/unmatched (the root is never a zoom target).
func findNode(root *d3hierarchy.Node, id string) *d3hierarchy.Node {
	if id == "" {
		return nil
	}
	var found *d3hierarchy.Node
	root.Each(func(n *d3hierarchy.Node) {
		if found == nil && n.Data.(CirclePackingNode).ID == id {
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
