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
// v2 scope: SVG only, static render. The interactive "slices" layer (hover
// tooltip) arrives with the Phase 5 client-side layer, so it is omitted here
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
// The "slices" layer is interactive-only and not rendered in the v2 static
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
	Data    []map[string]any
	Keys    []string
	IndexBy string

	Width  float64
	Height float64
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
	Curve       core.CurveFactoryId

	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	GridLevels      int
	GridShape       GridShape
	GridLabelOffset float64

	// EnableDots gates the per-point dots. nil → true (nivo default).
	EnableDots      *bool
	DotSize         float64
	DotColor        colors.InheritedColorConfig
	DotBorderWidth  float64
	DotBorderColor  colors.InheritedColorConfig
	EnableDotLabel  bool
	DotLabelYOffset float64

	Colors      colors.OrdinalColorScaleConfig
	FillOpacity float64

	// Interactive enables the client-side hover layer (charts/interact): each
	// dot emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []RadarLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

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

// BoolPtr returns a pointer to b — a helper for *bool props like EnableDots.
func BoolPtr(b bool) *bool { return &b }

// FloatPtr returns a pointer to f — a helper for *float64 props like MaxValue.
func FloatPtr(f float64) *float64 { return &f }
