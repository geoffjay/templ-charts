package scales

import (
	"fmt"
	"math"
	"strconv"
	"time"

	d3scale "github.com/geoffjay/templ-charts/internal/d3/scale"
)

// Scale is a resolved scale wrapper. The underlying d3-scale instance is held
// in Impl; Call dispatches by type.
type scaleImpl struct {
	typ     ScaleType
	linear  *d3scale.Linear
	log     *d3scale.Log
	symlog  *d3scale.Symlog
	band    *d3scale.Band
	point   *d3scale.Point
	time    *d3scale.Time
	stacked bool
	useUTC  bool
}

func (s *scaleImpl) Type() ScaleType { return s.typ }

func (s *scaleImpl) Call(v any) float64 {
	switch s.typ {
	case ScaleTypeLinear:
		return s.linear.Call(toFloat(v))
	case ScaleTypeLog:
		return s.log.Call(toFloat(v))
	case ScaleTypeSymlog:
		return s.symlog.Call(toFloat(v))
	case ScaleTypeBand:
		return s.band.Call(toString(v))
	case ScaleTypePoint:
		return s.point.Call(toString(v))
	case ScaleTypeTime:
		if t, ok := toTime(v); ok {
			return s.time.Call(t)
		}
		return 0
	}
	return 0
}

func (s *scaleImpl) Bandwidth() float64 {
	if s.typ == ScaleTypeBand {
		return s.band.Bandwidth()
	}
	if s.typ == ScaleTypePoint {
		return s.point.Bandwidth()
	}
	return 0
}

func (s *scaleImpl) Step() float64 {
	if s.typ == ScaleTypeBand {
		return s.band.Step()
	}
	if s.typ == ScaleTypePoint {
		return s.point.Step()
	}
	return 0
}

func (s *scaleImpl) Round() bool {
	if s.typ == ScaleTypeBand {
		return s.band.Round()
	}
	if s.typ == ScaleTypePoint {
		return s.point.Round()
	}
	return false
}

func (s *scaleImpl) Ticks(count int) []any {
	switch s.typ {
	case ScaleTypeLinear:
		ts := s.linear.Ticks(count)
		out := make([]any, len(ts))
		for i, t := range ts {
			out[i] = t
		}
		return out
	case ScaleTypeLog:
		ts := s.log.Ticks(count)
		out := make([]any, len(ts))
		for i, t := range ts {
			out[i] = t
		}
		return out
	case ScaleTypeSymlog:
		ts := s.symlog.Ticks(count)
		out := make([]any, len(ts))
		for i, t := range ts {
			out[i] = t
		}
		return out
	case ScaleTypeTime:
		ts := s.time.Ticks(count)
		out := make([]any, len(ts))
		for i, t := range ts {
			out[i] = t
		}
		return out
	case ScaleTypeBand, ScaleTypePoint:
		// Discrete scales: domain values are the ticks.
		if s.typ == ScaleTypeBand {
			out := make([]any, len(s.band.Domain()))
			for i, d := range s.band.Domain() {
				out[i] = d
			}
			return out
		}
		out := make([]any, len(s.point.Domain()))
		for i, d := range s.point.Domain() {
			out[i] = d
		}
		return out
	}
	return nil
}

func (s *scaleImpl) Domain() []any {
	switch s.typ {
	case ScaleTypeLinear:
		d := s.linear.Domain()
		return []any{d[0], d[1]}
	case ScaleTypeLog:
		d := s.log.Domain()
		return []any{d[0], d[1]}
	case ScaleTypeSymlog:
		d := s.symlog.Domain()
		return []any{d[0], d[1]}
	case ScaleTypeTime:
		d := s.time.Domain()
		return []any{d[0], d[1]}
	case ScaleTypeBand:
		return toAnySlice(s.band.Domain())
	case ScaleTypePoint:
		return toAnySlice(s.point.Domain())
	}
	return nil
}

