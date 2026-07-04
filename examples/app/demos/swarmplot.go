package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
)

// SwarmPlotDemo is one swarmplot tile on the /swarmplot page (static SVG; the
// force layout is deterministic, so the render is stable).
type SwarmPlotDemo struct {
	ID          string
	Title       string
	Description string
	Props       swarmplot.SwarmPlotProps
}

// swarmSample is the shared value spread, sourced from the public samples
// package (samples.SwarmPlot also returns the group list).
func swarmSample() []swarmplot.SwarmPlotDatum {
	d, _ := samples.SwarmPlot()
	return d
}

// SwarmPlotDemos returns the swarmplot demos for the /swarmplot page.
func SwarmPlotDemos() []SwarmPlotDemo {
	data := swarmSample()
	groups := []string{"group A", "group B", "group C"}
	margin := core.Margin{Top: 20, Right: 20, Bottom: 60, Left: 60}
	return []SwarmPlotDemo{
		{
			ID:          "swarm-vertical",
			Title:       "Vertical (default)",
			Description: "Value on the y axis, groups spread along x, relaxed with ForceX/ForceY + ForceCollide (internal/d3/force). Hover a node for its value.",
			Props: swarmplot.SwarmPlotProps{
				Width: 500, Height: 460, Margin: margin,
				Data: data, Groups: groups,
				Size:        8,
				Interactive: true,
			},
		},
		{
			ID:          "swarm-horizontal",
			Title:       "Horizontal",
			Description: "layout: horizontal — value on the x axis, groups along y.",
			Props: swarmplot.SwarmPlotProps{
				Width: 500, Height: 460, Margin: margin,
				Data: data, Groups: groups,
				Size:        8,
				Layout:      "horizontal",
				Interactive: true,
			},
		},
		{
			ID:          "swarm-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover resolves to the nearest node via an accurate Voronoi mesh (internal/d3/delaunay), matching nivo's detection.",
			Props: swarmplot.SwarmPlotProps{
				Width: 500, Height: 460, Margin: margin,
				Data: data, Groups: groups,
				Size:        8,
				Interactive: true,
				UseMesh:     true,
			},
		},
	}
}
