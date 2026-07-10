package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/examples/app/demos"
)

func init() {
	register(ChartEntry{
		Slug:        "geo",
		Title:       "Geo",
		Description: "Choropleth: a {id,value} dataset bound onto GeoJSON features, colored by a quantize scale. The projection is auto-fit to the frame (fitExtent) — no manual scale/center.",
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
    // Fit auto-scales/centers the projection to the frame (d3-geo fitExtent).
    GeoBase:  geo.GeoBase{Width: 400, Height: 300, Fit: true},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			// The real embedded world map with deterministic per-country values —
			// a much better palette showcase than abstract polygons.
			features := demos.WorldFeatures()
			p := geo.ChoroplethProps{
				Features: features,
				Data:     demos.WorldChoroplethData(features),
				GeoBase: geo.GeoBase{
					Width: 720, Height: 440, Responsive: true,
					Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
					// Explicit Natural Earth projection spanning the inner width
					// (Fit/fitExtent NaNs on the world MultiPolygons; the /geo demo
					// page uses this same explicit-scale approach).
					ProjectionType:  geo.ProjectionNaturalEarth1,
					ProjectionScale: (720 - 20) / (2 * 2.73),
					BorderWidth:     0.4,
					BorderColor:     "#152238",
					Theme:           theme,
				},
				ValueFormat: ",.0f",
			}
			// Colors is a quantize scheme id: gradient palettes map directly onto
			// it; categorical picks keep the default scheme.
			if pal, ok := colors.LookupPalette(palette); ok && pal.Kind != colors.KindCategorical {
				p.Colors = string(palette)
			}
			p.Animate = animate
			return render.String(geo.Choropleth(p))
		},
	})
}
