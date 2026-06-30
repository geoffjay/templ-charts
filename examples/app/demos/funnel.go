package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/funnel"
)

// FunnelDemo is one funnel tile on the /funnel page.
type FunnelDemo struct {
	ID          string
	Title       string
	Description string
	Props       funnel.FunnelProps
}

// FunnelDemos returns the funnel demos for the /funnel page.
func FunnelDemos() []FunnelDemo {
	data := []funnel.FunnelDatum{
		{ID: "sent", Label: "Sent", Value: 60000},
		{ID: "viewed", Label: "Viewed", Value: 38000},
		{ID: "clicked", Label: "Clicked", Value: 22000},
		{ID: "cart", Label: "Add to Cart", Value: 12000},
		{ID: "purchased", Label: "Purchased", Value: 7000},
	}
	return []FunnelDemo{
		{
			ID:          "funnel-vertical",
			Title:       "Vertical (smooth)",
			Description: "Five funnel parts with smooth (curveBasis) edges, side borders, separators, and centered labels.",
			Props: funnel.FunnelProps{
				Width: commonChartWidth, Height: 460,
				Margin: core.Margin{Top: 20, Right: 30, Bottom: 20, Left: 30},
				Data:   data,
			},
		},
		{
			ID:          "funnel-horizontal",
			Title:       "Horizontal + linear",
			Description: "direction=horizontal with linear interpolation (straight trapezoid edges).",
			Props: funnel.FunnelProps{
				Width: commonChartWidth, Height: 320,
				Margin:        core.Margin{Top: 20, Right: 30, Bottom: 20, Left: 30},
				Data:          data,
				Direction:     funnel.FunnelDirectionHorizontal,
				Interpolation: funnel.FunnelInterpolationLinear,
			},
		},
	}
}
