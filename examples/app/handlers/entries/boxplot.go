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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			var data []boxplot.BoxPlotDatum
			for _, g := range []string{"Alpha", "Beta"} {
				for i := 0; i < 8; i++ {
					data = append(data, boxplot.BoxPlotDatum{Group: g, Value: float64((i*7)%20) + 20})
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
			return render.String(boxplot.BoxPlot(p))
		},
	})
}
