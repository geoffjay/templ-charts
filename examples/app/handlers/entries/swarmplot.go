package entries

import (
	"fmt"
	"math"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scales"
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
		ScaleRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool, logScale bool) (string, error) {
			data, groups := swarmScaleData()
			p := swarmplot.SwarmPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data:   data,
				Groups: groups,
				Size:   8,
				Theme:  theme,
			}
			// nil ValueScale defaults to linear; a log spec swaps the value axis to
			// base-10 log so the magnitude bands separate instead of piling up.
			if logScale {
				p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.AutoFloat(), Max: scales.AutoFloat()}
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(swarmplot.SwarmPlot(p))
		},
	})
}

// swarmScaleData builds four groups whose values sit in different magnitude
// bands (~2 up to ~13000), so the value-scale toggle shows a log axis pulling
// apart bands that a linear axis stacks against the floor.
func swarmScaleData() ([]swarmplot.SwarmPlotDatum, []string) {
	bands := []struct {
		id   string
		base float64
	}{
		{"nano", 2}, {"micro", 30}, {"milli", 450}, {"kilo", 9000},
	}
	groups := make([]string, len(bands))
	var data []swarmplot.SwarmPlotDatum
	for gi, band := range bands {
		groups[gi] = band.id
		for i := 0; i < 14; i++ {
			mult := 1 + 0.5*math.Sin(float64(i)*0.9+float64(gi))
			data = append(data, swarmplot.SwarmPlotDatum{
				ID:    fmt.Sprintf("%s-%d", band.id, i),
				Group: band.id,
				Value: math.Round(band.base*mult*10) / 10,
			})
		}
	}
	return data, groups
}
