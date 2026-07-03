package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "geo",
		Title:       "Geo",
		Description: "Choropleth: a {id,value} dataset bound onto GeoJSON features, colored by a quantize scale.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/geo"
    "github.com/geoffjay/templ-charts/charts/render"
)

features := []geo.Feature{
    {Type: "Feature", ID: "AAA", Geometry: geo.Geometry{
        Type:        geo.TypePolygon,
        Coordinates: [][][2]float64{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
    }},
}

svg, _ := render.String(geo.Choropleth(geo.ChoroplethProps{
    Features: features,
    Data:     []geo.ChoroplethDatum{{ID: "AAA", Value: 10}},
    GeoBase:  geo.GeoBase{Width: 400, Height: 300, ProjectionScale: 60},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			features := []geo.Feature{
				{
					Type: "Feature",
					ID:   "AAA",
					Geometry: geo.Geometry{
						Type: geo.TypePolygon,
						Coordinates: [][][2]float64{{
							{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
						}},
					},
				},
				{
					Type: "Feature",
					ID:   "BBB",
					Geometry: geo.Geometry{
						Type: geo.TypePolygon,
						Coordinates: [][][2]float64{{
							{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 20},
						}},
					},
				},
			}
			p := geo.ChoroplethProps{
				Features: features,
				Data: []geo.ChoroplethDatum{
					{ID: "AAA", Value: 10},
					{ID: "BBB", Value: 90},
				},
				GeoBase: geo.GeoBase{
					Width: 400, Height: 300, Responsive: true,
					ProjectionScale: 60,
					BorderWidth:     0.4,
					BorderColor:     "#152238",
					Theme:           theme,
				},
			}
			p.Animate = animate
			return render.String(geo.Choropleth(p))
		},
	})
}
