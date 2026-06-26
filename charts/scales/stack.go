package scales

import (
	"math"
	"sort"
	"strconv"
	"time"
)

// SerieDatum is a single x/y data point with optional stacked values. Mirrors
// @nivo/scales SerieDatum.
type SerieDatum struct {
	X        any // number | string | time.Time | nil
	XStacked *float64
	Y        any // number | string | time.Time | nil
	YStacked *float64
}

// Serie is a chart series: extra metadata + a slice of data points. The data
// field carries the original points; nesting (data.data) is nivo-internal and
// not needed for the Go port.
type Serie struct {
	Data []SerieDatum
	// Extra carries series metadata (id, color, etc.).
	Extra any
}

// ComputedSerie is a Serie with computed x/y positions for each datum.
type ComputedSerie struct {
	Serie
	Data []ComputedDatum
}

// ComputedDatum is a SerieDatum with computed x/y pixel positions.
type ComputedDatum struct {
	SerieDatum
	Position DatumPosition
}

// DatumPosition is the computed x/y pixel position of a datum.
type DatumPosition struct {
	X *float64
	Y *float64
}

// ComputedXYScales is the result of ComputeXYScalesForSeries: the per-axis
// computed data, the resolved x/y scales, and the series with positions.
type ComputedXYScales struct {
	X      ComputedSerieAxis
	Y      ComputedSerieAxis
	XScale Scale
	YScale Scale
	Series []ComputedSerie
}

// ComputeXYScalesForSeries mirrors @nivo/scales computeXYScalesForSeries:
// computes the x/y axes (all/min/max), optionally stacks values, builds the
// x/y scales, and assigns pixel positions to each datum.
func ComputeXYScalesForSeries(
	series []Serie,
	xScaleSpec ScaleSpec,
	yScaleSpec ScaleSpec,
	width, height float64,
) ComputedXYScales {
	xy := generateSeriesXY(series, xScaleSpec, yScaleSpec)

	// Stack x values if x scale is stacked.
	if ls, ok := xScaleSpec.(ScaleLinearSpec); ok && ls.Stacked {
		stackAxis(ScaleAxisX, &xy.X, series)
	}
	// Stack y values if y scale is stacked.
	if ls, ok := yScaleSpec.(ScaleLinearSpec); ok && ls.Stacked {
		stackAxis(ScaleAxisY, &xy.Y, series)
	}

	xScale := ComputeScale(xScaleSpec, xy.X, width, ScaleAxisX)
	yScale := ComputeScale(yScaleSpec, xy.Y, height, ScaleAxisY)

	computed := make([]ComputedSerie, len(series))
	for i, s := range series {
		data := make([]ComputedDatum, len(s.Data))
		for j, d := range s.Data {
			data[j] = ComputedDatum{
				SerieDatum: d,
				Position: DatumPosition{
					X: getDatumAxisPosition(d, ScaleAxisX, xScale),
					Y: getDatumAxisPosition(d, ScaleAxisY, yScale),
				},
			}
		}
		computed[i] = ComputedSerie{Serie: s, Data: data}
	}

	return ComputedXYScales{
		X: xy.X, Y: xy.Y, XScale: xScale, YScale: yScale, Series: computed,
	}
}

// seriesXY holds the computed per-axis data before stacking.
type seriesXY struct {
	X ComputedSerieAxis
	Y ComputedSerieAxis
}

func generateSeriesXY(series []Serie, xSpec ScaleSpec, ySpec ScaleSpec) seriesXY {
	return seriesXY{
		X: generateSeriesAxis(series, ScaleAxisX, xSpec),
		Y: generateSeriesAxis(series, ScaleAxisY, ySpec),
	}
}

// generateSeriesAxis mirrors @nivo/scales generateSeriesAxis: normalizes the
// axis values per scale type (linear → float, time → Date), collects unique
// sorted values, and computes min/max.
func generateSeriesAxis(series []Serie, axis ScaleAxis, spec ScaleSpec) ComputedSerieAxis {
	values := []any{}
	for _, s := range series {
		for _, d := range s.Data {
			var v any
			if axis == ScaleAxisX {
				v = d.X
			} else {
				v = d.Y
			}
			if v == nil {
				continue
			}
			// Normalize linear values to float.
			if _, ok := spec.(ScaleLinearSpec); ok {
				if f, ok := toFloatOK(v); ok {
					v = f
				}
			}
			values = append(values, v)
		}
	}

	switch spec.(type) {
	case ScaleLinearSpec:
		floats := uniqueFloats(values)
		sort.Float64s(floats)
		all := make([]any, len(floats))
		for i, f := range floats {
			all[i] = f
		}
		if len(all) == 0 {
			return ComputedSerieAxis{All: all, Min: 0.0, Max: 0.0}
		}
		return ComputedSerieAxis{All: all, Min: all[0], Max: all[len(all)-1]}
	case ScaleTimeSpec:
		times := uniqueTimes(values)
		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
		all := make([]any, len(times))
		for i, t := range times {
			all[i] = t
		}
		if len(all) == 0 {
			return ComputedSerieAxis{All: all, Min: time.Time{}, Max: time.Time{}}
		}
		return ComputedSerieAxis{All: all, Min: all[0], Max: all[len(all)-1]}
	default: // point / band
		all := uniqueAny(values)
		if len(all) == 0 {
			return ComputedSerieAxis{All: all}
		}
		return ComputedSerieAxis{All: all, Min: all[0], Max: all[len(all)-1]}
	}
}

