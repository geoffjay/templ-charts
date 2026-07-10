package radar

import (
	"fmt"
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/internal/curves"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseRadar mirrors @nivo/radar's useRadar: it derives the indices, the linear
// radius scale, the center/angle step, the per-key color mapping, and the
// projected closed polygons (one per key). props.Width/Height are the inner
// (margin-subtracted) dimensions, matching the Radar templ component's
// contract.
func UseRadar(props RadarProps) RadarResult {
	theme := resolveTheme(props.Theme)

	// Indices (one per row) and the value domain max.
	indices := make([]string, 0, len(props.Data))
	maxVal := math.Inf(-1)
	hasVal := false
	for _, row := range props.Data {
		indices = append(indices, toString(row[props.IndexBy]))
		for _, key := range props.Keys {
			v := toFloat(row[key])
			if v > maxVal {
				maxVal = v
			}
			hasVal = true
		}
	}
	if !hasVal {
		maxVal = 0
	}
	if props.MaxValue != nil {
		maxVal = *props.MaxValue
	}

	radius := math.Min(props.Width, props.Height) / 2
	radiusScale := scales.ComputeScale(
		scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.FloatVal(maxVal)},
		scales.ComputedSerieAxis{All: []any{0.0, maxVal}, Min: 0.0, Max: maxVal},
		radius, scales.ScaleAxisX,
	)
	centerX := props.Width / 2
	centerY := props.Height / 2
	rotation := arcs.DegToRad(props.Rotation)
	angleStep := 0.0
	if len(props.Data) > 0 {
		angleStep = (math.Pi * 2) / float64(len(props.Data))
	}

	// Per-key color mapping (ordinal scale keyed by key, like nivo).
	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(k string) string { return k })
	colorByKey := make(map[string]string, len(props.Keys))
	for _, key := range props.Keys {
		colorByKey[key] = getColor(key)
	}

	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	getDotColor := colors.GetInheritedColorGenerator(props.DotColor, theme)
	getDotBorderColor := colors.GetInheritedColorGenerator(props.DotBorderColor, theme)
	format := valueFormatter(props.ValueFormat)

	lineGen := d3shape.NewLine().Curve(curves.CurveFromProp(props.Curve))

	series := make([]ComputedSerie, 0, len(props.Keys))
	allPoints := make([]ComputedPoint, 0, len(props.Keys)*len(props.Data))
	for _, key := range props.Keys {
		color := colorByKey[key]
		pts := make([]ComputedPoint, 0, len(props.Data))
		linePoints := make([]d3shape.Point2D, 0, len(props.Data))
		for i, row := range props.Data {
			value := toFloat(row[key])
			angle := rotation + float64(i)*angleStep - math.Pi/2
			p := arcs.PositionFromAngle(angle, radiusScale.Call(value))
			ctx := map[string]any{"color": color, "key": key}
			cp := ComputedPoint{
				Key:            key,
				Index:          indices[i],
				Value:          value,
				FormattedValue: format(value),
				X:              p.X,
				Y:              p.Y,
				Color:          getDotColor(ctx),
				BorderColor:    getDotBorderColor(ctx),
			}
			pts = append(pts, cp)
			allPoints = append(allPoints, cp)
			linePoints = append(linePoints, d3shape.Point2D{p.X, p.Y})
		}
		stroke := getBorderColor(map[string]any{"color": color, "key": key})
		series = append(series, ComputedSerie{
			Key:    key,
			Path:   lineGen.Call(linePoints),
			Color:  color,
			Stroke: stroke,
			Points: pts,
		})
	}

	legendData := make([]legends.Datum, 0, len(props.Keys))
	for _, key := range props.Keys {
		legendData = append(legendData, legends.Datum{ID: key, Label: key, Color: colorByKey[key]})
	}

	return RadarResult{
		Indices:     indices,
		Keys:        props.Keys,
		ColorByKey:  colorByKey,
		Series:      series,
		Points:      allPoints,
		RadiusScale: radiusScale,
		Radius:      radius,
		CenterX:     centerX,
		CenterY:     centerY,
		Rotation:    rotation,
		AngleStep:   angleStep,
		MaxValue:    maxVal,
		LegendData:  legendData,
	}
}

// valueFormatter returns a number→string formatter. Empty spec uses %g; a
// non-empty spec is a d3-format spec.
func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color (or "" for transparent).
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		if f, err := strconv.ParseFloat(x, 64); err == nil {
			return f
		}
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
