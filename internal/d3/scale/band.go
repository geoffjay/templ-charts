// band.go — port of d3-scale's scaleBand (and scalePoint via Point()).
//
// Discrete scale: maps a string domain to evenly-spaced bands in a numeric
// range. The algorithm is a faithful transliteration of d3-scale's band.js:
//
//	step = (stop - start) / max(1, n - paddingInner + 2*paddingOuter)
//	start += (stop - start - step*(n - paddingInner)) * align
//	bandwidth = step * (1 - paddingInner)
//
// with optional rounding of step/start/bandwidth. Call(x) returns the band's
// range start. The Point scale is a Band with paddingInner=1 (bands collapse
// to points) and exposes paddingOuter as "padding".
package d3scale

import "math"

// Band is a discrete band scale mapping string → float64.
type Band struct {
	domain       []string
	index        map[string]int // position of each domain value
	r0, r1       float64
	step         float64
	bandwidth    float64
	round        bool
	paddingInner float64
	paddingOuter float64
	align        float64
	values       []float64 // cached range positions per domain index
}

// NewBand constructs a band scale with empty domain, range [0,1], and
// d3 defaults (round=false, paddingInner=0, paddingOuter=0, align=0.5).
func NewBand() *Band {
	b := &Band{
		domain:       []string{},
		index:        map[string]int{},
		r0:           0,
		r1:           1,
		paddingInner: 0,
		paddingOuter: 0,
		align:        0.5,
	}
	b.rescale()
	return b
}

// Type returns "band".
func (b *Band) Type() string { return "band" }

// Domain returns the domain slice.
func (b *Band) Domain() []string { return append([]string(nil), b.domain...) }

// SetDomain sets the domain (dedup preserving first-seen order, matching
// d3's InternMap behavior) and rescales.
func (b *Band) SetDomain(domain []string) *Band {
	b.domain = b.domain[:0]
	b.index = map[string]int{}
	for _, v := range domain {
		if _, ok := b.index[v]; ok {
			continue
		}
		b.index[v] = len(b.domain)
		b.domain = append(b.domain, v)
	}
	b.rescale()
	return b
}

// Range returns the current range.
func (b *Band) Range() [2]float64 { return [2]float64{b.r0, b.r1} }

// SetRange sets the range and rescales.
func (b *Band) SetRange(r0, r1 float64) *Band {
	b.r0 = r0
	b.r1 = r1
	b.rescale()
	return b
}

// SetRangeRound sets the range and enables rounding.
func (b *Band) SetRangeRound(r0, r1 float64) *Band {
	b.r0 = r0
	b.r1 = r1
	b.round = true
	b.rescale()
	return b
}

// Bandwidth returns the band width.
func (b *Band) Bandwidth() float64 { return b.bandwidth }

// Step returns the step between band starts.
func (b *Band) Step() float64 { return b.step }

// Round reports whether rounding is enabled.
func (b *Band) Round() bool { return b.round }

// SetRound enables/disables integer rounding.
func (b *Band) SetRound(r bool) *Band {
	b.round = r
	b.rescale()
	return b
}

// PaddingInner returns the inner padding ratio.
func (b *Band) PaddingInner() float64 { return b.paddingInner }

// SetPaddingInner sets inner padding (clamped to [0,1]).
func (b *Band) SetPaddingInner(p float64) *Band {
	b.paddingInner = math.Min(1, p)
	b.rescale()
	return b
}

// PaddingOuter returns the outer padding ratio.
func (b *Band) PaddingOuter() float64 { return b.paddingOuter }

// SetPaddingOuter sets outer padding.
func (b *Band) SetPaddingOuter(p float64) *Band {
	b.paddingOuter = p
	b.rescale()
	return b
}

// SetPadding sets both inner and outer padding (d3's .padding()).
func (b *Band) SetPadding(p float64) *Band {
	b.paddingInner = math.Min(1, p)
	b.paddingOuter = p
	b.rescale()
	return b
}

// Align returns the alignment (0..1).
func (b *Band) Align() float64 { return b.align }

// SetAlign sets alignment (clamped to [0,1]).
func (b *Band) SetAlign(a float64) *Band {
	b.align = math.Max(0, math.Min(1, a))
	b.rescale()
	return b
}

