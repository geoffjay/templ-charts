package scales

import (
	"math"
	"testing"
	"time"
)

func linearData(vals ...float64) ComputedSerieAxis {
	all := make([]any, len(vals))
	for i, v := range vals {
		all[i] = v
	}
	return ComputedSerieAxis{All: all, Min: vals[0], Max: vals[len(vals)-1]}
}

// --- linear ------------------------------------------------------------------

func TestCreateLinearScale_XAxisAndTicks(t *testing.T) {
	s := ComputeScale(ScaleLinearSpec{Min: FloatVal(0), Max: FloatVal(20)}, linearData(0, 20), 200, ScaleAxisX)
	if got := s.Call(0); got != 0 {
		t.Errorf("linear x(0) = %v, want 0", got)
	}
	if got := s.Call(20); got != 200 {
		t.Errorf("linear x(20) = %v, want 200", got)
	}
	if got := s.Call(10); got != 100 {
		t.Errorf("linear x(10) = %v, want 100", got)
	}
	si := s.(*scaleImpl)
	ticks := si.Ticks(5)
	if len(ticks) == 0 {
		t.Fatal("no ticks")
	}
	if ticks[0].(float64) != 0 || ticks[len(ticks)-1].(float64) != 20 {
		t.Errorf("ticks = %v, want to span 0..20", ticks)
	}
	dom := si.Domain()
	if dom[0].(float64) != 0 || dom[1].(float64) != 20 {
		t.Errorf("domain = %v, want [0 20]", dom)
	}
	// A continuous scale has no bandwidth/step/round.
	if si.Bandwidth() != 0 || si.Step() != 0 || si.Round() {
		t.Errorf("linear bandwidth/step/round = %v/%v/%v, want zeros", si.Bandwidth(), si.Step(), si.Round())
	}
}

func TestCreateLinearScale_ReverseClampRoundNiceInt(t *testing.T) {
	spec := ScaleLinearSpec{Min: FloatVal(0), Max: FloatVal(10), Reverse: true, Clamp: true, Round: true, Nice: 5}
	s := ComputeScale(spec, linearData(0, 10), 100, ScaleAxisX)
	// Reversed: domain is [10, 0].
	if got := s.Call(10); got != 0 {
		t.Errorf("reversed(10) = %v, want 0", got)
	}
	if got := s.Call(0); got != 100 {
		t.Errorf("reversed(0) = %v, want 100", got)
	}
	// Clamped: values beyond the domain pin to the range ends.
	if got := s.Call(-5); got != 100 {
		t.Errorf("clamp(-5) = %v, want 100", got)
	}
	if got := s.Call(15); got != 0 {
		t.Errorf("clamp(15) = %v, want 0", got)
	}
}

func TestCreateLinearScale_StackedUsesStackedMinMax(t *testing.T) {
	mn, mx := -3.0, 33.0
	data := ComputedSerieAxis{All: []any{0.0, 10.0}, Min: 0.0, Max: 10.0, MinStacked: &mn, MaxStacked: &mx}
	spec := ScaleLinearSpec{Min: AutoFloat(), Max: AutoFloat(), Stacked: true, Nice: false}
	s := ComputeScale(spec, data, 100, ScaleAxisX)
	dom := s.(*scaleImpl).Domain()
	if dom[0].(float64) != -3 || dom[1].(float64) != 33 {
		t.Errorf("stacked domain = %v, want [-3 33]", dom)
	}
}

// --- log ----------------------------------------------------------------------

func TestCreateLogScale(t *testing.T) {
	spec := ScaleLogSpec{Min: FloatVal(1), Max: FloatVal(100), Nice: false}
	s := ComputeScale(spec, linearData(1, 100), 100, ScaleAxisX)
	if s.Type() != ScaleTypeLog {
		t.Fatalf("type = %v", s.Type())
	}
	if got := s.Call(1); !approxF(got, 0) {
		t.Errorf("log(1) = %v, want 0", got)
	}
	if got := s.Call(10); !approxF(got, 50) {
		t.Errorf("log(10) = %v, want 50", got)
	}
	if got := s.Call(100); !approxF(got, 100) {
		t.Errorf("log(100) = %v, want 100", got)
	}
	si := s.(*scaleImpl)
	ticks := si.Ticks(10)
	if len(ticks) == 0 {
		t.Error("log scale should produce ticks")
	}
	dom := si.Domain()
	if dom[0].(float64) != 1 || dom[1].(float64) != 100 {
		t.Errorf("log domain = %v", dom)
	}
}

