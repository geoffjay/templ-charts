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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := line.LineProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []line.LineSeries{
					{ID: "USA", Data: []line.LinePointData{
						{X: "2015", Y: 22.0},
						{X: "2016", Y: 28.0},
						{X: "2017", Y: 34.0},
						{X: "2018", Y: 31.0},
						{X: "2019", Y: 40.0},
						{X: "2020", Y: 26.0},
						{X: "2021", Y: 45.0},
						{X: "2022", Y: 52.0},
						{X: "2023", Y: 58.0},
						{X: "2024", Y: 63.0},
					}},
					{ID: "Japan", Data: []line.LinePointData{
						{X: "2015", Y: 35.0},
						{X: "2016", Y: 33.0},
						{X: "2017", Y: 38.0},
						{X: "2018", Y: 42.0},
						{X: "2019", Y: 39.0},
						{X: "2020", Y: 30.0},
						{X: "2021", Y: 36.0},
						{X: "2022", Y: 41.0},
						{X: "2023", Y: 44.0},
						{X: "2024", Y: 47.0},
					}},
					{ID: "France", Data: []line.LinePointData{
						{X: "2015", Y: 18.0},
						{X: "2016", Y: 21.0},
						{X: "2017", Y: 25.0},
						{X: "2018", Y: 29.0},
						{X: "2019", Y: 27.0},
						{X: "2020", Y: 19.0},
						{X: "2021", Y: 24.0},
						{X: "2022", Y: 31.0},
						{X: "2023", Y: 35.0},
						{X: "2024", Y: 33.0},
					}},
					{ID: "Germany", Data: []line.LinePointData{
						{X: "2015", Y: 41.0},
						{X: "2016", Y: 44.0},
						{X: "2017", Y: 39.0},
						{X: "2018", Y: 47.0},
						{X: "2019", Y: 51.0},
						{X: "2020", Y: 43.0},
						{X: "2021", Y: 49.0},
						{X: "2022", Y: 46.0},
						{X: "2023", Y: 53.0},
						{X: "2024", Y: 56.0},
					}},
					{ID: "Brazil", Data: []line.LinePointData{
						{X: "2015", Y: 12.0},
						{X: "2016", Y: 15.0},
						{X: "2017", Y: 14.0},
						{X: "2018", Y: 19.0},
						{X: "2019", Y: 23.0},
						{X: "2020", Y: 17.0},
						{X: "2021", Y: 22.0},
						{X: "2022", Y: 27.0},
						{X: "2023", Y: 25.0},
						{X: "2024", Y: 30.0},
					}},
					{ID: "UK", Data: []line.LinePointData{
						{X: "2015", Y: 29.0},
						{X: "2016", Y: 32.0},
						{X: "2017", Y: 30.0},
						{X: "2018", Y: 36.0},
						{X: "2019", Y: 34.0},
						{X: "2020", Y: 24.0},
						{X: "2021", Y: 33.0},
						{X: "2022", Y: 38.0},
						{X: "2023", Y: 42.0},
						{X: "2024", Y: 40.0},
					}},
					{ID: "Canada", Data: []line.LinePointData{
						{X: "2015", Y: 25.0},
						{X: "2016", Y: 23.0},
						{X: "2017", Y: 28.0},
						{X: "2018", Y: 26.0},
						{X: "2019", Y: 31.0},
						{X: "2020", Y: 28.0},
						{X: "2021", Y: 35.0},
						{X: "2022", Y: 33.0},
						{X: "2023", Y: 39.0},
						{X: "2024", Y: 43.0},
					}},
					{ID: "Australia", Data: []line.LinePointData{
						{X: "2015", Y: 15.0},
						{X: "2016", Y: 18.0},
						{X: "2017", Y: 22.0},
						{X: "2018", Y: 20.0},
						{X: "2019", Y: 24.0},
						{X: "2020", Y: 21.0},
						{X: "2021", Y: 27.0},
						{X: "2022", Y: 25.0},
						{X: "2023", Y: 29.0},
						{X: "2024", Y: 34.0},
					}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(line.Line(p))
		},
	})
}
