package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/marimekko"
)

// MarimekkoDemo is one marimekko tile on the /marimekko page (static render).
type MarimekkoDemo struct {
	ID          string
	Title       string
	Description string
	Props       marimekko.MarimekkoProps
}

// MarimekkoDemos returns the marimekko demos for the /marimekko page.
func MarimekkoDemos() []MarimekkoDemo {
	dims := []marimekko.MarimekkoDimension{
		{ID: "agree strongly", Key: "agreeStrongly"},
		{ID: "agree", Key: "agree"},
		{ID: "disagree", Key: "disagree"},
		{ID: "disagree strongly", Key: "disagreeStrongly"},
	}
	data := []marimekko.MarimekkoDatum{
		{ID: "France", Value: 42, Dimensions: map[string]float64{"agreeStrongly": 18, "agree": 12, "disagree": 8, "disagreeStrongly": 4}},
		{ID: "Japan", Value: 28, Dimensions: map[string]float64{"agreeStrongly": 6, "agree": 10, "disagree": 8, "disagreeStrongly": 4}},
		{ID: "USA", Value: 63, Dimensions: map[string]float64{"agreeStrongly": 30, "agree": 18, "disagree": 10, "disagreeStrongly": 5}},
		{ID: "Germany", Value: 35, Dimensions: map[string]float64{"agreeStrongly": 12, "agree": 14, "disagree": 6, "disagreeStrongly": 3}},
	}
	return []MarimekkoDemo{
		{
			ID:          "marimekko-basic",
			Title:       "Variable-width stacked bars",
			Description: "Bar width is proportional to sample size; each bar is split by response via d3.Stack (offset none). Hover a segment for its value.",
			Props: marimekko.MarimekkoProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 20, Right: 160, Bottom: 50, Left: 60},
				Data:        data,
				Dimensions:  dims,
				Interactive: true,
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 130},
				},
			},
		},
		{
			ID:          "marimekko-expand",
			Title:       "Normalized (expand offset)",
			Description: "offset=expand normalizes every bar to full height, so each column shows response share; widths still encode sample size.",
			Props: marimekko.MarimekkoProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 20, Right: 160, Bottom: 50, Left: 60},
				Data:        data,
				Dimensions:  dims,
				Offset:      marimekko.OffsetExpand,
				BorderWidth: 1,
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 130},
				},
			},
		},
	}
}
