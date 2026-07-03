// Package radialbar mirrors @nivo/radial-bar: stacked bars drawn as arcs in
// polar space, one radius band per serie (group) and one arc per category,
// stacked along an angle scale. It reuses charts/arcs (ArcGenerator + ArcsLayer
// + ArcLabelsLayer), charts/polar-axes (PolarGrid, CircularAxis, RadialAxis),
// charts/scales (the new band/linear range constructors), charts/colors
// (ordinal color per category), charts/legends, charts/theming and charts/core.
//
// Angle convention: the value scale maps [0, maxValue] onto [startAngle,
// endAngle] in degrees (0 = top, clockwise); arc angles are those degrees
// converted to radians and handed to the d3-shape arc generator (which places
// angle 0 at the top).
//
// v2 scope: SVG only, static render. The interactive hover tooltip arrives with
// the Phase 5 client-side layer.
package radialbar

import (
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// RadialBarDatum is one category's value within a serie. Mirrors
// @nivo/radial-bar RadialBarDatum ({x, y}).
type RadialBarDatum struct {
	X string
	Y float64
}

// RadialBarSerie is one group: an id plus its per-category data. Mirrors
// @nivo/radial-bar RadialBarSerie.
type RadialBarSerie struct {
	ID   string
	Data []RadialBarDatum
}

// ComputedBar is one positioned, colored arc. Mirrors @nivo/radial-bar
// ComputedBar.
type ComputedBar struct {
	ID             string
	Data           RadialBarDatum
	GroupID        string
	Category       string
	Value          float64
	FormattedValue string
	Color          string
	StackedValue   float64
	Arc            arcs.Arc
}

// RadialBarTrackDatum is one background track arc. Mirrors @nivo/radial-bar
// RadialBarTrackDatum.
type RadialBarTrackDatum struct {
	ID    string
	Color string
	Arc   arcs.Arc
}

// RadialBarLayerId enumerates the render layers. Mirrors @nivo/radial-bar
// RadialBarLayerId.
type RadialBarLayerId string

const (
	RadialBarLayerGrid    RadialBarLayerId = "grid"
	RadialBarLayerTracks  RadialBarLayerId = "tracks"
	RadialBarLayerBars    RadialBarLayerId = "bars"
	RadialBarLayerLabels  RadialBarLayerId = "labels"
	RadialBarLayerLegends RadialBarLayerId = "legends"
)

// DefaultLayers mirrors @nivo/radial-bar commonDefaultProps.layers.
var DefaultLayers = []RadialBarLayerId{
	RadialBarLayerGrid, RadialBarLayerTracks, RadialBarLayerBars, RadialBarLayerLabels, RadialBarLayerLegends,
}

// RadialBarProps mirrors @nivo/radial-bar RadialBarSvgProps (the supported
// subset). Fields left zero fall back to Defaults via applyDefaults.
type RadialBarProps struct {
	Data []RadialBarSerie

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// MaxValue is the value-scale domain max. nil → "auto" (max group total).
	MaxValue    *float64
	ValueFormat string // d3-format spec; empty → %g

	StartAngle   float64 // degrees
	EndAngle     float64 // degrees
	InnerRadius  float64 // ratio in [0,1] of the outer radius
	Padding      float64 // band padding between series rings
	PadAngle     float64 // degrees
	CornerRadius float64

	// EnableTracks gates the background track arcs. nil → true (nivo default).
	EnableTracks *bool
	TracksColor  string

	// Grid toggles. nil → nivo default (radial+circular grids on; radial-start
	// axis + circular-outer axis on; radial-end + circular-inner off).
	EnableRadialGrid      *bool
	EnableCircularGrid    *bool
	ShowRadialAxisStart   *bool
	ShowRadialAxisEnd     *bool
	ShowCircularAxisInner *bool
	ShowCircularAxisOuter *bool

	// Interactive enables the client-side hover layer (charts/interact): each
	// bar arc emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	// EnableLabels gates per-arc labels. nil → false (nivo default).
	EnableLabels       *bool
	Label              string // datum path; empty → "formattedValue"
	LabelsSkipAngle    float64
	LabelsRadiusOffset float64
	LabelsTextColor    colors.InheritedColorConfig

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []RadialBarLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// RadialBarResult is the computed model produced by UseRadialBar.
type RadialBarResult struct {
	Center       [2]float64
	InnerRadius  float64
	OuterRadius  float64
	Bars         []ComputedBar
	Tracks       []RadialBarTrackDatum
	ArcGenerator *arcs.ArcGenerator
	RadiusScale  scales.Scale
	ValueScale   scales.Scale
	StartAngle   float64 // degrees, clamped
	EndAngle     float64 // degrees, clamped
	LegendData   []legends.Datum
	Categories   []string
}

// TracksEnabled resolves EnableTracks (nil → true).
func (p RadialBarProps) TracksEnabled() bool { return p.EnableTracks == nil || *p.EnableTracks }

// LabelsEnabled resolves EnableLabels (nil → false).
func (p RadialBarProps) LabelsEnabled() bool { return p.EnableLabels != nil && *p.EnableLabels }

// BoolPtr returns a pointer to b — a helper for *bool props.
func BoolPtr(b bool) *bool { return &b }

// FloatPtr returns a pointer to f — a helper for *float64 props (MaxValue).
func FloatPtr(f float64) *float64 { return &f }
