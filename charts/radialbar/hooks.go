package radialbar

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

// UseRadialBar mirrors @nivo/radial-bar's useRadialBar: it derives the center,
// inner/outer radii, the value (angle) + radius (band) scales, the stacked arc
// bars, the background tracks, and the legend data. props.Width/Height are the
// inner (margin-subtracted) dimensions.
func UseRadialBar(props RadialBarProps) RadialBarResult {
	theme := resolveTheme(props.Theme)

	center := [2]float64{props.Width / 2, props.Height / 2}
	outerRadius := math.Min(center[0], center[1])
	innerRadius := outerRadius * math.Min(props.InnerRadius, 1)

	startAngle, endAngle := clampArc(props.StartAngle, props.EndAngle, 360)

	// Extract serie ids, categories (first-seen), group totals.
	serieIDs := make([]string, 0, len(props.Data))
	categories := make([]string, 0)
	seenCat := map[string]bool{}
	type group struct {
		id    string
		total float64
		data  []RadialBarDatum
	}
	groups := make([]group, 0, len(props.Data))
	for _, serie := range props.Data {
		serieIDs = append(serieIDs, serie.ID)
		total := 0.0
		for _, d := range serie.Data {
			if !seenCat[d.X] {
				seenCat[d.X] = true
				categories = append(categories, d.X)
			}
			total += d.Y
		}
		groups = append(groups, group{id: serie.ID, total: total, data: serie.Data})
	}

	maxValue := 0.0
	if props.MaxValue != nil {
		maxValue = *props.MaxValue
	} else {
		for _, g := range groups {
			if g.total > maxValue {
				maxValue = g.total
			}
		}
	}

	valueScale := scales.NewLinearScaleWithRange(0, maxValue, startAngle, endAngle)
	radiusScale := scales.NewBandScaleWithRange(serieIDs, innerRadius, outerRadius, props.Padding, false)
	bandwidth := 0.0
	if bw, ok := radiusScale.(scales.ScaleWithBandwidth); ok {
		bandwidth = bw.Bandwidth()
	}

	arcGen := arcs.CreateArcGenerator(props.CornerRadius, arcs.DegToRad(props.PadAngle))
	getColor := colors.GetOrdinalColorScale[ComputedBar](props.Colors, func(b ComputedBar) string { return b.Category })
	format := valueFormatter(props.ValueFormat)

	bars := make([]ComputedBar, 0)
	for _, g := range groups {
		currentValue := 0.0
		arcInner := radiusScale.Call(g.id)
		arcOuter := arcInner + bandwidth
		for _, d := range g.data {
			stacked := currentValue + d.Y
			bar := ComputedBar{
				ID:             g.id + "." + d.X,
				Data:           d,
				GroupID:        g.id,
				Category:       d.X,
				Value:          d.Y,
				FormattedValue: format(d.Y),
				StackedValue:   stacked,
				Arc: arcs.Arc{
					StartAngle:  arcs.DegToRad(valueScale.Call(currentValue)),
					EndAngle:    arcs.DegToRad(valueScale.Call(stacked)),
					InnerRadius: arcInner,
					OuterRadius: arcOuter,
				},
			}
			bar.Color = getColor(bar)
			currentValue += d.Y
			bars = append(bars, bar)
		}
	}

	startRad := arcs.DegToRad(startAngle)
	endRad := arcs.DegToRad(endAngle)
	var tracks []RadialBarTrackDatum
	if props.TracksEnabled() {
		for _, value := range scales.GetScaleTicks(radiusScale, scales.TicksSpec{}) {
			trackRadius := radiusScale.Call(value)
			tracks = append(tracks, RadialBarTrackDatum{
				ID:    toStr(value),
				Color: props.TracksColor,
				Arc: arcs.Arc{
					StartAngle:  startRad,
					EndAngle:    endRad,
					InnerRadius: trackRadius,
					OuterRadius: trackRadius + bandwidth,
				},
			})
		}
	}

	legendData := make([]legends.Datum, 0, len(categories))
	for _, cat := range categories {
		color := ""
		for _, b := range bars {
			if b.Category == cat {
				color = b.Color
				break
			}
		}
		legendData = append(legendData, legends.Datum{ID: cat, Label: cat, Color: color})
	}

	_ = theme
	return RadialBarResult{
		Center:       center,
		InnerRadius:  innerRadius,
		OuterRadius:  outerRadius,
		Bars:         bars,
		Tracks:       tracks,
		ArcGenerator: arcGen,
		RadiusScale:  radiusScale,
		ValueScale:   valueScale,
		StartAngle:   startAngle,
		EndAngle:     endAngle,
		LegendData:   legendData,
		Categories:   categories,
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

func toStr(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(toF(v), 'g', -1, 64)
}

func toF(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	}
	return 0
}
