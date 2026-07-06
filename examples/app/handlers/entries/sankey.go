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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := sankey.SankeyProps{
				Width: 720, Height: 440, Responsive: true,
				// Energy flows (TWh): six sources feed the grid or storage,
				// which fan out to five end uses.
				Nodes: []sankey.SankeyInputNode{
					{ID: "Coal"},
					{ID: "Gas"},
					{ID: "Nuclear"},
					{ID: "Solar"},
					{ID: "Wind"},
					{ID: "Hydro"},
					{ID: "Grid"},
					{ID: "Storage"},
					{ID: "Homes"},
					{ID: "Industry"},
					{ID: "Commercial"},
					{ID: "Transport"},
					{ID: "Losses"},
				},
				Links: []sankey.SankeyInputLink{
					{Source: "Coal", Target: "Grid", Value: 22},
					{Source: "Gas", Target: "Grid", Value: 34},
					{Source: "Nuclear", Target: "Grid", Value: 26},
					{Source: "Solar", Target: "Grid", Value: 14},
					{Source: "Solar", Target: "Storage", Value: 6},
					{Source: "Wind", Target: "Grid", Value: 18},
					{Source: "Wind", Target: "Storage", Value: 8},
					{Source: "Hydro", Target: "Grid", Value: 12},
					{Source: "Grid", Target: "Homes", Value: 44},
					{Source: "Grid", Target: "Industry", Value: 38},
					{Source: "Grid", Target: "Commercial", Value: 26},
					{Source: "Grid", Target: "Transport", Value: 9},
					{Source: "Grid", Target: "Losses", Value: 9},
					{Source: "Storage", Target: "Homes", Value: 6},
					{Source: "Storage", Target: "Transport", Value: 8},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(sankey.Sankey(p))
		},
	})
}
