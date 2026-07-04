package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// HeatmapDemos returns the heatmap demos for the /heatmap page. They are
// registered with the htmx registry so cell hover shows a tooltip + active
// highlight, exactly like the bar/line/pie demos.
func HeatmapDemos() []Demo {
	margin := core.Margin{Top: 60, Right: 90, Bottom: 30, Left: 90}
	return []Demo{
		{
			ID:          "heatmap-sequential",
			Title:       "Sequential (hover a cell)",
			Description: "Default sequential color scale with value labels. Hover a cell for a tooltip; the rest dim.",
			Kind:        htmx.KindHeatmap,
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin: margin, Data: heatmapData(),
			},
		},
		{
			ID:          "heatmap-diverging",
			Title:       "Diverging + border + legend",
			Description: "Diverging red-yellow-blue scale, cell borders, and a bottom-anchored continuous legend.",
			Kind:        htmx.KindHeatmap,
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 60, Right: 90, Bottom: 80, Left: 90},
				Data:        heatmapData(),
				Colors:      heatmap.HeatMapColorConfig{Type: "diverging", Scheme: "red_yellow_blue"},
				BorderWidth: 1,
				Legends: []heatmap.HeatMapLegend{
					{Anchor: legends.LegendAnchorBottom, TranslateY: 64, Length: 260, Thickness: 12, Title: "Value"},
				},
			},
		},
		{
			ID:          "heatmap-no-labels",
			Title:       "Labels off, squared cells",
			Description: "EnableLabels=false and ForceSquare=true for a compact square grid. Hover still works.",
			Kind:        htmx.KindHeatmap,
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:       margin,
				Data:         heatmapData(),
				EnableLabels: heatmap.BoolPtr(false),
				ForceSquare:  true,
			},
		},
		{
			ID:          "heatmap-canvas",
			Title:       "Canvas backend",
			Description: "The same grid rendered with the Canvas engine (charts/canvas): each cell is a FillRect in a <canvas> draw-list while axes stay SVG — for large-N grids where one <rect> per cell is too many DOM nodes.",
			Kind:        htmx.KindHeatmap,
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:  margin,
				Data:    heatmapData(),
				Render:  theming.EngineCanvas,
				ChartID: "heatmap-canvas",
			},
		},
	}
}

// heatmapData is a deterministic 5-country × 6-metric dataset.
func heatmapData() []heatmap.HeatMapSerie {
	countries := []string{"Japan", "France", "USA", "Germany", "Brazil"}
	metrics := []string{"Train", "Subway", "Bus", "Car", "Bike", "Walk"}
	out := make([]heatmap.HeatMapSerie, len(countries))
	for i, c := range countries {
		data := make([]heatmap.HeatMapDatum, len(metrics))
		for j, m := range metrics {
			v := float64((i*13+j*29)%100 + 5)
			data[j] = heatmap.HeatMapDatum{X: m, Y: &v}
		}
		out[i] = heatmap.HeatMapSerie{ID: c, Data: data}
	}
	return out
}
