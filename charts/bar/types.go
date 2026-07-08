// Package bar mirrors @nivo/bar: BarDatum, ComputedDatum, ComputedBarDatum,
// BarLegendProps, BarLayerId, BarProps/BarSvgProps, the compute pipeline
// (grouped/stacked/legends/totals), the UseBar orchestrator, and the
// Bar/BarItem/BarTotals/BarAnnotations/BarLegends templ components.
//
// Renders server-side SVG only (no Canvas). Interactivity is routed
// through charts/htmx (HTMX attrs emitted on bars/legend items).
package bar

import (
	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/bar/compute"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// BarDatum is a single row of bar-chart data: a map from key to value
// (string | number). Mirrors @nivo/bar BarDatum = Record<string, string|number>.
type BarDatum = map[string]any

// ComputedDatum is one key's value within an index, after normalization.
// Mirrors @nivo/bar ComputedDatum<D>. Value is nil when the raw value was
// null/missing (nivo: `value: rawValue === null ? rawValue : value`).
//
// Defined in compute (to avoid an import cycle) and re-exported here.
type ComputedDatum = compute.ComputedDatum

// ComputedBarDatum is a ComputedDatum resolved to a bar rectangle (x/y/w/h)
// plus color and label. Mirrors @nivo/bar ComputedBarDatum<D>.
type ComputedBarDatum = compute.ComputedBarDatum

// ComputedBarDatumWithValue is a ComputedBarDatum whose Data.Value is non-nil.
// Used by the render path (nivo's barsWithValue filter). Modelled as a
// pointer alias of ComputedBarDatum for simplicity — the renderer checks
// Data.Value != nil.
type ComputedBarDatumWithValue = ComputedBarDatum

// LegendData is one legend entry. Mirrors @nivo/bar LegendData.
type LegendData = compute.LegendData

// BarTotalData is one totals label. Mirrors @nivo/bar BarTotalsData.
type BarTotalData = compute.BarTotalData

// BarLegendProps extends LegendProps with `dataFrom` ("indexes"|"keys").
// Mirrors @nivo/bar BarLegendProps.
type BarLegendProps struct {
	legends.LegendProps
	DataFrom string // "indexes" | "keys" (default "keys")
}

// BarLayerId enumerates the built-in bar layers.
type BarLayerId string

const (
	BarLayerGrid        BarLayerId = "grid"
	BarLayerAxes        BarLayerId = "axes"
	BarLayerBars        BarLayerId = "bars"
	BarLayerMarkers     BarLayerId = "markers"
	BarLayerLegends     BarLayerId = "legends"
	BarLayerAnnotations BarLayerId = "annotations"
	BarLayerTotals      BarLayerId = "totals"
)

// DefaultLayers is the default bar layer order. Mirrors svgDefaultProps.layers.
var DefaultLayers = []BarLayerId{
	BarLayerGrid, BarLayerAxes, BarLayerBars, BarLayerTotals,
	BarLayerMarkers, BarLayerLegends, BarLayerAnnotations,
}

// GroupMode is "stacked" or "grouped".
type GroupMode string

const (
	GroupModeStacked GroupMode = "stacked"
	GroupModeGrouped GroupMode = "grouped"
)

// Layout is "vertical" or "horizontal".
type Layout string

const (
	LayoutVertical   Layout = "vertical"
	LayoutHorizontal Layout = "horizontal"
)

// ColorBy is "id" or "indexValue".
type ColorBy string

const (
	ColorByID         ColorBy = "id"
	ColorByIndexValue ColorBy = "indexValue"
)

// LabelPosition is "start", "middle", or "end".
type LabelPosition string

const (
	LabelPositionStart  LabelPosition = "start"
	LabelPositionMiddle LabelPosition = "middle"
	LabelPositionEnd    LabelPosition = "end"
)

// BarProps mirrors @nivo/bar BarCommonProps + BarSvgExtraProps (the fields
// used by the SVG rendering path). Fields with zero values fall back to
// Defaults via UseBar.
type BarProps struct {
	// Data & mapping.
	Data    []BarDatum
	IndexBy core.PropertyAccessor[BarDatum, string]
	Keys    []string

	// Geometry.
	Width  float64
	Height float64
	Margin core.Margin

	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool

	GroupMode    GroupMode
	Layout       Layout
	Padding      float64
	InnerPadding float64

	// Scales. ValueScale accepts any continuous scale spec — a
	// scales.ScaleLinearSpec (the default), scales.ScaleLogSpec, or
	// scales.ScaleSymlogSpec. A nil value falls back to the linear default.
	ValueScale scales.ScaleSpec
	IndexScale scales.ScaleBandSpec

	// Grid.
	EnableGridX bool
	GridXValues []any
	EnableGridY bool
	GridYValues []any

	// Axes.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Colors.
	Colors       colors.OrdinalColorScaleConfig
	ColorBy      ColorBy
	BorderRadius float64
	BorderWidth  float64
	BorderColor  colors.InheritedColorConfig

	// Labels.
	EnableLabel     bool
	Label           core.PropertyAccessor[ComputedDatum, string]
	LabelPosition   LabelPosition
	LabelOffset     float64
	LabelFormat     core.ValueFormat[string]
	LabelSkipWidth  float64
	LabelSkipHeight float64
	LabelTextColor  colors.InheritedColorConfig

	// Value formatting.
	ValueFormat  core.ValueFormat[float64]
	LegendLabel  core.PropertyAccessor[map[string]any, string]
	TooltipLabel core.PropertyAccessor[ComputedDatum, string]

	// Interactivity / a11y.
	Interactive      bool
	InitialHiddenIDs []string

	// Legends.
	Legends []BarLegendProps

	// Annotations.
	Annotations []annotations.AnnotationSpec[ComputedBarDatum]

	// Markers.
	Markers []core.CartesianMarker

	// Defs / fill.
	Defs []core.Def
	Fill []core.DefRule

	// Totals.
	EnableTotals bool
	TotalsOffset float64

	// Layers.
	Layers []BarLayerId

	// Motion.
	core.MotionProps

	// A11y.
	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

	// Theme.
	Theme *theming.Theme

	// HTMX chart instance ID; when non-empty, bars/legends emit hx-* attrs.
	ChartID string

	// HoveredKey, when set by the htmx hover endpoint, marks the bar with this
	// key as active (rendered with reduced opacity on siblings / highlighted
	// border). Cleared on mouseleave.
	HoveredKey string
}

// BarSvgProps is an alias of BarProps (nivo splits common/svg; this unifies them).
type BarSvgProps = BarProps

// BarResult is the output of UseBar: the computed bars, scales, color/label
// accessors, legend data, and totals — everything the templ components need.
type BarResult struct {
	Bars              []ComputedBarDatum
	BarsWithValue     []ComputedBarDatum
	XScale            scales.Scale
	YScale            scales.Scale
	GetIndex          func(BarDatum) string
	GetLabel          func(ComputedDatum) string
	GetTooltipLabel   func(ComputedDatum) string
	FormatValue       func(float64) string
	GetColor          func(ComputedDatum) string
	GetBorderColor    func(any) string
	GetLabelColor     func(any) string
	ShouldRenderLabel func(width, height float64) bool
	HiddenIDs         []string
	LegendsWithData   []LegendWithData
	BarTotals         []BarTotalData
	// ComputeLabelLayout returns the label position/alignment for one bar's
	// (width, height). Set by UseBar from layout/labelPosition/labelOffset/
	// valueScale.reverse.
	ComputeLabelLayout func(width, height float64) compute.BarLabelLayout
}

// LegendWithData pairs a BarLegendProps with its computed LegendData items.
type LegendWithData struct {
	Props BarLegendProps
	Data  []LegendData
}
