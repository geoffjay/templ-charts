package polaraxes

import (
	"math"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func TestFmtN_EdgeCases(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{1.5, "1.5"},
		{10, "10"},
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

func TestStrokeFromExtra_Branches(t *testing.T) {
	if got := strokeFromExtra(nil); got != "" {
		t.Errorf("nil map = %q, want empty", got)
	}
	if got := strokeFromExtra(map[string]any{"other": "x"}); got != "" {
		t.Errorf("missing key = %q, want empty", got)
	}
	if got := strokeFromExtra(map[string]any{"stroke": "#123"}); got != "#123" {
		t.Errorf("string = %q, want #123", got)
	}
	if got := strokeFromExtra(map[string]any{"stroke": 42}); got != "42" {
		t.Errorf("non-string = %q, want 42", got)
	}
}

func TestStrokeWidthFromExtra_Branches(t *testing.T) {
	if got := strokeWidthFromExtra(nil); got != 0 {
		t.Errorf("nil map = %v, want 0", got)
	}
	if got := strokeWidthFromExtra(map[string]any{"strokeWidth": 1.5}); got != 1.5 {
		t.Errorf("float = %v, want 1.5", got)
	}
	if got := strokeWidthFromExtra(map[string]any{"strokeWidth": 2}); got != 2 {
		t.Errorf("int = %v, want 2", got)
	}
	if got := strokeWidthFromExtra(map[string]any{"strokeWidth": "x"}); got != 0 {
		t.Errorf("bad type = %v, want 0", got)
	}
	if got := strokeWidthFromExtra(map[string]any{"other": 1}); got != 0 {
		t.Errorf("missing key = %v, want 0", got)
	}
}

func TestResolveThemes_Nil(t *testing.T) {
	got := resolveAxisTheme(nil, nil)
	if got.Domain.Line.Extra != nil || got.Ticks.Line.Extra != nil || got.Legend.Text.Fill != "" {
		t.Errorf("nil theme should yield zero axis theme: %+v", got)
	}
	grid := resolveGridTheme(nil)
	if grid.Line.Extra != nil {
		t.Errorf("nil theme should yield zero grid theme")
	}
	// Non-nil path passes through the theme's blocks.
	axis := resolveAxisTheme(&theming.DefaultTheme, nil)
	if strokeFromExtra(axis.Ticks.Line.Extra) != "#777777" {
		t.Errorf("default theme ticks line stroke not resolved")
	}
	grid = resolveGridTheme(&theming.DefaultTheme)
	if strokeFromExtra(grid.Line.Extra) != "#dddddd" {
		t.Errorf("default theme grid line stroke not resolved")
	}
}

func TestTickValues_SpecHandling(t *testing.T) {
	scale := scales.NewLinearScaleWithRange(0, 100, 0, 100)
	// Explicit values pass through untouched.
	vals := tickValues(scale, scales.TicksSpec{HasValues: true, Values: []any{1.0, 2.0}})
	if len(vals) != 2 {
		t.Fatalf("explicit values len = %d, want 2", len(vals))
	}
	// Empty spec defaults to count 10.
	vals = tickValues(scale, scales.TicksSpec{})
	if len(vals) == 0 {
		t.Fatalf("default spec produced no ticks")
	}
	// Explicit count honored (linear ticks are approximate but non-empty and
	// no larger than a generous bound).
	vals = tickValues(scale, scales.TicksSpec{Count: 2, HasCount: true})
	if len(vals) == 0 || len(vals) > 5 {
		t.Fatalf("count=2 ticks len = %d, want small non-zero", len(vals))
	}
}

func TestTickLabel_Formats(t *testing.T) {
	if got := tickLabel(1.5, nil); got != "1.5" {
		t.Errorf("nil format = %q, want 1.5", got)
	}
	if got := tickLabel(1.5, func(v any) string { return "X" }); got != "X" {
		t.Errorf("func format = %q, want X", got)
	}
}

func TestArcPath_LargeArcFlag(t *testing.T) {
	if p := arcPath(10, 0, 90); !strings.Contains(p, " 0 0 1 ") {
		t.Errorf("quarter arc should use large-arc-flag 0: %s", p)
	}
	if p := arcPath(10, 0, 270); !strings.Contains(p, " 0 1 1 ") {
		t.Errorf("3/4 arc should use large-arc-flag 1: %s", p)
	}
}

func TestLinePositionsAndTextPosition(t *testing.T) {
	// Angle 0 (already offset) points along +x.
	x1, y1, x2, y2 := linePositions(0, 10, 20)
	if x1 != 10 || y1 != 0 || x2 != 20 || y2 != 0 {
		t.Errorf("linePositions(0,10,20) = (%v,%v,%v,%v), want (10,0,20,0)", x1, y1, x2, y2)
	}
	// Angle -90 points up (negative y).
	tx, ty := textPosition(-90, 30)
	if math.Abs(tx) > 1e-9 || math.Abs(ty+30) > 1e-9 {
		t.Errorf("textPosition(-90,30) = (%v,%v), want (0,-30)", tx, ty)
	}
}

func TestAngleScaleToTicks_BandCentered(t *testing.T) {
	scale := scales.NewBandScaleWithRange([]string{"a", "b"}, 0, 360, 0, false)
	ticks := angleScaleToTicks(scale)
	if len(ticks) != 2 {
		t.Fatalf("tick count = %d, want 2", len(ticks))
	}
	// Bands of 180 → centers 90 and 270, minus the 90 offset → 0 and 180.
	if ticks[0].angle != 0 || ticks[1].angle != 180 {
		t.Errorf("band angles = %v, %v, want 0, 180", ticks[0].angle, ticks[1].angle)
	}
}

func TestRadialTickPositions_Linear(t *testing.T) {
	scale := scales.NewLinearScaleWithRange(0, 10, 0, 100)
	positions := radialTickPositions(scale, scales.TicksSpec{HasValues: true, Values: []any{0.0, 5.0, 10.0}})
	if len(positions) != 3 {
		t.Fatalf("positions len = %d, want 3", len(positions))
	}
	if positions[1].position != 50 {
		t.Errorf("mid position = %v, want 50", positions[1].position)
	}
}

func TestAnimateRotate_Scaffold(t *testing.T) {
	if got := animateRotate(true, 45); got != "" {
		t.Errorf("animateRotate is a scaffold and should emit nothing, got %q", got)
	}
}
