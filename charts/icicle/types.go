// Package icicle mirrors @nivo/icicle: a hierarchy laid out with the partition
// layout as nested rectangles banded by depth, oriented top/bottom/left/right.
// It reuses internal/d3/hierarchy (Hierarchy + Partition), charts/colors,
// charts/core and charts/theming.
//
// SVG only, with optional client-side hover tooltips per node (charts/interact)
// via Interactive and opt-in click-to-zoom drill-down via EnableZooming +
// ChartID (through the charts/htmx server-round-trip layer).
package icicle

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// IcicleNode is one node of the input hierarchy.
type IcicleNode struct {
	ID       string
	Value    float64
	Children []IcicleNode
}

// Orientation is the growth direction of the depth axis. Mirrors @nivo/icicle
// orientation.
type Orientation string

const (
	OrientationBottom Orientation = "bottom"
	OrientationTop    Orientation = "top"
	OrientationLeft   Orientation = "left"
	OrientationRight  Orientation = "right"
)

// ComputedRect is one positioned, colored rectangle.
type ComputedRect struct {
	ID             string
	Depth          int
	Value          float64
	FormattedValue string
	X, Y           float64
	Width, Height  float64
	Color          string
}

// IcicleProps mirrors @nivo/icicle IcicleSvgProps (the supported subset).
type IcicleProps struct {
	// Data is the root of the input hierarchy laid out as nested rectangles;
	// each leaf's area is proportional to its Value.
	Data IcicleNode

	// Width is the total chart width in pixels (including Margin).
	Width float64
	// Height is the total chart height in pixels (including Margin).
	Height float64
	// Margin reserves space around the partitioned rectangles.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Orientation is the growth direction of the depth axis (bottom, top, left,
	// or right). Default bottom.
	Orientation Orientation
	// GapX is the horizontal gap between adjacent rectangles, in pixels.
	// Default 1.
	GapX float64
	// GapY is the vertical gap between adjacent rectangles, in pixels. Default 1.
	GapY float64

	// Colors is the ordinal color scale; rectangles are banded by depth.
	// Default is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// BorderWidth is the stroke width of each rectangle's border, in pixels.
	// Default 0 (no border drawn).
	BorderWidth float64
	// BorderColor resolves each rectangle's border color, typically inherited
	// from the fill. Only applied when BorderWidth > 0.
	BorderColor colors.InheritedColorConfig
	// BorderRadius is the corner radius of the rectangles, in pixels. Default 0.
	BorderRadius float64

	ValueFormat string // d3-format spec; empty → %g

	EnableLabels *bool  // nil → false
	Label        string // "id" | "value" | "formattedValue"

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool

	// EnableZooming turns each rect into a click-to-zoom target and renders a
	// breadcrumb when focused. It only takes effect when ChartID is also set
	// (htmx mode): clicking a node re-renders the chart focused on that node's
	// subtree via GET /charts/{ChartID}/zoom?node=<id>. nil → true (nivo
	// default). Zoom wiring is still only emitted when ChartID is set, so
	// standalone/static renders remain byte-identical regardless.
	EnableZooming *bool

	// ChartID is the htmx registry instance id. When set (with Interactive /
	// EnableZooming) the chart emits hx-* wiring scoped to this id, mirroring
	// bar/line/pie. Empty for standalone/static renders.
	ChartID string

	// FocusID is the id of the currently-focused node (the zoom target). Empty
	// (or the root id) shows the full chart. It is set by the htmx zoom handler
	// on a clone of the props; consumers normally leave it zero.
	FocusID string

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// rect. Default false: the rendered SVG is byte-identical to the un-animated
	// output. MotionStagger delays each successive rect by that many seconds
	// (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

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
}

// Crumb is one segment of the zoom breadcrumb: the node id (used to build the
// zoom-back hx-get) and a display label.
type Crumb struct {
	ID    string
	Label string
}

// IcicleResult is the computed model produced by UseIcicle.
type IcicleResult struct {
	Rects []ComputedRect
	// Breadcrumb is the root→focus ancestor path, non-nil only when focused.
	Breadcrumb []Crumb
}

// LabelsEnabled resolves EnableLabels (nil → false).
func (p IcicleProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }

// ZoomingEnabled resolves EnableZooming (nil → true, nivo default).
func (p IcicleProps) ZoomingEnabled() bool { return p.EnableZooming == nil || *p.EnableZooming }
