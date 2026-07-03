package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/polarbar"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "polar-bar",
		Title:       "Polar bar",
		Description: "Stacked bars wrapped into a full circle: angle band per index, radius-stacked keys.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/polarbar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(polarbar.PolarBar(polarbar.PolarBarProps{
    Width: 720, Height: 440,
    Keys: []string{"walk", "bus", "bike"},
    Data: []polarbar.PolarBarDatum{
        {Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3}},
        {Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4}},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := polarbar.PolarBarProps{
				Width: 720, Height: 440, Responsive: true,
				Keys: []string{"walk", "bus", "bike"},
				Data: []polarbar.PolarBarDatum{
					{Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3}},
					{Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(polarbar.PolarBar(p))
		},
	})
}
