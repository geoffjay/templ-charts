package demos

import (
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/core"
)

// CirclePackingDemo is one tile on the /circle-packing page (static render).
type CirclePackingDemo struct {
	ID          string
	Title       string
	Description string
	Props       cp.CirclePackingProps
}

func circlePackingSample() cp.CirclePackingNode {
	return cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
		{ID: "A", Children: []cp.CirclePackingNode{
			{ID: "a1", Value: 20}, {ID: "a2", Value: 12}, {ID: "a3", Value: 8},
		}},
		{ID: "B", Children: []cp.CirclePackingNode{
			{ID: "b1", Value: 16}, {ID: "b2", Value: 24},
		}},
		{ID: "C", Children: []cp.CirclePackingNode{
			{ID: "c1", Value: 10}, {ID: "c2", Value: 6}, {ID: "c3", Value: 14}, {ID: "c4", Value: 5},
		}},
	}}
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
				EnableLabels: cp.BoolPtr(true),
			},
		},
	}
}
