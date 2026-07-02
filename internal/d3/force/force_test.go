package force

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// TestPhyllotaxisInit asserts unset nodes are seeded on d3's phyllotaxis spiral
// (deterministic, no Math.random) and pinned nodes take their fixed positions.
func TestPhyllotaxisInit(t *testing.T) {
	nodes := []*Node{NewNode(nil), NewNode(nil), NewNode(nil)}
	NewSimulation(nodes)
	for i, node := range nodes {
		if node.Index != i {
			t.Fatalf("node %d Index = %d", i, node.Index)
		}
		radius := initialRadius * math.Sqrt(0.5+float64(i))
		angle := float64(i) * initialAngle
		wantX := radius * math.Cos(angle)
		wantY := radius * math.Sin(angle)
		if !approx(node.X, wantX) || !approx(node.Y, wantY) {
			t.Errorf("node %d = (%v,%v), want (%v,%v)", i, node.X, node.Y, wantX, wantY)
		}
		if node.Vx != 0 || node.Vy != 0 {
			t.Errorf("node %d velocity = (%v,%v), want 0", i, node.Vx, node.Vy)
		}
	}
}

func TestFixedPosition(t *testing.T) {
	fx, fy := 5.0, 7.0
	pinned := NewNode(nil)
	pinned.Fx, pinned.Fy = &fx, &fy
	sim := NewSimulation([]*Node{pinned}).
		Force("charge", ForceManyBody()).
		Tick(50)
	got := sim.Nodes()[0]
	if got.X != fx || got.Y != fy {
		t.Errorf("pinned node moved to (%v,%v), want (%v,%v)", got.X, got.Y, fx, fy)
	}
}

// newNetworkSim builds a small deterministic network like the network chart:
// link + charge + center forces, then a fixed tick count.
func newNetworkSim() *Simulation {
	nodes := make([]*Node, 5)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	links := []*Link{
		NewLink(nodes[0], nodes[1], nil),
		NewLink(nodes[0], nodes[2], nil),
		NewLink(nodes[1], nodes[3], nil),
		NewLink(nodes[2], nodes[4], nil),
	}
	return NewSimulation(nodes).
		Force("link", ForceLink(links).DistanceConst(30).StrengthConst(1)).
		Force("charge", ForceManyBody().StrengthConst(-10).DistanceMin(1)).
		Force("center", ForceCenter(100, 100)).
		Stop()
}

// TestDeterministic asserts two identical simulations produce byte-identical
// positions after the same tick count — the core determinism guarantee.
func TestDeterministic(t *testing.T) {
	a := newNetworkSim().Tick(120).Nodes()
	b := newNetworkSim().Tick(120).Nodes()
	for i := range a {
		if a[i].X != b[i].X || a[i].Y != b[i].Y {
			t.Fatalf("node %d not deterministic: (%v,%v) vs (%v,%v)", i, a[i].X, a[i].Y, b[i].X, b[i].Y)
		}
	}
}

// TestCenterForce asserts the centering force keeps the mean node position at
// the requested center (it is exact each tick, ignoring alpha).
func TestCenterForce(t *testing.T) {
	nodes := newNetworkSim().Tick(120).Nodes()
	var sx, sy float64
	for _, n := range nodes {
		sx += n.X
		sy += n.Y
	}
	cx := sx / float64(len(nodes))
	cy := sy / float64(len(nodes))
	// forceCenter pins the mean each tick, but the velocity-integration step
	// that follows shifts it by the mean velocity; near convergence that drift
	// is tiny, so the mean sits within a pixel of the requested center.
	if math.Abs(cx-100) > 1 || math.Abs(cy-100) > 1 {
		t.Errorf("center of mass = (%v,%v), want ~(100,100)", cx, cy)
	}
}

// TestLinkDistanceRelaxes asserts linked pairs settle near the link distance.
func TestLinkDistanceRelaxes(t *testing.T) {
	sim := newNetworkSim().Tick(300)
	nodes := sim.Nodes()
	dist := func(a, b *Node) float64 {
		return math.Hypot(a.X-b.X, a.Y-b.Y)
	}
	// The 0–1 edge should be within a reasonable band of the 30px target.
	d := dist(nodes[0], nodes[1])
	if d < 10 || d > 80 {
		t.Errorf("edge 0-1 length = %v, expected roughly ~30", d)
	}
}

// TestForceRemoval asserts a nil force detaches a previously-added one.
func TestForceRemoval(t *testing.T) {
	sim := newNetworkSim()
	before := len(sim.forces)
	sim.Force("charge", nil)
	if len(sim.forces) != before-1 {
		t.Fatalf("force count = %d, want %d", len(sim.forces), before-1)
	}
}

// TestSwarmCollideNoOverlap asserts collide + positioning forces spread nodes
// so circles no longer overlap (the swarmplot use case).
func TestSwarmCollideNoOverlap(t *testing.T) {
	const size = 6.0
	nodes := make([]*Node, 20)
	for i := range nodes {
		nodes[i] = NewNode(i)
	}
	radius := func(*Node) float64 { return size/2 + 2.0/2 } // size/2 + spacing/2
	sim := NewSimulation(nodes).
		Force("x", ForceX(func(*Node) float64 { return 200 }).Strength(1)).
		Force("y", ForceY(func(*Node) float64 { return 100 })).
		Force("collide", ForceCollide(radius)).
		Stop().
		Tick(120)
	got := sim.Nodes()
	// After settling, no two nodes should be closer than ~ (2*size/2) apart,
	// allowing a small tolerance for the finite iteration count.
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
