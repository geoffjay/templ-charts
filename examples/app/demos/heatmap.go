package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/legends"
)

// HeatmapDemo is one heatmap tile on the /heatmap page. Heatmaps render
// statically (no HTMX) in v2 — interactivity arrives with the Phase 5 client
// layer.
type HeatmapDemo struct {
	ID          string
	Title       string
	Description string
	Props       heatmap.HeatMapProps
}

// HeatmapDemos returns the heatmap demos for the /heatmap page.
func HeatmapDemos() []HeatmapDemo {
	margin := core.Margin{Top: 60, Right: 90, Bottom: 30, Left: 90}
	return []HeatmapDemo{
		{
			ID:          "heatmap-sequential",
			Title:       "Sequential",
			Description: "Default sequential color scale (brown→blue-green) with value labels.",
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin: margin, Data: heatmapData(),
			},
		},
		{
			ID:          "heatmap-diverging",
			Title:       "Diverging + border",
			Description: "Diverging red-yellow-blue scale, cell borders, and a continuous legend.",
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      margin,
				Data:        heatmapData(),
				Colors:      heatmap.HeatMapColorConfig{Type: "diverging", Scheme: "red_yellow_blue"},
				BorderWidth: 1,
				Legends: []heatmap.HeatMapLegend{
					{Anchor: legends.LegendAnchorBottom, TranslateY: 30, Length: 240, Thickness: 12, Title: "Value"},
				},
			},
		},
		{
			ID:          "heatmap-no-labels",
			Title:       "Labels off, squared cells",
			Description: "EnableLabels=false and ForceSquare=true for a compact square grid.",
			Props: heatmap.HeatMapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:       margin,
				Data:         heatmapData(),
				EnableLabels: heatmap.BoolPtr(false),
				ForceSquare:  true,
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
