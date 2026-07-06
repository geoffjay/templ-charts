package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/radialbar"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "radial-bar",
		Title:       "Radial bar",
		Description: "Stacked bars drawn as polar arcs with tracks and axes.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/radialbar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(radialbar.RadialBar(radialbar.RadialBarProps{
    Width: 720, Height: 440,
    Data: []radialbar.RadialBarSerie{
        {ID: "Supermarket", Data: []radialbar.RadialBarDatum{
            {X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18},
        }},
        {ID: "Online", Data: []radialbar.RadialBarDatum{
            {X: "Vegetables", Y: 8}, {X: "Fruits", Y: 15},
        }},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := radialbar.RadialBarProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []radialbar.RadialBarSerie{
					{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
						{X: "Vegetables", Y: 25},
						{X: "Fruits", Y: 18},
						{X: "Meat", Y: 14},
						{X: "Dairy", Y: 12},
						{X: "Bakery", Y: 9},
						{X: "Seafood", Y: 6},
						{X: "Frozen", Y: 8},
						{X: "Beverages", Y: 11},
					}},
					{ID: "Online", Data: []radialbar.RadialBarDatum{
						{X: "Vegetables", Y: 8},
						{X: "Fruits", Y: 15},
						{X: "Meat", Y: 5},
						{X: "Dairy", Y: 7},
						{X: "Bakery", Y: 4},
						{X: "Seafood", Y: 9},
						{X: "Frozen", Y: 13},
						{X: "Beverages", Y: 6},
					}},
					{ID: "Farmers market", Data: []radialbar.RadialBarDatum{
						{X: "Vegetables", Y: 19},
						{X: "Fruits", Y: 12},
						{X: "Meat", Y: 8},
						{X: "Dairy", Y: 5},
						{X: "Bakery", Y: 7},
						{X: "Seafood", Y: 4},
						{X: "Frozen", Y: 2},
						{X: "Beverages", Y: 3},
					}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(radialbar.RadialBar(p))
		},
	})
}
