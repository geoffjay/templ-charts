package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/samples"
)

// LineDemos returns the set of line chart demos for the /line page.
func LineDemos() []Demo {
	return []Demo{
		{
			ID:          "line-single",
			Title:       "Single series",
			Description: "One line series over 7 years.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: defaultMargin(),
				Data:   []line.LineSeries{lineData()[0]},
			},
		},
		{
			ID:          "line-multi",
			Title:       "Multiple series + legend",
			Description: "Three series with a toggle legend (HTMX).",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: defaultMargin(),
				Curve:  core.CurveMonotoneX,
				Data:   lineData(),
				Legends: []legends.LegendProps{
					{
						Anchor:    legends.LegendAnchorTopRight,
						Direction: legends.LegendDirectionColumn,
						ItemWidth: 100, ItemHeight: 20,
						TranslateX: 10,
					},
				},
			},
		},
		{
			ID:          "line-area-points",
			Title:       "Area + points",
			Description: "Filled area under each line with point dots.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:        commonChartWidth,
				Height:       commonChartHeight,
				Margin:       defaultMargin(),
				Curve:        core.CurveMonotoneX,
				EnableArea:   true,
				AreaOpacity:  0.2,
				EnablePoints: true,
				Data:         lineData(),
			},
		},
		{
			ID:          "line-slices",
			Title:       "Slices + slice tooltip",
			Description: "Hover a vertical slice for a per-series tooltip — resolved client-side (charts/interact), no server round-trip.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:           commonChartWidth,
				Height:          commonChartHeight,
				Margin:          defaultMargin(),
				Curve:           core.CurveMonotoneX,
				EnablePoints:    true,
				EnableSlices:    line.EnableSlicesX,
				EnableCrosshair: true,
				ClientHover:     true,
				Data:            lineData(),
			},
		},
		{
			ID:          "line-mesh",
			Title:       "Mesh hover + crosshair",
			Description: "useMesh=true with the client interactivity layer: nearest-point tooltip + crosshair tracked in the browser (resolves the per-mousemove round-trip).",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:           commonChartWidth,
				Height:          commonChartHeight,
				Margin:          defaultMargin(),
				Curve:           core.CurveMonotoneX,
				UseMesh:         true,
				EnableCrosshair: true,
				ClientHover:     true,
				Data:            lineData(),
			},
		},
	}
}

// lineData is the shared line dataset, sourced from the public samples package.
func lineData() []line.LineSeries {
	return samples.Line()
}
