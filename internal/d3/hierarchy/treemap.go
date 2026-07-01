// treemap.go — port of d3-hierarchy's treemap layout and its tiling
// algorithms (squarify, binary, dice, slice, sliceDice). Treemap recursively
// subdivides a rectangle so each leaf's area is proportional to its value.
package d3hierarchy

import "math"

// Phi is the golden ratio, the default squarify target aspect ratio.
var Phi = (1 + math.Sqrt(5)) / 2

// TileFunc positions a node's children within the rectangle (x0,y0,x1,y1),
// setting their X0/Y0/X1/Y1. Mirrors d3-hierarchy's tiling function contract.
type TileFunc func(parent *Node, x0, y0, x1, y1 float64)

// Treemap is the treemap layout generator. Configure via the setters, then call
// Layout(root) after root.Sum()/Sort(). Mirrors d3-hierarchy treemap().
type Treemap struct {
	tile          TileFunc
	round         bool
	dx, dy        float64
	paddingStack  []float64
	paddingInner  func(*Node) float64
	paddingTop    func(*Node) float64
	paddingRight  func(*Node) float64
	paddingBottom func(*Node) float64
	paddingLeft   func(*Node) float64
}

func constPad(v float64) func(*Node) float64 { return func(*Node) float64 { return v } }

// NewTreemap returns a treemap with squarify tiling, size 1×1, no padding, no
// rounding — matching d3-hierarchy defaults.
func NewTreemap() *Treemap {
	return &Treemap{
		tile:          TreemapSquarify,
		dx:            1,
		dy:            1,
		paddingStack:  []float64{0},
		paddingInner:  constPad(0),
		paddingTop:    constPad(0),
		paddingRight:  constPad(0),
		paddingBottom: constPad(0),
		paddingLeft:   constPad(0),
	}
}

// Size sets the layout size [width, height].
func (t *Treemap) Size(w, h float64) *Treemap { t.dx, t.dy = w, h; return t }

// Round toggles rounding of node coordinates to integers.
func (t *Treemap) Round(b bool) *Treemap { t.round = b; return t }

// Tile sets the tiling function (TreemapSquarify, TreemapBinary, …).
func (t *Treemap) Tile(fn TileFunc) *Treemap {
	if fn != nil {
		t.tile = fn
	}
	return t
}

// PaddingInner sets the padding between a node's children.
func (t *Treemap) PaddingInner(v float64) *Treemap { t.paddingInner = constPad(v); return t }

// PaddingOuter sets the padding around a node (all four sides).
func (t *Treemap) PaddingOuter(v float64) *Treemap {
	t.paddingTop, t.paddingRight, t.paddingBottom, t.paddingLeft = constPad(v), constPad(v), constPad(v), constPad(v)
	return t
}

// PaddingTop/Right/Bottom/Left set per-side outer padding.
func (t *Treemap) PaddingTop(v float64) *Treemap    { t.paddingTop = constPad(v); return t }
func (t *Treemap) PaddingRight(v float64) *Treemap  { t.paddingRight = constPad(v); return t }
func (t *Treemap) PaddingBottom(v float64) *Treemap { t.paddingBottom = constPad(v); return t }
func (t *Treemap) PaddingLeft(v float64) *Treemap   { t.paddingLeft = constPad(v); return t }

// Layout runs the treemap on root (which must already have values via Sum or
// Count), setting X0/Y0/X1/Y1 on every node.
func (t *Treemap) Layout(root *Node) *Node {
	root.X0, root.Y0, root.X1, root.Y1 = 0, 0, t.dx, t.dy
	t.paddingStack = []float64{0}
	root.EachBefore(t.positionNode)
	t.paddingStack = []float64{0}
	if t.round {
		root.EachBefore(roundNode)
	}
	return root
}

func (t *Treemap) getPad(depth int) float64 {
	if depth < len(t.paddingStack) {
		return t.paddingStack[depth]
	}
	return 0
}

func (t *Treemap) setPad(depth int, v float64) {
	for len(t.paddingStack) <= depth {
		t.paddingStack = append(t.paddingStack, 0)
	}
	t.paddingStack[depth] = v
}

func (t *Treemap) positionNode(node *Node) {
	p := t.getPad(node.Depth)
	x0 := node.X0 + p
	y0 := node.Y0 + p
	x1 := node.X1 - p
	y1 := node.Y1 - p
	if x1 < x0 {
		x0 = (x0 + x1) / 2
		x1 = x0
	}
	if y1 < y0 {
		y0 = (y0 + y1) / 2
		y1 = y0
	}
	node.X0, node.Y0, node.X1, node.Y1 = x0, y0, x1, y1
	if len(node.Children) > 0 {
		p = t.paddingInner(node) / 2
		t.setPad(node.Depth+1, p)
		x0 += t.paddingLeft(node) - p
		y0 += t.paddingTop(node) - p
		x1 -= t.paddingRight(node) - p
		y1 -= t.paddingBottom(node) - p
		if x1 < x0 {
			x0 = (x0 + x1) / 2
			x1 = x0
		}
		if y1 < y0 {
			y0 = (y0 + y1) / 2
			y1 = y0
		}
		t.tile(node, x0, y0, x1, y1)
	}
}

func roundNode(node *Node) {
	node.X0 = math.Round(node.X0)
	node.Y0 = math.Round(node.Y0)
	node.X1 = math.Round(node.X1)
	node.Y1 = math.Round(node.Y1)
}

// --- tiling algorithms ----------------------------------------------------

