// Package line mirrors @nivo/line: LineSeries, ComputedDatum, ComputedSeries,
// Point, SliceData, LineLayerId, the UseLine orchestrator (line/area
// generators, points, slices), and the Line/Lines/LinesItem/Areas/Points/
// Slices/SlicesItem/Mesh/PointTooltip/SliceTooltip templ components.
//
// Renders server-side SVG only (no Canvas). Interactivity (point/slice/
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
	// Data is the set of line series to plot; each LineSeries carries an id and
	// its ordered list of {x,y} points.
	Data []LineSeries

	// Width is the total outer width of the chart in pixels (plot area plus
	// margins).
	Width float64
	// Height is the total outer height of the chart in pixels (plot area plus
	// margins).
	Height float64
	// Margin is the space reserved around the plot area for axes, tick labels,
	// and legends; the inner plotted area is Width/Height minus these margins.
	Margin core.Margin

	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool

	// Scales. YScale accepts any continuous scale spec — a
	// scales.ScaleLinearSpec (the default), scales.ScaleLogSpec, or
	// scales.ScaleSymlogSpec. A nil value falls back to the linear default.
	XScale scales.ScaleSpec
	YScale scales.ScaleSpec
	// XFormat is an optional formatter applied to x values for tick labels and
	// tooltips. Nil leaves values unformatted (stringified as-is).
	XFormat core.ValueFormat[any]
	// YFormat is an optional formatter applied to y values for tick labels and
	// tooltips. Nil leaves values unformatted (stringified as-is).
	YFormat core.ValueFormat[any]

	// Curve is the interpolation used to connect points into a path (e.g.
	// linear, monotoneX, step, cardinal, basis). Default core.CurveLinear.
	Curve core.CurveFactoryId
	// LineWidth is the stroke width of each series line in pixels. Default 2.
	LineWidth float64

	// EnableArea toggles the filled area beneath each line. nil → false (nivo
	// default).
	EnableArea *bool // nil → false (nivo default)
	// AreaBaselineValue is the y data value the area fill is anchored to.
	// Default 0.
	AreaBaselineValue float64
	// AreaOpacity is the fill opacity of the area, in the range 0–1. Only
	// applies when EnableArea is true. Default 0.2.
	AreaOpacity float64
	// AreaBlendMode is the CSS mix-blend-mode applied to overlapping area fills.
	// Default core.MixBlendNormal.
	AreaBlendMode core.CssMixBlendMode

	// Points.
	EnablePoints *bool // nil → true (nivo default)
	// PointSize is the diameter of each point circle in pixels. Default 6.
	PointSize float64
	// PointColor determines each point's fill color; by default it inherits the
	// series color ("series.color").
	PointColor colors.InheritedColorConfig
	// PointBorderWidth is the stroke width of the point border in pixels.
	// Default 0 (no visible border).
	PointBorderWidth float64
	// PointBorderColor determines each point's border color; by default it
	// resolves to the theme background color.
	PointBorderColor colors.InheritedColorConfig
	// EnablePointLabel toggles rendering a text label next to each point.
	// nil → false (nivo default).
	EnablePointLabel *bool // nil → false (nivo default)
	// PointLabel selects the value used for point labels; accepts a datum
	// property path or a function returning the label. Default "data.yFormatted".
	PointLabel core.PropertyAccessor[Point, string]
	// PointLabelYOffset is the vertical offset in pixels of the point label from
	// the point shape. Default 0.
	PointLabelYOffset float64

	// Colors configures the ordinal color scale mapping series to colors.
	// Default is the "nivo" categorical scheme.
	Colors colors.OrdinalColorScaleConfig

	// Grid. nil → true (nivo default) for both.
	EnableGridX *bool
	// GridXValues restricts the x grid lines to these specific x values; nil
	// uses the x scale's default tick positions.
	GridXValues []any
	// EnableGridY toggles horizontal grid lines. nil → true (nivo default).
	EnableGridY *bool
	// GridYValues restricts the y grid lines to these specific y values; nil
	// uses the y scale's default tick positions.
	GridYValues []any

	// AxisTop configures the top axis; nil hides it (no default axis).
	AxisTop *axes.AxisProps
	// AxisRight configures the right axis; nil hides it (no default axis).
	AxisRight *axes.AxisProps
	// AxisBottom configures the bottom axis; when nil, the render path
	// substitutes the default axis (tickSize/tickPadding 5).
	AxisBottom *axes.AxisProps
	// AxisLeft configures the left axis; when nil, the render path substitutes
	// the default axis (tickSize/tickPadding 5).
	AxisLeft *axes.AxisProps

	// Markers are reference lines/annotations drawn at fixed x or y values.
	Markers []core.CartesianMarker

	// Legends are the chart's legends; empty renders none.
	Legends []legends.LegendProps

	// Interactive enables hover interactivity (points/slices/mesh). Default
	// true.
	Interactive bool
	// ClientHover is retained for compatibility: mesh/slice hover now routes
	// through the client interactivity layer (charts/interact) by default, so
	// setting this is a no-op. The client mesh emits a data-tc-mesh points array
	// (nearest-point + crosshair handled in the browser) and slices emit
	// data-tc-tooltip — no per-mousemove server round-trip. See ServerHover to
	// opt back into the legacy htmx path.
	ClientHover bool
	// ServerHover opts back into the legacy htmx per-mousemove server round-trip
	// for mesh/slice hover (a full server SVG re-render on every throttled
	// mousemove), for the genuinely-JS-limited case where htmx is present but the
	// charts/interact client script is not. Default false: the per-mousemove
	// fallback is retired and the client path is the default. Requires ChartID.
	ServerHover bool
	// UseMesh enables a voronoi mesh for nearest-point mouse detection; requires
	// EnableSlices to be disabled. Default false.
	UseMesh bool
	// EnableSlices groups points into hover slices along the "x" or "y" axis
	// (automatically disabling the mesh). Empty string = disabled (default).
	EnableSlices EnableSlices
	// DebugSlices renders the slice detection rectangles for debugging. Default
	// false.
	DebugSlices bool
	// EnableCrosshair toggles the crosshair lines shown on hover. nil → true
	// (nivo default).
	EnableCrosshair *bool // nil → true (nivo default)
	// CrosshairType selects which crosshair lines to draw relative to the cursor
	// (e.g. bottom-left, cross, x, y); forced to the slice axis when slices are
	// enabled. Default tooltip.CrosshairTypeBottomLeft.
	CrosshairType tooltip.CrosshairType
	// EnableTouchCrosshair allows the crosshair to be dragged on touch screens.
	// nil → false (nivo default).
	EnableTouchCrosshair *bool // nil → false (nivo default)
	// DebugMesh renders the voronoi mesh cells for debugging. Default false.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds voronoi-mesh hit-testing to this pixel
	// distance from the cursor (nivo's mesh detectionRadius). Client-hover only.
	DetectionRadius float64
	// InitialHiddenIDs lists series ids that start hidden (excluded from the plot
	// but still shown as toggleable legend entries).
	InitialHiddenIDs []string

	// Defs are reusable SVG definitions (gradients/patterns) referenced by Fill.
	Defs []core.Def
	// Fill maps Defs to series via matching rules, letting lines/areas use a
	// gradient or pattern fill instead of a flat color.
	Fill []core.DefRule

	// Layers sets the render order of the built-in (and any custom) line layers.
	// Defaults to DefaultLayers.
	Layers []LineLayerId

	// MotionProps carries the shared animation settings (Animate, MotionConfig).
	core.MotionProps

	// Role is the ARIA role of the root SVG element. Default "img".
	Role string
	// AriaLabel sets aria-label on the root SVG element.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the root SVG element.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the root SVG element.
	AriaDescribedBy string
	// Title is the SVG <title> element text, used for accessibility.
	Title string
	// Desc is the SVG <desc> element text, used for accessibility.
	Desc string
	// IsFocusable makes the root SVG element (and each point) focusable for
	// keyboard navigation. Default false.
	IsFocusable bool

	// Theme overrides the chart theme (colors, fonts, axis styling); nil uses
	// the default theme.
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

