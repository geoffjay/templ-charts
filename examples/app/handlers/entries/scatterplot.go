package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/theming"
)

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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := scatterplot.ScatterPlotProps{
				Width: 720, Height: 440, Responsive: true,
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
				Theme:       theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(scatterplot.ScatterPlot(p))
		},
	})
}
