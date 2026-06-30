// Package line mirrors @nivo/line: LineSeries, ComputedDatum, ComputedSeries,
// Point, SliceData, LineLayerId, the UseLine orchestrator (line/area
// generators, points, slices), and the Line/Lines/LinesItem/Areas/Points/
// Slices/SlicesItem/Mesh/PointTooltip/SliceTooltip templ components.
//
// v1 renders server-side SVG only (no Canvas). Interactivity (point/slice/
// mesh modes) is routed through charts/htmx; crosshair is rendered server-side
// from the current hover state.
package line

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/tooltip"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// AllowedValue is the union of x/y datum values: number | string | time.Time
// | nil. Modeled as any for Go ergonomics.
type AllowedValue any

// LinePointData is one raw point in a series: x and y AllowedValues.
type LinePointData struct {
	X AllowedValue
	Y AllowedValue
}

// LineSeries is one line series: id + data points. Mirrors @nivo/line
// LineSeries.
type LineSeries struct {
	ID   string
	Data []LinePointData
}

// ComputedDatum is one raw point with its computed pixel position. Mirrors
// @nivo/line ComputedDatum.
type ComputedDatum struct {
	Data     LinePointData
	Position struct{ X, Y float64 }
}

// ComputedSeries is a LineSeries with computed positions + resolved color.
// Mirrors @nivo/line ComputedSeries.
type ComputedSeries struct {
	ID    string
	Data  []ComputedDatum
	Color string
}

// PointDatum is the formatted raw data carried on a Point.
type PointDatum struct {
	X          AllowedValue
	Y          AllowedValue
	XFormatted string
	YFormatted string
}

// Point is one rendered point (dot) on a line. Mirrors @nivo/line Point.
type Point struct {
	ID            string
	IndexInSeries int
	AbsIndex      int
	SeriesIndex   int
	SeriesID      string
	SeriesColor   string
	X, Y          float64
	Color         string
	BorderColor   string
	Data          PointDatum
}

// SliceData is one hover slice (x or y axis). Mirrors @nivo/line SliceData.
type SliceData struct {
	ID     string
	X0, X  float64
	Y0, Y  float64
	Width  float64
	Height float64
	Points []Point
}

// PointColorContext is the input to pointColor inherited-color generators.
type PointColorContext struct {
	Series ComputedSeries
	Point  PointColorContextPoint
}

// PointColorContextPoint is the point subset passed to pointColor (omits
// color/borderColor to avoid cycles).
type PointColorContextPoint struct {
	ID            string
	IndexInSeries int
	AbsIndex      int
	SeriesIndex   int
	SeriesID      string
	SeriesColor   string
	X, Y          float64
	Data          PointDatum
}

// LineLayerId enumerates the built-in line layers.
type LineLayerId string

const (
	LineLayerGrid      LineLayerId = "grid"
	LineLayerMarkers   LineLayerId = "markers"
	LineLayerAxes      LineLayerId = "axes"
	LineLayerAreas     LineLayerId = "areas"
	LineLayerCrosshair LineLayerId = "crosshair"
	LineLayerLines     LineLayerId = "lines"
	LineLayerPoints    LineLayerId = "points"
	LineLayerSlices    LineLayerId = "slices"
	LineLayerMesh      LineLayerId = "mesh"
	LineLayerLegends   LineLayerId = "legends"
)

// DefaultLayers is the default line layer order. Mirrors
// commonDefaultProps.layers.
var DefaultLayers = []LineLayerId{
	LineLayerGrid, LineLayerMarkers, LineLayerAxes, LineLayerAreas,
	LineLayerCrosshair, LineLayerLines, LineLayerPoints, LineLayerSlices,
	LineLayerMesh, LineLayerLegends,
}

// EnableSlices is "x" | "y" | false. Empty string = false.
type EnableSlices string

const (
	EnableSlicesX     EnableSlices = "x"
	EnableSlicesY     EnableSlices = "y"
	EnableSlicesFalse EnableSlices = ""
)

// LineGenerator is a func that produces an SVG path-data string from a slice
// of {x,y} points. Mirrors d3-shape's Line generator signature.
type LineGenerator func(points []PointXY) string

