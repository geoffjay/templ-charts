package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/radialbar"
)

// RadialBarDemo is one radial-bar tile on the /radial-bar page (static render).
type RadialBarDemo struct {
	ID          string
	Title       string
	Description string
	Props       radialbar.RadialBarProps
}

// RadialBarDemos returns the radial-bar demos for the /radial-bar page.
func RadialBarDemos() []RadialBarDemo {
	data := []radialbar.RadialBarSerie{
		{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18}, {X: "Meat", Y: 12},
		}},
		{ID: "Combini", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 12}, {X: "Fruits", Y: 9}, {X: "Meat", Y: 7},
		}},
		{ID: "Online", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 8}, {X: "Fruits", Y: 15}, {X: "Meat", Y: 4},
		}},
	}
	return []RadialBarDemo{
		{
			ID:          "radialbar-basic",
			Title:       "Stacked rings + tracks",
			Description: "Three sources, three categories each, stacked over a 0–270° angle scale with background tracks, a radial axis, and the outer circular axis.",
			Props: radialbar.RadialBarProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin: core.Margin{Top: 40, Right: 60, Bottom: 40, Left: 60},
				Data:   data,
			},
		},
		{
			ID:          "radialbar-labels-legend",
			Title:       "Labels + rounded corners + legend",
			Description: "enableLabels=true, cornerRadius=4, and a category legend.",
			Props: radialbar.RadialBarProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:       core.Margin{Top: 40, Right: 180, Bottom: 40, Left: 60},
				Data:         data,
				CornerRadius: 4,
				EnableLabels: radialbar.BoolPtr(true),
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 140},
				},
			},
		},
	}
}
