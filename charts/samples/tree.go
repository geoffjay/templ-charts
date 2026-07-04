package samples

import "github.com/geoffjay/templ-charts/charts/tree"

// Tree returns a small three-level hierarchy ready to assign to
// tree.TreeProps.Data. Mirrors the demo's treeSample so the dendrogram/tidy
// layouts have a stable, golden-friendly shape.
func Tree() tree.TreeNode {
	return tree.TreeNode{ID: "root", Children: []tree.TreeNode{
		{ID: "A", Children: []tree.TreeNode{
			{ID: "A.1", Children: []tree.TreeNode{{ID: "A.1.a"}, {ID: "A.1.b"}}},
			{ID: "A.2"},
		}},
		{ID: "B", Children: []tree.TreeNode{{ID: "B.1"}, {ID: "B.2"}, {ID: "B.3"}}},
		{ID: "C"},
	}}
}
