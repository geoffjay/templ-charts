package polarbar

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UsePolarBar mirrors @nivo/polar-bar's usePolarBar: it derives the center,
// inner/outer radii, the angle (band over indices) + radius (linear over
// stacked totals) scales, and the stacked arcs (one per index/key). Colors are
// assigned per key. props.Width/Height are the inner dimensions.
func UsePolarBar(props PolarBarProps) PolarBarResult {
	center := [2]float64{props.Width / 2, props.Height / 2}
	outerRadius := math.Min(center[0], center[1])
	innerRadius := outerRadius * math.Min(props.InnerRadius, 1)

	startAngle, endAngle := clampArc(props.StartAngle, props.EndAngle, 360)

	keys := props.Keys
	if len(keys) == 0 {
		keys = []string{"value"}
	}

	indices := make([]string, 0, len(props.Data))
	for _, d := range props.Data {
		indices = append(indices, d.Index)
	}

	// Max stacked total across indices (radius-scale domain max).
	maxValue := 0.0
	for _, d := range props.Data {
		total := 0.0
		for _, k := range keys {
			total += d.Values[k]
		}
		if total > maxValue {
			maxValue = total
		}
	}
	if props.MaxValue != nil {
		maxValue = *props.MaxValue
	}
	if maxValue == 0 {
		maxValue = 1
	}

	angleScale := scales.NewBandScaleWithRange(indices, startAngle, endAngle, props.Padding, false)
	radiusScale := scales.NewLinearScaleWithRange(0, maxValue, innerRadius, outerRadius)
	bandwidth := 0.0
	if bw, ok := angleScale.(scales.ScaleWithBandwidth); ok {
		bandwidth = bw.Bandwidth()
	}

	arcGen := arcs.CreateArcGenerator(props.CornerRadius, arcs.DegToRad(props.PadAngle))
	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(k string) string { return k })
	format := valueFormatter(props.ValueFormat)

	computed := make([]ComputedArc, 0, len(props.Data)*len(keys))
	for _, d := range props.Data {
		bandStart := angleScale.Call(d.Index)
		lo := 0.0
		for _, k := range keys {
			v := d.Values[k]
			hi := lo + v
			arc := ComputedArc{
				ID:             d.Index + "." + k,
				Index:          d.Index,
				Key:            k,
				Value:          v,
				FormattedValue: format(v),
				StackLo:        lo,
				StackHi:        hi,
				Color:          getColor(k),
				Arc: arcs.Arc{
					StartAngle:  arcs.DegToRad(bandStart),
					EndAngle:    arcs.DegToRad(bandStart + bandwidth),
					InnerRadius: radiusScale.Call(lo),
					OuterRadius: radiusScale.Call(hi),
				},
			}
			computed = append(computed, arc)
			lo = hi
		}
	}

	legendData := make([]legends.Datum, 0, len(keys))
	for _, k := range keys {
		legendData = append(legendData, legends.Datum{ID: k, Label: k, Color: getColor(k)})
	}

	return PolarBarResult{
		Center:       center,
		InnerRadius:  innerRadius,
		OuterRadius:  outerRadius,
		Arcs:         computed,
		ArcGenerator: arcGen,
		AngleScale:   angleScale,
		RadiusScale:  radiusScale,
		StartAngle:   startAngle,
		EndAngle:     endAngle,
		LegendData:   legendData,
		Indices:      indices,
		MaxValue:     maxValue,
	}
}

// clampArc mirrors @nivo/core clampArc: limits the angular span to `length`
// degrees, preserving the start angle.
func clampArc(startAngle, endAngle, length float64) (float64, float64) {
	clampedEnd := endAngle
	if math.Abs(endAngle-startAngle) > length {
		if endAngle > startAngle {
			clampedEnd = startAngle + length
		} else {
			clampedEnd = startAngle - length
		}
	}
	return startAngle, clampedEnd
}

// valueFormatter returns a number→string formatter. Empty spec uses %g.
func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
