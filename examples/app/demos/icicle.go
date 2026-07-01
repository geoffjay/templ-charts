package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/icicle"
)

// IcicleDemo is one icicle tile on the /icicle page (static render).
type IcicleDemo struct {
	ID          string
	Title       string
	Description string
	Props       icicle.IcicleProps
}

func icicleSample() icicle.IcicleNode {
	return icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
		{ID: "analytics", Children: []icicle.IcicleNode{
			{ID: "charts", Value: 40}, {ID: "reports", Value: 22},
		}},
		{ID: "billing", Children: []icicle.IcicleNode{
			{ID: "invoices", Value: 18}, {ID: "refunds", Value: 9}, {ID: "plans", Value: 15},
		}},
		{ID: "auth", Value: 28},
	}}
}

// IcicleDemos returns the icicle demos for the /icicle page.
func IcicleDemos() []IcicleDemo {
	return []IcicleDemo{
		{
			ID:          "icicle-bottom",
			Title:       "Depth bands (bottom)",
			Description: "The partition layout as depth-banded rectangles growing downward, with per-node hover.",
			Props: icicle.IcicleProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:        icicleSample(),
				Interactive: true,
			},
		},
		{
			ID:          "icicle-right-labels",
			Title:       "Right orientation + labels",
			Description: "orientation=right (depth grows across the width) with node id labels.",
			Props: icicle.IcicleProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:       core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:         icicleSample(),
				Orientation:  icicle.OrientationRight,
				EnableLabels: icicle.BoolPtr(true),
			},
		},
	}
}
