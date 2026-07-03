package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/marimekko"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "marimekko",
		Title:       "Marimekko",
		Description: "Variable-width stacked bars: width by value, segments stacked via d3.Stack.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/marimekko"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(marimekko.Marimekko(marimekko.MarimekkoProps{
    Width: 720, Height: 440,
    Dimensions: []marimekko.MarimekkoDimension{
        {ID: "agree", Key: "agree"},
        {ID: "disagree", Key: "disagree"},
    },
    Data: []marimekko.MarimekkoDatum{
        {ID: "France", Value: 42, Dimensions: map[string]float64{"agree": 30, "disagree": 12}},
        {ID: "USA", Value: 63, Dimensions: map[string]float64{"agree": 48, "disagree": 15}},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := marimekko.MarimekkoProps{
				Width: 720, Height: 440, Responsive: true,
				Dimensions: []marimekko.MarimekkoDimension{
					{ID: "agree", Key: "agree"},
					{ID: "disagree", Key: "disagree"},
				},
				Data: []marimekko.MarimekkoDatum{
					{ID: "France", Value: 42, Dimensions: map[string]float64{"agree": 30, "disagree": 12}},
					{ID: "USA", Value: 63, Dimensions: map[string]float64{"agree": 48, "disagree": 15}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(marimekko.Marimekko(p))
		},
	})
}
