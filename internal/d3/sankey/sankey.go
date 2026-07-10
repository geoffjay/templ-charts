// Package sankey ports d3-sankey (v0.12): a deterministic, single-pass layout
// for flow diagrams. Given a graph of nodes and directed, weighted links it
// assigns each node a rectangle (x0,x1,y0,y1) and each link a vertical band
// (y0,y1 center offsets + width) so that flow is conserved and crossings are
// reduced by a small fixed number of relaxation iterations.
//
// The port mirrors d3-sankey's algorithm verbatim — computeNodeLinks →
// computeNodeValues → computeNodeDepths → computeNodeHeights →
// computeNodeBreadths (initialize + relax left/right + resolve collisions) →
// computeLinkBreadths — including the alignment functions (Left/Right/Center/
// Justify). It is fully deterministic (no randomness); stable sorts are used
// wherever d3 relies on a stable Array.sort so goldens are byte-stable.
//
// Charts consume this via charts/sankey, which runs the layout then renders the
// node rects and the link ribbons (a line generator with curveMonotoneX/Y over
// the ribbon outline, matching @nivo/sankey).
package sankey

import (
	"math"
	"sort"
)

// Node is one sankey node. ID is the caller-set identity used to resolve link
// endpoints and (by default) as the layout key. FixedValue, when non-nil,
// overrides the flow-derived Value. The remaining fields are populated by
// Layout: Index (input order), Depth (distance from a source), Height
// (distance from a sink), Layer (assigned column), Value, the rectangle
// X0/Y0/X1/Y1, and the incident links.
type Node struct {
	ID         string
	FixedValue *float64

	Index  int
	Depth  int
	Height int
	Layer  int
	Value  float64
	X0     float64
	Y0     float64
	X1     float64
	Y1     float64

	SourceLinks []*Link
	TargetLinks []*Link
}

// Link is one directed, weighted link. SourceID/TargetID reference Node.ID in
// the input; Layout resolves them into Source/Target pointers. Y0/Y1 are the
// link's center y at the source and target edges and Width its thickness — all
// set by Layout.
type Link struct {
	SourceID string
	TargetID string
	Value    float64

	Index  int
	Source *Node
	Target *Node
	Y0     float64
	Y1     float64
	Width  float64
}

// Graph is the mutable node/link set Layout operates on.
type Graph struct {
	Nodes []*Node
	Links []*Link
}

// AlignFunc returns a node's target column index given the total column count.
// Mirrors d3-sankey's nodeAlign functions.
type AlignFunc func(node *Node, n int) float64

// SankeyLeft aligns every node to its depth (flow distance from a source).
func SankeyLeft(node *Node, _ int) float64 { return float64(node.Depth) }

// SankeyRight aligns every node against the right edge by its height.
func SankeyRight(node *Node, n int) float64 { return float64(n - 1 - node.Height) }

// SankeyJustify keeps sources at their depth and pushes sinks to the last column.
func SankeyJustify(node *Node, n int) float64 {
	if len(node.SourceLinks) > 0 {
		return float64(node.Depth)
	}
	return float64(n - 1)
}

// SankeyCenter keeps nodes with inputs at their depth, pulls pure sources one
// column left of their nearest target, and leaves true roots at column 0.
func SankeyCenter(node *Node, _ int) float64 {
	if len(node.TargetLinks) > 0 {
		return float64(node.Depth)
	}
	if len(node.SourceLinks) > 0 {
		min := node.SourceLinks[0].Target.Depth
		for _, l := range node.SourceLinks[1:] {
			if l.Target.Depth < min {
				min = l.Target.Depth
			}
		}
		return float64(min - 1)
	}
	return 0
}

// nodeSortMode selects how columns are (re)ordered. Mirrors d3-sankey's
// tri-state nodeSort: auto (undefined → re-sort by breadth during relaxation),
// input (null → never re-sort), custom (a comparator → sort once, no re-sort).
type nodeSortMode int

const (
	nodeSortAuto nodeSortMode = iota
	nodeSortInput
	nodeSortCustom
)

