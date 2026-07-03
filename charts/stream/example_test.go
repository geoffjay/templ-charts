package stream_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/stream"
)

// ExampleStream renders a small stacked stream chart to an SVG string with the
// charts/render helper.
func ExampleStream() {
	keys := []string{"Raoul", "Josiane", "Marcel"}
	svg, err := render.String(stream.Stream(stream.StreamProps{
		Width: 400, Height: 300,
		Keys: keys,
		Data: []stream.StreamDatum{
			{"Raoul": 10, "Josiane": 20, "Marcel": 30},
			{"Raoul": 15, "Josiane": 18, "Marcel": 25},
			{"Raoul": 12, "Josiane": 22, "Marcel": 28},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
