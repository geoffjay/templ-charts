package scales

import (
	"math"
	"testing"
	"time"
)

func TestComputeXYScalesForSeries_StackedX(t *testing.T) {
	series := []Serie{
		{Data: []SerieDatum{{X: 1.0, Y: "a"}, {X: 2.0, Y: "b"}}},
		{Data: []SerieDatum{{X: 3.0, Y: "a"}, {X: 4.0, Y: "b"}}},
	}
	xSpec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Stacked: true, Nice: false}
	res := ComputeXYScalesForSeries(series, xSpec, ScaleBandSpec{}, 100, 100)
	// Stacked x: for "a" the stacked values are 1 and 1+3=4; for "b" 2 and 6.
	if res.X.MaxStacked == nil || *res.X.MaxStacked != 6 {
		t.Fatalf("MaxStacked = %v, want 6", res.X.MaxStacked)
	}
	if res.X.MinStacked == nil || *res.X.MinStacked != 1 {
		t.Fatalf("MinStacked = %v, want 1", res.X.MinStacked)
	}
	d := res.Series[1].Data[0]
	if d.XStacked == nil || *d.XStacked != 4 {
		t.Errorf("series[1] a XStacked = %v, want 4", d.XStacked)
	}
	// The stacked scale positions by the stacked value.
	if d.Position.X == nil {
		t.Fatal("stacked datum should have an x position")
	}
	if got, want := *d.Position.X, res.XScale.Call(4.0); got != want {
		t.Errorf("stacked position = %v, want scale(4) = %v", got, want)
	}
}

func TestComputeXYScalesForSeries_StackedY_MissingDatum(t *testing.T) {
	// Second series has no datum at x="two": the stack skips it.
	series := []Serie{
		{Data: []SerieDatum{{X: "one", Y: 10.0}, {X: "two", Y: 20.0}}},
		{Data: []SerieDatum{{X: "one", Y: 5.0}}},
	}
	ySpec := ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Stacked: true, Nice: false}
	res := ComputeXYScalesForSeries(series, ScalePointSpec{}, ySpec, 100, 100)
	if res.Y.MaxStacked == nil || *res.Y.MaxStacked != 20 {
		t.Fatalf("MaxStacked = %v, want 20 (one: 10+5=15, two: 20)", res.Y.MaxStacked)
	}
	if got := res.Series[1].Data[0].YStacked; got == nil || *got != 15 {
		t.Errorf("series[1] one YStacked = %v, want 15", got)
	}
	// A datum with nil Y gets no position.
	series2 := []Serie{{Data: []SerieDatum{{X: "one", Y: nil}}}}
	res2 := ComputeXYScalesForSeries(series2, ScalePointSpec{}, ScaleLinearSpec{Min: FloatVal(0), Max: FloatVal(10)}, 100, 100)
	if res2.Series[0].Data[0].Position.Y != nil {
		t.Error("nil Y should yield nil position")
	}
}

func TestComputeXYScalesForSeries_TimeAxis(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	series := []Serie{
		{Data: []SerieDatum{{X: t0, Y: 1.0}, {X: t1, Y: 2.0}, {X: t0, Y: 3.0}}},
	}
	xSpec := ScaleTimeSpec{Format: "native", Min: "auto", Max: "auto", UseUTC: true}
	res := ComputeXYScalesForSeries(series, xSpec, ScaleLinearSpec{Min: FloatVal(0), Max: AutoFloat(), Nice: false}, 100, 100)
	// Duplicate t0 is deduped.
	if len(res.X.All) != 2 {
		t.Fatalf("time axis All = %v, want 2 unique times", res.X.All)
	}
	if !res.X.Min.(time.Time).Equal(t0) || !res.X.Max.(time.Time).Equal(t1) {
		t.Errorf("time min/max = %v/%v", res.X.Min, res.X.Max)
	}
	if res.XScale.Type() != ScaleTypeTime {
		t.Errorf("x scale type = %v, want time", res.XScale.Type())
	}
}

