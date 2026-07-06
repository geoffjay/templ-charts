package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/voronoi"
)

func init() {
	register(ChartEntry{
		Slug:        "voronoi",
		Title:       "Voronoi",
		Description: "Delaunay triangulation + Voronoi cells: links, cells, points, bounds.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/voronoi"
)

svg, _ := render.String(voronoi.Voronoi(voronoi.VoronoiProps{
    Width: 440, Height: 440,
    Data: []voronoi.VoronoiDatum{
        {ID: "0", X: 0.12, Y: 0.18},
        {ID: "1", X: 0.42, Y: 0.10},
        {ID: "2", X: 0.78, Y: 0.22},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := voronoi.VoronoiProps{
				Width: 440, Height: 440, Responsive: true,
				// Thirty-two scattered sites spread over the unit square so the
				// tessellation reads as a real point field.
				Data: []voronoi.VoronoiDatum{
					{ID: "0", X: 0.12, Y: 0.18},
					{ID: "1", X: 0.42, Y: 0.10},
					{ID: "2", X: 0.78, Y: 0.22},
					{ID: "3", X: 0.55, Y: 0.78},
					{ID: "4", X: 0.30, Y: 0.55},
					{ID: "5", X: 0.07, Y: 0.42},
					{ID: "6", X: 0.21, Y: 0.33},
					{ID: "7", X: 0.35, Y: 0.24},
					{ID: "8", X: 0.51, Y: 0.31},
					{ID: "9", X: 0.66, Y: 0.09},
					{ID: "10", X: 0.88, Y: 0.13},
					{ID: "11", X: 0.93, Y: 0.34},
					{ID: "12", X: 0.71, Y: 0.41},
					{ID: "13", X: 0.59, Y: 0.52},
					{ID: "14", X: 0.44, Y: 0.44},
					{ID: "15", X: 0.17, Y: 0.62},
					{ID: "16", X: 0.05, Y: 0.79},
					{ID: "17", X: 0.26, Y: 0.86},
					{ID: "18", X: 0.39, Y: 0.68},
					{ID: "19", X: 0.47, Y: 0.91},
					{ID: "20", X: 0.63, Y: 0.66},
					{ID: "21", X: 0.74, Y: 0.84},
					{ID: "22", X: 0.86, Y: 0.61},
					{ID: "23", X: 0.95, Y: 0.77},
					{ID: "24", X: 0.83, Y: 0.94},
					{ID: "25", X: 0.58, Y: 0.94},
					{ID: "26", X: 0.14, Y: 0.94},
					{ID: "27", X: 0.03, Y: 0.09},
					{ID: "28", X: 0.29, Y: 0.05},
					{ID: "29", X: 0.53, Y: 0.05},
					{ID: "30", X: 0.09, Y: 0.28},
					{ID: "31", X: 0.97, Y: 0.49},
				},
				Theme:   theme,
				Animate: animate,
			}
			return render.String(voronoi.Voronoi(p))
		},
	})
}
