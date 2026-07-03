package waffle_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/waffle"
)

// ExampleWaffle renders a small part-of-whole waffle chart to an SVG string
// with the charts/render helper.
func ExampleWaffle() {
	svg, err := render.String(waffle.Waffle(waffle.WaffleProps{
		Width:   300,
		Height:  300,
		Total:   100,
		Rows:    10,
		Columns: 10,
		Data: []waffle.WaffleDatum{
			{ID: "men", Label: "Men", Value: 32},
			{ID: "women", Label: "Women", Value: 41},
			{ID: "children", Label: "Children", Value: 19},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
