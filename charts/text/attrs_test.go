package text

import "testing"

func TestSvgTextAttrsFrom(t *testing.T) {
	a := SvgTextAttrsFrom("center", "top")
	if a.TextAnchor != "middle" || a.DominantBaseline != "text-before-edge" {
		t.Fatalf("got %+v, want anchor=middle baseline=text-before-edge", a)
	}
	a = SvgTextAttrsFrom("start", "bottom")
	if a.TextAnchor != "start" || a.DominantBaseline != "text-after-edge" {
		t.Fatalf("got %+v", a)
	}
	a = SvgTextAttrsFrom("end", "center")
	if a.TextAnchor != "end" || a.DominantBaseline != "middle" {
		t.Fatalf("got %+v", a)
	}
}

func TestTruncateTickAt(t *testing.T) {
	if got := TruncateTickAt("January", 30, "0", 6); got == "January" {
		t.Fatalf("expected truncation, got %q", got)
	}
	if got := TruncateTickAt("Jan", 100, "0", 6); got != "Jan" {
		t.Fatalf("short label should be unchanged, got %q", got)
	}
	if got := TruncateTickAt("January", 30, "45", 6); got != "Janua…" {
		t.Fatalf("rotated truncation got %q, want Janua…", got)
	}
	if got := TruncateTickAt("January", 0, "0", 6); got != "January" {
		t.Fatalf("length<=0 should disable, got %q", got)
	}
}
