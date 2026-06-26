// path.go — port of d3-path's Path serializer.
//
// Accumulates SVG path-data commands and serializes them to a string via
// String(). Coordinates are rounded to `digits` decimal places (default 3,
// matching d3-shape's withPath default). The method signatures mirror d3-path's
// Path class (moveTo, lineTo, arc, bezierCurveTo, etc.) so d3-shape curve/generator
// ports can call them with the same semantics.
package d3shape

import (
	"math"
	"strconv"
)

// Path is a port of d3-path's Path serializer.
type Path struct {
	sx, sy   float64 // start of current subpath
	ex, ey   float64 // end of current subpath
	buf      string  // accumulated path-data string
	digits   int     // decimal places; -1 = no rounding
	hasStart bool    // tracks whether a subpath has been started
}

// NewPath returns a Path that rounds coordinates to 3 decimal places
// (d3-shape's default for generators).
func NewPath() *Path { return &Path{digits: 3} }

// NewPathDigits returns a Path with custom rounding precision.
// digits < 0 disables rounding (full float precision).
func NewPathDigits(digits int) *Path {
	if digits < 0 {
		digits = -1
	}
	return &Path{digits: digits}
}

// fmtNum formats a coordinate per the Path's digit setting.
func (p *Path) fmtNum(v float64) string {
	if p.digits < 0 {
		return floatToShortest(v)
	}
	k := math.Pow10(p.digits)
	return floatToShortest(math.Round(v*k) / k)
}

// floatToShortest formats v with the shortest round-trip representation,
// matching d3's Number.prototype.toString coercion. Uses strconv 'f' -1 for
// fixed-point (chart coordinates are well within fixed-point range); falls
// back to 'g' for very large/small magnitudes.
func floatToShortest(v float64) string {
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 1) {
		return "Infinity"
	}
	if math.IsInf(v, -1) {
		return "-Infinity"
	}
	if v == 0 {
		return "0"
	}
	// 'g' with -1 precision matches JS Number.toString for most values and
	// handles exponential notation for extreme magnitudes.
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// moveTo emits "Mx,y" and starts a new subpath.
func (p *Path) moveTo(x, y float64) {
	p.sx = x
	p.sy = y
	p.ex = x
	p.ey = y
	p.hasStart = true
	p.buf += "M" + p.fmtNum(x) + "," + p.fmtNum(y)
}

// closePath emits "Z" and returns the current point to the subpath start.
func (p *Path) closePath() {
	if p.hasStart {
		p.ex = p.sx
		p.ey = p.sy
		p.buf += "Z"
	}
}

// lineTo emits "Lx,y".
func (p *Path) lineTo(x, y float64) {
	p.ex = x
	p.ey = y
	p.hasStart = true
	p.buf += "L" + p.fmtNum(x) + "," + p.fmtNum(y)
}

// quadraticCurveTo emits "Qx1,y1,x,y".
func (p *Path) quadraticCurveTo(x1, y1, x, y float64) {
	p.ex = x
	p.ey = y
	p.hasStart = true
	p.buf += "Q" + p.fmtNum(x1) + "," + p.fmtNum(y1) + "," + p.fmtNum(x) + "," + p.fmtNum(y)
}

// bezierCurveTo emits "Cx1,y1,x2,y2,x,y".
func (p *Path) bezierCurveTo(x1, y1, x2, y2, x, y float64) {
	p.ex = x
	p.ey = y
	p.hasStart = true
	p.buf += "C" + p.fmtNum(x1) + "," + p.fmtNum(y1) + "," + p.fmtNum(x2) + "," + p.fmtNum(y2) + "," + p.fmtNum(x) + "," + p.fmtNum(y)
}

