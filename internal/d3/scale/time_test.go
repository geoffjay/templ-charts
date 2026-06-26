package d3scale

import (
	"math"
	"testing"
	"time"
)

func timesEqual(a, b []time.Time) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func mustParse(t *testing.T, s string) time.Time {
	t.Helper()
	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return tt
}

func parseTimes(t *testing.T, ss ...string) []time.Time {
	t.Helper()
	out := make([]time.Time, len(ss))
	for i, s := range ss {
		out[i] = mustParse(t, s)
	}
	return out
}

// d3: scaleUtc().domain([2020-01-01, 2020-12-31]).range([0,960])
//
//	Call(2020-06-01) === 399.7808219178082
func TestTimeCall(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-12-31T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 960)
	got := s.Call(mustParse(t, "2020-06-01T00:00:00Z"))
	if !approxEps(got, 399.7808219178082, 1e-6) {
		t.Errorf("Call(2020-06-01) = %v want 399.7808", got)
	}
}

// d3: scaleUtc()...invert(480) === 2020-07-01T12:00:00Z
func TestTimeInvert(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-12-31T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 960)
	got := s.Invert(480)
	if !got.Equal(mustParse(t, "2020-07-01T12:00:00Z")) {
		t.Errorf("Invert(480) = %v want 2020-07-01T12:00:00Z", got.UTC().Format(time.RFC3339))
	}
}

// d3: ticks(5) over a year span → quarterly: Jan, Apr, Jul, Oct
func TestTimeTicksQuarterly(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-12-31T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 960)
	got := s.Ticks(5)
	want := parseTimes(t,
		"2020-01-01T00:00:00Z",
		"2020-04-01T00:00:00Z",
		"2020-07-01T00:00:00Z",
		"2020-10-01T00:00:00Z",
	)
	if !timesEqual(got, want) {
		t.Errorf("ticks(5) = %v want %v", formatTimes(got), formatTimes(want))
	}
}

// d3: ticks(10) over a year span → monthly: Jan..Dec
func TestTimeTicksMonthly(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-12-31T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 960)
	got := s.Ticks(10)
	want := parseTimes(t,
		"2020-01-01T00:00:00Z", "2020-02-01T00:00:00Z",
		"2020-03-01T00:00:00Z", "2020-04-01T00:00:00Z",
		"2020-05-01T00:00:00Z", "2020-06-01T00:00:00Z",
		"2020-07-01T00:00:00Z", "2020-08-01T00:00:00Z",
		"2020-09-01T00:00:00Z", "2020-10-01T00:00:00Z",
		"2020-11-01T00:00:00Z", "2020-12-01T00:00:00Z",
	)
	if !timesEqual(got, want) {
		t.Errorf("ticks(10) mismatch:\n got=%v\nwant=%v", formatTimes(got), formatTimes(want))
	}
}

// d3: hour-scale ticks(5) → every 3h; ticks(10) → every 1h
func TestTimeTicksHourly(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-01-01T12:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 100)
	got5 := s.Ticks(5)
	want5 := parseTimes(t,
		"2020-01-01T00:00:00Z", "2020-01-01T03:00:00Z",
		"2020-01-01T06:00:00Z", "2020-01-01T09:00:00Z",
		"2020-01-01T12:00:00Z",
	)
	if !timesEqual(got5, want5) {
		t.Errorf("hour ticks(5) = %v want %v", formatTimes(got5), formatTimes(want5))
	}
}

// d3: second-scale ticks(5) → every 1s (span is 10s; tickInterval picks 1s)
func TestTimeTicksSeconds(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-01-01T00:00:10Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 100)
	got := s.Ticks(5)
	want := parseTimes(t,
		"2020-01-01T00:00:00Z", "2020-01-01T00:00:01Z",
		"2020-01-01T00:00:02Z", "2020-01-01T00:00:03Z",
		"2020-01-01T00:00:04Z", "2020-01-01T00:00:05Z",
		"2020-01-01T00:00:06Z", "2020-01-01T00:00:07Z",
		"2020-01-01T00:00:08Z", "2020-01-01T00:00:09Z",
		"2020-01-01T00:00:10Z",
	)
	if !timesEqual(got, want) {
		t.Errorf("second ticks(5) mismatch:\n got=%v\nwant=%v", formatTimes(got), formatTimes(want))
	}
}

