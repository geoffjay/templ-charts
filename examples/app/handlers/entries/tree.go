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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := tree.TreeProps{
				Width: 720, Height: 440, Responsive: true,
				// A project layout: seven top-level areas with varied children.
				Data: tree.TreeNode{ID: "app", Children: []tree.TreeNode{
					{ID: "api", Children: []tree.TreeNode{
						{ID: "handlers"}, {ID: "middleware"}, {ID: "routes"},
					}},
					{ID: "ui", Children: []tree.TreeNode{
						{ID: "components"}, {ID: "views"}, {ID: "hooks"}, {ID: "styles"},
					}},
					{ID: "data", Children: []tree.TreeNode{
						{ID: "models"}, {ID: "migrations"},
					}},
					{ID: "auth", Children: []tree.TreeNode{
						{ID: "sessions"}, {ID: "tokens"}, {ID: "oauth"},
					}},
					{ID: "jobs", Children: []tree.TreeNode{
						{ID: "scheduler"}, {ID: "workers"},
					}},
					{ID: "infra", Children: []tree.TreeNode{
						{ID: "docker"}, {ID: "ci"}, {ID: "terraform"},
					}},
					{ID: "docs", Children: []tree.TreeNode{
						{ID: "guides"}, {ID: "reference"},
					}},
				}},
				Theme:   theme,
				Animate: animate,
			}
			return render.String(tree.Tree(p))
		},
	})
}
