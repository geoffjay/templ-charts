package rects

import "testing"

func TestBuildRoundedRectPath_Square(t *testing.T) {
	got := BuildRoundedRectPath(0, 0, 10, 10, 0, 0, 0, 0)
	want := "M0,0H10V10H0V0Z"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestBuildRoundedRectPath_AllCorners(t *testing.T) {
	got := BuildRoundedRectPath(0, 0, 10, 10, 2, 2, 2, 2)
	// Must start at (x+tl, y) and contain 4 arc commands.
	if want := "M2,0"; got[:len(want)] != want {
		t.Fatalf("path should start at %q, got %q", want, got)
	}
	for _, sub := range []string{"A2,2 0 0 1", "A2,2 0 0 1", "A2,2 0 0 1", "A2,2 0 0 1"} {
		if !contains(got, sub) {
			t.Fatalf("path %q missing arc %q", got, sub)
		}
	}
	if !endsWith(got, "Z") {
		t.Fatalf("path should end with Z, got %q", got)
	}
}

func TestBuildRoundedRectPath_Clamp(t *testing.T) {
	// Radius exceeding half the width should be clamped.
	got := BuildRoundedRectPath(0, 0, 4, 10, 5, 5, 5, 5)
	// Each radius clamped to w/2 = 2.
	for _, sub := range []string{"A2,2 0 0 1"} {
		if n := count(got, sub); n != 4 {
			t.Fatalf("expected 4 clamped arcs of %q, got %d in %q", sub, n, got)
		}
	}
}

func TestBorderRadius_Resolved(t *testing.T) {
	b := BorderRadiusFromFloat(8)
	c := b.Resolved(20, 20)
	if c.TopLeft != 8 || c.TopRight != 8 || c.BottomLeft != 8 || c.BottomRight != 8 {
		t.Fatalf("uniform resolved got %+v", c)
	}
	// Clamped to half-min side.
	c = b.Resolved(10, 10)
	if c.TopLeft != 5 {
		t.Fatalf("expected clamp to 5, got %v", c.TopLeft)
	}
}

func TestRectLabelPosition_Vertical(t *testing.T) {
	r := Rect{X: 0, Y: 0, Width: 10, Height: 20}
	p := RectLabelPosition(r, AnchorTop, 4, false)
	if p.Y != -4 || p.TextAnchor != "middle" || p.DominantBaseline != "auto" {
		t.Fatalf("top anchor got %+v", p)
	}
	p = RectLabelPosition(r, AnchorBottom, 4, false)
	if p.Y != 24 || p.DominantBaseline != "hanging" {
		t.Fatalf("bottom anchor got %+v", p)
	}
	p = RectLabelPosition(r, AnchorMiddle, 0, false)
	if p.Y != 10 || p.DominantBaseline != "central" {
		t.Fatalf("middle anchor got %+v", p)
	}
}

func TestRectLabelPosition_Horizontal(t *testing.T) {
	r := Rect{X: 0, Y: 0, Width: 30, Height: 10}
	p := RectLabelPosition(r, AnchorTop, 4, true)
	if p.X != 34 || p.TextAnchor != "middle" {
		t.Fatalf("horizontal top got %+v", p)
	}
	p = RectLabelPosition(r, AnchorBottom, 4, true)
	if p.X != -4 {
		t.Fatalf("horizontal bottom got %+v", p)
	}
}

func TestRender_NoAnimate(t *testing.T) {
	s := Render(RoundedRectProps{X: 0, Y: 0, Width: 10, Height: 10, Fill: "#f00"})
	if !contains(s, `<path d="`) || !contains(s, `fill="#f00"`) {
		t.Fatalf("render output missing expected attrs: %s", s)
	}
	if contains(s, "<animate") {
		t.Fatalf("animate should be absent when Animate=false: %s", s)
	}
}

func TestRender_Animate(t *testing.T) {
	s := Render(RoundedRectProps{X: 0, Y: 0, Width: 10, Height: 10, Fill: "#f00", Animate: true})
	if !contains(s, "<animate") {
		t.Fatalf("expected <animate> in %s", s)
	}
	if !contains(s, `attributeName="height"`) {
		t.Fatalf("vertical animate should target height: %s", s)
	}
	s = Render(RoundedRectProps{X: 0, Y: 0, Width: 10, Height: 10, Fill: "#f00", Animate: true, Horizontal: true})
	if !contains(s, `attributeName="width"`) {
		t.Fatalf("horizontal animate should target width: %s", s)
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && indexOf(s, sub) >= 0 }
func endsWith(s, sub string) bool { return len(s) >= len(sub) && s[len(s)-len(sub):] == sub }
func count(s, sub string) int {
	n := 0
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			n++
			i += len(sub) - 1
		}
	}
	return n
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
