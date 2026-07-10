// Package bullet mirrors @nivo/bullet: a compact KPI chart of stacked
// qualitative ranges, one or more measure bars, and comparative markers, all on
// a shared value scale. It reuses charts/scales (linear value scale),
// charts/axes (the per-item axis), charts/colors (sequential color scales for
// ranges/measures/markers) and charts/theming.
//
// SVG only, static render. The interactive hover tooltip arrives with
// the charts/interact client layer.
package bullet

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// BulletLayout is horizontal (bars run left→right) or vertical (bottom→top).
type BulletLayout string

const (
	BulletLayoutHorizontal BulletLayout = "horizontal"
	BulletLayoutVertical   BulletLayout = "vertical"
)

// BulletItemDatum is one bullet row. Mirrors @nivo/bullet Datum.
type BulletItemDatum struct {
	ID       string
	Title    string // optional; falls back to ID
	Ranges   []float64
	Measures []float64
	Markers  []float64
}

// ComputedRect is one positioned range/measure rect.
type ComputedRect struct {
	X, Y, Width, Height float64
	Color               string
	V0, V1              float64
}

// ComputedMarker is one positioned marker line. (X1,Y1)→(X2,Y2).
type ComputedMarker struct {
	X1, Y1, X2, Y2 float64
	Color          string
	Value          float64
}

// ComputedBulletItem is one fully laid-out bullet row.
type ComputedBulletItem struct {
	ID       string
	Title    string
	OffsetX  float64 // item group translate
	OffsetY  float64
	Width    float64 // item dimensions (cross-axis = height of the bar band)
	Height   float64
	Ranges   []ComputedRect
	Measures []ComputedRect
	Markers  []ComputedMarker
	Scale    scales.Scale
	AxisX    float64 // axis group translate within the item
	AxisY    float64
	TitleX   float64
	TitleY   float64
}

// BulletProps mirrors @nivo/bullet BulletSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type BulletProps struct {
	// Data holds one entry per bullet row (its ranges, measures and markers).
	Data []BulletItemDatum

	// Width and Height are the total SVG dimensions in pixels; the plot area is
	// these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the plot area (for axes/titles).
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// Interactive enables the client-side hover layer (charts/interact): each
	// range/measure rect emits a data-tc-tooltip. Default false keeps the
	// static render.
	Interactive bool

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// range/measure rect. Default false: the rendered SVG is byte-identical to
	// the un-animated output. MotionStagger delays each successive row by that
	// many seconds (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

	// Layout orients the bars "horizontal" (default) or "vertical".
	Layout BulletLayout
	// Reverse flips the value-scale direction; default false.
	Reverse bool
	// Spacing is the gap in pixels between bullet rows; default 30.
	Spacing float64

	// MinValue/MaxValue bound the per-item value scale. nil → "auto" (derived
	// from the item's own values).
	MinValue *float64
	MaxValue *float64

	AxisPosition string // "before" | "after"

	RangeColors string // e.g. "seq:cool"
	// MeasureColors is the sequential color spec for measure bars; default
	// "seq:red_purple". MarkerColors is the same for markers; default
	// "seq:red_purple".
	MeasureColors string
	MarkerColors  string

	// RangeBorderWidth is the stroke width in pixels of range rects; default 0.
	RangeBorderWidth   float64
	MeasureSize        float64 // ratio of item height
	MeasureBorderWidth float64 // stroke width in pixels of measure rects; default 0
	MarkerSize         float64 // ratio of item height

	TitlePosition string // "before" | "after"
	TitleAlign    string // "start" | "middle" | "end"
	// TitleOffsetX and TitleOffsetY shift the title in pixels from its computed
	// position. TitleRotation rotates the title in degrees.
	TitleOffsetX  float64
	TitleOffsetY  float64
	TitleRotation float64

	// Theme overrides chart styling; nil → theming.DefaultTheme.
	Theme *theming.Theme

	// Role is the root SVG ARIA role; default "img".
	Role string
	// AriaLabel, AriaLabelledBy, and AriaDescribedBy set the corresponding SVG
	// accessibility attributes.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title>/<desc> elements.
	Title string
	Desc  string
	// IsFocusable makes the SVG keyboard-focusable.
	IsFocusable bool
}

// BulletResult is the computed model produced by UseBullet.
type BulletResult struct {
	Items []ComputedBulletItem
}
