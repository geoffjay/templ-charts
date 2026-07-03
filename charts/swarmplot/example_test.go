package swarmplot_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
)

// ExampleSwarmPlot renders a small swarm plot of grouped values to an SVG
// string with the charts/render helper.
func ExampleSwarmPlot() {
	svg, err := render.String(swarmplot.SwarmPlot(swarmplot.SwarmPlotProps{
		Width: 500, Height: 460,
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
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
