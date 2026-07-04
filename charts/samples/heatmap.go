package samples

import "github.com/geoffjay/templ-charts/charts/heatmap"

// Heatmap returns a deterministic 5-country × 6-metric grid, ready to assign to
// HeatMapProps.Data.
func Heatmap() []heatmap.HeatMapSerie {
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
