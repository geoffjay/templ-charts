package boxplot_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleBoxPlot renders a small box-and-whisker chart to an SVG string with
// the charts/render helper.
func ExampleBoxPlot() {
	var data []boxplot.BoxPlotDatum
	for _, g := range []string{"Alpha", "Beta"} {
		for i := 0; i < 8; i++ {
			data = append(data, boxplot.BoxPlotDatum{Group: g, Value: float64((i*7)%20) + 20})
		}
	}
	svg, err := render.String(boxplot.BoxPlot(boxplot.BoxPlotProps{
		Width: 400, Height: 300,
		Data:    data,
		ColorBy: "group",
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
