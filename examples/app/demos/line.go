package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/scales"
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
				EnableArea:   core.BoolPtr(true),
				AreaOpacity:  0.2,
				EnablePoints: core.BoolPtr(true),
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
				EnablePoints:    core.BoolPtr(true),
				EnableSlices:    line.EnableSlicesX,
				EnableCrosshair: core.BoolPtr(true),
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
				EnableCrosshair: core.BoolPtr(true),
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

// LineScaleDemo returns the value-scale demo card (id "line-scale"): two
// exponentially growing series (~×3/year), rendered on a linear or log Y axis.
// The page toggles logScale so the same data can be compared under both scales.
func LineScaleDemo(logScale bool) Demo {
	p := line.LineProps{
		Width:  commonChartWidth,
		Height: commonChartHeight,
		Margin: defaultMargin(),
		Data: []line.LineSeries{
			{ID: "signups", Data: []line.LinePointData{
				{X: "2016", Y: 120.0},
				{X: "2017", Y: 640.0},
				{X: "2018", Y: 3100.0},
				{X: "2019", Y: 9800.0},
				{X: "2020", Y: 26000.0},
				{X: "2021", Y: 74000.0},
				{X: "2022", Y: 210000.0},
				{X: "2023", Y: 560000.0},
				{X: "2024", Y: 1400000.0},
			}},
			{ID: "paying", Data: []line.LinePointData{
				{X: "2016", Y: 6.0},
				{X: "2017", Y: 34.0},
				{X: "2018", Y: 190.0},
				{X: "2019", Y: 720.0},
				{X: "2020", Y: 2400.0},
				{X: "2021", Y: 8100.0},
				{X: "2022", Y: 26000.0},
				{X: "2023", Y: 71000.0},
				{X: "2024", Y: 190000.0},
			}},
		},
	}
	desc := "Two series growing ~3×/year. On a linear axis the early years flatten to zero; on a log axis constant growth reads as a straight line."
	if logScale {
		p.YScale = scales.ScaleLogSpec{Base: 10, Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	}
	return Demo{ID: "line-scale", Title: "Value scale (linear / log)", Description: desc, Kind: htmx.KindLine, Props: p}
}
