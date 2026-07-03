package heatmap_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleHeatMap renders a small heatmap to an SVG string with the
// charts/render helper.
func ExampleHeatMap() {
	v := func(f float64) *float64 { return &f }
	svg, err := render.String(heatmap.HeatMap(heatmap.HeatMapProps{
		Width:  400,
		Height: 300,
		Data: []heatmap.HeatMapSerie{
			{ID: "Japan", Data: []heatmap.HeatMapDatum{
				{X: "Train", Y: v(42)}, {X: "Bus", Y: v(18)},
			}},
			{ID: "France", Data: []heatmap.HeatMapDatum{
				{X: "Train", Y: v(33)}, {X: "Bus", Y: v(51)},
			}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
