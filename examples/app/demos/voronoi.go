package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/voronoi"
)

// VoronoiDemo is one voronoi tile on the /voronoi page (static render).
type VoronoiDemo struct {
	ID          string
	Title       string
	Description string
	Props       voronoi.VoronoiProps
}

// voronoiSample is a fixed scatter of points in the unit square (the default
// xDomain/yDomain), chosen to give visibly distinct cells.
func voronoiSample() []voronoi.VoronoiDatum {
	return []voronoi.VoronoiDatum{
		{ID: "0", X: 0.12, Y: 0.18},
		{ID: "1", X: 0.42, Y: 0.10},
		{ID: "2", X: 0.78, Y: 0.22},
		{ID: "3", X: 0.90, Y: 0.55},
		{ID: "4", X: 0.66, Y: 0.48},
		{ID: "5", X: 0.30, Y: 0.40},
		{ID: "6", X: 0.08, Y: 0.62},
		{ID: "7", X: 0.24, Y: 0.86},
		{ID: "8", X: 0.55, Y: 0.78},
		{ID: "9", X: 0.84, Y: 0.90},
		{ID: "10", X: 0.48, Y: 0.55},
		{ID: "11", X: 0.70, Y: 0.15},
	}
}

// VoronoiDemos returns the voronoi demos for the /voronoi page.
func VoronoiDemos() []VoronoiDemo {
	data := voronoiSample()
	margin := core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20}
	side := 460.0
	return []VoronoiDemo{
		{
			ID:          "voronoi-cells",
			Title:       "Cells + points (default)",
			Description: "The Voronoi dual clipped to the chart bounds, with the site points. Per-cell hover via Interactive.",
			Props: voronoi.VoronoiProps{
				Width: side, Height: side, Margin: margin,
				Data:        data,
				Interactive: true,
			},
		},
		{
			ID:          "voronoi-links",
			Title:       "Delaunay links + cells",
			Description: "enableLinks: the Delaunay triangulation edges drawn beneath the Voronoi cells.",
			Props: voronoi.VoronoiProps{
				Width: side, Height: side, Margin: margin,
				Data:        data,
				EnableLinks: voronoi.BoolPtr(true),
			},
		},
		{
			ID:          "voronoi-points-only",
			Title:       "Links only (mesh)",
			Description: "Cells off, links on — the bare Delaunay mesh (as used for voronoi-mesh hit-testing).",
			Props: voronoi.VoronoiProps{
				Width: side, Height: side, Margin: margin,
				Data:        data,
				EnableLinks: voronoi.BoolPtr(true),
				EnableCells: voronoi.BoolPtr(false),
			},
		},
	}
}
