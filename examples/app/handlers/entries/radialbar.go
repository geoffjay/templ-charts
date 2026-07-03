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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := radialbar.RadialBarProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []radialbar.RadialBarSerie{
					{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
						{X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18},
					}},
					{ID: "Online", Data: []radialbar.RadialBarDatum{
						{X: "Vegetables", Y: 8}, {X: "Fruits", Y: 15},
					}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(radialbar.RadialBar(p))
		},
	})
}
