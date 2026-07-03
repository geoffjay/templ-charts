package d3hierarchy

import (
	"fmt"
	"testing"
)

// benchTree builds a root datum with n leaf children whose values are
// deterministic (d3's LCG constants), so the layout benchmarks are reproducible.
func benchTree(n int) *datum {
	state := uint32(0x9e3779b9)
	next := func() float64 {
		state = state*1664525 + 1013904223
		return float64(state) / 4294967296.0
	}
	kids := make([]*datum, n)
	for i := range kids {
		kids[i] = &datum{id: "c", val: 1 + next()*99}
	}
	return &datum{id: "root", children: kids}
}

// buildSummed constructs the summed hierarchy fresh for each timed iteration:
// Layout mutates node coordinates in place, so a fresh tree keeps each run
// representative (and the Sum/Hierarchy construction is excluded from timing).
func buildSummed(d *datum) *Node {
	return Hierarchy(d, childrenOf).Sum(func(x any) float64 { return x.(*datum).val })
}

// BenchmarkTreemapSquarify times the default squarify tiling over n leaves.
func BenchmarkTreemapSquarify(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			tree := benchTree(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				root := buildSummed(tree)
				b.StartTimer()
				NewTreemap().Size(1000, 1000).Layout(root)
			}
		})
	}
}

// BenchmarkPack times the enclosing-circle pack layout over n leaves — the
// per-row enclosing-circle step is the non-linear cost.
func BenchmarkPack(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			tree := benchTree(n)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				root := buildSummed(tree)
				b.StartTimer()
				NewPack().Size(1000, 1000).Layout(root)
			}
		})
	}
}
