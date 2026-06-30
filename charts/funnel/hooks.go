package funnel

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// position is a simple (x, y) pair.
type position struct{ X, Y float64 }

// areaPoint carries the boundary coordinates for one funnel area vertex. For a
// vertical funnel x0/x1 vary and y is fixed per row; for horizontal y0/y1 vary
// and x is fixed.
type areaPoint struct {
	X, X0, X1 float64
	Y, Y0, Y1 float64
}

// UseFunnel mirrors @nivo/funnel's useFunnel: it builds the band + linear
// scales, lays out each part as a trapezoid (with shape-blended curved sides),
// generates the area + border paths, and computes the separators.
// props.Width/Height are the inner (margin-subtracted) dimensions.
func UseFunnel(props FunnelProps) FunnelResult {
	theme := resolveTheme(props.Theme)
	vertical := props.Direction != FunnelDirectionHorizontal
	n := len(props.Data)
	if n == 0 {
		return FunnelResult{}
	}

	paddingBefore := 0.0
	if props.BeforeSeparatorsEnabled() {
		paddingBefore = props.BeforeSeparatorOffset
	}
	paddingAfter := 0.0
	if props.AfterSeparatorsEnabled() {
		paddingAfter = props.AfterSeparatorOffset
	}

	var innerWidth, innerHeight float64
	if vertical {
		innerWidth = props.Width - paddingBefore - paddingAfter
		innerHeight = props.Height
	} else {
		innerWidth = props.Width
		innerHeight = props.Height - paddingBefore - paddingAfter
	}

	bandSize := innerHeight
	linearSize := innerWidth
	if !vertical {
		bandSize = innerWidth
		linearSize = innerHeight
	}
	bandwidth := (bandSize - props.Spacing*float64(n-1)) / float64(n)
	bandScale := func(i int) float64 { return props.Spacing*float64(i) + bandwidth*float64(i) }

	maxValue := 0.0
	for _, d := range props.Data {
		if d.Value > maxValue {
			maxValue = d.Value
		}
	}
	linearScale := scales.NewLinearScaleWithRange(0, maxValue, 0, linearSize)

	getColor := colors.GetOrdinalColorScale[FunnelDatum](props.Colors, func(d FunnelDatum) string { return d.ID })
	getBorderColor := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	getLabelColor := colors.GetInheritedColorGenerator(props.LabelColor, theme)
	format := valueFormatter(props.ValueFormat)
	shapeBlending := props.ShapeBlending / 2

	curve := d3shape.CurveBasis
	if props.Interpolation == FunnelInterpolationLinear {
		curve = d3shape.CurveLinear
	}

	// Geometry pass: positions per part.
	geoms := make([]geom, n)
	for i, d := range props.Data {
		var g geom
		if vertical {
			g.w = linearScale.Call(d.Value)
			g.h = bandwidth
			g.x0 = paddingBefore + (innerWidth-g.w)*0.5
			g.y0 = bandScale(i)
		} else {
			g.w = bandwidth
			g.h = linearScale.Call(d.Value)
			g.x0 = bandScale(i)
			g.y0 = paddingBefore + (innerHeight-g.h)*0.5
		}
		g.x1 = g.x0 + g.w
		g.y1 = g.y0 + g.h
		g.x = g.x0 + g.w*0.5
		g.y = g.y0 + g.h*0.5
		geoms[i] = g
	}

	parts := make([]ComputedPart, 0, n)
	for i, d := range props.Data {
		g := geoms[i]
		var next *geom
		if i+1 < n {
			next = &geoms[i+1]
		}

		var pts [4]position
		var areaPts []areaPoint
		if vertical {
			pts[0] = position{g.x0, g.y0}
			pts[1] = position{g.x1, g.y0}
			if next != nil {
				pts[2] = position{next.x1, g.y1}
				pts[3] = position{next.x0, g.y1}
			} else {
				pts[2] = position{pts[1].X, g.y1}
				pts[3] = position{pts[0].X, g.y1}
			}
			areaPts = []areaPoint{
				{X0: pts[0].X, X1: pts[1].X, Y: g.y0},
				{X0: pts[0].X, X1: pts[1].X, Y: g.y0 + g.h*shapeBlending},
				{X0: pts[3].X, X1: pts[2].X, Y: g.y1 - g.h*shapeBlending},
				{X0: pts[3].X, X1: pts[2].X, Y: g.y1},
			}
		} else {
			pts[0] = position{g.x0, g.y0}
			if next != nil {
				pts[1] = position{g.x1, next.y0}
				pts[2] = position{g.x1, next.y1}
			} else {
				pts[1] = position{g.x1, g.y0}
				pts[2] = position{g.x1, g.y1}
			}
			pts[3] = position{g.x0, g.y1}
			areaPts = []areaPoint{
				{X: g.x0, Y0: pts[0].Y, Y1: pts[3].Y},
				{X: g.x0 + g.w*shapeBlending, Y0: pts[0].Y, Y1: pts[3].Y},
				{X: g.x1 - g.w*shapeBlending, Y0: pts[1].Y, Y1: pts[2].Y},
				{X: g.x1, Y0: pts[1].Y, Y1: pts[2].Y},
			}
		}

		areaPath := buildAreaPath(areaPts, vertical, curve)
		borderA, borderB := buildBorderPaths(areaPts, vertical, curve)

		ctx := map[string]any{"color": getColor(d)}
		color := getColor(d)
		parts = append(parts, ComputedPart{
			ID:             d.ID,
			FormattedValue: format(d.Value),
			Label:          resolveLabel(d),
			Color:          color,
			FillOpacity:    props.FillOpacity,
			BorderColor:    getBorderColor(ctx),
			BorderWidth:    props.BorderWidth,
			BorderOpacity:  props.BorderOpacity,
			LabelColor:     getLabelColor(ctx),
			X:              g.x,
			Y:              g.y,
			AreaPath:       areaPath,
			BorderPathA:    borderA,
			BorderPathB:    borderB,
		})
	}

	before, after := computeSeparators(props, geoms, vertical, innerWidth, innerHeight)

	return FunnelResult{Parts: parts, BeforeSeparators: before, AfterSeparators: after}
}

