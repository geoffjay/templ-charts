package bullet

import (
	"math"
	"sort"
	"strings"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// UseBullet mirrors @nivo/bullet's useEnhancedData + BulletItem compute: for
// each row it derives a clamped linear value scale, stacks the ranges/measures
// into colored rects, places the markers, and positions the item, its axis and
// its title. props.Width/Height are the inner (margin-subtracted) dimensions.
func UseBullet(props BulletProps) BulletResult {
	n := len(props.Data)
	if n == 0 {
		return BulletResult{}
	}
	horizontal := props.Layout != BulletLayoutVertical

	var valueSize, bandThickness float64
	if horizontal {
		valueSize = props.Width
		bandThickness = (props.Height - props.Spacing*float64(n-1)) / float64(n)
	} else {
		valueSize = props.Height
		bandThickness = (props.Width - props.Spacing*float64(n-1)) / float64(n)
	}
	measureThickness := bandThickness * props.MeasureSize
	markerThickness := bandThickness * props.MarkerSize

	rangeColors := parseColorSpec(props.RangeColors)
	measureColors := parseColorSpec(props.MeasureColors)
	markerColors := parseColorSpec(props.MarkerColors)

	items := make([]ComputedBulletItem, 0, n)
	for i, d := range props.Data {
		all := append(append(append([]float64{}, d.Ranges...), d.Measures...), d.Markers...)
		minV := props.minValueFor(all)
		maxV := props.maxValueFor(all)

		// Per-item value scale (clamped, rounded), oriented by layout/reverse.
		axisOrient := scales.ScaleAxisX
		if horizontal {
			if props.Reverse {
				axisOrient = scales.ScaleAxisY
			}
		} else {
			axisOrient = scales.ScaleAxisY
			if props.Reverse {
				axisOrient = scales.ScaleAxisX
			}
		}
		scale := scales.ComputeScale(
			scales.ScaleLinearSpec{Min: scales.FloatVal(minV), Max: scales.FloatVal(maxV), Clamp: true, Round: true},
			scales.ComputedSerieAxis{All: []any{minV, maxV}, Min: minV, Max: maxV},
			valueSize, axisOrient,
		)

		rangeColor := sequentialColor(rangeColors, minV, maxV)
		measureColor := sequentialColor(measureColors, minV, maxV)
		markerColor := sequentialColor(markerColors, minV, maxV)

		ranges := buildRects(stackValues(d.Ranges, minV, maxV, false), scale, horizontal, props.Reverse, bandThickness, 0, rangeColor)
		measureOffset := (bandThickness - measureThickness) / 2
		measures := buildRects(stackValues(d.Measures, minV, maxV, true), scale, horizontal, props.Reverse, measureThickness, measureOffset, measureColor)

		markerOffset := (bandThickness - markerThickness) / 2
		markers := make([]ComputedMarker, 0, len(d.Markers))
		for _, mv := range d.Markers {
			pos := scale.Call(mv)
			m := ComputedMarker{Color: markerColor(mv), Value: mv}
			if horizontal {
				m.X1, m.X2 = pos, pos
				m.Y1, m.Y2 = markerOffset, markerOffset+markerThickness
			} else {
				m.Y1, m.Y2 = pos, pos
				m.X1, m.X2 = markerOffset, markerOffset+markerThickness
			}
			markers = append(markers, m)
		}

		title := d.Title
		if title == "" {
			title = d.ID
		}
		titleX, titleY := titlePosition(props, horizontal, valueSize, bandThickness)

		var offsetX, offsetY float64
		if horizontal {
			offsetY = (bandThickness + props.Spacing) * float64(i)
		} else {
			offsetX = (bandThickness + props.Spacing) * float64(i)
		}

		axisX, axisY := 0.0, 0.0
		if horizontal && props.AxisPosition == "after" {
			axisY = bandThickness
		} else if !horizontal && props.AxisPosition == "after" {
			axisX = bandThickness
		}

		items = append(items, ComputedBulletItem{
			ID:       d.ID,
			Title:    title,
			OffsetX:  offsetX,
			OffsetY:  offsetY,
			Width:    valueSize,
			Height:   bandThickness,
			Ranges:   ranges,
			Measures: measures,
			Markers:  markers,
			Scale:    scale,
			AxisX:    axisX,
			AxisY:    axisY,
			TitleX:   titleX,
			TitleY:   titleY,
		})
	}

	return BulletResult{Items: items}
}

// rangeDatum is one stacked [v0,v1] segment before pixel layout.
type rangeDatum struct {
	V0, V1 float64
}

// stackValues mirrors @nivo/bullet stackValues (useAverage=false): for ranges
// it appends max (unless already present) so the last band reaches the domain
// end; for measures it stacks from min upward. Zeros are dropped, values sorted
// ascending, then each segment runs from the previous v1 (or min) to v1.
func stackValues(values []float64, minV, maxV float64, isMeasure bool) []rangeDatum {
	norm := append([]float64{}, values...)
	extra := maxV
	if isMeasure || containsFloat(values, maxV) {
		extra = 0
	}
	norm = append(norm, extra)
	filtered := norm[:0]
	for _, v := range norm {
		if v != 0 {
			filtered = append(filtered, v)
		}
	}
	sort.Float64s(filtered)
	out := make([]rangeDatum, 0, len(filtered))
	prev := minV
	for _, v1 := range filtered {
		out = append(out, rangeDatum{V0: prev, V1: v1})
		prev = v1
	}
	return out
}

// buildRects positions stacked [v0,v1] segments into pixel rects, mirroring
// @nivo/bullet getComputeRect. crossOffset shifts the band along the cross axis
// (used to center measures).
func buildRects(data []rangeDatum, scale scales.Scale, horizontal, reverse bool, thickness, crossOffset float64, color func(float64) string) []ComputedRect {
	out := make([]ComputedRect, 0, len(data))
	for _, d := range data {
		p0 := scale.Call(d.V0)
		p1 := scale.Call(d.V1)
		var r ComputedRect
		r.Color = color(d.V1)
		r.V0, r.V1 = d.V0, d.V1
		if horizontal {
			if reverse {
				r.X = p1
				r.Width = p0 - p1
			} else {
				r.X = p0
				r.Width = p1 - p0
			}
			r.Y = crossOffset
			r.Height = thickness
		} else {
			if reverse {
				r.Y = p0
				r.Height = p1 - p0
			} else {
				r.Y = p1
				r.Height = p0 - p1
			}
			r.X = crossOffset
			r.Width = thickness
		}
		out = append(out, r)
	}
	return out
}

// titlePosition mirrors @nivo/bullet BulletItem's title placement.
func titlePosition(props BulletProps, horizontal bool, valueSize, bandThickness float64) (x, y float64) {
	if horizontal {
		if props.TitlePosition == "after" {
			x = valueSize + props.TitleOffsetX
		} else {
			x = props.TitleOffsetX
		}
		y = bandThickness/2 + props.TitleOffsetY
	} else {
		x = bandThickness/2 + props.TitleOffsetX
		if props.TitlePosition == "after" {
			y = valueSize + props.TitleOffsetY
		} else {
			y = props.TitleOffsetY
		}
	}
	return x, y
}

// parseColorSpec strips a "seq:" prefix and returns the scheme/interpolator id.
func parseColorSpec(spec string) string {
	if strings.HasPrefix(spec, "seq:") {
		return strings.TrimPrefix(spec, "seq:")
	}
	return spec
}

// sequentialColor builds a value→color function over [minV,maxV] for the given
// sequential scheme/interpolator id.
func sequentialColor(scheme string, minV, maxV float64) func(float64) string {
	return colors.GetSequentialColorScale(
		colors.SequentialColorScaleConfig{Type: "sequential", Scheme: scheme},
		colors.SequentialColorScaleValues{Min: minV, Max: maxV},
	)
}

func (p BulletProps) minValueFor(all []float64) float64 {
	if p.MinValue != nil {
		return *p.MinValue
	}
	return minFloat(all)
}

func (p BulletProps) maxValueFor(all []float64) float64 {
	if p.MaxValue != nil {
		return *p.MaxValue
	}
	return maxFloat(all)
}

func containsFloat(s []float64, v float64) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func minFloat(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	m := s[0]
	for _, v := range s[1:] {
		m = math.Min(m, v)
	}
	return m
}

func maxFloat(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	m := s[0]
	for _, v := range s[1:] {
		m = math.Max(m, v)
	}
	return m
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
