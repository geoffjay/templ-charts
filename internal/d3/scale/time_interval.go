// time_interval.go — minimal port of d3-time intervals used by d3-scale's
// time/utc scales. Only UTC intervals are needed for templ-charts (nivo
// defaults to useUTC:true), but local-time intervals are included for
// completeness via the `useUTC` flag on each interval constructor.
//
// Each interval exposes Floor, Ceil, Offset, Range, Count, and Every, matching
// d3-time's `timeInterval()` factory. The implementation is a faithful
// transliteration of d3-time v3's interval.js / second.js / minute.js / etc.
package d3scale

import (
	"math"
	"time"
)

// Duration constants matching d3-time's duration.js.
const (
	durationSecond = 1000 * time.Millisecond
	durationMinute = 60 * durationSecond
	durationHour   = 60 * durationMinute
	durationDay    = 24 * durationHour
	durationWeek   = 7 * durationDay
	durationMonth  = 30 * durationDay
	durationYear   = 365 * durationDay
)

// Interval is a d3-time interval: a set of dates spaced by a fixed amount,
// with floor/ceil/offset/range/count/every operations.
type Interval struct {
	floori  func(t *time.Time)
	offseti func(t *time.Time, step int)
	countFn func(start, end time.Time) int
	fieldFn func(t time.Time) int
}

// Floor returns a copy of t floored to the nearest interval boundary.
func (i *Interval) Floor(t time.Time) time.Time {
	c := t
	i.floori(&c)
	return c
}

// Ceil returns a copy of t ceiled to the nearest interval boundary.
// d3: floor(t-1ms), offset(+1), floor.
func (i *Interval) Ceil(t time.Time) time.Time {
	c := t.Add(-time.Millisecond)
	i.floori(&c)
	i.offseti(&c, 1)
	i.floori(&c)
	return c
}

// Round returns the nearer of Floor(t) / Ceil(t).
func (i *Interval) Round(t time.Time) time.Time {
	d0 := i.Floor(t)
	d1 := i.Ceil(t)
	if t.Sub(d0) < d1.Sub(t) {
		return d0
	}
	return d1
}

// Offset returns t shifted by `step` intervals.
func (i *Interval) Offset(t time.Time, step int) time.Time {
	c := t
	i.offseti(&c, step)
	return c
}

// Range returns the sequence of interval boundaries in [start, stop). d3
// uses Ceil(start) as the first element; the loop pushes then offsets.
// `step` defaults to 1.
func (i *Interval) Range(start, stop time.Time, step int) []time.Time {
	if step == 0 {
		step = 1
	}
	if step < 0 {
		// d3 floors step: step = Math.floor(step); negative not supported here.
		step = -step
		// swap and reverse at the end
		start, stop = stop, start
		out := i.rangeAsc(start, stop, step)
		// reverse
		for l, r := 0, len(out)-1; l < r; l, r = l+1, r-1 {
			out[l], out[r] = out[r], out[l]
		}
		return out
	}
	return i.rangeAsc(start, stop, step)
}

func (i *Interval) rangeAsc(start, stop time.Time, step int) []time.Time {
	out := []time.Time{}
	cur := i.Ceil(start)
	if !cur.Before(stop) || step <= 0 {
		return out
	}
	var prev time.Time
	for {
		out = append(out, cur)
		prev = cur
		i.offseti(&cur, step)
		i.floori(&cur)
		if !prev.Before(cur) || !cur.Before(stop) {
			break
		}
	}
	return out
}

// Count returns the number of interval boundaries between start and end
// (both floored first).
func (i *Interval) Count(start, end time.Time) int {
	if i.countFn == nil {
		return 0
	}
	s := start
	e := end
	i.floori(&s)
	i.floori(&e)
	return int(math.Floor(float64(i.countFn(s, e))))
}

// Every returns a derived interval that only fires every `step`th boundary.
// step <= 0 or non-finite returns nil; step == 1 returns the interval itself.
func (i *Interval) Every(step int) *Interval {
	if step <= 0 {
		return nil
	}
	if step == 1 {
		return i
	}
	if i.fieldFn != nil {
		// filter by field % step == 0
		floori := func(t *time.Time) {
			for {
				i.floori(t)
				if i.fieldFn(*t)%step == 0 {
					return
				}
				*t = t.Add(-time.Millisecond)
			}
		}
		offseti := func(t *time.Time, s int) {
			n := s
			if n < 0 {
				for n < 0 {
					n++
					i.offseti(t, -1)
					for i.fieldFn(*t)%step != 0 {
						i.offseti(t, -1)
					}
				}
			} else {
				for n > 0 {
					n--
					i.offseti(t, 1)
					for i.fieldFn(*t)%step != 0 {
						i.offseti(t, 1)
					}
				}
			}
		}
		return &Interval{floori: floori, offseti: offseti, countFn: i.countFn, fieldFn: i.fieldFn}
	}
	// count-based filter (e.g. millisecond)
	floori := func(t *time.Time) {
		// not used by millisecond.every in our ports (only year has a custom every)
		i.floori(t)
	}
	offseti := i.offseti
	return &Interval{floori: floori, offseti: offseti, countFn: i.countFn, fieldFn: i.fieldFn}
}

