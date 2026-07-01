// Package tree mirrors @nivo/tree: a node-link diagram laid out with the tidy
// tree (mode 'tree') or dendrogram (mode 'dendogram') algorithm, links drawn as
// smooth bump curves (curveBumpX/Y from internal/d3/shape). It reuses
// internal/d3/hierarchy (Hierarchy + Tree/Cluster), charts/colors, charts/core
// and charts/theming.
//
// v3 scope: SVG only, static render with optional client-side hover tooltips
// per node (charts/interact) via Interactive. The accurate voronoi-mesh hover
// layer arrives with the Phase 3 d3-delaunay port (docs/PLAN-v3.md §4.3).
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
	Data TreeNode

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Mode   Mode
	Layout LayoutDir

	NodeSize      float64
	NodeColor     colors.OrdinalColorScaleConfig
	LinkThickness float64
	LinkOpacity   float64

	EnableLabel *bool // nil → true
	LabelOffset float64

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool
	// UseMesh is accepted for API parity; the accurate voronoi-mesh layer is
	// deferred to the Phase 3 d3-delaunay port.
	UseMesh bool

	Theme *theming.Theme

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

	core.MotionProps
}

// TreeResult is the computed model produced by UseTree.
type TreeResult struct {
	Nodes []ComputedNode
	Links []ComputedLink
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// LabelEnabled resolves EnableLabel (nil → true).
func (p TreeProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }
