package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/radar"
)

// RadarDemo is one radar tile on the /radar page (static render).
type RadarDemo struct {
	ID          string
	Title       string
	Description string
	Props       radar.RadarProps
}

// RadarDemos returns the radar demos for the /radar page.
func RadarDemos() []RadarDemo {
	data := []map[string]any{
		{"taste": "fruity", "chardonnay": 93.0, "carmenere": 61.0, "syrah": 114.0},
		{"taste": "bitter", "chardonnay": 91.0, "carmenere": 37.0, "syrah": 72.0},
		{"taste": "heavy", "chardonnay": 56.0, "carmenere": 95.0, "syrah": 99.0},
		{"taste": "strong", "chardonnay": 64.0, "carmenere": 90.0, "syrah": 30.0},
		{"taste": "sunny", "chardonnay": 119.0, "carmenere": 94.0, "syrah": 103.0},
	}
	keys := []string{"chardonnay", "carmenere", "syrah"}
	return []RadarDemo{
		{
			ID:          "radar-circular",
			Title:       "Circular grid + dots",
			Description: "Three wine varieties scored on five taste axes; default circular grid levels, filled areas, and per-point dots.",
			Props: radar.RadarProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin:      core.Margin{Top: 70, Right: 90, Bottom: 50, Left: 90},
				Data:        data,
				Keys:        keys,
				IndexBy:     "taste",
				Interactive: true,
			},
		},
		{
			ID:          "radar-linear-legend",
			Title:       "Polygon grid + legend",
			Description: "gridShape=linear (polygon rings), a 30° rotation, and a legend built from the keys.",
			Props: radar.RadarProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:    core.Margin{Top: 70, Right: 160, Bottom: 50, Left: 90},
				Data:      data,
				Keys:      keys,
				IndexBy:   "taste",
				GridShape: radar.GridShapeLinear,
				Rotation:  30,
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 110},
				},
			},
		},
	}
}
