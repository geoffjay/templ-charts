// Package funnel mirrors @nivo/funnel: ordered parts whose widths interpolate
// between successive values, drawn as smooth (curveBasis) or linear trapezoid
// bands with optional separators and centered labels. It reuses charts/scales
// (linear value scale), charts/colors (ordinal part color + inherited
// border/label colors), the internal/d3/shape Area + Line generators and
// charts/theming.
//
// v2 scope: SVG only, static render. The interactive hover tooltip arrives with
// the Phase 5 client layer; annotations are deferred.
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
	Data []FunnelDatum

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	Direction     FunnelDirection
	Interpolation FunnelInterpolation
	Spacing       float64
	ShapeBlending float64

	Colors      colors.OrdinalColorScaleConfig
	FillOpacity float64
	ValueFormat string

	BorderWidth   float64
	BorderColor   colors.InheritedColorConfig
	BorderOpacity float64

	// EnableLabel gates centered part labels. nil → true (nivo default).
	EnableLabel *bool
	LabelColor  colors.InheritedColorConfig

	// EnableBeforeSeparators / EnableAfterSeparators gate the separator lines.
	// nil → true (nivo default).
	EnableBeforeSeparators *bool
	BeforeSeparatorOffset  float64
	EnableAfterSeparators  *bool
	AfterSeparatorOffset   float64

	Theme  *theming.Theme
	Layers []FunnelLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	IsFocusable     bool

	core.MotionProps
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

// BoolPtr returns a pointer to b — a helper for *bool props.
func BoolPtr(b bool) *bool { return &b }
