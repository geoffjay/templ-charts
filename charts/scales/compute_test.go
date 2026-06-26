package scales

import (
	"testing"
	"time"
)

func TestComputeScale_Linear(t *testing.T) {
	data := ComputedSerieAxis{All: []any{0.0, 10.0, 20.0}, Min: 0.0, Max: 20.0}
	spec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Nice: true}
	s := ComputeScale(spec, data, 200, ScaleAxisY)
	if got := s.Call(0); got != 200 {
		t.Fatalf("linear(0) = %v, want 200 (y axis inverted)", got)
	}
	if got := s.Call(20); got != 0 {
		t.Fatalf("linear(20) = %v, want 0", got)
	}
}

func TestComputeScale_Band(t *testing.T) {
	data := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	s := ComputeScale(ScaleBandSpec{}, data, 300, ScaleAxisX)
	// band should produce increasing positions.
	a := s.Call("a")
	c := s.Call("c")
	if a >= c {
		t.Fatalf("band(a)=%v >= band(c)=%v, want a<c", a, c)
	}
	if s.Type() != ScaleTypeBand {
		t.Fatalf("type = %v, want band", s.Type())
	}
}

func TestComputeScale_Point(t *testing.T) {
	data := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	s := ComputeScale(ScalePointSpec{}, data, 300, ScaleAxisX)
	a := s.Call("a")
	c := s.Call("c")
	if a >= c {
		t.Fatalf("point(a)=%v >= point(c)=%v", a, c)
	}
	bw := s.(ScaleWithBandwidth).Bandwidth()
	if bw != 0 {
		t.Fatalf("point bandwidth = %v, want 0", bw)
	}
}

func TestComputeXYScalesForSeries_LinearLinear(t *testing.T) {
	series := []Serie{
		{Data: []SerieDatum{{X: 0.0, Y: 1.0}, {X: 1.0, Y: 2.0}, {X: 2.0, Y: 3.0}}},
	}
	xSpec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat()}
	ySpec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat()}
	res := ComputeXYScalesForSeries(series, xSpec, ySpec, 100, 200)
	if len(res.Series) != 1 {
		t.Fatalf("series len = %d, want 1", len(res.Series))
	}
	if len(res.Series[0].Data) != 3 {
		t.Fatalf("data len = %d, want 3", len(res.Series[0].Data))
	}
	d0 := res.Series[0].Data[0]
	if d0.Position.X == nil || d0.Position.Y == nil {
		t.Fatal("expected non-nil positions")
	}
	if *d0.Position.X != 0 {
		t.Fatalf("x[0] = %v, want 0", *d0.Position.X)
	}
	// y axis inverted: y=1 → 200 - 1/3*200 = 133.33
	if *d0.Position.Y <= 0 || *d0.Position.Y >= 200 {
		t.Fatalf("y[0] = %v, want within (0, 200)", *d0.Position.Y)
	}
}

func TestComputeXYScalesForSeries_StackedY(t *testing.T) {
	series := []Serie{
		{Data: []SerieDatum{{X: "a", Y: 1.0}, {X: "b", Y: 2.0}}},
		{Data: []SerieDatum{{X: "a", Y: 3.0}, {X: "b", Y: 4.0}}},
	}
	xSpec := ScaleBandSpec{}
	ySpec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Stacked: true}
	res := ComputeXYScalesForSeries(series, xSpec, ySpec, 100, 200)
	if res.Y.MinStacked == nil || res.Y.MaxStacked == nil {
		t.Fatal("expected stacked min/max")
	}
	// stack at "a": 1 + 3 = 4; at "b": 2 + 4 = 6 → max 6
	if *res.Y.MaxStacked != 6 {
		t.Fatalf("y maxStacked = %v, want 6", *res.Y.MaxStacked)
	}
	if res.Series[1].Data[0].YStacked == nil {
		t.Fatal("expected YStacked on second series")
	}
	if *res.Series[1].Data[0].YStacked != 4 {
		t.Fatalf("yStacked[1][a] = %v, want 4", *res.Series[1].Data[0].YStacked)
	}
}

func TestCenterScale_Band(t *testing.T) {
	data := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	s := ComputeScale(ScaleBandSpec{}, data, 300, ScaleAxisX)
	centered := CenterScale(s)
	a := s.Call("a")
	ca := centered("a")
	bw := s.(ScaleWithBandwidth).Bandwidth()
	if ca != a+bw/2 {
		t.Fatalf("centered(a) = %v, want %v", ca, a+bw/2)
	}
}

func TestGetScaleTicks_Linear(t *testing.T) {
	data := ComputedSerieAxis{All: []any{0.0, 100.0}, Min: 0.0, Max: 100.0}
	s := ComputeScale(ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Nice: true}, data, 200, ScaleAxisX)
	ticks := GetScaleTicks(s, TicksSpec{Count: 5, HasCount: true})
	if len(ticks) < 2 {
		t.Fatalf("ticks len = %d, want >= 2", len(ticks))
	}
}

func TestCreateDateNormalizer_Native(t *testing.T) {
	n := CreateDateNormalizer("native", TimePrecisionDay, true)
	got := n(time.Date(2024, 3, 15, 10, 30, 45, 123456789, time.UTC))
	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
		t.Fatalf("day precision should zero time: %v", got)
	}
	if got.Day() != 15 {
		t.Fatalf("day preserved: %v", got)
	}
}
