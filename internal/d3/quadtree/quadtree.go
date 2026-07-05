// Package quadtree is a port of d3-quadtree: a two-dimensional recursive spatial
// subdivision used to accelerate spatial queries. templ-charts needs it as the
// shared index behind the Barnes–Hut many-body approximation and the pruned
// collision force in internal/d3/force — the O(n log n) replacements for the
// exact O(n²) charge/collide loops.
//
// The structure mirrors d3-quadtree exactly so its output is bit-comparable:
// each node is either an internal node with up to four children (indexed
// bottom<<1|right, i.e. 0=top-left, 1=top-right, 2=bottom-left, 3=bottom-right)
// or a leaf holding one datum plus a linked list (Next) of any datums that share
// its exact coordinates. The tree covers a square extent that grows by doubling
// (Cover) as points outside it are added.
//
// Barnes–Hut consumers walk the tree post-order (VisitAfter) and accumulate a
// charge/centre-of-mass on each internal node; the Value/X/Y fields on Node are
// scratch space reserved for exactly that, matching how d3-force annotates the
// quadtree.
package quadtree

import "math"

// Node is a quadtree node: either internal (children set, IsLeaf false) or a
// leaf (Data set, IsLeaf true, with Next chaining coincident datums). Value/X/Y
// are scratch fields reserved for consumers that aggregate over the tree (e.g.
// Barnes–Hut charge + centre of mass); the quadtree itself never reads them.
type Node struct {
	children [4]*Node
	leaf     bool
	data     any
	next     *Node

	// Value/X/Y are consumer scratch (e.g. Barnes–Hut). Unused by the quadtree.
	Value float64
	X, Y  float64
}

// IsLeaf reports whether n is a leaf (holds a datum) rather than an internal
// node (holds children).
func (n *Node) IsLeaf() bool { return n.leaf }

// Data returns the leaf's datum (nil for an internal node).
func (n *Node) Data() any { return n.data }

// Next returns the next coincident-point leaf sharing this leaf's coordinates,
// or nil. Only meaningful for leaves.
func (n *Node) Next() *Node { return n.next }

// Child returns the i-th child (0..3) of an internal node, or nil.
func (n *Node) Child(i int) *Node { return n.children[i] }

// Quadtree is a d3-quadtree over data of any type, positioned by the x/y
// accessors supplied at construction. The zero value is not usable; build one
// with New or From.
type Quadtree struct {
	x0, y0, x1, y1 float64 // covered extent; NaN until the first point is covered
	root           *Node
	xAcc           func(any) float64
	yAcc           func(any) float64
}

// New returns an empty quadtree that positions each datum via the x and y
// accessors. The extent is uncovered until the first Add/Cover.
func New(x, y func(any) float64) *Quadtree {
	nan := math.NaN()
	return &Quadtree{x0: nan, y0: nan, x1: nan, y1: nan, xAcc: x, yAcc: y}
}

// From builds a quadtree over data with the given accessors (New + AddAll).
func From(data []any, x, y func(any) float64) *Quadtree {
	q := New(x, y)
	q.AddAll(data)
	return q
}

// Root returns the tree's root node (nil when empty).
func (q *Quadtree) Root() *Node { return q.root }

// Extent returns the covered square extent and whether the tree has been
// covered (false when still empty/uncovered).
func (q *Quadtree) Extent() (x0, y0, x1, y1 float64, ok bool) {
	return q.x0, q.y0, q.x1, q.y1, !math.IsNaN(q.x0)
}

func mid(a, b float64) float64 { return (a + b) / 2 }

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Cover expands the extent (by repeatedly doubling toward the point) so it
// contains (x, y), re-rooting the tree as it grows. Matches d3's cover. NaN
// coordinates are ignored.
func (q *Quadtree) Cover(x, y float64) *Quadtree {
	if math.IsNaN(x) || math.IsNaN(y) {
		return q
	}
	x0, y0, x1, y1 := q.x0, q.y0, q.x1, q.y1
	if math.IsNaN(x0) {
		// Uncovered: seed a unit square around the floored point.
		x0 = math.Floor(x)
		x1 = x0 + 1
		y0 = math.Floor(y)
		y1 = y0 + 1
	} else {
		// Grow toward the point, doubling the extent each step and wrapping the
		// current root as a child of a new internal node.
		z := x1 - x0
		if z == 0 {
			z = 1
		}
		node := q.root
		for x0 > x || x >= x1 || y0 > y || y >= y1 {
			i := b2i(y < y0)<<1 | b2i(x < x0)
			parent := &Node{}
			parent.children[i] = node
			node = parent
			z *= 2
			switch i {
			case 0:
				x1 = x0 + z
				y1 = y0 + z
			case 1:
				x0 = x1 - z
				y1 = y0 + z
			case 2:
				x1 = x0 + z
				y0 = y1 - z
			case 3:
				x0 = x1 - z
				y0 = y1 - z
			}
		}
		// Only re-root when the existing root is an internal node: a lone leaf
		// carries its position in its datum, not its quadrant, so it need not move.
		if q.root != nil && !q.root.leaf {
			q.root = node
		}
	}
	q.x0, q.y0, q.x1, q.y1 = x0, y0, x1, y1
	return q
}

