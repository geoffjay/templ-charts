package entries

import (
	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "bump",
		Title:       "Bump",
		Description: "Ranking over time: smooth or linear lines with end labels and point hover.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/bump"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(bump.Bump(bump.BumpProps{
    Width: 720, Height: 440,
    Data: []bump.BumpSerie{
        {ID: "React", Data: []bump.BumpDatum{
            {X: "2020", Y: 1}, {X: "2021", Y: 2}, {X: "2022", Y: 1},
        }},
        {ID: "Vue", Data: []bump.BumpDatum{
            {X: "2020", Y: 2}, {X: "2021", Y: 1}, {X: "2022", Y: 2},
        }},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := bump.BumpProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []bump.BumpSerie{
					{ID: "React", Data: []bump.BumpDatum{
						{X: "2020", Y: 1}, {X: "2021", Y: 2}, {X: "2022", Y: 1},
					}},
					{ID: "Vue", Data: []bump.BumpDatum{
						{X: "2020", Y: 2}, {X: "2021", Y: 1}, {X: "2022", Y: 2},
					}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(bump.Bump(p))
		},
	})
}
