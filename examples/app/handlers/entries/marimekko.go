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
					{ID: "strongly agree", Key: "strongly agree"},
					{ID: "agree", Key: "agree"},
					{ID: "somewhat agree", Key: "somewhat agree"},
					{ID: "neutral", Key: "neutral"},
					{ID: "somewhat disagree", Key: "somewhat disagree"},
					{ID: "disagree", Key: "disagree"},
					{ID: "strongly disagree", Key: "strongly disagree"},
				},
				Data: []marimekko.MarimekkoDatum{
					{ID: "France", Value: 70, Dimensions: map[string]float64{
						"strongly agree": 12, "agree": 18, "somewhat agree": 14, "neutral": 9,
						"somewhat disagree": 7, "disagree": 6, "strongly disagree": 4,
					}},
					{ID: "USA", Value: 84, Dimensions: map[string]float64{
						"strongly agree": 20, "agree": 22, "somewhat agree": 15, "neutral": 11,
						"somewhat disagree": 8, "disagree": 5, "strongly disagree": 3,
					}},
					{ID: "Japan", Value: 71, Dimensions: map[string]float64{
						"strongly agree": 8, "agree": 14, "somewhat agree": 12, "neutral": 15,
						"somewhat disagree": 10, "disagree": 7, "strongly disagree": 5,
					}},
					{ID: "Germany", Value: 80, Dimensions: map[string]float64{
						"strongly agree": 15, "agree": 19, "somewhat agree": 13, "neutral": 10,
						"somewhat disagree": 9, "disagree": 8, "strongly disagree": 6,
					}},
					{ID: "Brazil", Value: 65, Dimensions: map[string]float64{
						"strongly agree": 18, "agree": 16, "somewhat agree": 11, "neutral": 8,
						"somewhat disagree": 6, "disagree": 4, "strongly disagree": 2,
					}},
					{ID: "UK", Value: 74, Dimensions: map[string]float64{
						"strongly agree": 11, "agree": 17, "somewhat agree": 14, "neutral": 12,
						"somewhat disagree": 8, "disagree": 7, "strongly disagree": 5,
					}},
					{ID: "Canada", Value: 69, Dimensions: map[string]float64{
						"strongly agree": 14, "agree": 20, "somewhat agree": 12, "neutral": 9,
						"somewhat disagree": 6, "disagree": 5, "strongly disagree": 3,
					}},
					{ID: "Australia", Value: 68, Dimensions: map[string]float64{
						"strongly agree": 10, "agree": 15, "somewhat agree": 13, "neutral": 11,
						"somewhat disagree": 9, "disagree": 6, "strongly disagree": 4,
					}},
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
