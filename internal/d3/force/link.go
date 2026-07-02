// link.go — port of d3-force's forceLink: a spring between each linked pair
// pulling them toward a target distance.
package force

import "math"

// LinkForce is d3's forceLink. Distance/Strength default to d3's (30, and the
// degree-based 1/min(deg(source),deg(target))); callers override with
// Distance/DistanceConst and Strength/StrengthConst.
type LinkForce struct {
	links      []*Link
	random     *lcg
	distanceFn func(*Link) float64
	strengthFn func(*Link) float64
	strengths  []float64
	distances  []float64
	bias       []float64
	count      []int
	iterations int
}

// ForceLink creates a link force over links. Source/Target must be nodes that
// belong to the simulation (their Index is read to compute degree bias).
func ForceLink(links []*Link) *LinkForce {
	return &LinkForce{
		links:      links,
		iterations: 1,
		distanceFn: func(*Link) float64 { return 30 },
	}
}

// Distance sets a per-link target distance accessor. Returns f for chaining.
func (f *LinkForce) Distance(fn func(*Link) float64) *LinkForce {
	f.distanceFn = fn
	return f
}

// DistanceConst sets a constant target distance. Returns f for chaining.
func (f *LinkForce) DistanceConst(d float64) *LinkForce {
	f.distanceFn = func(*Link) float64 { return d }
	return f
}

// Strength sets a per-link strength accessor, overriding the degree-based
// default. Returns f for chaining.
func (f *LinkForce) Strength(fn func(*Link) float64) *LinkForce {
	f.strengthFn = fn
	return f
}

// StrengthConst sets a constant strength, overriding the degree-based default.
// Returns f for chaining.
func (f *LinkForce) StrengthConst(s float64) *LinkForce {
	f.strengthFn = func(*Link) float64 { return s }
	return f
}

// Iterations sets the number of relaxation passes per tick (d3 default 1).
func (f *LinkForce) Iterations(n int) *LinkForce {
	if n > 0 {
		f.iterations = n
	}
	return f
}

func (f *LinkForce) initialize(nodes []*Node, random *lcg) {
	f.random = random
	n := len(nodes)
	m := len(f.links)
	f.count = make([]int, n)
	for i, link := range f.links {
		link.Index = i
		f.count[link.Source.Index]++
		f.count[link.Target.Index]++
	}
	f.bias = make([]float64, m)
	for i, link := range f.links {
		cs := float64(f.count[link.Source.Index])
		ct := float64(f.count[link.Target.Index])
		f.bias[i] = cs / (cs + ct)
	}
	f.strengths = make([]float64, m)
	for i, link := range f.links {
		if f.strengthFn != nil {
			f.strengths[i] = f.strengthFn(link)
		} else {
			f.strengths[i] = 1.0 / math.Min(float64(f.count[link.Source.Index]), float64(f.count[link.Target.Index]))
		}
	}
	f.distances = make([]float64, m)
	for i, link := range f.links {
		f.distances[i] = f.distanceFn(link)
	}
}

func (f *LinkForce) apply(alpha float64) {
	for k := 0; k < f.iterations; k++ {
		for i, link := range f.links {
			source, target := link.Source, link.Target
			x := target.X + target.Vx - source.X - source.Vx
			if x == 0 {
				x = jiggle(f.random)
			}
			y := target.Y + target.Vy - source.Y - source.Vy
			if y == 0 {
				y = jiggle(f.random)
			}
			l := math.Sqrt(x*x + y*y)
			l = (l - f.distances[i]) / l * alpha * f.strengths[i]
			x *= l
			y *= l
			b := f.bias[i]
			target.Vx -= x * b
			target.Vy -= y * b
			b = 1 - b
			source.Vx += x * b
			source.Vy += y * b
		}
	}
}
