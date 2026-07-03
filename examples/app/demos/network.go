package demos

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/network"
)

// NetworkDemo is one network tile on the /network page (static SVG; the layout
// is deterministic, so the render is stable).
type NetworkDemo struct {
	ID          string
	Title       string
	Description string
	Props       network.NetworkProps
}

var networkPalette = []string{
	"#e8c1a0", "#f47560", "#f1e15b", "#e8a838", "#61cdbb", "#97e3d5",
}

// networkSample builds a small clustered graph: three hubs each with a few
// leaves, plus a couple of cross-links between hubs.
func networkSample() ([]network.NetworkInputNode, []network.NetworkInputLink) {
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

// NetworkDemos returns the network demos for the /network page — a comparison
// that isolates repulsivity (the many-body charge). All three tiles share the
// same graph and a generous LinkDistance (so even the low tile isn't cramped)
// and run the raw force layout — no FitView, so the spread you see is the
// genuine physics: stronger repulsivity pushes the nodes apart and fills more
// of the frame. (FitView would rescale each layout to fill, normalizing away
// exactly this effect — see network.NetworkProps.FitView.)
func NetworkDemos() []NetworkDemo {
	nodes, links := networkSample()
	margin := core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20}
	side := 460.0
	const linkDistance = 90.0
	return []NetworkDemo{
		{
			ID:          "network-low-repulsivity",
			Title:       "Low repulsivity (compact)",
			Description: "repulsivity 40: a gentle charge, so the links keep the graph relatively compact near the center.",
			Props: network.NetworkProps{
				Width: side, Height: side, Margin: margin,
				Nodes: nodes, Links: links,
				LinkDistance: linkDistance,
				Repulsivity:  40,
				Interactive:  true,
			},
		},
		{
			ID:          "network-medium-repulsivity",
			Title:       "Medium repulsivity",
			Description: "repulsivity 120: a stronger charge pushes the clusters apart, spreading the graph across more of the frame.",
			Props: network.NetworkProps{
				Width: side, Height: side, Margin: margin,
				Nodes: nodes, Links: links,
				LinkDistance: linkDistance,
				Repulsivity:  120,
				Interactive:  true,
			},
		},
		{
			ID:          "network-high-repulsivity",
			Title:       "High repulsivity (spread out)",
			Description: "repulsivity 300: the charge dominates, pushing every node strongly apart so the graph fills most of the frame. Hover a node to highlight it, its links and its neighbours (dimming the rest) — a scoped CSS :has() block, no JS — plus a tooltip with its id.",
			Props: network.NetworkProps{
				Width: side, Height: side, Margin: margin,
				Nodes: nodes, Links: links,
				LinkDistance: linkDistance,
				Repulsivity:  300,
				Interactive:  true,
			},
		},
		{
			ID:          "network-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover anywhere resolves to the nearest node via an accurate Voronoi mesh (internal/d3/delaunay) built over the settled node centers, rather than needing to land on a node.",
			Props: network.NetworkProps{
				Width: side, Height: side, Margin: margin,
				Nodes: nodes, Links: links,
				LinkDistance: linkDistance,
				Repulsivity:  120,
				Interactive:  true,
				UseMesh:      true,
			},
		},
	}
}
