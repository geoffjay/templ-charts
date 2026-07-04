package demos

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/voronoi"
)

// VoronoiDemo is one voronoi tile on the /voronoi page (static render).
type VoronoiDemo struct {
	ID          string
	Title       string
	Description string
	Props       voronoi.VoronoiProps
}

// voronoiSample is the shared point scatter, sourced from the public samples
// package.
func voronoiSample() []voronoi.VoronoiDatum {
	return samples.Voronoi()
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
		{
			ID:          "voronoi-colored",
			Title:       "Colored cells",
			Description: "enableCellFill: each cell is filled with its site's ordinal color (the \"set3\" scheme, keyed by datum id) with white cell borders — a Voronoi tessellation map. Hover a cell for its coordinates.",
			Props: voronoi.VoronoiProps{
				Width: side, Height: side, Margin: margin,
				Data:           data,
				EnableCellFill: voronoi.BoolPtr(true),
				Colors:         colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "set3"},
				CellLineColor:  "#ffffff",
				CellLineWidth:  2,
				PointColor:     "#ffffff",
				Interactive:    true,
			},
		},
	}
}
