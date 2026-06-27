package compute

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/scales"
	d3scale "github.com/geoffjay/templ-charts/internal/d3/scale"
)

// GenerateGroupedBarsResult is the output of GenerateGroupedBars.
type GenerateGroupedBarsResult struct {
	XScale scales.Scale
	YScale scales.Scale
	Bars   []ComputedBarDatum
}

// GroupedParams is the input to GenerateGroupedBars. Mirrors the destructured
// props of @nivo/bar generateGroupedBars.
type GroupedParams struct {
	Data            []map[string]any
	Keys            []string
	Layout          string
	Width           float64
	Height          float64
	Padding         float64
	InnerPadding    float64
	ValueScale      scales.ScaleLinearSpec
	IndexScale      scales.ScaleBandSpec
	GetIndex        func(map[string]any) string
	GetColor        func(ComputedDatum) string
	GetTooltipLabel func(ComputedDatum) string
	FormatValue     func(float64) string
	Margin          MarginLike
	HiddenIDs       []string
}

// GenerateGroupedBars mirrors @nivo/bar generateGroupedBars. Produces the x/y
// scales and the per-bar geometry for a grouped bar chart (vertical or
// horizontal).
func GenerateGroupedBars(p GroupedParams) GenerateGroupedBarsResult {
	keys := filterKeys(p.Keys, p.HiddenIDs)
	data := NormalizeData(p.Data, keys)

	var axis, otherAxis scales.ScaleAxis
	var size float64
	if p.Layout == "vertical" {
		axis, otherAxis, size = scales.ScaleAxisY, scales.ScaleAxisX, p.Width
	} else {
		axis, otherAxis, size = scales.ScaleAxisX, scales.ScaleAxisY, p.Height
	}

	indexScale := GetIndexScale(p.Data, p.GetIndex, p.Padding, p.IndexScale, size, otherAxis)
	ixScale := indexScale.(IndexScaleWithDomain)

	// Gather all values, compute min/max, clamp min to 0 when valueScale.min=="auto".
	clampMin := func(v float64) float64 { return v }
	if p.ValueScale.Min.Auto {
		clampMin = clampToZero
	}
	values := []float64{}
	for _, entry := range data {
		for _, k := range keys {
			if v, ok := toFloatOK(entry[k]); ok {
				if v != 0 {
					values = append(values, v)
				}
			}
		}
	}
	min := clampMin(mathMin(values))
	max := zeroIfNotFinite(mathMax(values))

	valueAxisSize := p.Height
	if axis == scales.ScaleAxisX {
		valueAxisSize = p.Width
	}
	valueScale := scales.ComputeScale(p.ValueScale, scales.ComputedSerieAxis{
		All: valuesToAny(values), Min: min, Max: max,
	}, valueAxisSize, axis)

	var xScale, yScale scales.Scale
	if p.Layout == "vertical" {
		xScale, yScale = indexScale, valueScale
	} else {
		xScale, yScale = valueScale, indexScale
	}

	bw := ixScale.Bandwidth()
	// nivo: bandwidth = (bandwidth - innerPadding*(keys.length-1)) / keys.length
	if len(keys) > 0 {
		bw = (bw - p.InnerPadding*float64(len(keys)-1)) / float64(len(keys))
	}

	var bars []ComputedBarDatum
	if bw > 0 {
		reverse := p.ValueScale.Reverse
		if p.Layout == "vertical" {
			bars = generateVerticalGroupedBars(p, data, keys, xScale, yScale, ixScale, bw, reverse, valueScale.Call(0))
		} else {
			bars = generateHorizontalGroupedBars(p, data, keys, xScale, yScale, ixScale, bw, reverse, valueScale.Call(0))
		}
	}

	return GenerateGroupedBarsResult{XScale: xScale, YScale: yScale, Bars: bars}
}

func generateVerticalGroupedBars(
	p GroupedParams,
	data []map[string]any,
	keys []string,
	xScale, yScale scales.Scale,
	ixScale IndexScaleWithDomain,
	barWidth float64,
	reverse bool,
	yRef float64,
) []ComputedBarDatum {
	compare := func(a, b float64) bool { return a > b }
	if reverse {
		compare = func(a, b float64) bool { return a < b }
	}
	getY := func(d float64) float64 {
		if compare(d, 0) {
			return yScale.Call(d)
		}
		return yRef
	}
	getHeight := func(d, y float64) float64 {
		if compare(d, 0) {
			return yRef - y
		}
		return yScale.Call(d) - yRef
	}
	domain := ixScale.Domain()
	bars := []ComputedBarDatum{}
	for i, key := range keys {
		for index := 0; index < len(domain); index++ {
			rawValue, value := CoerceValue(data[index][key])
			indexValue := domain[index]
			x := xScale.Call(indexValue) + barWidth*float64(i) + p.InnerPadding*float64(i)
			y := getY(value)
			barHeight := getHeight(value, y)
			bd := ComputedDatum{
				ID:             key,
				Value:          rawValue,
				FormattedValue: p.FormatValue(value),
				Hidden:         false,
				Index:          index,
				IndexValue:     indexValue,
				Data:           FilterNullValues(data[index]),
			}
			bars = append(bars, ComputedBarDatum{
				Key:    key + "." + indexValue,
				Index:  len(bars),
				Data:   bd,
				X:      x,
				Y:      y,
				AbsX:   p.Margin.Left + x,
				AbsY:   p.Margin.Top + y,
				Width:  barWidth,
				Height: barHeight,
				Color:  p.GetColor(bd),
				Label:  p.GetTooltipLabel(bd),
			})
		}
	}
	return bars
}

func generateHorizontalGroupedBars(
	p GroupedParams,
	data []map[string]any,
	keys []string,
	xScale, yScale scales.Scale,
	ixScale IndexScaleWithDomain,
	barHeight float64,
	reverse bool,
	xRef float64,
) []ComputedBarDatum {
	compare := func(a, b float64) bool { return a > b }
	if reverse {
		compare = func(a, b float64) bool { return a < b }
	}
	getX := func(d float64) float64 {
		if compare(d, 0) {
			return xRef
		}
		return xScale.Call(d)
	}
	getWidth := func(d, x float64) float64 {
		if compare(d, 0) {
			return xScale.Call(d) - xRef
		}
		return xRef - x
	}
	domain := ixScale.Domain()
	bars := []ComputedBarDatum{}
	for i, key := range keys {
		for index := 0; index < len(domain); index++ {
			rawValue, value := CoerceValue(data[index][key])
			indexValue := domain[index]
			x := getX(value)
			y := yScale.Call(indexValue) + barHeight*float64(i) + p.InnerPadding*float64(i)
			barWidth := getWidth(value, x)
			bd := ComputedDatum{
				ID:             key,
				Value:          rawValue,
				FormattedValue: p.FormatValue(value),
				Hidden:         false,
				Index:          index,
				IndexValue:     indexValue,
				Data:           FilterNullValues(data[index]),
			}
			bars = append(bars, ComputedBarDatum{
				Key:    key + "." + indexValue,
				Index:  len(bars),
				Data:   bd,
				X:      x,
				Y:      y,
				AbsX:   p.Margin.Left + x,
				AbsY:   p.Margin.Top + y,
				Width:  barWidth,
				Height: barHeight,
				Color:  p.GetColor(bd),
				Label:  p.GetTooltipLabel(bd),
			})
		}
	}
	return bars
}

// --- helpers shared by grouped/stacked ---

func filterKeys(keys []string, hidden []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if !containsString(hidden, k) {
			out = append(out, k)
		}
	}
	return out
}

func containsString(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func clampToZero(v float64) float64 {
	if v > 0 {
		return 0
	}
	return v
}

func zeroIfNotFinite(v float64) float64 {
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0
	}
	return v
}

func mathMin(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func mathMax(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func valuesToAny(xs []float64) []any {
	out := make([]any, len(xs))
	for i, x := range xs {
		out[i] = x
	}
	return out
}

// filterZerosIfLog mirrors nivo's filterZerosIfLog (log scales can't handle 0).
func filterZerosIfLog(xs []float64, scaleType string) []float64 {
	if scaleType != "log" {
		return xs
	}
	out := xs[:0:0]
	for _, x := range xs {
		if x != 0 {
			out = append(out, x)
		}
	}
	return out
}

// logScaleBase returns the spec base or 10 if zero. Used by stacked.
func logScaleBase(spec scales.ScaleLogSpec) float64 {
	if spec.Base == 0 {
		return 10
	}
	return spec.Base
}

// flattenStacked mirrors nivo's flattenDeep over stackedData.
func flattenStacked(series []stackedSeries) []float64 {
	var out []float64
	for _, s := range series {
		for _, p := range s.Stats {
			out = append(out, p.Lo, p.Hi)
		}
	}
	return out
}

// Compile-time check that bandScale satisfies the interfaces we use.
var (
	_ IndexScaleWithDomain      = bandScale{}
	_ scales.ScaleWithBandwidth = bandScale{}
	_ scales.Scale              = bandScale{}
)

// reachLog unused import guard
var _ = d3scale.NewBand
