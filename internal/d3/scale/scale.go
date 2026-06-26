// Package d3scale provides a minimal port of d3-scale covering the scale
// types nivo uses: linear, band, point, time, log, and symlog.
//
// Each scale mirrors the d3-scale fluent-builder API: methods set state and
// return the scale (as a typed pointer, e.g. *Linear), and `Call(x)` produces
// the output value. Discrete scales (band, point) also expose `Bandwidth`,
// `Step`, and `Padding*` accessors.
//
// Numerical tick generation is shared via d3array.Ticks; time ticks use a
// dedicated interval-based algorithm in time_ticks.go.
package d3scale

import (
	"math"

	d3array "github.com/geoffjay/templ-charts/internal/d3/array"
)

// --- shared helpers --------------------------------------------------------

// interpolateNumber linearly interpolates between a and b: returns a function
// f(t) with t in [0,1]. Mirrors d3-interpolate's interpolateNumber.
func interpolateNumber(a, b float64) func(float64) float64 {
	return func(t float64) float64 { return a + (b-a)*t }
}

// interpolateRound is interpolateNumber with rounding to the nearest integer.
func interpolateRound(a, b float64) func(float64) float64 {
	return func(t float64) float64 { return float64(roundInt(a + (b-a)*t)) }
}

func roundInt(v float64) int { return int(math.Floor(v + 0.5)) }

// deinterpolate returns a function f(x) producing the normalized position of
// x within [d0, d1] (i.e. (x - d0) / (d1 - d0)). Used by Call/Invert.
func deinterpolate(d0, d1 float64) func(float64) float64 {
	d := d1 - d0
	if d == 0 {
		return func(float64) float64 { return 0 }
	}
	return func(x float64) float64 { return (x - d0) / d }
}

// clamp01 clamps to [0,1].
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// tickCount returns a sane count, defaulting to 10 (d3 default).
func tickCount(count int) int {
	if count <= 0 {
		return 10
	}
	return count
}

// ticksFloats is a thin wrapper over d3array.Ticks returning []float64.
func ticksFloats(start, stop float64, count int) []float64 {
	return d3array.Ticks(start, stop, count)
}

// tickIncrement is d3-scale's tickIncrement (the step size d3 uses to extend
// the domain via nice()). Delegates to d3array.TickIncrement which implements
// the same 1/2/5 × 10^k algorithm.
func tickIncrement(start, stop float64, count int) float64 {
	return d3array.TickIncrement(start, stop, count)
}
