// Package pie mirrors @nivo/pie: DefaultRawDatum, ComputedDatum, PieArc,
// PieProps/PieSvgProps, the NormalizeData/PieArcs/PieFromBox compute
// functions, and the Pie/Arcs/ArcLinkLabels/ArcLabels/PieLegends/PieTooltip
// templ components.
//
// v1 renders server-side SVG only (no Canvas). Interactivity (hover/click
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

// TransitionMode mirrors @nivo/arcs ArcTransitionMode. v1 supports
// "innerRadius" (grow outer radius from innerRadius).
type TransitionMode string

const (
	TransitionModeInnerRadius TransitionMode = "innerRadius"
)

// PieProps mirrors @nivo/pie CommonPieProps + PieSvgProps (the SVG subset).
// Fields with zero values fall back to Defaults via applyDefaults.
type PieProps struct {
	// Data.
	Data []any

	// Mapping.
	ID    core.PropertyAccessor[any, DatumId]
	Value core.PropertyAccessor[any, float64]

	// Value formatting.
	ValueFormat core.ValueFormat[float64]

	// Geometry.
	Width  float64
	Height float64
	Margin core.Margin

	SortByValue             bool
	InnerRadius             float64
	PadAngle                float64 // degrees
	CornerRadius            float64
	StartAngle              float64 // degrees
	EndAngle                float64 // degrees
	Fit                     bool
	ActiveInnerRadiusOffset float64
	ActiveOuterRadiusOffset float64
	ActiveID                DatumId

	// Colors.
	Colors      colors.OrdinalColorScaleConfig
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	// Arc labels.
	EnableArcLabels       bool
	ArcLabel              core.PropertyAccessor[ComputedDatum, string]
	ArcLabelsSkipAngle    float64
	ArcLabelsSkipRadius   float64
	ArcLabelsRadiusOffset float64
	ArcLabelsTextColor    colors.InheritedColorConfig

	// Arc link labels.
	EnableArcLinkLabels         bool
	ArcLinkLabel                core.PropertyAccessor[ComputedDatum, string]
	ArcLinkLabelsSkipAngle      float64
	ArcLinkLabelsOffset         float64
	ArcLinkLabelsDiagonalLength float64
	ArcLinkLabelsStraightLength float64
	ArcLinkLabelsThickness      float64
	ArcLinkLabelsTextOffset     float64
	ArcLinkLabelsTextColor      colors.InheritedColorConfig
	ArcLinkLabelsColor          colors.InheritedColorConfig

	// Legends.
	Legends []legends.LegendProps

	// Defs / fill.
	Defs []core.Def
	Fill []core.DefRule

	// Interactivity.
	IsInteractive    bool
	InitialHiddenIDs []DatumId

	// Layers.
	Layers []PieLayerId

	// Motion.
	core.MotionProps
	TransitionMode TransitionMode

	// A11y.
	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string

	// Theme.
	Theme *theming.Theme

	// HTMX chart instance ID; when non-empty, arcs emit hx-* attrs.
	ChartID string
}

// PieSvgProps is an alias of PieProps (nivo splits common/svg; v1 unifies).
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
