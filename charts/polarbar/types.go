// Package polarbar mirrors @nivo/polar-bar: a stacked bar chart wrapped into a
// full circle. Each index occupies an angular band (a band scale over the
// indices maps to [startAngle, endAngle]); within a band the keys are stacked
// radially from the inner radius outward, mapped through a linear radius scale
// over the stacked totals. Distinct from radial-bar (which stacks along the
// angle over a 0–270° sweep with a per-serie radius ring); polar-bar stacks
// along the radius over a full 0–360° circle with a per-index angle band.
//
// It reuses charts/arcs (ArcGenerator + ArcsLayer + ArcLabelsLayer),
// charts/polar-axes (PolarGrid, CircularAxis, RadialAxis), charts/scales (the
// band/linear range constructors), charts/colors (ordinal color per key),
// charts/legends, charts/theming and charts/core.
//
// SVG only, static render with optional client-side hover tooltips
// per arc (charts/interact) via Interactive.
package polarbar

import (
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// PolarBarDatum is one index row: an index id plus a value per key. Mirrors a
// @nivo/polar-bar data object (indexBy → Index; keys → Values entries).
type PolarBarDatum struct {
	Index  string
	Values map[string]float64
}

// ComputedArc is one positioned, colored stacked arc. Mirrors @nivo/polar-bar
// ComputedBar.
type ComputedArc struct {
	ID             string
	Index          string
	Key            string
	Value          float64
	FormattedValue string
	StackLo        float64
	StackHi        float64
	Color          string
	Arc            arcs.Arc
}

// PolarBarLayerId enumerates the render layers. Mirrors @nivo/polar-bar.
type PolarBarLayerId string

const (
	PolarBarLayerGrid    PolarBarLayerId = "grid"
	PolarBarLayerArcs    PolarBarLayerId = "arcs"
	PolarBarLayerAxes    PolarBarLayerId = "axes"
	PolarBarLayerLabels  PolarBarLayerId = "labels"
	PolarBarLayerLegends PolarBarLayerId = "legends"
)

// DefaultLayers mirrors @nivo/polar-bar defaultProps.layers.
var DefaultLayers = []PolarBarLayerId{
	PolarBarLayerGrid, PolarBarLayerArcs, PolarBarLayerAxes, PolarBarLayerLabels, PolarBarLayerLegends,
}

// PolarBarProps mirrors @nivo/polar-bar PolarBarSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type PolarBarProps struct {
	Data    []PolarBarDatum
	Keys    []string
	IndexBy string

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	StartAngle   float64 // degrees
	EndAngle     float64 // degrees
	InnerRadius  float64 // ratio in [0,1] of the outer radius
	CornerRadius float64
	PadAngle     float64 // degrees
	Padding      float64 // band padding between index bands

	// MaxValue overrides the radius-scale domain max. nil → max stacked total.
	MaxValue    *float64
	ValueFormat string // d3-format spec; empty → %g

	// Grid toggles. nil → nivo default (both on).
	EnableRadialGrid   *bool
	EnableCircularGrid *bool

	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	// EnableArcLabels gates per-arc labels. nil → false (nivo default).
	EnableArcLabels       *bool
	ArcLabel              string // datum path; empty → "formattedValue"
	ArcLabelsRadiusOffset float64
	ArcLabelsSkipAngle    float64

	// Interactive enables per-arc client-side hover tooltips (charts/interact).
	Interactive bool

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []PolarBarLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// bar arc (600ms). MotionStagger delays successive arcs by that many
	// seconds (0 ⇒ all arcs enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/polar-bar's enter transition.
	Animate       bool
	MotionStagger float64
}

// PolarBarResult is the computed model produced by UsePolarBar.
type PolarBarResult struct {
	Center       [2]float64
	InnerRadius  float64
	OuterRadius  float64
	Arcs         []ComputedArc
	ArcGenerator *arcs.ArcGenerator
	AngleScale   scales.Scale
	RadiusScale  scales.Scale
	StartAngle   float64 // degrees, clamped
	EndAngle     float64 // degrees, clamped
	LegendData   []legends.Datum
	Indices      []string
	MaxValue     float64
}

// BoolPtr returns a pointer to b — a helper for the *bool props.
func BoolPtr(b bool) *bool { return &b }

// FloatPtr returns a pointer to f — a helper for the *float64 props.
func FloatPtr(f float64) *float64 { return &f }

// RadialGridEnabled resolves EnableRadialGrid (nil → true).
func (p PolarBarProps) RadialGridEnabled() bool {
	return p.EnableRadialGrid == nil || *p.EnableRadialGrid
}

// CircularGridEnabled resolves EnableCircularGrid (nil → true).
func (p PolarBarProps) CircularGridEnabled() bool {
	return p.EnableCircularGrid == nil || *p.EnableCircularGrid
}

// ArcLabelsEnabled resolves EnableArcLabels (nil → false).
func (p PolarBarProps) ArcLabelsEnabled() bool { return p.EnableArcLabels != nil && *p.EnableArcLabels }
