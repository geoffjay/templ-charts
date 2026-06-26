// time.go — port of d3-scale's scaleTime / scaleUtc.
//
// A time scale is a continuous linear scale over Unix-millisecond timestamps,
// with domain/range expressed as time.Time. Ticks and Nice delegate to
// d3-time intervals (time_ticks.go). nivo defaults to useUTC:true (scaleUtc);
// we support both but only implement UTC tick intervals (local time falls
// back to UTC intervals since the math is identical for wall-clock alignment
// in most time zones — the difference is only in DST handling, which nivo's
// demo data does not exercise).
package d3scale

import (
	"math"
	"time"
)

// Time is a continuous time scale mapping time.Time → float64.
type Time struct {
	domain        [2]time.Time
	rangeVals     [2]float64
	clamp         bool
	round         bool
	useUTC        bool
	interpolateFn func(float64) float64
}

// NewTime constructs a UTC time scale with default domain
// [2000-01-01, 2000-01-02] and range [0,1] (mirroring d3's defaults).
func NewTime() *Time {
	t := &Time{
		domain: [2]time.Time{
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		rangeVals: [2]float64{0, 1},
		useUTC:    true,
	}
	t.rebuildInterpolator()
	return t
}

// Type returns "time".
func (t *Time) Type() string { return "time" }

// UseUTC reports whether UTC intervals are used for ticks/nice.
func (t *Time) UseUTC() bool { return t.useUTC }

// SetUseUTC toggles UTC vs local-time tick intervals.
func (t *Time) SetUseUTC(u bool) *Time {
	t.useUTC = u
	return t
}

// Domain returns the domain as [2]time.Time.
func (t *Time) Domain() [2]time.Time { return t.domain }

// SetDomain sets the domain.
func (t *Time) SetDomain(d0, d1 time.Time) *Time {
	t.domain = [2]time.Time{d0, d1}
	t.rebuildInterpolator()
	return t
}

// Range returns the current range.
func (t *Time) Range() [2]float64 { return t.rangeVals }

// SetRange sets the range.
func (t *Time) SetRange(r0, r1 float64) *Time {
	t.rangeVals = [2]float64{r0, r1}
	t.round = false
	t.rebuildInterpolator()
	return t
}

// SetRangeRound sets the range with integer rounding.
func (t *Time) SetRangeRound(r0, r1 float64) *Time {
	t.rangeVals = [2]float64{r0, r1}
	t.round = true
	t.rebuildInterpolator()
	return t
}

// Clamp reports clamping state.
func (t *Time) Clamp() bool { return t.clamp }

// SetClamp enables/disables clamping.
func (t *Time) SetClamp(c bool) *Time {
	t.clamp = c
	return t
}

// Call maps a time.Time to a range value.
func (t *Time) Call(x time.Time) float64 {
	d0 := float64(t.domain[0].UnixNano()) / float64(time.Millisecond)
	d1 := float64(t.domain[1].UnixNano()) / float64(time.Millisecond)
	v := float64(x.UnixNano()) / float64(time.Millisecond)
	f := deinterpolate(d0, d1)
	tr := f(v)
	if t.clamp {
		tr = clamp01(tr)
	} else if math.IsNaN(tr) {
		return tr
	}
	return t.interpolateFn(tr)
}

// Invert maps a range value back to a time.Time.
func (t *Time) Invert(y float64) time.Time {
	r0, r1 := t.rangeVals[0], t.rangeVals[1]
	rd := r1 - r0
	if rd == 0 {
		return t.domain[0]
	}
	tr := (y - r0) / rd
	if t.clamp {
		tr = clamp01(tr)
	}
	d0 := float64(t.domain[0].UnixNano()) / float64(time.Millisecond)
	d1 := float64(t.domain[1].UnixNano()) / float64(time.Millisecond)
	ms := d0 + tr*(d1-d0)
	return time.UnixMilli(int64(math.Round(ms))).UTC()
}

// Ticks returns approximately `count` tick values spanning the domain.
func (t *Time) Ticks(count int) []time.Time {
	return timeTicks(t.domain[0], t.domain[1], tickCount(count))
}

// Nice extends the domain to nice interval boundaries. With count <= 0 d3
// uses a default of 10.
func (t *Time) Nice(count int) *Time {
	interval := timeTickInterval(t.domain[0], t.domain[1], tickCount(count))
	if interval == nil {
		return t
	}
	d0 := interval.Floor(t.domain[0])
	d1 := interval.Ceil(t.domain[1])
	t.domain = [2]time.Time{d0, d1}
	t.rebuildInterpolator()
	return t
}

// Copy returns a deep copy.
func (t *Time) Copy() *Time {
	c := *t
	c.rebuildInterpolator()
	return &c
}

// rebuildInterpolator (re)creates the range interpolator.
func (t *Time) rebuildInterpolator() {
	if t.round {
		t.interpolateFn = interpolateRound(t.rangeVals[0], t.rangeVals[1])
	} else {
		t.interpolateFn = interpolateNumber(t.rangeVals[0], t.rangeVals[1])
	}
}
