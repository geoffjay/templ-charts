// Package d3hierarchy is a port of d3-hierarchy: the Node tree structure and
// its traversal/aggregation helpers (sum, count, sort, each/eachBefore/
// eachAfter, descendants, leaves, ancestors, links) plus the treemap,
// partition, pack, and tree/cluster layouts. It is the geometry engine behind
// the treemap/sunburst/icicle/circle-packing/tree charts.
//
// The algorithms are faithful transliterations of d3-hierarchy v3, so node
// coordinates match d3 for the same input. All layouts are deterministic; the
// pack layout ports d3's built-in LCG (see lcg.go) instead of Math.random so
// runs are reproducible server-side.
package d3hierarchy

import (
	"math"
	"sort"
)

// Node is a node in a hierarchy. Data holds the caller's datum; the traversal
// helpers and layouts populate the remaining fields. Coordinate fields are set
// by the layouts: treemap/partition set X0,Y0,X1,Y1; pack sets X,Y,R; tree/
// cluster set X,Y.
type Node struct {
	Data     any
	Parent   *Node
	Children []*Node
	Depth    int
	Height   int
	Value    float64

	// Rectangular layouts (treemap, partition).
	X0, Y0, X1, Y1 float64
	// Circular layout (pack) and node layouts (tree/cluster).
	X, Y, R float64
}

// Link is a parent→child edge (used by tree/dendrogram link rendering).
type Link struct {
	Source *Node
	Target *Node
}

// Hierarchy constructs a Node tree from root data, using children to enumerate
// each datum's child data (return nil/empty for a leaf). It computes Depth and
// Height for every node. Mirrors d3-hierarchy's hierarchy().
func Hierarchy(data any, children func(any) []any) *Node {
	root := &Node{Data: data}
	nodes := []*Node{root}
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		childs := children(node.Data)
		if n := len(childs); n > 0 {
			node.Children = make([]*Node, n)
			for i := n - 1; i >= 0; i-- {
				child := &Node{Data: childs[i], Parent: node, Depth: node.Depth + 1}
				node.Children[i] = child
				nodes = append(nodes, child)
			}
		}
	}
	return root.EachBefore(computeHeight)
}

// computeHeight sets node.Height (the longest downward path to a leaf) by
// walking up from each pre-order-visited node. Mirrors d3-hierarchy computeHeight.
func computeHeight(node *Node) {
	height := 0
	n := node
	for {
		n.Height = height
		n = n.Parent
		height++
		if n == nil || n.Height >= height {
			break
		}
	}
}

// EachBefore invokes fn on this node then its descendants in pre-order (depth
// first). Returns the receiver for chaining. Mirrors d3-hierarchy eachBefore.
func (root *Node) EachBefore(fn func(*Node)) *Node {
	nodes := []*Node{root}
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		fn(node)
		ch := node.Children
		for i := len(ch) - 1; i >= 0; i-- {
			nodes = append(nodes, ch[i])
		}
	}
	return root
}

// EachAfter invokes fn on this node's descendants then itself in post-order.
// Returns the receiver for chaining. Mirrors d3-hierarchy eachAfter.
func (root *Node) EachAfter(fn func(*Node)) *Node {
	nodes := []*Node{root}
	var next []*Node
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		next = append(next, node)
		nodes = append(nodes, node.Children...)
	}
	for i := len(next) - 1; i >= 0; i-- {
		fn(next[i])
	}
	return root
}

// Each invokes fn on this node then its descendants in breadth-first order.
// Returns the receiver for chaining. Mirrors d3-hierarchy each.
func (root *Node) Each(fn func(*Node)) *Node {
	queue := []*Node{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fn(node)
		queue = append(queue, node.Children...)
	}
	return root
}

// Sum computes node.Value for every node as value(data) plus the summed values
// of its children (post-order). Mirrors d3-hierarchy node.sum.
func (root *Node) Sum(value func(any) float64) *Node {
	return root.EachAfter(func(node *Node) {
		sum := value(node.Data)
		if math.IsNaN(sum) {
			sum = 0
		}
		for _, c := range node.Children {
			sum += c.Value
		}
		node.Value = sum
	})
}

// Count sets node.Value to the number of leaves under each node. Mirrors
// d3-hierarchy node.count.
func (root *Node) Count() *Node {
	return root.EachAfter(func(node *Node) {
		if len(node.Children) == 0 {
			node.Value = 1
			return
		}
		sum := 0.0
		for _, c := range node.Children {
			sum += c.Value
		}
		node.Value = sum
	})
}

// Sort sorts each node's children by less (pre-order). Mirrors d3-hierarchy
// node.sort; a stable sort is used so equal siblings keep input order.
func (root *Node) Sort(less func(a, b *Node) bool) *Node {
	return root.EachBefore(func(node *Node) {
		if len(node.Children) > 1 {
			ch := node.Children
			sort.SliceStable(ch, func(i, j int) bool { return less(ch[i], ch[j]) })
		}
	})
}

// Descendants returns this node and all its descendants in breadth-first order.
func (root *Node) Descendants() []*Node {
	out := make([]*Node, 0)
	root.Each(func(n *Node) { out = append(out, n) })
	return out
}

// Leaves returns all leaf nodes (no children) under this node, in pre-order.
func (root *Node) Leaves() []*Node {
	out := make([]*Node, 0)
	root.EachBefore(func(n *Node) {
		if len(n.Children) == 0 {
			out = append(out, n)
		}
	})
	return out
}

// Ancestors returns this node and its ancestors up to the root.
func (root *Node) Ancestors() []*Node {
	out := []*Node{root}
	for n := root.Parent; n != nil; n = n.Parent {
		out = append(out, n)
	}
	return out
}

// Links returns the parent→child edges of the subtree rooted here.
func (root *Node) Links() []Link {
	links := make([]Link, 0)
	root.Each(func(node *Node) {
		if node != root {
			links = append(links, Link{Source: node.Parent, Target: node})
		}
	})
	return links
}
