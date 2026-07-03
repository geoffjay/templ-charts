package treemap

import (
	"strconv"
	"strings"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3hierarchy "github.com/geoffjay/templ-charts/internal/d3/hierarchy"
)

// UseTreemap mirrors @nivo/treemap's useTreeMap: it builds a hierarchy from the
// data, sums leaf values, runs the treemap layout, and produces positioned,
// colored nodes. Colors are grouped by the depth-1 ancestor (nivo's
// colorBy:'pathComponents.1'). props.Width/Height are the inner dimensions.
func UseTreemap(props TreemapProps) TreemapResult {
	root := d3hierarchy.Hierarchy(props.Data, treemapChildren).
		Sum(func(d any) float64 { return d.(TreemapNode).Value })

	tm := d3hierarchy.NewTreemap().
		Tile(tileFunc(props.Tile)).
		Size(props.Width, props.Height).
		PaddingInner(props.InnerPadding).
		PaddingOuter(props.OuterPadding)
	tm.Layout(root)

	// Focused zoom: re-run the treemap layout on the subtree rooted at the
	// focused node over the full inner rect (re-sum is unnecessary — Sum already
	// set the subtree values). Only the focus subtree is rendered. focus == nil
	// ⇒ identity, so the un-focused output is byte-identical to before.
	focus := findNode(root, props.FocusID)
	layoutRoot := root
	if focus != nil && focus.Depth > 0 {
		tm.Layout(focus)
		layoutRoot = focus
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	nodes := make([]ComputedNode, 0)
	layoutRoot.Each(func(n *d3hierarchy.Node) {
		id := n.Data.(TreemapNode).ID
		nodes = append(nodes, ComputedNode{
			ID:             id,
			Path:           nodePath(n),
			Depth:          n.Depth,
			Value:          n.Value,
			FormattedValue: format(n.Value),
			X:              n.X0,
			Y:              n.Y0,
			Width:          n.X1 - n.X0,
			Height:         n.Y1 - n.Y0,
			Color:          getColor(colorGroup(n)),
			IsParent:       len(n.Children) > 0,
		})
	})

	return TreemapResult{Nodes: nodes, Breadcrumb: breadcrumbOf(focus, func(n *d3hierarchy.Node) string {
		return n.Data.(TreemapNode).ID
	})}
}

// findNode returns the node whose TreemapNode.ID matches id, or nil if id is
// empty/unmatched (the root is never a zoom target).
func findNode(root *d3hierarchy.Node, id string) *d3hierarchy.Node {
	if id == "" {
		return nil
	}
	var found *d3hierarchy.Node
	root.Each(func(n *d3hierarchy.Node) {
		if found == nil && n.Data.(TreemapNode).ID == id {
			found = n
		}
	})
	return found
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

// treemapChildren adapts TreemapNode children for the hierarchy builder.
func treemapChildren(d any) []any {
	n := d.(TreemapNode)
	if len(n.Children) == 0 {
		return nil
	}
	out := make([]any, len(n.Children))
	for i, c := range n.Children {
		out[i] = c
	}
	return out
}

// colorGroup returns the id of the node's depth-1 ancestor (pathComponents.1);
// nodes at depth <= 1 use their own id. This groups a subtree under one color.
func colorGroup(n *d3hierarchy.Node) string {
	if n.Depth <= 1 {
		return n.Data.(TreemapNode).ID
	}
	cur := n
	for cur.Depth > 1 {
		cur = cur.Parent
	}
	return cur.Data.(TreemapNode).ID
}

// nodePath returns the dot-joined ancestor ids from root to node.
func nodePath(n *d3hierarchy.Node) string {
	var parts []string
	for cur := n; cur != nil; cur = cur.Parent {
		parts = append(parts, cur.Data.(TreemapNode).ID)
	}
	// reverse
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

func tileFunc(t TileType) d3hierarchy.TileFunc {
	switch t {
	case TileBinary:
		return d3hierarchy.TreemapBinary
	case TileDice:
		return d3hierarchy.TreemapDice
	case TileSlice:
		return d3hierarchy.TreemapSlice
	case TileSliceDice:
		return d3hierarchy.TreemapSliceDice
	default:
		return d3hierarchy.TreemapSquarify
	}
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
