// Package bump mirrors @nivo/bump: a ranking-over-time chart. Each serie holds
// one {x, y} point per column; x is a category placed on a point scale and y is
// the serie's rank there, placed on a linear ranking scale (rank 1 at the top).
// The series are drawn as smooth "bump" lines (curveBumpX) connecting their
// ranks across columns, with a point per rank and optional start/end labels.
//
// It reuses charts/scales (point x scale + linear y scale via
// ComputeXYScalesForSeries), charts/axes (grid + axes), charts/colors (ordinal
// color per serie), charts/core (SvgWrapper, DotsItem), charts/legends,
// charts/theming, and internal/d3/shape (the line generator + curveBumpX/Y).
//
// SVG only. Direct point hover is available via Interactive (the
// charts/interact client layer), with an accurate voronoi-mesh layer backed by
// the internal/d3/delaunay port.
package bump

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// BumpDatum is one column's rank for a serie: x (the column, number | string |
// time) and y (the rank; nil marks a missing rank so the line breaks). Mirrors
// @nivo/bump's serie datum.
type BumpDatum struct {
	X any
	Y any
}

// BumpSerie is one ranked serie: an id plus its per-column ranks. Mirrors
// @nivo/bump BumpSerie.
type BumpSerie struct {
	ID   string
	Data []BumpDatum
}

// BumpPoint is one positioned rank dot. Mirrors @nivo/bump BumpPoint.
type BumpPoint struct {
	ID           string
	SerieID      string
	IndexInSerie int
	X            float64
	Y            float64
	XValue       any
	YValue       any
	FormattedX   string
	FormattedY   string
	Color        string
	Size         float64
	// Defined is false when the rank was nil (the point is not drawn but the
	// line still breaks around it).
	Defined bool
}

// ComputedSerie is one positioned, colored serie ready to render: its line path
// plus its drawn points. Mirrors @nivo/bump ComputedSerie.
type ComputedSerie struct {
	ID        string
	Color     string
	LineWidth float64
	Opacity   float64
	LinePath  string
	Points    []BumpPoint
}

// BumpLayerId enumerates the render layers. Mirrors @nivo/bump BumpLayerId.
type BumpLayerId string

const (
	BumpLayerGrid   BumpLayerId = "grid"
	BumpLayerAxes   BumpLayerId = "axes"
	BumpLayerLabels BumpLayerId = "labels"
	BumpLayerLines  BumpLayerId = "lines"
	BumpLayerPoints BumpLayerId = "points"
	BumpLayerMesh   BumpLayerId = "mesh"
)

// DefaultLayers mirrors @nivo/bump svgDefaultProps.layers.
var DefaultLayers = []BumpLayerId{
	BumpLayerGrid, BumpLayerAxes, BumpLayerLabels, BumpLayerLines, BumpLayerPoints, BumpLayerMesh,
}

// Interpolation selects the line interpolation: "smooth" (curveBumpX) or
// "linear". Mirrors @nivo/bump interpolation.
type Interpolation string

const (
	InterpolationSmooth Interpolation = "smooth"
	InterpolationLinear Interpolation = "linear"
)

// BumpProps mirrors @nivo/bump BumpSvgProps (the supported subset). Fields left
// zero fall back to Defaults via applyDefaults.
type BumpProps struct {
	// Data holds one BumpSerie per ranked serie, each with its per-column ranks.
	Data []BumpSerie

	// Width and Height are the total SVG dimensions in pixels; the plot area is
	// these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the plot area (for axes/labels/legends).
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// Interpolation selects the line shape: "smooth" (curveBumpX, default) or
	// "linear".
	Interpolation Interpolation
	// XPadding is the inner padding fraction (0..1) of the x point scale; default
	// 0.6. XOuterPadding and YOuterPadding are the outer padding fractions of the
	// x and y scales; both default 0.5.
	XPadding      float64
	XOuterPadding float64
	YOuterPadding float64

	// LineWidth is the serie line width in pixels; default 2. ActiveLineWidth
	// (default 4) and InactiveLineWidth (default 1) apply to hovered vs.
	// non-hovered series when Interactive.
	LineWidth         float64
	ActiveLineWidth   float64
	InactiveLineWidth float64
	// Opacity is the serie line/point opacity; default 1. ActiveOpacity (default
	// 1) and InactiveOpacity (default 0.3) apply to hovered vs. non-hovered
	// series when Interactive.
	Opacity         float64
	ActiveOpacity   float64
	InactiveOpacity float64

	// StartLabel / EndLabel toggle the serie-id labels at the first / last
	// column. nil → Defaults (start off, end on).
	StartLabel *bool
	EndLabel   *bool

	// PointSize is the rank dot diameter in pixels; default 6. ActivePointSize
	// (default 8) and InactivePointSize (default 4) apply to hovered vs.
	// non-hovered series when Interactive.
	PointSize         float64
	ActivePointSize   float64
	InactivePointSize float64

	// Colors is the ordinal color scale mapping each serie to a color; default
	// the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig

	XFormat string // d3-format spec; empty → %g
	YFormat string

	// Interactive enables the client-side hover layer (charts/interact): each
	// point emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool
	// UseMesh enables the voronoi-mesh hover layer for nearest-point detection.
	UseMesh bool
	// DebugMesh renders the voronoi mesh cell outlines for debugging.
	DebugMesh bool

	EnableGridX *bool // nil → true (nivo default)
	EnableGridY *bool // nil → true (nivo default)
	// AxisTop/Right/Bottom/Left configure each axis; nil hides that axis.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Legends configures zero or more legends; empty → no legend.
	Legends []legends.LegendProps

	// Theme overrides chart styling; nil → theming.DefaultTheme.
	Theme *theming.Theme
	// Layers sets the render order of chart layers; empty → DefaultLayers.
	Layers []BumpLayerId

	// Role is the root SVG ARIA role; default "img".
	Role string
	// AriaLabel, AriaLabelledBy, and AriaDescribedBy set the corresponding SVG
	// accessibility attributes.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title>/<desc> elements.
	Title string
	Desc  string
	// IsFocusable makes the SVG keyboard-focusable.
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// series line (600ms). MotionStagger delays successive series by that many
	// seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/bump's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// BumpResult is the computed model produced by UseBump.
type BumpResult struct {
	Series     []ComputedSerie
	XScale     scales.Scale
	YScale     scales.Scale
	LegendData []legends.Datum
}

// StartLabelEnabled resolves StartLabel (nil → false).
func (p BumpProps) StartLabelEnabled() bool { return p.StartLabel != nil && *p.StartLabel }

// EndLabelEnabled resolves EndLabel (nil → true).
func (p BumpProps) EndLabelEnabled() bool { return p.EndLabel == nil || *p.EndLabel }

// GridXEnabled resolves EnableGridX (nil → true, nivo default).
func (p BumpProps) GridXEnabled() bool { return p.EnableGridX == nil || *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p BumpProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }
