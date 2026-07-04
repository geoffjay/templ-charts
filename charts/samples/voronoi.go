package samples

import "github.com/geoffjay/templ-charts/charts/voronoi"

// Voronoi returns a fixed scatter of 12 points in the unit square (the default
// x/y domains), chosen to give visibly distinct cells. Ready to assign to
// voronoi.VoronoiProps.Data. Mirrors the demo's voronoiSample.
func Voronoi() []voronoi.VoronoiDatum {
	return []voronoi.VoronoiDatum{
		{ID: "0", X: 0.12, Y: 0.18},
		{ID: "1", X: 0.42, Y: 0.10},
		{ID: "2", X: 0.78, Y: 0.22},
		{ID: "3", X: 0.90, Y: 0.55},
		{ID: "4", X: 0.66, Y: 0.48},
		{ID: "5", X: 0.30, Y: 0.40},
		{ID: "6", X: 0.08, Y: 0.62},
		{ID: "7", X: 0.24, Y: 0.86},
		{ID: "8", X: 0.55, Y: 0.78},
		{ID: "9", X: 0.84, Y: 0.90},
		{ID: "10", X: 0.48, Y: 0.55},
		{ID: "11", X: 0.70, Y: 0.15},
	}
}
