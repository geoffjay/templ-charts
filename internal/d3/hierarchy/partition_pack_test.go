package d3hierarchy

import (
	"math"
	"testing"
)

func TestPartition_Bands(t *testing.T) {
	root := summed(flat(6, 4))
	NewPartition().Size(100, 100).Layout(root)
	// root occupies the top band [0,50]; leaves the bottom band [50,100].
	if root.Y0 != 0 || root.Y1 != 50 {
		t.Errorf("root y = [%v,%v], want [0,50]", root.Y0, root.Y1)
	}
	c0, c1 := root.Children[0], root.Children[1]
	if c0.Y0 != 50 || c0.Y1 != 100 || c0.X0 != 0 || c0.X1 != 60 {
		t.Errorf("c0 = [%v,%v,%v,%v], want [0,50,60,100]", c0.X0, c0.Y0, c0.X1, c0.Y1)
	}
	if c1.X0 != 60 || c1.X1 != 100 {
		t.Errorf("c1 x = [%v,%v], want [60,100]", c1.X0, c1.X1)
	}
}

func TestPack_NoOverlapWithinBounds(t *testing.T) {
	root := summed(flat(6, 6, 4, 3, 2, 2, 1, 5, 8))
	NewPack().Size(200, 200).Layout(root)

	leaves := root.Leaves()
	// Leaf radius ∝ sqrt(value) up to the layout scale factor: check the ratio
	// between two leaves matches sqrt(value) ratio.
	for _, l := range leaves {
		if l.R <= 0 {
			t.Errorf("leaf radius non-positive: %v", l.R)
		}
		// within the layout square (allow a small epsilon).
		if l.X-l.R < -1e-6 || l.Y-l.R < -1e-6 || l.X+l.R > 200+1e-6 || l.Y+l.R > 200+1e-6 {
			t.Errorf("leaf circle out of bounds: x=%v y=%v r=%v", l.X, l.Y, l.R)
		}
	}
	// No two leaf circles overlap (allow tiny epsilon).
	for i := 0; i < len(leaves); i++ {
		for j := i + 1; j < len(leaves); j++ {
			a, b := leaves[i], leaves[j]
			d := math.Hypot(a.X-b.X, a.Y-b.Y)
			if d < a.R+b.R-1e-6 {
				t.Errorf("leaves %d and %d overlap: d=%v r0+r1=%v", i, j, d, a.R+b.R)
			}
		}
	}
}

func TestPack_Deterministic(t *testing.T) {
	build := func() []float64 {
		root := summed(flat(6, 6, 4, 3, 2, 2, 1, 5, 8))
		NewPack().Size(200, 200).Layout(root)
		var out []float64
		for _, l := range root.Leaves() {
			out = append(out, l.X, l.Y, l.R)
		}
		return out
	}
	a, b := build(), build()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("pack not deterministic at %d: %v vs %v", i, a[i], b[i])
		}
	}
}

func TestPack_RadiusRatio(t *testing.T) {
	// Two leaves with values 4 and 1 → radius ratio 2:1 (sqrt(4):sqrt(1)),
	// preserved through the uniform scale.
	root := summed(flat(4, 1))
	NewPack().Size(100, 100).Layout(root)
	r0 := root.Children[0].R
	r1 := root.Children[1].R
	if math.Abs(r0/r1-2) > 1e-9 {
		t.Errorf("radius ratio = %v, want 2", r0/r1)
	}
}
