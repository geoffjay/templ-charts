// tree.go — port of d3-hierarchy's tidy-tree (Reingold–Tilford / Buchheim et
// al.) and cluster (dendrogram) layouts. Both set X and Y on every node in a
// unit-ish space scaled by Size; charts rescale/orient as needed.
package d3hierarchy

import "math"

// Separation returns the desired spacing between two adjacent nodes, in units
// of node breadth. The default is 1 for siblings, 2 otherwise.
type Separation func(a, b *Node) float64

func defaultSeparation(a, b *Node) float64 {
	if a.Parent == b.Parent {
		return 1
	}
	return 2
}

// --- Cluster (dendrogram) -------------------------------------------------

// Cluster is the cluster/dendrogram layout: all leaves at the same depth.
// Mirrors d3-hierarchy cluster().
type Cluster struct {
	separation Separation
	dx, dy     float64
	nodeSize   bool
}

// NewCluster returns a cluster layout with size 1×1 and default separation.
func NewCluster() *Cluster { return &Cluster{separation: defaultSeparation, dx: 1, dy: 1} }

// Size sets the layout size [dx, dy].
func (c *Cluster) Size(dx, dy float64) *Cluster { c.dx, c.dy = dx, dy; c.nodeSize = false; return c }

// Separation sets the sibling separation function.
func (c *Cluster) Separation(s Separation) *Cluster {
	if s != nil {
		c.separation = s
	}
	return c
}

// Layout runs the cluster layout, setting X and Y on every node.
func (c *Cluster) Layout(root *Node) *Node {
	var previousNode *Node
	x := 0.0

	root.EachAfter(func(node *Node) {
		children := node.Children
		if len(children) > 0 {
			// mean of children x; max child y + 1
			sum := 0.0
			maxY := 0.0
			for _, ch := range children {
				sum += ch.X
				maxY = math.Max(maxY, ch.Y)
			}
			node.X = sum / float64(len(children))
			node.Y = 1 + maxY
		} else {
			if previousNode != nil {
				x += c.separation(node, previousNode)
			} else {
				x = 0
			}
			node.X = x
			node.Y = 0
			previousNode = node
		}
	})

	left := leafLeft(root)
	right := leafRight(root)
	x0 := left.X - c.separation(left, right)/2
	x1 := right.X + c.separation(right, left)/2

	return root.EachAfter(func(node *Node) {
		node.X = (node.X - x0) / (x1 - x0) * c.dx
		var yn float64
		if root.Y != 0 {
			yn = node.Y / root.Y
		} else {
			yn = 1
		}
		node.Y = (1 - yn) * c.dy
	})
}

func leafLeft(node *Node) *Node {
	for len(node.Children) > 0 {
		node = node.Children[0]
	}
	return node
}

func leafRight(node *Node) *Node {
	for len(node.Children) > 0 {
		node = node.Children[len(node.Children)-1]
	}
	return node
}

// --- Tree (tidy tree) -----------------------------------------------------

// Tree is the Reingold–Tilford tidy-tree layout. Mirrors d3-hierarchy tree().
type Tree struct {
	separation Separation
	dx, dy     float64
}

// NewTree returns a tidy-tree layout with size 1×1 and default separation.
func NewTree() *Tree { return &Tree{separation: defaultSeparation, dx: 1, dy: 1} }

// Size sets the layout size [dx, dy].
func (t *Tree) Size(dx, dy float64) *Tree { t.dx, t.dy = dx, dy; return t }

// Separation sets the sibling separation function.
func (t *Tree) Separation(s Separation) *Tree {
	if s != nil {
		t.separation = s
	}
	return t
}

type treeNode struct {
	node     *Node
	parent   *treeNode
	children []*treeNode
	A        *treeNode // default ancestor
	a        *treeNode // ancestor
	z        float64   // prelim
	m        float64   // mod
	c        float64   // change
	s        float64   // shift
	t        *treeNode // thread
	i        int       // index among siblings
}

func newTreeNode(n *Node, i int) *treeNode {
	tn := &treeNode{node: n, i: i}
	tn.a = tn
	return tn
}

func treeRootOf(root *Node) *treeNode {
	t := newTreeNode(root, 0)
	stack := []*treeNode{t}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		ch := node.node.Children
		if n := len(ch); n > 0 {
			node.children = make([]*treeNode, n)
			for i := n - 1; i >= 0; i-- {
				child := newTreeNode(ch[i], i)
				child.parent = node
				node.children[i] = child
				stack = append(stack, child)
			}
		}
	}
	p := newTreeNode(nil, 0)
	p.children = []*treeNode{t}
	t.parent = p
	return t
}

