package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "icicle",
		Title:       "Icicle",
		Description: "Depth-banded partition rectangles (d3-hierarchy), oriented four ways.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/icicle"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(icicle.Icicle(icicle.IcicleProps{
    Width: 720, Height: 440,
    Data: icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
        {ID: "analytics", Value: 40},
        {ID: "billing", Value: 22},
    }},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := icicle.IcicleProps{
				Width: 720, Height: 440, Responsive: true,
				Data: icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
					{ID: "analytics", Value: 40},
					{ID: "billing", Value: 22},
				}},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(icicle.Icicle(p))
		},
	})
}
