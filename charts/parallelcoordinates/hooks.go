package parallelcoordinates

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseParallelCoordinates mirrors @nivo/parallel-coordinates' useParallelCoordinates:
// it positions each variable along the layout axis, resolves a scale per
// variable (linear or point), assigns an ordinal color per datum, and builds a
// polyline per datum across the variable axes. props.Width/Height are the inner
// (margin-subtracted) dimensions.
func UseParallelCoordinates(props PCProps) PCResult {
	horizontal := props.Layout != PCLayoutVertical
	layoutLen := props.Width
	crossLen := props.Height
	if !horizontal {
		layoutLen = props.Height
		crossLen = props.Width
	}

	n := len(props.Variables)
	vars := make([]ComputedVariable, n)
	for i, v := range props.Variables {
		pos := layoutLen / 2
		if n > 1 {
			pos = float64(i) * (layoutLen / float64(n-1))
		}
		label := v.Label
		if label == "" {
			label = v.Key
		}
		vars[i] = ComputedVariable{
			Key:   v.Key,
			Label: label,
			Type:  v.Type,
			Pos:   pos,
			Scale: buildVariableScale(v, props.Data, horizontal, crossLen),
		}
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })

	gen := d3shape.NewLineTyped[d3shape.Point2D](
		func(d d3shape.Point2D, _ int, _ []d3shape.Point2D) float64 { return d[0] },
		func(d d3shape.Point2D, _ int, _ []d3shape.Point2D) float64 { return d[1] },
	)
	gen.Defined(func(d d3shape.Point2D, _ int, _ []d3shape.Point2D) bool {
		return !math.IsNaN(d[0]) && !math.IsNaN(d[1])
	})
	gen.Curve(core.CurveFromProp(props.Curve))

	lines := make([]ComputedLine, 0, len(props.Data))
	for di, datum := range props.Data {
		id := datum.ID
		if id == "" {
			id = itoa(di)
		}
		color := getColor(id)
		pts := make([]d3shape.Point2D, n)
		for i, v := range vars {
			raw, ok := datum.Values[v.Key]
			if !ok || raw == nil {
				pts[i] = d3shape.Point2D{math.NaN(), math.NaN()}
				continue
			}
			cross := v.Scale.Call(raw)
			if horizontal {
				pts[i] = d3shape.Point2D{v.Pos, cross}
			} else {
				pts[i] = d3shape.Point2D{cross, v.Pos}
			}
		}
		line := ComputedLine{ID: id, Color: color, Path: gen.Call(pts)}
		for _, p := range pts {
			if math.IsNaN(p[0]) || math.IsNaN(p[1]) {
				continue
			}
			line.Points = append(line.Points, [2]float64{p[0], p[1]})
		}
		if props.Interactive {
			line.Tooltip = interact.TooltipHTML(color, id, "")
		}
		lines = append(lines, line)
	}

	legendData := make([]legends.Datum, 0, len(props.Data))
	for di, datum := range props.Data {
		id := datum.ID
		if id == "" {
			id = itoa(di)
		}
		legendData = append(legendData, legends.Datum{ID: id, Label: id, Color: getColor(id)})
	}

	return PCResult{Variables: vars, Lines: lines, LegendData: legendData}
}

// buildVariableScale resolves a variable's scale. For horizontal layout the
// cross axis is vertical (value min at the bottom); for vertical layout it is
// horizontal (value min at the left).
func buildVariableScale(v PCVariable, data []PCDatum, horizontal bool, crossLen float64) scales.Scale {
	if v.Type == PCScalePoint {
		cats := v.Values
		if len(cats) == 0 {
			seen := map[string]bool{}
			for _, d := range data {
				if s, ok := d.Values[v.Key].(string); ok && !seen[s] {
					seen[s] = true
					cats = append(cats, s)
				}
			}
		}
		if horizontal {
			return scales.NewPointScaleWithRange(cats, 0, crossLen, v.Padding, false)
		}
		return scales.NewPointScaleWithRange(cats, 0, crossLen, v.Padding, false)
	}
	// linear
	min, max := math.Inf(1), math.Inf(-1)
	for _, d := range data {
		if f, ok := toFloatOK(d.Values[v.Key]); ok {
			min = math.Min(min, f)
			max = math.Max(max, f)
		}
	}
	if math.IsInf(min, 1) {
		min, max = 0, 1
	}
	if v.Min != nil {
		min = *v.Min
	}
	if v.Max != nil {
		max = *v.Max
	}
	if min == max {
		max = min + 1
	}
	if horizontal {
		// vertical cross axis: max at the top (pixel 0), min at the bottom.
		return scales.NewLinearScaleWithRange(min, max, crossLen, 0)
	}
	// horizontal cross axis: min at the left (pixel 0), max at the right.
	return scales.NewLinearScaleWithRange(min, max, 0, crossLen)
}

func toFloatOK(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	p := len(buf)
	for i > 0 {
		p--
		buf[p] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		p--
		buf[p] = '-'
	}
	return string(buf[p:])
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
