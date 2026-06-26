// Package scales mirrors @nivo/scales: ScaleSpec types, ComputedSerieAxis,
// computeScale, computeXYScalesForSeries, stack helpers, getScaleTicks,
// centerScale, and time helpers. Backed by internal/d3/scale.
package scales

import (
	"time"
)

// ScaleAxis is "x" or "y".
type ScaleAxis string

const (
	ScaleAxisX ScaleAxis = "x"
	ScaleAxisY ScaleAxis = "y"
)

// ScaleValue is the union of value types a scale can accept (number, string,
// time.Time, or nil). Modeled as any for Go ergonomics.
type ScaleValue any

// ScaleType enumerates the supported scale types.
type ScaleType string

const (
	ScaleTypeLinear ScaleType = "linear"
	ScaleTypeLog    ScaleType = "log"
	ScaleTypeSymlog ScaleType = "symlog"
	ScaleTypePoint  ScaleType = "point"
	ScaleTypeBand   ScaleType = "band"
	ScaleTypeTime   ScaleType = "time"
)

// AutoOrFloat is a float64 value or the sentinel "auto". The Zero value
// (Auto=false, Value=0) means 0; Auto=true means "auto".
type AutoOrFloat struct {
	Auto  bool
	Value float64
}

// AutoFloat returns an "auto" sentinel.
func AutoFloat() AutoOrFloat { return AutoOrFloat{Auto: true} }

// FloatVal returns a concrete float sentinel.
func FloatVal(v float64) AutoOrFloat { return AutoOrFloat{Value: v} }

// ScaleSpec is the interface implemented by all scale spec types. Mirrors
// @nivo/scales ScaleSpec (a discriminated union in TS).
type ScaleSpec interface {
	ScaleType() ScaleType
}

// ScaleLinearSpec mirrors @nivo/scales ScaleLinearSpec.
type ScaleLinearSpec struct {
	Min     AutoOrFloat // default 0
	Max     AutoOrFloat // default "auto"
	Stacked bool
	Reverse bool
	Clamp   bool
	Nice    any // bool | int
	Round   bool
}

func (ScaleLinearSpec) ScaleType() ScaleType { return ScaleTypeLinear }

// ScaleLogSpec mirrors @nivo/scales ScaleLogSpec.
type ScaleLogSpec struct {
	Base    float64 // default 10
	Min     AutoOrFloat
	Max     AutoOrFloat
	Round   bool
	Reverse bool
	Nice    any
}

func (ScaleLogSpec) ScaleType() ScaleType { return ScaleTypeLog }

// ScaleSymlogSpec mirrors @nivo/scales ScaleSymlogSpec.
type ScaleSymlogSpec struct {
	Constant float64 // default 1
	Min      AutoOrFloat
	Max      AutoOrFloat
	Round    bool
	Reverse  bool
	Nice     any
}

func (ScaleSymlogSpec) ScaleType() ScaleType { return ScaleTypeSymlog }

// ScalePointSpec mirrors @nivo/scales ScalePointSpec.
type ScalePointSpec struct{}

func (ScalePointSpec) ScaleType() ScaleType { return ScaleTypePoint }

// ScaleBandSpec mirrors @nivo/scales ScaleBandSpec.
type ScaleBandSpec struct {
	Round bool
}

func (ScaleBandSpec) ScaleType() ScaleType { return ScaleTypeBand }

// ScaleTimeSpec mirrors @nivo/scales ScaleTimeSpec.
type ScaleTimeSpec struct {
	Format    string // "native" or a d3-time-format spec
	Precision TimePrecision
	Min       any // "auto" | time.Time | string
	Max       any // "auto" | time.Time | string
	UseUTC    bool
	Nice      any // bool | time.Time
}

func (ScaleTimeSpec) ScaleType() ScaleType { return ScaleTypeTime }

// TicksSpec mirrors @nivo/scales TicksSpec: either a tick count (int), a time
// interval string ("every 2 weeks"), or an explicit list of values.
type TicksSpec struct {
	Count    int
	Interval string
	Values   []any
	// HasCount / HasInterval / HasValues discriminate which form is set.
	HasCount    bool
	HasInterval bool
	HasValues   bool
}

// ComputedSerieAxis mirrors @nivo/scales ComputedSerieAxis: the all/min/max
// (and optional minStacked/maxStacked) for one axis.
type ComputedSerieAxis struct {
	All        []any
	Min        any
	Max        any
	MinStacked *float64
	MaxStacked *float64
}

