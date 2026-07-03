package tree_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/tree"
)

// ExampleTree renders a small node-link tree diagram to an SVG string with the
// charts/render helper.
func ExampleTree() {
	svg, err := render.String(tree.Tree(tree.TreeProps{
		Width: 400, Height: 300,
		Data: tree.TreeNode{ID: "root", Children: []tree.TreeNode{
			{ID: "A"},
			{ID: "B"},
		}},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
