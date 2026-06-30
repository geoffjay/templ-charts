// Package heatmap mirrors @nivo/heatmap: a 2D grid of cells whose color
// encodes a value, laid out with band scales on both axes and colored through
// a continuous (sequential/diverging) color scale. Reuses charts/scales,
// charts/axes, charts/colors, charts/legends, charts/theming, and charts/core.
//
// v2 scope: SVG only, static render (no Canvas, interactivity arrives with the
// Phase 5 client-side layer). Band-scale paddings other than 0 (nivo's
// xInnerPadding/…) are not yet supported.
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
// @nivo/heatmap's `colors` we support in v2.
type HeatMapColorConfig struct {
	Type      string // "sequential" | "diverging"
	Scheme    string
	MinValue  *float64
	MaxValue  *float64
	DivergeAt *float64 // diverging only
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
	Width  float64
	Height float64
	Margin core.Margin
	Data   []HeatMapSerie

	ForceSquare bool

	Colors     HeatMapColorConfig
	EmptyColor string

	EnableGridX bool
	EnableGridY bool
	AxisTop     *axes.AxisProps
	AxisRight   *axes.AxisProps
	AxisBottom  *axes.AxisProps
	AxisLeft    *axes.AxisProps

	Opacity         float64
	ActiveOpacity   float64
	InactiveOpacity float64
	BorderWidth     float64
	BorderColor     colors.InheritedColorConfig
	BorderRadius    float64

	// EnableLabels gates per-cell value labels. nil → true (nivo default);
	// set to a pointer to false to disable. Modeled as a pointer so the
	// default-on behavior survives Go's bool zero value.
	EnableLabels   *bool
	LabelTextColor colors.InheritedColorConfig
	ValueFormat    string // d3-format spec; empty → %g

	Legends []HeatMapLegend

	Theme  *theming.Theme
	Layers []HeatMapLayerId

	// Interactivity (htmx). When IsInteractive is true and ChartID is set,
	// cells with data emit hx-* hover attributes. HoveredKey is the cell id
	// (serieId.x) currently hovered, which drives the active/dim opacity.
	IsInteractive bool
	ChartID       string
	HoveredKey    string

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	IsFocusable     bool

	core.MotionProps
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

// BoolPtr returns a pointer to b — a helper for setting *bool props like
// EnableLabels (e.g. heatmap.BoolPtr(false)).
func BoolPtr(b bool) *bool { return &b }
