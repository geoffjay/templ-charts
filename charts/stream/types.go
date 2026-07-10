// Package stream mirrors @nivo/stream: stacked areas over an index axis with a
// configurable stack offset (wiggle/silhouette/expand/diverging/none) and a
// smooth curve. It reuses charts/scales (point x-scale + linear y-scale),
// charts/axes (grid + axes), charts/colors (ordinal color per layer), the
// internal/d3/shape Stack + Area generators, charts/legends and charts/theming.
//
// SVG only, static render. The interactive slices/dots layers arrive
// with the charts/interact client layer; dots default to off in nivo anyway.
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
	// Data holds one row per index position; each row maps a layer key to its
	// value at that index.
	Data []StreamDatum
	// Keys names the layers (data keys) to stack, in stacking order.
	Keys []string

	// Width and Height are the overall SVG dimensions in pixels.
	Width  float64
	Height float64
	// Margin reserves space around the plot area (top/right/bottom/left) in
	// pixels, e.g. for axes and legends.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container; see
	// core.SvgWrapperProps.Responsive.
	Responsive bool

	// OffsetType is the stack offset applied to the layers
	// (wiggle/silhouette/expand/diverging/none). Default wiggle.
	OffsetType core.StackOffset
	// Order is the stacking order of the layers. Default none (input order).
	Order core.StackOrder
	// Curve is the interpolation used for the area outlines. Default
	// catmullRom.
	Curve core.CurveFactoryId
	// ValueFormat is a d3-format spec applied to values (e.g. in tooltips);
	// empty → %g.
	ValueFormat string

	// Interactive enables the client-side hover layer (charts/interact): each
	// area layer emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Colors is the ordinal color scale mapping each layer to a color. Default
	// is the "nivo" scheme.
	Colors colors.OrdinalColorScaleConfig
	// FillOpacity is the fill opacity of each layer area, in [0,1]. Default 1.
	FillOpacity float64
	// BorderWidth is the stroke width of each layer outline, in pixels.
	// Default 0.
	BorderWidth float64
	// BorderColor is the layer border color, resolved via the inherited-color
	// system. Default derives from each layer's fill darkened by 1.0.
	BorderColor colors.InheritedColorConfig

	EnableGridX *bool // nil → false (nivo default)
	EnableGridY *bool // nil → true (nivo default)
	// GridXValues and GridYValues override the tick positions of the x/y grid
	// lines; nil lets the scale choose them.
	GridXValues []any
	GridYValues []any
	// AxisTop, AxisRight, AxisBottom and AxisLeft configure the four axes; a nil
	// pointer hides that axis.
	AxisTop    *axes.AxisProps
	AxisRight  *axes.AxisProps
	AxisBottom *axes.AxisProps
	AxisLeft   *axes.AxisProps

	// Legends configures the chart legends; empty means no legend.
	Legends []legends.LegendProps

	// Theme overrides the styling theme; nil uses the default theme.
	Theme *theming.Theme
	// Layers selects which layers are rendered and in what order. Defaults to
	// DefaultLayers.
	Layers []StreamLayerId

	// Role is the SVG root's ARIA role. Default "img".
	Role string
	// AriaLabel, AriaLabelledBy and AriaDescribedBy set the matching ARIA
	// attributes on the SVG root for accessibility.
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	// Title and Desc set the SVG <title> and <desc> elements for accessibility.
	Title string
	Desc  string
	// IsFocusable, when true, makes the SVG root keyboard-focusable.
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// stream layer area (600ms). MotionStagger delays successive layers by that
	// many seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/stream's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// StreamResult is the computed model produced by UseStream.
type StreamResult struct {
	Layers     []ComputedLayer
	XScale     scales.Scale
	YScale     scales.Scale
	LegendData []legends.Datum
}

// GridXEnabled resolves EnableGridX (nil → false, nivo default).
func (p StreamProps) GridXEnabled() bool { return p.EnableGridX != nil && *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → true, nivo default).
func (p StreamProps) GridYEnabled() bool { return p.EnableGridY == nil || *p.EnableGridY }
