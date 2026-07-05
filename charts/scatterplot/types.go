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
	Data []ScatterPlotSerie

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	XScale  scales.ScaleSpec
	YScale  scales.ScaleSpec
	XFormat string // d3-format spec; empty → %g
	YFormat string

	NodeSize float64
	Colors   colors.OrdinalColorScaleConfig

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

	EnableGridX bool
	EnableGridY bool
	GridXValues []any
	GridYValues []any
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	Markers []core.CartesianMarker
	Legends []legends.LegendProps

	Theme  *theming.Theme
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

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

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
