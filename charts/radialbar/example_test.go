package radialbar_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/radialbar"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleRadialBar renders a small stacked radial-bar chart to an SVG string
// with the charts/render helper.
func ExampleRadialBar() {
	svg, err := render.String(radialbar.RadialBar(radialbar.RadialBarProps{
		Width:  400,
		Height: 400,
		Data: []radialbar.RadialBarSerie{
			{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
				{X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18},
			}},
			{ID: "Online", Data: []radialbar.RadialBarDatum{
				{X: "Vegetables", Y: 8}, {X: "Fruits", Y: 15},
			}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
