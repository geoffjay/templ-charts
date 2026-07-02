package demos

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/core"
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

// swarmSample builds a spread of values across three groups (a small
// pseudo-random but fixed sequence so the demo is stable).
func swarmSample() []swarmplot.SwarmPlotDatum {
	groups := []string{"group A", "group B", "group C"}
	// A fixed value sequence per group index, enough points to show the swarm.
	seq := []float64{
		12, 47, 23, 68, 34, 89, 5, 56, 78, 30, 41, 62,
		19, 73, 27, 51, 44, 8, 95, 60, 15, 38, 82, 49,
		22, 66, 11, 90, 33, 57, 71, 4, 84, 26, 53, 40,
	}
	data := make([]swarmplot.SwarmPlotDatum, 0, len(seq))
	for i, v := range seq {
		g := groups[i%len(groups)]
		data = append(data, swarmplot.SwarmPlotDatum{
			ID: "n" + strconv.Itoa(i), Group: g, Value: v,
		})
	}
	return data
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
