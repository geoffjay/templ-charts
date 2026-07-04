package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ScatterPlotDemo is one scatterplot tile on the /scatterplot page.
type ScatterPlotDemo struct {
	ID          string
	Title       string
	Description string
	Props       scatterplot.ScatterPlotProps
}

// ScatterPlotDemos returns the scatterplot demos for the /scatterplot page.
func ScatterPlotDemos() []ScatterPlotDemo {
	data := []scatterplot.ScatterPlotSerie{
		{ID: "group A", Data: []scatterplot.ScatterPlotDatum{
			{X: 8.0, Y: 14.0}, {X: 22.0, Y: 40.0}, {X: 35.0, Y: 9.0}, {X: 51.0, Y: 60.0}, {X: 64.0, Y: 28.0}, {X: 78.0, Y: 73.0},
		}},
		{ID: "group B", Data: []scatterplot.ScatterPlotDatum{
			{X: 12.0, Y: 55.0}, {X: 27.0, Y: 22.0}, {X: 40.0, Y: 78.0}, {X: 58.0, Y: 33.0}, {X: 70.0, Y: 90.0}, {X: 88.0, Y: 47.0},
		}},
	}
	return []ScatterPlotDemo{
		{
			ID:          "scatter-basic",
			Title:       "Two series",
			Description: "Two groups on linear x/y scales with grid, axes, and per-point dots (direct point hover arrives with the Phase 5 client layer).",
			Props: scatterplot.ScatterPlotProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
				Data:        data,
				EnableGridX: true,
				EnableGridY: true,
				Interactive: true,
				Title:       "Scatterplot: group A vs group B",
				Desc:        "Two series of x/y points on linear scales; group A and group B each have six observations.",
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorBottomRight, Direction: legends.LegendDirectionColumn, TranslateX: 0, TranslateY: 0},
				},
			},
		},
		{
			ID:          "scatter-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover resolves to the nearest node via an accurate Voronoi mesh (internal/d3/delaunay), matching nivo's detection.",
			Props: scatterplot.ScatterPlotProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
				Data:        data,
				EnableGridX: true,
				EnableGridY: true,
				Interactive: true,
				UseMesh:     true,
				Title:       "Scatterplot: Voronoi-mesh hover",
				Desc:        "Two series of x/y points; hover anywhere resolves to the nearest node via a Voronoi mesh.",
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorBottomRight, Direction: legends.LegendDirectionColumn, TranslateX: 0, TranslateY: 0},
				},
			},
		},
		{
			ID:          "scatter-canvas",
			Title:       "Canvas backend",
			Description: "The same series rendered with the Canvas engine (charts/canvas): the dots are drawn into a <canvas> draw-list while grid and axes stay SVG. Same output, one node instead of one SVG element per point — for large-N data.",
			Props: scatterplot.ScatterPlotProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
				Data:        data,
				EnableGridX: true,
				EnableGridY: true,
				Render:      theming.EngineCanvas,
				ChartID:     "scatter-canvas",
				Title:       "Scatterplot: Canvas backend",
				Desc:        "Two series of x/y points drawn into a canvas draw-list; grid and axes remain SVG panes.",
			},
		},
	}
}