// d3: year-span ticks(5) over 5 years → yearly
func TestTimeTicksYearly(t *testing.T) {
	d0 := mustParse(t, "2018-01-01T00:00:00Z")
	d1 := mustParse(t, "2023-01-01T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 500)
	got := s.Ticks(5)
	want := parseTimes(t,
		"2018-01-01T00:00:00Z", "2019-01-01T00:00:00Z",
		"2020-01-01T00:00:00Z", "2021-01-01T00:00:00Z",
		"2022-01-01T00:00:00Z", "2023-01-01T00:00:00Z",
	)
	if !timesEqual(got, want) {
		t.Errorf("year ticks(5) = %v want %v", formatTimes(got), formatTimes(want))
	}
}

// d3: nice(5) on [2020-01-01, 2020-12-31] → [2020-01-01, 2021-01-01]
func TestTimeNice(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-12-31T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).Nice(5)
	d := s.Domain()
	if !d[0].Equal(mustParse(t, "2020-01-01T00:00:00Z")) {
		t.Errorf("nice d0 = %v want 2020-01-01", d[0].UTC().Format(time.RFC3339))
	}
	if !d[1].Equal(mustParse(t, "2021-01-01T00:00:00Z")) {
		t.Errorf("nice d1 = %v want 2021-01-01", d[1].UTC().Format(time.RFC3339))
	}
}

// d3: clamp(true) maps out-of-domain to range endpoints
func TestTimeClamp(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-01-02T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 100).SetClamp(true)
	if got := s.Call(mustParse(t, "2020-01-05T00:00:00Z")); got != 100 {
		t.Errorf("clamp future = %v want 100", got)
	}
	if got := s.Call(mustParse(t, "2019-01-01T00:00:00Z")); got != 0 {
		t.Errorf("clamp past = %v want 0", got)
	}
}

func TestTimeType(t *testing.T) {
	if NewTime().Type() != "time" {
		t.Error("Type() != time")
	}
}

func TestTimeUseUTC(t *testing.T) {
	if !NewTime().UseUTC() {
		t.Error("default UseUTC should be true")
	}
}

func TestTimeCopy(t *testing.T) {
	d0 := mustParse(t, "2020-01-01T00:00:00Z")
	d1 := mustParse(t, "2020-01-02T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 100)
	c := s.Copy()
	c.SetRange(0, 999)
	if s.Range() != [2]float64{0, 100} {
		t.Errorf("original Range after copy mutate = %v want [0 100]", s.Range())
	}
}

func formatTimes(ts []time.Time) []string {
	out := make([]string, len(ts))
	for i, t := range ts {
		out[i] = t.UTC().Format(time.RFC3339)
	}
	return out
}

// --- Interval unit tests (sanity checks for the port) ---------------------

func TestIntervalUTCDayFloorCeil(t *testing.T) {
	d := utcDayInterval()
	mid := mustParse(t, "2020-03-15T13:45:30Z")
	want := mustParse(t, "2020-03-15T00:00:00Z")
	if !d.Floor(mid).Equal(want) {
		t.Errorf("day.floor = %v want %v", d.Floor(mid), want)
	}
}

func TestIntervalUTCMonthEvery3(t *testing.T) {
	q := utcMonthInterval().Every(3)
	if q == nil {
		t.Fatal("month.every(3) = nil")
	}
	// floor of Apr 15 2020 → Apr 1 (month index 3, 3%3==0)
	got := q.Floor(mustParse(t, "2020-04-15T00:00:00Z"))
	want := mustParse(t, "2020-04-01T00:00:00Z")
	if !got.Equal(want) {
		t.Errorf("month.every(3).floor(Apr 15) = %v want %v", got, want)
	}
	// floor of May 15 2020 → Apr 1 (month index 4, 4%3!=0 → step back to Apr)
	got = q.Floor(mustParse(t, "2020-05-15T00:00:00Z"))
	if !got.Equal(want) {
		t.Errorf("month.every(3).floor(May 15) = %v want %v", got, want)
	}
}

func TestIntervalMillisecondEvery(t *testing.T) {
	m := millisecondEvery(50)
	if m == nil {
		t.Fatal("ms.every(50) = nil")
	}
	got := m.Floor(mustParse(t, "2020-01-01T00:00:00.073Z"))
	// 73ms floored to 50ms = 50ms
	want := mustParse(t, "2020-01-01T00:00:00.050Z")
	if !got.Equal(want) {
		t.Errorf("ms.every(50).floor = %v want %v", got.UTC().Format(time.RFC3339Nano), want.UTC().Format(time.RFC3339Nano))
	}
}

func TestTimeReversedDomain(t *testing.T) {
	d0 := mustParse(t, "2020-01-02T00:00:00Z")
	d1 := mustParse(t, "2020-01-01T00:00:00Z")
	s := NewTime().SetDomain(d0, d1).SetRange(0, 100)
	got := s.Call(mustParse(t, "2020-01-01T12:00:00Z"))
	// midpoint → 50
	if !approxEps(got, 50, 1e-6) {
		t.Errorf("reversed Call(midpoint) = %v want 50", got)
	}
	if math.IsNaN(got) {
		t.Errorf("got NaN")
	}
}