func toAnySlice(s []string) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0
		}
		return f
	case time.Time:
		return float64(x.UnixMilli())
	case nil:
		return 0
	}
	return 0
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func toTime(v any) (time.Time, bool) {
	switch x := v.(type) {
	case time.Time:
		return x, true
	case *time.Time:
		if x == nil {
			return time.Time{}, false
		}
		return *x, true
	case string:
		if t, err := time.Parse(time.RFC3339, x); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// --- ComputeScale ----------------------------------------------------------

// LinearScaleDefaults mirrors @nivo/scales linearScaleDefaults.
var LinearScaleDefaults = ScaleLinearSpec{
	Min: FloatVal(0), Max: AutoFloat(), Stacked: false, Reverse: false,
	Clamp: false, Nice: true, Round: false,
}

// BandScaleDefaults mirrors @nivo/scales bandScaleDefaults.
var BandScaleDefaults = ScaleBandSpec{Round: false}

// LogScaleDefaults mirrors @nivo/scales logScaleDefaults.
var LogScaleDefaults = ScaleLogSpec{
	Base: 10, Min: AutoFloat(), Max: AutoFloat(), Round: false, Reverse: false, Nice: true,
}

// SymlogScaleDefaults mirrors @nivo/scales symlogScaleDefaults.
var SymlogScaleDefaults = ScaleSymlogSpec{
	Constant: 1, Min: AutoFloat(), Max: AutoFloat(), Round: false, Reverse: false, Nice: true,
}

// TimeScaleDefaults mirrors @nivo/scales timeScaleDefaults.
var TimeScaleDefaults = ScaleTimeSpec{
	Format: "native", Precision: TimePrecisionMillisecond, Min: "auto", Max: "auto", UseUTC: true, Nice: false,
}

// PointScaleDefaults is nivo's implicit point-scale defaults.
var PointScaleDefaults = ScalePointSpec{}

// ComputeScale builds a Scale from a spec + ComputedSerieAxis + size + axis.
// Mirrors @nivo/scales computeScale.
func ComputeScale(spec ScaleSpec, data ComputedSerieAxis, size float64, axis ScaleAxis) Scale {
	switch s := spec.(type) {
	case ScaleLinearSpec:
		return createLinearScale(s, data, size, axis)
	case ScalePointSpec:
		return createPointScale(s, data, size)
	case ScaleBandSpec:
		return createBandScale(s, data, size, axis)
	case ScaleTimeSpec:
		return createTimeScale(s, data, size)
	case ScaleLogSpec:
		return createLogScale(s, data, size, axis)
	case ScaleSymlogSpec:
		return createSymlogScale(s, data, size, axis)
	}
	return nil
}

func axisRange(size float64, axis ScaleAxis) (float64, float64) {
	if axis == ScaleAxisX {
		return 0, size
	}
	return size, 0
}

func niceBool(v any) (bool, int) {
	switch x := v.(type) {
	case bool:
		return x, 0
	case int:
		return true, x
	case float64:
		return true, int(x)
	}
	return false, 0
}

func createLinearScale(spec ScaleLinearSpec, data ComputedSerieAxis, size float64, axis ScaleAxis) Scale {
	min := data.toFloat(data.Min)
	if !spec.Min.Auto {
		min = spec.Min.Value
	} else if spec.Stacked && data.MinStacked != nil {
		min = *data.MinStacked
	}
	max := data.toFloat(data.Max)
	if !spec.Max.Auto {
		max = spec.Max.Value
	} else if spec.Stacked && data.MaxStacked != nil {
		max = *data.MaxStacked
	}
	s := d3scale.NewLinear()
	r0, r1 := axisRange(size, axis)
	s.SetRange(r0, r1)
	if spec.Round {
		s.SetRangeRound(r0, r1)
	}
	if spec.Reverse {
		s.SetDomain(max, min)
	} else {
		s.SetDomain(min, max)
	}
	if spec.Clamp {
		s.SetClamp(true)
	}
	if ok, n := niceBool(spec.Nice); ok {
		if n > 0 {
			s.Nice(n)
		} else {
			s.Nice(10)
		}
	}
	return &scaleImpl{typ: ScaleTypeLinear, linear: s, stacked: spec.Stacked}
}

func createLogScale(spec ScaleLogSpec, data ComputedSerieAxis, size float64, axis ScaleAxis) Scale {
	min := data.toFloat(data.Min)
	if !spec.Min.Auto {
		min = spec.Min.Value
	}
	max := data.toFloat(data.Max)
	if !spec.Max.Auto {
		max = spec.Max.Value
	}
	base := spec.Base
	if base == 0 {
		base = LogScaleDefaults.Base
	}
	s := d3scale.NewLog().SetBase(base)
	r0, r1 := axisRange(size, axis)
	if spec.Round {
		s.SetRangeRound(r0, r1)
	} else {
		s.SetRange(r0, r1)
	}
	if spec.Reverse {
		s.SetDomain(max, min)
	} else {
		s.SetDomain(min, max)
	}
	if ok, _ := niceBool(spec.Nice); ok {
		s.Nice(10)
	}
	return &scaleImpl{typ: ScaleTypeLog, log: s}
}

func createSymlogScale(spec ScaleSymlogSpec, data ComputedSerieAxis, size float64, axis ScaleAxis) Scale {
	min := data.toFloat(data.Min)
	if !spec.Min.Auto {
		min = spec.Min.Value
	}
	max := data.toFloat(data.Max)
	if !spec.Max.Auto {
		max = spec.Max.Value
	}
	constant := spec.Constant
	if constant == 0 {
		constant = SymlogScaleDefaults.Constant
	}
	s := d3scale.NewSymlog().SetConstant(constant)
	r0, r1 := axisRange(size, axis)
	if spec.Round {
		s.SetRangeRound(r0, r1)
	} else {
		s.SetRange(r0, r1)
	}
	if spec.Reverse {
		s.SetDomain(max, min)
	} else {
		s.SetDomain(min, max)
	}
	if ok, n := niceBool(spec.Nice); ok {
		if n > 0 {
			s.Nice(n)
		} else {
			s.Nice(10)
		}
	}
	return &scaleImpl{typ: ScaleTypeSymlog, symlog: s}
}

func createBandScale(spec ScaleBandSpec, data ComputedSerieAxis, size float64, axis ScaleAxis) Scale {
	domain := data.toStringSlice()
	s := d3scale.NewBand().SetDomain(domain)
	r0, r1 := axisRange(size, axis)
	s.SetRange(r0, r1)
	s.SetRound(spec.Round)
	return &scaleImpl{typ: ScaleTypeBand, band: s}
}

func createPointScale(spec ScalePointSpec, data ComputedSerieAxis, size float64) Scale {
	domain := data.toStringSlice()
	s := d3scale.NewPoint().SetDomain(domain)
	s.SetRange(0, size)
	return &scaleImpl{typ: ScaleTypePoint, point: s}
}

func createTimeScale(spec ScaleTimeSpec, data ComputedSerieAxis, size float64) Scale {
	normalize := CreateDateNormalizer(spec.Format, spec.Precision, spec.UseUTC)
	var minT, maxT time.Time
	if spec.Min == "auto" || spec.Min == nil {
		if t, ok := toTime(data.Min); ok {
			minT = normalize(t)
		}
	} else if t, ok := toTime(spec.Min); ok {
		minT = normalize(t)
	} else if s, ok := spec.Min.(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			minT = normalize(t)
		}
	}
	if spec.Max == "auto" || spec.Max == nil {
		if t, ok := toTime(data.Max); ok {
			maxT = normalize(t)
		}
	} else if t, ok := toTime(spec.Max); ok {
		maxT = normalize(t)
	} else if s, ok := spec.Max.(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			maxT = normalize(t)
		}
	}
	useUTC := spec.UseUTC
	if useUTC {
		// d3 scaleUtc not in port; use local time scale with UTC-normalized
		// instants. For v1 charts this is acceptable.
	}
	s := d3scale.NewTime().SetUseUTC(useUTC)
	s.SetRange(0, size)
	if !minT.IsZero() && !maxT.IsZero() {
		s.SetDomain(minT, maxT)
	}
	if ok, _ := niceBool(spec.Nice); ok {
		s.Nice(10)
	}
	return &scaleImpl{typ: ScaleTypeTime, time: s, useUTC: useUTC}
}

