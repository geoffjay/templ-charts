package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
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
			Description: "Hover a vertical slice to see a per-series tooltip (HTMX).",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:        commonChartWidth,
				Height:       commonChartHeight,
				Margin:       defaultMargin(),
				Curve:        core.CurveMonotoneX,
				EnablePoints: true,
				EnableSlices: line.EnableSlicesX,
				Data:         lineData(),
			},
		},
		{
			ID:          "line-mesh",
			Title:       "Mesh hover",
			Description: "useMesh=true captures hover across the full plot area.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				Margin:  defaultMargin(),
				Curve:   core.CurveMonotoneX,
				UseMesh: true,
				Data:    lineData(),
			},
		},
	}
}

// lineData is the shared 3 countries × 7 years dataset for the line demos.
func lineData() []line.LineSeries {
	years := []string{"2018", "2019", "2020", "2021", "2022", "2023", "2024"}
	countries := []string{"USA", "Germany", "France"}
	series := make([]line.LineSeries, len(countries))
	for i, c := range countries {
		data := make([]line.LinePointData, len(years))
		for j, y := range years {
			data[j] = line.LinePointData{
				X: y,
				Y: float64((i*3+j*11)%90 + 20),
			}
		}
		series[i] = line.LineSeries{ID: c, Data: data}
	}
	return series
}
