package tree_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/tree"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() tree.TreeNode {
	return tree.TreeNode{ID: "root", Children: []tree.TreeNode{
		{ID: "A", Children: []tree.TreeNode{{ID: "a1"}, {ID: "a2"}}},
		{ID: "B", Children: []tree.TreeNode{{ID: "b1"}, {ID: "b2"}, {ID: "b3"}}},
	}}
}

func baseProps() tree.TreeProps {
	return tree.TreeProps{
		Width: 400, Height: 300,
		Margin: core.Margin{Top: 30, Right: 30, Bottom: 30, Left: 30},
		Data:   sampleData(),
	}
}

func render(t *testing.T, props tree.TreeProps) string {
	t.Helper()
	var b strings.Builder
	if err := tree.Tree(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Tree.Render: %v", err)
	}
	return b.String()
}

func TestTree_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestTree_NodeAndLinkCount(t *testing.T) {
	// 8 nodes (root + A,B + a1,a2,b1,b2,b3) → 7 links.
	out := render(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 8 {
		t.Errorf("node count = %d, want 8", got)
	}
	if got := strings.Count(out, `fill="none" stroke=`); got != 7 {
		t.Errorf("link count = %d, want 7", got)
	}
}

func TestTree_BumpLinks(t *testing.T) {
	// Links use cubic bump curves → C commands in the path.
	out := render(t, baseProps())
	if !strings.Contains(out, `d="M`) || !strings.Contains(out, "C") {
		t.Errorf("expected bump link paths with cubic segments")
	}
}

func TestTree_TreeMode(t *testing.T) {
	p := baseProps()
	p.Mode = tree.ModeTree
	out := render(t, p)
	if strings.Count(out, "<circle") != 8 {
		t.Errorf("tree mode should still render 8 nodes")
	}
}

func TestTree_Layouts(t *testing.T) {
	for _, l := range []tree.LayoutDir{tree.LayoutTopToBottom, tree.LayoutBottomToTop, tree.LayoutLeftToRight, tree.LayoutRightToLeft} {
		p := baseProps()
		p.Layout = l
		if !strings.HasPrefix(render(t, p), "<svg") {
			t.Errorf("layout %q failed", l)
		}
	}
}

func TestTree_Golden(t *testing.T) {
	golden.Assert(t, "tree-basic", render(t, baseProps()))
}

func TestTree_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Org"
	p.Desc = "Org chart."
	out := render(t, p)
	if !strings.Contains(out, "<title>Org</title>") || !strings.Contains(out, "<desc>Org chart.</desc>") {
		t.Errorf("expected title/desc")
	}
}

func TestTree_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(render(t, p), "data-tc-tooltip") {
		t.Errorf("interactive tree should emit data-tc-tooltip")
	}
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive tree must not emit data-tc-tooltip")
	}
}
