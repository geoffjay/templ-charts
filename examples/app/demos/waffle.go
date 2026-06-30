package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/waffle"
)

// WaffleDemo is one waffle tile on the /waffle page (static render).
type WaffleDemo struct {
	ID          string
	Title       string
	Description string
	Props       waffle.WaffleProps
}

// WaffleDemos returns the waffle demos for the /waffle page.
func WaffleDemos() []WaffleDemo {
	data := []waffle.WaffleDatum{
		{ID: "men", Label: "Men", Value: 32},
		{ID: "women", Label: "Women", Value: 41},
		{ID: "children", Label: "Children", Value: 19},
	}
	return []WaffleDemo{
		{
			ID:          "waffle-basic",
			Title:       "Part of whole",
			Description: "100-cell grid (total=100); each cell is one unit. Remaining cells are empty.",
			Props: waffle.WaffleProps{
				Width: commonChartHeight, Height: commonChartHeight,
				Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Total:  100, Rows: 10, Columns: 10, Data: data,
			},
		},
		{
			ID:          "waffle-legend",
			Title:       "Fill from bottom + legend",
			Description: "fillDirection=bottom, cell borders, and a legend built from the data.",
			Props: waffle.WaffleProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:        core.Margin{Top: 10, Right: 140, Bottom: 10, Left: 10},
				Total:         100,
				Rows:          10,
				Columns:       10,
				Data:          data,
				Interactive:   true,
				FillDirection: grid.GridFillBottom,
				BorderWidth:   2,
				Legends: []waffle.WaffleLegend{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 120},
				},
			},
		},
	}
}
