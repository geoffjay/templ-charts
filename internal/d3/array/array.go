// Package d3array provides a minimal port of the d3-array helpers used by
// templ-charts. Most nivo data manipulation uses lodash; d3-array is
// primarily consumed by d3-scale's tick generation. We port the subset
// actually needed: extent, ticks (the "nice tick count" algorithm),
// bisect/bisector, sum, max, min, mean, median, range, ascending/descending,
// and the Numberable filter used by quantize.
package d3array

import (
	"math"
	"sort"
)

// Extent returns the [min, max] of vs. Empty input returns [math.NaN(),
// math.NaN()]. NaN values are ignored.
func Extent(vs []float64) [2]float64 {
	min, max := math.NaN(), math.NaN()
	first := true
	for _, v := range vs {
		if math.IsNaN(v) {
			continue
		}
		if first {
			min, max = v, v
			first = false
			continue
		}
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return [2]float64{min, max}
}

// Max returns the maximum of vs or NaN if empty (NaNs ignored).
func Max(vs []float64) float64 {
	e := Extent(vs)
	return e[1]
}

// Min returns the minimum of vs or NaN if empty (NaNs ignored).
func Min(vs []float64) float64 {
	e := Extent(vs)
	return e[0]
}

// Sum returns the sum of vs (NaNs treated as 0).
func Sum(vs []float64) float64 {
	var s float64
	for _, v := range vs {
		if !math.IsNaN(v) {
			s += v
		}
	}
	return s
}

// Mean returns the arithmetic mean of vs, or NaN if empty (NaNs ignored).
func Mean(vs []float64) float64 {
	var s float64
	n := 0
	for _, v := range vs {
		if math.IsNaN(v) {
			continue
		}
		s += v
		n++
	}
	if n == 0 {
		return math.NaN()
	}
	return s / float64(n)
}

// Median returns the median of vs (NaNs ignored). Returns NaN if empty.
// Allocates a sorted copy.
func Median(vs []float64) float64 {
	cleaned := filterNaNs(vs)
	if len(cleaned) == 0 {
		return math.NaN()
	}
	sorted := make([]float64, len(cleaned))
	copy(sorted, cleaned)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 1 {
		return sorted[n/2]
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// Range returns a slice [start, start+step, ..., stop) (half-open).
// Step defaults to 1; if step is 0 returns an empty slice.
// Matches d3-array's range semantics.
func Range(args ...float64) []float64 {
	var start, stop, step float64
	switch len(args) {
	case 1:
		start, stop, step = 0, args[0], 1
	case 2:
		start, stop, step = args[0], args[1], 1
	case 3:
		start, stop, step = args[0], args[1], args[2]
	default:
		return nil
	}
	if step == 0 {
		return nil
	}
	// d3 computes the count as ceil((stop - start) / step) and rejects
	// non-positive counts.
	n := int(math.Ceil((stop - start) / step))
	if n <= 0 {
		return []float64{}
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = start + float64(i)*step
	}
	return out
}

// Ascending is the d3-array comparator: -1, 0, 1, with NaN treated as last.
func Ascending(a, b float64) int {
	if math.IsNaN(a) {
		if math.IsNaN(b) {
			return 0
		}
		return 1
	}
	if math.IsNaN(b) {
		return -1
	}
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// Descending is the inverse of Ascending.
func Descending(a, b float64) int {
	return -Ascending(a, b)
}

// Bisector returns a function that finds the insertion index for a value
// in an already-sorted slice using the given comparator. d3-array's
// bisector wraps a comparator and exposes left/right variants; here we
// expose BisectLeft / BisectRight directly for clarity.
//
// All bisect operations assume `vs` is sorted ascending per the comparator.
type Bisector func(a, b float64) int

// BisectLeft returns the index at which to insert x to maintain sorted
// order; equal values go before existing ones (i.e. leftmost insertion
// point).
func BisectLeft(vs []float64, x float64) int {
	return bisectLeft(vs, x, 0, len(vs))
}

// BisectRight returns the index at which to insert x to maintain sorted
// order; equal values go after existing ones (i.e. rightmost insertion
// point).
func BisectRight(vs []float64, x float64) int {
	return bisectRight(vs, x, 0, len(vs))
}

func bisectLeft(vs []float64, x float64, lo, hi int) int {
	for lo < hi {
		mid := (lo + hi) >> 1
		if Ascending(vs[mid], x) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func bisectRight(vs []float64, x float64, lo, hi int) int {
	for lo < hi {
		mid := (lo + hi) >> 1
		if Ascending(x, vs[mid]) < 0 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

// Ticks returns an array of approximately `count` "nice" tick values that
// span [start, stop]. This is d3-array's `ticks()` algorithm (also used
// internally by d3-scale's linear.ticks()).
//
// It computes a step size as 10^k * {1,2,5,10} (whichever is closest to the
// raw step), then walks from ceil(start/step)*step to stop, inclusive,
// accumulating tick values. Returns an empty slice if start or stop is NaN
// or if count is non-positive.
func Ticks(start, stop float64, count int) []float64 {
	if math.IsNaN(start) || math.IsNaN(stop) || count <= 0 {
		return []float64{}
	}
	span := stop - start
	if math.IsInf(span, 0) {
		return []float64{}
	}
	if span == 0 {
		return []float64{start}
	}
	step := tickStep(start, stop, count)
	// If the step is too small to be useful (sub-divides below the input
	// precision), fall back to returning just the endpoints.
	if step == 0 || math.IsNaN(step) {
		return []float64{start, stop}
	}
	// Generate ticks: walk from ceil(start/step)*step, advancing by step,
	// until we exceed stop. d3 uses an incremental accumulator to avoid
	// floating-point drift.
	out := []float64{}
	v := math.Ceil(start/step) * step
	if v < start {
		v = start
	}
	// include the starting tick
	out = append(out, v)
	for {
		v += step
		if v > stop {
			break
		}
		out = append(out, v)
	}
	// ensure the final tick is exactly stop (or close to it) when the
	// accumulated value crossed it
	if len(out) > 0 && out[len(out)-1] < stop && math.Abs(out[len(out)-1]+step-stop) < step*1e-6 {
		out = append(out, stop)
	}
	return out
}

// tickStep computes a "nice" step size (one of 1, 2, 5 × 10^k) closest to
// (stop-start)/count. This is the core of d3-array's ticks algorithm.
func tickStep(start, stop float64, count int) float64 {
	step0 := math.Abs(stop-start) / float64(count)
	step1 := math.Pow10(int(math.Floor(math.Log10(step0))))
	err := step0 / step1
	switch {
	case err >= 7.5:
		return step1 * 10
	case err >= 3.5:
		return step1 * 5
	case err >= 1.5:
		return step1 * 2
	default:
		return step1
	}
}

// TickIncrement returns the increment d3-scale uses for ticks() given a
// count; same as tickStep but returns the rounded step. Exported for
// scale.nice().
func TickIncrement(start, stop float64, count int) float64 {
	return tickStep(start, stop, count)
}

// Numberable reports whether v is a finite number (used by quantize to
// filter out undefined/NaN inputs).
func Numberable(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// filterNaNs returns a new slice with NaNs removed.
func filterNaNs(vs []float64) []float64 {
	out := make([]float64, 0, len(vs))
	for _, v := range vs {
		if !math.IsNaN(v) {
			out = append(out, v)
		}
	}
	return out
}

// Shuffle returns a new slice with the elements of vs in a Fisher-Yates
// shuffled order using the provided rand function. If rand is nil, uses
// math/rand's default. (Used rarely; included for completeness since the
// nivo generators package consumes it.)
//
// Not currently used by templ-charts v1; preserved as a helper for any
// future synthetic data generator.
func Shuffle(vs []float64, rand func(n int) int) []float64 {
	out := make([]float64, len(vs))
	copy(out, vs)
	if rand == nil {
		// simple deterministic shuffle without math/rand dep
		seed := uint64(1)
		next := func() uint64 {
			seed = seed*1103515245 + 12345
			return seed
		}
		rand = func(n int) int { return int(next() % uint64(n)) }
	}
	for i := len(out) - 1; i > 0; i-- {
		j := rand(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return out
}
