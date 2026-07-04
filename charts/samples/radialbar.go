package samples

import "github.com/geoffjay/templ-charts/charts/radialbar"

// RadialBar returns ready-to-render radial-bar data: three sources, each
// scored across three grocery categories, to be stacked over an angle scale.
func RadialBar() []radialbar.RadialBarSerie {
	return []radialbar.RadialBarSerie{
		{ID: "Supermarket", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 25}, {X: "Fruits", Y: 18}, {X: "Meat", Y: 12},
		}},
		{ID: "Combini", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 12}, {X: "Fruits", Y: 9}, {X: "Meat", Y: 7},
		}},
		{ID: "Online", Data: []radialbar.RadialBarDatum{
			{X: "Vegetables", Y: 8}, {X: "Fruits", Y: 15}, {X: "Meat", Y: 4},
		}},
	}
}
