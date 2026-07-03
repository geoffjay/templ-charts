package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/tree"
)

func init() {
	register(ChartEntry{
		Slug:        "tree",
		Title:       "Tree",
		Description: "Tidy-tree / dendrogram node-link diagrams with bump links.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/tree"
)

svg, _ := render.String(tree.Tree(tree.TreeProps{
    Width: 720, Height: 440,
    Data: tree.TreeNode{ID: "root", Children: []tree.TreeNode{
        {ID: "A"},
        {ID: "B"},
    }},
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := tree.TreeProps{
				Width: 720, Height: 440, Responsive: true,
				Data: tree.TreeNode{ID: "root", Children: []tree.TreeNode{
					{ID: "A"},
					{ID: "B"},
				}},
				Theme: theme,
			}
			return render.String(tree.Tree(p))
		},
	})
}
