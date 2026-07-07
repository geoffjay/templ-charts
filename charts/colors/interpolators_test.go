package colors

import "testing"

func TestInterpolateRgbBasisEdgeCases(t *testing.T) {
	// Empty stop list → black.
	if got := interpolateRgbBasis(nil)(0.5); got != "#000000" {
		t.Errorf("empty stops = %s, want #000000", got)
	}
	// Single stop → constant.
	f := interpolateRgbBasis([]string{"#123456"})
	for _, tt := range []float64{0, 0.5, 1} {
		if got := f(tt); got != "#123456" {
			t.Errorf("single stop f(%v) = %s, want #123456", tt, got)
		}
	}
	// Endpoints are exact; out-of-range t clamps to them.
	f = interpolateRgbBasis([]string{"#ff0000", "#0000ff"})
	if got := f(-1); got != "#ff0000" {
		t.Errorf("f(-1) = %s, want first stop", got)
	}
	if got := f(2); got != "#0000ff" {
		t.Errorf("f(2) = %s, want last stop", got)
	}
	if got := f(0.5); !hexRe.MatchString(got) {
		t.Errorf("f(0.5) = %q, not a hex color", got)
	}
}

func TestParseHexRGB(t *testing.T) {
	r, g, b := parseHexRGB("#ff8000")
	if r != 1 || !approxEqualF(g, 128.0/255) || b != 0 {
		t.Errorf("parseHexRGB(#ff8000) = %v %v %v", r, g, b)
	}
	// Malformed inputs fall back to black.
	for _, s := range []string{"", "red", "#fff", "ff8000"} {
		r, g, b := parseHexRGB(s)
		if r != 0 || g != 0 || b != 0 {
			t.Errorf("parseHexRGB(%q) = %v %v %v, want zeros", s, r, g, b)
		}
	}
}

func TestRgbToHexClamps(t *testing.T) {
	if got := rgbToHex(-1, 2, 128.0/255); got != "#00ff80" {
		t.Errorf("rgbToHex(-1, 2, 0.5) = %s, want #00ff80", got)
	}
	if got := rgbToHex(0, 0, 0); got != "#000000" {
		t.Errorf("rgbToHex(0,0,0) = %s", got)
	}
	if got := rgbToHex(1, 1, 1); got != "#ffffff" {
		t.Errorf("rgbToHex(1,1,1) = %s", got)
	}
}

func TestHex2Clamps(t *testing.T) {
	if got := hex2(-5); got != "00" {
		t.Errorf("hex2(-5) = %q, want 00", got)
	}
	if got := hex2(300); got != "ff" {
		t.Errorf("hex2(300) = %q, want ff", got)
	}
	if got := hex2(171); got != "ab" {
		t.Errorf("hex2(171) = %q, want ab", got)
	}
}

func TestNewOrdinalScaleEmptyColors(t *testing.T) {
	// An empty color list degrades to empty strings rather than panicking.
	s := GetOrdinalColorScale[string](OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: nil}, "")
	if got := s("a"); got != "" {
		t.Errorf("empty-colors scale = %q, want empty", got)
	}
}

func TestGetQuantizeColorScale_DefaultStepsWithScheme(t *testing.T) {
	// Steps <= 0 falls back to the default step count of 7.
	s := GetQuantizeColorScale(QuantizeColorScaleConfig{Scheme: "viridis", Steps: 0}, SequentialColorScaleValues{Min: 0, Max: 7})
	if got, want := s(-1), ColorInterpolators["viridis"](0); got != want {
		t.Errorf("first bucket = %s, want viridis(0) %s", got, want)
	}
	// An unknown scheme falls back to turbo.
	s = GetQuantizeColorScale(QuantizeColorScaleConfig{Scheme: "bogus"}, SequentialColorScaleValues{Min: 0, Max: 1})
	if got, want := s(-1), ColorInterpolators["turbo"](0); got != want {
		t.Errorf("bogus scheme first bucket = %s, want turbo(0) %s", got, want)
	}
}