// Scale is the interface implemented by resolved scales. Mirrors the subset
// of d3-scale operations nivo uses. Discrete scales (band/point) also expose
// Bandwidth/Step.
type Scale interface {
	Type() ScaleType
	Call(v any) float64
}

// ScaleWithBandwidth is implemented by band/point scales.
type ScaleWithBandwidth interface {
	Scale
	Bandwidth() float64
	Step() float64
	Round() bool
}

// TimePrecision mirrors @nivo/scales TIME_PRECISION.
type TimePrecision string

const (
	TimePrecisionMillisecond TimePrecision = "millisecond"
	TimePrecisionSecond      TimePrecision = "second"
	TimePrecisionMinute      TimePrecision = "minute"
	TimePrecisionHour        TimePrecision = "hour"
	TimePrecisionDay         TimePrecision = "day"
	TimePrecisionMonth       TimePrecision = "month"
	TimePrecisionYear        TimePrecision = "year"
)

// TimePrecisions is the ordered list of time precisions.
var TimePrecisions = []TimePrecision{
	TimePrecisionMillisecond, TimePrecisionSecond, TimePrecisionMinute,
	TimePrecisionHour, TimePrecisionDay, TimePrecisionMonth, TimePrecisionYear,
}

// precisionCutOffs applies the field-zeroing for a given precision.
var precisionCutOffs = []func(*time.Time){
	func(t *time.Time) { t.Add(0) }, // placeholder; actual mods below
}

func init() {
	precisionCutOffs = []func(*time.Time){
		func(t *time.Time) { *t = t.Add(-time.Duration(t.Nanosecond())) }, // ms → 0 ns? actually set ms=0 means zero sub-ms
	}
}

// CreateDateNormalizer mirrors @nivo/scales createDateNormalizer: given a
// format ("native" or spec), precision, and useUTC flag, returns a func that
// parses/rounds a value to a time.Time.
func CreateDateNormalizer(format string, precision TimePrecision, useUTC bool) func(any) time.Time {
	precisionFn := CreatePrecisionMethod(precision)
	return func(value any) time.Time {
		if v, ok := value.(time.Time); ok {
			return precisionFn(v)
		}
		if s, ok := value.(string); ok && format != "" && format != "native" {
			if t, err := parseTimeFormat(s, format, useUTC); err == nil {
				return precisionFn(t)
			}
		}
		if t, ok := value.(time.Time); ok {
			return precisionFn(t)
		}
		return time.Time{}
	}
}

// CreatePrecisionMethod returns a func that zeroes the fields below the given
// precision. Mirrors @nivo/scales createPrecisionMethod.
func CreatePrecisionMethod(precision TimePrecision) func(time.Time) time.Time {
	cutOffs := precisionCutOffsByType[precision]
	return func(t time.Time) time.Time {
		for _, cut := range cutOffs {
			cut(&t)
		}
		return t
	}
}

var precisionCutOffsByType = map[TimePrecision][]func(*time.Time){
	TimePrecisionMillisecond: {},
	TimePrecisionSecond:      {setMillisZero},
	TimePrecisionMinute:      {setMillisZero, setSecondsZero},
	TimePrecisionHour:        {setMillisZero, setSecondsZero, setMinutesZero},
	TimePrecisionDay:         {setMillisZero, setSecondsZero, setMinutesZero, setHoursZero},
	TimePrecisionMonth:       {setMillisZero, setSecondsZero, setMinutesZero, setHoursZero, setDayOne},
	TimePrecisionYear:        {setMillisZero, setSecondsZero, setMinutesZero, setHoursZero, setDayOne, setMonthZero},
}

func setMillisZero(t *time.Time)  { *t = time.Unix(t.Unix(), 0).In(t.Location()) }
func setSecondsZero(t *time.Time) { *t = t.Truncate(time.Minute) }
func setMinutesZero(t *time.Time) { *t = t.Truncate(time.Hour) }
func setHoursZero(t *time.Time) {
	*t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func setDayOne(t *time.Time)    { *t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()) }
func setMonthZero(t *time.Time) { *t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()) }

// parseTimeFormat parses s using a d3-time-format spec. Delegates to
// internal/d3/timeformat.
func parseTimeFormat(s, spec string, useUTC bool) (time.Time, error) {
	// d3-time-format parsing is not implemented in the port (only formatting).
	// For v1 we support RFC3339 and common Go reference layouts as a fallback.
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, errTimeParse
}

var errTimeParse = &parseError{msg: "time parse error"}

type parseError struct{ msg string }

func (e *parseError) Error() string { return e.msg }
