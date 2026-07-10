// Package radar mirrors @nivo/radar: a polar line/area chart plotting one value
// per key around a set of indices arranged on a circle. It reuses charts/scales
// (linear radius scale), charts/colors (ordinal color per key), charts/core
// (DotsItem, SvgWrapper, curve interpolation via internal/d3/shape with the
// linearClosed curve), charts/legends, and charts/theming. The polar grid
// (axis rays, concentric levels, index labels) is computed here following
// @nivo/radar's RadarGrid geometry.
//
// Angle convention (matching @nivo/radar): index i sits at angle
// rotation + i·angleStep, projected to cartesian via positionFromAngle with a
// -π/2 offset so index 0 is at the top and indices proceed clockwise.
//
// SVG only, static render. The interactive "slices" layer (hover
// tooltip) arrives with the charts/interact client-side layer, so it is omitted here
// even though it is in nivo's default layer list.
package radar

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// GridShape is the radar grid level shape: concentric circles or polygons.
// Mirrors @nivo/radar GridShape.
type GridShape string

const (
	GridShapeCircular GridShape = "circular"
	GridShapeLinear   GridShape = "linear"
)

// RadarLayerId enumerates the render layers. Mirrors @nivo/radar RadarLayerId.
// The "slices" layer is interactive-only and not rendered in the static
// pipeline.
type RadarLayerId string

const (
	RadarLayerGrid    RadarLayerId = "grid"
	RadarLayerLayers  RadarLayerId = "layers"
	RadarLayerSlices  RadarLayerId = "slices"
	RadarLayerDots    RadarLayerId = "dots"
	RadarLayerLegends RadarLayerId = "legends"
)

// DefaultLayers mirrors @nivo/radar svgDefaultProps.layers.
var DefaultLayers = []RadarLayerId{
	RadarLayerGrid, RadarLayerLayers, RadarLayerSlices, RadarLayerDots, RadarLayerLegends,
}

// RadarProps mirrors @nivo/radar RadarSvgProps (the supported subset). Data is
// a list of objects (rows); each row carries one index value (read via IndexBy)
// and one numeric value per key in Keys. Fields left zero fall back to Defaults
// via applyDefaults.
type RadarProps struct {
	// Data is the list of rows; each row holds one index value and one numeric
	// value per key.
	Data []map[string]any
	// Keys are the value keys, each drawn as its own closed polygon.
	Keys []string
	// IndexBy is the row field used as the index (one axis ray per index).
	IndexBy string

	// Width and Height are the outer SVG dimensions in pixels; the inner plot
	// area is these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the inner plot area (for labels and
	// legends), in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// MaxValue is the radius-scale domain max. nil → "auto" (max over all
	// key values). Modeled as a pointer so the default-auto behavior survives
	// the float64 zero value.
	MaxValue *float64
	// Rotation rotates the whole radar by this many degrees (clockwise).
	Rotation    float64
	ValueFormat string // d3-format spec for dot labels; empty → %g
	// Curve is the d3-shape curve factory used to connect a series' points.
	// Default linearClosed.
	Curve core.CurveFactoryId

	// BorderWidth is the series polygon stroke width in pixels. Default 2.
	BorderWidth float64
	// BorderColor resolves the series polygon stroke color; by default inherits
	// the series color.
	BorderColor colors.InheritedColorConfig

	// GridLevels is the number of concentric grid levels. Default 5.
	GridLevels int
	// GridShape is the grid level shape, "circular" (default) or "linear"
	// (polygons).
	GridShape GridShape
	// GridLabelOffset is the distance in pixels between the outermost grid level
	// and the index labels. Default 16.
	GridLabelOffset float64

	// EnableDots gates the per-point dots. nil → true (nivo default).
	EnableDots *bool
	// DotSize is the dot diameter in pixels. Default 6.
	DotSize float64
	// DotColor resolves the dot fill color; by default inherits the series color.
	DotColor colors.InheritedColorConfig
	// DotBorderWidth is the dot stroke width in pixels. Default 0.
	DotBorderWidth float64
	// DotBorderColor resolves the dot stroke color; by default inherits the
	// series color.
	DotBorderColor colors.InheritedColorConfig
	EnableDotLabel *bool // nil → false (nivo default)
	// DotLabelYOffset is the vertical offset in pixels applied to dot labels.
	// Default -12.
	DotLabelYOffset float64

	// Colors is the ordinal color scale used to color each key. Default is the
	// "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// FillOpacity is the series polygon fill opacity (0..1). Default 0.25.
	FillOpacity float64

	// Interactive enables the client-side hover layer (charts/interact): each
	// dot emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Legends configures zero or more legends describing the keys.
	Legends []legends.LegendProps

	// Theme overrides the styling theme; nil uses the default theme.
	Theme *theming.Theme
	// Layers is the ordered list of render layers; default draws grid, layers,
	// slices, dots then legends.
	Layers []RadarLayerId

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

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// series polygon (600ms). MotionStagger delays successive series by that
	// many seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/radar's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// ComputedPoint is one projected key/index point: its cartesian position
// (relative to the radar center), value, and resolved colors.
type ComputedPoint struct {
	Key            string
	Index          string
	Value          float64
	FormattedValue string
	X              float64
	Y              float64
	Color          string
	BorderColor    string
}

// ComputedSerie is one key's closed polygon: its line path (already projected
// + closed via the configured curve), fill/stroke colors, and its points.
type ComputedSerie struct {
	Key    string
	Path   string
	Color  string
	Stroke string
	Points []ComputedPoint
}

// RadarResult is the computed model produced by UseRadar.
type RadarResult struct {
	Indices     []string
	Keys        []string
	ColorByKey  map[string]string
	Series      []ComputedSerie
	Points      []ComputedPoint
	RadiusScale scales.Scale
	Radius      float64
	CenterX     float64
	CenterY     float64
	Rotation    float64 // radians
	AngleStep   float64 // radians
	MaxValue    float64
	LegendData  []legends.Datum
}

// DotsEnabled resolves EnableDots: unset (nil) means true (nivo default).
func (p RadarProps) DotsEnabled() bool {
	return p.EnableDots == nil || *p.EnableDots
}

// DotLabelEnabled resolves EnableDotLabel (nil → false, nivo default).
func (p RadarProps) DotLabelEnabled() bool {
	return p.EnableDotLabel != nil && *p.EnableDotLabel
}
