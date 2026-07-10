// Package voronoi mirrors @nivo/voronoi: it maps a set of 2-D data points from
// user domains into the chart area, builds their Delaunay triangulation and
// Voronoi diagram (internal/d3/delaunay), and renders the links (Delaunay
// edges), cells (Voronoi polygons), points, and bounding rectangle. It reuses
// charts/scales (linear x/y), charts/core (SvgWrapper) and charts/theming.
//
// SVG only, static render with optional per-cell client-side hover tooltips
// (charts/interact) via Interactive. The same delaunay port powers the accurate
// voronoi-mesh hover retrofitted into line/scatterplot/bump/tree.
package voronoi

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/d3/delaunay"
)

// VoronoiDatum is one input point. ID is optional (used for tooltips/keys).
type VoronoiDatum struct {
	ID string
	X  float64
	Y  float64
}

// LayerId names a render layer. Mirrors @nivo/voronoi layers.
type LayerId string

const (
	LayerLinks  LayerId = "links"
	LayerCells  LayerId = "cells"
	LayerPoints LayerId = "points"
	LayerBounds LayerId = "bounds"
)

// DefaultLayers is nivo's default layer order.
var DefaultLayers = []LayerId{LayerLinks, LayerCells, LayerPoints, LayerBounds}

// ComputedPoint is one input datum positioned in pixel space.
type ComputedPoint struct {
	ID    string
	X, Y  float64 // pixel coordinates within the inner chart area
	Color string  // resolved cell/point color (from Colors, keyed by ID)
	Data  VoronoiDatum
}

// VoronoiProps mirrors @nivo/voronoi VoronoiProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type VoronoiProps struct {
	// Data is the set of input points to triangulate. Each datum's X/Y are in
	// the user domains (see XDomain/YDomain).
	Data []VoronoiDatum

	// Width is the outer chart width in pixels (including Margin).
	Width float64
	// Height is the outer chart height in pixels (including Margin).
	Height float64
	// Margin is the space reserved around the plot area, in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// XDomain/YDomain map data to the chart area. Zero value → [0,1].
	XDomain [2]float64
	YDomain [2]float64

	// Layers is the ordered list of render layers to draw. Empty falls back to
	// DefaultLayers (links, cells, points, bounds).
	Layers []LayerId

	EnableLinks *bool // nil → false
	// LinkLineWidth is the stroke width in pixels of the Delaunay link edges.
	// Defaults to 1.
	LinkLineWidth float64
	// LinkLineColor is the stroke color of the link edges. Defaults to
	// "#bbbbbb".
	LinkLineColor string

	EnableCells *bool // nil → true
	// CellLineWidth is the stroke width in pixels of the Voronoi cell borders.
	// Defaults to 2.
	CellLineWidth float64
	// CellLineColor is the stroke color of the cell borders. Defaults to
	// "#000000".
	CellLineColor string

	// EnableCellFill fills each Voronoi cell with its site's color (from Colors)
	// rather than leaving it hollow. nil → false. When on, cells are rendered as
	// one path per cell (needed for per-cell fills).
	EnableCellFill  *bool
	CellFillOpacity float64 // 0 → default (Defaults.CellFillOpacity)

	// Colors assigns a color per site, keyed by datum ID (index when ID is
	// empty). Used for cell fills and available on ComputedPoint.Color.
	Colors colors.OrdinalColorScaleConfig

	EnablePoints *bool // nil → true
	// PointSize is the diameter in pixels of each site marker (the drawn
	// radius is PointSize/2). Defaults to 4.
	PointSize float64
	// PointColor is the fill color of the site markers. Defaults to "#666666".
	PointColor string

	// Interactive enables per-cell client-side hover tooltips (charts/interact).
	Interactive bool

	// Theme overrides the styling theme (colors, fonts). Nil uses the package
	// default theme.
	Theme *theming.Theme

	// Role is the ARIA role of the root svg element. Defaults to "img".
	Role string
	// AriaLabel sets the aria-label on the root svg element. Empty by default.
	AriaLabel string
	// AriaLabelledBy sets the aria-labelledby on the root svg element. Empty
	// by default.
	AriaLabelledBy string
	// AriaDescribedBy sets the aria-describedby on the root svg element.
	// Empty by default.
	AriaDescribedBy string
	// Title sets the svg <title> element. Empty by default.
	Title string
	// Desc sets the svg <desc> element. Empty by default.
	Desc string
	// IsFocusable makes the svg keyboard-focusable. Defaults to false.
	IsFocusable bool

	// Animate emits a SMIL enter animation (600ms) — cells fade in and points
	// scale their radius from 0 — staggered by MotionStagger seconds per item.
	// Defaults off, so static output is unchanged.
	Animate       bool
	MotionStagger float64
}

// VoronoiResult is the computed model produced by UseVoronoi.
type VoronoiResult struct {
	Points   []ComputedPoint
	Delaunay *delaunay.Delaunay
	Voronoi  *delaunay.Voronoi
}

// LinksEnabled resolves EnableLinks (nil → false).
func (p VoronoiProps) LinksEnabled() bool { return p.EnableLinks != nil && *p.EnableLinks }

// CellsEnabled resolves EnableCells (nil → true).
func (p VoronoiProps) CellsEnabled() bool { return p.EnableCells == nil || *p.EnableCells }

// PointsEnabled resolves EnablePoints (nil → true).
func (p VoronoiProps) PointsEnabled() bool { return p.EnablePoints == nil || *p.EnablePoints }

// CellFillEnabled resolves EnableCellFill (nil → false).
func (p VoronoiProps) CellFillEnabled() bool { return p.EnableCellFill != nil && *p.EnableCellFill }
