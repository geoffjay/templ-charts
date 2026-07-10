// Package funnel mirrors @nivo/funnel: ordered parts whose widths interpolate
// between successive values, drawn as smooth (curveBasis) or linear trapezoid
// bands with optional separators and centered labels. It reuses charts/scales
// (linear value scale), charts/colors (ordinal part color + inherited
// border/label colors), the internal/d3/shape Area + Line generators and
// charts/theming.
//
// SVG only, static render. The interactive hover tooltip arrives with
// the charts/interact client layer; annotations are deferred.
package funnel

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// FunnelDirection is the funnel flow direction.
type FunnelDirection string

const (
	FunnelDirectionVertical   FunnelDirection = "vertical"
	FunnelDirectionHorizontal FunnelDirection = "horizontal"
)

// FunnelInterpolation is the part-edge interpolation.
type FunnelInterpolation string

const (
	FunnelInterpolationSmooth FunnelInterpolation = "smooth"
	FunnelInterpolationLinear FunnelInterpolation = "linear"
)

// FunnelLayerId enumerates the render layers. Mirrors @nivo/funnel FunnelLayerId.
type FunnelLayerId string

const (
	FunnelLayerSeparators  FunnelLayerId = "separators"
	FunnelLayerParts       FunnelLayerId = "parts"
	FunnelLayerLabels      FunnelLayerId = "labels"
	FunnelLayerAnnotations FunnelLayerId = "annotations"
)

// DefaultLayers mirrors @nivo/funnel svgDefaultProps.layers.
var DefaultLayers = []FunnelLayerId{
	FunnelLayerSeparators, FunnelLayerParts, FunnelLayerLabels, FunnelLayerAnnotations,
}

// FunnelDatum is one part: an id, a value, and an optional label. Mirrors
// @nivo/funnel FunnelDatum.
type FunnelDatum struct {
	ID    string
	Value float64
	Label string // optional; falls back to ID
}

// ComputedPart is one laid-out funnel part ready to render.
type ComputedPart struct {
	ID             string
	FormattedValue string
	Label          string
	Color          string
	FillOpacity    float64
	BorderColor    string
	BorderWidth    float64
	BorderOpacity  float64
	LabelColor     string
	X              float64 // center
	Y              float64
	AreaPath       string
	BorderPathA    string // one side
	BorderPathB    string // other side
}

// Separator is one before/after separator line.
type Separator struct {
	X0, Y0, X1, Y1 float64
}

// FunnelProps mirrors @nivo/funnel FunnelSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type FunnelProps struct {
	// Data is the ordered list of parts; each part's width is derived from its
	// Value relative to its neighbors.
	Data []FunnelDatum

	// Width is the total chart width in pixels (including Margin).
	Width float64
	// Height is the total chart height in pixels (including Margin).
	Height float64
	// Margin reserves space around the funnel for labels and separators.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// Interactive enables the client-side hover layer (charts/interact): each
	// part emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Direction is the funnel flow direction (vertical or horizontal). Default
	// vertical.
	Direction FunnelDirection
	// Interpolation is the part-edge interpolation: smooth (curveBasis) or
	// linear (straight trapezoids). Default smooth.
	Interpolation FunnelInterpolation
	// Spacing is the gap between successive parts, in pixels. Default 0.
	Spacing float64
	// ShapeBlending controls how strongly each part's edges bend toward its
	// neighbors' widths, from 0 (angular steps) to 1 (fully blended). Default
	// 0.66.
	ShapeBlending float64

	// Colors is the ordinal color scale assigning a fill to each part. Default
	// is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// FillOpacity is the fill opacity of each part, from 0 to 1. Default 1.
	FillOpacity float64
	// ValueFormat is a d3-format spec for part values; empty → %g.
	ValueFormat string

	// BorderWidth is the stroke width of each part's side borders, in pixels.
	// Default 6.
	BorderWidth float64
	// BorderColor resolves each part's border color, by default inherited from
	// the part's fill.
	BorderColor colors.InheritedColorConfig
	// BorderOpacity is the opacity of the part borders, from 0 to 1. Default
	// 0.66.
	BorderOpacity float64

	// EnableLabel gates centered part labels. nil → true (nivo default).
	EnableLabel *bool
	// LabelColor resolves the centered label color, by default the theme
	// background color.
	LabelColor colors.InheritedColorConfig

	// EnableBeforeSeparators / EnableAfterSeparators gate the separator lines.
	// nil → true (nivo default).
	EnableBeforeSeparators *bool
	// BeforeSeparatorOffset is the gap, in pixels, between a part and its
	// leading (before) separator line. Default 0.
	BeforeSeparatorOffset float64
	EnableAfterSeparators *bool
	// AfterSeparatorOffset is the gap, in pixels, between a part and its
	// trailing (after) separator line. Default 0.
	AfterSeparatorOffset float64

	// Theme overrides the styling theme (colors, fonts, label styles). Nil uses
	// the default theme.
	Theme *theming.Theme
	// Layers selects which render layers to draw and their order. Default
	// DefaultLayers (separators, parts, labels, annotations).
	Layers []FunnelLayerId

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

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// funnel part (600ms). MotionStagger delays successive parts by that many
	// seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/funnel's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// FunnelResult is the computed model produced by UseFunnel.
type FunnelResult struct {
	Parts            []ComputedPart
	BeforeSeparators []Separator
	AfterSeparators  []Separator
}

// LabelEnabled resolves EnableLabel (nil → true).
func (p FunnelProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }

// BeforeSeparatorsEnabled resolves EnableBeforeSeparators (nil → true).
func (p FunnelProps) BeforeSeparatorsEnabled() bool {
	return p.EnableBeforeSeparators == nil || *p.EnableBeforeSeparators
}

// AfterSeparatorsEnabled resolves EnableAfterSeparators (nil → true).
func (p FunnelProps) AfterSeparatorsEnabled() bool {
	return p.EnableAfterSeparators == nil || *p.EnableAfterSeparators
}
