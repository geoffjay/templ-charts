// Package sunburst mirrors @nivo/sunburst: a hierarchy laid out with the
// partition layout mapped onto polar coordinates — each node is an arc whose
// angular span is proportional to its value and whose radius reflects its
// depth. It reuses internal/d3/hierarchy (Hierarchy + Partition), charts/arcs
// (ArcGenerator + ArcsLayer + ArcLabelsLayer), charts/colors, charts/core and
// charts/theming.
//
// The partition runs on [2π, r²]; each node's arc is
// {startAngle:x0, endAngle:x1, innerRadius:√y0, outerRadius:√y1}.
//
// SVG only, static render with optional client-side hover tooltips
// per arc (charts/interact) via Interactive.
package sunburst

import (
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// SunburstNode is one node of the input hierarchy.
type SunburstNode struct {
	ID       string
	Value    float64
	Children []SunburstNode
}

// ComputedArc is one positioned, colored arc.
type ComputedArc struct {
	ID             string
	Depth          int
	Value          float64
	FormattedValue string
	Color          string
	Arc            arcs.Arc
}

// SunburstProps mirrors @nivo/sunburst SunburstSvgProps (the supported subset).
type SunburstProps struct {
	Data SunburstNode

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	CornerRadius float64

	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
	BorderColor string

	ValueFormat string // d3-format spec; empty → %g

	EnableArcLabels       *bool // nil → false
	ArcLabelsRadiusOffset float64
	ArcLabelsSkipAngle    float64

	// Interactive enables per-arc client-side hover tooltips (charts/interact).
	Interactive bool

	// EnableZooming turns each arc into a click-to-zoom target and renders a
	// breadcrumb when focused. It only takes effect when ChartID is also set
	// (htmx mode): clicking an arc re-renders the chart focused on that node's
	// subtree via GET /charts/{ChartID}/zoom?node=<id>. Default false → the
	// rendered SVG is byte-identical to the un-zoomable output.
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

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// arc (600ms). MotionStagger delays successive arcs by that many seconds
	// (0 ⇒ all arcs enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/sunburst's enter transition.
	Animate       bool
	MotionStagger float64
}

// Crumb is one segment of the zoom breadcrumb: the node id (used to build the
// zoom-back hx-get) and a display label.
type Crumb struct {
	ID    string
	Label string
}

// SunburstResult is the computed model produced by UseSunburst.
type SunburstResult struct {
	Center [2]float64
	Arcs   []ComputedArc
	// Breadcrumb is the root→focus ancestor path, non-nil only when focused.
	Breadcrumb []Crumb
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// ArcLabelsEnabled resolves EnableArcLabels (nil → false).
func (p SunburstProps) ArcLabelsEnabled() bool { return p.EnableArcLabels != nil && *p.EnableArcLabels }
