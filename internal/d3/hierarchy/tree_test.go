package d3hierarchy

import (
	"math"
	"testing"
)

// balanced returns a full binary tree of the given depth (leaves carry val 1).
func balanced() *datum {
	leaf := func(id string) *datum { return &datum{id: id, val: 1} }
	return &datum{id: "root", children: []*datum{
		{id: "A", children: []*datum{leaf("a1"), leaf("a2")}},
		{id: "B", children: []*datum{leaf("b1"), leaf("b2")}},
	}}
}

func TestCluster_LeavesAligned(t *testing.T) {
	root := Hierarchy(balanced(), childrenOf)
	NewCluster().Size(100, 60).Layout(root)
	if root.Y != 0 {
		t.Errorf("cluster root Y = %v, want 0 (top)", root.Y)
	}
	for _, leaf := range root.Leaves() {
		if math.Abs(leaf.Y-60) > 1e-9 {
			t.Errorf("cluster leaf Y = %v, want 60 (all leaves aligned at bottom)", leaf.Y)
		}
		if leaf.X < -1e-9 || leaf.X > 100+1e-9 {
			t.Errorf("cluster leaf X out of [0,100]: %v", leaf.X)
		}
	}
	// symmetric tree: root centered.
	if math.Abs(root.X-50) > 1e-9 {
		t.Errorf("cluster root X = %v, want 50 (centered)", root.X)
	}
}

func TestTree_ParentCentered(t *testing.T) {
	root := Hierarchy(balanced(), childrenOf)
	NewTree().Size(100, 60).Layout(root)
	// Y by depth: root 0, mid 30, leaves 60.
	if root.Y != 0 {
		t.Errorf("tree root Y = %v, want 0", root.Y)
	}
	a := root.Children[0]
	if math.Abs(a.Y-30) > 1e-9 {
		t.Errorf("tree A Y = %v, want 30", a.Y)
	}
	for _, leaf := range root.Leaves() {
		if math.Abs(leaf.Y-60) > 1e-9 {
			t.Errorf("tree leaf Y = %v, want 60", leaf.Y)
		}
	}
	// Each parent centered over its children.
	mid := (a.Children[0].X + a.Children[1].X) / 2
	if math.Abs(a.X-mid) > 1e-9 {
		t.Errorf("tree A X = %v, want midpoint of children %v", a.X, mid)
	}
	if math.Abs(root.X-50) > 1e-9 {
		t.Errorf("tree root X = %v, want 50 (centered)", root.X)
	}
}

func TestTree_Deterministic(t *testing.T) {
	build := func() []float64 {
		root := Hierarchy(balanced(), childrenOf)
		NewTree().Size(100, 60).Layout(root)
		var out []float64
		for _, n := range root.Descendants() {
			out = append(out, n.X, n.Y)
		}
		return out
	}
	a, b := build(), build()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("tree not deterministic at %d", i)
		}
	}
}