// Call maps a domain value to its band's range start. Returns NaN for
// unknown values (d3 returns undefined; we use NaN since float64 has no
// nil). Callers can distinguish via the index map if needed.
func (b *Band) Call(x string) float64 {
	i, ok := b.index[x]
	if !ok || i >= len(b.values) {
		return math.NaN()
	}
	return b.values[i]
}

// Copy returns a deep copy.
func (b *Band) Copy() *Band {
	c := &Band{
		domain:       append([]string(nil), b.domain...),
		index:        make(map[string]int, len(b.domain)),
		r0:           b.r0,
		r1:           b.r1,
		step:         b.step,
		bandwidth:    b.bandwidth,
		round:        b.round,
		paddingInner: b.paddingInner,
		paddingOuter: b.paddingOuter,
		align:        b.align,
	}
	for k, v := range b.index {
		c.index[k] = v
	}
	c.rescale()
	return c
}

// rescale recomputes step, bandwidth, and per-domain range values.
// Faithful transliteration of d3-scale band.js rescale().
func (b *Band) rescale() {
	n := len(b.domain)
	reverse := b.r1 < b.r0
	start, stop := b.r0, b.r1
	if reverse {
		start, stop = b.r1, b.r0
	}
	denom := float64(n) - b.paddingInner + b.paddingOuter*2
	if denom < 1 {
		denom = 1
	}
	b.step = (stop - start) / denom
	if b.round {
		b.step = math.Floor(b.step)
	}
	// shift start by the alignment slack
	start += (stop - start - b.step*(float64(n)-b.paddingInner)) * b.align
	b.bandwidth = b.step * (1 - b.paddingInner)
	if b.round {
		start = math.Round(start)
		b.bandwidth = math.Round(b.bandwidth)
	}
	b.values = make([]float64, n)
	for i := 0; i < n; i++ {
		b.values[i] = start + b.step*float64(i)
	}
	if reverse {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			b.values[i], b.values[j] = b.values[j], b.values[i]
		}
	}
}

// --- Point scale -----------------------------------------------------------
//
// scalePoint = scaleBand with paddingInner=1 and paddingOuter exposed as
// "padding". Bandwidth is always 0.

// Point is a discrete point scale mapping string → float64.
type Point struct {
	*Band
}

// NewPoint constructs a point scale with empty domain, range [0,1].
func NewPoint() *Point {
	p := &Point{Band: NewBand()}
	p.Band.paddingInner = 1
	p.Band.rescale()
	return p
}

// Type returns "point".
func (p *Point) Type() string { return "point" }

// SetDomain sets the domain.
func (p *Point) SetDomain(domain []string) *Point {
	p.Band.SetDomain(domain)
	return p
}

// SetRange sets the range.
func (p *Point) SetRange(r0, r1 float64) *Point {
	p.Band.SetRange(r0, r1)
	return p
}

// SetRangeRound sets the range with rounding.
func (p *Point) SetRangeRound(r0, r1 float64) *Point {
	p.Band.SetRangeRound(r0, r1)
	return p
}

// SetRound enables/disables rounding.
func (p *Point) SetRound(r bool) *Point {
	p.Band.SetRound(r)
	return p
}

// Padding returns the outer padding (d3 point.padding === band.paddingOuter).
func (p *Point) Padding() float64 { return p.Band.paddingOuter }

// SetPadding sets the outer padding for the point scale (d3 point.padding).
func (p *Point) SetPadding(pad float64) *Point {
	p.Band.paddingOuter = pad
	p.Band.rescale()
	return p
}

// SetAlign sets alignment.
func (p *Point) SetAlign(a float64) *Point {
	p.Band.SetAlign(a)
	return p
}

// Bandwidth returns 0 (points have no width).
func (p *Point) Bandwidth() float64 { return 0 }

// Copy returns a deep copy.
func (p *Point) Copy() *Point {
	c := &Point{Band: NewBand()}
	c.Band.domain = append([]string(nil), p.Band.domain...)
	c.Band.index = make(map[string]int, len(p.Band.domain))
	for k, v := range p.Band.index {
		c.Band.index[k] = v
	}
	c.Band.r0 = p.Band.r0
	c.Band.r1 = p.Band.r1
	c.Band.round = p.Band.round
	c.Band.paddingInner = 1 // always 1 for point
	c.Band.paddingOuter = p.Band.paddingOuter
	c.Band.align = p.Band.align
	c.Band.rescale()
	return c
}
