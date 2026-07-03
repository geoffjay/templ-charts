package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "swarmplot",
		Title:       "Swarmplot",
		Description: "Grouped value distribution relaxed with d3-force (ForceX/Y + collide).",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/swarmplot"
)

svg, _ := render.String(swarmplot.SwarmPlot(swarmplot.SwarmPlotProps{
    Width: 720, Height: 440,
    Data: []swarmplot.SwarmPlotDatum{
        {ID: "n0", Group: "A", Value: 12},
        {ID: "n1", Group: "B", Value: 47},
    },
    Groups: []string{"A", "B", "C"},
    Size:   8,
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := swarmplot.SwarmPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []swarmplot.SwarmPlotDatum{
					{ID: "n0", Group: "A", Value: 12},
					{ID: "n1", Group: "B", Value: 47},
					{ID: "n2", Group: "C", Value: 23},
					{ID: "n3", Group: "A", Value: 68},
					{ID: "n4", Group: "B", Value: 34},
					{ID: "n5", Group: "C", Value: 55},
				},
				Groups: []string{"A", "B", "C"},
				Size:   8,
				Theme:  theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(swarmplot.SwarmPlot(p))
		},
	})
}
