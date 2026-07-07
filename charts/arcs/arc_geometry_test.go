package arcs_test

import (
	"math"
	"testing"

	"github.com/geoffjay/templ-charts/charts/arcs"
)

func approx(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func TestCentroid(t *testing.T) {
	g := arcs.CreateArcGenerator(0, 0)
	// Quarter arc 0..π/2, radii 0..10: mid angle (after the -π/2 d3 offset)
	// is -π/4, mid radius 5 → centroid (5·cos(-π/4), 5·sin(-π/4)).
	c := g.Centroid(arcs.Arc{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10})
	approx(t, "centroid.x", c[0], 5*math.Cos(-math.Pi/4))
	approx(t, "centroid.y", c[1], 5*math.Sin(-math.Pi/4))
}

func TestPositionFromAngle(t *testing.T) {
	p := arcs.PositionFromAngle(0, 5)
	approx(t, "x@0", p.X, 5)
	approx(t, "y@0", p.Y, 0)
	p = arcs.PositionFromAngle(math.Pi/2, 2)
	approx(t, "x@90", p.X, 0)
	approx(t, "y@90", p.Y, 2)
	p = arcs.PositionFromAngle(math.Pi, 3)
	approx(t, "x@180", p.X, -3)
	approx(t, "y@180", p.Y, 0)
}

func TestNormalizeAngleDegrees(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{-10, 350}, {370, 10}, {360, 0}, {0, 0}, {-720, 0}, {725, 5},
	}
	for _, c := range cases {
		if got := arcs.NormalizeAngleDegrees(c.in); got != c.want {
			t.Errorf("NormalizeAngleDegrees(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestRadToDeg(t *testing.T) {
	approx(t, "RadToDeg(π)", arcs.RadToDeg(math.Pi), 180)
	approx(t, "RadToDeg(π/2)", arcs.RadToDeg(math.Pi/2), 90)
}

type testDatum struct{ a arcs.Arc }

func (d testDatum) GetArc() arcs.Arc { return d.a }

func TestFilterDataBySkipAngle_NonPositiveReturnsCopy(t *testing.T) {
	data := []arcs.DatumWithArc{
		testDatum{arcs.Arc{StartAngle: 0, EndAngle: arcs.DegToRad(5)}},
		testDatum{arcs.Arc{StartAngle: arcs.DegToRad(5), EndAngle: arcs.DegToRad(10)}},
	}
	out := arcs.FilterDataBySkipAngle(data, 0)
	if len(out) != 2 {
		t.Fatalf("skipAngle=0 should keep all %d arcs, got %d", len(data), len(out))
	}
	// The result is a copy, not the same backing slice.
	out[0] = testDatum{arcs.Arc{EndAngle: 1}}
	if data[0].GetArc().EndAngle == 1.0 {
		t.Errorf("FilterDataBySkipAngle should return a copy")
	}
}

func TestFindArcUnderCursor_DonutHole(t *testing.T) {
	ring := []arcs.Arc{{StartAngle: 0, EndAngle: 2 * math.Pi, InnerRadius: 5, OuterRadius: 10}}
	if idx := arcs.FindArcUnderCursor(ring, 0, 0, 1, 1); idx != -1 {
		t.Errorf("cursor in donut hole got idx %d, want -1", idx)
	}
	if idx := arcs.FindArcUnderCursor(ring, 0, 0, 0, -7); idx != 0 {
		t.Errorf("cursor on ring got idx %d, want 0", idx)
	}
}

func TestFindArcUnderCursor_ReversedAngles(t *testing.T) {
	// Start > End: the search normalizes so lo <= hi.
	a := []arcs.Arc{{StartAngle: math.Pi / 2, EndAngle: 0, InnerRadius: 0, OuterRadius: 10}}
	if idx := arcs.FindArcUnderCursor(a, 0, 0, 5, -5); idx != 0 {
		t.Errorf("reversed-angle arc got idx %d, want 0", idx)
	}
}

func TestFindArcUnderCursor_NegativeStartAngle(t *testing.T) {
	// Arc spanning [-π/2, π/2] (left-top-right through 0). A cursor on the
	// far left side (cursor angle ≈ 3π/2) only matches after unwrapping by
	// -2π.
	a := []arcs.Arc{{StartAngle: -math.Pi / 2, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10}}
	if idx := arcs.FindArcUnderCursor(a, 0, 0, -5, -0.5); idx != 0 {
		t.Errorf("unwrapped cursor got idx %d, want 0", idx)
	}
	// A point in the bottom half is genuinely outside the arc.
	if idx := arcs.FindArcUnderCursor(a, 0, 0, 0, 5); idx != -1 {
		t.Errorf("bottom cursor got idx %d, want -1", idx)
	}
}

func TestFindArcUnderCursor_OffsetCenter(t *testing.T) {
	// Same top-right quarter, but the pie center sits at (100, 100).
	a := []arcs.Arc{{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10}}
	if idx := arcs.FindArcUnderCursor(a, 100, 100, 105, 95); idx != 0 {
		t.Errorf("offset-center cursor got idx %d, want 0", idx)
	}
	if idx := arcs.FindArcUnderCursor(a, 100, 100, 5, -5); idx != -1 {
		t.Errorf("cursor near origin is far from the offset center, got %d", idx)
	}
}

func TestComputeArcBoundingBox_ExcludeCenter(t *testing.T) {
	// A narrow arc from 80° to 100° at radius 10 hugs the right edge. Without
	// the center the box is a thin sliver; including the origin stretches it
	// back to x=0.
	x, _, w, _ := arcs.ComputeArcBoundingBox(0, 0, 10, 80, 100, false)
	if x < 9 {
		t.Errorf("sliver min x = %v, want close to the arc (>= 9)", x)
	}
	if w > 1.5 {
		t.Errorf("sliver width = %v, want < 1.5", w)
	}
	xc, _, wc, _ := arcs.ComputeArcBoundingBox(0, 0, 10, 80, 100, true)
	approx(t, "with-center min x", xc, 0)
	if wc <= w {
		t.Errorf("including center should widen the box: %v <= %v", wc, w)
	}
}

func TestComputeArcBoundingBox_ReversedAngles(t *testing.T) {
	// Swapped start/end must yield the same box as the forward quarter.
	x1, y1, w1, h1 := arcs.ComputeArcBoundingBox(0, 0, 10, 0, 90, true)
	x2, y2, w2, h2 := arcs.ComputeArcBoundingBox(0, 0, 10, 90, 0, true)
	approx(t, "x", x2, x1)
	approx(t, "y", y2, y1)
	approx(t, "w", w2, w1)
	approx(t, "h", h2, h1)
}

func TestFindArcUnderCursor_WrapAroundArc(t *testing.T) {
	// Arc spanning [3pi/2, 5pi/2] wraps past 12 o'clock. A cursor slightly
	// right of the top (cursor angle ~0.2 rad) only matches after wrapping
	// the cursor angle forward by 2pi.
	a := []arcs.Arc{{StartAngle: 3 * math.Pi / 2, EndAngle: 5 * math.Pi / 2, InnerRadius: 0, OuterRadius: 10}}
	if idx := arcs.FindArcUnderCursor(a, 0, 0, 1, -4.9); idx != 0 {
		t.Errorf("wrap-around cursor got idx %d, want 0", idx)
	}
	// The bottom half is outside the wrapped arc.
	if idx := arcs.FindArcUnderCursor(a, 0, 0, 0, 5); idx != -1 {
		t.Errorf("bottom cursor got idx %d, want -1", idx)
	}
}
