package demos

import (
	"fmt"
	"math"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/scales"
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

// SwarmPlotScaleDemo returns the value-scale demo tile (id "swarm-scale"): four
// groups whose values sit in different magnitude bands (~2 to ~13000), on a
// linear or log value axis so the toggle can compare both.
func SwarmPlotScaleDemo(logScale bool) SwarmPlotDemo {
	bands := []struct {
		id   string
		base float64
	}{
		{"nano", 2}, {"micro", 30}, {"milli", 450}, {"kilo", 9000},
	}
	groups := make([]string, len(bands))
	var data []swarmplot.SwarmPlotDatum
	for gi, band := range bands {
		groups[gi] = band.id
		for i := 0; i < 14; i++ {
			mult := 1 + 0.5*math.Sin(float64(i)*0.9+float64(gi))
			data = append(data, swarmplot.SwarmPlotDatum{
				ID:    fmt.Sprintf("%s-%d", band.id, i),
				Group: band.id,
				Value: math.Round(band.base*mult*10) / 10,
			})
		}
	}
	p := swarmplot.SwarmPlotProps{
		Width: commonChartWidth, Height: commonChartHeight,
		Margin: core.Margin{Top: 20, Right: 20, Bottom: 60, Left: 60},
		Data:   data, Groups: groups,
		Size:        8,
		Interactive: true,
	}
	if logScale {
		p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	}
	desc := "Groups in different magnitude bands (~2 to ~13000). A log value axis pulls apart bands a linear axis stacks against the floor."
	return SwarmPlotDemo{ID: "swarm-scale", Title: "Value scale (linear / log)", Description: desc, Props: p}
}