func TestCreateLogScale_ReverseRoundNiceDefaultBase(t *testing.T) {
	// Base 0 falls back to 10; y axis inverts the range.
	spec := ScaleLogSpec{Min: FloatVal(1), Max: FloatVal(100), Reverse: true, Round: true, Nice: true}
	s := ComputeScale(spec, linearData(1, 100), 100, ScaleAxisY)
	// Reversed domain [100, 1] on y range [100, 0]: log(100) → 100.
	if got := s.Call(100); !approxF(got, 100) {
		t.Errorf("reversed log(100) = %v, want 100", got)
	}
	if got := s.Call(1); !approxF(got, 0) {
		t.Errorf("reversed log(1) = %v, want 0", got)
	}
}

// --- symlog -------------------------------------------------------------------

func TestCreateSymlogScale(t *testing.T) {
	spec := ScaleSymlogSpec{Min: FloatVal(0), Max: FloatVal(100), Nice: false}
	s := ComputeScale(spec, linearData(0, 100), 100, ScaleAxisX)
	if s.Type() != ScaleTypeSymlog {
		t.Fatalf("type = %v", s.Type())
	}
	if got := s.Call(0); !approxF(got, 0) {
		t.Errorf("symlog(0) = %v, want 0", got)
	}
	if got := s.Call(100); !approxF(got, 100) {
		t.Errorf("symlog(100) = %v, want 100", got)
	}
	// Monotonic in between.
	if !(s.Call(10) < s.Call(50)) {
		t.Errorf("symlog not monotonic: f(10)=%v f(50)=%v", s.Call(10), s.Call(50))
	}
	si := s.(*scaleImpl)
	if got := si.Ticks(10); len(got) == 0 {
		t.Error("symlog should produce ticks")
	}
	dom := si.Domain()
	if dom[0].(float64) != 0 || dom[1].(float64) != 100 {
		t.Errorf("symlog domain = %v", dom)
	}
}

func TestCreateSymlogScale_ReverseRoundConstantNiceInt(t *testing.T) {
	spec := ScaleSymlogSpec{Constant: 2, Min: FloatVal(0), Max: FloatVal(10), Reverse: true, Round: true, Nice: 5}
	s := ComputeScale(spec, linearData(0, 10), 100, ScaleAxisX)
	// Reversed: min value maps to the top of the range.
	if !(s.Call(0) > s.Call(10)) {
		t.Errorf("reversed symlog: f(0)=%v should exceed f(10)=%v", s.Call(0), s.Call(10))
	}
}

// --- time ---------------------------------------------------------------------

func TestCreateTimeScale_AutoDomainFromData(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	data := ComputedSerieAxis{All: []any{t0, t1}, Min: t0, Max: t1}
	spec := ScaleTimeSpec{Format: "native", Min: "auto", Max: "auto", UseUTC: true}
	s := ComputeScale(spec, data, 100, ScaleAxisX)
	if s.Type() != ScaleTypeTime {
		t.Fatalf("type = %v", s.Type())
	}
	if got := s.Call(t0); !approxF(got, 0) {
		t.Errorf("time(min) = %v, want 0", got)
	}
	if got := s.Call(t1); !approxF(got, 100) {
		t.Errorf("time(max) = %v, want 100", got)
	}
	mid := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	if got := s.Call(mid); !approxF(got, 50) {
		t.Errorf("time(mid) = %v, want 50", got)
	}
	// A string RFC3339 value is coerced.
	if got := s.Call("2024-01-01T00:00:00Z"); !approxF(got, 0) {
		t.Errorf("time(string) = %v, want 0", got)
	}
	// Bad values map to 0.
	if got := s.Call("not-a-time"); got != 0 {
		t.Errorf("time(bad string) = %v, want 0", got)
	}
	var nilT *time.Time
	if got := s.Call(nilT); got != 0 {
		t.Errorf("time(nil ptr) = %v, want 0", got)
	}
	// Pointer values are dereferenced.
	if got := s.Call(&t1); !approxF(got, 100) {
		t.Errorf("time(*time.Time) = %v, want 100", got)
	}
	si := s.(*scaleImpl)
	if got := si.Ticks(5); len(got) == 0 {
		t.Error("time ticks empty")
	}
	dom := si.Domain()
	if len(dom) != 2 {
		t.Fatalf("time domain = %v", dom)
	}
	if !dom[0].(time.Time).Equal(t0) || !dom[1].(time.Time).Equal(t1) {
		t.Errorf("time domain = %v, want [%v %v]", dom, t0, t1)
	}
}

func TestCreateTimeScale_ExplicitMinMaxAndNice(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
	data := ComputedSerieAxis{All: []any{t0, t1}, Min: t0, Max: t1}
	// Explicit time.Time bounds.
	spec := ScaleTimeSpec{Format: "native", Min: t0, Max: t1, UseUTC: true, Nice: true}
	s := ComputeScale(spec, data, 100, ScaleAxisX)
	if got := s.Call(t0); got < 0 {
		t.Errorf("explicit time scale f(min) = %v", got)
	}
	// Explicit RFC3339 string bounds.
	spec2 := ScaleTimeSpec{Format: "native", Min: "2024-01-01T00:00:00Z", Max: "2024-01-11T00:00:00Z", UseUTC: false}
	s2 := ComputeScale(spec2, data, 100, ScaleAxisX)
	if got := s2.Call(t1); !approxF(got, 100) {
		t.Errorf("string-bounded time scale f(max) = %v, want 100", got)
	}
}

// --- band / point --------------------------------------------------------------

func TestBandScale_StepRoundDomainTicks(t *testing.T) {
	data := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	s := ComputeScale(ScaleBandSpec{Round: true}, data, 300, ScaleAxisX)
	si := s.(*scaleImpl)
	if !si.Round() {
		t.Error("band Round() = false, want true")
	}
	if bw := si.Bandwidth(); bw != 100 {
		t.Errorf("bandwidth = %v, want 100", bw)
	}
	if st := si.Step(); st != 100 {
		t.Errorf("step = %v, want 100", st)
	}
	ticks := si.Ticks(10)
	if len(ticks) != 3 || ticks[0].(string) != "a" || ticks[2].(string) != "c" {
		t.Errorf("band ticks = %v, want domain values", ticks)
	}
	dom := si.Domain()
	if len(dom) != 3 || dom[1].(string) != "b" {
		t.Errorf("band domain = %v", dom)
	}
	// Non-string values are stringified by Call.
	dataNum := ComputedSerieAxis{All: []any{1.0, 2.0}}
	sn := ComputeScale(ScaleBandSpec{}, dataNum, 100, ScaleAxisX)
	if got := sn.Call(1.0); got != 0 {
		t.Errorf("band Call(1.0) = %v, want 0", got)
	}
	// Unknown keys (nil stringifies to "") map to NaN, matching d3 scaleBand.
	if got := sn.Call(nil); !math.IsNaN(got) {
		t.Errorf("band Call(nil) = %v, want NaN (unknown key)", got)
	}
}

func TestPointScale_StepDomainTicks(t *testing.T) {
	data := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	s := ComputeScale(ScalePointSpec{}, data, 200, ScaleAxisX)
	si := s.(*scaleImpl)
	if got := si.Step(); got != 100 {
		t.Errorf("point step = %v, want 100", got)
	}
	if si.Round() {
		t.Error("point Round() = true, want false")
	}
	ticks := si.Ticks(10)
	if len(ticks) != 3 {
		t.Errorf("point ticks = %v", ticks)
	}
	dom := si.Domain()
	if len(dom) != 3 || dom[0].(string) != "a" {
		t.Errorf("point domain = %v", dom)
	}
}

// --- helpers -------------------------------------------------------------------

