package circlepacking_test

import (
	"fmt"
	"log"
	"strings"

	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleCirclePacking renders a small circle-packing chart to an SVG string
// with the charts/render helper.
func ExampleCirclePacking() {
	svg, err := render.String(cp.CirclePacking(cp.CirclePackingProps{
		Width: 400, Height: 400,
		Data: cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
			{ID: "A", Value: 20},
			{ID: "B", Value: 24},
		}},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
