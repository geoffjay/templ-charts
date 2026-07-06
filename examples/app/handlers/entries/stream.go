package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/stream"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "stream",
		Title:       "Stream",
		Description: "Stacked areas with wiggle/silhouette/expand offsets and a smooth curve.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/stream"
)

svg, _ := render.String(stream.Stream(stream.StreamProps{
    Width: 720, Height: 440,
    Keys: []string{"Raoul", "Josiane", "Marcel"},
    Data: []stream.StreamDatum{
        {"Raoul": 10, "Josiane": 20, "Marcel": 30},
        {"Raoul": 15, "Josiane": 18, "Marcel": 25},
        {"Raoul": 12, "Josiane": 22, "Marcel": 28},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := stream.StreamProps{
				Width: 720, Height: 440, Responsive: true,
				Keys: []string{"Raoul", "Josiane", "Marcel", "René", "Paul", "Jacques", "Simone", "Colette"},
				Data: []stream.StreamDatum{
					{"Raoul": 10, "Josiane": 20, "Marcel": 30, "René": 14, "Paul": 22, "Jacques": 8, "Simone": 17, "Colette": 12},
					{"Raoul": 15, "Josiane": 18, "Marcel": 25, "René": 20, "Paul": 16, "Jacques": 13, "Simone": 21, "Colette": 9},
					{"Raoul": 12, "Josiane": 22, "Marcel": 28, "René": 25, "Paul": 11, "Jacques": 19, "Simone": 14, "Colette": 16},
					{"Raoul": 20, "Josiane": 15, "Marcel": 21, "René": 30, "Paul": 18, "Jacques": 24, "Simone": 10, "Colette": 22},
					{"Raoul": 26, "Josiane": 12, "Marcel": 17, "René": 24, "Paul": 27, "Jacques": 31, "Simone": 13, "Colette": 18},
					{"Raoul": 18, "Josiane": 19, "Marcel": 14, "René": 17, "Paul": 33, "Jacques": 26, "Simone": 20, "Colette": 25},
					{"Raoul": 13, "Josiane": 27, "Marcel": 19, "René": 12, "Paul": 25, "Jacques": 18, "Simone": 28, "Colette": 30},
					{"Raoul": 17, "Josiane": 31, "Marcel": 24, "René": 15, "Paul": 20, "Jacques": 12, "Simone": 34, "Colette": 23},
					{"Raoul": 23, "Josiane": 25, "Marcel": 29, "René": 21, "Paul": 14, "Jacques": 16, "Simone": 27, "Colette": 19},
					{"Raoul": 28, "Josiane": 21, "Marcel": 33, "René": 26, "Paul": 19, "Jacques": 22, "Simone": 18, "Colette": 15},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(stream.Stream(p))
		},
	})
}
