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
	"github.com/geoffjay/templ-charts/charts/bar/internal/compute"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// BarDatum is a single row of bar-chart data: a map from key to value
// (string | number). Mirrors @nivo/bar BarDatum = Record<string, string|number>.
// It is intentionally an open, key-indexed map (not a fixed struct): rows are
// addressed dynamically by IndexBy and Keys, so callers can supply arbitrary
// series without a predeclared schema.
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

	// Data is the set of rows to plot; each BarDatum is a map keyed by the
	// IndexBy field and one entry per key in Keys. Required.
	Data []BarDatum
	// IndexBy selects the datum field used to index (group and label) each
	// bar/stack along the index axis. Default "id".
	IndexBy core.PropertyAccessor[BarDatum, string]
	// Keys lists the datum fields plotted as series; each should point to a
	// numeric value. Default ["value"].
	Keys []string

	// Geometry.

	// Width is the outer chart width in pixels (including Margin).
	Width float64
	// Height is the outer chart height in pixels (including Margin).
	Height float64
	// Margin is the space reserved around the plot area for axes and legends.
	Margin core.Margin

	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool

	// GroupMode controls how bars for multiple keys are combined:
	// "stacked" (default) or "grouped" (side by side).
	GroupMode GroupMode
	// Layout is the bar orientation: "vertical" (default) or "horizontal".
	Layout Layout
	// Padding is the spacing between each index band, as a ratio of the band
	// width (0–0.9). Default 0.1.
	Padding float64
	// InnerPadding is the spacing in pixels between grouped/stacked bars
	// within a single index. Default 0.
	InnerPadding float64

	// Scales. ValueScale accepts any continuous scale spec — a
	// scales.ScaleLinearSpec (the default), scales.ScaleLogSpec, or
	// scales.ScaleSymlogSpec. A nil value falls back to the linear default.
	ValueScale scales.ScaleSpec
	// IndexScale is the band scale spec applied to the index axis (controls
	// band rounding). Default ScaleBandSpec{Round: false}.
	IndexScale scales.ScaleBandSpec

	// Grid. nil → use default (X false, Y true).
	EnableGridX *bool
	// GridXValues optionally overrides the tick values at which vertical
	// gridlines are drawn; nil uses the value/index scale ticks.
	GridXValues []any
	EnableGridY *bool
	// GridYValues optionally overrides the tick values at which horizontal
	// gridlines are drawn; nil uses the value/index scale ticks.
	GridYValues []any

	// Axes. A nil Axis* disables that side's axis.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Colors.

	// Colors is the ordinal color scale used to color bars. Default is the
	// "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// ColorBy selects the property that drives bar color: "id" (default, color
	// per series key) or "indexValue" (color per index).
	ColorBy ColorBy
	// BorderRadius is the bar corner radius in pixels. Default 0.
	BorderRadius float64
	// BorderWidth is the bar border width in pixels. Default 0.
	BorderWidth float64
	// BorderColor computes each bar's border color. Default derives it from
	// the bar's own fill color.
	BorderColor colors.InheritedColorConfig

	// Labels.
	EnableLabel *bool // nil → true (nivo default)
	// Label defines how each bar's label text is computed from its
	// ComputedDatum. Default "formattedValue" (the formatted value).
	Label core.PropertyAccessor[ComputedDatum, string]
	// LabelPosition places the label relative to its bar: "start", "middle"
	// (default), or "end".
	LabelPosition LabelPosition
	// LabelOffset shifts the label along the bar's value axis (vertical or
	// horizontal depending on Layout), in pixels (-16–16). Default 0.
	LabelOffset float64
	// LabelFormat optionally formats the label string after it is computed.
	LabelFormat core.ValueFormat[string]
	// LabelSkipWidth hides a bar's label when the bar is narrower than this
	// many pixels; ignored when 0. Default 0.
	LabelSkipWidth float64
	// LabelSkipHeight hides a bar's label when the bar is shorter than this
	// many pixels; ignored when 0. Default 0.
	LabelSkipHeight float64
	// LabelTextColor computes the label text color. Default is the theme's
	// labels.text.fill.
	LabelTextColor colors.InheritedColorConfig

	// Value formatting.

	// ValueFormat optionally formats numeric values for labels, tooltips, and
	// totals. nil leaves values unformatted.
	ValueFormat core.ValueFormat[float64]
	// LegendLabel computes the text for each legend item. Default is the item
	// "id" (dataFrom "keys") or "indexValue" (dataFrom "indexes").
	LegendLabel core.PropertyAccessor[map[string]any, string]
	// TooltipLabel computes the tooltip label for a bar. Default is
	// "{id} - {indexValue}".
	TooltipLabel core.PropertyAccessor[ComputedDatum, string]

	// Interactivity / a11y.

	// Interactive enables hover/HTMX interactivity on bars and legends.
	// Default true.
	Interactive bool
	// InitialHiddenIDs lists series keys hidden on first render (toggleable
	// via the legend). Default empty.
	InitialHiddenIDs []string

	// Legends is the set of legends to render for the chart.
	Legends []BarLegendProps

	// Annotations are the annotations overlaid on matching bars.
	Annotations []annotations.AnnotationSpec[ComputedBarDatum]

	// Markers are reference lines/regions drawn on the cartesian plane.
	Markers []core.CartesianMarker

	// Defs / fill.

	// Defs declares reusable SVG paint definitions (gradients/patterns)
	// referenced by Fill.
	Defs []core.Def
	// Fill maps Defs to bars by matcher, overriding their solid fill.
	Fill []core.DefRule

	// Totals. nil → false (nivo default).
	EnableTotals *bool
	// TotalsOffset is the gap in pixels between a bar/stack edge and its total
	// label (0–40). Default 10.
	TotalsOffset float64

	// Layers is the render order of the built-in layers. Default DefaultLayers.
	Layers []BarLayerId

	// Motion.
	core.MotionProps

	// A11y.

	// Role is the ARIA role of the root svg element. Default "img".
	Role string
	// AriaLabel sets aria-label on the root svg element.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the root svg element.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the root svg element.
	AriaDescribedBy string
	// Title sets the svg <title> element text.
	Title string
	// Desc sets the svg <desc> element text.
	Desc string
	// IsFocusable makes the root svg and each bar item focusable for keyboard
	// navigation. Default false.
	IsFocusable bool

	// Theme overrides the chart theme; nil uses the default theme.
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

// GridXEnabled resolves EnableGridX (nil → false, nivo default).
func (p BarProps) GridXEnabled() bool { return p.EnableGridX != nil && *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p BarProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }

// LabelEnabled resolves EnableLabel (nil → true, nivo default).
func (p BarProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }

// TotalsEnabled resolves EnableTotals (nil → false, nivo default).
func (p BarProps) TotalsEnabled() bool { return p.EnableTotals != nil && *p.EnableTotals }
