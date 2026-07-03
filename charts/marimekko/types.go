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
// v3 scope: SVG only, static render with optional client-side hover tooltips
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
	Data       []MarimekkoDatum
	Dimensions []MarimekkoDimension

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Layout       MarimekkoLayout
	Offset       OffsetType
	OuterPadding float64
	InnerPadding float64

	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
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

	EnableGridX bool
	EnableGridY bool
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []MarimekkoLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
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
