package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// TreemapDemo is one treemap tile on the /treemap page (static render).
type TreemapDemo struct {
	ID          string
	Title       string
	Description string
	Props       treemap.TreemapProps
}

func treemapSample() treemap.TreemapNode {
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

// TreemapDemos returns the treemap demos for the /treemap page.
func TreemapDemos() []TreemapDemo {
	return []TreemapDemo{
		{
			ID:          "treemap-squarify",
			Title:       "Squarify tiling + labels",
			Description: "A three-level hierarchy tiled with the default squarify algorithm, leaf value labels, parent labels, and per-node hover.",
			Props: treemap.TreemapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:        treemapSample(),
				Interactive: true,
			},
		},
		{
			ID:          "treemap-binary",
			Title:       "Binary tiling + padding",
			Description: "tile=binary with inner/outer padding for a clearer nesting.",
			Props: treemap.TreemapProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:       core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
				Data:         treemapSample(),
				Tile:         treemap.TileBinary,
				InnerPadding: 4,
				OuterPadding: 4,
			},
		},
	}
}
