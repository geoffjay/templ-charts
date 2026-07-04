package samples

import "github.com/geoffjay/templ-charts/charts/waffle"

// Waffle returns a three-part breakdown (men/women/children) summing to 92 of a
// 100-cell grid, ready to assign to a waffle.WaffleProps.Data field.
func Waffle() []waffle.WaffleDatum {
	return []waffle.WaffleDatum{
		{ID: "men", Label: "Men", Value: 32},
		{ID: "women", Label: "Women", Value: 41},
		{ID: "children", Label: "Children", Value: 19},
	}
}
