// Package tree mirrors @nivo/tree: a node-link diagram laid out with the tidy
// tree (mode 'tree') or dendrogram (mode 'dendogram') algorithm, links drawn as
// smooth bump curves (curveBumpX/Y from internal/d3/shape). It reuses
// internal/d3/hierarchy (Hierarchy + Tree/Cluster), charts/colors, charts/core
// and charts/theming.
//
// SVG only, static render with optional client-side hover tooltips per node
// (charts/interact) via Interactive, plus an accurate voronoi-mesh hover layer
// backed by the internal/d3/delaunay port.
package tree

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// TreeNode is one node of the input hierarchy (structure only; tree layouts do
// not use values).
type TreeNode struct {
	ID       string
	Children []TreeNode
}

// Mode selects the layout algorithm. Mirrors @nivo/tree mode.
type Mode string

const (
	ModeDendogram Mode = "dendogram"
	ModeTree      Mode = "tree"
)

// LayoutDir is the growth direction. Mirrors @nivo/tree layout.
type LayoutDir string

const (
	LayoutTopToBottom LayoutDir = "top-to-bottom"
	LayoutBottomToTop LayoutDir = "bottom-to-top"
	LayoutLeftToRight LayoutDir = "left-to-right"
	LayoutRightToLeft LayoutDir = "right-to-left"
)

// ComputedNode is one positioned, colored node.
type ComputedNode struct {
	ID    string
	Depth int
	X, Y  float64
	Color string
}

// ComputedLink is one parent→child edge with its bump path.
type ComputedLink struct {
	SourceID string
	TargetID string
	Path     string
	Color    string
}

// TreeProps mirrors @nivo/tree TreeSvgProps (the supported subset).
type TreeProps struct {
	// Data is the root of the input hierarchy; only its structure (ids and
	// nesting) is used, values are ignored.
	Data TreeNode

	// Width is the outer chart width in pixels (including Margin).
	Width float64
	// Height is the outer chart height in pixels (including Margin).
	Height float64
	// Margin is the space reserved around the plot area (for labels), in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Mode selects the layout algorithm: dendogram (the default) or tree.
	Mode Mode
	// Layout is the growth direction of the diagram. Defaults to
	// top-to-bottom.
	Layout LayoutDir

	// NodeSize is the diameter of each node circle in pixels (the drawn radius
	// is NodeSize/2). Defaults to 12.
	NodeSize float64
	// Colors is the ordinal color scale config used to color nodes (keyed by
	// node id). Matches the Colors field on other chart families.
	Colors colors.OrdinalColorScaleConfig
	// LinkThickness is the stroke width of the parent→child links in pixels.
	// Defaults to 1.
	LinkThickness float64
	// LinkOpacity is the stroke opacity of the links, in [0,1]. Defaults to
	// 0.4.
	LinkOpacity float64

	// EnableLabel toggles the per-node id labels. nil → true (labels shown).
	EnableLabel *bool
	// LabelOffset is the gap in pixels between a node and its label. Defaults
	// to 6.
	LabelOffset float64

	// Interactive enables client-side hover (charts/interact). With UseMesh
	// (the default), the whole area resolves to the nearest node via an accurate
	// voronoi mesh (internal/d3/delaunay); with UseMesh off, each node carries
	// its own hover tooltip.
	Interactive bool
	// UseMesh selects the voronoi-mesh hover path over per-node tooltips.
	UseMesh bool
	// DebugMesh draws the voronoi cells as a faint guide when UseMesh is on.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds mesh hit-testing to this pixel distance.
	DetectionRadius float64

	// Theme overrides the styling theme (colors, fonts, label styles). Nil
	// uses the package default theme.
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

	// Animate emits a SMIL enter animation (600ms) — links fade in and nodes
	// scale their radius from 0 — staggered by MotionStagger seconds per item.
	// Defaults off, so static output is unchanged.
	Animate       bool
	MotionStagger float64
}

// TreeResult is the computed model produced by UseTree.
type TreeResult struct {
	Nodes []ComputedNode
	Links []ComputedLink
}

// LabelEnabled resolves EnableLabel (nil → true).
func (p TreeProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }
