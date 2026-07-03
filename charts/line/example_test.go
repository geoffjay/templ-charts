package line_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleLine renders a small single-series line chart to an SVG string with
// the charts/render helper.
func ExampleLine() {
	svg, err := render.String(line.Line(line.LineProps{
		Width:  400,
		Height: 300,
		Data: []line.LineSeries{
			{ID: "USA", Data: []line.LinePointData{
				{X: "2020", Y: 20.0},
				{X: "2021", Y: 45.0},
				{X: "2022", Y: 33.0},
			}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
