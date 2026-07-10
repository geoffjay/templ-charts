package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/samples"
)

// IcicleDemo is one icicle tile on the /icicle page (static render).
type IcicleDemo struct {
	ID          string
	Title       string
	Description string
	Props       icicle.IcicleProps
}

// icicleSample is the shared hierarchy, sourced from the public samples package.
func icicleSample() icicle.IcicleNode {
	return samples.Icicle()
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
				EnableLabels: core.BoolPtr(true),
			},
		},
	}
}
