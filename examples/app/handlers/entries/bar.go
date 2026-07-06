package entries

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "bar",
		Title:       "Bar",
		Description: "Stacked/grouped bars over a band index scale.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/bar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(bar.Bar(bar.BarProps{
    Width: 720, Height: 440,
    IndexBy: "country",
    Keys:    []string{"hot dogs", "burgers"},
    Data: []bar.BarDatum{
        {"country": "USA", "hot dogs": 42, "burgers": 33},
        {"country": "Japan", "hot dogs": 28, "burgers": 51},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := bar.BarProps{
				Width: 720, Height: 440, Responsive: true,
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
				Data: []bar.BarDatum{
					{"country": "USA", "hot dogs": 42.0, "burgers": 33.0, "sandwich": 21.0, "kebab": 12.0, "fries": 55.0, "donut": 38.0},
					{"country": "Japan", "hot dogs": 28.0, "burgers": 51.0, "sandwich": 34.0, "kebab": 18.0, "fries": 26.0, "donut": 15.0},
					{"country": "France", "hot dogs": 15.0, "burgers": 60.0, "sandwich": 48.0, "kebab": 41.0, "fries": 32.0, "donut": 9.0},
					{"country": "Germany", "hot dogs": 36.0, "burgers": 44.0, "sandwich": 27.0, "kebab": 58.0, "fries": 39.0, "donut": 17.0},
					{"country": "Brazil", "hot dogs": 24.0, "burgers": 37.0, "sandwich": 19.0, "kebab": 8.0, "fries": 46.0, "donut": 29.0},
					{"country": "UK", "hot dogs": 19.0, "burgers": 49.0, "sandwich": 53.0, "kebab": 35.0, "fries": 61.0, "donut": 22.0},
					{"country": "Mexico", "hot dogs": 47.0, "burgers": 29.0, "sandwich": 16.0, "kebab": 11.0, "fries": 34.0, "donut": 25.0},
					{"country": "Italy", "hot dogs": 12.0, "burgers": 31.0, "sandwich": 44.0, "kebab": 26.0, "fries": 23.0, "donut": 13.0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(bar.Bar(p))
		},
	})
}
