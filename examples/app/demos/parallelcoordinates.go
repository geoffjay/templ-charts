package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
)

// ParallelCoordinatesDemo is one tile on the /parallel-coordinates page.
type ParallelCoordinatesDemo struct {
	ID          string
	Title       string
	Description string
	Props       pc.PCProps
}

// ParallelCoordinatesDemos returns the demos for the /parallel-coordinates page.
func ParallelCoordinatesDemos() []ParallelCoordinatesDemo {
	vars := []pc.PCVariable{
		{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
		{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
		{Key: "weight", Type: pc.PCScaleLinear, Label: "weight"},
		{Key: "volume", Type: pc.PCScaleLinear, Label: "volume"},
		{Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
	}
	data := []pc.PCDatum{
		{ID: "batch A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "weight": 30.0, "volume": 8.0, "grade": "B"}},
		{ID: "batch B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "weight": 12.0, "volume": 15.0, "grade": "A"}},
		{ID: "batch C", Values: map[string]any{"temp": 12.0, "cost": 2.0, "weight": 25.0, "volume": 3.0, "grade": "C"}},
		{ID: "batch D", Values: map[string]any{"temp": 28.0, "cost": 7.0, "weight": 18.0, "volume": 11.0, "grade": "B"}},
	}
	return []ParallelCoordinatesDemo{
		{
			ID:          "parallelcoordinates-basic",
			Title:       "Four linear variables + a point variable",
			Description: "Each variable gets its own axis; every record is a polyline across them. Hover a line for its id.",
			Props: pc.PCProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 50, Right: 60, Bottom: 50, Left: 60},
				Data:        data,
				Variables:   vars,
				Interactive: true,
			},
		},
		{
			ID:          "parallelcoordinates-vertical",
			Title:       "Vertical layout",
			Description: "layout=vertical stacks the variable axes down the height, with lines running left-to-right.",
			Props: pc.PCProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:    core.Margin{Top: 30, Right: 90, Bottom: 30, Left: 90},
				Data:      data,
				Variables: vars,
				Layout:    pc.PCLayoutVertical,
			},
		},
		{
			ID:          "parallelcoordinates-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover anywhere resolves to the nearest datum via an accurate Voronoi mesh (internal/d3/delaunay) built over each line's per-axis vertices, rather than needing to land on a thin polyline.",
			Props: pc.PCProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 50, Right: 60, Bottom: 50, Left: 60},
				Data:        data,
				Variables:   vars,
				Interactive: true,
				UseMesh:     true,
			},
		},
	}
}
