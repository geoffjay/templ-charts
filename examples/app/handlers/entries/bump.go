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
						{X: "2017", Y: 1},
						{X: "2018", Y: 1},
						{X: "2019", Y: 1},
						{X: "2020", Y: 1},
						{X: "2021", Y: 2},
						{X: "2022", Y: 1},
						{X: "2023", Y: 1},
						{X: "2024", Y: 1},
					}},
					{ID: "Vue", Data: []bump.BumpDatum{
						{X: "2017", Y: 3},
						{X: "2018", Y: 2},
						{X: "2019", Y: 2},
						{X: "2020", Y: 2},
						{X: "2021", Y: 1},
						{X: "2022", Y: 2},
						{X: "2023", Y: 3},
						{X: "2024", Y: 3},
					}},
					{ID: "Angular", Data: []bump.BumpDatum{
						{X: "2017", Y: 2},
						{X: "2018", Y: 3},
						{X: "2019", Y: 3},
						{X: "2020", Y: 4},
						{X: "2021", Y: 5},
						{X: "2022", Y: 5},
						{X: "2023", Y: 6},
						{X: "2024", Y: 6},
					}},
					{ID: "Svelte", Data: []bump.BumpDatum{
						{X: "2017", Y: 6},
						{X: "2018", Y: 6},
						{X: "2019", Y: 4},
						{X: "2020", Y: 3},
						{X: "2021", Y: 3},
						{X: "2022", Y: 3},
						{X: "2023", Y: 2},
						{X: "2024", Y: 2},
					}},
					{ID: "Solid", Data: []bump.BumpDatum{
						{X: "2017", Y: 7},
						{X: "2018", Y: 7},
						{X: "2019", Y: 7},
						{X: "2020", Y: 6},
						{X: "2021", Y: 4},
						{X: "2022", Y: 4},
						{X: "2023", Y: 4},
						{X: "2024", Y: 5},
					}},
					{ID: "Preact", Data: []bump.BumpDatum{
						{X: "2017", Y: 5},
						{X: "2018", Y: 4},
						{X: "2019", Y: 5},
						{X: "2020", Y: 5},
						{X: "2021", Y: 6},
						{X: "2022", Y: 7},
						{X: "2023", Y: 7},
						{X: "2024", Y: 7},
					}},
					{ID: "Ember", Data: []bump.BumpDatum{
						{X: "2017", Y: 4},
						{X: "2018", Y: 5},
						{X: "2019", Y: 6},
						{X: "2020", Y: 7},
						{X: "2021", Y: 8},
						{X: "2022", Y: 8},
						{X: "2023", Y: 8},
						{X: "2024", Y: 8},
					}},
					{ID: "Qwik", Data: []bump.BumpDatum{
						{X: "2017", Y: 8},
						{X: "2018", Y: 8},
						{X: "2019", Y: 8},
						{X: "2020", Y: 8},
						{X: "2021", Y: 7},
						{X: "2022", Y: 6},
						{X: "2023", Y: 5},
						{X: "2024", Y: 4},
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
