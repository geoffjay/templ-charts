// Package parallelcoordinates mirrors @nivo/parallel-coordinates: a set of
// variables, each drawn as its own axis, with every datum rendered as a
// polyline connecting its scaled value on each variable's axis.
//
// Variables are laid out along the layout axis (horizontal → vertical axes
// spread across the width; vertical → horizontal axes spread down the height).
// Each variable owns its own scale (linear or point). Lines are drawn with the
// shared d3-shape line generator + the configured curve.
//
// It reuses charts/scales (the linear/point range constructors), charts/axes
// (one Axis per variable), charts/colors (ordinal color per datum),
// charts/core (SvgWrapper), charts/legends and charts/theming.
//
// SVG only, static render with optional client-side hover tooltips
// per line (charts/interact) via Interactive.
package parallelcoordinates

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// PCScaleType selects a variable's scale kind.
type PCScaleType string

const (
	PCScaleLinear PCScaleType = "linear"
	PCScalePoint  PCScaleType = "point"
)

// PCLayout is the layout orientation. Mirrors @nivo/parallel-coordinates layout.
type PCLayout string

const (
	PCLayoutHorizontal PCLayout = "horizontal"
	PCLayoutVertical   PCLayout = "vertical"
)

// PCVariable describes one variable (axis). Mirrors @nivo/parallel-coordinates
// variable spec.
type PCVariable struct {
	// Key is the datum value key this variable's axis reads.
	Key string
	// Type is the variable's scale kind: "linear" or "point".
	Type  PCScaleType
	Label string // axis legend; empty → Key

	// Linear variable domain overrides. nil → auto (data extent).
	Min *float64
	Max *float64

	// Point variable categories (first-seen order used when empty).
	Values []string

	TickCount int // linear tick count hint; 0 → default
	// Padding is the point scale padding (fraction of step, 0..1) for point
	// variables; ignored for linear variables.
	Padding float64
}

// PCDatum is one record: an id (for color/tooltip) plus a value per variable
// key. Linear variables read float64/int values; point variables read strings.
type PCDatum struct {
	ID     string
	Values map[string]any
}

// PCLayerId enumerates the render layers. Mirrors @nivo/parallel-coordinates.
type PCLayerId string

const (
	PCLayerAxes    PCLayerId = "axes"
	PCLayerLines   PCLayerId = "lines"
	PCLayerMesh    PCLayerId = "mesh"
	PCLayerLegends PCLayerId = "legends"
)

// DefaultLayers mirrors @nivo/parallel-coordinates defaultProps.layers.
var DefaultLayers = []PCLayerId{PCLayerLines, PCLayerAxes, PCLayerMesh, PCLayerLegends}

// ComputedVariable is a positioned variable with its resolved scale.
type ComputedVariable struct {
	Key   string
	Label string
	Type  PCScaleType
	// Pos is the variable's position along the layout axis (x for horizontal,
	// y for vertical), in pixels.
	Pos   float64
	Scale scales.Scale
}

// ComputedLine is one datum's polyline.
type ComputedLine struct {
	ID    string
	Color string
	Path  string
	// Points are the datum's pixel vertices where the line crosses each
	// variable axis (skipping undefined values). These are the mesh hover
	// targets when UseMesh is on.
	Points [][2]float64
	// Tooltip HTML (emitted only when Interactive).
	Tooltip string
}

// PCProps mirrors @nivo/parallel-coordinates ParallelCoordinatesSvgProps (the
// supported subset). Fields left zero fall back to Defaults via applyDefaults.
type PCProps struct {
	// Data is the set of records, each drawn as one polyline.
	Data []PCDatum
	// Variables define the axes; each owns its scale and layout position.
	Variables []PCVariable

	// Width and Height are the outer SVG dimensions in pixels; the inner plot
	// area is these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the inner plot area (for axes and
	// legends), in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Layout is the axis spread orientation, "horizontal" (default, vertical
	// axes across the width) or "vertical".
	Layout PCLayout
	// Curve is the d3-shape curve factory used to connect a datum's points.
	// Default linear.
	Curve core.CurveFactoryId
	// LineWidth is the polyline stroke width in pixels. Default 2.
	LineWidth float64
	// LineOpacity is the polyline stroke opacity (0..1). Default 0.5.
	LineOpacity       float64
	AxesTicksPosition string // "before" | "after"

	// Colors is the ordinal color scale used to color each datum's line.
	// Default is the "category10" scheme.
	Colors colors.OrdinalColorScaleConfig

	// Interactive enables per-line client-side hover tooltips (charts/interact).
	Interactive bool

	// UseMesh routes hover through an accurate voronoi mesh (charts/interact,
	// backed by internal/d3/delaunay) built over the per-axis line vertices
	// rather than per-line tooltips: hovering anywhere resolves to the nearest
	// datum vertex. Active only when Interactive.
	UseMesh bool
	// DebugMesh draws the voronoi cells as a faint guide when UseMesh is on.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds mesh hit-testing to this pixel distance.
	DetectionRadius float64

	// Legends configures zero or more legends describing the data.
	Legends []legends.LegendProps

	// Theme overrides the styling theme; nil uses the default theme.
	Theme *theming.Theme
	// Layers is the ordered list of render layers; default draws lines, axes,
	// mesh then legends.
	Layers []PCLayerId

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
	// datum polyline (600ms). MotionStagger delays successive lines by that many
	// seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/parallel-coordinates' enter transition.
	Animate       bool
	MotionStagger float64
}

// PCResult is the computed model produced by UseParallelCoordinates.
type PCResult struct {
	Variables  []ComputedVariable
	Lines      []ComputedLine
	LegendData []legends.Datum
}