// Add inserts one datum, covering its coordinates first. NaN coordinates are
// ignored. Matches d3's add (including coincident-point chaining and the
// subdivide-until-separate leaf split).
func (q *Quadtree) Add(d any) *Quadtree {
	x, y := q.xAcc(d), q.yAcc(d)
	q.Cover(x, y)
	return q.add(x, y, d)
}

func (q *Quadtree) add(x, y float64, d any) *Quadtree {
	if math.IsNaN(x) || math.IsNaN(y) {
		return q
	}
	leaf := &Node{leaf: true, data: d}
	if q.root == nil {
		q.root = leaf
		return q
	}
	x0, y0, x1, y1 := q.x0, q.y0, q.x1, q.y1
	var parent *Node
	node := q.root
	var i int
	// Descend into internal nodes toward (x, y).
	for !node.leaf {
		xm := mid(x0, x1)
		right := x >= xm
		if right {
			x0 = xm
		} else {
			x1 = xm
		}
		ym := mid(y0, y1)
		bottom := y >= ym
		if bottom {
			y0 = ym
		} else {
			y1 = ym
		}
		parent = node
		i = b2i(bottom)<<1 | b2i(right)
		node = node.children[i]
		if node == nil {
			parent.children[i] = leaf
			return q
		}
	}
	// node is the leaf occupying this cell.
	xp, yp := q.xAcc(node.data), q.yAcc(node.data)
	if x == xp && y == yp {
		// Coincident: chain the new leaf ahead of the existing one.
		leaf.next = node
		if parent != nil {
			parent.children[i] = leaf
		} else {
			q.root = leaf
		}
		return q
	}
	// Subdivide until the new point and the existing leaf fall in distinct quads.
	for {
		if parent != nil {
			n := &Node{}
			parent.children[i] = n
			parent = n
		} else {
			n := &Node{}
			q.root = n
			parent = n
		}
		xm := mid(x0, x1)
		right := x >= xm
		if right {
			x0 = xm
		} else {
			x1 = xm
		}
		ym := mid(y0, y1)
		bottom := y >= ym
		if bottom {
			y0 = ym
		} else {
			y1 = ym
		}
		i = b2i(bottom)<<1 | b2i(right)
		j := b2i(yp >= ym)<<1 | b2i(xp >= xm)
		if i != j {
			parent.children[j] = node
			parent.children[i] = leaf
			return q
		}
	}
}

// AddAll inserts many datums. It pre-covers the bounding box of the batch (as
// d3 does) so the tree is not repeatedly re-rooted mid-insert.
func (q *Quadtree) AddAll(data []any) *Quadtree {
	n := len(data)
	xs := make([]float64, n)
	ys := make([]float64, n)
	x0, y0 := math.Inf(1), math.Inf(1)
	x1, y1 := math.Inf(-1), math.Inf(-1)
	for i, d := range data {
		x, y := q.xAcc(d), q.yAcc(d)
		if math.IsNaN(x) || math.IsNaN(y) {
			xs[i] = math.NaN()
			ys[i] = math.NaN()
			continue
		}
		xs[i], ys[i] = x, y
		if x < x0 {
			x0 = x
		}
		if x > x1 {
			x1 = x
		}
		if y < y0 {
			y0 = y
		}
		if y > y1 {
			y1 = y
		}
	}
	if x0 > x1 || y0 > y1 {
		return q // all NaN / empty batch
	}
	q.Cover(x0, y0).Cover(x1, y1)
	for i, d := range data {
		q.add(xs[i], ys[i], d)
	}
	return q
}

