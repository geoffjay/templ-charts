// Package sankey mirrors @nivo/sankey: a flow diagram whose nodes are sized by
// throughput and connected by variable-thickness ribbons. The layout is
// produced server-side by internal/d3/sankey (a faithful d3-sankey port) and
// the ribbons are drawn with internal/d3/shape's line generator + curveMonotoneX
// /Y over the ribbon outline, exactly as nivo does.
//
// It reuses internal/d3/sankey (layout), internal/d3/shape (line + monotone
// curves for ribbons), charts/colors (ordinal node colors + inherited border/
// label colors), charts/core (SvgWrapper + gradient defs), charts/legends,
// charts/theming and charts/interact (optional client-side hover tooltips).
//
// SVG only; the layout is deterministic, so goldens are byte-stable. Interactive
// enables per-node/link hover tooltips (charts/interact) plus a chord-style
// hover-highlight (scoped CSS :has() — dims others, re-lights the connected
// node/links). Enter animation is available via Animate.
package sankey

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// SankeyInputNode is one input node, identified by ID.
type SankeyInputNode struct {
	// ID uniquely identifies the node; link Source/Target reference it and it is
	// used as the node's label (see Label) and ordinal color key. Required.
	ID string
}

// SankeyInputLink is one directed, weighted link between two node ids.
type SankeyInputLink struct {
	// Source is the ID of the node the flow originates from.
	Source string
	// Target is the ID of the node the flow terminates at.
	Target string
	// Value is the flow magnitude; it sets the link ribbon thickness and
	// contributes to each connected node's total throughput (and thus size).
	Value float64
}

// SankeyLayout is the flow orientation. Mirrors @nivo/sankey layout.
type SankeyLayout string

const (
	SankeyLayoutHorizontal SankeyLayout = "horizontal"
	SankeyLayoutVertical   SankeyLayout = "vertical"
)

// SankeyAlign selects the node-alignment strategy. Mirrors @nivo/sankey align.
type SankeyAlign string

const (
	SankeyAlignCenter  SankeyAlign = "center"
	SankeyAlignJustify SankeyAlign = "justify"
	SankeyAlignStart   SankeyAlign = "start"
	SankeyAlignEnd     SankeyAlign = "end"
)

// SankeySort selects node/link ordering. Mirrors @nivo/sankey sort.
type SankeySort string

const (
	SankeySortAuto       SankeySort = "auto"
	SankeySortInput      SankeySort = "input"
	SankeySortAscending  SankeySort = "ascending"
	SankeySortDescending SankeySort = "descending"
)

// SankeyLabelPosition places labels inside (toward the node) or outside.
// Mirrors @nivo/sankey labelPosition.
type SankeyLabelPosition string

const (
	SankeyLabelInside  SankeyLabelPosition = "inside"
	SankeyLabelOutside SankeyLabelPosition = "outside"
)

// SankeyLabelOrientation is the label text orientation.
type SankeyLabelOrientation string

const (
	SankeyLabelHorizontal SankeyLabelOrientation = "horizontal"
	SankeyLabelVertical   SankeyLabelOrientation = "vertical"
)

// ComputedNode is one positioned, colored node ready to render. X/Y/Width/
// Height are the render rectangle (already accounting for layout + inner
// padding); X0/Y0/X1/Y1 are the edge coordinates the link ribbons anchor to.
type ComputedNode struct {
	ID             string
	Label          string
	FormattedValue string
	Value          float64
	X              float64
	Y              float64
	Width          float64
	Height         float64
	X0             float64
	Y0             float64
	X1             float64
	Y1             float64
	Color          string
}

// ComputedLink is one positioned ribbon. Path is the SVG path string; Pos0/Pos1
// are the ribbon center offsets at the source/target edges and Thickness its
// (contracted) width.
type ComputedLink struct {
	Index          int
	Source         string
	Target         string
	Value          float64
	FormattedValue string
	Path           string
	Color          string
	StartColor     string
	EndColor       string
	Pos0           float64
	Pos1           float64
	Thickness      float64
	// GradX0/GradY0 and GradX1/GradY1 are the absolute (userSpaceOnUse) anchor
	// points of the ribbon's flow used for the link gradient: the center of the
	// source edge and the center of the target edge respectively.
	GradX0 float64
	GradY0 float64
	GradX1 float64
	GradY1 float64
}

// SankeyLayerId enumerates the render layers. Mirrors @nivo/sankey LayerId.
type SankeyLayerId string

const (
	SankeyLayerLinks   SankeyLayerId = "links"
	SankeyLayerNodes   SankeyLayerId = "nodes"
	SankeyLayerLabels  SankeyLayerId = "labels"
	SankeyLayerLegends SankeyLayerId = "legends"
)

// DefaultLayers mirrors @nivo/sankey svgDefaultProps.layers.
var DefaultLayers = []SankeyLayerId{SankeyLayerLinks, SankeyLayerNodes, SankeyLayerLabels, SankeyLayerLegends}