// Sankey is the configured layout generator. Construct with New, configure with
// the chainable setters, then call Layout.
type Sankey struct {
	x0, y0, x1, y1 float64 // extent
	dx             float64 // node width
	dy             float64 // node padding
	py             float64 // resolved padding (computed in Layout)
	iterations     int

	align        AlignFunc
	nodeSortMode nodeSortMode
	nodeCmp      func(a, b *Node) int
	// linkSortInput mirrors d3's linkSort===null (input order, no reordering).
	// The default (false) is d3's undefined: reorder links by breadth.
	linkSortInput bool
}

// New returns a Sankey with d3-sankey's defaults: extent [[0,0],[1,1]],
// nodeWidth 24, nodePadding 8, 6 iterations, justify alignment.
func New() *Sankey {
	return &Sankey{
		x0: 0, y0: 0, x1: 1, y1: 1,
		dx:         24,
		dy:         8,
		iterations: 6,
		align:      SankeyJustify,
	}
}

// NodeWidth sets the node rectangle thickness (dx). Mirrors nodeWidth.
func (s *Sankey) NodeWidth(w float64) *Sankey { s.dx = w; return s }

// NodePadding sets the vertical gap between nodes in a column. Mirrors nodePadding.
func (s *Sankey) NodePadding(p float64) *Sankey { s.dy = p; return s }

// Size sets the extent to [[0,0],[w,h]]. Mirrors size([w,h]).
func (s *Sankey) Size(w, h float64) *Sankey {
	s.x0, s.y0, s.x1, s.y1 = 0, 0, w, h
	return s
}

// Extent sets the layout rectangle. Mirrors extent([[x0,y0],[x1,y1]]).
func (s *Sankey) Extent(x0, y0, x1, y1 float64) *Sankey {
	s.x0, s.y0, s.x1, s.y1 = x0, y0, x1, y1
	return s
}

// Iterations sets the relaxation pass count (d3 default 6).
func (s *Sankey) Iterations(n int) *Sankey { s.iterations = n; return s }

// NodeAlign sets the alignment function (default SankeyJustify).
func (s *Sankey) NodeAlign(fn AlignFunc) *Sankey { s.align = fn; return s }

// NodeSortAuto is d3's default node ordering: columns are re-sorted by vertical
// breadth during each relaxation pass.
func (s *Sankey) NodeSortAuto() *Sankey { s.nodeSortMode = nodeSortAuto; s.nodeCmp = nil; return s }

// NodeSortInput preserves the input node order within columns (no reordering).
func (s *Sankey) NodeSortInput() *Sankey { s.nodeSortMode = nodeSortInput; s.nodeCmp = nil; return s }

// NodeSort sets a custom column comparator (applied once; disables breadth
// re-sorting). cmp returns <0, 0, >0 like sort's convention.
func (s *Sankey) NodeSort(cmp func(a, b *Node) int) *Sankey {
	s.nodeSortMode = nodeSortCustom
	s.nodeCmp = cmp
	return s
}

// LinkSortInput preserves the input link order (d3's linkSort===null): links are
// never reordered. The default reorders links by breadth (d3's undefined).
func (s *Sankey) LinkSortInput(v bool) *Sankey { s.linkSortInput = v; return s }

// Layout runs the full d3-sankey pipeline over g, mutating its nodes and links
// in place.
func (s *Sankey) Layout(g *Graph) {
	s.computeNodeLinks(g)
	s.computeNodeValues(g)
	s.computeNodeDepths(g)
	s.computeNodeHeights(g)
	s.computeNodeBreadths(g)
	s.computeLinkBreadths(g)
}

