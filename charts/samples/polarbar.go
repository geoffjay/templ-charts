package samples

import "github.com/geoffjay/templ-charts/charts/polarbar"

// PolarBar returns ready-to-render polar-bar data: five weekday indices, each
// with three transport-mode keys to stack radially. The returned keys slice is
// the stack order.
func PolarBar() ([]polarbar.PolarBarDatum, []string) {
	data := []polarbar.PolarBarDatum{
		{Index: "Mon", Values: map[string]float64{"walk": 8, "bus": 5, "bike": 3}},
		{Index: "Tue", Values: map[string]float64{"walk": 6, "bus": 7, "bike": 4}},
		{Index: "Wed", Values: map[string]float64{"walk": 10, "bus": 4, "bike": 6}},
		{Index: "Thu", Values: map[string]float64{"walk": 7, "bus": 9, "bike": 2}},
		{Index: "Fri", Values: map[string]float64{"walk": 12, "bus": 3, "bike": 8}},
	}
	keys := []string{"walk", "bus", "bike"}
	return data, keys
}
