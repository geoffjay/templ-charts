package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/network"
)

var benchNetworkPalette = []string{
	"#e8c1a0", "#f47560", "#f1e15b", "#e8a838", "#61cdbb", "#97e3d5",
}

// benchNetworkData builds a deterministic graph of n nodes: one hub plus n-1
// leaves, each leaf linked to the hub, with a few extra deterministic
// cross-links so the layout does non-trivial work. The structure is fully
// index-derived, so it is stable across runs.
func benchNetworkData(n int) ([]network.NetworkInputNode, []network.NetworkInputLink) {
	nodes := make([]network.NetworkInputNode, n)
	nodes[0] = network.NetworkInputNode{ID: "node-0", Size: 20, Color: benchNetworkPalette[0]}
	for i := 1; i < n; i++ {
		nodes[i] = network.NetworkInputNode{
			ID:    fmt.Sprintf("node-%d", i),
			Size:  10,
			Color: benchNetworkPalette[i%len(benchNetworkPalette)],
		}
	}

	links := make([]network.NetworkInputLink, 0, n)
	for i := 1; i < n; i++ {
		links = append(links, network.NetworkInputLink{
			Source: "node-0",
			Target: fmt.Sprintf("node-%d", i),
		})
		// A deterministic cross-link every 5th node to add graph structure.
		if i%5 == 0 && i >= 5 {
			links = append(links, network.NetworkInputLink{
				Source: fmt.Sprintf("node-%d", i),
				Target: fmt.Sprintf("node-%d", i-4),
			})
		}
	}
	return nodes, links
}

// benchNetworkProps assembles valid NetworkProps for an n-node graph.
func benchNetworkProps(n int) network.NetworkProps {
	nodes, links := benchNetworkData(n)
	return network.NetworkProps{
		Width:        460,
		Height:       460,
		Margin:       core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
		Nodes:        nodes,
		Links:        links,
		LinkDistance: 90,
		Repulsivity:  120,
	}
}

// BenchmarkNetwork measures the SVG render cost of the network chart as the
// number of nodes (n) grows. Each render runs the force layout, so sizes are
// kept modest.
func BenchmarkNetwork(b *testing.B) {
	for _, n := range []int{20, 50, 100} {
		props := benchNetworkProps(n)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var sb strings.Builder
				if err := network.Network(props).Render(context.Background(), &sb); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
