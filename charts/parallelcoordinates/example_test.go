package parallelcoordinates_test

import (
	"fmt"
	"log"
	"strings"

	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleParallelCoordinates renders a small parallel-coordinates chart to an
// SVG string with the charts/render helper.
func ExampleParallelCoordinates() {
	svg, err := render.String(pc.ParallelCoordinates(pc.PCProps{
		Width: 400, Height: 300,
		Variables: []pc.PCVariable{
			{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
			{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
			{Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
		},
		Data: []pc.PCDatum{
			{ID: "batch A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "grade": "B"}},
			{ID: "batch B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "grade": "A"}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
