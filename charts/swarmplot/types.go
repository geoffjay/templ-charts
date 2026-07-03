// Package swarmplot mirrors @nivo/swarmplot: points grouped along one axis and
// positioned by value along the other, then relaxed with a force simulation so
// they don't overlap. It reuses internal/d3/force (ForceX/ForceY + ForceCollide,
// run for a fixed iteration count), charts/scales (a linear value scale + a band
// group scale), charts/axes (grid + axes), charts/colors (ordinal color per
// group or id), charts/core, charts/interact (voronoi-mesh hover) and charts/theming.
//
// v3 scope: SVG only, static render. The layout is deterministic (phyllotaxis
// seeding + d3's LCG + a fixed SimulationIterations count), so goldens are
// byte-stable. Annotations are deferred.
package swarmplot

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// SwarmPlotDatum is one input point: an id, a group, and a numeric value.
// Mirrors @nivo/swarmplot's default datum (id/group/value accessors).
type SwarmPlotDatum struct {
	ID    string
	Group string
	Value float64
}

// ComputedNode is one positioned, colored node ready to render.
type ComputedNode struct {
	ID             string
	Group          string
	Value          float64
	FormattedValue string
	X, Y           float64
	Size           float64
	Color          string
}

// SwarmPlotLayerId enumerates the render layers. Mirrors @nivo/swarmplot
// SwarmPlotLayerId.
type SwarmPlotLayerId string

const (
	SwarmPlotLayerGrid        SwarmPlotLayerId = "grid"
	SwarmPlotLayerAxes        SwarmPlotLayerId = "axes"
	SwarmPlotLayerCircles     SwarmPlotLayerId = "circles"
	SwarmPlotLayerAnnotations SwarmPlotLayerId = "annotations"
	SwarmPlotLayerMesh        SwarmPlotLayerId = "mesh"
)

// DefaultLayers mirrors @nivo/swarmplot defaultProps.layers.
var DefaultLayers = []SwarmPlotLayerId{
	SwarmPlotLayerGrid, SwarmPlotLayerAxes, SwarmPlotLayerCircles,
	SwarmPlotLayerAnnotations, SwarmPlotLayerMesh,
}

// SwarmPlotProps mirrors @nivo/swarmplot SwarmPlotSvgProps (the supported
// subset). Fields left zero fall back to Defaults via applyDefaults.
type SwarmPlotProps struct {
	Data []SwarmPlotDatum
	// Groups fixes the group order/set; when empty it is derived from Data in
	// first-appearance order.
	Groups []string

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	ValueScale  scales.ScaleSpec // default linear {min:0, max:auto}
	ValueFormat string           // d3-format spec; empty → %g

	Size    float64 // node diameter
	Spacing float64 // extra gap between nodes (added to the collide radius)
	Layout  string  // "vertical" (default) | "horizontal"
	Gap     float64 // extra gap between groups

	ForceStrength        float64
	SimulationIterations int

	Colors  colors.OrdinalColorScaleConfig
	ColorBy string // "group" (default) | "id"

	BorderWidth float64
	BorderColor string

	EnableGridX bool
	EnableGridY bool
	GridXValues []any
	GridYValues []any
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	// Interactive enables the client-side hover layer (charts/interact): each
	// node emits a data-tc-tooltip the script shows on hover.
	Interactive bool
	// UseMesh routes hover through an accurate voronoi mesh (charts/interact,
	// backed by internal/d3/delaunay). Active only when Interactive.
	UseMesh bool
	// DebugMesh draws the voronoi cells as a faint guide when UseMesh is on.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds mesh hit-testing to this pixel distance.
	DetectionRadius float64

	Theme  *theming.Theme
	Layers []SwarmPlotLayerId

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
	// so static output is unchanged.
	Animate       bool
	MotionStagger float64
}

// SwarmPlotResult is the computed model produced by UseSwarmPlot.
type SwarmPlotResult struct {
	Nodes  []ComputedNode
	XScale scales.Scale
	YScale scales.Scale
}
