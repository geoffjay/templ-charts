package chord

import (
	"math"
	"testing"
)

// sampleMatrix mirrors the @nivo/chord test fixture.
var sampleMatrix = [][]float64{
	{0, 1, 0, 1},
	{1, 0, 1, 0},
	{0, 1, 0, 1},
	{1, 0, 1, 0},
}

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestChord_GroupsSpanFullCircle(t *testing.T) {
	res := New().Compute(sampleMatrix)
	if len(res.Groups) != 4 {
		t.Fatalf("groups = %d, want 4", len(res.Groups))
	}
	// With padAngle 0, groups tile [0, 2π) contiguously.
	if !approx(res.Groups[0].StartAngle, 0) {
		t.Errorf("first group start = %v, want 0", res.Groups[0].StartAngle)
	}
	last := res.Groups[len(res.Groups)-1]
	if !approx(last.EndAngle, tau) {
		t.Errorf("last group end = %v, want 2π", last.EndAngle)
	}
	for i := 1; i < len(res.Groups); i++ {
		if !approx(res.Groups[i-1].EndAngle, res.Groups[i].StartAngle) {
			t.Errorf("group %d end (%v) != group %d start (%v)", i-1, res.Groups[i-1].EndAngle, i, res.Groups[i].StartAngle)
		}
	}
}

func TestChord_GroupValues(t *testing.T) {
	res := New().Compute(sampleMatrix)
	// Each row sums to 2; total 8. Each group spans 2/8 * 2π = π/2.
	for i, g := range res.Groups {
		if g.Value != 2 {
			t.Errorf("group %d value = %v, want 2", i, g.Value)
		}
		if !approx(g.EndAngle-g.StartAngle, math.Pi/2) {
			t.Errorf("group %d span = %v, want π/2", i, g.EndAngle-g.StartAngle)
		}
	}
}

func TestChord_RibbonCount(t *testing.T) {
	res := New().Compute(sampleMatrix)
	if len(res.Ribbons) != 4 {
		t.Fatalf("ribbons = %d, want 4", len(res.Ribbons))
	}
}

func TestChord_RibbonSourceIsLarger(t *testing.T) {
	// Asymmetric matrix so source/target ordering matters.
	m := [][]float64{
		{0, 5, 0},
		{2, 0, 3},
		{0, 4, 0},
	}
	res := New().Compute(m)
	for _, r := range res.Ribbons {
		if r.Source.Value < r.Target.Value {
			t.Errorf("ribbon source value %v < target %v; source must be the larger", r.Source.Value, r.Target.Value)
		}
	}
}

func TestChord_PadAngleReservesGaps(t *testing.T) {
	pad := 0.1
	res := New().PadAngle(pad).Compute(sampleMatrix)
	// Total used angle = 2π - pad*n for the arcs, plus pad*n gaps = 2π.
	var arcSpan float64
	for _, g := range res.Groups {
		arcSpan += g.EndAngle - g.StartAngle
	}
	want := tau - pad*float64(len(res.Groups))
	if !approx(arcSpan, want) {
		t.Errorf("total arc span = %v, want %v", arcSpan, want)
	}
}

func TestChord_Deterministic(t *testing.T) {
	a := New().Compute(sampleMatrix)
	b := New().Compute(sampleMatrix)
	if len(a.Ribbons) != len(b.Ribbons) {
		t.Fatalf("ribbon count differs between runs")
	}
	for i := range a.Ribbons {
		if a.Ribbons[i] != b.Ribbons[i] {
			t.Errorf("ribbon %d differs between runs", i)
		}
	}
}

func TestRibbon_ClosedPath(t *testing.T) {
	res := New().Compute(sampleMatrix)
	gen := NewRibbon(100)
	r := res.Ribbons[0]
	path := gen.Call(RibbonArg{
		Source: RibbonEndpoint{StartAngle: r.Source.StartAngle, EndAngle: r.Source.EndAngle},
		Target: RibbonEndpoint{StartAngle: r.Target.StartAngle, EndAngle: r.Target.EndAngle},
	})
	if len(path) == 0 || path[0] != 'M' {
		t.Fatalf("ribbon path should start with a moveto, got %q", path)
	}
	if path[len(path)-1] != 'Z' {
		t.Errorf("ribbon path should be closed (end with Z), got %q", path)
	}
}
