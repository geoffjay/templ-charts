package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "heatmap",
		Title:       "Heatmap",
		Description: "2D value grid colored through a sequential/diverging scale.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/heatmap"
    "github.com/geoffjay/templ-charts/charts/render"
)

v := func(f float64) *float64 { return &f }
svg, _ := render.String(heatmap.HeatMap(heatmap.HeatMapProps{
    Width: 720, Height: 440,
    Data: []heatmap.HeatMapSerie{
        {ID: "Japan", Data: []heatmap.HeatMapDatum{
            {X: "Train", Y: v(42)}, {X: "Bus", Y: v(18)},
        }},
        {ID: "France", Data: []heatmap.HeatMapDatum{
            {X: "Train", Y: v(33)}, {X: "Bus", Y: v(51)},
        }},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			v := func(f float64) *float64 { return &f }
			p := heatmap.HeatMapProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []heatmap.HeatMapSerie{
					{ID: "Japan", Data: []heatmap.HeatMapDatum{
						{X: "Train", Y: v(42)}, {X: "Bus", Y: v(18)},
					}},
					{ID: "France", Data: []heatmap.HeatMapDatum{
						{X: "Train", Y: v(33)}, {X: "Bus", Y: v(51)},
					}},
				},
				Theme: theme,
			}
			p.Animate = animate
			return render.String(heatmap.HeatMap(p))
		},
		CanvasRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			v := func(f float64) *float64 { return &f }
			p := heatmap.HeatMapProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []heatmap.HeatMapSerie{
					{ID: "Japan", Data: []heatmap.HeatMapDatum{
						{X: "Train", Y: v(42)}, {X: "Bus", Y: v(18)},
					}},
					{ID: "France", Data: []heatmap.HeatMapDatum{
						{X: "Train", Y: v(33)}, {X: "Bus", Y: v(51)},
					}},
				},
				Theme:   theme,
				Render:  theming.EngineCanvas,
				ChartID: "detail-heatmap-canvas",
			}
			return render.String(heatmap.HeatMap(p))
		},
	})
}