func (t *treeNode) eachAfter(fn func(*treeNode)) {
	stack := []*treeNode{t}
	var next []*treeNode
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		next = append(next, n)
		stack = append(stack, n.children...)
	}
	for i := len(next) - 1; i >= 0; i-- {
		fn(next[i])
	}
}

func (t *treeNode) eachBefore(fn func(*treeNode)) {
	stack := []*treeNode{t}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fn(n)
		for i := len(n.children) - 1; i >= 0; i-- {
			stack = append(stack, n.children[i])
		}
	}
}

func nextLeft(v *treeNode) *treeNode {
	if len(v.children) > 0 {
		return v.children[0]
	}
	return v.t
}

func nextRight(v *treeNode) *treeNode {
	if len(v.children) > 0 {
		return v.children[len(v.children)-1]
	}
	return v.t
}

func moveSubtree(wm, wp *treeNode, shift float64) {
	change := shift / float64(wp.i-wm.i)
	wp.c -= change
	wp.s += shift
	wm.c += change
	wp.z += shift
	wp.m += shift
}

func executeShifts(v *treeNode) {
	shift := 0.0
	change := 0.0
	children := v.children
	for i := len(children) - 1; i >= 0; i-- {
		w := children[i]
		w.z += shift
		w.m += shift
		change += w.c
		shift += w.s + change
	}
}

func nextAncestor(vim, v, ancestor *treeNode) *treeNode {
	if vim.a.parent == v.parent {
		return vim.a
	}
	return ancestor
}

// Layout runs the tidy-tree layout, setting X and Y on every node.
func (t *Tree) Layout(root *Node) *Node {
	tr := treeRootOf(root)

	firstWalk := func(v *treeNode) {
		children := v.children
		siblings := v.parent.children
		var w *treeNode
		if v.i > 0 {
			w = siblings[v.i-1]
		}
		if len(children) > 0 {
			executeShifts(v)
			midpoint := (children[0].z + children[len(children)-1].z) / 2
			if w != nil {
				v.z = w.z + t.separation(v.node, w.node)
				v.m = v.z - midpoint
			} else {
				v.z = midpoint
			}
		} else if w != nil {
			v.z = w.z + t.separation(v.node, w.node)
		}
		def := v.parent.A
		if def == nil {
			def = siblings[0]
		}
		v.parent.A = t.apportion(v, w, def)
	}

	tr.eachAfter(firstWalk)
	tr.parent.m = -tr.z
	tr.eachBefore(func(v *treeNode) {
		v.node.X = v.z + v.parent.m
		v.m += v.parent.m
	})

	// Normalize to [dx, dy] by depth (nodeSize disabled — nivo uses size()).
	left, right, bottom := root, root, root
	root.EachBefore(func(node *Node) {
		if node.X < left.X {
			left = node
		}
		if node.X > right.X {
			right = node
		}
		if node.Depth > bottom.Depth {
			bottom = node
		}
	})
	s := 1.0
	if left != right {
		s = t.separation(left, right) / 2
	}
	tx := s - left.X
	kx := t.dx / (right.X + s + tx)
	ky := t.dy
	if bottom.Depth != 0 {
		ky = t.dy / float64(bottom.Depth)
	}
	root.EachBefore(func(node *Node) {
		node.X = (node.X + tx) * kx
		node.Y = float64(node.Depth) * ky
	})
	return root
}

func (t *Tree) apportion(v, w, ancestor *treeNode) *treeNode {
	if w == nil {
		return ancestor
	}
	vip := v
	vop := v
	vim := w
	vom := vip.parent.children[0]
	sip := vip.m
	sop := vop.m
	sim := vim.m
	som := vom.m
	for {
		vim = nextRight(vim)
		vip = nextLeft(vip)
		if vim == nil || vip == nil {
			break
		}
		vom = nextLeft(vom)
		vop = nextRight(vop)
		vop.a = v
		shift := vim.z + sim - vip.z - sip + t.separation(vim.node, vip.node)
		if shift > 0 {
			moveSubtree(nextAncestor(vim, v, ancestor), v, shift)
			sip += shift
			sop += shift
		}
		sim += vim.m
		sip += vip.m
		som += vom.m
		sop += vop.m
	}
	if vim != nil && nextRight(vop) == nil {
		vop.t = vim
		vop.m += sim - sop
	}
	if vip != nil && nextLeft(vom) == nil {
		vom.t = vip
		vom.m += sip - som
		ancestor = v
	}
	return ancestor
}
