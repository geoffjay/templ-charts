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
				Data: []voronoi.VoronoiDatum{
					{ID: "0", X: 0.12, Y: 0.18},
					{ID: "1", X: 0.42, Y: 0.10},
					{ID: "2", X: 0.78, Y: 0.22},
					{ID: "3", X: 0.55, Y: 0.78},
					{ID: "4", X: 0.30, Y: 0.55},
				},
				Theme:   theme,
				Animate: animate,
			}
			return render.String(voronoi.Voronoi(p))
		},
	})
}
