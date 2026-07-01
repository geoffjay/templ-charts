// Package circlepacking mirrors @nivo/circle-packing: a hierarchy laid out as
// nested circles (Welzl enclosing), each leaf's area proportional to its value.
// It reuses internal/d3/hierarchy (Hierarchy + Pack — deterministic via the
// ported LCG), charts/colors (colored by depth), charts/core and charts/theming.
//
// v3 scope: SVG only, static render with optional client-side hover tooltips
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

// CirclePackingResult is the computed model produced by UseCirclePacking.
type CirclePackingResult struct {
	Circles []ComputedCircle
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// LabelsEnabled resolves EnableLabels (nil → false).
func (p CirclePackingProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }
