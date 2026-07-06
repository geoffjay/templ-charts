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
				Keys:    []string{"chardonnay", "syrah", "merlot", "riesling", "malbec"},
				Data: []map[string]any{
					{"taste": "fruity", "chardonnay": 93.0, "syrah": 74.0, "merlot": 82.0, "riesling": 118.0, "malbec": 66.0},
					{"taste": "bitter", "chardonnay": 41.0, "syrah": 92.0, "merlot": 68.0, "riesling": 28.0, "malbec": 87.0},
					{"taste": "heavy", "chardonnay": 36.0, "syrah": 109.0, "merlot": 84.0, "riesling": 22.0, "malbec": 102.0},
					{"taste": "spicy", "chardonnay": 28.0, "syrah": 96.0, "merlot": 54.0, "riesling": 34.0, "malbec": 78.0},
					{"taste": "sweet", "chardonnay": 72.0, "syrah": 31.0, "merlot": 47.0, "riesling": 111.0, "malbec": 38.0},
					{"taste": "floral", "chardonnay": 84.0, "syrah": 26.0, "merlot": 42.0, "riesling": 97.0, "malbec": 29.0},
					{"taste": "oaky", "chardonnay": 66.0, "syrah": 88.0, "merlot": 73.0, "riesling": 18.0, "malbec": 91.0},
					{"taste": "earthy", "chardonnay": 33.0, "syrah": 79.0, "merlot": 88.0, "riesling": 24.0, "malbec": 95.0},
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
