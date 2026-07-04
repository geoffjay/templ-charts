package samples

import "github.com/geoffjay/templ-charts/charts/geo"

// Choropleth returns a small, self-contained synthetic world: a 3×2 grid of
// rectangular polygon "regions" (ids A–F) spanning lon [-90,90] × lat [-40,40],
// plus per-region values bound by id. It is deliberately independent of the
// bundled world-countries GeoJSON so the samples package stays lightweight.
//
// The value for region "F" is intentionally omitted so a choropleth rendered
// from this data exercises its unknownColor path (as the demo does for
// Antarctica/Greenland). Assign the features to geo.ChoroplethProps.Features
// and the data to geo.ChoroplethProps.Data.
func Choropleth() ([]geo.Feature, []geo.ChoroplethDatum) {
	// Grid bands: 3 longitude columns × 2 latitude rows.
	lonEdges := []float64{-90, -30, 30, 90}
	latEdges := []float64{-40, 0, 40}
	ids := [][]string{
		{"A", "B", "C"}, // bottom row (lat -40..0)
		{"D", "E", "F"}, // top row    (lat 0..40)
	}

	var features []geo.Feature
	for r := 0; r < len(latEdges)-1; r++ {
		lat0, lat1 := latEdges[r], latEdges[r+1]
		for c := 0; c < len(lonEdges)-1; c++ {
			lon0, lon1 := lonEdges[c], lonEdges[c+1]
			// A single closed outer ring, [lon,lat] pairs.
			ring := [][2]float64{
				{lon0, lat0}, {lon1, lat0}, {lon1, lat1}, {lon0, lat1}, {lon0, lat0},
			}
			features = append(features, geo.Feature{
				Type: "Feature",
				ID:   ids[r][c],
				Geometry: geo.Geometry{
					Type:        geo.TypePolygon,
					Coordinates: [][][2]float64{ring},
				},
			})
		}
	}

	// Deterministic values per region; "F" is omitted on purpose.
	data := []geo.ChoroplethDatum{
		{ID: "A", Value: 12},
		{ID: "B", Value: 34},
		{ID: "C", Value: 56},
		{ID: "D", Value: 78},
		{ID: "E", Value: 91},
	}
	return features, data
}