// AreaGenerator is a func that produces an SVG area path from points.
type AreaGenerator func(points []PointXY) string

// PointXY is a simple {x, y} pair used by the line/area generators.
type PointXY struct {
	X, Y float64
}

// LineProps mirrors @nivo/line CommonLineProps + LineSvgExtraProps (the SVG
// subset). Fields with zero values fall back to Defaults via UseLine/Bar.
type LineProps struct {
	// Data.
	Data []LineSeries

	// Geometry.
	Width  float64
	Height float64
	Margin core.Margin

	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool

	// Scales.
	XScale  scales.ScaleSpec
	YScale  scales.ScaleLinearSpec
	XFormat core.ValueFormat[any]
	YFormat core.ValueFormat[any]

	// Line.
	Curve     core.CurveFactoryId
	LineWidth float64

	// Area.
	EnableArea        bool
	AreaBaselineValue float64
	AreaOpacity       float64
	AreaBlendMode     core.CssMixBlendMode

	// Points.
	EnablePoints      bool
	PointSize         float64
	PointColor        colors.InheritedColorConfig
	PointBorderWidth  float64
	PointBorderColor  colors.InheritedColorConfig
	EnablePointLabel  bool
	PointLabel        core.PropertyAccessor[Point, string]
	PointLabelYOffset float64

	// Colors.
	Colors colors.OrdinalColorScaleConfig

	// Grid.
	EnableGridX bool
	GridXValues []any
	EnableGridY bool
	GridYValues []any

	// Axes.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Markers.
	Markers []core.CartesianMarker

	// Legends.
	Legends []legends.LegendProps

	// Interactivity.
	IsInteractive bool
	// ClientHover routes mesh/slice hover through the client interactivity
	// layer (charts/interact) instead of the htmx server round-trip: the mesh
	// emits a data-tc-mesh points array (nearest-point + crosshair handled in
	// the browser) and slices emit data-tc-tooltip. Resolves the per-mousemove
	// round-trip the v1 NOTES.md flagged. Default false keeps the htmx path.
	ClientHover          bool
	UseMesh              bool
	EnableSlices         EnableSlices
	DebugSlices          bool
	EnableCrosshair      bool
	CrosshairType        tooltip.CrosshairType
	EnableTouchCrosshair bool
	DebugMesh            bool
	InitialHiddenIDs     []string

	// Defs / fill.
	Defs []core.Def
	Fill []core.DefRule

	// Layers.
	Layers []LineLayerId

	// Motion.
	core.MotionProps

	// A11y.
	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	IsFocusable     bool

	// Theme.
	Theme *theming.Theme

	// HTMX chart instance ID; when non-empty, points/slices emit hx-* attrs.
	ChartID string

	// HoverX/HoverY, when set by the htmx hover endpoint (mesh mode), position
	// the crosshair lines. HasHover gates whether the crosshair renders, so a
	// legitimate hover at the origin (0,0) is not mistaken for "no hover".
	// Cleared on mouseleave.
	HoverX   float64
	HoverY   float64
	HasHover bool
}

// LineSvgProps is an alias of LineProps (nivo splits common/svg; v1 unifies).
type LineSvgProps = LineProps

// LineResult is the output of UseLine: the computed series, points, slices,
// generators, scales, and legend data — everything the templ components need.
type LineResult struct {
	Series        []ComputedSeries
	Points        []Point
	Slices        []SliceData
	LegendData    []LegendDatum
	XScale        scales.Scale
	YScale        scales.Scale
	LineGenerator LineGenerator
	AreaGenerator AreaGenerator
	GetColor      func(LineSeries) string
	HiddenIDs     []string
}

// LegendDatum is one legend entry (id/label/color/hidden). Mirrors nivo's
// line legendData.
type LegendDatum struct {
	ID     string
	Label  string
	Color  string
	Hidden bool
}

// pointXYFromPosition converts a ComputedDatum to a PointXY for the line/area
// generators.
func pointXYFromPosition(d ComputedDatum) PointXY {
	return PointXY{X: d.Position.X, Y: d.Position.Y}
}

// Compile-time check that shape import is used.
var _ = d3shape.CurveLinear
