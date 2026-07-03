package bump

import (
	"fmt"
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// bumpXY is one point fed to the line generator; missing ranks carry NaN so the
// generator's Defined predicate breaks the line there.
type bumpXY struct{ X, Y float64 }

// UseBump mirrors @nivo/bump's useBump: it builds the point x scale + linear
// ranking y scale via ComputeXYScalesForSeries, positions each serie's rank
// points, assigns an ordinal color per serie, and produces the smooth line
// paths. props.Width/Height are the inner (margin-subtracted) dimensions.
func UseBump(props BumpProps) BumpResult {
	// Rank domain min/max (with yOuterPadding), scanned across all numeric ranks.
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, s := range props.Data {
		for _, d := range s.Data {
			if v, ok := toFloatOK(d.Y); ok {
				minY = math.Min(minY, v)
				maxY = math.Max(maxY, v)
			}
		}
	}
	if math.IsInf(minY, 1) {
		minY, maxY = 0, 1
	}
	yMin := minY - props.YOuterPadding
	yMax := maxY + props.YOuterPadding

	seriesInput := make([]scales.Serie, len(props.Data))
	for i, s := range props.Data {
		data := make([]scales.SerieDatum, len(s.Data))
		for j, d := range s.Data {
			data[j] = scales.SerieDatum{X: d.X, Y: d.Y}
		}
		seriesInput[i] = scales.Serie{Data: data, Extra: s.ID}
	}

	xy := scales.ComputeXYScalesForSeries(
		seriesInput,
		scales.ScalePointSpec{},
		// Reverse so the smallest rank (best) sits at the top.
		scales.ScaleLinearSpec{Min: scales.FloatVal(yMin), Max: scales.FloatVal(yMax), Reverse: true},
		props.Width, props.Height,
	)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	formatX := valueFormatter(props.XFormat)
	formatY := valueFormatter(props.YFormat)

	curve := d3shape.CurveLinear
	if props.Interpolation == InterpolationSmooth {
		curve = d3shape.CurveBumpX
	}
	gen := d3shape.NewLineTyped[bumpXY](
		func(d bumpXY, _ int, _ []bumpXY) float64 { return d.X },
		func(d bumpXY, _ int, _ []bumpXY) float64 { return d.Y },
	)
	gen.Defined(func(d bumpXY, _ int, _ []bumpXY) bool { return !math.IsNaN(d.X) && !math.IsNaN(d.Y) })
	gen.Curve(curve)

	series := make([]ComputedSerie, 0, len(props.Data))
	for _, s := range props.Data {
		var rs *scales.ComputedSerie
		for i := range xy.Series {
			if id, ok := xy.Series[i].Extra.(string); ok && id == s.ID {
				rs = &xy.Series[i]
				break
			}
		}
		if rs == nil {
			continue
		}
		color := getColor(s.ID)
		pts := make([]BumpPoint, len(rs.Data))
		linePts := make([]bumpXY, len(rs.Data))
		for j, d := range rs.Data {
			x, okX := ptrToFloat(d.Position.X)
			y, okY := ptrToFloat(d.Position.Y)
			defined := okX && okY
			pts[j] = BumpPoint{
				ID:           s.ID + "." + strconv.Itoa(j),
				SerieID:      s.ID,
				IndexInSerie: j,
				X:            x,
				Y:            y,
				XValue:       s.Data[j].X,
				YValue:       s.Data[j].Y,
				FormattedX:   formatX(s.Data[j].X),
				FormattedY:   formatY(s.Data[j].Y),
				Color:        color,
				Size:         props.PointSize,
				Defined:      defined,
			}
			if defined {
				linePts[j] = bumpXY{X: x, Y: y}
			} else {
				linePts[j] = bumpXY{X: math.NaN(), Y: math.NaN()}
			}
		}
		series = append(series, ComputedSerie{
			ID:        s.ID,
			Color:     color,
			LineWidth: props.LineWidth,
			Opacity:   props.Opacity,
			LinePath:  gen.Call(linePts),
			Points:    pts,
		})
	}

	legendData := make([]legends.Datum, 0, len(props.Data))
	for _, s := range props.Data {
		legendData = append(legendData, legends.Datum{ID: s.ID, Label: s.ID, Color: getColor(s.ID)})
	}

	return BumpResult{
		Series:     series,
		XScale:     xy.XScale,
		YScale:     xy.YScale,
		LegendData: legendData,
	}
}

// valueFormatter returns an any→string formatter. Empty spec falls back to a
// generic formatter; a non-empty spec is a d3-format spec for numeric values.
func valueFormatter(spec string) func(any) string {
	if spec == "" {
		return func(v any) string {
			switch x := v.(type) {
			case float64:
				return strconv.FormatFloat(x, 'g', -1, 64)
			case int:
				return strconv.Itoa(x)
			case nil:
				return ""
			default:
				return fmt.Sprintf("%v", x)
			}
		}
	}
	return func(v any) string {
		if f, ok := toFloatOK(v); ok {
			return d3format.FormatString(spec, f)
		}
		return fmt.Sprintf("%v", v)
	}
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

// ptrToFloat returns (*p, true) or (NaN, false) when p is nil.
func ptrToFloat(p *float64) (float64, bool) {
	if p == nil {
		return math.NaN(), false
	}
	return *p, true
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
