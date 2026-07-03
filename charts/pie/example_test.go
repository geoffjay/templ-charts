package pie_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExamplePie renders a small pie chart to an SVG string with the charts/render
// helper.
func ExamplePie() {
	svg, err := render.String(pie.Pie(pie.PieProps{
		Width:  400,
		Height: 300,
		Data: []any{
			map[string]any{"id": "Go", "value": 40.0},
			map[string]any{"id": "Rust", "value": 25.0},
			map[string]any{"id": "Python", "value": 35.0},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
