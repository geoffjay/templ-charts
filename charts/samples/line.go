package samples

import "github.com/geoffjay/templ-charts/charts/line"

// Line returns ready-to-render line data: three country series over seven
// years, with deterministic y values.
func Line() []line.LineSeries {
	years := []string{"2018", "2019", "2020", "2021", "2022", "2023", "2024"}
	countries := []string{"USA", "Germany", "France"}
	series := make([]line.LineSeries, len(countries))
	for i, c := range countries {
		data := make([]line.LinePointData, len(years))
		for j, y := range years {
			data[j] = line.LinePointData{X: y, Y: float64((i*3+j*11)%90 + 20)}
		}
		series[i] = line.LineSeries{ID: c, Data: data}
	}
	return series
}
