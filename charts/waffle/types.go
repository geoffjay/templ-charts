// Package waffle mirrors @nivo/waffle: a part-of-whole chart that fills a grid
// of cells, each cell representing one unit of the total. Built on the
// charts/grid layout primitives (GenerateGrid) plus charts/colors ordinal
// scales, charts/legends, charts/theming, and charts/core.
//
// SVG only, static render. The cells + legends layers are on by
// default; the polygon "areas" layer (one union-outline polygon per datum,
// instead of per-cell rects) is a supported opt-in layer via props.Layers.
package waffle

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// WaffleDatum is one input slice of the whole. Color overrides the ordinal
// scale color when non-empty. Mirrors @nivo/waffle Datum.
type WaffleDatum struct {
	ID    string
	Label string
	Value float64
	Color string
}

// ComputedDatum is an enhanced datum: its color and the half-open cell range
// [StartAt, EndAt) it occupies in fill order. Mirrors @nivo/waffle
// ComputedDatum.
type ComputedDatum struct {
	ID             string
	Label          string
	Value          float64
	FormattedValue string
	GroupIndex     int
	StartAt        int
	EndAt          int
	Color          string
	BorderColor    string
	Hidden         bool
}

// WaffleCell is one positioned grid cell. When HasData is true it belongs to a
// datum (DatumID/Color); otherwise it renders with the empty color.
type WaffleCell struct {
	Key         string
	Index       int
	X           float64
	Y           float64
	Width       float64
	Height      float64
	Color       string
	Opacity     float64
	BorderColor string
	DatumID     string
	Label       string
	HasData     bool
}

// GetX/GetY/GetWidth/GetHeight let WaffleCell satisfy grid.GridCellLike so it
// can be passed to grid.GetCellsPolygons (used by the "areas" layer).
func (c WaffleCell) GetX() float64      { return c.X }
func (c WaffleCell) GetY() float64      { return c.Y }
func (c WaffleCell) GetWidth() float64  { return c.Width }
func (c WaffleCell) GetHeight() float64 { return c.Height }

// WaffleLayerId enumerates the render layers. Mirrors @nivo/waffle layer ids.
// The "areas" layer is a supported opt-in (not in DefaultLayers).
type WaffleLayerId string

const (
	WaffleLayerCells   WaffleLayerId = "cells"
	WaffleLayerAreas   WaffleLayerId = "areas"
	WaffleLayerLegends WaffleLayerId = "legends"
)

// DefaultLayers is the default layer order. The "areas" layer is opt-in and is
// intentionally omitted here so existing renders stay byte-identical.
var DefaultLayers = []WaffleLayerId{WaffleLayerCells, WaffleLayerLegends}

// WaffleLegend configures a discrete legend.
type WaffleLegend struct {
	// Anchor positions the legend relative to the chart area.
	Anchor legends.LegendAnchor
	// Direction lays the items out in a row or column. Defaults to column.
	Direction legends.LegendDirection
	// TranslateX offsets the legend horizontally from its anchor, in pixels.
	TranslateX float64
	// TranslateY offsets the legend vertically from its anchor, in pixels.
	TranslateY float64
	// ItemWidth is the width of each legend item in pixels. Defaults to 100.
	ItemWidth float64
	// ItemHeight is the height of each legend item in pixels. Defaults to 20.
	ItemHeight float64
	// SymbolShape is the marker shape for each item. Defaults to square.
	SymbolShape legends.SymbolShape
}

// WaffleProps is the input to the Waffle component. Mirrors @nivo/waffle
// CommonProps + svgDefaultProps (supported subset).
type WaffleProps struct {
	// Width is the outer chart width in pixels (including Margin).
	Width float64
	// Height is the outer chart height in pixels (including Margin).
	Height float64
	// Margin is the space reserved around the grid (for legends), in pixels.
	Margin core.Margin
	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool
	// Data is the set of slices of the whole; each datum's Value fills a
	// proportional run of cells.
	Data []WaffleDatum

	// Total is the value representing the full grid (rows×columns cells).
	Total float64
	// Rows is the number of cell rows in the grid.
	Rows int
	// Columns is the number of cell columns in the grid.
	Columns int
	// FillDirection is the order cells are filled. Defaults to top (bottom-up).
	FillDirection grid.GridFillDirection
	// Padding is the gap in pixels between adjacent cells. Defaults to 1.
	Padding float64

	// Colors is the ordinal color scale config used to color the data cells
	// (a datum's own Color overrides it when set).
	Colors colors.OrdinalColorScaleConfig
	// EmptyColor is the fill of cells not covered by any datum. Defaults to
	// "#cccccc".
	EmptyColor string
	// EmptyOpacity is the fill-opacity of empty cells, in [0,1]. Defaults to 1.
	EmptyOpacity float64
	// BorderRadius is the corner radius of each cell in pixels. Defaults to 0.
	BorderRadius float64
	// BorderWidth is the stroke width of each cell in pixels; borders are
	// drawn only when > 0. Defaults to 0.
	BorderWidth float64
	// BorderColor is the inherited-color config for the cell border, resolved
	// relative to each cell's fill color.
	BorderColor colors.InheritedColorConfig

	// HiddenIDs lists datum ids to omit from the fill; they occupy no cells and
	// are flagged Hidden in the legend.
	HiddenIDs []string
	// ValueFormat is a d3-format spec for tooltip values; empty leaves values
	// unformatted.
	ValueFormat string

	// Interactive enables the client-side hover layer (charts/interact): each
	// data cell emits a data-tc-tooltip. Default false keeps the static render.
	Interactive bool

	// Animate enables the nivo-style enter animation (opacity fade-in) on each
	// cell. Default false: the rendered SVG is byte-identical to the
	// un-animated output. MotionStagger delays each successive cell by that many
	// seconds (0 = all cells enter together). See core.SMILFadeIn / StaggerBegin.
	Animate       bool
	MotionStagger float64

	// Legends configures zero or more discrete legends drawn from the data.
	Legends []WaffleLegend

	// Theme overrides the styling theme (colors, fonts). Nil uses the package
	// default theme.
	Theme *theming.Theme
	// Layers is the ordered list of render layers. Empty falls back to
	// DefaultLayers (cells, legends); "areas" is an opt-in layer.
	Layers []WaffleLayerId

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
	// IsFocusable makes the svg keyboard-focusable. Defaults to false.
	IsFocusable bool
}

// WaffleResult is the computed model produced by UseWaffle.
type WaffleResult struct {
	Cells        []WaffleCell
	ComputedData []ComputedDatum
	GridX        float64
	GridY        float64
}
