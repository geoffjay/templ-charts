package arcs

import "testing"

// Internal-package tests for the small unexported helpers.

func TestFmtF(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{1.5, "1.5"},
		{2, "2"},
		{1.2345678, "1.235"}, // 3 dp
		{10.100, "10.1"},     // trailing zeros trimmed
		{-3.25, "-3.25"},
		{0, "0"},
		{-0.0001, "0"}, // tiny negatives normalize to "0", not "-0"
	}
	for _, c := range cases {
		if got := fmtF(c.in); got != c.want {
			t.Errorf("fmtF(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExtent(t *testing.T) {
	if min, max := extent(nil); min != 0 || max != 0 {
		t.Errorf("extent(nil) = (%v, %v), want (0, 0)", min, max)
	}
	if min, max := extent([]float64{3, -1, 7, 2}); min != -1 || max != 7 {
		t.Errorf("extent = (%v, %v), want (-1, 7)", min, max)
	}
	if min, max := extent([]float64{5}); min != 5 || max != 5 {
		t.Errorf("extent single = (%v, %v), want (5, 5)", min, max)
	}
}
