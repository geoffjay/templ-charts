package marimekko

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseMarimekko mirrors @nivo/marimekko's useMarimekko: it derives each column's
// value-driven thickness, stacks the dimensions with the shared d3-shape Stack
// generator (honoring the configured offset), and positions one segment rect
// per (column, dimension). props.Width/Height are the inner dimensions.
func UseMarimekko(props MarimekkoProps) MarimekkoResult {
	horizontal := props.Layout == MarimekkoLayoutHorizontal

	dimKeys := make([]string, len(props.Dimensions))
	dimByKey := make(map[string]MarimekkoDimension, len(props.Dimensions))
	for i, d := range props.Dimensions {
		dimKeys[i] = d.Key
		dimByKey[d.Key] = d
	}

	// Thickness driver: total across columns.
	totalValue := 0.0
	for _, d := range props.Data {
		totalValue += d.Value
	}
	if totalValue == 0 {
		totalValue = 1
	}

	// Stack the dimensions across columns.
	stack := d3shape.NewStack[MarimekkoDatum]().
		Keys(dimKeys).
		Value(func(d MarimekkoDatum, key string, _ int, _ []MarimekkoDatum) float64 {
			return d.Dimensions[key]
		}).
		Offset(offsetFunc(props.Offset))
	series := stack.Call(props.Data)

	stackMax := 0.0
	for _, s := range series {
		for _, p := range s.Stats {
			if p.Hi > stackMax {
				stackMax = p.Hi
			}
		}
	}
	if stackMax == 0 {
		stackMax = 1
	}

	// Thickness axis span (width for vertical, height for horizontal).
	n := len(props.Data)
	thicknessLen := props.Width
	stackLen := props.Height
	if horizontal {
		thicknessLen = props.Height
		stackLen = props.Width
	}
	gaps := 0.0
	if n > 1 {
		gaps = props.InnerPadding * float64(n-1)
	}
	available := thicknessLen - 2*props.OuterPadding - gaps
	if available < 0 {
		available = 0
	}

	// Stack value scale: [0, stackMax] → pixels along the stack axis.
	var stackScale, thicknessScale scales.Scale
	if horizontal {
		// horizontal layout: stack runs left→right, thickness down the height.
		stackScale = scales.NewLinearScaleWithRange(0, stackMax, 0, stackLen)
		thicknessScale = scales.NewLinearScaleWithRange(0, totalValue, 0, thicknessLen)
	} else {
		// vertical layout: stack runs bottom→top (0 at the bottom).
		stackScale = scales.NewLinearScaleWithRange(0, stackMax, stackLen, 0)
		thicknessScale = scales.NewLinearScaleWithRange(0, totalValue, 0, thicknessLen)
	}

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)

	bars := make([]ComputedBar, 0, n*len(dimKeys))
	offset := props.OuterPadding
	for i, datum := range props.Data {
		thickness := datum.Value / totalValue * available
		start := offset
		for di, key := range dimKeys {
			dim := dimByKey[key]
			p := series[di].Stats[i]
			raw := datum.Dimensions[key]
			color := getColor(dim.ID)
			var bar ComputedBar
			if horizontal {
				// thickness along y (top→bottom), stack along x.
				x := stackScale.Call(p.Lo)
				w := stackScale.Call(p.Hi) - x
				bar = ComputedBar{
					X: x, Y: start, Width: w, Height: thickness,
				}
			} else {
				// thickness along x, stack along y (top of segment = Hi).
				yTop := stackScale.Call(p.Hi)
				h := stackScale.Call(p.Lo) - yTop
				bar = ComputedBar{
					X: start, Y: yTop, Width: thickness, Height: h,
				}
			}
			bar.ID = datum.ID + "." + dim.ID
			bar.Index = datum.ID
			bar.DimensionID = dim.ID
			bar.Value = raw
			bar.FormattedValue = format(raw)
			bar.Color = color
			bars = append(bars, bar)
		}
		offset = start + thickness + props.InnerPadding
	}

	legendData := make([]legends.Datum, 0, len(props.Dimensions))
	for _, d := range props.Dimensions {
		legendData = append(legendData, legends.Datum{ID: d.ID, Label: d.ID, Color: getColor(d.ID)})
	}

	return MarimekkoResult{
		Bars:           bars,
		LegendData:     legendData,
		ThicknessMax:   totalValue,
		StackMax:       stackMax,
		ThicknessScale: thicknessScale,
		StackScale:     stackScale,
	}
}

// offsetFunc maps an OffsetType to the d3-shape stack offset function.
func offsetFunc(o OffsetType) d3shape.StackOffsetFunc {
	switch o {
	case OffsetExpand:
		return d3shape.StackOffsetExpand
	case OffsetDiverging:
		return d3shape.StackOffsetDiverging
	case OffsetSilhouette:
		return d3shape.StackOffsetSilhouette
	case OffsetWiggle:
		return d3shape.StackOffsetWiggle
	default:
		return d3shape.StackOffsetNone
	}
}

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