func (s *Sankey) computeNodeLinks(g *Graph) {
	byID := make(map[string]*Node, len(g.Nodes))
	for i, n := range g.Nodes {
		n.Index = i
		n.SourceLinks = n.SourceLinks[:0]
		n.TargetLinks = n.TargetLinks[:0]
		byID[n.ID] = n
	}
	for i, l := range g.Links {
		l.Index = i
		l.Source = byID[l.SourceID]
		l.Target = byID[l.TargetID]
		if l.Source == nil || l.Target == nil {
			continue // skip links referencing unknown nodes
		}
		l.Source.SourceLinks = append(l.Source.SourceLinks, l)
		l.Target.TargetLinks = append(l.Target.TargetLinks, l)
	}
}

func (s *Sankey) computeNodeValues(g *Graph) {
	for _, n := range g.Nodes {
		if n.FixedValue != nil {
			n.Value = *n.FixedValue
			continue
		}
		n.Value = math.Max(sumLinks(n.SourceLinks), sumLinks(n.TargetLinks))
	}
}

func sumLinks(links []*Link) float64 {
	total := 0.0
	for _, l := range links {
		total += l.Value
	}
	return total
}

// computeNodeDepths assigns each node its BFS distance from a source. Mirrors
// d3's set-based breadth traversal along source links.
func (s *Sankey) computeNodeDepths(g *Graph) {
	n := len(g.Nodes)
	current := append([]*Node(nil), g.Nodes...)
	x := 0
	for len(current) > 0 {
		seen := make(map[*Node]bool)
		var next []*Node
		for _, node := range current {
			node.Depth = x
			for _, l := range node.SourceLinks {
				if !seen[l.Target] {
					seen[l.Target] = true
					next = append(next, l.Target)
				}
			}
		}
		x++
		if x > n {
			// More iterations than nodes means the links form a cycle. Sankey
			// diagrams are acyclic by definition, but the input is user-supplied
			// and not validated upstream. Rather than panic (which would take
			// down the caller's request), stop the traversal and lay out a
			// best-effort diagram from the depths assigned so far.
			break
		}
		current = next
	}
}

// computeNodeHeights assigns each node its BFS distance from a sink (along
// target links).
func (s *Sankey) computeNodeHeights(g *Graph) {
	n := len(g.Nodes)
	current := append([]*Node(nil), g.Nodes...)
	x := 0
	for len(current) > 0 {
		seen := make(map[*Node]bool)
		var next []*Node
		for _, node := range current {
			node.Height = x
			for _, l := range node.TargetLinks {
				if !seen[l.Source] {
					seen[l.Source] = true
					next = append(next, l.Source)
				}
			}
		}
		x++
		if x > n {
			// Cyclic input (see computeNodeDepths): stop rather than panic and
			// lay out a best-effort diagram from the heights assigned so far.
			break
		}
		current = next
	}
}

// computeNodeLayers partitions nodes into columns by alignment, sets their x
// range, and returns the columns (nil entries for empty layers, preserving
// index → layer correspondence, as d3's sparse array does).
func (s *Sankey) computeNodeLayers(g *Graph) [][]*Node {
	x := 0
	for _, n := range g.Nodes {
		if n.Depth+1 > x {
			x = n.Depth + 1
		}
	}
	kx := 0.0
	if x > 1 {
		// d3 divides by (x-1); guard the single-column case to avoid 0*Inf=NaN.
		kx = (s.x1 - s.x0 - s.dx) / float64(x-1)
	}
	columns := make([][]*Node, x)
	for _, node := range g.Nodes {
		i := int(math.Floor(s.align(node, x)))
		if i < 0 {
			i = 0
		}
		if i > x-1 {
			i = x - 1
		}
		node.Layer = i
		node.X0 = s.x0 + float64(i)*kx
		node.X1 = node.X0 + s.dx
		columns[i] = append(columns[i], node)
	}
	if s.nodeSortMode == nodeSortCustom && s.nodeCmp != nil {
		for _, col := range columns {
			sortNodes(col, s.nodeCmp)
		}
	}
	return columns
}

