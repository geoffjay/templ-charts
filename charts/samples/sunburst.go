package samples

import "github.com/geoffjay/templ-charts/charts/sunburst"

// Sunburst returns a small three-level hierarchy (fruit/veg subtrees plus a
// leaf "grain"), ready to assign to a sunburst.SunburstProps.Data field.
func Sunburst() sunburst.SunburstNode {
	return sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
		{ID: "fruit", Children: []sunburst.SunburstNode{
			{ID: "apple", Value: 30}, {ID: "pear", Value: 18}, {ID: "grape", Value: 22},
		}},
		{ID: "veg", Children: []sunburst.SunburstNode{
			{ID: "carrot", Value: 20}, {ID: "pea", Value: 14},
		}},
		{ID: "grain", Value: 26},
	}}
}
