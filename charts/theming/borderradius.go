package theming

import (
	"math"
	"strconv"
)

// BorderRadiusCorners holds the four explicit corner radii.
type BorderRadiusCorners struct {
	TopLeft, TopRight, BottomRight, BottomLeft float64
}

// BorderRadiusObject is the object form: any combination of top/bottom/left/
// right groups + explicit corner values.
type BorderRadiusObject struct {
	Top         *float64
	Bottom      *float64
	Left        *float64
	Right       *float64
	TopLeft     *float64
	TopRight    *float64
	BottomRight *float64
	BottomLeft  *float64
}

// BorderRadius is either a uniform number or an object mixing group+corner
// values. The Uniform field, when non-nil, indicates a uniform number;
// otherwise Object is consulted. Modeled as a struct so Go can express the
// union cleanly.
type BorderRadius struct {
	Uniform *float64
	Object  BorderRadiusObject
}

// NewUniformBorderRadius returns a uniform BorderRadius.
func NewUniformBorderRadius(v float64) BorderRadius {
	return BorderRadius{Uniform: &v}
}

// NormalizeBorderRadius resolves a BorderRadius into explicit corner values.
// Priority: uniform → explicit corner → group (top/bottom/left/right) → 0.
// Mirrors @nivo/theming normalizeBorderRadius.
func NormalizeBorderRadius(r BorderRadius) BorderRadiusCorners {
	if r.Uniform != nil {
		v := *r.Uniform
		return BorderRadiusCorners{TopLeft: v, TopRight: v, BottomRight: v, BottomLeft: v}
	}
	o := r.Object
	uniform := 0.0
	tl := orFallback(o.TopLeft, orGroup(o.Top, o.Left, uniform))
	tr := orFallback(o.TopRight, orGroup(o.Top, o.Right, uniform))
	br := orFallback(o.BottomRight, orGroup(o.Bottom, o.Right, uniform))
	bl := orFallback(o.BottomLeft, orGroup(o.Bottom, o.Left, uniform))
	return BorderRadiusCorners{TopLeft: tl, TopRight: tr, BottomRight: br, BottomLeft: bl}
}

func orFallback(v *float64, fallback float64) float64 {
	if v != nil {
		return *v
	}
	return fallback
}

func orGroup(a, b *float64, fallback float64) float64 {
	if a != nil {
		return *a
	}
	if b != nil {
		return *b
	}
	return fallback
}

// ConstrainBorderRadius adjusts corner radii so they never exceed half of
// the width/height or sum constraints. Mirrors @nivo/theming
// constrainBorderRadius.
func ConstrainBorderRadius(r BorderRadius, width, height float64) BorderRadiusCorners {
	c := NormalizeBorderRadius(r)
	tl := math.Max(0, c.TopLeft)
	tr := math.Max(0, c.TopRight)
	br := math.Max(0, c.BottomRight)
	bl := math.Max(0, c.BottomLeft)

	if sum := tl + tr; sum > width {
		k := width / sum
		tl *= k
		tr *= k
	}
	if sum := bl + br; sum > width {
		k := width / sum
		bl *= k
		br *= k
	}
	if sum := tl + bl; sum > height {
		k := height / sum
		tl *= k
		bl *= k
	}
	if sum := tr + br; sum > height {
		k := height / sum
		tr *= k
		br *= k
	}
	return BorderRadiusCorners{TopLeft: tl, TopRight: tr, BottomRight: br, BottomLeft: bl}
}

// BorderRadiusToCss renders a CSS-compatible border-radius string
// (e.g. "4px 4px 0 0"). Mirrors @nivo/theming borderRadiusToCss.
func BorderRadiusToCss(c BorderRadiusCorners) string {
	return formatPx(c.TopLeft) + " " + formatPx(c.TopRight) + " " + formatPx(c.BottomRight) + " " + formatPx(c.BottomLeft)
}

func formatPx(v float64) string {
	if v == 0 {
		return "0"
	}
	return ftoa(v) + "px"
}

// ftoa formats a float without trailing zeros.
func ftoa(v float64) string {
	s := ""
	if v == float64(int(v)) {
		s = itoa(int(v))
	} else {
		s = ftoaTrim(v)
	}
	return s
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func ftoaTrim(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}
