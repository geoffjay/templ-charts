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
				Keys:    []string{"hot dogs", "burgers"},
				Data: []bar.BarDatum{
					{"country": "USA", "hot dogs": 42.0, "burgers": 33.0},
					{"country": "Japan", "hot dogs": 28.0, "burgers": 51.0},
					{"country": "France", "hot dogs": 15.0, "burgers": 60.0},
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
