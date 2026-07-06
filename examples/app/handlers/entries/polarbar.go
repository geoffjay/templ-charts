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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := polarbar.PolarBarProps{
				Width: 720, Height: 440, Responsive: true,
				Keys: []string{"walk", "bus", "bike", "car", "train", "tram", "scooter", "ferry"},
				Data: []polarbar.PolarBarDatum{
					{Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3, "car": 6, "train": 4, "tram": 2, "scooter": 1, "ferry": 1}},
					{Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4, "car": 5, "train": 5, "tram": 3, "scooter": 2, "ferry": 1}},
					{Index: "Wed", Values: map[string]float64{"walk": 7, "bus": 6, "bike": 5, "car": 4, "train": 6, "tram": 2, "scooter": 3, "ferry": 2}},
					{Index: "Thu", Values: map[string]float64{"walk": 5, "bus": 8, "bike": 6, "car": 7, "train": 3, "tram": 4, "scooter": 2, "ferry": 1}},
					{Index: "Fri", Values: map[string]float64{"walk": 9, "bus": 4, "bike": 7, "car": 8, "train": 5, "tram": 3, "scooter": 4, "ferry": 2}},
					{Index: "Sat", Values: map[string]float64{"walk": 12, "bus": 3, "bike": 9, "car": 10, "train": 2, "tram": 1, "scooter": 5, "ferry": 3}},
					{Index: "Sun", Values: map[string]float64{"walk": 11, "bus": 2, "bike": 8, "car": 9, "train": 1, "tram": 1, "scooter": 3, "ferry": 4}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(polarbar.PolarBar(p))
		},
	})
}
