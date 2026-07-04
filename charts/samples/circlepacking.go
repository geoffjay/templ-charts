package samples

import cp "github.com/geoffjay/templ-charts/charts/circlepacking"

// CirclePacking returns a small three-branch hierarchy of leaf values, ready to
// assign to CirclePackingProps.Data.
func CirclePacking() cp.CirclePackingNode {
	return cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
		{ID: "A", Children: []cp.CirclePackingNode{
			{ID: "a1", Value: 20}, {ID: "a2", Value: 12}, {ID: "a3", Value: 8},
		}},
		{ID: "B", Children: []cp.CirclePackingNode{
			{ID: "b1", Value: 16}, {ID: "b2", Value: 24},
		}},
		{ID: "C", Children: []cp.CirclePackingNode{
			{ID: "c1", Value: 10}, {ID: "c2", Value: 6}, {ID: "c3", Value: 14}, {ID: "c4", Value: 5},
		}},
	}}
}
