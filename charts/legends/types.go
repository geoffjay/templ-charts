// Package legends mirrors @nivo/legends: LegendProps, the LegendAnchor /
// LegendDirection / LegendItemDirection enums, the Datum type, the
// BoxLegendSvg / LegendSvg / LegendSvgItem / ContinuousColorsLegendSvg templ
// components, the ComputeDimensions / ComputePositionFromAnchor /
// ComputeItemLayout / ComputeContinuousColorsLegend helpers, the SymbolShape
// factories, and the legend defaults.
package legends

// LegendAnchor is which corner/side of the chart the legend box attaches to.
type LegendAnchor string

const (
	LegendAnchorTopLeft     LegendAnchor = "top-left"
	LegendAnchorTop         LegendAnchor = "top"
	LegendAnchorTopRight    LegendAnchor = "top-right"
	LegendAnchorRight       LegendAnchor = "right"
	LegendAnchorBottomRight LegendAnchor = "bottom-right"
	LegendAnchorBottom      LegendAnchor = "bottom"
	LegendAnchorBottomLeft  LegendAnchor = "bottom-left"
	LegendAnchorLeft        LegendAnchor = "left"
)

// LegendDirection is the layout direction of legend items (row/column).
type LegendDirection string

const (
	LegendDirectionRow    LegendDirection = "row"
	LegendDirectionColumn LegendDirection = "column"
)

// LegendItemDirection is the layout direction within a single legend item
// (symbol → label).
type LegendItemDirection string

const (
	LegendItemDirectionLeftToRight LegendItemDirection = "left-to-right"
	LegendItemDirectionRightToLeft LegendItemDirection = "right-to-left"
	LegendItemDirectionTopToBottom LegendItemDirection = "top-to-bottom"
	LegendItemDirectionBottomToTop LegendItemDirection = "bottom-to-top"
)

// SymbolShape enumerates the legend symbol shapes.
type SymbolShape string

const (
	SymbolShapeCircle   SymbolShape = "circle"
	SymbolShapeDiamond  SymbolShape = "diamond"
	SymbolShapeSquare   SymbolShape = "square"
	SymbolShapeTriangle SymbolShape = "triangle"
)

// Datum is one legend entry: a label, color, and hidden flag.
type Datum struct {
	ID     string
	Label  string
	Color  string
	Hidden bool
}

// LegendProps mirrors @nivo/legends BoxLegendSvg props (SVG subset).
type LegendProps struct {
	Anchor            LegendAnchor
	Direction         LegendDirection
	Items             []Datum
	TranslateX        float64
	TranslateY        float64
	ItemDirection     LegendItemDirection
	ItemWidth         float64
	ItemHeight        float64
	ItemOpacity       float64
	SymbolShape       SymbolShape
	SymbolSize        float64
	SymbolSpacing     float64
	SymbolBorderColor string
	SymbolBorderWidth float64
	// Toggle is the HTMX verb for series toggle (emitted on click). When
	// non-empty, each item emits an hx-post attribute.
	Toggle string
	// ChartID + ItemIDs back-reference the htmx registry instance.
	ChartID string
}

// LegendDefaults mirrors @nivo/legends legendDefaults.
var LegendDefaults = LegendProps{
	Anchor:            LegendAnchorTopLeft,
	Direction:         LegendDirectionColumn,
	ItemDirection:     LegendItemDirectionLeftToRight,
	ItemWidth:         100,
	ItemHeight:        20,
	ItemOpacity:       1,
	SymbolShape:       SymbolShapeCircle,
	SymbolSize:        18,
	SymbolSpacing:     6,
	SymbolBorderColor: "#ffffff",
	SymbolBorderWidth: 0,
}

// ContinuousColorsLegendProps mirrors @nivo/legends ContinuousColorsLegendSvg.
type ContinuousColorsLegendProps struct {
	Scale      func(float64) string
	Min        float64
	Max        float64
	Width      float64
	Height     float64
	Anchor     LegendAnchor
	TranslateX float64
	TranslateY float64
	Title      string
	// Number of discrete color samples along the gradient bar.
	Samples int
}

// ContinuousColorsLegendDefaults mirrors @nivo/legends
// continuousColorsLegendDefaults.
var ContinuousColorsLegendDefaults = ContinuousColorsLegendProps{
	Anchor:  LegendAnchorTopRight,
	Width:   100,
	Height:  10,
	Samples: 32,
}