func (s *Sankey) computeNodeBreadths(g *Graph) {
	columns := s.computeNodeLayers(g)
	maxLen := 0
	for _, c := range columns {
		if len(c) > maxLen {
			maxLen = len(c)
		}
	}
	s.py = math.Min(s.dy, (s.y1-s.y0)/float64(maxLen-1))
	s.initializeNodeBreadths(columns)
	for i := 0; i < s.iterations; i++ {
		alpha := math.Pow(0.99, float64(i))
		beta := math.Max(1-alpha, float64(i+1)/float64(s.iterations))
		s.relaxRightToLeft(columns, alpha, beta)
		s.relaxLeftToRight(columns, alpha, beta)
	}
}

func (s *Sankey) initializeNodeBreadths(columns [][]*Node) {
	ky := math.Inf(1)
	for _, c := range columns {
		if c == nil {
			continue
		}
		sumValue := 0.0
		for _, n := range c {
			sumValue += n.Value
		}
		k := (s.y1 - s.y0 - float64(len(c)-1)*s.py) / sumValue
		if k < ky {
			ky = k
		}
	}
	for _, nodes := range columns {
		if nodes == nil {
			continue
		}
		y := s.y0
		for _, node := range nodes {
			node.Y0 = y
			node.Y1 = y + node.Value*ky
			y = node.Y1 + s.py
			for _, link := range node.SourceLinks {
				link.Width = link.Value * ky
			}
		}
		// Center the column vertically in the remaining space.
		y = (s.y1 - y + s.py) / float64(len(nodes)+1)
		for i, node := range nodes {
			node.Y0 += y * float64(i+1)
			node.Y1 += y * float64(i+1)
		}
		s.reorderLinks(nodes)
	}
}

// relaxLeftToRight repositions each node based on its incoming (target) links.
func (s *Sankey) relaxLeftToRight(columns [][]*Node, alpha, beta float64) {
	for i := 1; i < len(columns); i++ {
		column := columns[i]
		if column == nil {
			continue
		}
		for _, target := range column {
			var y, w float64
			for _, l := range target.TargetLinks {
				source := l.Source
				v := l.Value * float64(target.Layer-source.Layer)
				y += s.targetTop(source, target) * v
				w += v
			}
			if !(w > 0) {
				continue
			}
			dy := (y/w - target.Y0) * alpha
			target.Y0 += dy
			target.Y1 += dy
			s.reorderNodeLinks(target)
		}
		if s.nodeSortMode == nodeSortAuto {
			sortByBreadth(column)
		}
		s.resolveCollisions(column, beta)
	}
}

// relaxRightToLeft repositions each node based on its outgoing (source) links.
func (s *Sankey) relaxRightToLeft(columns [][]*Node, alpha, beta float64) {
	for i := len(columns) - 2; i >= 0; i-- {
		column := columns[i]
		if column == nil {
			continue
		}
		for _, source := range column {
			var y, w float64
			for _, l := range source.SourceLinks {
				target := l.Target
				v := l.Value * float64(target.Layer-source.Layer)
				y += s.sourceTop(source, target) * v
				w += v
			}
			if !(w > 0) {
				continue
			}
			dy := (y/w - source.Y0) * alpha
			source.Y0 += dy
			source.Y1 += dy
			s.reorderNodeLinks(source)
		}
		if s.nodeSortMode == nodeSortAuto {
			sortByBreadth(column)
		}
		s.resolveCollisions(column, beta)
	}
}

func (s *Sankey) resolveCollisions(nodes []*Node, alpha float64) {
	i := len(nodes) >> 1
	subject := nodes[i]
	s.resolveCollisionsBottomToTop(nodes, subject.Y0-s.py, i-1, alpha)
	s.resolveCollisionsTopToBottom(nodes, subject.Y1+s.py, i+1, alpha)
	s.resolveCollisionsBottomToTop(nodes, s.y1, len(nodes)-1, alpha)
	s.resolveCollisionsTopToBottom(nodes, s.y0, 0, alpha)
}

// resolveCollisionsTopToBottom pushes any overlapping nodes down.
func (s *Sankey) resolveCollisionsTopToBottom(nodes []*Node, y float64, i int, alpha float64) {
	for ; i < len(nodes); i++ {
		node := nodes[i]
		dy := (y - node.Y0) * alpha
		if dy > 1e-6 {
			node.Y0 += dy
			node.Y1 += dy
		}
		y = node.Y1 + s.py
	}
}

