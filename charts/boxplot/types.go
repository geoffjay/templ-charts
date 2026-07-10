// Package boxplot mirrors @nivo/boxplot: raw observations grouped by
// group/subGroup, summarized to quantiles (via internal/d3/array.Quantile) and
// drawn as box + whisker + median glyphs on a shared value scale. It reuses
// charts/scales (band index scale via NewBandScaleWithRange + linear value
// scale), charts/axes (grid + axes), charts/colors (ordinal color), charts/legends
// and charts/theming.
//
// SVG only, static render. The interactive hover tooltip arrives with
// the charts/interact client layer.
package boxplot

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// BoxPlotLayout is the orientation of the boxes.
type BoxPlotLayout string

const (
	BoxPlotLayoutVertical   BoxPlotLayout = "vertical"
	BoxPlotLayoutHorizontal BoxPlotLayout = "horizontal"
)

// BoxPlotDatum is one raw observation: its group, optional sub-group, and
// numeric value. Mirrors @nivo/boxplot BoxPlotDatum.
type BoxPlotDatum struct {
	Group    string
	SubGroup string
	Value    float64
}

// BoxPlotSummary is the quantile summary of one group/subGroup distribution.
// Mirrors @nivo/boxplot BoxPlotSummary.
type BoxPlotSummary struct {
	Group         string
	SubGroup      string
	GroupIndex    int
	SubGroupIndex int
	N             int
	Extrema       [2]float64 // min, max
	Quantiles     []float64  // the quantile fractions (e.g. 0.1..0.9)
	Values        []float64  // quantile values
	Mean          float64
}

// ComputedBox is one fully laid-out box-plot glyph (SVG-ready, in the item's
// local coordinate frame plus a group transform). Mirrors the subset of
// @nivo/boxplot ComputedBoxPlotSummary the renderer needs.
type ComputedBox struct {
	Key          string
	Group        string
	SubGroup     string
	Summary      BoxPlotSummary
	Transform    string
	Bandwidth    float64
	Color        string
	BorderColor  string
	MedianColor  string
	WhiskerColor string
	Opacity      float64
	// Local-frame distances from the median (values[2]) in pixels.
	RectY              float64 // vertical → VD3, horizontal → VD1
	ValueInterval      float64
	VD0, VD1, VD3, VD4 float64
	WhiskerEnd         float64
}

// BoxPlotLayerId enumerates the render layers. Mirrors @nivo/boxplot LayerId.
type BoxPlotLayerId string

const (
	BoxPlotLayerGrid        BoxPlotLayerId = "grid"
	BoxPlotLayerAxes        BoxPlotLayerId = "axes"
	BoxPlotLayerBoxPlots    BoxPlotLayerId = "boxPlots"
	BoxPlotLayerMarkers     BoxPlotLayerId = "markers"
	BoxPlotLayerLegends     BoxPlotLayerId = "legends"
	BoxPlotLayerAnnotations BoxPlotLayerId = "annotations"
)

// DefaultLayers mirrors @nivo/boxplot svgDefaultProps.layers.
var DefaultLayers = []BoxPlotLayerId{
	BoxPlotLayerGrid, BoxPlotLayerAxes, BoxPlotLayerBoxPlots, BoxPlotLayerMarkers, BoxPlotLayerLegends, BoxPlotLayerAnnotations,
}

// BoxPlotProps mirrors @nivo/boxplot BoxPlotSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type BoxPlotProps struct {
	// Data holds the raw observations; each datum is bucketed by Group/SubGroup
	// and summarized to quantiles.
	Data []BoxPlotDatum

	// Groups / SubGroups give an explicit ordering; nil → first-seen order.
	Groups    []string
	SubGroups []string
	// Quantiles are the quantile fractions computed per distribution; nil →
	// [0.1,0.25,0.5,0.75,0.9]. Five values are required to render a glyph.
	Quantiles []float64

	// Width and Height are the total SVG dimensions in pixels; the plot area is
	// these minus Margin.
	Width  float64
	Height float64
	// Margin is the space reserved around the plot area (for axes/legends).
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// Layout orients the boxes "vertical" (default) or "horizontal".
	Layout BoxPlotLayout

	// Interactive enables the client-side hover layer (charts/interact): each
	// box glyph emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// box rect. Default false: the rendered SVG is byte-identical to the
	// un-animated output. MotionStagger delays each successive box by that many
	// seconds (0 = all enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

	// MinValue and MaxValue override the value-axis domain; nil → derived from
	// the data's quantile values.
	MinValue *float64
	MaxValue *float64

	// Padding is the fraction (0..1) of the band scale left as gaps between
	// groups; default 0.1. InnerPadding is the gap in pixels between sub-group
	// boxes within a group; default 6.
	Padding      float64
	InnerPadding float64

	// Opacity is the fill opacity of each box; default 1. ActiveOpacity and
	// InactiveOpacity are the opacities of hovered vs. non-hovered boxes when
	// Interactive (defaults 1 and 0.25).
	Opacity         float64
	ActiveOpacity   float64
	InactiveOpacity float64

	EnableGridX *bool // nil → false (nivo default)
	EnableGridY *bool // nil → true (nivo default)
	// AxisTop/Right/Bottom/Left configure each axis; nil hides that axis. When
	// all four are nil, bottom and left axes are shown by default.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// ValueFormat is a format spec applied to displayed values (e.g. tooltips).
	ValueFormat string

	// ColorBy is "subGroup" (default) or "group".
	ColorBy string
	// Colors is the ordinal color scale mapping the ColorBy key to a box fill;
	// default the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig

	// BorderRadius is the corner radius of the box rect in pixels; default 0.
	BorderRadius float64
	// BorderWidth is the box stroke width in pixels; default 0 (no border).
	BorderWidth float64
	// BorderColor resolves the box border color; default inherits the box color.
	BorderColor colors.InheritedColorConfig

	// MedianWidth is the median line stroke width in pixels; default 2.
	MedianWidth float64
	// MedianColor resolves the median line color; default the box color darkened.
	MedianColor colors.InheritedColorConfig

	// WhiskerWidth is the whisker line stroke width in pixels; default 2.
	WhiskerWidth float64
	// WhiskerColor resolves the whisker color; default inherits the box color.
	WhiskerColor colors.InheritedColorConfig
	// WhiskerEndSize is the whisker cap width as a fraction (0..1) of the box
	// width; default 0.6.
	WhiskerEndSize float64

	// Legends configures zero or more legends; empty → no legend.
	Legends []legends.LegendProps

	// Theme overrides chart styling; nil → theming.DefaultTheme.
	Theme *theming.Theme
	// Layers sets the render order of chart layers; empty → DefaultLayers.
	Layers []BoxPlotLayerId

	// Role is the root SVG ARIA role; default "img".
	Role string
	// AriaLabel, AriaLabelledBy, and AriaDescribedBy set the corresponding SVG
	// accessibility attributes.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title>/<desc> elements.
	Title string
	Desc  string
	// IsFocusable makes the SVG keyboard-focusable.
	IsFocusable bool
}

// BoxPlotResult is the computed model produced by UseBoxPlot.
type BoxPlotResult struct {
	Boxes      []ComputedBox
	IndexScale scales.Scale
	ValueScale scales.Scale
	XScale     scales.Scale
	YScale     scales.Scale
	LegendData []legends.Datum
}

// GridXEnabled resolves EnableGridX (nil → false, nivo default).
func (p BoxPlotProps) GridXEnabled() bool { return p.EnableGridX != nil && *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p BoxPlotProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }
