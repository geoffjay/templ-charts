package samples

import "github.com/geoffjay/templ-charts/charts/bullet"

// Bullet returns ready-to-render bullet data: three KPI rows, each with
// stacked qualitative ranges, one or more measure bars, and comparative
// markers.
func Bullet() []bullet.BulletItemDatum {
	return []bullet.BulletItemDatum{
		{ID: "temperature", Ranges: []float64{0, 20, 60, 100}, Measures: []float64{55}, Markers: []float64{72}},
		{ID: "power", Ranges: []float64{0, 30, 70, 120}, Measures: []float64{48, 88}, Markers: []float64{95}},
		{ID: "volume", Ranges: []float64{0, 40, 80, 140}, Measures: []float64{110}, Markers: []float64{100}},
	}
}
