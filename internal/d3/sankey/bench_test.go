package sankey

import (
	"fmt"
	"testing"
)

// buildGraph builds a deterministic layered DAG of n nodes arranged in layers,
// with each node linking to a couple of nodes in the next layer. Link values
// are deterministic (d3's LCG constants) so the benchmark is reproducible.
// The graph is rebuilt for each timed iteration because Layout mutates it.
func buildGraph(n int) *Graph {
	state := uint32(0x9e3779b9)
	next := func() float64 {
		state = state*1664525 + 1013904223
		return float64(state) / 4294967296.0
	}
	// Roughly sqrt(n) layers so the DAG is neither a chain nor a single column.
	layers := 1
	for layers*layers < n {
		layers++
	}
	perLayer := (n + layers - 1) / layers

	nodes := make([]*Node, n)
	id := func(i int) string { return fmt.Sprintf("n%d", i) }
	for i := range nodes {
		nodes[i] = &Node{ID: id(i)}
	}

	var links []*Link
	for i := 0; i < n; i++ {
		layer := i / perLayer
		if layer >= layers-1 {
			continue // last layer are sinks
		}
		nextStart := (layer + 1) * perLayer
		nextEnd := nextStart + perLayer
		if nextEnd > n {
			nextEnd = n
		}
		if nextStart >= n {
			continue
		}
		span := nextEnd - nextStart
		// Two deterministic links into the next layer.
		for k := 0; k < 2; k++ {
			target := nextStart + int(next()*float64(span))
			if target >= n {
				target = n - 1
			}
			links = append(links, &Link{
				SourceID: id(i),
				TargetID: id(target),
				Value:    1 + next()*9,
			})
		}
	}
	return &Graph{Nodes: nodes, Links: links}
}

// BenchmarkLayout times the full sankey layout pass over an n-node DAG.
func BenchmarkLayout(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				g := buildGraph(n)
				b.StartTimer()
				New().NodeAlign(SankeyJustify).NodeWidth(12).NodePadding(4).Size(1200, 800).Layout(g)
			}
		})
	}
}
