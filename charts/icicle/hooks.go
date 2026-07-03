package icicle

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3hierarchy "github.com/geoffjay/templ-charts/internal/d3/hierarchy"
)

// UseIcicle mirrors @nivo/icicle's useIcicle: it builds a hierarchy, sums leaf
// values, runs the partition layout, and maps each node to a rectangle oriented
// per props.Orientation. Colors are grouped by the depth-1 ancestor.
//
// When props.FocusID names a non-root node, the layout is rescaled d3's
// zoomable-icicle way: the spread axis is remapped so the focused node's
// [x0,x1] fills the whole spread extent, and the depth axis so the focused
// node's depth becomes the near edge and the remaining depth fills the area.
// Nodes outside the focused subtree (except its ancestors, which drive the
// breadcrumb) collapse to zero size and are skipped.
func UseIcicle(props IcicleProps) IcicleResult {
	root := d3hierarchy.Hierarchy(props.Data, icicleChildren).
		Sum(func(d any) float64 { return d.(IcicleNode).Value })

	horizontal := props.Orientation == OrientationLeft || props.Orientation == OrientationRight
	// partition spread axis is x; depth axis is y. For left/right orientations
	// the spread runs down the height and depth runs across the width.
	spread, depth := props.Width, props.Height
	if horizontal {
		spread, depth = props.Height, props.Width
	}
	d3hierarchy.NewPartition().Size(spread, depth).Layout(root)

	// Focused rescale (classic d3 zoomable-icicle). focus == nil ⇒ identity, so
	// the un-focused branch is byte-identical to the pre-zoom output.
	focus := findNode(root, props.FocusID)
	if focus != nil && focus.Depth > 0 {
		fx0, fx1 := focus.X0, focus.X1
		xk := spread
		if fx1 > fx0 {
			xk = spread / (fx1 - fx0)
		}
		fy0 := focus.Y0
		yk := depth
		if depth > fy0 {
			yk = depth / (depth - fy0)
		}
		root.Each(func(n *d3hierarchy.Node) {
			n.X0 = (n.X0 - fx0) * xk
			n.X1 = (n.X1 - fx0) * xk
			n.Y0 = (n.Y0 - fy0) * yk
			n.Y1 = (n.Y1 - fy0) * yk
		})
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	visible := visibleSet(focus)

	rects := make([]ComputedRect, 0)
	root.Each(func(n *d3hierarchy.Node) {
		if visible != nil && !visible[n] {
			return
		}
		var x, y, w, h float64
		switch props.Orientation {
		case OrientationTop:
			x, y = n.X0, depth-n.Y1
			w, h = n.X1-n.X0, n.Y1-n.Y0
		case OrientationLeft:
			x, y = n.Y0, n.X0
			w, h = n.Y1-n.Y0, n.X1-n.X0
		case OrientationRight:
			x, y = depth-n.Y1, n.X0
			w, h = n.Y1-n.Y0, n.X1-n.X0
		default: // bottom
			x, y = n.X0, n.Y0
			w, h = n.X1-n.X0, n.Y1-n.Y0
		}
		w = math.Max(0, w-props.GapX)
		h = math.Max(0, h-props.GapY)
		rects = append(rects, ComputedRect{
			ID:             n.Data.(IcicleNode).ID,
			Depth:          n.Depth,
			Value:          n.Value,
			FormattedValue: format(n.Value),
			X:              x, Y: y, Width: w, Height: h,
			Color: getColor(colorGroup(n)),
		})
	})
	return IcicleResult{Rects: rects, Breadcrumb: breadcrumbOf(focus, func(n *d3hierarchy.Node) string {
		return n.Data.(IcicleNode).ID
	})}
}

// findNode returns the node whose IcicleNode.ID matches id, or nil if id is
// empty or unmatched (including the root, which is never a zoom target).
func findNode(root *d3hierarchy.Node, id string) *d3hierarchy.Node {
	if id == "" {
		return nil
	}
	var found *d3hierarchy.Node
	root.Each(func(n *d3hierarchy.Node) {
		if found == nil && n.Data.(IcicleNode).ID == id {
			found = n
		}
	})
	return found
}

// visibleSet returns the set of nodes to render when focused: the focus node
// plus all its descendants. Returns nil when focus is nil (render everything).
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

// breadcrumbOf returns the root→focus ancestor path as crumbs (each carrying
// its node id), or nil when focus is nil/root. The label is the node id.
func breadcrumbOf(focus *d3hierarchy.Node, id func(*d3hierarchy.Node) string) []Crumb {
	if focus == nil || focus.Depth == 0 {
		return nil
	}
	anc := focus.Ancestors() // focus … root
	crumbs := make([]Crumb, 0, len(anc))
	for i := len(anc) - 1; i >= 0; i-- { // root … focus
		crumbs = append(crumbs, Crumb{ID: id(anc[i]), Label: id(anc[i])})
	}
	return crumbs
}

func icicleChildren(d any) []any {
	n := d.(IcicleNode)
	if len(n.Children) == 0 {
		return nil
	}
	out := make([]any, len(n.Children))
	for i, c := range n.Children {
		out[i] = c
	}
	return out
}

func colorGroup(n *d3hierarchy.Node) string {
	cur := n
	for cur.Depth > 1 {
		cur = cur.Parent
	}
	return cur.Data.(IcicleNode).ID
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