// resolveCollisionsBottomToTop pushes any overlapping nodes up.
func (s *Sankey) resolveCollisionsBottomToTop(nodes []*Node, y float64, i int, alpha float64) {
	for ; i >= 0; i-- {
		node := nodes[i]
		dy := (node.Y1 - y) * alpha
		if dy > 1e-6 {
			node.Y0 -= dy
			node.Y1 -= dy
		}
		y = node.Y0 - s.py
	}
}

func (s *Sankey) reorderNodeLinks(node *Node) {
	if s.linkSortInput {
		return
	}
	for _, l := range node.TargetLinks {
		sortByTargetBreadth(l.Source.SourceLinks)
	}
	for _, l := range node.SourceLinks {
		sortBySourceBreadth(l.Target.TargetLinks)
	}
}

func (s *Sankey) reorderLinks(nodes []*Node) {
	if s.linkSortInput {
		return
	}
	for _, node := range nodes {
		sortByTargetBreadth(node.SourceLinks)
		sortBySourceBreadth(node.TargetLinks)
	}
}

// targetTop returns the target.Y0 that would produce an ideal link from source
// to target (d3-sankey targetTop).
func (s *Sankey) targetTop(source, target *Node) float64 {
	y := source.Y0 - float64(len(source.SourceLinks)-1)*s.py/2
	for _, l := range source.SourceLinks {
		if l.Target == target {
			break
		}
		y += l.Width + s.py
	}
	for _, l := range target.TargetLinks {
		if l.Source == source {
			break
		}
		y -= l.Width
	}
	return y
}

// sourceTop returns the source.Y0 that would produce an ideal link from source
// to target (d3-sankey sourceTop).
func (s *Sankey) sourceTop(source, target *Node) float64 {
	y := target.Y0 - float64(len(target.TargetLinks)-1)*s.py/2
	for _, l := range target.TargetLinks {
		if l.Source == source {
			break
		}
		y += l.Width + s.py
	}
	for _, l := range source.SourceLinks {
		if l.Target == target {
			break
		}
		y -= l.Width
	}
	return y
}

func (s *Sankey) computeLinkBreadths(g *Graph) {
	for _, node := range g.Nodes {
		y0 := node.Y0
		y1 := node.Y0
		for _, link := range node.SourceLinks {
			link.Y0 = y0 + link.Width/2
			y0 += link.Width
		}
		for _, link := range node.TargetLinks {
			link.Y1 = y1 + link.Width/2
			y1 += link.Width
		}
	}
}

// --- sort helpers (stable, matching d3's stable Array.sort) ---

func sortNodes(nodes []*Node, cmp func(a, b *Node) int) {
	sort.SliceStable(nodes, func(i, j int) bool { return cmp(nodes[i], nodes[j]) < 0 })
}

// sortByBreadth sorts a column by ascending Y0 (ascendingBreadth). Stable to
// match d3's stable sort tie-breaking.
func sortByBreadth(nodes []*Node) {
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].Y0 < nodes[j].Y0 })
}

// sortByTargetBreadth mirrors ascendingTargetBreadth: by target Y0 then index.
func sortByTargetBreadth(links []*Link) {
	sort.SliceStable(links, func(i, j int) bool {
		if links[i].Target.Y0 != links[j].Target.Y0 {
			return links[i].Target.Y0 < links[j].Target.Y0
		}
		return links[i].Index < links[j].Index
	})
}

// sortBySourceBreadth mirrors ascendingSourceBreadth: by source Y0 then index.
func sortBySourceBreadth(links []*Link) {
	sort.SliceStable(links, func(i, j int) bool {
		if links[i].Source.Y0 != links[j].Source.Y0 {
			return links[i].Source.Y0 < links[j].Source.Y0
		}
		return links[i].Index < links[j].Index
	})
}