// LineSvgProps is an alias of LineProps (nivo splits common/svg; this unifies them).
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

// AreaEnabled resolves EnableArea (nil → false, nivo default).
func (p LineProps) AreaEnabled() bool { return p.EnableArea != nil && *p.EnableArea }

// PointsEnabled resolves EnablePoints (nil → true, nivo default).
func (p LineProps) PointsEnabled() bool { return p.EnablePoints == nil || *p.EnablePoints }

// PointLabelEnabled resolves EnablePointLabel (nil → false, nivo default).
func (p LineProps) PointLabelEnabled() bool { return p.EnablePointLabel != nil && *p.EnablePointLabel }

// GridXEnabled resolves EnableGridX (nil → true, nivo default).
func (p LineProps) GridXEnabled() bool { return p.EnableGridX == nil || *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p LineProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }

// CrosshairEnabled resolves EnableCrosshair (nil → true, nivo default).
func (p LineProps) CrosshairEnabled() bool { return p.EnableCrosshair == nil || *p.EnableCrosshair }

// TouchCrosshairEnabled resolves EnableTouchCrosshair (nil → false, nivo default).
func (p LineProps) TouchCrosshairEnabled() bool {
	return p.EnableTouchCrosshair != nil && *p.EnableTouchCrosshair
}
