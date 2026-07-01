// pack.go — port of d3-hierarchy's circle-packing layout (pack + packSiblings
// front-chain). Enclosing-circle math lives in enclose.go; the shared LCG (from
// lcg.go) makes packing deterministic. Pack sets X, Y, R on every node.
package d3hierarchy

import "math"

// Pack is the circle-packing layout generator. Mirrors d3-hierarchy pack().
type Pack struct {
	dx, dy  float64
	padding float64
}

// NewPack returns a pack layout with size 1×1 and no padding.
func NewPack() *Pack { return &Pack{dx: 1, dy: 1} }

// Size sets the layout size [width, height].
func (p *Pack) Size(w, h float64) *Pack { p.dx, p.dy = w, h; return p }

// Padding sets the padding between circles.
func (p *Pack) Padding(v float64) *Pack { p.padding = v; return p }

// Layout runs circle packing on root (which must already have values via Sum),
// setting X, Y and R on every node. Leaf radius defaults to sqrt(value).
func (p *Pack) Layout(root *Node) *Node {
	random := newLCG().next
	root.X = p.dx / 2
	root.Y = p.dy / 2
	minSize := math.Min(p.dx, p.dy)

	root.EachBefore(radiusLeaf(defaultRadius))
	root.EachAfter(packChildren(0, 1, random))
	root.EachAfter(packChildren(p.padding, root.R/minSize, random))
	root.EachBefore(translateChild(minSize / (2 * root.R)))
	return root
}

func defaultRadius(node *Node) float64 { return math.Sqrt(node.Value) }

func radiusLeaf(radius func(*Node) float64) func(*Node) {
	return func(node *Node) {
		if len(node.Children) == 0 {
			r := radius(node)
			if math.IsNaN(r) || r < 0 {
				r = 0
			}
			node.R = r
		}
	}
}

func packChildren(padding, k float64, random func() float64) func(*Node) {
	return func(node *Node) {
		children := node.Children
		if len(children) == 0 {
			return
		}
		n := len(children)
		r := padding * k
		if math.IsNaN(r) {
			r = 0
		}
		if r != 0 {
			for i := 0; i < n; i++ {
				children[i].R += r
			}
		}
		e := packSiblingsRandom(children, random)
		if r != 0 {
			for i := 0; i < n; i++ {
				children[i].R -= r
			}
		}
		node.R = e + r
	}
}

func translateChild(k float64) func(*Node) {
	return func(node *Node) {
		parent := node.Parent
		node.R *= k
		if parent != nil {
			node.X = parent.X + k*node.X
			node.Y = parent.Y + k*node.Y
		}
	}
}

// --- packSiblings front-chain ---------------------------------------------

// place positions circle c tangent to a and b (a and b already placed).
func place(b, a, c *Node) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	d2 := dx*dx + dy*dy
	if d2 != 0 {
		a2 := a.R + c.R
		a2 *= a2
		b2 := b.R + c.R
		b2 *= b2
		if a2 > b2 {
			x := (d2 + b2 - a2) / (2 * d2)
			y := math.Sqrt(math.Max(0, b2/d2-x*x))
			c.X = b.X - x*dx - y*dy
			c.Y = b.Y - x*dy + y*dx
		} else {
			x := (d2 + a2 - b2) / (2 * d2)
			y := math.Sqrt(math.Max(0, a2/d2-x*x))
			c.X = a.X + x*dx - y*dy
			c.Y = a.Y + x*dy + y*dx
		}
	} else {
		c.X = a.X + c.R
		c.Y = a.Y
	}
}

func intersects(a, b *Node) bool {
	dr := a.R + b.R - 1e-6
	dx := b.X - a.X
	dy := b.Y - a.Y
	return dr > 0 && dr*dr > dx*dx+dy*dy
}

type chainNode struct {
	c          *Node
	next, prev *chainNode
}

func chainScore(node *chainNode) float64 {
	a := node.c
	b := node.next.c
	ab := a.R + b.R
	dx := (a.X*b.R + b.X*a.R) / ab
	dy := (a.Y*b.R + b.Y*a.R) / ab
	return dx*dx + dy*dy
}

// packSiblingsRandom places the given circles (nodes with R set) so none
// overlap, then centers them, returning the enclosing radius. Mirrors
// d3-hierarchy packSiblingsRandom.
func packSiblingsRandom(circles []*Node, random func() float64) float64 {
	n := len(circles)
	if n == 0 {
		return 0
	}

	a := circles[0]
	a.X, a.Y = 0, 0
	if n <= 1 {
		return a.R
	}

	b := circles[1]
	a.X, b.X, b.Y = -b.R, a.R, 0
	if n <= 2 {
		return a.R + b.R
	}

	c := circles[2]
	place(b, a, c)

	na := &chainNode{c: a}
	nb := &chainNode{c: b}
	nc := &chainNode{c: c}
	na.next, nc.prev = nb, nb
	nb.next, na.prev = nc, nc
	nc.next, nb.prev = na, na

	i := 3
	for i < n {
		place(na.c, nb.c, circles[i])
		nc = &chainNode{c: circles[i]}

		j := nb.next
		k := na.prev
		sj := nb.c.R
		sk := na.c.R
		placed := true
		for {
			if sj <= sk {
				if intersects(j.c, nc.c) {
					nb = j
					na.next = nb
					nb.prev = na
					i--
					placed = false
					break
				}
				sj += j.c.R
				j = j.next
			} else {
				if intersects(k.c, nc.c) {
					na = k
					na.next = nb
					nb.prev = na
					i--
					placed = false
					break
				}
				sk += k.c.R
				k = k.prev
			}
			if j == k.next {
				break
			}
		}
		if !placed {
			i++
			continue
		}

		// Insert the new node nc between na and nb.
		nc.prev = na
		nc.next = nb
		na.next = nc
		nb.prev = nc

		// Recompute the closest circle pair to the centroid: scan every node
		// on the front chain (from nc.next around back to nc), keeping the one
		// with the smallest score in na.
		aa := chainScore(na)
		for cc := nc.next; cc != nc; cc = cc.next {
			if ca := chainScore(cc); ca < aa {
				na = cc
				aa = ca
			}
		}
		nb = na.next
		i++
	}

	// Enclosing circle of the front chain.
	chain := []*Node{nb.c}
	for cc := nb.next; cc != nb; cc = cc.next {
		chain = append(chain, cc.c)
	}
	enc := packEncloseRandom(chain, random)

	// Translate so the enclosing circle is centered at the origin.
	for idx := 0; idx < n; idx++ {
		circles[idx].X -= enc.x
		circles[idx].Y -= enc.y
	}
	return enc.r
}
