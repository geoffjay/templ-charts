package stream

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/internal/curves"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// streamPoint is one area vertex (x + lower/upper y in pixels).
type streamPoint struct {
	X  float64
	Y0 float64
	Y1 float64
}

// UseStream mirrors @nivo/stream's useStream: it stacks the data with the
// configured offset/order, builds a point x-scale over the indices and a linear
// y-scale over the stacked extent, generates one smooth area path per layer,
// and resolves per-layer colors. props.Width/Height are the inner dimensions.
func UseStream(props StreamProps) StreamResult {
	theme := resolveTheme(props.Theme)

	stack := d3shape.NewStack[StreamDatum]().
		Keys(props.Keys).
		Value(func(d StreamDatum, key string, _ int, _ []StreamDatum) float64 { return d[key] }).
		Offset(curves.StackOffsetFromProp(props.OffsetType)).
		Order(curves.StackOrderFromProp(props.Order))

	series := stack.Call(props.Data)

	// Value domain over the stacked extent.
	minValue := math.Inf(1)
	maxValue := math.Inf(-1)
	for _, s := range series {
		for _, p := range s.Stats {
			if p.Lo < minValue {
				minValue = p.Lo
			}
			if p.Hi > maxValue {
				maxValue = p.Hi
			}
		}
	}
	if math.IsInf(minValue, 1) {
		minValue, maxValue = 0, 0
	}

	n := len(props.Data)
	indices := make([]any, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	xScale := scales.ComputeScale(scales.ScalePointSpec{}, scales.ComputedSerieAxis{All: indices}, props.Width, scales.ScaleAxisX)
	yScale := scales.ComputeScale(
		scales.ScaleLinearSpec{Min: scales.FloatVal(minValue), Max: scales.FloatVal(maxValue)},
		scales.ComputedSerieAxis{All: []any{minValue, maxValue}, Min: minValue, Max: maxValue},
		props.Height, scales.ScaleAxisY,
	)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)

	areaGen := d3shape.NewAreaTyped[streamPoint](
		func(p streamPoint, _ int, _ []streamPoint) float64 { return p.X },
		func(p streamPoint, _ int, _ []streamPoint) float64 { return p.Y0 },
		func(p streamPoint, _ int, _ []streamPoint) float64 { return p.Y1 },
	).Curve(curves.CurveFromProp(props.Curve))

	layers := make([]ComputedLayer, 0, len(series))
	for li, s := range series {
		key := props.Keys[li]
		pts := make([]streamPoint, len(s.Stats))
		for i, p := range s.Stats {
			pts[i] = streamPoint{
				X:  xScale.Call(i),
				Y0: yScale.Call(p.Lo),
				Y1: yScale.Call(p.Hi),
			}
		}
		color := getColor(key)
		layers = append(layers, ComputedLayer{
			ID:          key,
			Label:       key,
			Path:        areaGen.Call(pts),
			Color:       color,
			BorderColor: getBorderColor(map[string]any{"color": color}),
		})
	}

	legendData := make([]legends.Datum, 0, len(props.Keys))
	for _, key := range props.Keys {
		legendData = append(legendData, legends.Datum{ID: key, Label: key, Color: getColor(key)})
	}

	return StreamResult{
		Layers:     layers,
		XScale:     xScale,
		YScale:     yScale,
		LegendData: legendData,
	}
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
