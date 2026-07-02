// manybody.go — port of d3-force's forceManyBody (n-body charge). This uses an
// exact all-pairs sum rather than d3's Barnes–Hut quadtree approximation: for
// chart-scale node counts the O(n²) cost is negligible and the result is more
// accurate (equivalent to θ=0). The per-pair math — coincident-node jiggle,
// distanceMin/Max clamping and the strength/alpha/l³ falloff — matches d3's
// "process points directly" leaf branch exactly.
package force

import "math"

// ManyBodyForce is d3's forceManyBody. Strength defaults to d3's -30 (negative
// = repulsion); DistanceMin/Max default to 1 and +Inf.
type ManyBodyForce struct {
	nodes        []*Node
	random       *lcg
	strengthFn   func(*Node) float64
	strengths    []float64
	distanceMin2 float64
	distanceMax2 float64
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

func (f *ManyBodyForce) initialize(nodes []*Node, random *lcg) {
	f.nodes = nodes
	f.random = random
	f.strengths = make([]float64, len(nodes))
	for _, node := range nodes {
		f.strengths[node.Index] = f.strengthFn(node)
	}
}

func (f *ManyBodyForce) apply(alpha float64) {
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
