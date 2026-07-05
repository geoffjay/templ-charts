package network

import (
	"hash/fnv"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
)

// applyDefaults fills zero-valued NetworkProps fields from Defaults.
func applyDefaults(p NetworkProps) NetworkProps {
	if p.LinkDistance == 0 {
		p.LinkDistance = Defaults.LinkDistance
	}
	if p.CenteringStrength == 0 {
		p.CenteringStrength = Defaults.CenteringStrength
	}
	if p.Repulsivity == 0 {
		p.Repulsivity = Defaults.Repulsivity
	}
	if p.DistanceMin == 0 {
		p.DistanceMin = Defaults.DistanceMin
	}
	if p.Iterations == 0 {
		p.Iterations = Defaults.Iterations
	}
	if p.NodeSize == 0 {
		p.NodeSize = Defaults.NodeSize
	}
	if p.NodeColor == "" {
		p.NodeColor = Defaults.NodeColor
	}
	if p.LinkThickness == 0 {
		p.LinkThickness = Defaults.LinkThickness
	}
	if p.NodeHoverOpacity == 0 {
		p.NodeHoverOpacity = Defaults.NodeHoverOpacity
	}
	if p.NodeHoverOthersOpacity == 0 {
		p.NodeHoverOthersOpacity = Defaults.NodeHoverOthersOpacity
	}
	if p.LinkHoverOpacity == 0 {
		p.LinkHoverOpacity = Defaults.LinkHoverOpacity
	}
	if p.LinkHoverOthersOpacity == 0 {
		p.LinkHoverOthersOpacity = Defaults.LinkHoverOthersOpacity
	}
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderLayers renders the enabled layers as an inner SVG string. Annotations
// are deferred (no network annotation specs).
func renderLayers(props NetworkProps, result NetworkResult) string {
	// When interactive (and not routing hover through the mesh overlay), a scoped
	// <style> block drives the chord/sankey-style hover-highlight via CSS :has().
	cid := networkChartID(result.Nodes)
	var b strings.Builder
	if hoverHighlightOn(props) {
		b.WriteString(hoverStyleBlock(props, result, cid))
	}
	for _, layer := range props.Layers {
		switch layer {
		case NetworkLayerLinks:
			b.WriteString(renderLinksLayer(props, result, cid))
		case NetworkLayerNodes:
			b.WriteString(renderNodesLayer(props, result, cid))
		case NetworkLayerMesh:
			b.WriteString(renderMeshLayer(props, result))
		case NetworkLayerAnnotations:
			// annotations: deferred.
		}
	}
	return b.String()
}

// hoverHighlightOn reports whether the CSS :has() hover-highlight is active: it
// needs Interactive and is mutually exclusive with the voronoi mesh overlay
// (which sits atop the nodes and captures the pointer, so :has(circle:hover)
// would never fire).
func hoverHighlightOn(props NetworkProps) bool {
	return props.Interactive && !props.UseMesh
}

// networkChartID derives a stable, CSS-safe token from the node ids so the
// hover-highlight <style> rules only affect this chart instance.
func networkChartID(nodes []ComputedNode) string {
	h := fnv.New32a()
	for _, n := range nodes {
		_, _ = h.Write([]byte(n.ID))
		_, _ = h.Write([]byte{0})
	}
	return strconv.FormatUint(uint64(h.Sum32()), 16)
}

// hoverStyleBlock builds the scoped hover-highlight CSS via the shared
// charts/interact helper. Hovering a node dims every node/link, then re-lights
// the hovered node, its incident links and its neighbour nodes; hovering a link
// re-lights it and its two endpoints.
func hoverStyleBlock(props NetworkProps, result NetworkResult, cid string) string {
	scope := ".tc-nw" + cid
	idx := make(map[string]int, len(result.Nodes))
	for i, n := range result.Nodes {
		idx[n.ID] = i
	}
	// Adjacency: for each node, the neighbour node indices reached by a link.
	neighbours := make([][]int, len(result.Nodes))
	for _, l := range result.Links {
		s, sok := idx[l.Source]
		t, tok := idx[l.Target]
		if sok && tok {
			neighbours[s] = append(neighbours[s], t)
			neighbours[t] = append(neighbours[t], s)
		}
	}
	no := fmtF(props.NodeHoverOpacity)
	noo := fmtF(props.NodeHoverOthersOpacity)
	lo := fmtF(props.LinkHoverOpacity)
	loo := fmtF(props.LinkHoverOthersOpacity)

	hh := interact.HoverHighlight{Scope: scope}

	// Hover a node → dim all, then re-light the node + its neighbours + its links.
	for i := range result.Nodes {
		k := strconv.Itoa(i)
		nodeSels := []string{`.tc-nw-node.n` + k}
		for _, n := range neighbours[i] {
			nodeSels = append(nodeSels, `.tc-nw-node.n`+strconv.Itoa(n))
		}
		hh.Groups = append(hh.Groups, interact.HoverGroup{
			Trigger: `.n` + k,
			Rules: []interact.HoverRule{
				{Sels: []string{`.tc-nw-node`}, Body: `opacity:` + noo},
				{Sels: []string{`.tc-nw-link`}, Body: `opacity:` + loo},
				{Sels: []string{`.tc-nw-link.s` + k, `.tc-nw-link.t` + k}, Body: `opacity:` + lo},
				{Sels: nodeSels, Body: `opacity:` + no},
			},
		})
	}

	// Hover a link → dim all, then re-light that link and its two endpoints.
	for i, l := range result.Links {
		li := strconv.Itoa(i)
		a := strconv.Itoa(idx[l.Source])
		t := strconv.Itoa(idx[l.Target])
		hh.Groups = append(hh.Groups, interact.HoverGroup{
			Trigger: `.l` + li,
			Rules: []interact.HoverRule{
				{Sels: []string{`.tc-nw-node`}, Body: `opacity:` + noo},
				{Sels: []string{`.tc-nw-link`}, Body: `opacity:` + loo},
				{Sels: []string{`.tc-nw-link.l` + li}, Body: `opacity:` + lo},
				{Sels: []string{`.tc-nw-node.n` + a, `.tc-nw-node.n` + t}, Body: `opacity:` + no},
			},
		})
	}

	return hh.Style()
}

func renderLinksLayer(props NetworkProps, result NetworkResult, cid string) string {
	var b strings.Builder
	var idx map[string]int
	if hoverHighlightOn(props) {
		idx = make(map[string]int, len(result.Nodes))
		for i, n := range result.Nodes {
			idx[n.ID] = i
		}
	}
	for i, l := range result.Links {
		b.WriteString(`<line`)
		if hoverHighlightOn(props) {
			b.WriteString(` class="tc-nw`)
			b.WriteString(cid)
			b.WriteString(` tc-nw-link l`)
			b.WriteString(strconv.Itoa(i))
			b.WriteString(` s`)
			b.WriteString(strconv.Itoa(idx[l.Source]))
			b.WriteString(` t`)
			b.WriteString(strconv.Itoa(idx[l.Target]))
			b.WriteString(`"`)
		}
		b.WriteString(` x1="`)
		b.WriteString(fmtF(l.X1))
		b.WriteString(`" y1="`)
		b.WriteString(fmtF(l.Y1))
		b.WriteString(`" x2="`)
		b.WriteString(fmtF(l.X2))
		b.WriteString(`" y2="`)
		b.WriteString(fmtF(l.Y2))
		b.WriteString(`" stroke="`)
		b.WriteString(l.Color)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(l.Thickness))
		b.WriteString(`"></line>`)
	}
	return b.String()
}

