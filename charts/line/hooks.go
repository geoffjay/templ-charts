package line

import (
	"fmt"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/internal/curves"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// UseLineGenerator mirrors @nivo/line useLineGenerator: returns a func that
// produces an SVG line path-data string from a slice of {x,y} points, using
// the given curve and honoring null positions (defined = x!=null && y!=null).
func UseLineGenerator(curve core.CurveFactoryId) LineGenerator {
	curveFactory := curves.CurveFromProp(curve)
	gen := d3shape.NewLineTyped[PointXY](
		func(d PointXY, _ int, _ []PointXY) float64 { return d.X },
		func(d PointXY, _ int, _ []PointXY) float64 { return d.Y },
	)
	gen.Defined(func(d PointXY, _ int, _ []PointXY) bool {
		return !isNilFloat(d.X) && !isNilFloat(d.Y)
	})
	gen.Curve(curveFactory)
	return func(points []PointXY) string {
		return gen.Call(points)
	}
}

// UseAreaGenerator mirrors @nivo/line useAreaGenerator: returns a func that
// produces an SVG area path from points, with the baseline at
// yScale(areaBaselineValue).
func UseAreaGenerator(curve core.CurveFactoryId, yScale scales.Scale, areaBaselineValue float64) AreaGenerator {
	curveFactory := curves.CurveFromProp(curve)
	baseline := yScale.Call(areaBaselineValue)
	gen := d3shape.NewAreaTyped[PointXY](
		func(d PointXY, _ int, _ []PointXY) float64 { return d.X },
		func(_ PointXY, _ int, _ []PointXY) float64 { return baseline },
		func(d PointXY, _ int, _ []PointXY) float64 { return d.Y },
	)
	gen.Defined(func(d PointXY, _ int, _ []PointXY) bool {
		return !isNilFloat(d.X) && !isNilFloat(d.Y)
	})
	gen.Curve(curveFactory)
	return func(points []PointXY) string {
		return gen.Call(points)
	}
}

// UsePoints mirrors @nivo/line usePoints: flattens the computed series into a
// single []Point, filtering out points with null positions, and resolving
// point color + border color.
func UsePoints(
	series []ComputedSeries,
	getPointColor func(PointColorContext) string,
	getPointBorderColor func(Point) string,
	formatX func(any) string,
	formatY func(any) string,
) []Point {
	out := []Point{}
	for seriesIndex, s := range series {
		for indexInSeries, d := range s.Data {
			if isNilFloat(d.Position.X) || isNilFloat(d.Position.Y) {
				continue
			}
			p := Point{
				ID:            fmt.Sprintf("%s.%d", s.ID, indexInSeries),
				IndexInSeries: indexInSeries,
				AbsIndex:      len(out) + indexInSeries,
				SeriesIndex:   seriesIndex,
				SeriesID:      s.ID,
				SeriesColor:   s.Color,
				X:             d.Position.X,
				Y:             d.Position.Y,
				Data: PointDatum{
					X:          d.Data.X,
					Y:          d.Data.Y,
					XFormatted: formatX(d.Data.X),
					YFormatted: formatY(d.Data.Y),
				},
			}
			p.Color = getPointColor(PointColorContext{
				Series: s,
				Point: PointColorContextPoint{
					ID: p.ID, IndexInSeries: p.IndexInSeries, AbsIndex: p.AbsIndex,
					SeriesIndex: p.SeriesIndex, SeriesID: p.SeriesID, SeriesColor: p.SeriesColor,
					X: p.X, Y: p.Y, Data: p.Data,
				},
			})
			p.BorderColor = getPointBorderColor(p)
			out = append(out, p)
		}
	}
	return out
}

// UseSlices mirrors @nivo/line useSlices: groups points by x ("x" mode) or y
// ("y" mode) and computes each slice's geometry (x0/width or y0/height as the
// midpoint to neighboring slices).
func UseSlices(componentID string, enableSlices EnableSlices, points []Point, width, height float64) []SliceData {
	if enableSlices == "x" {
		groups := map[float64][]Point{}
		var xs []float64
		for _, p := range points {
			if isNilAny(p.Data.X) || isNilAny(p.Data.Y) {
				continue
			}
			if _, ok := groups[p.X]; !ok {
				xs = append(xs, p.X)
			}
			groups[p.X] = append(groups[p.X], p)
		}
		sortFloats(xs)
		out := make([]SliceData, len(xs))
		for i, x := range xs {
			pts := groups[x]
			// reverse to match nivo's slicePoints.reverse()
			reversePoints(pts)
			var x0, sliceWidth float64
			if i == 0 {
				x0 = x
			} else {
				x0 = x - (x-xs[i-1])/2
			}
			if i == len(xs)-1 {
				sliceWidth = width - x0
			} else {
				sliceWidth = x - x0 + (xs[i+1]-x)/2
			}
			out[i] = SliceData{
				ID:     fmt.Sprintf("slice:%s:%g", componentID, x),
				X0:     x0,
				X:      x,
				Y0:     0,
				Y:      0,
				Width:  sliceWidth,
				Height: height,
				Points: pts,
			}
		}
		return out
	}
	if enableSlices == "y" {
		groups := map[float64][]Point{}
		var ys []float64
		for _, p := range points {
			if isNilAny(p.Data.X) || isNilAny(p.Data.Y) {
				continue
			}
			if _, ok := groups[p.Y]; !ok {
				ys = append(ys, p.Y)
			}
			groups[p.Y] = append(groups[p.Y], p)
		}
		sortFloats(ys)
		out := make([]SliceData, len(ys))
		for i, y := range ys {
			pts := groups[y]
			reversePoints(pts)
			var y0, sliceHeight float64
			if i == 0 {
				y0 = y
			} else {
				y0 = y - (y-ys[i-1])/2
			}
			if i == len(ys)-1 {
				sliceHeight = height - y0
			} else {
				sliceHeight = y - y0 + (ys[i+1]-y)/2
			}
			out[i] = SliceData{
				ID:     fmt.Sprintf("%g", y),
				X0:     0,
				X:      0,
				Y0:     y0,
				Y:      y,
				Width:  width,
				Height: sliceHeight,
				Points: pts,
			}
		}
		return out
	}
	return nil
}

// UseLine mirrors @nivo/line useLine: the full orchestrator. Resolves
// accessors, color scales, formatters; computes the x/y scales + series
// positions via scales.ComputeXYScalesForSeries; builds legend data, points,
// slices, and the line/area generators.
func UseLine(props LineProps) LineResult {
	// merge with defaults
	xScaleSpec := props.XScale
	if xScaleSpec == nil {
		xScaleSpec = Defaults.XScale
	}
	yScaleSpec := props.YScale
	if yScaleSpec == nil {
		yScaleSpec = Defaults.YScale
	}
	colorsCfg := props.Colors
	if colorsCfg.Type == 0 && colorsCfg.Scheme == "" && colorsCfg.Static == "" && len(colorsCfg.Colors) == 0 && colorsCfg.Func == nil && colorsCfg.DatumPath == "" {
		colorsCfg = Defaults.Colors
	}
	curve := props.Curve
	if curve == "" {
		curve = Defaults.Curve
	}
	areaBaselineValue := props.AreaBaselineValue
	initialHidden := props.InitialHiddenIDs
	if initialHidden == nil {
		initialHidden = Defaults.InitialHiddenIDs
	}
	enableSlices := props.EnableSlices

	theme := props.Theme
	if theme == nil {
		theme = &theming.DefaultTheme
	}

	formatX := core.GetValueFormatter[any](props.XFormat)
	formatY := core.GetValueFormatter[any](props.YFormat)

	// color scale: series id → color
	getColor := colors.GetOrdinalColorScale[LineSeries](colorsCfg, func(s LineSeries) string { return s.ID })

	// point color / border color
	pointColor := props.PointColor
	if pointColor.Type == 0 && pointColor.Static == "" && pointColor.ThemePath == "" && pointColor.FromPath == "" && pointColor.Func == nil {
		pointColor = Defaults.PointColor
	}
	pointBorderColor := props.PointBorderColor
	if pointBorderColor.Type == 0 && pointBorderColor.Static == "" && pointBorderColor.ThemePath == "" && pointBorderColor.FromPath == "" && pointBorderColor.Func == nil {
		pointBorderColor = Defaults.PointBorderColor
	}
	getPointColor := func(ctx PointColorContext) string {
		return colors.GetInheritedColorGenerator(pointColor, theme)(map[string]any{
			"series": map[string]any{
				"id":    ctx.Series.ID,
				"color": ctx.Series.Color,
			},
			"point": map[string]any{
				"id":            ctx.Point.ID,
				"indexInSeries": ctx.Point.IndexInSeries,
				"absIndex":      ctx.Point.AbsIndex,
				"seriesIndex":   ctx.Point.SeriesIndex,
				"seriesId":      ctx.Point.SeriesID,
				"seriesColor":   ctx.Point.SeriesColor,
				"x":             ctx.Point.X,
				"y":             ctx.Point.Y,
				"data":          ctx.Point.Data,
			},
		})
	}
	getPointBorderColor := func(p Point) string {
		return colors.GetInheritedColorGenerator(pointBorderColor, theme)(map[string]any{
			"id":            p.ID,
			"indexInSeries": p.IndexInSeries,
			"absIndex":      p.AbsIndex,
			"seriesIndex":   p.SeriesIndex,
			"seriesId":      p.SeriesID,
			"seriesColor":   p.SeriesColor,
			"x":             p.X,
			"y":             p.Y,
			"data":          p.Data,
		})
	}

	// filter hidden series
	visible := make([]LineSeries, 0, len(props.Data))
	for _, s := range props.Data {
		if !containsString(initialHidden, s.ID) {
			visible = append(visible, s)
		}
	}

	// build scales.Serie slice from LineSeries
	seriesInput := make([]scales.Serie, len(visible))
	for i, s := range visible {
		data := make([]scales.SerieDatum, len(s.Data))
		for j, d := range s.Data {
			data[j] = scales.SerieDatum{X: d.X, Y: d.Y}
		}
		seriesInput[i] = scales.Serie{Data: data, Extra: s.ID}
	}

	xy := scales.ComputeXYScalesForSeries(seriesInput, xScaleSpec, yScaleSpec, props.Width, props.Height)

	// build ComputedSeries: match rawSeries back to input + assign colors
	rawSeries := xy.Series
	computed := make([]ComputedSeries, 0, len(visible))
	for _, s := range visible {
		var rs *scales.ComputedSerie
		for i := range rawSeries {
			if id, ok := rawSeries[i].Extra.(string); ok && id == s.ID {
				rs = &rawSeries[i]
				break
			}
		}
		if rs == nil {
			continue
		}
		datums := make([]ComputedDatum, len(rs.Data))
		for k, d := range rs.Data {
			datums[k] = ComputedDatum{
				Data: s.Data[k],
				Position: struct{ X, Y float64 }{
					X: ptrToFloat(d.Position.X),
					Y: ptrToFloat(d.Position.Y),
				},
			}
		}
		computed = append(computed, ComputedSeries{
			ID:    s.ID,
			Data:  datums,
			Color: getColor(s),
		})
	}

	// legend data: all series, with hidden flag, reversed (nivo: legendData.reverse())
	legendData := make([]LegendDatum, len(props.Data))
	for i, s := range props.Data {
		legendData[i] = LegendDatum{
			ID:     s.ID,
			Label:  s.ID,
			Color:  getColor(s),
			Hidden: containsString(initialHidden, s.ID),
		}
	}
	reverseLegendData(legendData)

	points := UsePoints(computed, getPointColor, getPointBorderColor, formatX, formatY)
	slices := UseSlices("line", enableSlices, points, props.Width, props.Height)
	lineGen := UseLineGenerator(curve)
	areaGen := UseAreaGenerator(curve, xy.YScale, areaBaselineValue)

	return LineResult{
		Series:        computed,
		Points:        points,
		Slices:        slices,
		LegendData:    legendData,
		XScale:        xy.XScale,
		YScale:        xy.YScale,
		LineGenerator: lineGen,
		AreaGenerator: areaGen,
		GetColor:      getColor,
		HiddenIDs:     initialHidden,
	}
}

// --- helpers ---

func isNilFloat(v float64) bool {
	// NaN guard: d3 treats null positions as undefined. We encode "missing"
	// as NaN when a scale can't resolve a nil datum.
	return v != v
}

func isNilAny(v any) bool {
	return v == nil
}

func ptrToFloat(p *float64) float64 {
	if p == nil {
		return nilFloat()
	}
	return *p
}

// nilFloat returns the sentinel NaN used to represent a null position.
func nilFloat() float64 {
	var z float64
	return z / z // NaN
}

func containsString(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func reversePoints(ps []Point) {
	for i, j := 0, len(ps)-1; i < j; i, j = i+1, j-1 {
		ps[i], ps[j] = ps[j], ps[i]
	}
}

func reverseLegendData(d []LegendDatum) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}

func sortFloats(xs []float64) {
	// simple insertion sort (slice is usually small and already near-sorted)
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j] < xs[j-1]; j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}
