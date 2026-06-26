// linear.go — port of d3-scale's scaleLinear.
//
// d3-scale linear: domain [d0,d1] → range [r0,r1] via
//
//	r = r0 + (r1 - r0) * (x - d0) / (d1 - d0)
//
// with optional clamping and "nice" domain extension. Tick generation
// delegates to d3array.Ticks (the same algorithm d3-scale uses).
package d3scale

import "math"

// Linear is a continuous linear scale mapping float64 → float64.
type Linear struct {
	domain        [2]float64
	rangeVals     [2]float64
	clamp         bool
	round         bool
	interpolateFn func(float64) float64
}

// NewLinear constructs a linear scale with domain [0,1] and range [0,1].
func NewLinear() *Linear {
	l := &Linear{
		domain:    [2]float64{0, 1},
		rangeVals: [2]float64{0, 1},
	}
	l.rebuildInterpolator()
	return l
}

// Type returns "linear".
func (l *Linear) Type() string { return "linear" }

// Domain returns the domain.
func (l *Linear) Domain() [2]float64 { return l.domain }

// SetDomain sets the domain.
func (l *Linear) SetDomain(d0, d1 float64) *Linear {
	l.domain = [2]float64{d0, d1}
	l.rebuildInterpolator()
	return l
}

// Range returns the current range.
func (l *Linear) Range() [2]float64 { return l.rangeVals }

// SetRange sets the range and rebuilds the interpolator (non-rounding).
func (l *Linear) SetRange(r0, r1 float64) *Linear {
	l.rangeVals = [2]float64{r0, r1}
	l.round = false
	l.rebuildInterpolator()
	return l
}

// SetRangeRound sets the range and enables integer rounding of outputs.
func (l *Linear) SetRangeRound(r0, r1 float64) *Linear {
	l.rangeVals = [2]float64{r0, r1}
	l.round = true
	l.rebuildInterpolator()
	return l
}

// Clamp reports clamping state.
func (l *Linear) Clamp() bool { return l.clamp }

// SetClamp enables/disables clamping.
func (l *Linear) SetClamp(c bool) *Linear {
	l.clamp = c
	return l
}

// Call maps a domain value (float64) to a range value.
func (l *Linear) Call(x float64) float64 {
	t := deinterpolate(l.domain[0], l.domain[1])(x)
	if l.clamp {
		t = clamp01(t)
	} else if math.IsNaN(t) {
		return t
	}
	return l.interpolateFn(t)
}

// Invert maps a range value back to a domain value.
func (l *Linear) Invert(y float64) float64 {
	r0, r1 := l.rangeVals[0], l.rangeVals[1]
	rd := r1 - r0
	if rd == 0 {
		return l.domain[0]
	}
	t := (y - r0) / rd
	if l.clamp {
		t = clamp01(t)
	}
	return l.domain[0] + t*(l.domain[1]-l.domain[0])
}

// Ticks returns approximately `count` nice tick values spanning the domain.
func (l *Linear) Ticks(count int) []float64 {
	return ticksFloats(l.domain[0], l.domain[1], tickCount(count))
}

// Nice extends the domain to nice round values. With count <= 0 uses d3's
// default count of 10.
func (l *Linear) Nice(count int) *Linear {
	l.domain = niceLinear(l.domain, tickCount(count))
	l.rebuildInterpolator()
	return l
}

// Copy returns a deep copy.
func (l *Linear) Copy() *Linear {
	c := *l
	c.rebuildInterpolator()
	return &c
}

// rebuildInterpolator (re)creates the range interpolator based on round.
func (l *Linear) rebuildInterpolator() {
	if l.round {
		l.interpolateFn = interpolateRound(l.rangeVals[0], l.rangeVals[1])
	} else {
		l.interpolateFn = interpolateNumber(l.rangeVals[0], l.rangeVals[1])
	}
}

// niceLinear mirrors d3-scale's niceLinear: pick a tick step for the domain
// span, then extend both ends to the nearest multiple of that step.
func niceLinear(domain [2]float64, count int) [2]float64 {
	d0, d1 := domain[0], domain[1]
	reverse := d1 < d0
	if reverse {
		d0, d1 = d1, d0
	}
	step := tickIncrement(d0, d1, count)
	if step == 0 || math.IsNaN(step) {
		return domain
	}
	n0 := math.Floor(d0/step) * step
	n1 := math.Ceil(d1/step) * step
	if !reverse {
		return [2]float64{n0, n1}
	}
	return [2]float64{n1, n0}
}
