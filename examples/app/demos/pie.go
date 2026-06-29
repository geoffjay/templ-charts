package demos

import (
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/pie"
)

// PieDemos returns the set of pie chart demos for the /pie page.
func PieDemos() []Demo {
	pieMargin := defaultMargin()
	return []Demo{
		{
			ID:          "pie-plain",
			Title:       "Plain full pie",
			Description: "Default settings: full pie, arc + arc-link labels.",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:                   commonChartWidth,
				Height:                  commonChartHeight,
				Margin:                  pieMargin,
				Data:                    pieData(),
				ActiveOuterRadiusOffset: 6,
			},
		},
		{
			ID:          "pie-donut",
			Title:       "Donut",
			Description: "innerRadius=0.5 turns the pie into a donut.",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:                   commonChartWidth,
				Height:                  commonChartHeight,
				Margin:                  pieMargin,
				InnerRadius:             0.5,
				PadAngle:                0.5,
				CornerRadius:            5,
				Data:                    pieData(),
				ActiveOuterRadiusOffset: 8,
			},
		},
		{
			ID:          "pie-half",
			Title:       "Half pie",
			Description: "startAngle=0, endAngle=180 with fit rescaling.",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:      commonChartWidth,
				Height:     commonChartHeight,
				Margin:     pieMargin,
				StartAngle: 0,
				EndAngle:   180,
				Fit:        true,
				Data:       pieData(),
			},
		},
		{
			ID:          "pie-sorted",
			Title:       "Sorted by value",
			Description: "sortByValue=true orders slices largest-first.",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:       commonChartWidth,
				Height:      commonChartHeight,
				Margin:      pieMargin,
				SortByValue: true,
				InnerRadius: 0.5,
				PadAngle:    0.5,
				Data:        pieData(),
			},
		},
		{
			ID:          "pie-active",
			Title:       "Active-arc highlight",
			Description: "Hover an arc to pop it out (HTMX sets ActiveID).",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:                   commonChartWidth,
				Height:                  commonChartHeight,
				Margin:                  pieMargin,
				InnerRadius:             0.5,
				PadAngle:                0.5,
				CornerRadius:            3,
				ActiveOuterRadiusOffset: 8,
				Data:                    pieData(),
			},
		},
		{
			ID:          "pie-legend-toggle",
			Title:       "Legend + series toggle",
			Description: "Click a legend item to hide/show a slice (HTMX).",
			Kind:        htmx.KindPie,
			Props: pie.PieProps{
				Width:       commonChartWidth,
				Height:      commonChartHeight,
				Margin:      pieMargin,
				InnerRadius: 0.5,
				PadAngle:    0.5,
				Data:        pieData(),
				Legends: []legends.LegendProps{
					{
						Anchor:    legends.LegendAnchorRight,
						Direction: legends.LegendDirectionColumn,
						ItemWidth: 120, ItemHeight: 20,
						TranslateX: 10,
					},
				},
			},
		},
	}
}

// pieData is the shared programming-language dataset for the pie demos.
func pieData() []any {
	labels := []string{"Go", "Rust", "TypeScript", "Python", "Ruby", "Elixir", "C", "Zig", "Kotlin"}
	data := make([]any, len(labels))
	for i, l := range labels {
		data[i] = map[string]any{
			"id":    l,
			"value": float64((i*17+3)%80 + 10),
		}
	}
	return data
}
