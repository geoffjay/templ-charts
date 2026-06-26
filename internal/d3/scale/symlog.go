// symlog.go — port of d3-scale's scaleSymlog.
//
// A symlog (symmetric log) scale is a continuous scale that handles both
// positive and negative (and zero) domains by applying the transform
//
//	t(x) = sign(x) * log1p(|x / c|)
//
// before the linear domain→range mapping. The inverse is
//
//	x(t) = sign(t) * expm1(|t|) * c
//
// with `c` a configurable constant (default 1). Ticks and Nice reuse the
// linear-scale algorithm (linearish).
package d3scale

import "math"

// Symlog is a continuous symmetric-log scale.
type Symlog struct {
	domain        [2]float64
	rangeVals     [2]float64
	constant      float64
	clamp         bool
	round         bool
	interpolateFn func(float64) float64
	transformFn   func(float64) float64
	untransformFn func(float64) float64
}

// NewSymlog constructs a symlog scale with constant=1, domain [0, 1], range [0, 1].
func NewSymlog() *Symlog {
	s := &Symlog{
		domain:    [2]float64{0, 1},
		rangeVals: [2]float64{0, 1},
		constant:  1,
	}
	s.rebuildTransform()
	s.rebuildInterpolator()
	return s
}

// Type returns "symlog".
func (s *Symlog) Type() string { return "symlog" }

// Constant returns the symlog constant.
func (s *Symlog) Constant() float64 { return s.constant }

// SetConstant sets the symlog constant and rebuilds the transform.
func (s *Symlog) SetConstant(c float64) *Symlog {
	s.constant = c
	s.rebuildTransform()
	return s
}

// Domain returns the domain.
func (s *Symlog) Domain() [2]float64 { return s.domain }

// SetDomain sets the domain.
func (s *Symlog) SetDomain(d0, d1 float64) *Symlog {
	s.domain = [2]float64{d0, d1}
	return s
}

// Range returns the range.
func (s *Symlog) Range() [2]float64 { return s.rangeVals }

// SetRange sets the range.
func (s *Symlog) SetRange(r0, r1 float64) *Symlog {
	s.rangeVals = [2]float64{r0, r1}
	s.round = false
	s.rebuildInterpolator()
	return s
}

// SetRangeRound sets the range with rounding.
func (s *Symlog) SetRangeRound(r0, r1 float64) *Symlog {
	s.rangeVals = [2]float64{r0, r1}
	s.round = true
	s.rebuildInterpolator()
	return s
}

// Clamp reports clamping state.
func (s *Symlog) Clamp() bool { return s.clamp }

// SetClamp enables/disables clamping.
func (s *Symlog) SetClamp(c bool) *Symlog {
	s.clamp = c
	return s
}

// Call maps a domain value to a range value.
func (s *Symlog) Call(x float64) float64 {
	xc := x
	if s.clamp {
		xc = clampSymlogDomain(x, s.domain)
	}
	tx := s.transformFn(xc)
	t := deinterpolate(s.transformFn(s.domain[0]), s.transformFn(s.domain[1]))(tx)
	if s.clamp {
		t = clamp01(t)
	} else if math.IsNaN(t) {
		return t
	}
	return s.interpolateFn(t)
}

// Invert maps a range value back to a domain value.
func (s *Symlog) Invert(y float64) float64 {
	r0, r1 := s.rangeVals[0], s.rangeVals[1]
	rd := r1 - r0
	if rd == 0 {
		return s.domain[0]
	}
	t := (y - r0) / rd
	if s.clamp {
		t = clamp01(t)
	}
	td0 := s.transformFn(s.domain[0])
	td1 := s.transformFn(s.domain[1])
	return s.untransformFn(td0 + t*(td1-td0))
}

// Ticks returns approximately `count` tick values. Symlog reuses linearish
// ticks (d3array.Ticks over the raw domain).
func (s *Symlog) Ticks(count int) []float64 {
	return ticksFloats(s.domain[0], s.domain[1], tickCount(count))
}

// Nice extends the domain to nice round values (linearish Nice).
func (s *Symlog) Nice(count int) *Symlog {
	s.domain = niceLinear(s.domain, tickCount(count))
	return s
}

// Copy returns a deep copy.
func (s *Symlog) Copy() *Symlog {
	c := *s
	c.rebuildTransform()
	c.rebuildInterpolator()
	return &c
}

func (s *Symlog) rebuildTransform() {
	c := s.constant
	s.transformFn = func(x float64) float64 {
		return math.Copysign(1, x) * math.Log1p(math.Abs(x/c))
	}
	s.untransformFn = func(x float64) float64 {
		return math.Copysign(1, x) * math.Expm1(math.Abs(x)) * c
	}
}

func (s *Symlog) rebuildInterpolator() {
	if s.round {
		s.interpolateFn = interpolateRound(s.rangeVals[0], s.rangeVals[1])
	} else {
		s.interpolateFn = interpolateNumber(s.rangeVals[0], s.rangeVals[1])
	}
}

// clampSymlogDomain clamps x to the domain's min/max.
func clampSymlogDomain(x float64, domain [2]float64) float64 {
	lo, hi := domain[0], domain[1]
	if hi < lo {
		lo, hi = hi, lo
	}
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