// Remove deletes one datum (by value equality of the datum). Internal nodes are
// collapsed when a removal leaves a single leaf child, matching d3's remove.
func (q *Quadtree) Remove(d any) *Quadtree {
	node := q.root
	if node == nil {
		return q
	}
	x, y := q.xAcc(d), q.yAcc(d)
	x0, y0, x1, y1 := q.x0, q.y0, q.x1, q.y1
	var parent, retainer, previous *Node
	var i, j int
	if !node.leaf {
		for {
			xm := mid(x0, x1)
			right := x >= xm
			if right {
				x0 = xm
			} else {
				x1 = xm
			}
			ym := mid(y0, y1)
			bottom := y >= ym
			if bottom {
				y0 = ym
			} else {
				y1 = ym
			}
			parent = node
			i = b2i(bottom)<<1 | b2i(right)
			node = node.children[i]
			if node == nil {
				return q
			}
			if node.leaf {
				break
			}
			// Remember the deepest ancestor with more than one child; a collapse
			// promotes the surviving leaf up to it.
			if parent.children[(i+1)&3] != nil || parent.children[(i+2)&3] != nil || parent.children[(i+3)&3] != nil {
				retainer = parent
				j = i
			}
		}
	}
	// Walk the coincident chain to the exact datum.
	for node.data != d {
		previous = node
		node = node.next
		if node == nil {
			return q
		}
	}
	next := node.next
	node.next = nil
	// Removed from the middle/tail of a coincident chain: just relink.
	if previous != nil {
		previous.next = next
		return q
	}
	// Removed the head of a root leaf.
	if parent == nil {
		q.root = next
		return q
	}
	// Replace the cell with the chain remainder, or clear it.
	parent.children[i] = next
	// Collapse: if the parent now has exactly one child and it's a leaf, promote
	// it to the retainer (or root).
	only := firstChild(parent)
	if only != nil && only == lastChild(parent) && only.leaf {
		if retainer != nil {
			retainer.children[j] = only
		} else {
			q.root = only
		}
	}
	return q
}

// RemoveAll deletes each datum in data.
func (q *Quadtree) RemoveAll(data []any) *Quadtree {
	for _, d := range data {
		q.Remove(d)
	}
	return q
}

func firstChild(n *Node) *Node {
	for i := 0; i < 4; i++ {
		if n.children[i] != nil {
			return n.children[i]
		}
	}
	return nil
}

func lastChild(n *Node) *Node {
	for i := 3; i >= 0; i-- {
		if n.children[i] != nil {
			return n.children[i]
		}
	}
	return nil
}

type quad struct {
	node           *Node
	x0, y0, x1, y1 float64
}

// Visit traverses the tree pre-order (parent before children), invoking fn with
// each node and its cell bounds. Returning true from fn skips that node's
// children. Matches d3's visit (children are pushed 3,2,1,0 so they pop 0,1,2,3).
func (q *Quadtree) Visit(fn func(n *Node, x0, y0, x1, y1 float64) bool) *Quadtree {
	if q.root == nil {
		return q
	}
	quads := []quad{{q.root, q.x0, q.y0, q.x1, q.y1}}
	for len(quads) > 0 {
		cur := quads[len(quads)-1]
		quads = quads[:len(quads)-1]
		node := cur.node
		if !fn(node, cur.x0, cur.y0, cur.x1, cur.y1) && !node.leaf {
			xm := mid(cur.x0, cur.x1)
			ym := mid(cur.y0, cur.y1)
			if c := node.children[3]; c != nil {
				quads = append(quads, quad{c, xm, ym, cur.x1, cur.y1})
			}
			if c := node.children[2]; c != nil {
				quads = append(quads, quad{c, cur.x0, ym, xm, cur.y1})
			}
			if c := node.children[1]; c != nil {
				quads = append(quads, quad{c, xm, cur.y0, cur.x1, ym})
			}
			if c := node.children[0]; c != nil {
				quads = append(quads, quad{c, cur.x0, cur.y0, xm, ym})
			}
		}
	}
	return q
}

// VisitAfter traverses the tree post-order (children before parent), invoking fn
// with each node and its cell bounds. This is the order Barnes–Hut aggregation
// needs (a parent's charge is the sum of its children's). Matches d3's visitAfter.
func (q *Quadtree) VisitAfter(fn func(n *Node, x0, y0, x1, y1 float64)) *Quadtree {
	if q.root == nil {
		return q
	}
	quads := []quad{{q.root, q.x0, q.y0, q.x1, q.y1}}
	var next []quad
	for len(quads) > 0 {
		cur := quads[len(quads)-1]
		quads = quads[:len(quads)-1]
		node := cur.node
		if !node.leaf {
			xm := mid(cur.x0, cur.x1)
			ym := mid(cur.y0, cur.y1)
			if c := node.children[0]; c != nil {
				quads = append(quads, quad{c, cur.x0, cur.y0, xm, ym})
			}
			if c := node.children[1]; c != nil {
				quads = append(quads, quad{c, xm, cur.y0, cur.x1, ym})
			}
			if c := node.children[2]; c != nil {
				quads = append(quads, quad{c, cur.x0, ym, xm, cur.y1})
			}
			if c := node.children[3]; c != nil {
				quads = append(quads, quad{c, xm, ym, cur.x1, cur.y1})
			}
		}
		next = append(next, cur)
	}
	for i := len(next) - 1; i >= 0; i-- {
		fn(next[i].node, next[i].x0, next[i].y0, next[i].x1, next[i].y1)
	}
	return q
}

