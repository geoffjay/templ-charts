package voronoi_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/voronoi"
)

// ExampleVoronoi renders a small Voronoi diagram from a handful of points to an
// SVG string with the charts/render helper.
func ExampleVoronoi() {
	svg, err := render.String(voronoi.Voronoi(voronoi.VoronoiProps{
		Width: 400, Height: 400,
		Data: []voronoi.VoronoiDatum{
			{ID: "0", X: 0.12, Y: 0.18},
			{ID: "1", X: 0.42, Y: 0.10},
			{ID: "2", X: 0.78, Y: 0.22},
			{ID: "3", X: 0.55, Y: 0.78},
			{ID: "4", X: 0.30, Y: 0.55},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