// SankeyProps mirrors @nivo/sankey SankeySvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type SankeyProps struct {
	// Nodes is the set of graph nodes. Required.
	Nodes []SankeyInputNode
	// Links is the set of directed, weighted flows between nodes. Required.
	Links []SankeyInputLink

	// Width and Height are the inner drawing dimensions in pixels (the outer svg
	// adds Margin around them).
	Width  float64
	Height float64
	// Margin is the space reserved around the inner drawing area, e.g. for
	// outside labels and legends.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Layout is the flow orientation: horizontal (left→right) or vertical
	// (top→bottom). Empty → Defaults ("horizontal").
	Layout SankeyLayout
	// Align is the node-alignment strategy across columns: center, justify,
	// start, or end. Empty → Defaults ("center").
	Align SankeyAlign
	// Sort is the node/link ordering within each column: auto (d3 default),
	// input (preserve input order), ascending, or descending by node value.
	// Empty → Defaults ("auto").
	Sort SankeySort

	// Colors is the ordinal color scale used to color nodes by ID (links inherit
	// their source node color). Zero value → Defaults (scheme "nivo").
	Colors colors.OrdinalColorScaleConfig

	// NodeOpacity is the fill opacity of node rects, 0–1. Zero → Defaults (0.75).
	NodeOpacity float64
	// NodeThickness is the node width (cross-flow extent) in pixels. Zero →
	// Defaults (12).
	NodeThickness float64
	// NodeSpacing is the gap in pixels between nodes at the same level. Zero →
	// Defaults (12).
	NodeSpacing float64
	// NodeInnerPadding is the padding in pixels subtracted from each side of a
	// node (distance from the link), reducing the drawn node thickness. Zero →
	// Defaults (0).
	NodeInnerPadding float64
	// NodeBorderWidth is the node rect stroke width in pixels. Zero → Defaults
	// (1).
	NodeBorderWidth float64
	// NodeBorderColor is the node border color; may inherit from the node color.
	// Zero value → Defaults (node color darkened by 0.5).
	NodeBorderColor colors.InheritedColorConfig
	// NodeBorderRadius is the node rect corner radius in pixels. Zero → Defaults
	// (0, square corners).
	NodeBorderRadius float64

	// LinkOpacity is the fill opacity of link ribbons, 0–1. Zero → Defaults
	// (0.25).
	LinkOpacity float64
	// LinkContract shrinks each link ribbon by this many pixels on each side,
	// leaving a gap between adjacent ribbons. Zero → Defaults (0).
	LinkContract float64
	// LinkBlendMode is the CSS mix-blend-mode applied to link ribbons (e.g.
	// "multiply", "normal", "screen"). Empty → Defaults ("multiply").
	LinkBlendMode string
	// EnableLinkGradient draws each link ribbon with a per-link
	// <linearGradient> running from the source node color to the target node
	// color along the flow direction, instead of a flat source-color fill.
	// Mirrors @nivo/sankey enableLinkGradient (default false). nil → false.
	EnableLinkGradient *bool

	EnableLabels *bool  // nil → true
	Label        string // node field to use as label (only "id" supported)
	// LabelPosition places labels inside (toward the node) or outside it. Empty
	// → Defaults ("inside").
	LabelPosition SankeyLabelPosition
	// LabelPadding is the gap in pixels between a label and its node. Zero →
	// Defaults (9).
	LabelPadding float64
	// LabelOrientation is the label text orientation: horizontal or vertical.
	// Empty → Defaults ("horizontal").
	LabelOrientation SankeyLabelOrientation
	// LabelTextColor is the label color; may inherit from the node color. Zero
	// value → Defaults (node color darkened by 0.8).
	LabelTextColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	// Interactive enables per-node/link client-side hover tooltips
	// (charts/interact) AND the chord-style hover-highlight: hovering a node or
	// link dims the rest and re-lights the connected elements via a scoped CSS
	// :has() <style> block (no JS/server round-trip). The *HoverOpacity /
	// *HoverOthersOpacity fields control the highlighted / dimmed opacities;
	// zero values fall back to the Defaults.
	Interactive bool

	// Hover-highlight opacities (used when Interactive). NodeHoverOpacity /
	// LinkHoverOpacity are applied to the hovered element and its connected
	// elements; NodeHoverOthersOpacity / LinkHoverOthersOpacity dim everything
	// else. Mirror @nivo/sankey's nodeHover*/linkHover* props.
	NodeHoverOpacity       float64
	NodeHoverOthersOpacity float64
	LinkHoverOpacity       float64
	LinkHoverOthersOpacity float64

	// Legends are the optional legend blocks rendered from the node data.
	Legends []legends.LegendProps

	// Theme overrides the styling (fonts, colors, background). nil → the package
	// DefaultTheme.
	Theme *theming.Theme
	// Layers selects which render layers to draw and in what order (links, nodes,
	// labels, legends). Empty → Defaults (DefaultLayers, all four in order).
	Layers []SankeyLayerId

	// Role is the ARIA role of the root svg. Empty → Defaults ("img").
	Role string
	// AriaLabel sets the root svg aria-label attribute.
	AriaLabel string
	// AriaLabelledBy sets the root svg aria-labelledby attribute.
	AriaLabelledBy string
	// AriaDescribedBy sets the root svg aria-describedby attribute.
	AriaDescribedBy string
	// Title sets the svg <title> element (accessible name).
	Title string
	// Desc sets the svg <desc> element (accessible description).
	Desc string
	// IsFocusable makes the root svg keyboard-focusable (tabindex/focusable).
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// node rect and link path (600ms). MotionStagger delays successive items by
	// that many seconds (0 ⇒ all enter together). Defaults off, so static output
	// is unchanged. Mirrors @nivo/sankey's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// LabelsEnabled resolves EnableLabels (nil → true, matching nivo's default).
func (p SankeyProps) LabelsEnabled() bool { return p.EnableLabels == nil || *p.EnableLabels }

// LinkGradientEnabled resolves EnableLinkGradient (nil → false, nivo default).
func (p SankeyProps) LinkGradientEnabled() bool {
	return p.EnableLinkGradient != nil && *p.EnableLinkGradient
}

// SankeyResult is the computed model produced by UseSankey.
type SankeyResult struct {
	Nodes      []ComputedNode
	Links      []ComputedLink
	LegendData []legends.Datum
}
