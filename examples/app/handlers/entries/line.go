package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "line",
		Title:       "Line",
		Description: "One or more series over shared x/y scales.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/line"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(line.Line(line.LineProps{
    Width: 720, Height: 440,
    Data: []line.LineSeries{
        {ID: "USA", Data: []line.LinePointData{
            {X: "2020", Y: 20}, {X: "2021", Y: 45}, {X: "2022", Y: 33},
        }},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := line.LineProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []line.LineSeries{
					{ID: "USA", Data: []line.LinePointData{
						{X: "2020", Y: 20.0}, {X: "2021", Y: 45.0}, {X: "2022", Y: 33.0},
					}},
					{ID: "Japan", Data: []line.LinePointData{
						{X: "2020", Y: 30.0}, {X: "2021", Y: 25.0}, {X: "2022", Y: 52.0},
					}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(line.Line(p))
		},
	})
}