func TestToFloatCoercions(t *testing.T) {
	epoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		in   any
		want float64
	}{
		{1.5, 1.5},
		{2, 2},
		{int64(3), 3},
		{"4.5", 4.5},
		{"junk", 0},
		{nil, 0},
		{struct{}{}, 0},
		{epoch, float64(epoch.UnixMilli())},
	}
	for _, c := range cases {
		if got := toFloat(c.in); got != c.want {
			t.Errorf("toFloat(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToStringCoercions(t *testing.T) {
	if got := toString(nil); got != "" {
		t.Errorf("toString(nil) = %q", got)
	}
	if got := toString("x"); got != "x" {
		t.Errorf("toString(string) = %q", got)
	}
	if got := toString(1.5); got != "1.5" {
		t.Errorf("toString(1.5) = %q", got)
	}
}

func TestToTimeCoercions(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if got, ok := toTime(now); !ok || !got.Equal(now) {
		t.Errorf("toTime(time) = %v %v", got, ok)
	}
	if got, ok := toTime(&now); !ok || !got.Equal(now) {
		t.Errorf("toTime(*time) = %v %v", got, ok)
	}
	var nilT *time.Time
	if _, ok := toTime(nilT); ok {
		t.Error("toTime(nil *time) should fail")
	}
	if got, ok := toTime("2024-06-01T12:00:00Z"); !ok || !got.Equal(now) {
		t.Errorf("toTime(rfc3339) = %v %v", got, ok)
	}
	if _, ok := toTime("junk"); ok {
		t.Error("toTime(junk) should fail")
	}
	if _, ok := toTime(42); ok {
		t.Error("toTime(int) should fail")
	}
}

func TestNiceBool(t *testing.T) {
	if ok, n := niceBool(true); !ok || n != 0 {
		t.Errorf("niceBool(true) = %v %v", ok, n)
	}
	if ok, _ := niceBool(false); ok {
		t.Error("niceBool(false) should be off")
	}
	if ok, n := niceBool(7); !ok || n != 7 {
		t.Errorf("niceBool(7) = %v %v", ok, n)
	}
	if ok, n := niceBool(7.0); !ok || n != 7 {
		t.Errorf("niceBool(7.0) = %v %v", ok, n)
	}
	if ok, _ := niceBool(nil); ok {
		t.Error("niceBool(nil) should be off")
	}
}

func TestGetOtherAxis(t *testing.T) {
	if got := GetOtherAxis(ScaleAxisX); got != ScaleAxisY {
		t.Errorf("other of x = %v", got)
	}
	if got := GetOtherAxis(ScaleAxisY); got != ScaleAxisX {
		t.Errorf("other of y = %v", got)
	}
}

// fakeSpec / fakeScale exercise the unknown-type fallbacks.
type fakeSpec struct{}

func (fakeSpec) ScaleType() ScaleType { return ScaleType("fake") }

type fakeScale struct{}

func (fakeScale) Type() ScaleType    { return ScaleType("fake") }
func (fakeScale) Call(v any) float64 { return 7 }

func TestComputeScale_UnknownSpec(t *testing.T) {
	if got := ComputeScale(fakeSpec{}, ComputedSerieAxis{}, 100, ScaleAxisX); got != nil {
		t.Errorf("unknown spec = %v, want nil", got)
	}
}

func TestGetScaleTicks(t *testing.T) {
	data := linearData(0, 20)
	s := ComputeScale(ScaleLinearSpec{Min: FloatVal(0), Max: FloatVal(20)}, data, 100, ScaleAxisX)
	// Explicit values win.
	vals := []any{1.0, 2.0}
	got := GetScaleTicks(s, TicksSpec{HasValues: true, Values: vals})
	if len(got) != 2 || got[0].(float64) != 1 {
		t.Errorf("HasValues ticks = %v", got)
	}
	// Count is honored.
	if got := GetScaleTicks(s, TicksSpec{HasCount: true, Count: 2}); len(got) == 0 {
		t.Error("HasCount ticks empty")
	}
	// Interval falls back to default ticks.
	if got := GetScaleTicks(s, TicksSpec{HasInterval: true, Interval: "every day"}); len(got) == 0 {
		t.Error("HasInterval ticks empty")
	}
	// Default: scale ticks.
	if got := GetScaleTicks(s, TicksSpec{}); len(got) == 0 {
		t.Error("default ticks empty")
	}
	// Non-scaleImpl scale yields nil.
	if got := GetScaleTicks(fakeScale{}, TicksSpec{}); got != nil {
		t.Errorf("fake scale ticks = %v, want nil", got)
	}
}

func TestCenterScale(t *testing.T) {
	// Band: offset by bandwidth/2.
	data := ComputedSerieAxis{All: []any{"a", "b"}}
	band := ComputeScale(ScaleBandSpec{}, data, 100, ScaleAxisX)
	centered := CenterScale(band)
	if got, want := centered("a"), band.Call("a")+band.(*scaleImpl).Bandwidth()/2; got != want {
		t.Errorf("centered band(a) = %v, want %v", got, want)
	}
	// Rounded band: offset is rounded.
	data3 := ComputedSerieAxis{All: []any{"a", "b", "c"}}
	bandR := ComputeScale(ScaleBandSpec{Round: true}, data3, 100, ScaleAxisX)
	cR := CenterScale(bandR)
	if got := cR("a"); got != float64(int64(got)) {
		t.Errorf("rounded centered band(a) = %v, want integral offset applied", got)
	}
	// Point scale with a single value has bandwidth 0 → passthrough.
	point := ComputeScale(ScalePointSpec{}, data, 100, ScaleAxisX)
	cp := CenterScale(point)
	if got := cp("a"); got != point.Call("a") {
		t.Errorf("centered point(a) = %v, want passthrough %v", got, point.Call("a"))
	}
	// Continuous scale passes through.
	lin := ComputeScale(ScaleLinearSpec{Min: FloatVal(0), Max: FloatVal(10)}, linearData(0, 10), 100, ScaleAxisX)
	cl := CenterScale(lin)
	if got := cl(5.0); got != lin.Call(5.0) {
		t.Errorf("centered linear(5) = %v, want passthrough", got)
	}
	// Non-scaleImpl passes through.
	cf := CenterScale(fakeScale{})
	if got := cf(nil); got != 7 {
		t.Errorf("centered fake = %v, want 7", got)
	}
}

func approxF(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < 1e-6
}
