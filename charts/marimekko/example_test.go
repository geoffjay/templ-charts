package marimekko_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/marimekko"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleMarimekko renders a small variable-width stacked bar chart to an SVG
// string with the charts/render helper.
func ExampleMarimekko() {
	svg, err := render.String(marimekko.Marimekko(marimekko.MarimekkoProps{
		Width: 400, Height: 300,
		Dimensions: []marimekko.MarimekkoDimension{
			{ID: "agree", Key: "agree"},
			{ID: "disagree", Key: "disagree"},
		},
		Data: []marimekko.MarimekkoDatum{
			{ID: "France", Value: 42, Dimensions: map[string]float64{"agree": 30, "disagree": 12}},
			{ID: "USA", Value: 63, Dimensions: map[string]float64{"agree": 48, "disagree": 15}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
