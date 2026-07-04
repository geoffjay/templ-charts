package samples

import "github.com/geoffjay/templ-charts/charts/treemap"

// Treemap returns a small three-level weighted hierarchy ready to assign to
// treemap.TreemapProps.Data. Mirrors the demo's treemapSample.
func Treemap() treemap.TreemapNode {
	return treemap.TreemapNode{ID: "nivo", Children: []treemap.TreemapNode{
		{ID: "viz", Children: []treemap.TreemapNode{
			{ID: "stack", Value: 33}, {ID: "chart", Value: 41}, {ID: "xAxis", Value: 18},
		}},
		{ID: "colors", Children: []treemap.TreemapNode{
			{ID: "rgb", Value: 25}, {ID: "hsl", Value: 16},
		}},
		{ID: "utils", Children: []treemap.TreemapNode{
			{ID: "randomize", Value: 22}, {ID: "sortBy", Value: 30}, {ID: "format", Value: 12},
		}},
		{ID: "misc", Value: 45},
	}}
}
