package samples

import "github.com/geoffjay/templ-charts/charts/bar"

// Bar returns ready-to-render bar data: seven countries each measured across
// six food-sales keys, plus the key list to stack/group. Values are
// deterministic so renders are stable. The index-by property is "country".
func Bar() (data []bar.BarDatum, keys []string) {
	keys = []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"}
	countries := []string{"USA", "Germany", "France", "Japan", "Brazil", "India", "China"}
	data = make([]bar.BarDatum, len(countries))
	for i, c := range countries {
		row := map[string]any{"country": c}
		for j, k := range keys {
			row[k] = float64((i*7+j*13)%100 + 10)
		}
		data[i] = row
	}
	return data, keys
}
