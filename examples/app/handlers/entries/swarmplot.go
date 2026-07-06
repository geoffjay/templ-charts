package entries

import (
	"fmt"

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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			groups := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
			var data []swarmplot.SwarmPlotDatum
			for gi, g := range groups {
				center := 25.0 + float64((gi*11)%40)
				for i := 0; i < 12; i++ {
					v := center + float64((i*19+gi*7)%31) - 15.0
					data = append(data, swarmplot.SwarmPlotDatum{
						ID:    fmt.Sprintf("%s-%d", g, i),
						Group: g,
						Value: v,
					})
				}
			}
			p := swarmplot.SwarmPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data:   data,
				Groups: groups,
				Size:   8,
				Theme:  theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(swarmplot.SwarmPlot(p))
		},
	})
}