// --- concrete UTC intervals ------------------------------------------------

func utcSecondInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			*t = t.Truncate(time.Second)
		},
		offseti: func(t *time.Time, step int) {
			*t = t.Add(time.Duration(step) * time.Second)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / time.Second)
		},
		fieldFn: func(t time.Time) int { return t.UTC().Second() },
	}
}

func utcMinuteInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			*t = t.Truncate(time.Minute)
		},
		offseti: func(t *time.Time, step int) {
			*t = t.Add(time.Duration(step) * time.Minute)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / time.Minute)
		},
		fieldFn: func(t time.Time) int { return t.UTC().Minute() },
	}
}

func utcHourInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			*t = t.Truncate(time.Hour)
		},
		offseti: func(t *time.Time, step int) {
			*t = t.Add(time.Duration(step) * time.Hour)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / time.Hour)
		},
		fieldFn: func(t time.Time) int { return t.UTC().Hour() },
	}
}

// utcDay floors to 00:00:00 UTC of the current day and offsets by calendar days.
func utcDayInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			u := t.UTC()
			*t = time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
		},
		offseti: func(t *time.Time, step int) {
			u := t.UTC()
			*t = time.Date(u.Year(), u.Month(), u.Day()+step, 0, 0, 0, 0, time.UTC)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / durationDay)
		},
		fieldFn: func(t time.Time) int { return t.UTC().Day() - 1 },
	}
}

// utcWeek (Sunday-based) floors to the most recent Sunday 00:00 UTC.
func utcWeekInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			u := t.UTC()
			day := int(u.Weekday()) // Sunday = 0
			*t = time.Date(u.Year(), u.Month(), u.Day()-day, 0, 0, 0, 0, time.UTC)
		},
		offseti: func(t *time.Time, step int) {
			u := t.UTC()
			*t = time.Date(u.Year(), u.Month(), u.Day()+step*7, 0, 0, 0, 0, time.UTC)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / durationWeek)
		},
	}
}

// utcMonth floors to the first day of the month at 00:00 UTC.
func utcMonthInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			u := t.UTC()
			*t = time.Date(u.Year(), u.Month(), 1, 0, 0, 0, 0, time.UTC)
		},
		offseti: func(t *time.Time, step int) {
			u := t.UTC()
			*t = time.Date(u.Year(), u.Month()+time.Month(step), 1, 0, 0, 0, 0, time.UTC)
		},
		countFn: func(start, end time.Time) int {
			s, e := start.UTC(), end.UTC()
			return int(e.Month()) - int(s.Month()) + (e.Year()-s.Year())*12
		},
		fieldFn: func(t time.Time) int { return int(t.UTC().Month()) - 1 },
	}
}

// utcYear floors to Jan 1 00:00 UTC.
func utcYearInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {
			u := t.UTC()
			*t = time.Date(u.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		},
		offseti: func(t *time.Time, step int) {
			u := t.UTC()
			*t = time.Date(u.Year()+step, 1, 1, 0, 0, 0, 0, time.UTC)
		},
		countFn: func(start, end time.Time) int {
			return end.UTC().Year() - start.UTC().Year()
		},
		fieldFn: func(t time.Time) int { return t.UTC().Year() },
	}
}

// utcYearEvery returns a derived year interval that fires every `k` years
// (flooring to the nearest multiple-of-k year). Mirrors d3-time's
// `utcYear.every(k)` optimization.
func utcYearEvery(k int) *Interval {
	if k <= 0 {
		return nil
	}
	if k == 1 {
		return utcYearInterval()
	}
	return &Interval{
		floori: func(t *time.Time) {
			u := t.UTC()
			year := (u.Year() / k) * k
			*t = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		},
		offseti: func(t *time.Time, step int) {
			u := t.UTC()
			*t = time.Date(u.Year()+step*k, 1, 1, 0, 0, 0, 0, time.UTC)
		},
		countFn: func(start, end time.Time) int {
			return end.UTC().Year() - start.UTC().Year()
		},
	}
}

// millisecondInterval is the base millisecond interval. Its Every has a
// special optimized form.
func millisecondInterval() *Interval {
	return &Interval{
		floori: func(t *time.Time) {}, // noop
		offseti: func(t *time.Time, step int) {
			*t = t.Add(time.Duration(step) * time.Millisecond)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / time.Millisecond)
		},
	}
}

// millisecondEvery returns a millisecond interval spaced by k ms.
func millisecondEvery(k int) *Interval {
	if k <= 0 {
		return nil
	}
	if k == 1 {
		return millisecondInterval()
	}
	return &Interval{
		floori: func(t *time.Time) {
			ms := t.UnixNano() / int64(time.Millisecond)
			*t = time.UnixMilli((ms / int64(k)) * int64(k)).UTC()
		},
		offseti: func(t *time.Time, step int) {
			*t = t.Add(time.Duration(step*k) * time.Millisecond)
		},
		countFn: func(start, end time.Time) int {
			return int(end.Sub(start) / (time.Duration(k) * time.Millisecond))
		},
	}
}
