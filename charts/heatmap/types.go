// Package heatmap mirrors @nivo/heatmap: a 2D grid of cells whose color
// encodes a value, laid out with band scales on both axes and colored through
// a continuous (sequential/diverging) color scale. Reuses charts/scales,
// charts/axes, charts/colors, charts/legends, charts/theming, and charts/core.
//
// SVG by default with an opt-in Canvas backend (Render: theming.EngineCanvas)
// for large-N grids; hover tooltips come from the charts/interact client layer
// via Interactive. Band-scale paddings other than 0 (nivo's xInnerPadding/…)
// are not yet supported.
package heatmap

import (
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// HeatMapDatum is one cell's raw data: an x category and a nullable value.
// A nil Y renders as the empty cell color. Mirrors @nivo/heatmap datum.
type HeatMapDatum struct {
	X string
	Y *float64
}

// HeatMapSerie is one row of the heatmap: an id (the y-axis category) and its
// per-x cells. Mirrors @nivo/heatmap serie.
type HeatMapSerie struct {
	ID   string
	Data []HeatMapDatum
}

// ComputedCell is a positioned, colored cell ready to render. X/Y are the cell
// center (matching nivo). Mirrors @nivo/heatmap ComputedCell.
type ComputedCell struct {
	ID             string
	SerieID        string
	X              string // x category
	Value          *float64
	FormattedValue string
	XPos           float64 // center x
	YPos           float64 // center y
	Width          float64
	Height         float64
	Color          string
	Opacity        float64
	BorderColor    string
	Label          string
	LabelTextColor string
}

// HeatMapColorConfig selects the continuous color scale used to map cell
// values to colors. Type is "sequential" (default) or "diverging". Scheme is
// an interpolator id (default "brown_blueGreen"). Mirrors the subset of
// @nivo/heatmap's `colors` we support.
type HeatMapColorConfig struct {
	Type      string // "sequential" | "diverging"
	Scheme    string
	MinValue  *float64
	MaxValue  *float64
	DivergeAt *float64 // diverging only
	// Space selects the interpolation color space for the scale (additive; the
	// zero value colors.SpaceRGB reproduces today's gamma-sRGB interpolation
	// byte-for-byte). Set colors.SpaceLab/SpaceLch for perceptually-uniform
	// cell coloring.
	Space colors.Space
}

// HeatMapLayerId enumerates the render layers. Mirrors @nivo/heatmap LayerId.
type HeatMapLayerId string

const (
	HeatMapLayerGrid        HeatMapLayerId = "grid"
	HeatMapLayerAxes        HeatMapLayerId = "axes"
	HeatMapLayerCells       HeatMapLayerId = "cells"
	HeatMapLayerLegends     HeatMapLayerId = "legends"
	HeatMapLayerAnnotations HeatMapLayerId = "annotations"
)

// DefaultLayers mirrors @nivo/heatmap commonDefaultProps.layers.
var DefaultLayers = []HeatMapLayerId{
	HeatMapLayerGrid, HeatMapLayerAxes, HeatMapLayerCells, HeatMapLayerLegends, HeatMapLayerAnnotations,
}

// HeatMapLegend configures a continuous color legend. Mirrors the heatmap
// continuous legend subset.
type HeatMapLegend struct {
	Anchor     legends.LegendAnchor
	TranslateX float64
	TranslateY float64
	Length     float64
	Thickness  float64
	Title      string
}

// HeatMapProps is the input to the HeatMap component. Mirrors @nivo/heatmap
// HeatMapCommonProps + svgDefaultProps (fields not yet supported are omitted).
type HeatMapProps struct {
	// Width is the total chart width in pixels (including Margin).
	Width float64
	// Height is the total chart height in pixels (including Margin).
	Height float64
	// Margin reserves space around the grid for axes and legends.
	Margin core.Margin
	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool
	// Data is the list of series (rows); each serie is a y-axis category with
	// its per-x cells.
	Data []HeatMapSerie

	// ForceSquare forces cells to be square, sizing the grid to the smaller
	// axis. Default false.
	ForceSquare bool

	// Colors is the continuous color scale mapping cell values to colors (see
	// HeatMapColorConfig). Default is a sequential "brown_blueGreen" scale.
	Colors HeatMapColorConfig
	// EmptyColor is the CSS color used for cells with a nil value. Default
	// "#000000".
	EmptyColor string

	EnableGridX *bool // nil → false (nivo default)
	EnableGridY *bool // nil → false (nivo default)
	// AxisTop configures the top axis; nil hides it.
	AxisTop *axes.AxisProps
	// AxisRight configures the right axis; nil hides it.
	AxisRight *axes.AxisProps
	// AxisBottom configures the bottom axis; nil hides it.
	AxisBottom *axes.AxisProps
	// AxisLeft configures the left axis; nil hides it.
	AxisLeft *axes.AxisProps

	// Opacity is the base opacity of each cell, from 0 to 1. Default 1.
	Opacity float64
	// ActiveOpacity is the opacity of the hovered/active cell, from 0 to 1.
	// Default 1.
	ActiveOpacity float64
	// InactiveOpacity is the opacity of the non-active cells while another is
	// hovered, from 0 to 1. Default 0.15.
	InactiveOpacity float64
	// BorderWidth is the stroke width of each cell's border, in pixels. Default
	// 0 (no border drawn).
	BorderWidth float64
	// BorderColor resolves each cell's border color, by default the cell fill
	// darkened. Only applied when BorderWidth > 0.
	BorderColor colors.InheritedColorConfig
	// BorderRadius is the corner radius of the cells, in pixels. Default 0.
	BorderRadius float64

	// EnableLabels gates per-cell value labels. nil → true (nivo default);
	// set to a pointer to false to disable. Modeled as a pointer so the
	// default-on behavior survives Go's bool zero value.
	EnableLabels *bool
	// LabelTextColor resolves the per-cell label color, by default the cell
	// fill darkened.
	LabelTextColor colors.InheritedColorConfig
	ValueFormat    string // d3-format spec; empty → %g

	// Legends configures the continuous color legends to draw.
	Legends []HeatMapLegend

	// Theme overrides the styling theme (colors, fonts, label styles). Nil uses
	// the default theme.
	Theme *theming.Theme
	// Layers selects which render layers to draw and their order (grid, axes,
	// cells, legends, annotations). Default DefaultLayers.
	Layers []HeatMapLayerId

	// Render selects the backend: the zero value (or theming.EngineSVG) renders
	// SVG as always; theming.EngineCanvas draws the cells into a <canvas>
	// draw-list (charts/canvas) with grid/axes/legends kept as SVG panes. Default
	// SVG keeps every existing golden byte-stable.
	Render theming.Engine

	// Interactivity (htmx). When Interactive is true and ChartID is set,
	// cells with data emit hx-* hover attributes. HoveredKey is the cell id
	// (serieId.x) currently hovered, which drives the active/dim opacity.
	Interactive bool
	ChartID     string
	HoveredKey  string

	// Role is the ARIA role of the root svg element. Default "img".
	Role string
	// AriaLabel sets aria-label on the root svg element. Empty by default.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the root svg element. Empty by default.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the root svg element. Empty by default.
	AriaDescribedBy string
	// Title sets the svg <title> element. Empty by default.
	Title string
	// Desc sets the svg <desc> element. Empty by default.
	Desc string
	// IsFocusable makes the root svg keyboard-focusable. Defaults to false.
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// cell (600ms). MotionStagger delays successive cells by that many seconds
	// (0 ⇒ all cells enter together). Defaults off, so static output is
	// unchanged. Mirrors @nivo/heatmap's enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// HeatMapResult is the computed model produced by UseHeatMap and consumed by
// the render layers.
type HeatMapResult struct {
	Cells    []ComputedCell
	XScale   scales.Scale
	YScale   scales.Scale
	XValues  []any
	SerieIDs []any
	OffsetX  float64
	OffsetY  float64
	Width    float64 // laid-out grid width (≤ inner width when ForceSquare)
	Height   float64
	MinValue float64
	MaxValue float64
	// ColorScale maps a value to a color (for the continuous legend).
	ColorScale func(float64) string
}

// LabelsEnabled resolves EnableLabels: unset (nil) means true (nivo default).
func (p HeatMapProps) LabelsEnabled() bool {
	return p.EnableLabels == nil || *p.EnableLabels
}

// GridXEnabled resolves EnableGridX (nil → false, nivo default).
func (p HeatMapProps) GridXEnabled() bool { return p.EnableGridX != nil && *p.EnableGridX }

// GridYEnabled resolves EnableGridY (nil → false, nivo default).
func (p HeatMapProps) GridYEnabled() bool { return p.EnableGridY != nil && *p.EnableGridY }
