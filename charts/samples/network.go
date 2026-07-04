package samples

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/network"
)

// networkPalette colors the sample graph's hubs and their leaf clusters.
var networkPalette = []string{
	"#e8c1a0", "#f47560", "#f1e15b", "#e8a838", "#61cdbb", "#97e3d5",
}

// Network returns ready-to-render network data: a small clustered graph of
// three hubs (each with four leaves) plus cross-links between the hubs. The
// layout the chart computes from this is deterministic, so it renders stably.
func Network() ([]network.NetworkInputNode, []network.NetworkInputLink) {
	nodes := []network.NetworkInputNode{
		{ID: "hub-0", Size: 20, Color: networkPalette[0]},
		{ID: "hub-1", Size: 20, Color: networkPalette[1]},
		{ID: "hub-2", Size: 20, Color: networkPalette[2]},
	}
	links := []network.NetworkInputLink{
		{Source: "hub-0", Target: "hub-1"},
		{Source: "hub-1", Target: "hub-2"},
		{Source: "hub-2", Target: "hub-0"},
	}
	hubs := []string{"hub-0", "hub-1", "hub-2"}
	for h, hub := range hubs {
		for i := 0; i < 4; i++ {
			id := hub + "-leaf-" + strconv.Itoa(i)
			nodes = append(nodes, network.NetworkInputNode{ID: id, Size: 10, Color: networkPalette[3+h]})
			links = append(links, network.NetworkInputLink{Source: hub, Target: id})
		}
	}
	return nodes, links
}
