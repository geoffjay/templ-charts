// time_ticks.go — port of d3-time's ticker() which produces time ticks and
// the "tickInterval" selector used by d3-scale's time/utc scales.
//
// The algorithm: given a target span/count, bisect a table of candidate
// (interval, step, duration) tuples to pick the interval whose duration is
// closest to target = |stop - start| / count. For spans larger than a year,
// fall back to year.every(tickStep(start/year, stop/year, count)).
package d3scale

import (
	"sort"
	"time"

	d3array "github.com/geoffjay/templ-charts/internal/d3/array"
)

// tickIntervalEntry is one row of d3-time's tickIntervals table.
type tickIntervalEntry struct {
	interval *Interval
	step     int
	duration time.Duration
}

// utcTickIntervals is d3-time's tickIntervals table for UTC.
var utcTickIntervals = []tickIntervalEntry{
	{utcSecondInterval(), 1, durationSecond},
	{utcSecondInterval(), 5, 5 * durationSecond},
	{utcSecondInterval(), 15, 15 * durationSecond},
	{utcSecondInterval(), 30, 30 * durationSecond},
	{utcMinuteInterval(), 1, durationMinute},
	{utcMinuteInterval(), 5, 5 * durationMinute},
	{utcMinuteInterval(), 15, 15 * durationMinute},
	{utcMinuteInterval(), 30, 30 * durationMinute},
	{utcHourInterval(), 1, durationHour},
	{utcHourInterval(), 3, 3 * durationHour},
	{utcHourInterval(), 6, 6 * durationHour},
	{utcHourInterval(), 12, 12 * durationHour},
	{utcDayInterval(), 1, durationDay},
	{utcDayInterval(), 2, 2 * durationDay},
	{utcWeekInterval(), 1, durationWeek},
	{utcMonthInterval(), 1, durationMonth},
	{utcMonthInterval(), 3, 3 * durationMonth},
	{utcYearInterval(), 1, durationYear},
}

// timeTickInterval returns the d3-time interval best suited to produce
// approximately `count` ticks between start and stop. Returns nil if no
// suitable interval (e.g. span is 0).
//
// This is a faithful port of d3-time's tickInterval().
func timeTickInterval(start, stop time.Time, count int) *Interval {
	target := float64(absDur(stop.Sub(start))) / float64(count)
	// bisect by duration
	i := sort.Search(len(utcTickIntervals), func(j int) bool {
		return float64(utcTickIntervals[j].duration) > target
	})
	if i == len(utcTickIntervals) {
		// fall back to year.every(tickStep(start/year, stop/year, count))
		startYear := float64(start.UTC().Year())
		stopYear := float64(stop.UTC().Year())
		step := d3array.TickIncrement(startYear, stopYear, count)
		if step <= 0 {
			return nil
		}
		return utcYearEvery(int(step))
	}
	if i == 0 {
		// millisecond.every(max(tickStep(start, stop, count), 1))
		stepMs := d3array.TickIncrement(
			float64(start.UnixNano()/int64(time.Millisecond)),
			float64(stop.UnixNano()/int64(time.Millisecond)),
			count,
		)
		if stepMs < 1 {
			stepMs = 1
		}
		return millisecondEvery(int(stepMs))
	}
	// pick between i-1 and i based on which is closer to target
	var chosen int
	if target/float64(utcTickIntervals[i-1].duration) < float64(utcTickIntervals[i].duration)/target {
		chosen = i - 1
	} else {
		chosen = i
	}
	e := utcTickIntervals[chosen]
	return e.interval.Every(e.step)
}

// timeTicks returns approximately `count` ticks between start and stop as
// time.Time values. If `interval` is non-nil (has a Range method) it is used
// directly; otherwise timeTickInterval picks one.
//
// Faithful port of d3-time's ticks().
func timeTicks(start, stop time.Time, count int) []time.Time {
	reverse := stop.Before(start)
	if reverse {
		start, stop = stop, start
	}
	var interval *Interval
	interval = timeTickInterval(start, stop, count)
	if interval == nil {
		return nil
	}
	// d3 uses interval.range(start, +stop + 1) — inclusive of stop.
	ticks := interval.Range(start, stop.Add(time.Millisecond), 1)
	if reverse {
		for l, r := 0, len(ticks)-1; l < r; l, r = l+1, r-1 {
			ticks[l], ticks[r] = ticks[r], ticks[l]
		}
	}
	return ticks
}

func absDur(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
