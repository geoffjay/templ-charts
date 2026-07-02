package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/tree"
)

// TreeDemo is one tree tile on the /tree page (static render).
type TreeDemo struct {
	ID          string
	Title       string
	Description string
	Props       tree.TreeProps
}

func treeSample() tree.TreeNode {
	return tree.TreeNode{ID: "root", Children: []tree.TreeNode{
		{ID: "A", Children: []tree.TreeNode{
			{ID: "A.1", Children: []tree.TreeNode{{ID: "A.1.a"}, {ID: "A.1.b"}}},
			{ID: "A.2"},
		}},
		{ID: "B", Children: []tree.TreeNode{{ID: "B.1"}, {ID: "B.2"}, {ID: "B.3"}}},
		{ID: "C"},
	}}
}

// TreeDemos returns the tree demos for the /tree page.
func TreeDemos() []TreeDemo {
	return []TreeDemo{
		{
			ID:          "tree-dendogram",
			Title:       "Dendrogram (top-to-bottom)",
			Description: "The cluster layout (all leaves aligned) with smooth bump links (curveBumpY) and per-node hover.",
			Props: tree.TreeProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 40, Bottom: 40, Left: 40},
				Data:        treeSample(),
				Interactive: true,
			},
		},
		{
			ID:          "tree-tidy-left",
			Title:       "Tidy tree (left-to-right)",
			Description: "mode=tree with layout=left-to-right; links use curveBumpX.",
			Props: tree.TreeProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin: core.Margin{Top: 20, Right: 80, Bottom: 20, Left: 40},
				Data:   treeSample(),
				Mode:   tree.ModeTree,
				Layout: tree.LayoutLeftToRight,
			},
		},
		{
			ID:          "tree-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover anywhere resolves to the nearest node via an accurate Voronoi mesh (internal/d3/delaunay), instead of per-node targets.",
			Props: tree.TreeProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 40, Bottom: 40, Left: 40},
				Data:        treeSample(),
				Interactive: true,
				UseMesh:     true,
			},
		},
	}
}
