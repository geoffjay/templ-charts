package scatterplot_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
)

// ExampleScatterPlot renders two series of x/y points to an SVG string with the
// charts/render helper.
func ExampleScatterPlot() {
	svg, err := render.String(scatterplot.ScatterPlot(scatterplot.ScatterPlotProps{
		Width: 400, Height: 300,
		Data: []scatterplot.ScatterPlotSerie{
			{ID: "group A", Data: []scatterplot.ScatterPlotDatum{
				{X: 8.0, Y: 14.0}, {X: 22.0, Y: 40.0}, {X: 35.0, Y: 9.0},
			}},
			{ID: "group B", Data: []scatterplot.ScatterPlotDatum{
				{X: 12.0, Y: 55.0}, {X: 27.0, Y: 22.0}, {X: 40.0, Y: 78.0},
			}},
		},
		EnableGridX: core.BoolPtr(true),
		EnableGridY: core.BoolPtr(true),
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
