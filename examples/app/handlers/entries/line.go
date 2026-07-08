package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// lineScaleData is two exponentially growing series (roughly ×3/year) spanning
// ~2 to ~1.4M. On a linear Y axis the early years hug zero as a flat line; on a
// log Y axis the constant growth rate reads as a straight line — the point of
// the value-scale toggle.
func lineScaleData() []line.LineSeries {
	return []line.LineSeries{
		{ID: "signups", Data: []line.LinePointData{
			{X: "2016", Y: 120.0},
			{X: "2017", Y: 640.0},
			{X: "2018", Y: 3100.0},
			{X: "2019", Y: 9800.0},
			{X: "2020", Y: 26000.0},
			{X: "2021", Y: 74000.0},
			{X: "2022", Y: 210000.0},
			{X: "2023", Y: 560000.0},
			{X: "2024", Y: 1400000.0},
		}},
		{ID: "paying", Data: []line.LinePointData{
			{X: "2016", Y: 6.0},
			{X: "2017", Y: 34.0},
			{X: "2018", Y: 190.0},
			{X: "2019", Y: 720.0},
			{X: "2020", Y: 2400.0},
			{X: "2021", Y: 8100.0},
			{X: "2022", Y: 26000.0},
			{X: "2023", Y: 71000.0},
			{X: "2024", Y: 190000.0},
		}},
	}
}

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
		ScaleRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool, logScale bool) (string, error) {
			p := line.LineProps{
				Width: 720, Height: 440, Responsive: true,
				Data:  lineScaleData(),
				Theme: theme,
			}
			// nil YScale defaults to the linear spec {min:0, max:auto}; a log spec
			// swaps the Y axis to base-10 log with an auto (nonzero) floor.
			if logScale {
				p.YScale = scales.ScaleLogSpec{Base: 10, Min: scales.AutoFloat(), Max: scales.AutoFloat()}
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(line.Line(p))
		},
	})
}
