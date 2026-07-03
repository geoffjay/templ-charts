package bar_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleBar renders a small stacked bar chart to an SVG string with the
// charts/render helper. Every chart in the library follows this same shape:
// build a props value, pass the component to render.String.
func ExampleBar() {
	svg, err := render.String(bar.Bar(bar.BarProps{
		Width:   400,
		Height:  300,
		IndexBy: "country",
		Keys:    []string{"hot dogs", "burgers"},
		Data: []bar.BarDatum{
			{"country": "USA", "hot dogs": 42.0, "burgers": 33.0},
			{"country": "Japan", "hot dogs": 28.0, "burgers": 51.0},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
