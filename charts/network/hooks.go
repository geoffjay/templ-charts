package network

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/d3/force"
)

// UseNetwork mirrors @nivo/network's useNetwork: it seeds a force simulation
// with the input nodes/links, attaches the link, many-body (charge) and
// centering forces the same way nivo does, runs a fixed number of iterations,
// then reads back the final node positions and derives the link segments.
// props.Width/Height are the inner dimensions.
func UseNetwork(props NetworkProps) NetworkResult {
	n := len(props.Nodes)
	fnodes := make([]*force.Node, n)
	byID := make(map[string]*force.Node, n)
	idx := make(map[string]int, n)
	for i, in := range props.Nodes {
		node := force.NewNode(in)
		fnodes[i] = node
		byID[in.ID] = node
		idx[in.ID] = i
	}

	flinks := make([]*force.Link, 0, len(props.Links))
	for _, l := range props.Links {
		src, sok := byID[l.Source]
		tgt, tok := byID[l.Target]
		if !sok || !tok {
			continue // skip links referencing unknown nodes, as d3 would error
		}
		flinks = append(flinks, force.NewLink(src, tgt, l))
	}

	distanceMax := props.DistanceMax
	if distanceMax <= 0 {
		distanceMax = math.Inf(1)
	}

	sim := force.NewSimulation(fnodes).
		Force("link", force.ForceLink(flinks).
			Distance(func(l *force.Link) float64 {
				if d := l.Data.(NetworkInputLink).Distance; d > 0 {
					return d
				}
				return props.LinkDistance
			}).
			StrengthConst(props.CenteringStrength)).
		Force("charge", force.ForceManyBody().
			StrengthConst(-props.Repulsivity).
			DistanceMin(props.DistanceMin).
			DistanceMax(distanceMax)).
		Force("center", force.ForceCenter(props.Width/2, props.Height/2)).
		Stop()

	sim.Tick(props.Iterations)

	nodes := make([]ComputedNode, n)
	for i, fn := range sim.Nodes() {
		in := fn.Data.(NetworkInputNode)
		size := in.Size
		if size <= 0 {
			size = props.NodeSize
		}
		color := in.Color
		if color == "" {
			color = props.NodeColor
		}
		border := props.NodeBorderColor
		if border == "" {
			border = color
		}
		nodes[i] = ComputedNode{
			ID:          in.ID,
			Index:       fn.Index,
			X:           fn.X,
			Y:           fn.Y,
			Size:        size,
			Color:       color,
			BorderWidth: props.NodeBorderWidth,
			BorderColor: border,
		}
	}

	if props.FitView {
		fitNodes(nodes, props.Width, props.Height)
	}

	links := make([]ComputedLink, 0, len(flinks))
	for _, fl := range flinks {
		s := nodes[idx[fl.Source.Data.(NetworkInputNode).ID]]
		t := nodes[idx[fl.Target.Data.(NetworkInputNode).ID]]
		color := props.LinkColor
		if color == "" {
			color = s.Color
		}
		links = append(links, ComputedLink{
			ID:        s.ID + "." + t.ID,
			Source:    s.ID,
			Target:    t.ID,
			X1:        s.X,
			Y1:        s.Y,
			X2:        t.X,
			Y2:        t.Y,
			Thickness: props.LinkThickness,
			Color:     color,
		})
	}

	return NetworkResult{Nodes: nodes, Links: links}
}

// fitNodes rescales node centers to fill the width×height area, preserving
// aspect ratio and insetting by the largest node radius so no node clips. It
// mutates the nodes in place; the caller derives links afterwards so they
// follow. A degenerate (zero-span) axis is centered rather than scaled.
func fitNodes(nodes []ComputedNode, width, height float64) {
	if len(nodes) == 0 {
		return
	}
	minX, maxX := nodes[0].X, nodes[0].X
	minY, maxY := nodes[0].Y, nodes[0].Y
	maxR := 0.0
	for _, n := range nodes {
		if n.X < minX {
			minX = n.X
		}
		if n.X > maxX {
			maxX = n.X
		}
		if n.Y < minY {
			minY = n.Y
		}
		if n.Y > maxY {
			maxY = n.Y
		}
		if r := n.Size / 2; r > maxR {
			maxR = r
		}
	}

	availW := width - 2*maxR
	availH := height - 2*maxR
	if availW < 0 {
		availW = 0
	}
	if availH < 0 {
		availH = 0
	}
	spanX := maxX - minX
	spanY := maxY - minY

	scale := math.Inf(1)
	if spanX > 0 {
		scale = math.Min(scale, availW/spanX)
	}
	if spanY > 0 {
		scale = math.Min(scale, availH/spanY)
	}
	if math.IsInf(scale, 1) {
		scale = 1 // single point (or fully coincident): keep scale, just center
	}

	offX := (width - spanX*scale) / 2
	offY := (height - spanY*scale) / 2
	for i := range nodes {
		nodes[i].X = offX + (nodes[i].X-minX)*scale
		nodes[i].Y = offY + (nodes[i].Y-minY)*scale
	}
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
