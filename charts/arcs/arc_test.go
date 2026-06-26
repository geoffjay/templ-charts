package arcs

import (
	"math"
	"testing"
)

func TestDegToRad(t *testing.T) {
	if got := DegToRad(180); got != math.Pi {
		t.Fatalf("DegToRad(180)=%v, want π", got)
	}
}

func TestGetNormalizedAngle(t *testing.T) {
	if got := GetNormalizedAngle(-10); got != 350 {
		t.Fatalf("got %v, want 350", got)
	}
	if got := GetNormalizedAngle(370); got != 10 {
		t.Fatalf("got %v, want 10", got)
	}
}

func TestArcGenerator_FullCircle(t *testing.T) {
	g := CreateArcGenerator(0, 0)
	arc := Arc{StartAngle: 0, EndAngle: 2 * math.Pi, InnerRadius: 0, OuterRadius: 10}
	path := g.GenerateSvgArc(arc)
	if path == "" {
		t.Fatalf("expected non-empty path for full circle")
	}
	// A full circle path should contain two arc commands (outer + inner).
	if !containsStr(path, "A") {
		t.Fatalf("expected arc command in %s", path)
	}
}

func TestArcGenerator_Donut(t *testing.T) {
	g := CreateArcGenerator(0, 0)
	arc := Arc{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 5, OuterRadius: 10}
	path := g.GenerateSvgArc(arc)
	if path == "" {
		t.Fatalf("expected non-empty path for donut quarter")
	}
}

func TestComputeArcCenter(t *testing.T) {
	// Arc spanning 0 to π/2 (top to right). Mid-angle (d3) = π/4 - π/2 = -π/4.
	// With radius offset 0.5 and inner 0, outer 10: r = 5.
	arc := Arc{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10}
	x, y := ComputeArcCenter(arc, 0.5)
	// cos(-π/4)*5 ≈ 3.5355, sin(-π/4)*5 ≈ -3.5355
	if math.Abs(x-3.5355) > 0.01 || math.Abs(y-(-3.5355)) > 0.01 {
		t.Fatalf("center = (%.4f, %.4f), want (~3.5355, ~-3.5355)", x, y)
	}
}

func TestComputeArcBoundingBox_Quarter(t *testing.T) {
	// Quarter from 0° (top) to 90° (right) at radius 10, centered at (0,0).
	// The arc spans the top-right quadrant. Bounding box should be roughly
	// x∈[0, ~10], y∈[~-10, ~0] — includeCenter adds (0,0).
	x, y, w, h := ComputeArcBoundingBox(0, 0, 10, 0, 90, true)
	if x > 0.01 || y > 0.01 {
		t.Fatalf("min corner = (%.3f, %.3f), want <= 0", x, y)
	}
	if w < 9.9 || h < 9.9 {
		t.Fatalf("bbox w,h = (%.3f, %.3f), want ~10 each", w, h)
	}
}

func TestComputeArcBoundingBox_FullCircle(t *testing.T) {
	x, y, w, h := ComputeArcBoundingBox(5, 5, 10, 0, 360, true)
	// Full circle bbox: [-5, 15] in both axes → w=h=20.
	if math.Abs(x-(-5)) > 0.01 || math.Abs(y-(-5)) > 0.01 {
		t.Fatalf("bbox min = (%.3f, %.3f), want (-5, -5)", x, y)
	}
	if math.Abs(w-20) > 0.01 || math.Abs(h-20) > 0.01 {
		t.Fatalf("bbox w,h = (%.3f, %.3f), want (20, 20)", w, h)
	}
}

func TestComputeArcLink(t *testing.T) {
	// Right-side arc: user mid-angle ~0 (top) to π/2 (right) → mid π/4 →
	// after -π/2 offset = -π/4, cos > 0 → right.
	arc := Arc{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10}
	link := ComputeArcLink(arc, 0, 16, 24)
	if link.Side != "right" || link.TextAnchor != "start" {
		t.Fatalf("right-side arc link got side=%s anchor=%s", link.Side, link.TextAnchor)
	}
	if link.P2x <= link.P1x {
		t.Fatalf("right-side P2 (%.3f) should be > P1 (%.3f)", link.P2x, link.P1x)
	}
	// Left-side arc: user mid-angle 3π/2 (left, since 0=top clockwise).
	arc2 := Arc{StartAngle: math.Pi, EndAngle: 2 * math.Pi, InnerRadius: 0, OuterRadius: 10}
	link2 := ComputeArcLink(arc2, 0, 16, 24)
	if link2.Side != "left" || link2.TextAnchor != "end" {
		t.Fatalf("left-side arc link got side=%s anchor=%s", link2.Side, link2.TextAnchor)
	}
}

func TestComputeArcLinkTextAnchor(t *testing.T) {
	right := Arc{StartAngle: 0, EndAngle: math.Pi / 2} // mid π/4 → right
	if got := ComputeArcLinkTextAnchor(right); got != "start" {
		t.Fatalf("right arc anchor = %q, want start", got)
	}
	left := Arc{StartAngle: math.Pi, EndAngle: 2 * math.Pi} // mid 3π/2 → left
	if got := ComputeArcLinkTextAnchor(left); got != "end" {
		t.Fatalf("left arc anchor = %q, want end", got)
	}
}

func TestFindArcUnderCursor(t *testing.T) {
	arcs := []Arc{
		{StartAngle: 0, EndAngle: math.Pi / 2, InnerRadius: 0, OuterRadius: 10},       // top-right
		{StartAngle: math.Pi / 2, EndAngle: math.Pi, InnerRadius: 0, OuterRadius: 10}, // bottom-right
	}
	// Cursor in the top-right quadrant (x>0, y<0 in screen coords, i.e. y=-5).
	if idx := FindArcUnderCursor(arcs, 0, 0, 5, -5); idx != 0 {
		t.Fatalf("cursor top-right got idx %d, want 0", idx)
	}
	// Cursor in the bottom-right (x=5, y=5).
	if idx := FindArcUnderCursor(arcs, 0, 0, 5, 5); idx != 1 {
		t.Fatalf("cursor bottom-right got idx %d, want 1", idx)
	}
	// Cursor outside radius.
	if idx := FindArcUnderCursor(arcs, 0, 0, 20, 20); idx != -1 {
		t.Fatalf("cursor outside got idx %d, want -1", idx)
	}
}

func TestFilterDataBySkipAngle(t *testing.T) {
	data := []DatumWithArc{
		arcDatum{Arc{StartAngle: 0, EndAngle: DegToRad(45)}},            // 45° span
		arcDatum{Arc{StartAngle: DegToRad(45), EndAngle: DegToRad(50)}}, // 5° span
	}
	out := FilterDataBySkipAngle(data, 10)
	if len(out) != 1 {
		t.Fatalf("got %d arcs, want 1 (only the 45° arc)", len(out))
	}
}

type arcDatum struct{ a Arc }

func (d arcDatum) GetArc() Arc { return d.a }

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
