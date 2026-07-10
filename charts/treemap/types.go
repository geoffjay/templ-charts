// Package treemap mirrors @nivo/treemap: a hierarchy laid out as nested
// rectangles, each node's area proportional to its value, tiled with the
// configured algorithm (squarify by default). It reuses internal/d3/hierarchy
// (Hierarchy + Treemap), charts/colors (ordinal color, grouped by the depth-1
// ancestor), charts/core (SvgWrapper) and charts/theming.
//
// SVG only, with optional client-side hover tooltips per node (charts/interact)
// via Interactive and opt-in click-to-zoom drill-down via EnableZooming +
// ChartID (through the charts/htmx server-round-trip layer).
package treemap

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// TreemapNode is one node of the input hierarchy. Leaves carry a Value; a
// parent's value is the sum of its descendants.
type TreemapNode struct {
	ID       string
	Value    float64
	Children []TreemapNode
}

// TileType selects the tiling algorithm. Mirrors @nivo/treemap tile.
type TileType string

const (
	TileSquarify  TileType = "squarify"
	TileBinary    TileType = "binary"
	TileDice      TileType = "dice"
	TileSlice     TileType = "slice"
	TileSliceDice TileType = "sliceDice"
)

// ComputedNode is one positioned, colored rectangle.
type ComputedNode struct {
	ID             string
	Path           string // dot-joined ancestor ids (root..node)
	Depth          int
	Value          float64
	FormattedValue string
	X, Y           float64
	Width, Height  float64
	Color          string
	IsParent       bool
}

// TreemapProps mirrors @nivo/treemap TreeMapSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type TreemapProps struct {
	// Data is the root of the input hierarchy; leaves carry a Value and each
	// node's area is proportional to the sum of its descendants' values.
	Data TreemapNode

	// Width is the outer chart width in pixels (including Margin).
	Width float64
	// Height is the outer chart height in pixels (including Margin).
	Height float64
	// Margin is the space reserved around the tiled area, in pixels.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Tile selects the tiling algorithm. Defaults to squarify.
	Tile TileType
	// InnerPadding is the gap in pixels between sibling rectangles. Defaults
	// to 0.
	InnerPadding float64
	// OuterPadding is the padding in pixels inside each parent before its
	// children are laid out. Defaults to 0.
	OuterPadding float64

	// Colors is the ordinal color scale config for node rects, grouped by the
	// depth-1 ancestor.
	Colors colors.OrdinalColorScaleConfig
	// NodeOpacity is the fill-opacity of each rect, in [0,1]. Defaults to 0.33.
	NodeOpacity float64
	// BorderWidth is the stroke width of rect borders in pixels; borders are
	// drawn only when > 0. Defaults to 1.
	BorderWidth float64
	// BorderColor is the inherited-color config for the border stroke
	// (resolved relative to each node's color).
	BorderColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	EnableLabel *bool // nil → true
	// LabelSkipSize hides a leaf label when the smaller side of its rect
	// (min(width,height)) is below this many pixels. Defaults to 0 (never
	// skipped).
	LabelSkipSize float64

	EnableParentLabel *bool // nil → true
	// ParentLabelSize is the height budget in pixels for a parent label; a
	// parent label is shown only when its rect height is at least this.
	// Defaults to 20.
	ParentLabelSize float64
	// ParentLabelPadding is the inset in pixels of a parent label from the
	// rect's leading edge. Defaults to 6.
	ParentLabelPadding float64

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool

	// EnableZooming turns each rect into a click-to-zoom target and renders a
	// breadcrumb when focused. It only takes effect when ChartID is also set
	// (htmx mode): clicking a node re-runs the treemap layout on that node's
	// subtree via GET /charts/{ChartID}/zoom?node=<id>. Default false → the
	// rendered SVG is byte-identical to the un-zoomable output. nil → false.
	EnableZooming *bool

	// ChartID is the htmx registry instance id. When set the chart emits hx-*
	// wiring scoped to this id, mirroring bar/line/pie. Empty for standalone
	// renders.
	ChartID string

	// FocusID is the id of the currently-focused node (the zoom target). Empty
	// (or the root id) shows the full chart. Set by the htmx zoom handler on a
	// props clone; consumers normally leave it zero.
	FocusID string

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// node rect. Default false: the rendered SVG is byte-identical to the
	// un-animated output. MotionStagger delays each successive node by that many
	// seconds (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

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
}

// Crumb is one segment of the zoom breadcrumb: the node id (used to build the
// zoom-back hx-get) and a display label.
type Crumb struct {
	ID    string
	Label string
}

// TreemapResult is the computed model produced by UseTreemap.
type TreemapResult struct {
	Nodes []ComputedNode
	// Breadcrumb is the root→focus ancestor path, non-nil only when focused.
	Breadcrumb []Crumb
}

// LabelEnabled resolves EnableLabel (nil → true).
func (p TreemapProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }

// ParentLabelEnabled resolves EnableParentLabel (nil → true).
func (p TreemapProps) ParentLabelEnabled() bool {
	return p.EnableParentLabel == nil || *p.EnableParentLabel
}

// ZoomingEnabled resolves EnableZooming (nil → false, nivo default).
func (p TreemapProps) ZoomingEnabled() bool { return p.EnableZooming != nil && *p.EnableZooming }
