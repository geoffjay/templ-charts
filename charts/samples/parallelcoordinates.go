package samples

import pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"

// ParallelCoordinates returns ready-to-render parallel-coordinates data: four
// batches measured across four linear variables plus one categorical (point)
// variable. The returned variables slice describes each axis.
func ParallelCoordinates() ([]pc.PCDatum, []pc.PCVariable) {
	vars := []pc.PCVariable{
		{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
		{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
		{Key: "weight", Type: pc.PCScaleLinear, Label: "weight"},
		{Key: "volume", Type: pc.PCScaleLinear, Label: "volume"},
		{Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
	}
	data := []pc.PCDatum{
		{ID: "batch A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "weight": 30.0, "volume": 8.0, "grade": "B"}},
		{ID: "batch B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "weight": 12.0, "volume": 15.0, "grade": "A"}},
		{ID: "batch C", Values: map[string]any{"temp": 12.0, "cost": 2.0, "weight": 25.0, "volume": 3.0, "grade": "C"}},
		{ID: "batch D", Values: map[string]any{"temp": 28.0, "cost": 7.0, "weight": 18.0, "volume": 11.0, "grade": "B"}},
	}
	return data, vars
}
