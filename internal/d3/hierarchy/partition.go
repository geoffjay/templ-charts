// partition.go — port of d3-hierarchy's partition layout. Partition subdivides
// a rectangle into horizontal bands per depth (like icicle); the sunburst chart
// maps the rectangular output onto polar coordinates.
package d3hierarchy

// Partition is the partition layout generator. Mirrors d3-hierarchy partition().
type Partition struct {
	dx, dy  float64
	padding float64
	round   bool
}

// NewPartition returns a partition with size 1×1, no padding, no rounding.
func NewPartition() *Partition { return &Partition{dx: 1, dy: 1} }

// Size sets the layout size [width, height].
func (p *Partition) Size(w, h float64) *Partition { p.dx, p.dy = w, h; return p }

// Round toggles rounding of node coordinates to integers.
func (p *Partition) Round(b bool) *Partition { p.round = b; return p }

// Padding sets the padding between adjacent nodes.
func (p *Partition) Padding(v float64) *Partition { p.padding = v; return p }

// Layout runs the partition on root (which must already have values via Sum or
// Count), setting X0/Y0/X1/Y1 on every node.
func (p *Partition) Layout(root *Node) *Node {
	n := float64(root.Height + 1)
	root.X0 = p.padding
	root.Y0 = p.padding
	root.X1 = p.dx
	root.Y1 = p.dy / n
	root.EachBefore(func(node *Node) {
		if len(node.Children) > 0 {
			TreemapDice(node, node.X0, p.dy*float64(node.Depth+1)/n, node.X1, p.dy*float64(node.Depth+2)/n)
		}
		x0 := node.X0
		y0 := node.Y0
		x1 := node.X1 - p.padding
		y1 := node.Y1 - p.padding
		if x1 < x0 {
			x0 = (x0 + x1) / 2
			x1 = x0
		}
		if y1 < y0 {
			y0 = (y0 + y1) / 2
			y1 = y0
		}
		node.X0, node.Y0, node.X1, node.Y1 = x0, y0, x1, y1
	})
	if p.round {
		root.EachBefore(roundNode)
	}
	return root
}
