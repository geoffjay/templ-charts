package compute

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/scales"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// stackedSeries is a local alias of d3shape.StackSeries used by the flatten
// helper. Carries the key + the [lo,hi] points.
type stackedSeries = d3shape.StackSeries

// GenerateStackedBarsResult is the output of GenerateStackedBars.
type GenerateStackedBarsResult struct {
	XScale scales.Scale
	YScale scales.Scale
	Bars   []ComputedBarDatum
}

// StackedParams is the input to GenerateStackedBars. Mirrors the destructured
// props of @nivo/bar generateStackedBars.
type StackedParams struct {
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

// GenerateStackedBars mirrors @nivo/bar generateStackedBars. Stacks the keys
// via d3-shape's stack generator with the diverging offset, then builds the
// per-bar geometry for vertical/horizontal stacked layout.
func GenerateStackedBars(p StackedParams) GenerateStackedBarsResult {
	keys := filterKeys(p.Keys, p.HiddenIDs)
	data := NormalizeData(p.Data, keys)

	// d3-shape stack: keys, offset=diverging.
	stack := d3shape.NewStack[map[string]any]().
		Keys(keys).
		Offset(d3shape.StackOffsetDiverging).
		Value(func(d map[string]any, key string, _ int, _ []map[string]any) float64 {
			_, v := CoerceValue(d[key])
			return v
		})
	stackedData := stack.Call(data)

	var axis, otherAxis scales.ScaleAxis
	var size float64
	if p.Layout == "vertical" {
		axis, otherAxis, size = scales.ScaleAxisY, scales.ScaleAxisX, p.Width
	} else {
		axis, otherAxis, size = scales.ScaleAxisX, scales.ScaleAxisY, p.Height
	}

	indexScale := GetIndexScale(p.Data, p.GetIndex, p.Padding, p.IndexScale, size, otherAxis)
	ixScale := indexScale.(IndexScaleWithDomain)

	values := filterZerosIfLog(flattenStacked(stackedData), string(p.ValueScale.ScaleType()))
	min := mathMin(values)
	max := mathMax(values)

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

	innerPadding := p.InnerPadding
	bandwidth := ixScale.Bandwidth()
	reverse := p.ValueScale.Reverse

	var bars []ComputedBarDatum
	if bandwidth > 0 {
		if p.Layout == "vertical" {
			bars = generateVerticalStackedBars(p, stackedData, ixScale, xScale, yScale, bandwidth, innerPadding, reverse)
		} else {
			bars = generateHorizontalStackedBars(p, stackedData, ixScale, xScale, yScale, bandwidth, innerPadding, reverse)
		}
	}

	return GenerateStackedBarsResult{XScale: xScale, YScale: yScale, Bars: bars}
}

func generateVerticalStackedBars(
	p StackedParams,
	stackedData []stackedSeries,
	ixScale IndexScaleWithDomain,
	xScale, yScale scales.Scale,
	bandwidth, innerPadding float64,
	reverse bool,
) []ComputedBarDatum {
	domain := ixScale.Domain()
	bars := []ComputedBarDatum{}
	for _, si := range stackedData {
		for i, index := range domain {
			d := si.Stats[i]
			x := xScale.Call(index)
			y := yScale.Call(d.Hi)
			if reverse {
				y = yScale.Call(d.Lo)
			}
			y += innerPadding * 0.5
			barHeight := (yScale.Call(d.Lo) - y)
			if reverse {
				barHeight = yScale.Call(d.Hi) - y
			}
			barHeight -= innerPadding
			rawValue, value := CoerceValue(d.Data.(map[string]any)[si.Key])
			bd := ComputedDatum{
				ID:             si.Key,
				Value:          rawValue,
				FormattedValue: p.FormatValue(value),
				Hidden:         false,
				Index:          i,
				IndexValue:     index,
				Data:           FilterNullValues(d.Data.(map[string]any)),
			}
			bars = append(bars, ComputedBarDatum{
				Key:    si.Key + "." + index,
				Index:  len(bars),
				Data:   bd,
				X:      x,
				Y:      y,
				AbsX:   p.Margin.Left + x,
				AbsY:   p.Margin.Top + y,
				Width:  bandwidth,
				Height: barHeight,
				Color:  p.GetColor(bd),
				Label:  p.GetTooltipLabel(bd),
			})
		}
	}
	return bars
}

func generateHorizontalStackedBars(
	p StackedParams,
	stackedData []stackedSeries,
	ixScale IndexScaleWithDomain,
	xScale, yScale scales.Scale,
	bandwidth, innerPadding float64,
	reverse bool,
) []ComputedBarDatum {
	domain := ixScale.Domain()
	bars := []ComputedBarDatum{}
	for _, si := range stackedData {
		for i, index := range domain {
			d := si.Stats[i]
			y := yScale.Call(index)
			x := xScale.Call(d.Lo)
			if reverse {
				x = xScale.Call(d.Hi)
			}
			x += innerPadding * 0.5
			barWidth := xScale.Call(d.Hi) - x
			if reverse {
				barWidth = xScale.Call(d.Lo) - x
			}
			barWidth -= innerPadding
			rawValue, value := CoerceValue(d.Data.(map[string]any)[si.Key])
			bd := ComputedDatum{
				ID:             si.Key,
				Value:          rawValue,
				FormattedValue: p.FormatValue(value),
				Hidden:         false,
				Index:          i,
				IndexValue:     index,
				Data:           FilterNullValues(d.Data.(map[string]any)),
			}
			bars = append(bars, ComputedBarDatum{
				Key:    si.Key + "." + index,
				Index:  len(bars),
				Data:   bd,
				X:      x,
				Y:      y,
				AbsX:   p.Margin.Left + x,
				AbsY:   p.Margin.Top + y,
				Width:  barWidth,
				Height: bandwidth,
				Color:  p.GetColor(bd),
				Label:  p.GetTooltipLabel(bd),
			})
		}
	}
	return bars
}

// unused-guard so math import isn't dropped if future refactors remove calls.
var _ = math.Pi
