package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/sankey"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "sankey",
		Title:       "Sankey",
		Description: "Flow diagram (d3-sankey): variable-thickness monotone-curve ribbons.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/sankey"
)

svg, _ := render.String(sankey.Sankey(sankey.SankeyProps{
    Width: 720, Height: 440,
    Nodes: []sankey.SankeyInputNode{
        {ID: "Coal"}, {ID: "Grid"}, {ID: "Homes"},
    },
    Links: []sankey.SankeyInputLink{
        {Source: "Coal", Target: "Grid", Value: 12},
        {Source: "Grid", Target: "Homes", Value: 12},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := sankey.SankeyProps{
				Width: 720, Height: 440, Responsive: true,
				Nodes: []sankey.SankeyInputNode{
					{ID: "Coal"},
					{ID: "Grid"},
					{ID: "Homes"},
				},
				Links: []sankey.SankeyInputLink{
					{Source: "Coal", Target: "Grid", Value: 12},
					{Source: "Grid", Target: "Homes", Value: 12},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(sankey.Sankey(p))
		},
	})
}
