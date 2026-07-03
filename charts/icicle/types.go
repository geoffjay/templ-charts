// Package icicle mirrors @nivo/icicle: a hierarchy laid out with the partition
// layout as nested rectangles banded by depth, oriented top/bottom/left/right.
// It reuses internal/d3/hierarchy (Hierarchy + Partition), charts/colors,
// charts/core and charts/theming.
//
// v3 scope: SVG only, static render with optional client-side hover tooltips
// per node (charts/interact) via Interactive. Interactive zoom is deferred
// (see docs/PLAN-v3.md §9).
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
	Data IcicleNode

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Orientation Orientation
	GapX        float64
	GapY        float64

	Colors       colors.OrdinalColorScaleConfig
	BorderWidth  float64
	BorderColor  colors.InheritedColorConfig
	BorderRadius float64

	ValueFormat string // d3-format spec; empty → %g

	EnableLabels *bool  // nil → false
	Label        string // "id" | "value" | "formattedValue"

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
}

// IcicleResult is the computed model produced by UseIcicle.
type IcicleResult struct {
	Rects []ComputedRect
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// LabelsEnabled resolves EnableLabels (nil → false).
func (p IcicleProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }
