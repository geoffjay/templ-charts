package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/polarbar"
)

// PolarBarDemo is one polar-bar tile on the /polar-bar page (static render).
type PolarBarDemo struct {
	ID          string
	Title       string
	Description string
	Props       polarbar.PolarBarProps
}

// PolarBarDemos returns the polar-bar demos for the /polar-bar page.
func PolarBarDemos() []PolarBarDemo {
	data := []polarbar.PolarBarDatum{
		{Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3}},
		{Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4}},
		{Index: "Wed", Values: map[string]float64{"walk": 10, "bus": 4, "bike": 6}},
		{Index: "Thu", Values: map[string]float64{"walk": 7, "bus": 9, "bike": 2}},
		{Index: "Fri", Values: map[string]float64{"walk": 12, "bus": 3, "bike": 8}},
	}
	return []PolarBarDemo{
		{
			ID:          "polarbar-basic",
			Title:       "Stacked over a full circle",
			Description: "Five indices as angular bands over 0–360°; three keys stacked radially with radial + circular grids and index labels. Hover an arc for its value.",
			Props: polarbar.PolarBarProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:      core.Margin{Top: 40, Right: 60, Bottom: 40, Left: 60},
				Data:        data,
				Keys:        []string{"walk", "bus", "bike"},
				Interactive: true,
			},
		},
		{
			ID:          "polarbar-donut-labels",
			Title:       "Inner radius + arc labels + legend",
			Description: "innerRadius=0.3, cornerRadius=3, arc labels on, plus a key legend.",
			Props: polarbar.PolarBarProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:          core.Margin{Top: 40, Right: 160, Bottom: 40, Left: 60},
				Data:            data,
				Keys:            []string{"walk", "bus", "bike"},
				InnerRadius:     0.3,
				CornerRadius:    3,
				Padding:         0.1,
				EnableArcLabels: polarbar.BoolPtr(true),
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 120},
				},
			},
		},
	}
}