// toFloat coerces a ComputedSerieAxis value to float64.
func (a ComputedSerieAxis) toFloat(v any) float64 { return toFloat(v) }

func (a ComputedSerieAxis) toStringSlice() []string {
	out := make([]string, 0, len(a.All))
	for _, v := range a.All {
		out = append(out, toString(v))
	}
	return out
}

// GetOtherAxis returns the opposite axis.
func GetOtherAxis(axis ScaleAxis) ScaleAxis {
	if axis == ScaleAxisX {
		return ScaleAxisY
	}
	return ScaleAxisX
}

// --- GetScaleTicks / CenterScale -------------------------------------------

// GetScaleTicks mirrors @nivo/scales getScaleTicks: returns explicit tick
// values if spec.HasValues, a time-interval range if spec.HasInterval and the
// scale is a time scale, the scale's default ticks if HasCount, else the
// scale's default ticks.
func GetScaleTicks(scale Scale, spec TicksSpec) []any {
	if spec.HasValues {
		return spec.Values
	}
	si, ok := scale.(*scaleImpl)
	if !ok {
		return nil
	}
	if spec.HasInterval {
		// v1: time-interval parsing is simplified — return default ticks.
		return si.Ticks(10)
	}
	count := 10
	if spec.HasCount {
		count = spec.Count
	}
	return si.Ticks(count)
}

// CenterScale mirrors @nivo/scales centerScale: for band/point scales, returns
// a func that offsets the scale output by bandwidth/2 so ticks are centered.
func CenterScale(scale Scale) func(any) float64 {
	si, ok := scale.(*scaleImpl)
	if !ok {
		return func(v any) float64 { return scale.Call(v) }
	}
	if si.typ != ScaleTypeBand && si.typ != ScaleTypePoint {
		return func(v any) float64 { return scale.Call(v) }
	}
	bw := si.Bandwidth()
	if bw == 0 {
		return func(v any) float64 { return scale.Call(v) }
	}
	offset := bw / 2
	if si.Round() {
		offset = math.Round(offset)
	}
	return func(v any) float64 { return scale.Call(v) + offset }
}
