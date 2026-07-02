package sankey

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3sankey "github.com/geoffjay/templ-charts/internal/d3/sankey"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseSankey mirrors @nivo/sankey's useSankey/computeNodeAndLinks: it runs the
// d3-sankey layout over the input graph (with alignment/sort/thickness/spacing
// from props), then derives the render rectangles, colors, and ribbon paths.
// For a vertical layout the d3 coordinates are transposed exactly as nivo does.
// props.Width/Height are the inner dimensions.
func UseSankey(props SankeyProps) SankeyResult {
	horizontal := props.Layout != SankeyLayoutVertical

	g := &d3sankey.Graph{
		Nodes: make([]*d3sankey.Node, len(props.Nodes)),
		Links: make([]*d3sankey.Link, len(props.Links)),
	}
	for i, n := range props.Nodes {
		g.Nodes[i] = &d3sankey.Node{ID: n.ID}
	}
	for i, l := range props.Links {
		g.Links[i] = &d3sankey.Link{SourceID: l.Source, TargetID: l.Target, Value: l.Value}
	}

	// d3-sankey lays out in [width, height] for horizontal and [height, width]
	// for vertical (nivo transposes back afterwards).
	sw, sh := props.Width, props.Height
	if !horizontal {
		sw, sh = props.Height, props.Width
	}

	s := d3sankey.New().
		NodeAlign(alignFunc(props.Align)).
		NodeWidth(props.NodeThickness).
		NodePadding(props.NodeSpacing).
		Size(sw, sh)
	applySort(s, props.Sort)
	s.Layout(g)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	nodes := make([]ComputedNode, len(g.Nodes))
	nodeColor := make(map[string]string, len(g.Nodes))
	for i, n := range g.Nodes {
		color := getColor(n.ID)
		nodeColor[n.ID] = color
		cn := ComputedNode{
			ID:             n.ID,
			Label:          n.ID,
			FormattedValue: format(n.Value),
			Value:          n.Value,
			Color:          color,
		}
		if horizontal {
			cn.X = n.X0 + props.NodeInnerPadding
			cn.Y = n.Y0
			cn.Width = math.Max(n.X1-n.X0-props.NodeInnerPadding*2, 0)
			cn.Height = math.Max(n.Y1-n.Y0, 0)
			cn.X0, cn.Y0, cn.X1, cn.Y1 = n.X0, n.Y0, n.X1, n.Y1
		} else {
			// Transpose: d3's x-axis (layers) becomes screen y, d3's y-axis
			// becomes screen x. Mirrors @nivo/sankey's coordinate swap.
			cn.X = n.Y0
			cn.Y = n.X0 + props.NodeInnerPadding
			cn.Width = math.Max(n.Y1-n.Y0, 0)
			cn.Height = math.Max(n.X1-n.X0-props.NodeInnerPadding*2, 0)
			cn.X0, cn.Y0, cn.X1, cn.Y1 = n.Y0, n.X0, n.Y1, n.X1
		}
		nodes[i] = cn
	}
	nodeIdx := make(map[string]int, len(nodes))
	for i, n := range nodes {
		nodeIdx[n.ID] = i
	}

	links := make([]ComputedLink, 0, len(g.Links))
	for _, l := range g.Links {
		if l.Source == nil || l.Target == nil {
			continue
		}
		src := nodes[nodeIdx[l.Source.ID]]
		tgt := nodes[nodeIdx[l.Target.ID]]
		cl := ComputedLink{
			Index:          l.Index,
			Source:         l.Source.ID,
			Target:         l.Target.ID,
			Value:          l.Value,
			FormattedValue: format(l.Value),
			Color:          src.Color,
			StartColor:     src.Color,
			EndColor:       tgt.Color,
			Pos0:           l.Y0,
			Pos1:           l.Y1,
			Thickness:      l.Width,
		}
		cl.Path = linkPath(horizontal, src, tgt, cl.Pos0, cl.Pos1, cl.Thickness, props.LinkContract)
		links = append(links, cl)
	}

	legendData := make([]legends.Datum, 0, len(nodes))
	for _, n := range nodes {
		legendData = append(legendData, legends.Datum{ID: n.ID, Label: n.Label, Color: n.Color})
	}

	return SankeyResult{Nodes: nodes, Links: links, LegendData: legendData}
}

// linkPath builds the ribbon path exactly as @nivo/sankey does: a line
// generator with curveMonotoneX (horizontal) / curveMonotoneY (vertical) over a
// 9-point closed outline whose middle control points sit 12% in from each end,
// with a trailing "Z".
func linkPath(horizontal bool, src, tgt ComputedNode, pos0, pos1, thickness, contract float64) string {
	th := math.Max(1, thickness-contract*2)
	h := th / 2
	if horizontal {
		x0 := src.X1
		x1 := tgt.X0
		pad := (x1 - x0) * 0.12
		dots := []d3shape.Point2D{
			{x0, pos0 - h},
			{x0 + pad, pos0 - h},
			{x1 - pad, pos1 - h},
			{x1, pos1 - h},
			{x1, pos1 + h},
			{x1 - pad, pos1 + h},
			{x0 + pad, pos0 + h},
			{x0, pos0 + h},
			{x0, pos0 - h},
		}
		return d3shape.NewLine().Curve(d3shape.CurveMonotoneX).Call(dots) + "Z"
	}
	y0 := src.Y1
	y1 := tgt.Y0
	pad := (y1 - y0) * 0.12
	dots := []d3shape.Point2D{
		{pos0 + h, y0},
		{pos0 + h, y0 + pad},
		{pos1 + h, y1 - pad},
		{pos1 + h, y1},
		{pos1 - h, y1},
		{pos1 - h, y1 - pad},
		{pos0 - h, y0 + pad},
		{pos0 - h, y0},
		{pos0 + h, y0},
	}
	return d3shape.NewLine().Curve(d3shape.CurveMonotoneY).Call(dots) + "Z"
}

// alignFunc maps the Align prop to a d3-sankey alignment function.
func alignFunc(a SankeyAlign) d3sankey.AlignFunc {
	switch a {
	case SankeyAlignJustify:
		return d3sankey.SankeyJustify
	case SankeyAlignStart:
		return d3sankey.SankeyLeft
	case SankeyAlignEnd:
		return d3sankey.SankeyRight
	default:
		return d3sankey.SankeyCenter
	}
}

// applySort maps the Sort prop to the port's node/link sort mode, matching
// nivo: auto → d3 default, input → preserve input order, ascending/descending →
// order columns by node value.
func applySort(s *d3sankey.Sankey, sort SankeySort) {
	switch sort {
	case SankeySortInput:
		s.NodeSortInput()
		s.LinkSortInput(true)
	case SankeySortAscending:
		s.NodeSort(func(a, b *d3sankey.Node) int { return cmpFloat(a.Value, b.Value) })
	case SankeySortDescending:
		s.NodeSort(func(a, b *d3sankey.Node) int { return cmpFloat(b.Value, a.Value) })
	default:
		s.NodeSortAuto()
	}
}

func cmpFloat(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
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
