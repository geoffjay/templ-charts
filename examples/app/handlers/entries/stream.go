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
				Keys: []string{"Raoul", "Josiane", "Marcel"},
				Data: []stream.StreamDatum{
					{"Raoul": 10, "Josiane": 20, "Marcel": 30},
					{"Raoul": 15, "Josiane": 18, "Marcel": 25},
					{"Raoul": 12, "Josiane": 22, "Marcel": 28},
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
