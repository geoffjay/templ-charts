package sankey_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/sankey"
)

// ExampleSankey renders a small flow diagram to an SVG string with the
// charts/render helper.
func ExampleSankey() {
	svg, err := render.String(sankey.Sankey(sankey.SankeyProps{
		Width: 700, Height: 420,
		Nodes: []sankey.SankeyInputNode{
			{ID: "Coal"},
			{ID: "Grid"},
			{ID: "Homes"},
		},
		Links: []sankey.SankeyInputLink{
			{Source: "Coal", Target: "Grid", Value: 12},
			{Source: "Grid", Target: "Homes", Value: 12},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
