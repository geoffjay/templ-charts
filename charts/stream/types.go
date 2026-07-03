// Package stream mirrors @nivo/stream: stacked areas over an index axis with a
// configurable stack offset (wiggle/silhouette/expand/diverging/none) and a
// smooth curve. It reuses charts/scales (point x-scale + linear y-scale),
// charts/axes (grid + axes), charts/colors (ordinal color per layer), the
// internal/d3/shape Stack + Area generators, charts/legends and charts/theming.
//
// v2 scope: SVG only, static render. The interactive slices/dots layers arrive
// with the Phase 5 client layer; dots default to off in nivo anyway.
package stream

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// StreamDatum is one index row: a map of layer key → value. Mirrors
// @nivo/stream StreamDatum.
type StreamDatum = map[string]float64

// StreamLayerId enumerates the render layers. Mirrors @nivo/stream
// StreamLayerId.
type StreamLayerId string

const (
	StreamLayerGrid    StreamLayerId = "grid"
	StreamLayerAxes    StreamLayerId = "axes"
	StreamLayerLayers  StreamLayerId = "layers"
	StreamLayerDots    StreamLayerId = "dots"
	StreamLayerSlices  StreamLayerId = "slices"
	StreamLayerLegends StreamLayerId = "legends"
)

// DefaultLayers mirrors @nivo/stream svgDefaultProps.layers.
var DefaultLayers = []StreamLayerId{
	StreamLayerGrid, StreamLayerAxes, StreamLayerLayers, StreamLayerDots, StreamLayerSlices, StreamLayerLegends,
}

// ComputedLayer is one stacked area layer ready to render. Mirrors
// @nivo/stream StreamLayerData (the subset needed for SVG areas + legends).
type ComputedLayer struct {
	ID          string
	Label       string
	Path        string
	Color       string
	BorderColor string
}

// StreamProps mirrors @nivo/stream StreamSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type StreamProps struct {
	Data []StreamDatum
	Keys []string

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	OffsetType  core.StackOffset
	Order       core.StackOrder
	Curve       core.CurveFactoryId
	ValueFormat string

	// Interactive enables the client-side hover layer (charts/interact): each
	// area layer emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	Colors      colors.OrdinalColorScaleConfig
	FillOpacity float64
	BorderWidth float64
	BorderColor colors.InheritedColorConfig

	EnableGridX bool
	EnableGridY bool
	GridXValues []any
	GridYValues []any
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []StreamLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// StreamResult is the computed model produced by UseStream.
type StreamResult struct {
	Layers     []ComputedLayer
	XScale     scales.Scale
	YScale     scales.Scale
	LegendData []legends.Datum
}