func renderNodesLayer(props NetworkProps, result NetworkResult, cid string) string {
	var b strings.Builder
	for i, node := range result.Nodes {
		b.WriteString(`<circle`)
		if hoverHighlightOn(props) {
			b.WriteString(` class="tc-nw`)
			b.WriteString(cid)
			b.WriteString(` tc-nw-node n`)
			b.WriteString(strconv.Itoa(i))
			b.WriteString(`"`)
		}
		b.WriteString(` cx="`)
		b.WriteString(fmtF(node.X))
		b.WriteString(`" cy="`)
		b.WriteString(fmtF(node.Y))
		b.WriteString(`" r="`)
		b.WriteString(fmtF(node.Size / 2))
		b.WriteString(`" fill="`)
		b.WriteString(node.Color)
		b.WriteString(`"`)
		if node.BorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(node.BorderColor)
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(node.BorderWidth))
			b.WriteString(`"`)
		}
		// With UseMesh, the mesh overlay handles hover; skip per-node tooltips.
		if props.Interactive && !props.UseMesh {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(node.Color, node.ID, "")))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILAnimate("r", "0", fmtF(node.Size/2), core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</circle>`)
	}
	return b.String()
}

// renderMeshLayer emits the accurate voronoi-mesh hover overlay over the node
// centers (charts/interact, backed by internal/d3/delaunay). Active only when
// Interactive && UseMesh; DebugMesh draws the cells; DetectionRadius bounds
// hit-testing when > 0. props.Width/Height are the inner dimensions (set by the
// Network templ before renderLayers runs).
func renderMeshLayer(props NetworkProps, result NetworkResult) string {
	if !props.Interactive || !props.UseMesh {
		return ""
	}
	pts := make([]interact.MeshPoint, 0, len(result.Nodes))
	for _, n := range result.Nodes {
		pts = append(pts, interact.MeshPoint{
			X:    n.X,
			Y:    n.Y,
			HTML: interact.TooltipHTML(n.Color, n.ID, ""),
		})
	}
	return interact.MeshOverlay(pts, props.Width, props.Height, props.DebugMesh, props.DetectionRadius)
}

func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
