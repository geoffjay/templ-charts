package polarbar_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/polarbar"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExamplePolarBar renders a small stacked polar-bar chart to an SVG string with
// the charts/render helper.
func ExamplePolarBar() {
	svg, err := render.String(polarbar.PolarBar(polarbar.PolarBarProps{
		Width: 400, Height: 400,
		Keys: []string{"walk", "bus", "bike"},
		Data: []polarbar.PolarBarDatum{
			{Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3}},
			{Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
