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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := waffle.WaffleProps{
				Width: 720, Height: 440, Responsive: true,
				Total: 100, Rows: 10, Columns: 10,
				Data: []waffle.WaffleDatum{
					{ID: "men", Label: "Men", Value: 32},
					{ID: "women", Label: "Women", Value: 41},
					{ID: "children", Label: "Children", Value: 19},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(waffle.Waffle(p))
		},
	})
}
