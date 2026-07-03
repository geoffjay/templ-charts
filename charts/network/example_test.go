package network_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/network"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleNetwork renders a small force-directed node/link graph to an SVG
// string with the charts/render helper.
func ExampleNetwork() {
	svg, err := render.String(network.Network(network.NetworkProps{
		Width: 400, Height: 400,
		Nodes: []network.NetworkInputNode{
			{ID: "hub-0", Size: 20, Color: "#e8c1a0"},
			{ID: "leaf-0", Size: 10, Color: "#f47560"},
			{ID: "leaf-1", Size: 10, Color: "#f1e15b"},
		},
		Links: []network.NetworkInputLink{
			{Source: "hub-0", Target: "leaf-0"},
			{Source: "hub-0", Target: "leaf-1"},
		},
		LinkDistance: 90,
		Repulsivity:  120,
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
