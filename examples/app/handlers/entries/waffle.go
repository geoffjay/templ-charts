package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/waffle"
)

func init() {
	register(ChartEntry{
		Slug:        "waffle",
		Title:       "Waffle",
		Description: "Part-of-whole cell grid built on the charts/grid layout.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/waffle"
)

svg, _ := render.String(waffle.Waffle(waffle.WaffleProps{
    Width: 720, Height: 440,
    Total: 100, Rows: 10, Columns: 10,
    Data: []waffle.WaffleDatum{
        {ID: "men", Label: "Men", Value: 32},
        {ID: "women", Label: "Women", Value: 41},
        {ID: "children", Label: "Children", Value: 19},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := waffle.WaffleProps{
				Width: 720, Height: 440, Responsive: true,
				Total: 100, Rows: 10, Columns: 10,
				Data: []waffle.WaffleDatum{
					{ID: "engineering", Label: "Engineering", Value: 28},
					{ID: "marketing", Label: "Marketing", Value: 22},
					{ID: "sales", Label: "Sales", Value: 15},
					{ID: "support", Label: "Support", Value: 11},
					{ID: "operations", Label: "Operations", Value: 9},
					{ID: "hr", Label: "HR", Value: 6},
					{ID: "finance", Label: "Finance", Value: 5},
					{ID: "legal", Label: "Legal", Value: 4},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(waffle.Waffle(p))
		},
	})
}
