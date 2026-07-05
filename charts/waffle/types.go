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
	Anchor      legends.LegendAnchor
	Direction   legends.LegendDirection
	TranslateX  float64
	TranslateY  float64
	ItemWidth   float64
	ItemHeight  float64
	SymbolShape legends.SymbolShape
}

// WaffleProps is the input to the Waffle component. Mirrors @nivo/waffle
// CommonProps + svgDefaultProps (supported subset).
type WaffleProps struct {
	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the rendered svg scale fluidly to its container
	// (viewBox preserved, width:100%;height:auto) instead of a fixed pixel
	// size. See core.SvgWrapperProps.Responsive.
	Responsive bool
	Data       []WaffleDatum

	// Total is the value representing the full grid (rows×columns cells).
	Total         float64
	Rows          int
	Columns       int
	FillDirection grid.GridFillDirection
	Padding       float64

	Colors       colors.OrdinalColorScaleConfig
	EmptyColor   string
	EmptyOpacity float64
	BorderRadius float64
	BorderWidth  float64
	BorderColor  colors.InheritedColorConfig

	HiddenIDs   []string
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

	Legends []WaffleLegend

	Theme  *theming.Theme
	Layers []WaffleLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// WaffleResult is the computed model produced by UseWaffle.
type WaffleResult struct {
	Cells        []WaffleCell
	ComputedData []ComputedDatum
	GridX        float64
	GridY        float64
}
