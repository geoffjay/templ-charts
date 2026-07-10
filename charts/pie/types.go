// Package pie mirrors @nivo/pie: DefaultRawDatum, ComputedDatum, PieArc,
// PieProps/PieSvgProps, the NormalizeData/PieArcs/PieFromBox compute
// functions, and the Pie/Arcs/ArcLinkLabels/ArcLabels/PieLegends/PieTooltip
// templ components.
//
// Renders server-side SVG only (no Canvas). Interactivity (hover/click
// active-arc highlight, series toggle) is routed through charts/htmx.
//
// Angle convention: user-facing startAngle/endAngle/padAngle are in degrees
// (0 = top, clockwise). Converted to radians before d3-shape. nivo compensates
// with -90/-π2 offsets in ComputeArcBoundingBox and ComputeArcCenter.
package pie

import (
	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// DatumId is the union of string | number for datum ids. Modeled as string.
type DatumId = string

// DefaultRawDatum is the default datum shape: id + value. Mirrors @nivo/pie
// DefaultRawDatum.
type DefaultRawDatum struct {
	ID    DatumId
	Value float64
	Label string // optional; falls back to ID
}

// MayHaveLabel is the constraint on raw datum: any struct with an optional
// `label` field. We model this as `any` for Go ergonomics; the label accessor
// resolves it via PropertyAccessor.
type MayHaveLabel = any

// PieArc extends arcs.Arc with index, mid-angle, thickness, and padAngle.
// Mirrors @nivo/pie PieArc.
type PieArc struct {
	arcs.Arc
	Index     int
	Angle     float64 // |endAngle - startAngle| in radians
	AngleDeg  float64 // |end - start| in degrees
	Thickness float64 // outerRadius - innerRadius
	PadAngle  float64 // radians
}

// ComputedDatum is one normalized pie datum with its arc. Mirrors @nivo/pie
// ComputedDatum.
type ComputedDatum struct {
	ID             DatumId
	Label          DatumId
	Value          float64
	FormattedValue string
	Color          string
	Fill           string // only when gradient/pattern defs match
	Data           any
	Arc            PieArc
	Hidden         bool
}

// LegendDatum is one legend entry. Mirrors @nivo/pie LegendDatum.
type LegendDatum struct {
	ID     DatumId
	Label  DatumId
	Color  string
	Hidden bool
	Data   ComputedDatumNoArc
}

// ComputedDatumNoArc is ComputedDatum without the Arc/Fill fields (used in
// LegendDatum.Data).
type ComputedDatumNoArc struct {
	ID             DatumId
	Label          DatumId
	Value          float64
	FormattedValue string
	Color          string
	Data           any
	Hidden         bool
}

// PieLayerId enumerates the built-in pie layers.
type PieLayerId string

const (
	PieLayerArcs          PieLayerId = "arcs"
	PieLayerArcLinkLabels PieLayerId = "arcLinkLabels"
	PieLayerArcLabels     PieLayerId = "arcLabels"
	PieLayerLegends       PieLayerId = "legends"
)

// DefaultLayers is the default pie layer order. Mirrors defaultProps.layers.
var DefaultLayers = []PieLayerId{
	PieLayerArcs, PieLayerArcLinkLabels, PieLayerArcLabels, PieLayerLegends,
}

// TransitionMode mirrors @nivo/arcs ArcTransitionMode. templ-charts supports
// "innerRadius" (grow outer radius from innerRadius).
type TransitionMode string

const (
	TransitionModeInnerRadius TransitionMode = "innerRadius"
)

// PieProps mirrors @nivo/pie CommonPieProps + PieSvgProps (the SVG subset).
// Fields with zero values fall back to Defaults via applyDefaults.
type PieProps struct {
	// Data.

	// Data holds the raw datums, one arc per entry. It is intentionally
	// []any: datums are read through reflection-based accessors (ID, Value,
	// and the label accessors) rather than a fixed struct, so entries may be
	// maps or arbitrary custom structs, mirroring nivo's generic RawDatum.
	// Required; there is no default.
	Data []any

	// Mapping.

	// ID is the accessor for each datum's unique identifier (used for keys,
	// colors, and legends). A string selects that property on the datum; a
	// function computes it. Defaults to "id".
	ID core.PropertyAccessor[any, DatumId]
	// Value is the accessor for each datum's numeric value, which determines
	// its arc's angular size. A string selects that property; a function
	// computes it. Defaults to "value".
	Value core.PropertyAccessor[any, float64]

	// Value formatting.

	// ValueFormat optionally formats raw values for display in arc labels
	// and tooltips (a d3-format specifier string or a function). When unset,
	// the raw value is used as-is.
	ValueFormat core.ValueFormat[float64]

	// Geometry.

	// Width is the total chart width in pixels (including margins). Required.
	Width float64
	// Height is the total chart height in pixels (including margins). Required.
	Height float64
	// Margin reserves space around the plotted pie for labels and legends.
	// Defaults to zero on all sides.
	Margin core.Margin

	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool

	// SortByValue orders arcs by their value (largest first) instead of by
	// data order. Defaults to false.
	SortByValue bool
	// InnerRadius turns the pie into a donut. Expressed as a ratio of the
	// outer radius in the range 0..1 (0 = full pie). Defaults to 0.
	InnerRadius float64
	// PadAngle is the gap between adjacent arcs, in degrees. Defaults to 0.
	PadAngle float64
	// CornerRadius rounds the corners of each arc, in pixels. Defaults to 0.
	CornerRadius float64
	// StartAngle is the angle of the first arc, in degrees (0 = top,
	// clockwise). Useful for partial pies and gauges. Defaults to 0.
	StartAngle float64
	// EndAngle is the angle at which the last arc ends, in degrees. Useful
	// for partial pies and gauges. Defaults to 360.
	EndAngle float64
	// Fit, when true, scales a partial pie (start/end angle not spanning a
	// full circle) to occupy the available space. Defaults to true; set to
	// false to keep the pie at its natural size.
	Fit bool
	// ActiveInnerRadiusOffset extends the inner radius of the active
	// (hovered/selected) arc by this many pixels. Defaults to 0.
	ActiveInnerRadiusOffset float64
	// ActiveOuterRadiusOffset extends the outer radius of the active
	// (hovered/selected) arc by this many pixels. Defaults to 0.
	ActiveOuterRadiusOffset float64
	// ActiveID programmatically marks a single arc as active so the active
	// radius offsets are applied to it. Empty (the default) means no arc is
	// active.
	ActiveID DatumId

	// Colors.

	// Colors configures the ordinal color scale used to fill arcs. Defaults
	// to the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// BorderWidth is the stroke width of each arc's border, in pixels.
	// Defaults to 0 (no border).
	BorderWidth float64
	// BorderColor computes each arc's border color; it may be a fixed color,
	// a theme reference, or a value derived from the arc's fill color.
	// Defaults to the arc color darkened by 1.
	BorderColor colors.InheritedColorConfig

	// Arc labels.

	// EnableArcLabels toggles the value labels drawn inside each arc. Nil
	// (the default) enables them, matching the nivo default.
	EnableArcLabels *bool
	// ArcLabel is the accessor for the text drawn inside each arc; a string
	// selects a ComputedDatum property, a function computes it. Defaults to
	// "formattedValue".
	ArcLabel core.PropertyAccessor[ComputedDatum, string]
	// ArcLabelsSkipAngle hides the arc label when the arc's angle is smaller
	// than this value, in degrees. Defaults to 0 (never skip).
	ArcLabelsSkipAngle float64
	// ArcLabelsSkipRadius hides the arc label when the arc's radius is
	// smaller than this value, in pixels. Defaults to 0 (never skip).
	ArcLabelsSkipRadius float64
	// ArcLabelsRadiusOffset positions the arc label along the radius,
	// expressed as a ratio from the inner radius outward. Defaults to 0.5.
	ArcLabelsRadiusOffset float64
	// ArcLabelsTextColor computes the arc label text color. Defaults to the
	// theme's labels.text.fill.
	ArcLabelsTextColor colors.InheritedColorConfig

	// Arc link labels.

	// EnableArcLinkLabels toggles the labels connected to each arc by a link
	// line outside the pie. Nil (the default) enables them, matching the
	// nivo default.
	EnableArcLinkLabels *bool
	// ArcLinkLabel is the accessor for the arc link label text; a string
	// selects a ComputedDatum property, a function computes it. Defaults to
	// "id".
	ArcLinkLabel core.PropertyAccessor[ComputedDatum, string]
	// ArcLinkLabelsSkipAngle hides the link label when the arc's angle is
	// smaller than this value, in degrees. Defaults to 0 (never skip).
	ArcLinkLabelsSkipAngle float64
	// ArcLinkLabelsOffset shifts the link's start from the pie's outer
	// radius, in pixels (negative values overlap the slices). Defaults to 0.
	ArcLinkLabelsOffset float64
	// ArcLinkLabelsDiagonalLength is the length of the link's diagonal
	// segment, in pixels. Defaults to 16.
	ArcLinkLabelsDiagonalLength float64
	// ArcLinkLabelsStraightLength is the length of the link's horizontal
	// segment, in pixels. Defaults to 24.
	ArcLinkLabelsStraightLength float64
	// ArcLinkLabelsThickness is the stroke width of the link line, in pixels.
	// Defaults to 1.
	ArcLinkLabelsThickness float64
	// ArcLinkLabelsTextOffset is the horizontal gap between the end of the
	// link and its text, in pixels. Defaults to 6.
	ArcLinkLabelsTextOffset float64
	// ArcLinkLabelsTextColor computes the arc link label text color.
	// Defaults to the theme's labels.text.fill.
	ArcLinkLabelsTextColor colors.InheritedColorConfig
	// ArcLinkLabelsColor computes the color of the link line. Defaults to the
	// theme's axis.ticks.line.stroke.
	ArcLinkLabelsColor colors.InheritedColorConfig

	// Legends.

	// Legends configures zero or more legends rendered alongside the pie.
	// Defaults to none.
	Legends []legends.LegendProps

	// Defs / fill.

	// Defs declares reusable SVG gradient/pattern definitions that Fill rules
	// can reference. Defaults to none.
	Defs []core.Def
	// Fill maps datums to definitions declared in Defs, letting arcs be
	// filled with a gradient or pattern instead of a solid color. Defaults to
	// none.
	Fill []core.DefRule

	// Interactivity.

	// Interactive enables hover/click behavior on arcs (active-arc highlight
	// and, via ChartID, HTMX events). Defaults to true.
	Interactive bool
	// InitialHiddenIDs lists datum IDs that start hidden (excluded from the
	// pie and toggleable through the legend). Defaults to none.
	InitialHiddenIDs []DatumId

	// Layers.

	// Layers defines the render order of the built-in layers. Defaults to
	// DefaultLayers (arcs, arcLinkLabels, arcLabels, legends).
	Layers []PieLayerId

	// Motion.
	core.MotionProps
	// TransitionMode controls how arcs animate on enter/update. Defaults to
	// TransitionModeInnerRadius.
	TransitionMode TransitionMode

	// A11y.

	// Role is the ARIA role of the root svg element. Defaults to "img".
	Role string
	// AriaLabel sets the aria-label on the root svg element. Empty by default.
	AriaLabel string
	// AriaLabelledBy sets the aria-labelledby on the root svg element. Empty
	// by default.
	AriaLabelledBy string
	// AriaDescribedBy sets the aria-describedby on the root svg element.
	// Empty by default.
	AriaDescribedBy string
	// Title sets the svg <title> element. Empty by default.
	Title string
	// Desc sets the svg <desc> element. Empty by default.
	Desc string
	// IsFocusable makes the svg and its arcs keyboard-focusable. Defaults to
	// false.
	IsFocusable bool

	// Theme.

	// Theme overrides the styling theme (colors, fonts, label styles). Nil
	// (the default) uses the package default theme.
	Theme *theming.Theme

	// HTMX chart instance ID; when non-empty, arcs emit hx-* attrs.
	ChartID string
}

// PieSvgProps is an alias of PieProps (nivo splits common/svg; this unifies them).
type PieSvgProps = PieProps

// PieResult is the output of the compute pipeline: normalized data, arcs,
// legend data, the arc generator, and the computed center/radii.
type PieResult struct {
	DataWithArc  []ComputedDatum
	LegendData   []LegendDatum
	ArcGenerator *arcs.ArcGenerator
	CenterX      float64
	CenterY      float64
	Radius       float64
	InnerRadius  float64
	HiddenIDs    []DatumId
}

// ArcLabelsEnabled resolves EnableArcLabels (nil → true, nivo default).
func (p PieProps) ArcLabelsEnabled() bool { return p.EnableArcLabels == nil || *p.EnableArcLabels }

// ArcLinkLabelsEnabled resolves EnableArcLinkLabels (nil → true, nivo default).
func (p PieProps) ArcLinkLabelsEnabled() bool {
	return p.EnableArcLinkLabels == nil || *p.EnableArcLinkLabels
}
