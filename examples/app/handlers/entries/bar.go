package entries

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// barScaleData is a single-series dataset whose values span five orders of
// magnitude (12 → 1.2M requests/day per tier). On a linear value axis the
// small tiers collapse to invisible slivers; on a log axis every tier stays
// readable — the point of the value-scale toggle.
func barScaleData() []bar.BarDatum {
	return []bar.BarDatum{
		{"tier": "cache", "requests": 1200000.0},
		{"tier": "cdn", "requests": 340000.0},
		{"tier": "app", "requests": 42000.0},
		{"tier": "db", "requests": 3800.0},
		{"tier": "queue", "requests": 210.0},
		{"tier": "audit", "requests": 12.0},
	}
}

func init() {
	register(ChartEntry{
		Slug:        "bar",
		Title:       "Bar",
		Description: "Stacked/grouped bars over a band index scale.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/bar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(bar.Bar(bar.BarProps{
    Width: 720, Height: 440,
    IndexBy: "country",
    Keys:    []string{"hot dogs", "burgers"},
    Data: []bar.BarDatum{
        {"country": "USA", "hot dogs": 42, "burgers": 33},
        {"country": "Japan", "hot dogs": 28, "burgers": 51},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := bar.BarProps{
				Width: 720, Height: 440, Responsive: true,
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
				Data: []bar.BarDatum{
					{"country": "USA", "hot dogs": 42.0, "burgers": 33.0, "sandwich": 21.0, "kebab": 12.0, "fries": 55.0, "donut": 38.0},
					{"country": "Japan", "hot dogs": 28.0, "burgers": 51.0, "sandwich": 34.0, "kebab": 18.0, "fries": 26.0, "donut": 15.0},
					{"country": "France", "hot dogs": 15.0, "burgers": 60.0, "sandwich": 48.0, "kebab": 41.0, "fries": 32.0, "donut": 9.0},
					{"country": "Germany", "hot dogs": 36.0, "burgers": 44.0, "sandwich": 27.0, "kebab": 58.0, "fries": 39.0, "donut": 17.0},
					{"country": "Brazil", "hot dogs": 24.0, "burgers": 37.0, "sandwich": 19.0, "kebab": 8.0, "fries": 46.0, "donut": 29.0},
					{"country": "UK", "hot dogs": 19.0, "burgers": 49.0, "sandwich": 53.0, "kebab": 35.0, "fries": 61.0, "donut": 22.0},
					{"country": "Mexico", "hot dogs": 47.0, "burgers": 29.0, "sandwich": 16.0, "kebab": 11.0, "fries": 34.0, "donut": 25.0},
					{"country": "Italy", "hot dogs": 12.0, "burgers": 31.0, "sandwich": 44.0, "kebab": 26.0, "fries": 23.0, "donut": 13.0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(bar.Bar(p))
		},
		ScaleRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool, logScale bool) (string, error) {
			p := bar.BarProps{
				Width: 720, Height: 440, Responsive: true,
				IndexBy:   "tier",
				Keys:      []string{"requests"},
				Data:      barScaleData(),
				GroupMode: bar.GroupModeGrouped,
				Theme:     theme,
			}
			// nil ValueScale defaults to the linear spec {min:0, max:auto}; a log
			// spec swaps the value axis to base-10 log. Min is pinned to 1 (not
			// auto) so the smallest tier sits above the baseline with a visible bar
			// rather than exactly on it.
			if logScale {
				p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()}
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(bar.Bar(p))
		},
	})
}
