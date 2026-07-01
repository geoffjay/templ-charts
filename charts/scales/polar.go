package scales

import (
	d3scale "github.com/geoffjay/templ-charts/internal/d3/scale"
)

// NewLinearScaleWithRange builds a linear Scale mapping the domain
// [domainMin, domainMax] onto the explicit output range [rangeMin, rangeMax].
//
// Unlike ComputeScale, the range is given directly rather than derived from a
// size + axis orientation, so callers can map onto arbitrary targets — e.g. an
// angle span (degrees) for a polar/radial chart. The returned value is a full
// scaleImpl, so GetScaleTicks and CenterScale operate on it (polar-axes needs
// this for its radial-grid rays and circular-axis ticks).
func NewLinearScaleWithRange(domainMin, domainMax, rangeMin, rangeMax float64) Scale {
	s := d3scale.NewLinear()
	s.SetDomain(domainMin, domainMax)
	s.SetRange(rangeMin, rangeMax)
	return &scaleImpl{typ: ScaleTypeLinear, linear: s}
}

// NewBandScaleWithRange builds a band Scale over `domain` mapped onto the
// explicit output range [rangeMin, rangeMax] with the given padding (applied as
// both inner and outer padding, matching d3's scaleBand.padding()).
//
// Like NewLinearScaleWithRange, this exists because the range is not derivable
// from a size + axis: radial charts map series ids onto a radius band
// [innerRadius, outerRadius]. The returned value is a full scaleImpl with
// Bandwidth/Step, so it composes with polar-axes PolarGrid/CircularAxis.
func NewBandScaleWithRange(domain []string, rangeMin, rangeMax, padding float64, round bool) Scale {
	s := d3scale.NewBand().SetDomain(domain)
	s.SetRange(rangeMin, rangeMax)
	s.SetPadding(padding)
	if round {
		s.SetRound(true)
	}
	return &scaleImpl{typ: ScaleTypeBand, band: s}
}

// NewPointScaleWithRange builds a point Scale over `domain` mapped onto the
// explicit output range [rangeMin, rangeMax] with the given outer padding.
//
// Like the linear/band range constructors, this exists because the range is
// not derivable from a size + axis alone: charts that place one axis per
// variable (parallel-coordinates) map each variable's categories onto an
// arbitrary pixel span. The returned value is a full scaleImpl with
// Bandwidth/Step, so it composes with the axes tick machinery.
func NewPointScaleWithRange(domain []string, rangeMin, rangeMax, padding float64, round bool) Scale {
	s := d3scale.NewPoint().SetDomain(domain)
	s.SetRange(rangeMin, rangeMax)
	s.SetPadding(padding)
	if round {
		s.SetRound(true)
	}
	return &scaleImpl{typ: ScaleTypePoint, point: s}
}
