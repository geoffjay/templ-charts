package bump_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleBump renders a small ranking-over-time bump chart to an SVG string
// with the charts/render helper.
func ExampleBump() {
	svg, err := render.String(bump.Bump(bump.BumpProps{
		Width: 400, Height: 300,
		Data: []bump.BumpSerie{
			{ID: "React", Data: []bump.BumpDatum{
				{X: "2020", Y: 1}, {X: "2021", Y: 2}, {X: "2022", Y: 1},
			}},
			{ID: "Vue", Data: []bump.BumpDatum{
				{X: "2020", Y: 2}, {X: "2021", Y: 1}, {X: "2022", Y: 2},
			}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