// arcTo is a port of d3-path's Path.arcTo. Included for completeness; not
// used by line/area/arc generators (which use arc, not arcTo).
func (p *Path) arcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		panic("d3shape: negative radius")
	}
	x0, y0 := p.ex, p.ey
	x21 := x2 - x1
	y21 := y2 - y1
	x01 := x0 - x1
	y01 := y0 - y1
	l01_2 := x01*x01 + y01*y01

	if !p.hasStart {
		p.moveTo(x1, y1)
	} else if !(l01_2 > epsilon) {
		// coincident with (x0,y0): do nothing
	} else if !(math.Abs(y01*x21-y21*x01) > epsilon) || r == 0 {
		p.lineTo(x1, y1)
	} else {
		x20 := x2 - x0
		y20 := y2 - y0
		l21_2 := x21*x21 + y21*y21
		l20_2 := x20*x20 + y20*y20
		l21 := math.Sqrt(l21_2)
		l01 := math.Sqrt(l01_2)
		l := r * math.Tan((pi-math.Acos((l21_2+l01_2-l20_2)/(2*l21*l01)))/2)
		t01 := l / l01
		t21 := l / l21
		if math.Abs(t01-1) > epsilon {
			p.lineTo(x1+t01*x01, y1+t01*y01)
		}
		cw := 0
		if y01*x20 > x01*y20 {
			cw = 1
		}
		p.ex = x1 + t21*x21
		p.ey = y1 + t21*y21
		p.hasStart = true
		p.buf += "A" + p.fmtNum(r) + "," + p.fmtNum(r) + ",0,0," + intToStr(cw) + "," + p.fmtNum(p.ex) + "," + p.fmtNum(p.ey)
	}
}

// arc is a port of d3-path's Path.arc. Draws a circular arc centered at
// (x,y) with radius r from angle a0 to a1. ccw reverses direction.
func (p *Path) arc(x, y, r, a0, a1 float64, ccw bool) {
	if r < 0 {
		panic("d3shape: negative radius")
	}
	dx := r * math.Cos(a0)
	dy := r * math.Sin(a0)
	x0 := x + dx
	y0 := y + dy
	cw := 1
	if ccw {
		cw = 0
	}
	var da float64
	if ccw {
		da = a0 - a1
	} else {
		da = a1 - a0
	}

	if !p.hasStart {
		p.buf += "M" + p.fmtNum(x0) + "," + p.fmtNum(y0)
		p.hasStart = true
	} else if math.Abs(p.ex-x0) > epsilon || math.Abs(p.ey-y0) > epsilon {
		p.lineTo(x0, y0)
	}
	if r == 0 {
		return
	}
	if da < 0 {
		da = math.Mod(da, tau) + tau
	}
	if da > tauEpsilon {
		// complete circle: two arcs
		p.buf += "A" + p.fmtNum(r) + "," + p.fmtNum(r) + ",0,1," + intToStr(cw) + "," + p.fmtNum(x-dx) + "," + p.fmtNum(y-dy)
		p.ex = x0
		p.ey = y0
		p.buf += "A" + p.fmtNum(r) + "," + p.fmtNum(r) + ",0,1," + intToStr(cw) + "," + p.fmtNum(x0) + "," + p.fmtNum(y0)
	} else if da > epsilon {
		largeArc := 0
		if da >= pi {
			largeArc = 1
		}
		p.ex = x + r*math.Cos(a1)
		p.ey = y + r*math.Sin(a1)
		p.hasStart = true
		p.buf += "A" + p.fmtNum(r) + "," + p.fmtNum(r) + ",0," + intToStr(largeArc) + "," + intToStr(cw) + "," + p.fmtNum(p.ex) + "," + p.fmtNum(p.ey)
	}
}

// rect emits a closed rectangle (d3-path's rect).
func (p *Path) rect(x, y, w, h float64) {
	p.sx = x
	p.sy = y
	p.ex = x
	p.ey = y
	p.hasStart = true
	p.buf += "M" + p.fmtNum(x) + "," + p.fmtNum(y) + "h" + p.fmtNum(w) + "v" + p.fmtNum(h) + "h" + p.fmtNum(-w) + "Z"
}

// String returns the accumulated SVG path-data string.
func (p *Path) String() string { return p.buf }

// Reset clears the path for reuse.
func (p *Path) Reset() {
	p.buf = ""
	p.hasStart = false
	p.sx, p.sy, p.ex, p.ey = 0, 0, 0, 0
}

// intToStr formats a non-negative int.
func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	return strconv.Itoa(i)
}
