package sunburst_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/sunburst"
)

// ExampleSunburst renders a small sunburst to an SVG string with the
// charts/render helper.
func ExampleSunburst() {
	svg, err := render.String(sunburst.Sunburst(sunburst.SunburstProps{
		Width: 400, Height: 400,
		Data: sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
			{ID: "fruit", Value: 30},
			{ID: "veg", Value: 20},
		}},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
