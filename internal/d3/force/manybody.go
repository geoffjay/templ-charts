// manybody.go — port of d3-force's forceManyBody (n-body charge). By default it
// uses an exact all-pairs sum (equivalent to a Barnes–Hut θ of 0): for
// chart-scale node counts the O(n²) cost is negligible and the result is exact,
// and this is the path the network/swarmplot goldens encode. Calling Theta(θ)
// with θ>0 switches to d3's Barnes–Hut quadtree approximation (internal/d3/
// quadtree), which is O(n log n) for the large-N Canvas charts (docs/PLAN-v6.md
// §3.2). The exact path stays the default so existing goldens are byte-stable —
// a tree traversal sums forces in a different order than the all-pairs loop, so
// the two are not bit-identical even at θ=0. The per-pair math (coincident-node
// jiggle, distanceMin/Max clamping, strength/alpha/l falloff) matches d3 in
// both paths.
package force

import (
	"math"

	"github.com/geoffjay/templ-charts/internal/d3/quadtree"
)

// nodeX / nodeY are the quadtree position accessors over *Node (shared by the
// Barnes–Hut charge and the pruned collision force).
func nodeX(d any) float64 { return d.(*Node).X }
func nodeY(d any) float64 { return d.(*Node).Y }

// nodesToAny boxes a []*Node into the []any the quadtree consumes.
func nodesToAny(nodes []*Node) []any {
	a := make([]any, len(nodes))
	for i, n := range nodes {
		a[i] = n
	}
	return a
}

// ManyBodyForce is d3's forceManyBody. Strength defaults to d3's -30 (negative
// = repulsion); DistanceMin/Max default to 1 and +Inf.
type ManyBodyForce struct {
	nodes        []*Node
	random       *lcg
	strengthFn   func(*Node) float64
	strengths    []float64
	distanceMin2 float64
	distanceMax2 float64
	theta2       float64 // Barnes–Hut θ²; 0 (default) ⇒ exact all-pairs
}

// ForceManyBody creates a many-body (charge) force.
func ForceManyBody() *ManyBodyForce {
	return &ManyBodyForce{
		strengthFn:   func(*Node) float64 { return -30 },
		distanceMin2: 1,
		distanceMax2: math.Inf(1),
	}
}

// Strength sets a per-node charge accessor (negative repels). Returns f.
func (f *ManyBodyForce) Strength(fn func(*Node) float64) *ManyBodyForce {
	f.strengthFn = fn
	return f
}

// StrengthConst sets a constant charge (negative repels). Returns f.
func (f *ManyBodyForce) StrengthConst(s float64) *ManyBodyForce {
	f.strengthFn = func(*Node) float64 { return s }
	return f
}

// DistanceMin sets the minimum inter-node distance considered (below it the
// force is clamped, avoiding blow-ups). Returns f.
func (f *ManyBodyForce) DistanceMin(d float64) *ManyBodyForce {
	f.distanceMin2 = d * d
	return f
}

// DistanceMax sets the maximum distance over which charge acts (+Inf = no
// limit). Returns f.
func (f *ManyBodyForce) DistanceMax(d float64) *ManyBodyForce {
	if math.IsInf(d, 1) {
		f.distanceMax2 = math.Inf(1)
	} else {
		f.distanceMax2 = d * d
	}
	return f
}

// Theta sets the Barnes–Hut accuracy parameter (d3's theta, default there 0.9).
// θ>0 switches ForceManyBody to the O(n log n) quadtree approximation: a cell is
// treated as a single body when its width/distance ratio is below θ. θ=0 (the
// default) keeps the exact O(n²) all-pairs sum, which is what the existing
// force goldens encode — so set this only for large-N layouts. Returns f.
func (f *ManyBodyForce) Theta(theta float64) *ManyBodyForce {
	f.theta2 = theta * theta
	return f
}

func (f *ManyBodyForce) initialize(nodes []*Node, random *lcg) {
	f.nodes = nodes
	f.random = random
	f.strengths = make([]float64, len(nodes))
	for _, node := range nodes {
		f.strengths[node.Index] = f.strengthFn(node)
	}
}

