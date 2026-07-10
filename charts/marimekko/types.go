// Package marimekko mirrors @nivo/marimekko: a variable-width stacked bar
// chart. The primary dimension (bar thickness) is value-driven — each datum's
// bar width is proportional to its value — while the secondary dimension is the
// stacked category, computed with the shared d3-shape Stack generator (offset
// none / expand / … ). In the default vertical layout bars run left→right with
// widths ∝ value and dimension segments stacked up each bar.
//
// It reuses internal/d3/shape (Stack), charts/scales (linear range
// constructors), charts/axes (grid + axes), charts/colors (ordinal color per
// dimension), charts/core (SvgWrapper), charts/legends and charts/theming.
//
// SVG only, static render with optional client-side hover tooltips
// per segment (charts/interact) via Interactive.
package marimekko

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// MarimekkoDatum is one column: an index id, the value that drives the bar
// thickness, and a value per stacked dimension key.
type MarimekkoDatum struct {
	ID         string
	Value      float64
	Dimensions map[string]float64
}

// MarimekkoDimension is one stacked category: an id (color/legend) and the key
// read from each datum's Dimensions map.
type MarimekkoDimension struct {
	ID  string
	Key string
}

// MarimekkoLayout is the layout orientation. Mirrors @nivo/marimekko layout.
type MarimekkoLayout string

const (
	MarimekkoLayoutVertical   MarimekkoLayout = "vertical"
	MarimekkoLayoutHorizontal MarimekkoLayout = "horizontal"
)

// OffsetType selects the stack offset. Mirrors @nivo/marimekko offset.
type OffsetType string

const (
	OffsetNone       OffsetType = "none"
	OffsetExpand     OffsetType = "expand"
	OffsetDiverging  OffsetType = "diverging"
	OffsetSilhouette OffsetType = "silhouette"
	OffsetWiggle     OffsetType = "wiggle"
)

// ComputedBar is one positioned, colored segment rect. Mirrors @nivo/marimekko
// ComputedBarDatum.
type ComputedBar struct {
	ID             string
	Index          string
	DimensionID    string
	Value          float64
	FormattedValue string
	X              float64
	Y              float64
	Width          float64
	Height         float64
	Color          string
}

// MarimekkoLayerId enumerates the render layers. Mirrors @nivo/marimekko.
type MarimekkoLayerId string

const (
	MarimekkoLayerGrid    MarimekkoLayerId = "grid"
	MarimekkoLayerAxes    MarimekkoLayerId = "axes"
	MarimekkoLayerBars    MarimekkoLayerId = "bars"
	MarimekkoLayerLegends MarimekkoLayerId = "legends"
)

// DefaultLayers mirrors @nivo/marimekko defaultProps.layers.
var DefaultLayers = []MarimekkoLayerId{
	MarimekkoLayerGrid, MarimekkoLayerAxes, MarimekkoLayerBars, MarimekkoLayerLegends,
}

// MarimekkoProps mirrors @nivo/marimekko MarimekkoSvgProps (the supported
// subset). Fields left zero fall back to Defaults via applyDefaults.
type MarimekkoProps struct {
	// Data is the set of columns; each datum's Value drives its bar thickness.
	Data []MarimekkoDatum
	// Dimensions are the stacked categories, drawn as segments within each bar.
	Dimensions []MarimekkoDimension

	// Width and Height are the outer SVG dimensions in pixels; the inner plot
	// area is these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the inner plot area (for axes and
	// legends), in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Layout is the bar orientation, "vertical" (default) or "horizontal".
	Layout MarimekkoLayout
	// Offset is the d3-shape stack offset applied to the dimension segments:
	// none (default), expand, diverging, silhouette or wiggle.
	Offset OffsetType
	// OuterPadding is the gap in pixels before the first and after the last bar
	// along the thickness axis. Default 0.
	OuterPadding float64
	// InnerPadding is the gap in pixels between adjacent bars. Default 3.
	InnerPadding float64

	// Colors is the ordinal color scale used to color each dimension. Default
	// is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// BorderWidth is the segment rect stroke width in pixels. Default 0 (no
	// border).
	BorderWidth float64
	// BorderColor resolves the segment rect stroke color, optionally inherited
	// from the segment fill.
	BorderColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	// Interactive enables per-segment client-side hover tooltips (charts/interact).
	Interactive bool

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// segment rect. Default false: the rendered SVG is byte-identical to the
	// un-animated output. MotionStagger delays each successive segment by that
	// many seconds (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

	EnableGridX *bool // nil → false (nivo default)
	EnableGridY *bool // nil → true (nivo default)
	// AxisTop, AxisRight, AxisBottom and AxisLeft configure the four axes; nil
	// disables that side.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Legends configures zero or more legends describing the dimensions.
	Legends []legends.LegendProps

	// Theme overrides the styling theme; nil uses the default theme.
	Theme *theming.Theme
	// Layers is the ordered list of render layers; default draws grid, axes,
	// bars then legends.
	Layers []MarimekkoLayerId

	// Role is the SVG root ARIA role. Default "img".
	Role string
	// AriaLabel sets aria-label on the SVG root.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the SVG root.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the SVG root.
	AriaDescribedBy string
	// Title sets the SVG <title> element.
	Title string
	// Desc sets the SVG <desc> element.
	Desc string
	// IsFocusable sets tabindex/focusable on the SVG root.
	IsFocusable bool
}

// MarimekkoResult is the computed model produced by UseMarimekko.
type MarimekkoResult struct {
	Bars       []ComputedBar
	LegendData []legends.Datum
	// ThicknessMax is the summed value across all columns (the thickness-axis
	// domain max); StackMax is the largest stacked dimension total.
	ThicknessMax float64
	StackMax     float64
	// ThicknessScale maps [0, ThicknessMax] onto the thickness axis (x for
	// vertical); StackScale maps [0, StackMax] onto the stack axis (y for
	// vertical). Used by the grid + axes layers.
	ThicknessScale scales.Scale
	StackScale     scales.Scale
}

// GridXEnabled resolves EnableGridX (nil → false, nivo default).
func (p MarimekkoProps) GridXEnabled() bool { return p.EnableGridX != nil && *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p MarimekkoProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }
