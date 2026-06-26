// Package rects provides rect primitives shared by the chart packages: the
// Rect/NodeA11yProps/BorderRadius types, BuildRoundedRectPath (the primary
// rect path helper used by bar), the RoundedRect templ component (which emits
// SMIL <animate> enter animations when Animate is true), and the 9-anchor
// RectLabels label-position factories.
package rects

import (
	"fmt"
	"strings"
)

// Rect is a rectangle in svg units.
type Rect struct {
	X, Y, Width, Height float64
}

// NodeA11yProps carries the accessibility attributes emitted on chart nodes
// (bar, pie, …). Mirrors nivo's NodeA11yProps / a11yProps.
type NodeA11yProps struct {
	Tabindex        string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Role            string
	Focusable       bool
}

// BorderRadiusCorners allows per-corner radius. A zero value means square
// corners.
type BorderRadiusCorners struct {
	TopLeft     float64
	TopRight    float64
	BottomLeft  float64
	BottomRight float64
}

// BorderRadius is the union of a uniform radius (float) or per-corner radius.
// Uniform is used when non-zero; otherwise BorderRadiusCorners applies.
type BorderRadius struct {
	Uniform float64
	Corners BorderRadiusCorners
}

// NodeWithRect is implemented by any computed datum that has a rect (bar item,
// pie arc, …). Used by label-anchor factories.
type NodeWithRect interface {
	GetRect() Rect
}

// BorderRadiusFromFloat builds a uniform BorderRadius.
func BorderRadiusFromFloat(r float64) BorderRadius {
	return BorderRadius{Uniform: r}
}

// BorderRadiusFromCorners builds a per-corner BorderRadius.
func BorderRadiusFromCorners(c BorderRadiusCorners) BorderRadius {
	return BorderRadius{Corners: c}
}

// Resolved returns the effective per-corner radii, clamped so each radius
// cannot exceed half the shorter side of the rect (matches d3-shape's
// cornerRadius clamping).
func (b BorderRadius) Resolved(width, height float64) BorderRadiusCorners {
	if b.Uniform > 0 {
		return clampCorners(BorderRadiusCorners{
			TopLeft: b.Uniform, TopRight: b.Uniform,
			BottomLeft: b.Uniform, BottomRight: b.Uniform,
		}, width, height)
	}
	return clampCorners(b.Corners, width, height)
}

func clampCorners(c BorderRadiusCorners, w, h float64) BorderRadiusCorners {
	maxR := w
	if h < maxR {
		maxR = h
	}
	maxR /= 2
	if c.TopLeft > maxR {
		c.TopLeft = maxR
	}
	if c.TopRight > maxR {
		c.TopRight = maxR
	}
	if c.BottomLeft > maxR {
		c.BottomLeft = maxR
	}
	if c.BottomRight > maxR {
		c.BottomRight = maxR
	}
	return c
}

// fmtR formats a float for path output (3 dp, trimmed).
func fmtR(v float64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.3f", v), "0"), ".")
}

// BuildRoundedRectPath returns an SVG path-data string for a rounded rectangle
// at (x, y) with width w, height h, and per-corner radii tl/tr/br/bl. A zero
// radius produces a sharp corner. Mirrors @nivo/core's
// buildRoundedRectPathFromBorderRadius (which itself mirrors d3-shape's
// roundedRect path). All radii are clamped to min(r, w/2, h/2).
//
// The path is drawn clockwise from the top-left corner:
//
//	M (x+tl, y) H (x+w-tr) A tr,tr 0 0 1 (x+w, y+tr)
//	V (y+h-br) A br,br 0 0 1 (x+w-br, y+h) H (x+bl) A bl,bl 0 0 1 (x, y+h-bl)
//	V (y+tl)  A tl,tl 0 0 1 (x+tl, y) Z
func BuildRoundedRectPath(x, y, w, h, tl, tr, br, bl float64) string {
	// Clamp each radius to the smaller of w/2, h/2.
	maxR := w / 2
	if h2 := h / 2; h2 < maxR {
		maxR = h2
	}
	clamp := func(r float64) float64 {
		if r < 0 {
			return 0
		}
		if r > maxR {
			return maxR
		}
		return r
	}
	tl = clamp(tl)
	tr = clamp(tr)
	br = clamp(br)
	bl = clamp(bl)

	// Special-case: fully square corners.
	if tl == 0 && tr == 0 && br == 0 && bl == 0 {
		return fmt.Sprintf("M%s,%sH%sV%sH%sV%sZ", fmtR(x), fmtR(y), fmtR(x+w), fmtR(y+h), fmtR(x), fmtR(y))
	}

	var b strings.Builder
	x0, y0 := x, y
	x1, y1 := x+w, y+h
	fmt.Fprintf(&b, "M%s,%s", fmtR(x0+tl), fmtR(y0))
	fmt.Fprintf(&b, "H%s", fmtR(x1-tr))
	if tr > 0 {
		fmt.Fprintf(&b, "A%s,%s 0 0 1 %s,%s", fmtR(tr), fmtR(tr), fmtR(x1), fmtR(y0+tr))
	}
	fmt.Fprintf(&b, "V%s", fmtR(y1-br))
	if br > 0 {
		fmt.Fprintf(&b, "A%s,%s 0 0 1 %s,%s", fmtR(br), fmtR(br), fmtR(x1-br), fmtR(y1))
	}
	fmt.Fprintf(&b, "H%s", fmtR(x0+bl))
	if bl > 0 {
		fmt.Fprintf(&b, "A%s,%s 0 0 1 %s,%s", fmtR(bl), fmtR(bl), fmtR(x0), fmtR(y1-bl))
	}
	fmt.Fprintf(&b, "V%s", fmtR(y0+tl))
	if tl > 0 {
		fmt.Fprintf(&b, "A%s,%s 0 0 1 %s,%s", fmtR(tl), fmtR(tl), fmtR(x0+tl), fmtR(y0))
	}
	b.WriteString("Z")
	return b.String()
}
