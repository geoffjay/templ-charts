// Package boxplot mirrors @nivo/boxplot: raw observations grouped by
// group/subGroup, summarized to quantiles (via internal/d3/array.Quantile) and
// drawn as box + whisker + median glyphs on a shared value scale. It reuses
// charts/scales (band index scale via NewBandScaleWithRange + linear value
// scale), charts/axes (grid + axes), charts/colors (ordinal color), charts/legends
// and charts/theming.
//
// v2 scope: SVG only, static render. The interactive hover tooltip arrives with
// the Phase 5 client layer.
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
	Data []BoxPlotDatum

	// Groups / SubGroups give an explicit ordering; nil → first-seen order.
	Groups    []string
	SubGroups []string
	Quantiles []float64

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	Layout BoxPlotLayout

	MinValue *float64
	MaxValue *float64

	Padding      float64
	InnerPadding float64

	Opacity         float64
	ActiveOpacity   float64
	InactiveOpacity float64

	EnableGridX bool
	EnableGridY bool
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	ValueFormat string

	// ColorBy is "subGroup" (default) or "group".
	ColorBy string
	Colors  colors.OrdinalColorScaleConfig

	BorderRadius float64
	BorderWidth  float64
	BorderColor  colors.InheritedColorConfig

	MedianWidth float64
	MedianColor colors.InheritedColorConfig

	WhiskerWidth   float64
	WhiskerColor   colors.InheritedColorConfig
	WhiskerEndSize float64

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []BoxPlotLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	IsFocusable     bool

	core.MotionProps
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

// FloatPtr returns a pointer to f — a helper for *float64 props.
func FloatPtr(f float64) *float64 { return &f }
