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
func UseIcicle(props IcicleProps) IcicleResult {
	root := d3hierarchy.Hierarchy(props.Data, icicleChildren).
		Sum(func(d any) float64 { return d.(IcicleNode).Value })

	horizontal := props.Orientation == OrientationLeft || props.Orientation == OrientationRight
	// partition spread axis is x; depth axis is y. For left/right orientations
	// the spread runs down the height and depth runs across the width.
	if horizontal {
		d3hierarchy.NewPartition().Size(props.Height, props.Width).Layout(root)
	} else {
		d3hierarchy.NewPartition().Size(props.Width, props.Height).Layout(root)
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	rects := make([]ComputedRect, 0)
	root.Each(func(n *d3hierarchy.Node) {
		var x, y, w, h float64
		switch props.Orientation {
		case OrientationTop:
			x, y = n.X0, props.Height-n.Y1
			w, h = n.X1-n.X0, n.Y1-n.Y0
		case OrientationLeft:
			x, y = n.Y0, n.X0
			w, h = n.Y1-n.Y0, n.X1-n.X0
		case OrientationRight:
			x, y = props.Width-n.Y1, n.X0
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
	return IcicleResult{Rects: rects}
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
