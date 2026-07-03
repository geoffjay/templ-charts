package entries

import (
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "circle-packing",
		Title:       "Circle packing",
		Description: "Welzl enclosing-circle packing (d3-hierarchy), colored by depth.",
		Snippet: `import (
    cp "github.com/geoffjay/templ-charts/charts/circlepacking"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(cp.CirclePacking(cp.CirclePackingProps{
    Width: 720, Height: 440,
    Data: cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
        {ID: "A", Value: 20},
        {ID: "B", Value: 24},
    }},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := cp.CirclePackingProps{
				Width: 720, Height: 440, Responsive: true,
				Data: cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
					{ID: "A", Value: 20},
					{ID: "B", Value: 24},
				}},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(cp.CirclePacking(p))
		},
	})
}
