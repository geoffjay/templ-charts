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
// v3 scope: SVG only, static render. The layout is deterministic, so goldens
// are byte-stable. HTMX hover-others dimming and animated transitions are
// deferred; per-node/link hover tooltips are available via Interactive.
package sankey

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// SankeyInputNode is one input node, identified by ID.
type SankeyInputNode struct {
	ID string
}

// SankeyInputLink is one directed, weighted link between two node ids.
type SankeyInputLink struct {
	Source string
	Target string
	Value  float64
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
	Nodes []SankeyInputNode
	Links []SankeyInputLink

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Layout SankeyLayout
	Align  SankeyAlign
	Sort   SankeySort

	Colors colors.OrdinalColorScaleConfig

	NodeOpacity      float64
	NodeThickness    float64
	NodeSpacing      float64
	NodeInnerPadding float64
	NodeBorderWidth  float64
	NodeBorderColor  colors.InheritedColorConfig
	NodeBorderRadius float64

	LinkOpacity   float64
	LinkContract  float64
	LinkBlendMode string

	EnableLabels     *bool  // nil → true
	Label            string // node field to use as label (only "id" supported)
	LabelPosition    SankeyLabelPosition
	LabelPadding     float64
	LabelOrientation SankeyLabelOrientation
	LabelTextColor   colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	// Interactive enables per-node/link client-side hover tooltips (charts/interact).
	Interactive bool

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []SankeyLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// node rect and link path (600ms). MotionStagger delays successive items by
	// that many seconds (0 ⇒ all enter together). Defaults off, so static output
	// is unchanged. Mirrors @nivo/sankey's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// LabelsEnabled resolves EnableLabels (nil → true, matching nivo's default).
func (p SankeyProps) LabelsEnabled() bool { return p.EnableLabels == nil || *p.EnableLabels }

// SankeyResult is the computed model produced by UseSankey.
type SankeyResult struct {
	Nodes      []ComputedNode
	Links      []ComputedLink
	LegendData []legends.Datum
}
