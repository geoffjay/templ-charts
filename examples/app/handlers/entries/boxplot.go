package entries

import (
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "boxplot",
		Title:       "Box plot",
		Description: "Quantile box + whisker glyphs from raw observations.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/boxplot"
    "github.com/geoffjay/templ-charts/charts/render"
)

var data []boxplot.BoxPlotDatum
for _, g := range []string{"Alpha", "Beta"} {
    for i := 0; i < 8; i++ {
        data = append(data, boxplot.BoxPlotDatum{Group: g, Value: float64((i*7)%20) + 20})
    }
}
svg, _ := render.String(boxplot.BoxPlot(boxplot.BoxPlotProps{
    Width: 720, Height: 440,
    Data:    data,
    ColorBy: "group",
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			var data []boxplot.BoxPlotDatum
			groups := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta", "Eta", "Theta"}
			for gi, g := range groups {
				base := 20.0 + float64(gi*7)
				spread := 14.0 + float64((gi*5)%11)
				for i := 0; i < 14; i++ {
					v := base + float64((i*17+gi*13)%29)/28.0*spread
					data = append(data, boxplot.BoxPlotDatum{Group: g, Value: v})
				}
			}
			p := boxplot.BoxPlotProps{
				Width: 720, Height: 440, Responsive: true,
				Data:    data,
				ColorBy: "group",
				Theme:   theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(boxplot.BoxPlot(p))
		},
	})
}
