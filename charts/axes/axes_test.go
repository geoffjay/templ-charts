package axes

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func TestComputeCartesianTicks_Linear(t *testing.T) {
	data := scales.ComputedSerieAxis{All: []any{0.0, 10.0, 20.0}, Min: 0.0, Max: 20.0}
	scale := scales.ComputeScale(scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat(), Nice: true}, data, 200, scales.ScaleAxisX)
	props := AxisProps{Axis: "x", Scale: scale, Length: 200, TickSize: 5, TickPadding: 5, TicksPosition: "after"}
	computed := ComputeCartesianTicks(props, &theming.DefaultTheme)
	if len(computed.Ticks) == 0 {
		t.Fatalf("expected non-zero ticks")
	}
	// First tick should be at position 0 (min).
	if computed.Ticks[0].Position != 0 {
		t.Fatalf("first tick position = %v, want 0", computed.Ticks[0].Position)
	}
	// Label should be non-empty.
	if computed.Ticks[0].Label == "" {
		t.Fatalf("first tick label empty")
	}
	// Tick line: x-axis, ticksPosition="after" → y2 = +tickSize.
	if computed.Ticks[0].Y2 != 5 {
		t.Fatalf("tick y2 = %v, want 5", computed.Ticks[0].Y2)
	}
}

func TestComputeCartesianTicks_BandCentered(t *testing.T) {
	data := scales.ComputedSerieAxis{All: []any{"a", "b", "c"}}
	scale := scales.ComputeScale(scales.ScaleBandSpec{}, data, 300, scales.ScaleAxisX)
	props := AxisProps{Axis: "x", Scale: scale, Length: 300, TickSize: 5, TickPadding: 5}
	computed := ComputeCartesianTicks(props, &theming.DefaultTheme)
	if len(computed.Ticks) != 3 {
		t.Fatalf("expected 3 ticks, got %d", len(computed.Ticks))
	}
	// Band ticks should be centered (position > 0 for "a").
	if computed.Ticks[0].Position <= 0 {
		t.Fatalf("band tick 'a' position = %v, want > 0 (centered)", computed.Ticks[0].Position)
	}
	if computed.Ticks[0].Label != "a" {
		t.Fatalf("band tick label = %q, want 'a'", computed.Ticks[0].Label)
	}
}

func TestComputeCartesianTicks_YAxis(t *testing.T) {
	data := scales.ComputedSerieAxis{All: []any{0.0, 10.0}, Min: 0.0, Max: 10.0}
	scale := scales.ComputeScale(scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()}, data, 200, scales.ScaleAxisY)
	props := AxisProps{Axis: "y", Scale: scale, Length: 200, TickSize: 5, TickPadding: 5, TicksPosition: "after"}
	computed := ComputeCartesianTicks(props, &theming.DefaultTheme)
	if len(computed.Ticks) == 0 {
		t.Fatalf("no ticks")
	}
	// Y-axis: x2 = +tickSize (after).
	if computed.Ticks[0].X2 != 5 {
		t.Fatalf("y tick x2 = %v, want 5", computed.Ticks[0].X2)
	}
}

func TestComputeGridLines(t *testing.T) {
	data := scales.ComputedSerieAxis{All: []any{0.0, 5.0, 10.0}, Min: 0.0, Max: 10.0}
	scale := scales.ComputeScale(scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()}, data, 100, scales.ScaleAxisX)
	props := AxisProps{Axis: "x", Scale: scale}
	lines := ComputeGridLines(props, 100, 50)
	if len(lines) == 0 {
		t.Fatalf("expected grid lines")
	}
	// x-axis grid lines are vertical (y2 = height).
	if lines[0].Y2 != 50 {
		t.Fatalf("grid line y2 = %v, want 50", lines[0].Y2)
	}
}

func TestGetFormatter_Time(t *testing.T) {
	data := scales.ComputedSerieAxis{All: []any{}, Min: nil, Max: nil}
	scale := scales.ComputeScale(scales.ScaleTimeSpec{}, data, 100, scales.ScaleAxisX)
	f := GetFormatter(scale, &theming.DefaultTheme)
	if f == nil {
		t.Fatalf("nil formatter")
	}
}

func TestRotatedTextAttrs(t *testing.T) {
	if a, _ := rotatedTextAttrs(90, "x"); a != "middle" {
		t.Fatalf("90deg x anchor = %q, want middle", a)
	}
}

func TestLegendPosition_X(t *testing.T) {
	props := AxisProps{Axis: "x", Length: 200, LegendPosition: AxisLegendMiddle}
	x, y, rot, anchor := legendPosition(props)
	if x != 100 {
		t.Fatalf("middle x = %v, want 100", x)
	}
	if rot != 0 || anchor != "middle" {
		t.Fatalf("x legend rot=%v anchor=%q, want 0/middle", rot, anchor)
	}
	_ = y
}

func TestLegendPosition_Y(t *testing.T) {
	props := AxisProps{Axis: "y", Length: 200, LegendPosition: AxisLegendEnd, LegendOffset: 40}
	_, y, rot, anchor := legendPosition(props)
	if y != 200 {
		t.Fatalf("end y = %v, want 200", y)
	}
	if rot != -90 || anchor != "end" {
		t.Fatalf("y legend rot=%v anchor=%q, want -90/end", rot, anchor)
	}
}