// TreemapDice lays a node's children out in a single row (left→right), each
// width ∝ value. Mirrors d3-hierarchy treemapDice.
func TreemapDice(parent *Node, x0, y0, x1, y1 float64) {
	nodes := parent.Children
	var k float64
	if parent.Value != 0 {
		k = (x1 - x0) / parent.Value
	}
	x := x0
	for _, node := range nodes {
		node.Y0, node.Y1 = y0, y1
		node.X0 = x
		x += node.Value * k
		node.X1 = x
	}
}

// TreemapSlice lays a node's children out in a single column (top→bottom),
// each height ∝ value. Mirrors d3-hierarchy treemapSlice.
func TreemapSlice(parent *Node, x0, y0, x1, y1 float64) {
	nodes := parent.Children
	var k float64
	if parent.Value != 0 {
		k = (y1 - y0) / parent.Value
	}
	y := y0
	for _, node := range nodes {
		node.X0, node.X1 = x0, x1
		node.Y0 = y
		y += node.Value * k
		node.Y1 = y
	}
}

// TreemapSliceDice alternates slice (odd depth) and dice (even depth). Mirrors
// d3-hierarchy treemapSliceDice.
func TreemapSliceDice(parent *Node, x0, y0, x1, y1 float64) {
	if parent.Depth&1 == 1 {
		TreemapSlice(parent, x0, y0, x1, y1)
	} else {
		TreemapDice(parent, x0, y0, x1, y1)
	}
}

// TreemapSquarify is the default tiling: it packs children into rows whose
// aspect ratios approach the golden ratio. Mirrors d3-hierarchy treemapSquarify.
func TreemapSquarify(parent *Node, x0, y0, x1, y1 float64) {
	squarifyRatio(Phi, parent, x0, y0, x1, y1)
}

// SquarifyRatio returns a TileFunc squarifying to the given target ratio (>1).
func SquarifyRatio(ratio float64) TileFunc {
	if ratio < 1 {
		ratio = 1
	}
	return func(parent *Node, x0, y0, x1, y1 float64) {
		squarifyRatio(ratio, parent, x0, y0, x1, y1)
	}
}

func squarifyRatio(ratio float64, parent *Node, x0, y0, x1, y1 float64) {
	nodes := parent.Children
	n := len(nodes)
	value := parent.Value
	i0 := 0
	i1 := 0
	for i0 < n {
		dx := x1 - x0
		dy := y1 - y0

		// Find the next non-empty node.
		var sumValue float64
		for {
			sumValue = nodes[i1].Value
			i1++
			if sumValue != 0 || i1 >= n {
				break
			}
		}
		minValue := sumValue
		maxValue := sumValue
		alpha := math.Max(dy/dx, dx/dy) / (value * ratio)
		beta := sumValue * sumValue * alpha
		minRatio := math.Max(maxValue/beta, beta/minValue)

		// Keep adding nodes while the aspect ratio maintains or improves.
		for ; i1 < n; i1++ {
			nodeValue := nodes[i1].Value
			sumValue += nodeValue
			if nodeValue < minValue {
				minValue = nodeValue
			}
			if nodeValue > maxValue {
				maxValue = nodeValue
			}
			beta = sumValue * sumValue * alpha
			newRatio := math.Max(maxValue/beta, beta/minValue)
			if newRatio > minRatio {
				sumValue -= nodeValue
				break
			}
			minRatio = newRatio
		}

		// Position and record the row orientation.
		row := &Node{Value: sumValue, Children: nodes[i0:i1]}
		if dx < dy {
			// dice
			var yk float64
			if value != 0 {
				yk = y0 + dy*sumValue/value
			} else {
				yk = y1
			}
			TreemapDice(row, x0, y0, x1, yk)
			y0 = yk
		} else {
			// slice
			var xk float64
			if value != 0 {
				xk = x0 + dx*sumValue/value
			} else {
				xk = x1
			}
			TreemapSlice(row, x0, y0, xk, y1)
			x0 = xk
		}
		value -= sumValue
		i0 = i1
	}
}

// TreemapBinary splits children by value using a balanced binary partition.
// Mirrors d3-hierarchy treemapBinary.
func TreemapBinary(parent *Node, x0, y0, x1, y1 float64) {
	nodes := parent.Children
	n := len(nodes)
	sums := make([]float64, n+1)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += nodes[i].Value
		sums[i+1] = sum
	}

	var partition func(i, j int, value, x0, y0, x1, y1 float64)
	partition = func(i, j int, value, x0, y0, x1, y1 float64) {
		if i >= j-1 {
			node := nodes[i]
			node.X0, node.Y0, node.X1, node.Y1 = x0, y0, x1, y1
			return
		}
		valueOffset := sums[i]
		valueTarget := value/2 + valueOffset
		k := i + 1
		hi := j - 1
		for k < hi {
			mid := (k + hi) >> 1
			if sums[mid] < valueTarget {
				k = mid + 1
			} else {
				hi = mid
			}
		}
		if (valueTarget-sums[k-1]) < (sums[k]-valueTarget) && i+1 < k {
			k--
		}
		valueLeft := sums[k] - valueOffset
		valueRight := value - valueLeft
		if (x1 - x0) > (y1 - y0) {
			xk := x1
			if value != 0 {
				xk = (x0*valueRight + x1*valueLeft) / value
			}
			partition(i, k, valueLeft, x0, y0, xk, y1)
			partition(k, j, valueRight, xk, y0, x1, y1)
		} else {
			yk := y1
			if value != 0 {
				yk = (y0*valueRight + y1*valueLeft) / value
			}
			partition(i, k, valueLeft, x0, y0, x1, yk)
			partition(k, j, valueRight, x0, yk, x1, y1)
		}
	}
	partition(0, n, parent.Value, x0, y0, x1, y1)
}
