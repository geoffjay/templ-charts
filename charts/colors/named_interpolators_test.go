package colors

import "testing"

// TestNamedInterpolatorAnchors pins the named interpolators at t=0, 0.5, 1.
// turbo/cividis are polynomial (d3 formulas), viridis is a stepped lookup
// table, sinebow is trigonometric; the rest are RGB-basis ramps over their
// sampled gradients, so t=0/t=1 are the exact endpoint stops.
func TestNamedInterpolatorAnchors(t *testing.T) {
	cases := []struct {
		id             string
		at0, at05, at1 string
	}{
		{"turbo", "#23171b", "#95fb51", "#900c00"},
		{"cividis", "#002051", "#7f7c75", "#fdea45"},
		{"viridis", "#440154", "#21918c", "#ece51b"},
		{"inferno", "#000004", "#6b1fdd", "#fea265"},
		{"magma", "#000004", "#6b1fdd", "#fea265"},
		{"plasma", "#0d0887", "#300596", "#5003a2"},
		{"warm", "#4d1b3b", "#d724a8", "#ff42ff"},
		{"cool", "#3b8ca8", "#3b30a8", "#3b0076"},
		{"cubehelixDefault", "#1e1547", "#450bac", "#6c01ff"},
		{"sinebow", "#ff4040", "#00bfbf", "#ff4040"},
	}
	for _, c := range cases {
		f := ColorInterpolators[c.id]
		if f == nil {
			t.Errorf("no interpolator registered for %q", c.id)
			continue
		}
		if got := f(0); got != c.at0 {
			t.Errorf("%s(0) = %s, want %s", c.id, got, c.at0)
		}
		if got := f(0.5); got != c.at05 {
			t.Errorf("%s(0.5) = %s, want %s", c.id, got, c.at05)
		}
		if got := f(1); got != c.at1 {
			t.Errorf("%s(1) = %s, want %s", c.id, got, c.at1)
		}
	}
}

// TestSinebowExactValues verifies the trig formula by hand:
// t=0 → sin²(π/2)=1, sin²(π/2+π/3)=0.25, sin²(π/2+2π/3)=0.25 → #ff4040;
// t=0.5 → sin²(0)=0, sin²(π/3)=0.75, sin²(2π/3)=0.75 → #00bfbf.
// It is cyclical, so t=0 and t=1 agree.
func TestSinebowExactValues(t *testing.T) {
	if got := interpolateSinebow(0); got != "#ff4040" {
		t.Errorf("sinebow(0) = %s, want #ff4040", got)
	}
	if got := interpolateSinebow(0.5); got != "#00bfbf" {
		t.Errorf("sinebow(0.5) = %s, want #00bfbf", got)
	}
	if got := interpolateSinebow(1); got != interpolateSinebow(0) {
		t.Errorf("sinebow should be cyclical: f(1)=%s f(0)=%s", got, interpolateSinebow(0))
	}
}

// TestClampedInterpolators verifies out-of-range t is clamped to [0,1] for the
// polynomial interpolators.
func TestClampedInterpolators(t *testing.T) {
	for _, id := range []string{"turbo", "cividis"} {
		f := ColorInterpolators[id]
		if got, want := f(-2), f(0); got != want {
			t.Errorf("%s(-2) = %s, want clamp to %s", id, got, want)
		}
		if got, want := f(3), f(1); got != want {
			t.Errorf("%s(3) = %s, want clamp to %s", id, got, want)
		}
	}
}

// TestRainbowWraps verifies the cyclical rainbow interpolator wraps t outside
// [0,1] (wrap01) and produces well-formed colors.
func TestRainbowWraps(t *testing.T) {
	if a, b := interpolateRainbow(0.25), interpolateRainbow(1.25); a != b {
		t.Errorf("rainbow(1.25) = %s, want wrapped value %s", b, a)
	}
	if a, b := interpolateRainbow(0.75), interpolateRainbow(-0.25); a != b {
		t.Errorf("rainbow(-0.25) = %s, want wrapped value %s", b, a)
	}
	for _, tt := range []float64{0, 0.25, 0.5, 0.75, 1} {
		if got := interpolateRainbow(tt); !hexRe.MatchString(got) {
			t.Errorf("rainbow(%v) = %q, not a #rrggbb hex", tt, got)
		}
	}
}

func TestWrap01(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{0, 0}, {1, 1}, {0.25, 0.25}, {1.25, 0.25}, {-0.25, 0.75}, {2.5, 0.5},
	}
	for _, c := range cases {
		if got := wrap01(c.in); !approxEqualF(got, c.want) {
			t.Errorf("wrap01(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestClamp01(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{-1, 0}, {0, 0}, {0.4, 0.4}, {1, 1}, {2, 1},
	}
	for _, c := range cases {
		if got := clamp01(c.in); got != c.want {
			t.Errorf("clamp01(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// TestStepRamp verifies the discrete d3-style ramp: range[floor(t*n)] clamped.
func TestStepRamp(t *testing.T) {
	cols := []string{"#111111", "#222222", "#333333", "#444444"}
	cases := []struct {
		t    float64
		want string
	}{
		{-1, "#111111"}, // clamped low
		{0, "#111111"},  // first bucket
		{0.24, "#111111"},
		{0.25, "#222222"}, // bucket boundary
		{0.5, "#333333"},
		{0.99, "#444444"},
		{1, "#444444"}, // floor(1*4)=4 clamped to last
		{2, "#444444"}, // clamped high
	}
	for _, c := range cases {
		if got := stepRamp(cols, c.t); got != c.want {
			t.Errorf("stepRamp(%v) = %s, want %s", c.t, got, c.want)
		}
	}
	if got := stepRamp(nil, 0.5); got != "#000000" {
		t.Errorf("stepRamp(empty) = %s, want #000000", got)
	}
}

func approxEqualF(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 1e-9
}
