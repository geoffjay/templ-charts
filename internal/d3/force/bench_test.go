package force

import (
	"fmt"
	"testing"
)

// benchLCG is a tiny deterministic pseudo-random source (the same LCG constants
// d3 uses) so benchmark inputs are reproducible without math/rand.
type benchLCG struct{ state uint32 }

func (r *benchLCG) next() float64 {
	r.state = r.state*1664525 + 1013904223
	return float64(r.state) / 4294967296.0
}

// buildNetwork constructs a deterministic network of n nodes with a few links
// per node, plus the charge/link/center forces the network chart uses. It is
// rebuilt fresh for each timed iteration because Tick mutates node positions in
// place — reusing a settled simulation would not measure a representative run.
func buildNetwork(n int) *Simulation {
	rng := &benchLCG{state: 0x9e3779b9}
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	// A few deterministic links per node (a sparse graph, like a real network).
	var links []*Link
	for i := 0; i < n; i++ {
		degree := 1 + int(rng.next()*3) // 1..3 links out of each node
		for k := 0; k < degree; k++ {
			j := int(rng.next() * float64(n))
			if j != i {
				links = append(links, NewLink(nodes[i], nodes[j], nil))
			}
		}
	}
	return NewSimulation(nodes).
		Force("link", ForceLink(links).DistanceConst(30).StrengthConst(1)).
		Force("charge", ForceManyBody().StrengthConst(-30).DistanceMin(1)).
		Force("center", ForceCenter(300, 300)).
		Stop()
}

// BenchmarkSimulation times a fixed 120-tick run (nivo's default) of the full
// network simulation — the O(n^2) many-body sum dominates.
func BenchmarkSimulation(b *testing.B) {
	for _, n := range []int{50, 100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				sim := buildNetwork(n)
				b.StartTimer()
				sim.Tick(120)
			}
		})
	}
}

// buildCollide builds a swarm of n nodes with x/y positioning plus the O(n^2)
// collide force (the swarmplot use case).
func buildCollide(n int) *Simulation {
	const size = 6.0
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	radius := func(*Node) float64 { return size/2 + 1 }
	return NewSimulation(nodes).
		Force("x", ForceX(func(*Node) float64 { return 300 }).Strength(1)).
		Force("y", ForceY(func(*Node) float64 { return 300 })).
		Force("collide", ForceCollide(radius)).
		Stop()
}

// BenchmarkCollide times a 120-tick run of the collide-dominated swarm layout.
func BenchmarkCollide(b *testing.B) {
	for _, n := range []int{50, 100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				sim := buildCollide(n)
				b.StartTimer()
				sim.Tick(120)
			}
		})
	}
}

// buildNetworkTheta is buildNetwork with the Barnes–Hut charge (θ=0.9).
func buildNetworkTheta(n int, theta float64) *Simulation {
	rng := &benchLCG{state: 0x9e3779b9}
	nodes := make([]*Node, n)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	var links []*Link
	for i := 0; i < n; i++ {
		degree := 1 + int(rng.next()*3)
		for k := 0; k < degree; k++ {
			j := int(rng.next() * float64(n))
			if j != i {
				links = append(links, NewLink(nodes[i], nodes[j], nil))
			}
		}
	}
	return NewSimulation(nodes).
		Force("link", ForceLink(links).DistanceConst(30).StrengthConst(1)).
		Force("charge", ForceManyBody().StrengthConst(-30).DistanceMin(1).Theta(theta)).
		Force("center", ForceCenter(300, 300)).
		Stop()
}

// BenchmarkSimulationBarnesHut times the same run with the O(n log n) Barnes–Hut
// charge — compare against BenchmarkSimulation to see the asymptotic win at
// large n (the two diverge as n grows).
func BenchmarkSimulationBarnesHut(b *testing.B) {
	for _, n := range []int{50, 100, 500, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				sim := buildNetworkTheta(n, 0.9)
				b.StartTimer()
				sim.Tick(120)
			}
		})
	}
}

// BenchmarkCollideQuadtree times the swarm layout with the O(n log n) pruned
// collision pass — compare against BenchmarkCollide.
func BenchmarkCollideQuadtree(b *testing.B) {
	const size = 6.0
	for _, n := range []int{50, 100, 500, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				nodes := make([]*Node, n)
				for j := range nodes {
					nodes[j] = NewNode(j)
				}
				radius := func(*Node) float64 { return size/2 + 1 }
				sim := NewSimulation(nodes).
					Force("x", ForceX(func(*Node) float64 { return 300 }).Strength(1)).
					Force("y", ForceY(func(*Node) float64 { return 300 })).
					Force("collide", ForceCollide(radius).UseQuadtree()).
					Stop()
				b.StartTimer()
				sim.Tick(120)
			}
		})
	}
}
