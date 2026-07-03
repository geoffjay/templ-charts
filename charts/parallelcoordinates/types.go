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
// v3 scope: SVG only, static render with optional client-side hover tooltips
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
	Key   string
	Type  PCScaleType
	Label string // axis legend; empty → Key

	// Linear variable domain overrides. nil → auto (data extent).
	Min *float64
	Max *float64

	// Point variable categories (first-seen order used when empty).
	Values []string

	TickCount int // linear tick count hint; 0 → default
	Padding   float64
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
	Data      []PCDatum
	Variables []PCVariable

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Layout            PCLayout
	Curve             core.CurveFactoryId
	LineWidth         float64
	LineOpacity       float64
	AxesTicksPosition string // "before" | "after"

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

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []PCLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// PCResult is the computed model produced by UseParallelCoordinates.
type PCResult struct {
	Variables  []ComputedVariable
	Lines      []ComputedLine
	LegendData []legends.Datum
}