// stackAxis mirrors @nivo/scales stackAxis: accumulates stacked values across
// series for the given axis, populating xy[axis].MinStacked/MaxStacked and
// each datum's *Stacked field.
func stackAxis(axis ScaleAxis, xy *ComputedSerieAxis, series []Serie) {
	other := GetOtherAxis(axis)
	all := []float64{}
	for _, otherVal := range xyAllByOtherAxis(series, other) {
		stack := []float64{}
		for i := range series {
			s := &series[i]
			var datum *SerieDatum
			for j := range s.Data {
				d := &s.Data[j]
				if valuesEqual(axisValue(d, other), otherVal) {
					datum = d
					break
				}
			}
			var stackValue *float64
			if datum != nil {
				v := toFloat(axisValue(datum, axis))
				if v != 0 || !isNil(axisValue(datum, axis)) {
					var head *float64
					for k := len(stack) - 1; k >= 0; k-- {
						if stack[k] != 0 || true {
							head = &stack[k]
							break
						}
					}
					if head == nil {
						sv := v
						stackValue = &sv
					} else {
						sv := *head + v
						stackValue = &sv
					}
				}
				if axis == ScaleAxisX {
					datum.XStacked = stackValue
				} else {
					datum.YStacked = stackValue
				}
			}
			var pushed float64
			if stackValue != nil {
				pushed = *stackValue
			}
			stack = append(stack, pushed)
			if stackValue != nil {
				all = append(all, *stackValue)
			}
		}
	}
	if len(all) > 0 {
		min := all[0]
		max := all[0]
		for _, v := range all {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
		xy.MinStacked = &min
		xy.MaxStacked = &max
	} else {
		z := 0.0
		xy.MinStacked = &z
		xy.MaxStacked = &z
	}
}

func xyAllByOtherAxis(series []Serie, other ScaleAxis) []any {
	seen := map[string]bool{}
	out := []any{}
	for _, s := range series {
		for _, d := range s.Data {
			v := axisValue(&d, other)
			if v == nil {
				continue
			}
			key := fmtKey(v)
			if !seen[key] {
				seen[key] = true
				out = append(out, v)
			}
		}
	}
	return out
}

func axisValue(d *SerieDatum, axis ScaleAxis) any {
	if axis == ScaleAxisX {
		return d.X
	}
	return d.Y
}

func valuesEqual(a, b any) bool {
	return fmtKey(a) == fmtKey(b)
}

func isNil(v any) bool { return v == nil }

func fmtKey(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return "s:" + x
	case float64:
		return "f:" + strconv.FormatFloat(x, 'g', -1, 64)
	case int:
		return "f:" + strconv.Itoa(x)
	case time.Time:
		return "t:" + x.Format(time.RFC3339Nano)
	}
	return "x:" + toFloatAsString(v)
}

func toFloatAsString(v any) string {
	return strconv.FormatFloat(toFloat(v), 'g', -1, 64)
}

// getDatumAxisPosition returns the pixel position of a datum on an axis,
// honoring stacked scales. Mirrors @nivo/scales getDatumAxisPosition.
func getDatumAxisPosition(d SerieDatum, axis ScaleAxis, scale Scale) *float64 {
	si, ok := scale.(*scaleImpl)
	if ok && si.stacked {
		var stacked *float64
		if axis == ScaleAxisX {
			stacked = d.XStacked
		} else {
			stacked = d.YStacked
		}
		if stacked == nil {
			return nil
		}
		p := scale.Call(*stacked)
		return &p
	}
	var v any
	if axis == ScaleAxisX {
		v = d.X
	} else {
		v = d.Y
	}
	if v == nil {
		return nil
	}
	p := scale.Call(v)
	return &p
}

// --- unique helpers --------------------------------------------------------

func uniqueFloats(values []any) []float64 {
	seen := map[float64]bool{}
	out := []float64{}
	for _, v := range values {
		f, ok := toFloatOK(v)
		if !ok {
			continue
		}
		if !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	return out
}

func uniqueTimes(values []any) []time.Time {
	seen := map[int64]bool{}
	out := []time.Time{}
	for _, v := range values {
		t, ok := toTime(v)
		if !ok {
			continue
		}
		ns := t.UnixNano()
		if !seen[ns] {
			seen[ns] = true
			out = append(out, t)
		}
	}
	return out
}

func uniqueAny(values []any) []any {
	seen := map[string]bool{}
	out := []any{}
	for _, v := range values {
		k := fmtKey(v)
		if !seen[k] {
			seen[k] = true
			out = append(out, v)
		}
	}
	return out
}

func toFloatOK(v any) (float64, bool) {
	switch x := v.(type) {
	case nil:
		return 0, false
	case float64:
		if math.IsNaN(x) {
			return 0, false
		}
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
