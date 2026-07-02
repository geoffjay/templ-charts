// collide.go — port of d3-force's forceCollide: resolves overlap between nodes
// treated as circles of a given radius. d3 uses a quadtree purely to prune
// distant pairs; the impulse it applies is identical to the exact all-pairs
// version here, which visits each unordered pair once (Index j > Index i) so
// the symmetric push is applied a single time — matching d3's `data.index >
// node.index` guard. Strength defaults to 1; the force ignores alpha.
package force

import "math"

// CollideForce is d3's forceCollide.
type CollideForce struct {
	nodes      []*Node
	random     *lcg
	radiusFn   func(*Node) float64
	radii      []float64
	strength   float64
	iterations int
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

func (f *CollideForce) initialize(nodes []*Node, random *lcg) {
	f.nodes = nodes
	f.random = random
	f.radii = make([]float64, len(nodes))
	for _, node := range nodes {
		f.radii[node.Index] = f.radiusFn(node)
	}
}

func (f *CollideForce) apply(float64) {
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
