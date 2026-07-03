package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "sunburst",
		Title:       "Sunburst",
		Description: "Radial partition (d3-hierarchy): arcs by value, colors inherited down the tree.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/sunburst"
)

svg, _ := render.String(sunburst.Sunburst(sunburst.SunburstProps{
    Width: 720, Height: 440,
    Data: sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
        {ID: "fruit", Value: 30},
        {ID: "veg", Value: 20},
    }},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := sunburst.SunburstProps{
				Width: 720, Height: 440, Responsive: true,
				Data: sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
					{ID: "fruit", Value: 30},
					{ID: "veg", Value: 20},
				}},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(sunburst.Sunburst(p))
		},
	})
}
