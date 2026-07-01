package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/sunburst"
)

// SunburstDemo is one sunburst tile on the /sunburst page (static render).
type SunburstDemo struct {
	ID          string
	Title       string
	Description string
	Props       sunburst.SunburstProps
}

func sunburstSample() sunburst.SunburstNode {
	return sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
		{ID: "fruit", Children: []sunburst.SunburstNode{
			{ID: "apple", Value: 30}, {ID: "pear", Value: 18}, {ID: "grape", Value: 22},
		}},
		{ID: "veg", Children: []sunburst.SunburstNode{
			{ID: "carrot", Value: 20}, {ID: "pea", Value: 14},
		}},
		{ID: "grain", Value: 26},
	}}
}

// SunburstDemos returns the sunburst demos for the /sunburst page.
func SunburstDemos() []SunburstDemo {
	return []SunburstDemo{
		{
			ID:          "sunburst-basic",
			Title:       "Nested value breakdown",
			Description: "A hierarchy mapped to arcs via the partition layout on [2π, r²]; colors inherit from the top-level ancestor. Hover an arc for its value.",
			Props: sunburst.SunburstProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:      core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:        sunburstSample(),
				Interactive: true,
			},
		},
		{
			ID:          "sunburst-labels",
			Title:       "Arc labels + corner radius",
			Description: "enableArcLabels with a small cornerRadius and skip-angle.",
			Props: sunburst.SunburstProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:             core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:               sunburstSample(),
				CornerRadius:       2,
				EnableArcLabels:    sunburst.BoolPtr(true),
				ArcLabelsSkipAngle: 10,
			},
		},
	}
}
