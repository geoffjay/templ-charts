package funnel_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/funnel"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleFunnel renders a small conversion funnel to an SVG string with the
// charts/render helper.
func ExampleFunnel() {
	svg, err := render.String(funnel.Funnel(funnel.FunnelProps{
		Width: 400, Height: 400,
		Data: []funnel.FunnelDatum{
			{ID: "sent", Label: "Sent", Value: 60000},
			{ID: "viewed", Label: "Viewed", Value: 38000},
			{ID: "clicked", Label: "Clicked", Value: 22000},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
