package d3hierarchy

import (
	"math"
	"testing"
)

func flat(vals ...float64) *datum {
	kids := make([]*datum, len(vals))
	for i, v := range vals {
		kids[i] = &datum{id: "c", val: v}
	}
	return &datum{id: "root", children: kids}
}

func summed(d *datum) *Node {
	return Hierarchy(d, childrenOf).Sum(func(x any) float64 { return x.(*datum).val })
}

func TestTreemap_Dice(t *testing.T) {
	root := summed(flat(6, 4))
	NewTreemap().Tile(TreemapDice).Size(100, 50).Layout(root)
	c0, c1 := root.Children[0], root.Children[1]
	if c0.X0 != 0 || c0.X1 != 60 || c0.Y0 != 0 || c0.Y1 != 50 {
		t.Errorf("dice c0 = [%v,%v,%v,%v], want [0,0,60,50]", c0.X0, c0.Y0, c0.X1, c0.Y1)
	}
	if c1.X0 != 60 || c1.X1 != 100 || c1.Y0 != 0 || c1.Y1 != 50 {
		t.Errorf("dice c1 = [%v,%v,%v,%v], want [60,0,100,50]", c1.X0, c1.Y0, c1.X1, c1.Y1)
	}
}

func TestTreemap_Slice(t *testing.T) {
	root := summed(flat(6, 4))
	NewTreemap().Tile(TreemapSlice).Size(50, 100).Layout(root)
	c0, c1 := root.Children[0], root.Children[1]
	if c0.Y0 != 0 || c0.Y1 != 60 || c0.X0 != 0 || c0.X1 != 50 {
		t.Errorf("slice c0 = [%v,%v,%v,%v], want [0,0,50,60]", c0.X0, c0.Y0, c0.X1, c0.Y1)
	}
	if c1.Y0 != 60 || c1.Y1 != 100 {
		t.Errorf("slice c1 y = [%v,%v], want [60,100]", c1.Y0, c1.Y1)
	}
}

// area returns the rect area of a node.
func area(n *Node) float64 { return (n.X1 - n.X0) * (n.Y1 - n.Y0) }

// assertProportional checks each leaf's area ≈ value/total * totalArea.
func assertProportional(t *testing.T, root *Node, w, h float64) {
	t.Helper()
	total := root.Value
	totalArea := w * h
	for _, leaf := range root.Leaves() {
		want := leaf.Value / total * totalArea
		if math.Abs(area(leaf)-want) > 1e-6 {
			t.Errorf("leaf value %v: area = %v, want %v", leaf.Value, area(leaf), want)
		}
		if leaf.X0 < -1e-9 || leaf.Y0 < -1e-9 || leaf.X1 > w+1e-9 || leaf.Y1 > h+1e-9 {
			t.Errorf("leaf out of bounds: [%v,%v,%v,%v]", leaf.X0, leaf.Y0, leaf.X1, leaf.Y1)
		}
	}
}

func TestTreemap_SquarifyProportional(t *testing.T) {
	root := summed(flat(6, 6, 4, 3, 2, 2, 1))
	NewTreemap().Size(200, 100).Layout(root) // squarify default
	assertProportional(t, root, 200, 100)
}

func TestTreemap_BinaryProportional(t *testing.T) {
	root := summed(flat(6, 6, 4, 3, 2, 2, 1))
	NewTreemap().Tile(TreemapBinary).Size(200, 100).Layout(root)
	assertProportional(t, root, 200, 100)
}

func TestTreemap_Round(t *testing.T) {
	root := summed(flat(1, 1, 1))
	NewTreemap().Size(100, 100).Round(true).Layout(root)
	for _, n := range root.Descendants() {
		if n.X0 != math.Round(n.X0) || n.Y1 != math.Round(n.Y1) {
			t.Errorf("round: node coords not integral: [%v,%v,%v,%v]", n.X0, n.Y0, n.X1, n.Y1)
		}
	}
}

func TestTreemap_PaddingInner(t *testing.T) {
	// With inner padding the two children should not touch.
	root := summed(flat(1, 1))
	NewTreemap().Tile(TreemapDice).Size(100, 100).PaddingInner(10).Layout(root)
	c0, c1 := root.Children[0], root.Children[1]
	if !(c1.X0 > c0.X1) {
		t.Errorf("expected a gap between children with paddingInner: c0.X1=%v c1.X0=%v", c0.X1, c1.X0)
	}
}
