package d3hierarchy

import "testing"

// datum is a simple tree node for tests.
type datum struct {
	id       string
	val      float64
	children []*datum
}

func childrenOf(d any) []any {
	n := d.(*datum)
	if len(n.children) == 0 {
		return nil
	}
	out := make([]any, len(n.children))
	for i, c := range n.children {
		out[i] = c
	}
	return out
}

func sample() *datum {
	return &datum{id: "root", children: []*datum{
		{id: "A", children: []*datum{{id: "a1", val: 3}, {id: "a2", val: 2}}},
		{id: "B", children: []*datum{{id: "b1", val: 5}}},
	}}
}

func TestHierarchy_DepthHeight(t *testing.T) {
	root := Hierarchy(sample(), childrenOf)
	if root.Depth != 0 || root.Height != 2 {
		t.Errorf("root depth/height = %d/%d, want 0/2", root.Depth, root.Height)
	}
	a := root.Children[0]
	if a.Depth != 1 || a.Height != 1 {
		t.Errorf("A depth/height = %d/%d, want 1/1", a.Depth, a.Height)
	}
	a1 := a.Children[0]
	if a1.Depth != 2 || a1.Height != 0 {
		t.Errorf("a1 depth/height = %d/%d, want 2/0", a1.Depth, a1.Height)
	}
}

func TestHierarchy_Sum(t *testing.T) {
	root := Hierarchy(sample(), childrenOf).Sum(func(d any) float64 { return d.(*datum).val })
	if root.Value != 10 {
		t.Errorf("root sum = %v, want 10", root.Value)
	}
	if root.Children[0].Value != 5 || root.Children[1].Value != 5 {
		t.Errorf("A/B sums = %v/%v, want 5/5", root.Children[0].Value, root.Children[1].Value)
	}
}

func TestHierarchy_Count(t *testing.T) {
	root := Hierarchy(sample(), childrenOf).Count()
	if root.Value != 3 {
		t.Errorf("root count = %v, want 3", root.Value)
	}
	if root.Children[0].Value != 2 {
		t.Errorf("A count = %v, want 2", root.Children[0].Value)
	}
}

func TestHierarchy_LeavesAncestorsLinks(t *testing.T) {
	root := Hierarchy(sample(), childrenOf)
	if got := len(root.Leaves()); got != 3 {
		t.Errorf("leaves = %d, want 3", got)
	}
	if got := len(root.Descendants()); got != 6 {
		t.Errorf("descendants = %d, want 6", got)
	}
	if got := len(root.Links()); got != 5 {
		t.Errorf("links = %d, want 5", got)
	}
	a1 := root.Children[0].Children[0]
	anc := a1.Ancestors()
	if len(anc) != 3 || anc[0] != a1 || anc[2] != root {
		t.Errorf("a1 ancestors chain wrong: len=%d", len(anc))
	}
}

func TestHierarchy_Sort(t *testing.T) {
	root := Hierarchy(sample(), childrenOf).Sum(func(d any) float64 { return d.(*datum).val })
	// descending by value: a1(3) before a2(2).
	root.Sort(func(a, b *Node) bool { return a.Value > b.Value })
	if root.Children[0].Children[0].Data.(*datum).id != "a1" {
		t.Errorf("expected a1 first after descending sort")
	}
}

func TestHierarchy_EachOrder(t *testing.T) {
	root := Hierarchy(sample(), childrenOf)
	var pre, post, bfs []string
	root.EachBefore(func(n *Node) { pre = append(pre, n.Data.(*datum).id) })
	root.EachAfter(func(n *Node) { post = append(post, n.Data.(*datum).id) })
	root.Each(func(n *Node) { bfs = append(bfs, n.Data.(*datum).id) })
	wantPre := []string{"root", "A", "a1", "a2", "B", "b1"}
	wantPost := []string{"a1", "a2", "A", "b1", "B", "root"}
	wantBfs := []string{"root", "A", "B", "a1", "a2", "b1"}
	if !eqStr(pre, wantPre) {
		t.Errorf("preorder = %v, want %v", pre, wantPre)
	}
	if !eqStr(post, wantPost) {
		t.Errorf("postorder = %v, want %v", post, wantPost)
	}
	if !eqStr(bfs, wantBfs) {
		t.Errorf("bfs = %v, want %v", bfs, wantBfs)
	}
}

func eqStr(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
