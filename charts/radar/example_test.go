package radar_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/radar"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleRadar renders a small radar chart to an SVG string with the
// charts/render helper.
func ExampleRadar() {
	svg, err := render.String(radar.Radar(radar.RadarProps{
		Width:   400,
		Height:  400,
		IndexBy: "taste",
		Keys:    []string{"chardonnay", "syrah"},
		Data: []map[string]any{
			{"taste": "fruity", "chardonnay": 93.0, "syrah": 114.0},
			{"taste": "bitter", "chardonnay": 91.0, "syrah": 72.0},
			{"taste": "heavy", "chardonnay": 56.0, "syrah": 99.0},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
