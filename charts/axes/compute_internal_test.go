package axes

import (
	"math"
	"testing"
)

func TestFmtN(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{1.5, "1.5"},
		{10, "10"},
		{1.2345, "1.234"},
		{-2.5, "-2.5"},
		{-0.0001, "0"}, // -0 normalization
		{math.NaN(), "0"},
		{math.Inf(1), "0"},
		{math.Inf(-1), "0"},
	}
	for _, tc := range tests {
		if got := fmtN(tc.in); got != tc.want {
			t.Errorf("fmtN(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFontSizeAny(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{nil, ""},
		{12.5, "12.5"},
		{11, "11"},
		{"13px", "13px"},
		{float32(9), "9"}, // default %v branch
	}
	for _, tc := range tests {
		if got := fontSizeAny(tc.in); got != tc.want {
			t.Errorf("fontSizeAny(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestRotatedTextAttrs_Branches(t *testing.T) {
	tests := []struct {
		rotation     float64
		axis         string
		wantAnchor   string
		wantBaseline string
	}{
		{90, "x", "middle", "hanging"},
		{270, "x", "middle", "auto"},
		{45, "x", "end", "auto"},
		{-45, "x", "start", "auto"}, // normalized to 315 (> 270)
		{135, "x", "start", "auto"},
		{45, "y", "middle", "middle"}, // vertical axis fallback
	}
	for _, tc := range tests {
		anchor, baseline := rotatedTextAttrs(tc.rotation, tc.axis)
		if anchor != tc.wantAnchor || baseline != tc.wantBaseline {
			t.Errorf("rotatedTextAttrs(%v, %q) = (%q, %q), want (%q, %q)",
				tc.rotation, tc.axis, anchor, baseline, tc.wantAnchor, tc.wantBaseline)
		}
	}
}

func TestStringFromExtra(t *testing.T) {
	if got := stringFromExtra(nil, "stroke"); got != "" {
		t.Errorf("nil map = %q, want empty", got)
	}
	if got := stringFromExtra(map[string]any{"other": "x"}, "stroke"); got != "" {
		t.Errorf("missing key = %q, want empty", got)
	}
	if got := stringFromExtra(map[string]any{"stroke": "#fff"}, "stroke"); got != "#fff" {
		t.Errorf("string value = %q, want #fff", got)
	}
	if got := stringFromExtra(map[string]any{"stroke": 7}, "stroke"); got != "7" {
		t.Errorf("non-string value = %q, want 7", got)
	}
}

func TestFloatFromExtra(t *testing.T) {
	if got := floatFromExtra(nil, "strokeWidth"); got != 0 {
		t.Errorf("nil map = %v, want 0", got)
	}
	if got := floatFromExtra(map[string]any{"strokeWidth": 1.5}, "strokeWidth"); got != 1.5 {
		t.Errorf("float value = %v, want 1.5", got)
	}
	if got := floatFromExtra(map[string]any{"strokeWidth": 2}, "strokeWidth"); got != 2 {
		t.Errorf("int value = %v, want 2", got)
	}
	if got := floatFromExtra(map[string]any{"strokeWidth": "x"}, "strokeWidth"); got != 0 {
		t.Errorf("bad type = %v, want 0", got)
	}
	if got := floatFromExtra(map[string]any{"other": 1}, "strokeWidth"); got != 0 {
		t.Errorf("missing key = %v, want 0", got)
	}
}

func TestTickTransform(t *testing.T) {
	tick := Tick{X1: 10, Y1: 0, LabelX: 10, LabelY: 12}
	if got := tickTransform(AxisProps{}, tick); got != "translate(0,12)" {
		t.Errorf("no rotation transform = %q, want translate(0,12)", got)
	}
	if got := tickTransform(AxisProps{TickRotation: 45}, tick); got != "translate(0,12) rotate(45)" {
		t.Errorf("rotated transform = %q", got)
	}
}
