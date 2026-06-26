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
	Margin      Margin
	InnerWidth  float64
	InnerHeight float64
	OuterWidth  float64
	OuterHeight float64
}

// Margin is the chart container padding (nivo's defaultMargin = all zero).
type Margin struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// DefaultMargin is nivo's defaultMargin: all zero.
var DefaultMargin = Margin{}

// Box is a 2D rectangle in chart units (used by ComputeArcBoundingBox etc.).
type Box struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
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
	X float64
	Y float64
}

// Padding is uniform per-side padding (used by tooltip/legends layout).
type Padding struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
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
