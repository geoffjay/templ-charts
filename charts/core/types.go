// Package core provides foundational types, accessors, formatters, enums,
// SVG defs, and the SvgWrapper templ component shared by every chart package.
//
// It mirrors @nivo/core's non-React surface: types (Dimensions, Margin, Box,
// Point, Padding), PropertyAccessor/GetLabelGenerator, GetValueFormatter,
// the gradient/pattern defs system + BindDefs, the SvgWrapper + DotsItem +
// CartesianMarkers templ components, the curve/stack/blend-mode enums, and
// MotionProps.
package core

// Dimensions records the outer/inner geometry computed from width, height,
// and a Margin by UseDimensions.
type Dimensions struct {
	Margin      Margin  // resolved margin (partialMargin merged over DefaultMargin)
	InnerWidth  float64 // plotting-area width: OuterWidth - Margin.Left - Margin.Right
	InnerHeight float64 // plotting-area height: OuterHeight - Margin.Top - Margin.Bottom
	OuterWidth  float64 // full svg width (the width passed to UseDimensions)
	OuterHeight float64 // full svg height (the height passed to UseDimensions)
}

// Margin is the chart container padding (nivo's defaultMargin = all zero).
type Margin struct {
	Top    float64 // space above the inner plotting area, in px
	Right  float64 // space to the right of the inner plotting area, in px
	Bottom float64 // space below the inner plotting area, in px
	Left   float64 // space to the left of the inner plotting area, in px
}

// DefaultMargin is nivo's defaultMargin: all zero.
var DefaultMargin = Margin{}

// Box is a 2D rectangle in chart units (used by ComputeArcBoundingBox etc.).
type Box struct {
	X      float64 // left edge, in chart units
	Y      float64 // top edge, in chart units
	Width  float64 // horizontal extent, in chart units
	Height float64 // vertical extent, in chart units
}

// BoxAlign enumerates how a child box is positioned inside a parent box.
type BoxAlign string

const (
	BoxAlignCenter      BoxAlign = "center"
	BoxAlignTopLeft     BoxAlign = "top-left"
	BoxAlignTop         BoxAlign = "top"
	BoxAlignTopRight    BoxAlign = "top-right"
	BoxAlignRight       BoxAlign = "right"
	BoxAlignBottomRight BoxAlign = "bottom-right"
	BoxAlignBottom      BoxAlign = "bottom"
	BoxAlignBottomLeft  BoxAlign = "bottom-left"
	BoxAlignLeft        BoxAlign = "left"
)

// Point is a 2D cartesian point (svg units, y-down).
type Point struct {
	X float64 // horizontal position, in svg units
	Y float64 // vertical position, in svg units (y increases downward)
}

// Padding is uniform per-side padding (used by tooltip/legends layout).
type Padding struct {
	Top    float64 // inner padding above the content, in px
	Right  float64 // inner padding to the right of the content, in px
	Bottom float64 // inner padding below the content, in px
	Left   float64 // inner padding to the left of the content, in px
}

// DatumValue is a value that can be retrieved from a datum via a
// PropertyAccessor — number, string, or time.Time.
type DatumValue any

// UseDimensions mirrors nivo's useDimensions hook: merges partialMargin with
// DefaultMargin and returns the inner/outer geometry.
func UseDimensions(width, height float64, partialMargin Margin) Dimensions {
	margin := DefaultMargin
	if partialMargin.Top != 0 {
		margin.Top = partialMargin.Top
	}
	if partialMargin.Right != 0 {
		margin.Right = partialMargin.Right
	}
	if partialMargin.Bottom != 0 {
		margin.Bottom = partialMargin.Bottom
	}
	if partialMargin.Left != 0 {
		margin.Left = partialMargin.Left
	}
	return Dimensions{
		Margin:      margin,
		InnerWidth:  width - margin.Left - margin.Right,
		InnerHeight: height - margin.Top - margin.Bottom,
		OuterWidth:  width,
		OuterHeight: height,
	}
}
