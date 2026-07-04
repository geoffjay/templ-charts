package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// TreemapDemo is one treemap tile on the /treemap page (static render).
type TreemapDemo struct {
	ID          string
	Title       string
	Description string
	Props       treemap.TreemapProps
}

// treemapSample is the shared hierarchy, sourced from the public samples
// package.
func treemapSample() treemap.TreemapNode {
	return samples.Treemap()
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
