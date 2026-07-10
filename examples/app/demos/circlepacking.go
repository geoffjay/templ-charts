package demos

import (
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
)

// CirclePackingDemo is one tile on the /circle-packing page (static render).
type CirclePackingDemo struct {
	ID          string
	Title       string
	Description string
	Props       cp.CirclePackingProps
}

// circlePackingSample is the shared hierarchy, sourced from the public samples
// package.
func circlePackingSample() cp.CirclePackingNode {
	return samples.CirclePacking()
}

// CirclePackingDemos returns the demos for the /circle-packing page.
func CirclePackingDemos() []CirclePackingDemo {
	return []CirclePackingDemo{
		{
			ID:          "circlepacking-basic",
			Title:       "Nested packed circles",
			Description: "Welzl enclosing-circle packing (deterministic via the ported LCG), colored by depth, with per-node hover.",
			Props: cp.CirclePackingProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:      core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:        circlePackingSample(),
				Interactive: true,
			},
		},
		{
			ID:          "circlepacking-labels",
			Title:       "Padding + leaf labels",
			Description: "padding between circles, borders, and leaf id labels.",
			Props: cp.CirclePackingProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:       core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:         circlePackingSample(),
				Padding:      3,
				BorderWidth:  1,
				EnableLabels: core.BoolPtr(true),
			},
		},
	}
}
