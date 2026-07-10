// Package circlepacking mirrors @nivo/circle-packing: a hierarchy laid out as
// nested circles (Welzl enclosing), each leaf's area proportional to its value.
// It reuses internal/d3/hierarchy (Hierarchy + Pack — deterministic via the
// ported LCG), charts/colors (colored by depth), charts/core and charts/theming.
//
// SVG only, static render with optional client-side hover tooltips
// per node (charts/interact) via Interactive. Interactive zoom is deferred.
package circlepacking

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// CirclePackingNode is one node of the input hierarchy.
type CirclePackingNode struct {
	ID       string
	Value    float64
	Children []CirclePackingNode
}

// ComputedCircle is one positioned, colored circle.
type ComputedCircle struct {
	ID             string
	Depth          int
	Value          float64
	FormattedValue string
	X, Y, R        float64
	Color          string
	IsLeaf         bool
}

// CirclePackingProps mirrors @nivo/circle-packing CirclePackingSvgProps.
type CirclePackingProps struct {
	// Data is the root of the input hierarchy laid out as nested circles; each
	// leaf's area is proportional to its Value.
	Data CirclePackingNode

	// Width is the total chart width in pixels (including Margin).
	Width float64
	// Height is the total chart height in pixels (including Margin).
	Height float64
	// Margin reserves space around the packed circles.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Padding is the space inserted between adjacent circles by the pack
	// layout, in pixels. Default 0.
	Padding float64

	// Colors is the ordinal color scale; each node is colored by its depth.
	// Default is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// BorderWidth is the stroke width of each circle's border, in pixels.
	// Default 0 (no border drawn).
	BorderWidth float64
	// BorderColor resolves each circle's border color, typically inherited from
	// the circle's fill. Only applied when BorderWidth > 0.
	BorderColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	EnableLabels *bool // nil → false
	// LabelsSkipRadius hides a leaf's label when the circle's radius is smaller
	// than this value, in pixels. Default 8.
	LabelsSkipRadius float64

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool

	// EnableZooming turns each circle into a click-to-zoom target and renders a
	// breadcrumb when focused. It only takes effect when ChartID is also set
	// (htmx mode): clicking a circle applies the d3 zoomable-pack transform so
	// that node fills the viewport, via GET /charts/{ChartID}/zoom?node=<id>.
	// Default false (nil → false) → the rendered SVG is byte-identical to the
	// un-zoomable output.
	EnableZooming *bool

	// ChartID is the htmx registry instance id. When set the chart emits hx-*
	// wiring scoped to this id, mirroring bar/line/pie. Empty for standalone
	// renders.
	ChartID string

	// FocusID is the id of the currently-focused node (the zoom target). Empty
	// (or the root id) shows the full chart. Set by the htmx zoom handler on a
	// props clone; consumers normally leave it zero.
	FocusID string

	// Theme overrides the styling theme (colors, fonts, label styles). Nil uses
	// the default theme.
	Theme *theming.Theme

	// Role is the ARIA role of the root svg element. Default "img".
	Role string
	// AriaLabel sets aria-label on the root svg element. Empty by default.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the root svg element. Empty by default.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the root svg element. Empty by default.
	AriaDescribedBy string
	// Title sets the svg <title> element. Empty by default.
	Title string
	// Desc sets the svg <desc> element. Empty by default.
	Desc string
	// IsFocusable makes the root svg keyboard-focusable. Defaults to false.
	IsFocusable bool

	// Animate, when true, emits a SMIL enter animation on each circle scaling
	// its radius from 0 to its final value (600ms). MotionStagger delays
	// successive circles by that many seconds (0 ⇒ all enter together).
	// Defaults off, so static output is unchanged.
	Animate       bool
	MotionStagger float64
}

// Crumb is one segment of the zoom breadcrumb: the node id (used to build the
// zoom-back hx-get) and a display label.
type Crumb struct {
	ID    string
	Label string
}

// CirclePackingResult is the computed model produced by UseCirclePacking.
type CirclePackingResult struct {
	Circles []ComputedCircle
	// Breadcrumb is the root→focus ancestor path, non-nil only when focused.
	Breadcrumb []Crumb
}

// LabelsEnabled resolves EnableLabels (nil → false).
func (p CirclePackingProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }

// ZoomingEnabled resolves EnableZooming (nil → false, nivo default).
func (p CirclePackingProps) ZoomingEnabled() bool { return p.EnableZooming != nil && *p.EnableZooming }
