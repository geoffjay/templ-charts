package force

import (
	"math"
	"testing"
)

// chargeNodes builds n phyllotaxis-seeded nodes (distinct positions, so jiggle
// never fires and the exact vs Barnes–Hut paths differ only by summation order).
func chargeNodes(n int) []*Node {
	ns := make([]*Node, n)
	for i := range ns {
		ns[i] = NewNode(i)
	}
	return ns
}

func chargeSim(nodes []*Node, theta float64) *Simulation {
	charge := ForceManyBody().StrengthConst(-30)
	if theta > 0 {
		charge.Theta(theta)
	}
	return NewSimulation(nodes).Force("charge", charge).Stop()
}

// TestManyBodyBarnesHutConvergesToExact: with a vanishingly small θ the
// Barnes–Hut criterion never fires, so the tree is fully descended and the
// result equals the exact all-pairs sum up to floating-point summation order.
func TestManyBodyBarnesHutConvergesToExact(t *testing.T) {
	exact := chargeNodes(30)
	chargeSim(exact, 0).Tick(120)
	bh := chargeNodes(30)
	chargeSim(bh, 1e-9).Tick(120)
	for i := range exact {
		if math.Abs(exact[i].X-bh[i].X) > 1e-6 || math.Abs(exact[i].Y-bh[i].Y) > 1e-6 {
			t.Errorf("node %d: exact (%v,%v) vs θ→0 BH (%v,%v)", i, exact[i].X, exact[i].Y, bh[i].X, bh[i].Y)
		}
	}
}

// TestManyBodyBarnesHutDeterministic: the approximation is a pure function of
// its inputs (deterministic quadtree + LCG jiggle), so two runs match exactly.
func TestManyBodyBarnesHutDeterministic(t *testing.T) {
	a := chargeNodes(40)
	chargeSim(a, 0.9).Tick(120)
	b := chargeNodes(40)
	chargeSim(b, 0.9).Tick(120)
	for i := range a {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			t.Fatalf("node %d not deterministic under Barnes–Hut: (%v,%v) vs (%v,%v)", i, a[i].X, a[i].Y, b[i].X, b[i].Y)
		}
	}
}

// TestManyBodyBarnesHutApproxReasonable: at the usual θ=0.9 the approximation
// should still repel — a similar overall spread to the exact layout (not
// collapsed, not exploded). Loose bound: within 25% of the exact bounding box.
func TestManyBodyBarnesHutApproxReasonable(t *testing.T) {
	exact := chargeNodes(60)
	chargeSim(exact, 0).Tick(120)
	bh := chargeNodes(60)
	chargeSim(bh, 0.9).Tick(120)
	se := spread(exact)
	sb := spread(bh)
	if sb < se*0.75 || sb > se*1.25 {
		t.Errorf("Barnes–Hut spread %v not within 25%% of exact spread %v", sb, se)
	}
}

func spread(nodes []*Node) float64 {
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for _, n := range nodes {
		minX, maxX = math.Min(minX, n.X), math.Max(maxX, n.X)
		minY, maxY = math.Min(minY, n.Y), math.Max(maxY, n.Y)
	}
	return math.Hypot(maxX-minX, maxY-minY)
}

// gridNodes places n nodes on a tight grid so many overlap at the given radius,
// exercising the collision force. Positions are set after seeding so both the
// exact and quadtree passes start identically.
func gridNodes(n int) []*Node {
	ns := make([]*Node, n)
	for i := range ns {
		ns[i] = NewNode(i)
	}
	return ns
}

func placeGrid(nodes []*Node) {
	for i, node := range nodes {
		node.X = float64(i%5) * 8
		node.Y = float64(i/5) * 8
		node.Vx, node.Vy = 0, 0
	}
}

func collideSim(nodes []*Node, quadtree bool) *Simulation {
	c := ForceCollide(func(*Node) float64 { return 6 })
	if quadtree {
		c.UseQuadtree()
	}
	sim := NewSimulation(nodes).Force("collide", c).Stop()
	placeGrid(nodes) // override the phyllotaxis seed with an overlapping grid
	return sim
}

// TestCollideQuadtreeMatchesExact: the pruned pass visits the same overlapping
// pairs the exact pass does, in the same per-node order, so after one tick the
// positions match up to floating-point summation order.
func TestCollideQuadtreeMatchesExact(t *testing.T) {
	exact := gridNodes(25)
	collideSim(exact, false).Tick(1)
	quad := gridNodes(25)
	collideSim(quad, true).Tick(1)
	for i := range exact {
		if math.Abs(exact[i].X-quad[i].X) > 1e-9 || math.Abs(exact[i].Y-quad[i].Y) > 1e-9 {
			t.Errorf("node %d: exact (%v,%v) vs quadtree (%v,%v)", i, exact[i].X, exact[i].Y, quad[i].X, quad[i].Y)
		}
	}
}

func TestCollideQuadtreeDeterministic(t *testing.T) {
	a := gridNodes(25)
	collideSim(a, true).Tick(30)
	b := gridNodes(25)
	collideSim(b, true).Tick(30)
	for i := range a {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			t.Fatalf("node %d not deterministic under quadtree collide: (%v,%v) vs (%v,%v)", i, a[i].X, a[i].Y, b[i].X, b[i].Y)
		}
	}
}

// TestCollideQuadtreeNoOverlap mirrors TestSwarmCollideNoOverlap with the
// quadtree pass enabled: nodes should spread so circles no longer overlap.
func TestCollideQuadtreeNoOverlap(t *testing.T) {
	const size = 6.0
	nodes := make([]*Node, 20)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	radius := func(*Node) float64 { return size/2 + 2.0/2 }
	sim := NewSimulation(nodes).
		Force("x", ForceX(func(*Node) float64 { return 200 }).Strength(1)).
		Force("y", ForceY(func(*Node) float64 { return 100 })).
		Force("collide", ForceCollide(radius).UseQuadtree()).
		Stop().
		Tick(120)
	got := sim.Nodes()
	minGap := size - 1.0
	for i := 0; i < len(got); i++ {
		for j := i + 1; j < len(got); j++ {
			d := math.Hypot(got[i].X-got[j].X, got[i].Y-got[j].Y)
			if d < minGap {
				t.Errorf("nodes %d,%d overlap: distance %v < %v", i, j, d, minGap)
			}
		}
	}
}
