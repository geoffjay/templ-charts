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
	Data CirclePackingNode

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Padding float64

	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	EnableLabels     *bool // nil → false
	LabelsSkipRadius float64

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool

	// EnableZooming turns each circle into a click-to-zoom target and renders a
	// breadcrumb when focused. It only takes effect when ChartID is also set
	// (htmx mode): clicking a circle applies the d3 zoomable-pack transform so
	// that node fills the viewport, via GET /charts/{ChartID}/zoom?node=<id>.
	// Default false → the rendered SVG is byte-identical to the un-zoomable
	// output.
	EnableZooming bool

	// ChartID is the htmx registry instance id. When set the chart emits hx-*
	// wiring scoped to this id, mirroring bar/line/pie. Empty for standalone
	// renders.
	ChartID string

	// FocusID is the id of the currently-focused node (the zoom target). Empty
	// (or the root id) shows the full chart. Set by the htmx zoom handler on a
	// props clone; consumers normally leave it zero.
	FocusID string

	Theme *theming.Theme

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

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

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// LabelsEnabled resolves EnableLabels (nil → false).
func (p CirclePackingProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }
