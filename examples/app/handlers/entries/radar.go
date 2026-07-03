package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/radar"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "radar",
		Title:       "Radar",
		Description: "Polar line/area chart: values per key around shared indices.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/radar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(radar.Radar(radar.RadarProps{
    Width: 720, Height: 440,
    IndexBy: "taste",
    Keys:    []string{"chardonnay", "syrah"},
    Data: []map[string]any{
        {"taste": "fruity", "chardonnay": 93.0, "syrah": 114.0},
        {"taste": "bitter", "chardonnay": 91.0, "syrah": 72.0},
        {"taste": "heavy", "chardonnay": 56.0, "syrah": 99.0},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := radar.RadarProps{
				Width: 720, Height: 440, Responsive: true,
				IndexBy: "taste",
				Keys:    []string{"chardonnay", "syrah"},
				Data: []map[string]any{
					{"taste": "fruity", "chardonnay": 93.0, "syrah": 114.0},
					{"taste": "bitter", "chardonnay": 91.0, "syrah": 72.0},
					{"taste": "heavy", "chardonnay": 56.0, "syrah": 99.0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(radar.Radar(p))
		},
	})
}
