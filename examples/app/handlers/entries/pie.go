package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "pie",
		Title:       "Pie",
		Description: "Proportional arcs from id/value data.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/pie"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(pie.Pie(pie.PieProps{
    Width: 720, Height: 440,
    Data: []any{
        map[string]any{"id": "Go", "value": 40.0},
        map[string]any{"id": "Rust", "value": 25.0},
        map[string]any{"id": "Python", "value": 35.0},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := pie.PieProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []any{
					map[string]any{"id": "Go", "value": 40.0},
					map[string]any{"id": "Rust", "value": 25.0},
					map[string]any{"id": "Python", "value": 35.0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(pie.Pie(p))
		},
	})
}
