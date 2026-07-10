package tree

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3hierarchy "github.com/geoffjay/templ-charts/internal/d3/hierarchy"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseTree mirrors @nivo/tree's useTree: it builds a hierarchy, runs the tidy
// tree or dendrogram layout, orients the result per props.Layout, colors nodes
// by id, and builds smooth bump links from each parent to its children.
func UseTree(props TreeProps) TreeResult {
	root := d3hierarchy.Hierarchy(props.Data, treeChildren)

	horizontal := props.Layout == LayoutLeftToRight || props.Layout == LayoutRightToLeft
	// Layout spread axis is X, depth axis is Y. For horizontal orientations the
	// spread runs down the height and depth across the width.
	if horizontal {
		layoutRun(props.Mode, root, props.Height, props.Width)
	} else {
		layoutRun(props.Mode, root, props.Width, props.Height)
	}

	colorsCfg := props.Colors
	if colorsCfg.Type == 0 && colorsCfg.Scheme == "" && colorsCfg.Static == "" && len(colorsCfg.Colors) == 0 && colorsCfg.Func == nil && colorsCfg.DatumPath == "" {
		colorsCfg = Defaults.Colors
	}
	getColor := colors.GetOrdinalColorScale[string](colorsCfg, func(id string) string { return id })

	screen := make(map[*d3hierarchy.Node][2]float64)
	nodes := make([]ComputedNode, 0)
	root.Each(func(n *d3hierarchy.Node) {
		sx, sy := toScreen(props, n)
		screen[n] = [2]float64{sx, sy}
		nodes = append(nodes, ComputedNode{
			ID:    n.Data.(TreeNode).ID,
			Depth: n.Depth,
			X:     sx, Y: sy,
			Color: getColor(n.Data.(TreeNode).ID),
		})
	})

	curve := d3shape.CurveBumpY
	if horizontal {
		curve = d3shape.CurveBumpX
	}
	gen := d3shape.NewLine().Curve(curve)

	links := make([]ComputedLink, 0)
	for _, l := range root.Links() {
		s := screen[l.Source]
		t := screen[l.Target]
		links = append(links, ComputedLink{
			SourceID: l.Source.Data.(TreeNode).ID,
			TargetID: l.Target.Data.(TreeNode).ID,
			Path:     gen.Call([]d3shape.Point2D{{s[0], s[1]}, {t[0], t[1]}}),
			Color:    getColor(l.Source.Data.(TreeNode).ID),
		})
	}

	return TreeResult{Nodes: nodes, Links: links}
}

func layoutRun(mode Mode, root *d3hierarchy.Node, dx, dy float64) {
	if mode == ModeTree {
		d3hierarchy.NewTree().Size(dx, dy).Layout(root)
	} else {
		d3hierarchy.NewCluster().Size(dx, dy).Layout(root)
	}
}

// toScreen maps a laid-out node (X=spread, Y=depth) to screen coordinates per
// the orientation.
func toScreen(props TreeProps, n *d3hierarchy.Node) (float64, float64) {
	switch props.Layout {
	case LayoutBottomToTop:
		return n.X, props.Height - n.Y
	case LayoutLeftToRight:
		return n.Y, n.X
	case LayoutRightToLeft:
		return props.Width - n.Y, n.X
	default: // top-to-bottom
		return n.X, n.Y
	}
}

func treeChildren(d any) []any {
	n := d.(TreeNode)
	if len(n.Children) == 0 {
		return nil
	}
	out := make([]any, len(n.Children))
	for i, c := range n.Children {
		out[i] = c
	}
	return out
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