// buildAreaPath generates the trapezoid area path via the d3 area generator,
// configured for the funnel direction.
func buildAreaPath(pts []areaPoint, vertical bool, curve d3shape.CurveFactory) string {
	if vertical {
		gen := d3shape.NewAreaTyped[areaPoint](
			func(p areaPoint, _ int, _ []areaPoint) float64 { return p.X0 },
			func(p areaPoint, _ int, _ []areaPoint) float64 { return p.Y },
			func(p areaPoint, _ int, _ []areaPoint) float64 { return p.Y },
		)
		gen.X1(func(p areaPoint, _ int, _ []areaPoint) float64 { return p.X1 })
		gen.Y(func(p areaPoint, _ int, _ []areaPoint) float64 { return p.Y })
		gen.Curve(curve)
		return gen.Call(pts)
	}
	gen := d3shape.NewAreaTyped[areaPoint](
		func(p areaPoint, _ int, _ []areaPoint) float64 { return p.X },
		func(p areaPoint, _ int, _ []areaPoint) float64 { return p.Y0 },
		func(p areaPoint, _ int, _ []areaPoint) float64 { return p.Y1 },
	)
	gen.Curve(curve)
	return gen.Call(pts)
}

// buildBorderPaths generates the two side border line paths.
func buildBorderPaths(pts []areaPoint, vertical bool, curve d3shape.CurveFactory) (string, string) {
	a := make([]d3shape.Point2D, len(pts))
	b := make([]d3shape.Point2D, len(pts))
	for i, p := range pts {
		if vertical {
			a[i] = d3shape.Point2D{p.X0, p.Y}
			b[i] = d3shape.Point2D{p.X1, p.Y}
		} else {
			a[i] = d3shape.Point2D{p.X, p.Y0}
			b[i] = d3shape.Point2D{p.X, p.Y1}
		}
	}
	gen := d3shape.NewLine().Curve(curve)
	return gen.Call(a), gen.Call(b)
}

// computeSeparators mirrors @nivo/funnel computeSeparators (offsets honored;
// the separator length prop is fixed at 0 in this port).
func computeSeparators(props FunnelProps, geoms []geom, vertical bool, innerWidth, innerHeight float64) ([]Separator, []Separator) {
	var before, after []Separator
	n := len(geoms)
	if vertical {
		for _, g := range geoms {
			y := g.y0 - props.Spacing/2
			if props.BeforeSeparatorsEnabled() {
				before = append(before, Separator{X0: 0, Y0: y, X1: g.x0 - props.BeforeSeparatorOffset, Y1: y})
			}
			if props.AfterSeparatorsEnabled() {
				after = append(after, Separator{X0: g.x1 + props.AfterSeparatorOffset, Y0: y, X1: props.Width, Y1: y})
			}
		}
		last := geoms[n-1]
		if props.BeforeSeparatorsEnabled() {
			before = append(before, Separator{X0: 0, Y0: last.y1, X1: last.x0 - props.BeforeSeparatorOffset, Y1: last.y1})
		}
		if props.AfterSeparatorsEnabled() {
			after = append(after, Separator{X0: last.x1 + props.AfterSeparatorOffset, Y0: last.y1, X1: props.Width, Y1: last.y1})
		}
	} else {
		for _, g := range geoms {
			x := g.x0 - props.Spacing/2
			if props.BeforeSeparatorsEnabled() {
				before = append(before, Separator{X0: x, Y0: 0, X1: x, Y1: g.y0 - props.BeforeSeparatorOffset})
			}
			if props.AfterSeparatorsEnabled() {
				after = append(after, Separator{X0: x, Y0: g.y1 + props.AfterSeparatorOffset, X1: x, Y1: props.Height})
			}
		}
		last := geoms[n-1]
		if props.BeforeSeparatorsEnabled() {
			before = append(before, Separator{X0: last.x1, Y0: 0, X1: last.x1, Y1: last.y0 - props.BeforeSeparatorOffset})
		}
		if props.AfterSeparatorsEnabled() {
			after = append(after, Separator{X0: last.x1, Y0: last.y1 + props.AfterSeparatorOffset, X1: last.x1, Y1: props.Height})
		}
	}
	return before, after
}

// geom is the per-part pixel geometry (shared by UseFunnel + computeSeparators).
type geom struct {
	x0, y0, x1, y1, x, y, w, h float64
}

func resolveLabel(d FunnelDatum) string {
	if d.Label != "" {
		return d.Label
	}
	return d.ID
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
