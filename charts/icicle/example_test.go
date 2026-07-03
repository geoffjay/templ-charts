package icicle_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleIcicle renders a small icicle chart to an SVG string with the
// charts/render helper.
func ExampleIcicle() {
	svg, err := render.String(icicle.Icicle(icicle.IcicleProps{
		Width: 400, Height: 300,
		Data: icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
			{ID: "analytics", Value: 40},
			{ID: "billing", Value: 22},
		}},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
