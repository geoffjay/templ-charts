package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

func init() {
	register(ChartEntry{
		Slug:        "treemap",
		Title:       "Treemap",
		Description: "Nested rectangles (d3-hierarchy): squarify/binary tiling, leaf + parent labels.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/treemap"
)

svg, _ := render.String(treemap.Treemap(treemap.TreemapProps{
    Width: 720, Height: 440,
    Data: treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
        {ID: "viz", Value: 33},
        {ID: "colors", Value: 41},
    }},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := treemap.TreemapProps{
				Width: 720, Height: 440, Responsive: true,
				Data: treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
					{ID: "viz", Value: 33},
					{ID: "colors", Value: 41},
				}},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(treemap.Treemap(p))
		},
	})
}
