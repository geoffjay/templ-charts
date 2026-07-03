// Package treemap mirrors @nivo/treemap: a hierarchy laid out as nested
// rectangles, each node's area proportional to its value, tiled with the
// configured algorithm (squarify by default). It reuses internal/d3/hierarchy
// (Hierarchy + Treemap), charts/colors (ordinal color, grouped by the depth-1
// ancestor), charts/core (SvgWrapper) and charts/theming.
//
// v3 scope: SVG only, static render with optional client-side hover tooltips
// per node (charts/interact) via Interactive. Interactive zoom is deferred
// (see docs/PLAN-v3.md §9).
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
	Data TreemapNode

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	Tile         TileType
	InnerPadding float64
	OuterPadding float64

	Colors      colors.OrdinalColorScaleConfig
	NodeOpacity float64
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	EnableLabel   *bool // nil → true
	LabelSkipSize float64

	EnableParentLabel  *bool // nil → true
	ParentLabelSize    float64
	ParentLabelPadding float64

	// Interactive enables per-node client-side hover tooltips (charts/interact).
	Interactive bool

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// node rect. Default false: the rendered SVG is byte-identical to the
	// un-animated output. MotionStagger delays each successive node by that many
	// seconds (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

	Theme *theming.Theme

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// TreemapResult is the computed model produced by UseTreemap.
type TreemapResult struct {
	Nodes []ComputedNode
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// LabelEnabled resolves EnableLabel (nil → true).
func (p TreemapProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }

// ParentLabelEnabled resolves EnableParentLabel (nil → true).
func (p TreemapProps) ParentLabelEnabled() bool {
	return p.EnableParentLabel == nil || *p.EnableParentLabel
}
