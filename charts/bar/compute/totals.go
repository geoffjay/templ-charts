package compute

import (
	"github.com/geoffjay/templ-charts/charts/scales"
)

// ComputeBarTotals mirrors @nivo/bar computeBarTotals. Produces per-index
// totals labels positioned above (vertical) or beside (horizontal) the bars.
// For stacked charts the label sits at the top of the positive stack; for
// grouped charts it sits at the greatest single bar value.
func ComputeBarTotals(
	bars []ComputedBarDatum,
	xScale, yScale scales.Scale,
	layout string,
	groupMode string,
	totalsOffset float64,
	formatValue func(float64) string,
) []BarTotalData {
	if len(bars) == 0 {
		return nil
	}
	totalsByIndex := map[string]float64{}
	barWidth := bars[0].Width
	barHeight := bars[0].Height

	if groupMode == "stacked" {
		totalsPositivesByIndex := map[string]float64{}
		for _, b := range bars {
			iv := b.Data.IndexValue
			v := toFloatOrZero(b.Data.Value)
			totalsByIndex[iv] += v
			if v > 0 {
				totalsPositivesByIndex[iv] += v
			}
		}
		out := []BarTotalData{}
		for iv, positive := range totalsPositivesByIndex {
			total := totalsByIndex[iv]
			var x, y, animOff float64
			if layout == "vertical" {
				x = xScale.Call(iv)
				y = yScale.Call(positive)
				animOff = yScale.Call(positive / 2)
				x += barWidth / 2
				y -= totalsOffset
			} else {
				x = xScale.Call(positive)
				y = yScale.Call(iv)
				animOff = xScale.Call(positive / 2)
				x += totalsOffset
				y += barHeight / 2
			}
			out = append(out, BarTotalData{
				Key: "total_" + iv, X: x, Y: y,
				Value: total, FormattedValue: formatValue(total),
				AnimationOffset: animOff,
			})
		}
		return out
	}

	// grouped
	greatestByIndex := map[string]float64{}
	numBarsByIndex := map[string]int{}
	for _, b := range bars {
		iv := b.Data.IndexValue
		v := toFloatOrZero(b.Data.Value)
		totalsByIndex[iv] += v
		if v > greatestByIndex[iv] {
			greatestByIndex[iv] = v
		}
		numBarsByIndex[iv]++
	}
	out := []BarTotalData{}
	for iv, greatest := range greatestByIndex {
		total := totalsByIndex[iv]
		numBars := numBarsByIndex[iv]
		var x, y, animOff float64
		if layout == "vertical" {
			x = xScale.Call(iv)
			y = yScale.Call(greatest)
			animOff = yScale.Call(greatest / 2)
			x += float64(numBars) * barWidth / 2
			y -= totalsOffset
		} else {
			x = xScale.Call(greatest)
			y = yScale.Call(iv)
			animOff = xScale.Call(greatest / 2)
			x += totalsOffset
			y += float64(numBars) * barHeight / 2
		}
		out = append(out, BarTotalData{
			Key: "total_" + iv, X: x, Y: y,
			Value: total, FormattedValue: formatValue(total),
			AnimationOffset: animOff,
		})
	}
	return out
}

func toFloatOrZero(v any) float64 {
	if v == nil {
		return 0
	}
	f, ok := toFloatOK(v)
	if !ok {
		return 0
	}
	return f
}
