// Package scatterplot mirrors @nivo/scatterplot: series of {x,y} nodes placed
// on two continuous (or time) scales and drawn as dots. It reuses
// charts/scales (ComputeXYScalesForSeries), charts/axes (grid + axes),
// charts/colors (ordinal color per serie), charts/core (DotsItem, cartesian
// markers, SvgWrapper), charts/legends and charts/theming.
//
// Hover uses the charts/interact client layer via Interactive, with nivo's
// voronoi-mesh hit-testing backed by the internal/d3/delaunay port. An opt-in
// Canvas backend (Render: theming.EngineCanvas) renders the marks into a
// <canvas> draw-list for large-N datasets.
package scatterplot

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ScatterPlotDatum is one node's raw data: x and y (number | time | string).
// Mirrors @nivo/scatterplot ScatterPlotDatum.
type ScatterPlotDatum struct {
	X any
	Y any
}

// ScatterPlotSerie is one series: an id and its data points. Mirrors
// @nivo/scatterplot ScatterPlotRawSerie.
type ScatterPlotSerie struct {
	ID   string
	Data []ScatterPlotDatum
}

// ComputedNode is one positioned, colored node ready to render. Mirrors
// @nivo/scatterplot ScatterPlotNodeData.
type ComputedNode struct {
	ID         string
	SerieID    string
	Index      int
	SerieIndex int
	X          float64
	Y          float64
	XValue     any
	YValue     any
	FormattedX string
	FormattedY string
	Size       float64
	Color      string
}

// ScatterPlotLayerId enumerates the render layers. Mirrors @nivo/scatterplot
// ScatterPlotLayerId.
type ScatterPlotLayerId string

const (
	ScatterPlotLayerGrid        ScatterPlotLayerId = "grid"
	ScatterPlotLayerAxes        ScatterPlotLayerId = "axes"
	ScatterPlotLayerNodes       ScatterPlotLayerId = "nodes"
	ScatterPlotLayerMarkers     ScatterPlotLayerId = "markers"
	ScatterPlotLayerMesh        ScatterPlotLayerId = "mesh"
	ScatterPlotLayerLegends     ScatterPlotLayerId = "legends"
	ScatterPlotLayerAnnotations ScatterPlotLayerId = "annotations"
)

// DefaultLayers mirrors @nivo/scatterplot svgDefaultProps.layers.
var DefaultLayers = []ScatterPlotLayerId{
	ScatterPlotLayerGrid, ScatterPlotLayerAxes, ScatterPlotLayerNodes,
	ScatterPlotLayerMarkers, ScatterPlotLayerMesh, ScatterPlotLayerLegends, ScatterPlotLayerAnnotations,
}

// ScatterPlotProps mirrors @nivo/scatterplot ScatterPlotSvgProps (the supported
// subset). Fields left zero fall back to Defaults via applyDefaults.
type ScatterPlotProps struct {
	// Data is the set of series to plot; each serie's {x,y} points are drawn as
	// dots sharing that serie's color.
	Data []ScatterPlotSerie

	// Width and Height are the overall SVG dimensions in pixels.
	Width  float64
	Height float64
	// Margin reserves space around the plot area (top/right/bottom/left) in
	// pixels, e.g. for axes and legends.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// XScale and YScale configure the x/y value scales (linear, time, etc.).
	// Both default to a linear scale over [0, auto].
	XScale scales.ScaleSpec
	YScale scales.ScaleSpec
	// XFormat and YFormat are d3-format specs for the x/y values; empty → %g.
	XFormat string
	YFormat string

	// NodeSize is the diameter of each node dot, in pixels. Default 9.
	NodeSize float64
	// Colors is the ordinal color scale mapping each serie to a color. Default
	// is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig

	// Interactive enables the client-side hover layer (charts/interact): each
	// node emits a data-tc-tooltip the script shows on hover. Default false
	// keeps the static render (and goldens) unchanged.
	Interactive bool

	// UseMesh routes hover through an accurate voronoi mesh (charts/interact,
	// backed by internal/d3/delaunay) rather than per-node tooltips: the whole
	// plot area resolves to the nearest node. Active only when Interactive.
	UseMesh bool
	// DebugMesh draws the voronoi cells as a faint guide when UseMesh is on.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds mesh hit-testing to this pixel distance.
	DetectionRadius float64

	EnableGridX *bool // nil → true (nivo default)
	EnableGridY *bool // nil → true (nivo default)
	// GridXValues and GridYValues override the tick positions of the x/y grid
	// lines; nil lets the scale choose them.
	GridXValues []any
	GridYValues []any
	// AxisTop, AxisRight, AxisBottom and AxisLeft configure the four axes; a nil
	// pointer hides that axis. (nivo default shows bottom and left.)
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Markers draws reference lines/regions in data space over the plot.
	Markers []core.CartesianMarker
	// Legends configures the chart legends; empty means no legend.
	Legends []legends.LegendProps

	// Theme overrides the styling theme; nil uses the default theme.
	Theme *theming.Theme
	// Layers selects which layers are rendered and in what order. Defaults to
	// DefaultLayers.
	Layers []ScatterPlotLayerId

	// Render selects the backend: the zero value (or theming.EngineSVG) renders
	// SVG as always; theming.EngineCanvas draws the nodes into a <canvas>
	// draw-list (charts/canvas) with grid/axes/legends kept as SVG panes and the
	// hover mesh as a transparent overlay. Default SVG keeps every golden stable.
	Render theming.Engine
	// ChartID gives the <canvas> element a stable id (Canvas engine only);
	// defaults to "tc-scatterplot" when empty. Set it when embedding several
	// Canvas scatterplots on one page.
	ChartID string

	// Role is the SVG root's ARIA role. Default "img".
	Role string
	// AriaLabel, AriaLabelledBy and AriaDescribedBy set the matching ARIA
	// attributes on the SVG root for accessibility.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title> and <desc> elements for accessibility.
	Title string
	Desc  string
	// IsFocusable, when true, makes the SVG root keyboard-focusable.
	IsFocusable bool

	// Animate, when true, emits a SMIL enter animation on each node scaling its
	// radius from 0 to its final value (600ms). MotionStagger delays successive
	// nodes by that many seconds (0 ⇒ all nodes enter together). Defaults off,
	// so static output is unchanged. Mirrors @nivo's dot enter transition.
	Animate       bool
	MotionStagger float64
}

// ScatterPlotResult is the computed model produced by UseScatterPlot.
type ScatterPlotResult struct {
	Nodes      []ComputedNode
	XScale     scales.Scale
	YScale     scales.Scale
	LegendData []legends.Datum
}

// GridXEnabled resolves EnableGridX (nil → true, nivo default).
func (p ScatterPlotProps) GridXEnabled() bool { return p.EnableGridX == nil || *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p ScatterPlotProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }
