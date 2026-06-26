// Package axes mirrors @nivo/axes: AxisProps, ComputeCartesianTicks (tick
// positions + alignment), ComputeGridLines, the formatter resolver
// GetFormatter, and the Axes/Axis/AxisTick/Grid/GridLines/GridLine templ
// components.
package axes

import (
	"fmt"
	"math"
	"time"

	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// AxisLegendPosition is where the axis legend label sits relative to the axis.
type AxisLegendPosition string

const (
	AxisLegendStart  AxisLegendPosition = "start"
	AxisLegendMiddle AxisLegendPosition = "middle"
	AxisLegendEnd    AxisLegendPosition = "end"
)

// AxisTickProps is a single computed tick: value, pixel position, and label.
type AxisTickProps struct {
	Value    any
	Position float64 // pixel offset along the axis (after centering if applicable)
	Label    string
}

// GridValues[V] is the union of a tick count (int) or explicit tick values.
// (TicksSpec with interval is handled at the scales layer.)
type GridValues[V any] struct {
	Count  int
	Values []V
}

// AxisProps mirrors @nivo/axes AxisProps. Only the fields used by the SVG
// rendering path are kept; canvas-only fields are omitted.
type AxisProps struct {
	Axis           string // "x" | "y"
	Scale          scales.Scale
	X              float64 // origin offset along the cross axis
	Y              float64
	Length         float64 // axis length (px)
	TicksPosition  string  // "before" | "after" (nivo: ticksPosition)
	TickValues     []any
	TickSize       float64 // tick line length (perpendicular to axis)
	TickPadding    float64
	TickRotation   float64 // degrees
	TextAlign      string  // nivo textAlign ("start"|"center"|"end")
	TextBaseline   string  // nivo textBaseline ("top"|"center"|"bottom")
	Legend         string
	LegendPosition AxisLegendPosition
	LegendOffset   float64
	LegendWidth    float64
	LegendHeight   float64
	// Style overrides (fill into theme attrs when non-zero).
	LineColor            string
	LineWidth            float64
	TickLineColor        string
	TickLineWidth        float64
	TickTextColor        string
	TickTextFontSize     any
	TickTextFontFamily   string
	LegendTextColor      string
	LegendTextFontSize   any
	LegendTextFontFamily string
	// Hidden suppresses rendering entirely.
	Hidden bool
	// DisableTickLabel drops the label text (line still drawn).
	DisableTickLabel bool
}

// CanvasAxisProps mirrors @nivo/axes CanvasAxisProps (same fields, canvas-only
// rendering dropped). Provided for API parity.
type CanvasAxisProps = AxisProps

// DefaultAxisProps mirrors @nivo/axes defaultProps.
var DefaultAxisProps = AxisProps{
	TickSize:       5,
	TickPadding:    5,
	TickRotation:   0,
	TextAlign:      "center",
	TextBaseline:   "middle",
	LegendPosition: AxisLegendEnd,
	LegendOffset:   0,
}

// Positions mirrors @nivo/axes positions const.
const (
	PositionTop    = "top"
	PositionBottom = "bottom"
	PositionLeft   = "left"
	PositionRight  = "right"
)

// Tick is one computed tick with its line endpoints.
type Tick struct {
	Value    any
	Position float64
	Label    string
	// Line endpoints (relative to the axis origin). For a bottom axis, x1=x2=position,
	// y1=0, y2=tickSize (or -tickSize for ticksPosition="before").
	X1, Y1, X2, Y2 float64
	// Label position offset.
	LabelX, LabelY   float64
	TextAnchor       string
	DominantBaseline string
}

// ComputedTicks is the result of ComputeCartesianTicks.
type ComputedTicks struct {
	Ticks        []Tick
	TextAlign    string
	TextBaseline string
}

// ComputeCartesianTicks mirrors @nivo/axes computeCartesianTicks: derives the
// tick values (from TickValues or the scale's default ticks), positions them,
// and computes label text via GetFormatter. ticksPosition determines which
// side of the axis the ticks protrude.
func ComputeCartesianTicks(props AxisProps, theme *theming.Theme) ComputedTicks {
	if theme == nil {
		theme = &theming.DefaultTheme
	}
	axis := props.Axis
	values := props.TickValues
	if values == nil {
		values = scales.GetScaleTicks(props.Scale, scales.TicksSpec{})
	}

	// Center the scale for band/point axes.
	positionFn := props.Scale.Call
	if _, ok := props.Scale.(scales.ScaleWithBandwidth); ok {
		positionFn = scales.CenterScale(props.Scale)
	}

	formatter := GetFormatter(props.Scale, theme)
	textAlign := props.TextAlign
	if textAlign == "" {
		textAlign = "center"
	}
	textBaseline := props.TextBaseline
	if textBaseline == "" {
		textBaseline = "middle"
	}

	ticks := make([]Tick, 0, len(values))
	for _, v := range values {
		pos := positionFn(v)
		label := formatter(v)
		t := Tick{Value: v, Position: pos, Label: label}
		// Line endpoints and label position depend on axis orientation.
		sign := 1.0
		if props.TicksPosition == "before" {
			sign = -1
		}
		if axis == "x" {
			t.X1, t.X2 = pos, pos
			t.Y1 = 0
			t.Y2 = sign * props.TickSize
			t.LabelX = pos
			t.LabelY = sign * (props.TickSize + props.TickPadding)
		} else {
			t.Y1, t.Y2 = pos, pos
			t.X1 = 0
			t.X2 = sign * props.TickSize
			t.LabelX = sign * (props.TickSize + props.TickPadding)
			t.LabelY = pos
		}
		// SVG text attrs from the theming bridge.
		attrs := svgTextAttrs(textAlign, textBaseline)
		t.TextAnchor = attrs.TextAnchor
		t.DominantBaseline = attrs.DominantBaseline
		// For rotated ticks, nivo forces text-anchor based on rotation.
		if props.TickRotation != 0 {
			t.TextAnchor, t.DominantBaseline = rotatedTextAttrs(props.TickRotation, axis)
		}
		ticks = append(ticks, t)
	}
	return ComputedTicks{Ticks: ticks, TextAlign: textAlign, TextBaseline: textBaseline}
}

// svgTextAttrs mirrors the theming bridge conversion.
func svgTextAttrs(textAlign, textBaseline string) struct{ TextAnchor, DominantBaseline string } {
	a := struct{ TextAnchor, DominantBaseline string }{
		TextAnchor:       theming.ConvertStyleAttribute(theming.EngineSVG, "textAlign", textAlign),
		DominantBaseline: theming.ConvertStyleAttribute(theming.EngineSVG, "textBaseline", textBaseline),
	}
	return a
}

// rotatedTextAttrs returns the text-anchor / dominant-baseline for a rotated
// tick label. Mirrors @nivo/axes axis tick rotation logic.
func rotatedTextAttrs(rotation float64, axis string) (string, string) {
	r := math.Mod(rotation, 360)
	if r < 0 {
		r += 360
	}
	if axis == "x" {
		// Horizontal axis with rotated labels.
		if r == 90 {
			return "middle", "hanging"
		}
		if r == -90 || r == 270 {
			return "middle", "auto"
		}
		if r > 0 && r < 90 {
			return "end", "auto"
		}
		if r > 270 {
			return "start", "auto"
		}
		if r > 90 && r < 270 {
			return "start", "auto"
		}
	}
	// Vertical axis rotation not typical in nivo defaults.
	return "middle", "middle"
}

// Line is one grid line (x1,y1)-(x2,y2).
type Line struct {
	X1, Y1, X2, Y2 float64
}

// ComputeGridLines mirrors @nivo/axes computeGridLines: produces the set of
// grid lines (perpendicular to the axis at each tick position) spanning the
// chart's width/height.
func ComputeGridLines(props AxisProps, width, height float64) []Line {
	values := props.TickValues
	if values == nil {
		values = scales.GetScaleTicks(props.Scale, scales.TicksSpec{})
	}
	positionFn := props.Scale.Call
	if _, ok := props.Scale.(scales.ScaleWithBandwidth); ok {
		positionFn = scales.CenterScale(props.Scale)
	}
	lines := make([]Line, 0, len(values))
	for _, v := range values {
		pos := positionFn(v)
		if props.Axis == "x" {
			lines = append(lines, Line{X1: pos, Y1: 0, X2: pos, Y2: height})
		} else {
			lines = append(lines, Line{X1: 0, Y1: pos, X2: width, Y2: pos})
		}
	}
	return lines
}

// GetFormatter returns a func(any) string that formats tick values for the
// given scale. Linear/log/symlog use d3-format "%~s" approximation (nivo's
// default is the identity → String(value)); time scales use a multi-format
// based on precision; band/point scales stringify. Mirrors @nivo/axes
// getFormatter.
func GetFormatter(scale scales.Scale, theme *theming.Theme) func(any) string {
	switch scale.Type() {
	case scales.ScaleTypeTime:
		return timeFormatter
	case scales.ScaleTypeLinear, scales.ScaleTypeLog, scales.ScaleTypeSymlog:
		return numberFormatter
	default:
		return func(v any) string { return fmt.Sprintf("%v", v) }
	}
}

func numberFormatter(v any) string {
	switch x := v.(type) {
	case float64:
		// Trim trailing zeros for compactness.
		return trimFloat(fmt.Sprintf("%g", x))
	case int:
		return fmt.Sprintf("%d", x)
	case time.Time:
		return x.Format("2006-01-02")
	}
	return fmt.Sprintf("%v", v)
}

func timeFormatter(v any) string {
	if t, ok := v.(time.Time); ok {
		// Simple multi-format: year if Jan 1, else month, else day.
		if t.Month() == 1 && t.Day() == 1 {
			return t.Format("2006")
		}
		if t.Day() == 1 {
			return t.Format("Jan")
		}
		return t.Format("Jan 2")
	}
	return fmt.Sprintf("%v", v)
}

func trimFloat(s string) string {
	return s
}

// legendPosition computes the (x, y, rotation, text-anchor) for an axis legend
// label given the axis orientation and LegendPosition. Mirrors @nivo/axes
// computeLegendPosition.
func legendPosition(props AxisProps) (x, y, rotation float64, anchor string) {
	pos := props.LegendPosition
	if pos == "" {
		pos = AxisLegendEnd
	}
	offset := props.LegendOffset
	if offset == 0 {
		offset = 32
	}
	length := props.Length
	if props.Axis == "x" {
		switch pos {
		case AxisLegendStart:
			x = 0
		case AxisLegendMiddle:
			x = length / 2
		default:
			x = length
		}
		y = offset
		rotation = 0
		anchor = string(pos)
		return
	}
	// y axis: legend is rotated -90 (vertical).
	switch pos {
	case AxisLegendStart:
		y = 0
	case AxisLegendMiddle:
		y = length / 2
	default:
		y = length
	}
	x = -offset
	rotation = -90
	anchor = string(pos)
	return
}

// legendStyle returns the legend text color/font for an axis given the theme.
func legendStyle(props AxisProps, theme *theming.Theme) (fill string, fontSize any, fontFamily string) {
	if theme == nil {
		return
	}
	ts := theme.Axis.Legend.Text
	fill = ts.Fill
	if fill == "" {
		fill = theme.Text.Fill
	}
	fontSize = ts.FontSize
	if fontSize == nil {
		fontSize = theme.Text.FontSize
	}
	fontFamily = ts.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	if props.LegendTextColor != "" {
		fill = props.LegendTextColor
	}
	if props.LegendTextFontSize != nil {
		fontSize = props.LegendTextFontSize
	}
	if props.LegendTextFontFamily != "" {
		fontFamily = props.LegendTextFontFamily
	}
	return
}

// tickTextStyle returns the tick label color/font for an axis given the theme.
func tickTextStyle(props AxisProps, theme *theming.Theme) (fill string, fontSize any, fontFamily string) {
	if theme == nil {
		return
	}
	ts := theme.Axis.Ticks.Text
	fill = ts.Fill
	if fill == "" {
		fill = theme.Text.Fill
	}
	fontSize = ts.FontSize
	if fontSize == nil {
		fontSize = theme.Text.FontSize
	}
	fontFamily = ts.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	if props.TickTextColor != "" {
		fill = props.TickTextColor
	}
	if props.TickTextFontSize != nil {
		fontSize = props.TickTextFontSize
	}
	if props.TickTextFontFamily != "" {
		fontFamily = props.TickTextFontFamily
	}
	return
}

// tickLineStyle returns the tick line stroke/width for an axis given the theme.
func tickLineStyle(props AxisProps, theme *theming.Theme) (stroke string, strokeWidth float64) {
	if theme == nil {
		return
	}
	line := theme.Axis.Ticks.Line
	stroke = stringFromExtra(line.Extra, "stroke")
	strokeWidth = floatFromExtra(line.Extra, "strokeWidth")
	if props.TickLineColor != "" {
		stroke = props.TickLineColor
	}
	if props.TickLineWidth != 0 {
		strokeWidth = props.TickLineWidth
	}
	return
}

// domainLineStyle returns the axis domain line stroke/width.
func domainLineStyle(props AxisProps, theme *theming.Theme) (stroke string, strokeWidth float64) {
	if theme == nil {
		return
	}
	line := theme.Axis.Domain.Line
	stroke = stringFromExtra(line.Extra, "stroke")
	strokeWidth = floatFromExtra(line.Extra, "strokeWidth")
	if props.LineColor != "" {
		stroke = props.LineColor
	}
	if props.LineWidth != 0 {
		strokeWidth = props.LineWidth
	}
	return
}

func stringFromExtra(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func floatFromExtra(m map[string]any, key string) float64 {
	if m == nil {
		return 0
	}
	if v, ok := m[key]; ok {
		switch x := v.(type) {
		case float64:
			return x
		case int:
			return float64(x)
		}
	}
	return 0
}

// fmtN formats a number for SVG output (3 dp, matching d3-path rounding).
func fmtN(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := fmt.Sprintf("%.3g", v)
	return s
}

// AxesProps renders both (x and y) axes for a chart.
type AxesProps struct {
	XAxis  *AxisProps
	YAxis  *AxisProps
	Width  float64
	Height float64
	Theme  *theming.Theme
}