// Find returns the datum closest to (x, y), optionally within a search radius
// (pass a single radius; omit for unbounded). The bool is false when the tree is
// empty or nothing lies within the radius. Matches d3's find (the closest
// quadrant is visited last so its bound tightens the radius earliest).
func (q *Quadtree) Find(x, y float64, radius ...float64) (any, bool) {
	var data any
	found := false
	if q.root == nil {
		return nil, false
	}
	r := math.Inf(1)
	bx0, by0, bx1, by1 := q.x0, q.y0, q.x1, q.y1
	if len(radius) > 0 && !math.IsInf(radius[0], 1) {
		rr := radius[0]
		bx0, by0 = x-rr, y-rr
		bx1, by1 = x+rr, y+rr
		r = rr * rr
	}
	quads := []quad{{q.root, q.x0, q.y0, q.x1, q.y1}}
	for len(quads) > 0 {
		cur := quads[len(quads)-1]
		quads = quads[:len(quads)-1]
		node := cur.node
		if node == nil || cur.x0 > bx1 || cur.y0 > by1 || cur.x1 < bx0 || cur.y1 < by0 {
			continue
		}
		if !node.leaf {
			xm := mid(cur.x0, cur.x1)
			ym := mid(cur.y0, cur.y1)
			quads = append(quads,
				quad{node.children[3], xm, ym, cur.x1, cur.y1},
				quad{node.children[2], cur.x0, ym, xm, cur.y1},
				quad{node.children[1], xm, cur.y0, cur.x1, ym},
				quad{node.children[0], cur.x0, cur.y0, xm, ym},
			)
			// Visit the quadrant containing (x, y) last (popped first) so its
			// leaf tightens the radius before the others are tested.
			if i := b2i(y >= ym)<<1 | b2i(x >= xm); i != 0 {
				n := len(quads)
				quads[n-1], quads[n-1-i] = quads[n-1-i], quads[n-1]
			}
		} else {
			dx := x - q.xAcc(node.data)
			dy := y - q.yAcc(node.data)
			d2 := dx*dx + dy*dy
			if d2 < r {
				r = d2
				d := math.Sqrt(d2)
				bx0, by0 = x-d, y-d
				bx1, by1 = x+d, y+d
				data = node.data
				found = true
			}
		}
	}
	return data, found
}

// Size returns the number of datums in the tree (counting coincident points).
func (q *Quadtree) Size() int {
	size := 0
	q.Visit(func(n *Node, _, _, _, _ float64) bool {
		if n.leaf {
			for node := n; node != nil; node = node.next {
				size++
			}
		}
		return false
	})
	return size
}

// Data returns every datum in the tree, in traversal order.
func (q *Quadtree) Data() []any {
	var out []any
	q.Visit(func(n *Node, _, _, _, _ float64) bool {
		if n.leaf {
			for node := n; node != nil; node = node.next {
				out = append(out, node.data)
			}
		}
		return false
	})
	return out
}

// Copy returns a deep copy of the tree: node structure and coincident chains are
// cloned (so structural mutation of one tree doesn't affect the other), while
// the datums themselves are shared by reference — matching d3's copy. Consumer
// scratch (Value/X/Y) is copied too.
func (q *Quadtree) Copy() *Quadtree {
	c := &Quadtree{x0: q.x0, y0: q.y0, x1: q.x1, y1: q.y1, xAcc: q.xAcc, yAcc: q.yAcc}
	c.root = cloneNode(q.root)
	return c
}

func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	c := &Node{leaf: n.leaf, data: n.data, Value: n.Value, X: n.X, Y: n.Y}
	if n.leaf {
		c.next = cloneNode(n.next)
	} else {
		for i := 0; i < 4; i++ {
			c.children[i] = cloneNode(n.children[i])
		}
	}
	return c
}
