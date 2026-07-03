package delaunay

import (
	"fmt"
	"testing"
)

// benchPoints returns n deterministic pseudo-random points in [0,1000)^2 using
// d3's LCG constants, so the benchmark is reproducible without math/rand.
func benchPoints(n int) [][2]float64 {
	state := uint32(0x9e3779b9)
	next := func() float64 {
		state = state*1664525 + 1013904223
		return float64(state) / 4294967296.0
	}
	pts := make([][2]float64, n)
	for i := range pts {
		pts[i] = [2]float64{next() * 1000, next() * 1000}
	}
	return pts
}

// BenchmarkNewDelaunayFrom times the Delaunator sweep-hull triangulation
// (O(n log n)) over n deterministic points.
func BenchmarkNewDelaunayFrom(b *testing.B) {
	for _, n := range []int{50, 100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			pts := benchPoints(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = NewDelaunayFrom(pts)
			}
		})
	}
}

// BenchmarkVoronoi times Voronoi cell generation from a prebuilt triangulation.
func BenchmarkVoronoi(b *testing.B) {
	bounds := [4]float64{0, 0, 1000, 1000}
	for _, n := range []int{50, 100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			d := NewDelaunayFrom(benchPoints(n))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = d.Voronoi(bounds)
			}
		})
	}
}