func TestGenerateSeriesAxis_SkipsUnparseableLinearValues(t *testing.T) {
	series := []Serie{
		{Data: []SerieDatum{
			{X: 1.0, Y: 0.0},
			{X: "2", Y: 0.0},        // numeric string is coerced
			{X: "junk", Y: 0.0},     // non-numeric string skipped
			{X: math.NaN(), Y: 0.0}, // NaN skipped
			{X: nil, Y: 0.0},        // nil skipped
			{X: true, Y: 0.0},       // unsupported type skipped
			{X: int64(3), Y: 0.0},   // int64 coerced
		}},
	}
	axis := generateSeriesAxis(series, ScaleAxisX, ScaleLinearSpec{})
	if len(axis.All) != 3 {
		t.Fatalf("All = %v, want [1 2 3]", axis.All)
	}
	if axis.Min.(float64) != 1 || axis.Max.(float64) != 3 {
		t.Errorf("min/max = %v/%v, want 1/3", axis.Min, axis.Max)
	}
}

func TestGenerateSeriesAxis_EmptyAxes(t *testing.T) {
	empty := []Serie{{Data: nil}}
	lin := generateSeriesAxis(empty, ScaleAxisX, ScaleLinearSpec{})
	if len(lin.All) != 0 || lin.Min.(float64) != 0 || lin.Max.(float64) != 0 {
		t.Errorf("empty linear axis = %+v", lin)
	}
	tm := generateSeriesAxis(empty, ScaleAxisX, ScaleTimeSpec{})
	if len(tm.All) != 0 || !tm.Min.(time.Time).IsZero() {
		t.Errorf("empty time axis = %+v", tm)
	}
	band := generateSeriesAxis(empty, ScaleAxisY, ScaleBandSpec{})
	if len(band.All) != 0 || band.Min != nil {
		t.Errorf("empty band axis = %+v", band)
	}
}

func TestFmtKey(t *testing.T) {
	epoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		in   any
		want string
	}{
		{nil, ""},
		{"a", "s:a"},
		{1.5, "f:1.5"},
		{3, "f:3"},
		{epoch, "t:" + epoch.Format(time.RFC3339Nano)},
		{int64(4), "x:4"}, // falls back to float formatting with x: prefix
	}
	for _, c := range cases {
		if got := fmtKey(c.in); got != c.want {
			t.Errorf("fmtKey(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsNil(t *testing.T) {
	if !isNil(nil) {
		t.Error("isNil(nil) = false")
	}
	if isNil(0) {
		t.Error("isNil(0) = true")
	}
}

func TestToFloatOK(t *testing.T) {
	if _, ok := toFloatOK(nil); ok {
		t.Error("nil should fail")
	}
	if _, ok := toFloatOK(math.NaN()); ok {
		t.Error("NaN should fail")
	}
	if v, ok := toFloatOK(1.5); !ok || v != 1.5 {
		t.Errorf("float64 = %v %v", v, ok)
	}
	if v, ok := toFloatOK(2); !ok || v != 2 {
		t.Errorf("int = %v %v", v, ok)
	}
	if v, ok := toFloatOK(int64(3)); !ok || v != 3 {
		t.Errorf("int64 = %v %v", v, ok)
	}
	if v, ok := toFloatOK("4.5"); !ok || v != 4.5 {
		t.Errorf("string = %v %v", v, ok)
	}
	if _, ok := toFloatOK("junk"); ok {
		t.Error("bad string should fail")
	}
	if _, ok := toFloatOK(true); ok {
		t.Error("bool should fail")
	}
}

func TestUniqueTimes(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Hour)
	got := uniqueTimes([]any{t0, t1, t0, "junk", nil})
	if len(got) != 2 {
		t.Fatalf("uniqueTimes = %v, want 2 entries", got)
	}
}
