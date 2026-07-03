package treemap_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// ExampleTreemap renders a small treemap to an SVG string with the
// charts/render helper.
func ExampleTreemap() {
	svg, err := render.String(treemap.Treemap(treemap.TreemapProps{
		Width: 400, Height: 300,
		Data: treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
			{ID: "viz", Value: 33},
			{ID: "colors", Value: 41},
		}},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
