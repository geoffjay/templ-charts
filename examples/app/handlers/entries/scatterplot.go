package entries

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// scatterDemoSeries builds eight loosely clustered series (vehicle classes,
// weight vs efficiency-ish axes) so each legend entry gets its own visually
// distinct blob and an 8-10 color palette is fully exercised. Deterministic:
// points spiral out from each cluster center at a golden-angle step.
func scatterDemoSeries() []scatterplot.ScatterPlotSerie {
	clusters := []struct {
		id     string
		cx, cy float64
	}{
		{"sedan", 16, 28},
		{"suv", 34, 52},
		{"coupe", 46, 22},
		{"hatchback", 10, 62},
		{"wagon", 26, 76},
		{"minivan", 52, 68},
		{"pickup", 60, 40},
		{"roadster", 42, 88},
	}
	series := make([]scatterplot.ScatterPlotSerie, len(clusters))
	for si, c := range clusters {
		pts := make([]scatterplot.ScatterPlotDatum, 12)
		for i := range pts {
			a := float64(i)*2.399963 + float64(si)*0.7
			r := 1.5 + float64((i*5+si*3)%8)*0.9
			pts[i] = scatterplot.ScatterPlotDatum{
				X: math.Round((c.cx+r*math.Cos(a))*10) / 10,
				Y: math.Round((c.cy+1.6*r*math.Sin(a))*10) / 10,
			}
		}
		series[si] = scatterplot.ScatterPlotSerie{ID: c.id, Data: pts}
	}
	return series
}

func init() {
	register(ChartEntry{
		Slug:        "scatterplot",
		Title:       "Scatterplot",
		Description: "{x,y} nodes on linear scales with grid, axes, and legend.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/scatterplot"
)

svg, _ := render.String(scatterplot.ScatterPlot(scatterplot.ScatterPlotProps{
    Width: 720, Height: 440,
    Data: []scatterplot.ScatterPlotSerie{
        {ID: "group A", Data: []scatterplot.ScatterPlotDatum{
            {X: 8.0, Y: 14.0}, {X: 22.0, Y: 40.0}, {X: 35.0, Y: 9.0},
        }},
        {ID: "group B", Data: []scatterplot.ScatterPlotDatum{
            {X: 12.0, Y: 55.0}, {X: 27.0, Y: 22.0}, {X: 40.0, Y: 78.0},
        }},
    },
    EnableGridX: true,
    EnableGridY: true,
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := scatterplot.ScatterPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data:        scatterDemoSeries(),
				EnableGridX: true,
				EnableGridY: true,
				Theme:       theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(scatterplot.ScatterPlot(p))
		},
		CanvasRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := scatterplot.ScatterPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data:        scatterDemoSeries(),
				EnableGridX: true,
				EnableGridY: true,
				Theme:       theme,
				Render:      theming.EngineCanvas,
				ChartID:     "detail-scatterplot-canvas",
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(scatterplot.ScatterPlot(p))
		},
	})
}