func (f *ManyBodyForce) apply(alpha float64) {
	if f.theta2 > 0 {
		f.applyBarnesHut(alpha)
		return
	}
	nodes := f.nodes
	n := len(nodes)
	for i := 0; i < n; i++ {
		node := nodes[i]
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			other := nodes[j]
			x := other.X - node.X
			y := other.Y - node.Y
			l := x*x + y*y
			if l >= f.distanceMax2 {
				continue
			}
			// Limit forces for very close nodes; randomise direction if coincident.
			if x == 0 {
				x = jiggle(f.random)
				l += x * x
			}
			if y == 0 {
				y = jiggle(f.random)
				l += y * y
			}
			if l < f.distanceMin2 {
				l = math.Sqrt(f.distanceMin2 * l)
			}
			w := f.strengths[other.Index] * alpha / l
			node.Vx += x * w
			node.Vy += y * w
		}
	}
}

// applyBarnesHut is d3-force's quadtree charge: build a quadtree, accumulate the
// total charge and centre of charge on each cell (VisitAfter), then for each
// node walk the tree (Visit), treating a cell as a single body when the
// Barnes–Hut criterion w²/θ² < l holds and otherwise descending. Mirrors
// d3-force's accumulate/apply exactly.
func (f *ManyBodyForce) applyBarnesHut(alpha float64) {
	qt := quadtree.From(nodesToAny(f.nodes), nodeX, nodeY)

	// accumulate: each cell's charge (Value) and centre of charge (X,Y).
	qt.VisitAfter(func(n *quadtree.Node, _, _, _, _ float64) {
		var strength float64
		if !n.IsLeaf() {
			var weight, cx, cy float64
			for i := 0; i < 4; i++ {
				c := n.Child(i)
				if c == nil {
					continue
				}
				abs := math.Abs(c.Value)
				if abs == 0 {
					continue
				}
				strength += c.Value
				weight += abs
				cx += abs * c.X
				cy += abs * c.Y
			}
			n.X = cx / weight
			n.Y = cy / weight
		} else {
			p := n.Data().(*Node)
			n.X, n.Y = p.X, p.Y
			for leaf := n; leaf != nil; leaf = leaf.Next() {
				strength += f.strengths[leaf.Data().(*Node).Index]
			}
		}
		n.Value = strength
	})

	for _, node := range f.nodes {
		node := node
		qt.Visit(func(n *quadtree.Node, x0, _, x1, _ float64) bool {
			if n.Value == 0 {
				return true
			}
			x := n.X - node.X
			y := n.Y - node.Y
			w := x1 - x0
			l := x*x + y*y
			// Barnes–Hut: treat the whole cell as one body.
			if w*w/f.theta2 < l {
				if l < f.distanceMax2 {
					if x == 0 {
						x = jiggle(f.random)
						l += x * x
					}
					if y == 0 {
						y = jiggle(f.random)
						l += y * y
					}
					if l < f.distanceMin2 {
						l = math.Sqrt(f.distanceMin2 * l)
					}
					node.Vx += x * n.Value * alpha / l
					node.Vy += y * n.Value * alpha / l
				}
				return true
			}
			// Otherwise descend internal cells; skip far leaves.
			if !n.IsLeaf() || l >= f.distanceMax2 {
				return false
			}
			// Process the leaf's point(s) directly.
			head := n.Data().(*Node)
			if head != node || n.Next() != nil {
				if x == 0 {
					x = jiggle(f.random)
					l += x * x
				}
				if y == 0 {
					y = jiggle(f.random)
					l += y * y
				}
				if l < f.distanceMin2 {
					l = math.Sqrt(f.distanceMin2 * l)
				}
			}
			for leaf := n; leaf != nil; leaf = leaf.Next() {
				other := leaf.Data().(*Node)
				if other == node {
					continue
				}
				w := f.strengths[other.Index] * alpha / l
				node.Vx += x * w
				node.Vy += y * w
			}
			return true
		})
	}
}
