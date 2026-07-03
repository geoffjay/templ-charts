// collide.go — port of d3-force's forceCollide: resolves overlap between nodes
// treated as circles of a given radius. By default this uses an exact all-pairs
// pass, visiting each unordered pair once (Index j > Index i) so the symmetric
// push is applied a single time — matching d3's `data.index > node.index` guard.
// That exact pass is what the swarmplot golden encodes, so it stays the default.
//
// UseQuadtree() switches to d3's quadtree-pruned pass (internal/d3/quadtree),
// which skips whole cells that are too far to overlap — O(n log n) for the
// large-N Canvas charts (docs/PLAN-v6.md §3.2). The per-pair impulse is the same
// math, but (a) the traversal sums in a different order (not bit-identical) and
// (b), exactly like d3, a coincident-point leaf only collides via its head
// datum — so the pruned pass is opt-in, not the default. Strength defaults to 1;
// the force ignores alpha.
package force

import (
	"math"

	"github.com/geoffjay/templ-charts/internal/d3/quadtree"
)

// CollideForce is d3's forceCollide.
type CollideForce struct {
	nodes       []*Node
	random      *lcg
	radiusFn    func(*Node) float64
	radii       []float64
	strength    float64
	iterations  int
	useQuadtree bool // opt into the O(n log n) quadtree-pruned pass
}

// ForceCollide creates a collision force with a per-node radius accessor.
func ForceCollide(fn func(*Node) float64) *CollideForce {
	return &CollideForce{radiusFn: fn, strength: 1, iterations: 1}
}

// Strength sets the fraction of overlap resolved per pass (d3 default 1).
func (f *CollideForce) Strength(s float64) *CollideForce {
	f.strength = s
	return f
}

// Iterations sets the number of relaxation passes per tick (d3 default 1).
func (f *CollideForce) Iterations(n int) *CollideForce {
	if n > 0 {
		f.iterations = n
	}
	return f
}

// UseQuadtree switches to the O(n log n) quadtree-pruned collision pass. Keep it
// off (the default) for chart-scale layouts whose goldens encode the exact
// all-pairs pass; enable it for the large-N Canvas charts. Returns f.
func (f *CollideForce) UseQuadtree() *CollideForce {
	f.useQuadtree = true
	return f
}

func (f *CollideForce) initialize(nodes []*Node, random *lcg) {
	f.nodes = nodes
	f.random = random
	f.radii = make([]float64, len(nodes))
	for _, node := range nodes {
		f.radii[node.Index] = f.radiusFn(node)
	}
}

func (f *CollideForce) apply(float64) {
	if f.useQuadtree {
		f.applyQuadtree()
		return
	}
	nodes := f.nodes
	n := len(nodes)
	for k := 0; k < f.iterations; k++ {
		for i := 0; i < n; i++ {
			node := nodes[i]
			ri := f.radii[node.Index]
			ri2 := ri * ri
			xi := node.X + node.Vx
			yi := node.Y + node.Vy
			for j := 0; j < n; j++ {
				data := nodes[j]
				if data.Index <= node.Index {
					continue
				}
				rj := f.radii[data.Index]
				r := ri + rj
				x := xi - data.X - data.Vx
				y := yi - data.Y - data.Vy
				l := x*x + y*y
				if l >= r*r {
					continue
				}
				if x == 0 {
					x = jiggle(f.random)
					l += x * x
				}
				if y == 0 {
					y = jiggle(f.random)
					l += y * y
				}
				l = math.Sqrt(l)
				l = (r - l) / l * f.strength
				x *= l
				y *= l
				rj2 := rj * rj
				ratio := rj2 / (ri2 + rj2)
				node.Vx += x * ratio
				node.Vy += y * ratio
				inv := 1 - ratio
				data.Vx -= x * inv
				data.Vy -= y * inv
			}
		}
	}
}

// applyQuadtree is d3-force's quadtree-pruned collision: each relaxation pass
// rebuilds a quadtree, tags every cell (VisitAfter) with the maximum radius in
// its subtree (stored in Node.Value), then for each node walks the tree
// (Visit), pruning cells that cannot overlap and resolving the ones that do
// against the higher-indexed datum. Mirrors d3-force's prepare/apply.
func (f *CollideForce) applyQuadtree() {
	for k := 0; k < f.iterations; k++ {
		qt := quadtree.From(nodesToAny(f.nodes), nodeX, nodeY)
		// prepare: Node.Value = max radius in the cell.
		qt.VisitAfter(func(n *quadtree.Node, _, _, _, _ float64) {
			if n.IsLeaf() {
				n.Value = f.radii[n.Data().(*Node).Index]
				return
			}
			r := 0.0
			for i := 0; i < 4; i++ {
				if c := n.Child(i); c != nil && c.Value > r {
					r = c.Value
				}
			}
			n.Value = r
		})
		for _, node := range f.nodes {
			node := node
			ri := f.radii[node.Index]
			ri2 := ri * ri
			xi := node.X + node.Vx
			yi := node.Y + node.Vy
			qt.Visit(func(n *quadtree.Node, x0, y0, x1, y1 float64) bool {
				if n.IsLeaf() {
					data := n.Data().(*Node)
					rj := n.Value
					if data.Index > node.Index {
						x := xi - data.X - data.Vx
						y := yi - data.Y - data.Vy
						l := x*x + y*y
						r := ri + rj
						if l < r*r {
							if x == 0 {
								x = jiggle(f.random)
								l += x * x
							}
							if y == 0 {
								y = jiggle(f.random)
								l += y * y
							}
							l = math.Sqrt(l)
							l = (r - l) / l * f.strength
							x *= l
							y *= l
							rj2 := rj * rj
							ratio := rj2 / (ri2 + rj2)
							node.Vx += x * ratio
							node.Vy += y * ratio
							inv := 1 - ratio
							data.Vx -= x * inv
							data.Vy -= y * inv
						}
					}
					return true
				}
				// Internal cell: prune if it cannot reach node's collision radius.
				r := ri + n.Value
				return x0 > xi+r || x1 < xi-r || y0 > yi+r || y1 < yi-r
			})
		}
	}
}
