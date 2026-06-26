// arc.go — port of d3-shape's arc generator (src/arc.js).
//
// Arc() returns a generator that produces an SVG path-data string for a
// circular or annular sector defined by {startAngle, endAngle, innerRadius,
// outerRadius, padAngle}. The port includes the cornerRadius and padAngle
// logic, which is the densest math in d3-shape (~250 lines of corner-tangent
// intersection geometry).
//
// Angles are in radians, measured from the +x axis, but d3-shape offsets by
// -π/2 so that angle 0 is at the top (12 o'clock) and angles increase
// clockwise. The arc is drawn centered at the origin (0,0); callers translate
// to the desired center.
package d3shape

import "math"

// ArcDatum is the input to the arc generator.
type ArcDatum struct {
	StartAngle  float64
	EndAngle    float64
	InnerRadius float64
	OuterRadius float64
	PadAngle    float64
}

// ArcRadiusAccessor extracts a radius from an ArcDatum.
type (
	ArcRadiusAccessor func(d ArcDatum) float64
	ArcAngleAccessor  func(d ArcDatum) float64
)

// Arc is an arc generator producing SVG path-data strings.
type Arc struct {
	innerRadius  ArcRadiusAccessor
	outerRadius  ArcRadiusAccessor
	cornerRadius func(d ArcDatum) float64
	padRadius    func(d ArcDatum) float64
	startAngle   ArcAngleAccessor
	endAngle     ArcAngleAccessor
	padAngle     func(d ArcDatum) float64
	digits       int
}

// NewArc constructs an Arc with default accessors reading from ArcDatum
// fields, cornerRadius=0, padRadius=nil.
func NewArc() *Arc {
	return &Arc{
		innerRadius:  func(d ArcDatum) float64 { return d.InnerRadius },
		outerRadius:  func(d ArcDatum) float64 { return d.OuterRadius },
		cornerRadius: func(d ArcDatum) float64 { return 0 },
		padRadius:    nil,
		startAngle:   func(d ArcDatum) float64 { return d.StartAngle },
		endAngle:     func(d ArcDatum) float64 { return d.EndAngle },
		padAngle:     func(d ArcDatum) float64 { return d.PadAngle },
		digits:       3,
	}
}

// InnerRadius sets the inner radius accessor (or a constant if a number is
// passed via a closure).
func (a *Arc) InnerRadius(fn ArcRadiusAccessor) *Arc { a.innerRadius = fn; return a }

// OuterRadius sets the outer radius accessor.
func (a *Arc) OuterRadius(fn ArcRadiusAccessor) *Arc { a.outerRadius = fn; return a }

// CornerRadius sets the corner radius.
func (a *Arc) CornerRadius(r float64) *Arc {
	a.cornerRadius = func(d ArcDatum) float64 { return r }
	return a
}

// PadRadius sets the pad radius. 0 = use sqrt(r0² + r1²).
func (a *Arc) PadRadius(r float64) *Arc {
	if r == 0 {
		a.padRadius = nil
	} else {
		a.padRadius = func(d ArcDatum) float64 { return r }
	}
	return a
}

// StartAngle sets the start angle accessor.
func (a *Arc) StartAngle(fn ArcAngleAccessor) *Arc { a.startAngle = fn; return a }

// EndAngle sets the end angle accessor.
func (a *Arc) EndAngle(fn ArcAngleAccessor) *Arc { a.endAngle = fn; return a }

// PadAngle sets the pad angle.
func (a *Arc) PadAngle(pa float64) *Arc {
	a.padAngle = func(d ArcDatum) float64 { return pa }
	return a
}

// Digits sets the path rounding precision (default 3).
func (a *Arc) Digits(d int) *Arc { a.digits = d; return a }

// Call generates the SVG path-data string for the given arc datum.
func (a *Arc) Call(d ArcDatum) string {
	path := NewPathDigits(a.digits)
	a.draw(path, d)
	return path.String()
}

// Centroid computes the centroid [x, y] of the arc (d3-shape arc.centroid).
func (a *Arc) Centroid(d ArcDatum) [2]float64 {
	r := (a.innerRadius(d) + a.outerRadius(d)) / 2
	ang := (a.startAngle(d)+a.endAngle(d))/2 - pi/2
	return [2]float64{math.Cos(ang) * r, math.Sin(ang) * r}
}

// draw emits the arc path to the given Path. This is a faithful
// transliteration of d3-shape's arc.js.
func (a *Arc) draw(path *Path, d ArcDatum) {
	r0 := a.innerRadius(d)
	r1 := a.outerRadius(d)
	a0 := a.startAngle(d) - halfPi
	a1 := a.endAngle(d) - halfPi
	da := math.Abs(a1 - a0)
	cw := a1 > a0

	// Ensure outer >= inner
	if r1 < r0 {
		r0, r1 = r1, r0
	}

	// Is it a point?
	if !(r1 > epsilon) {
		path.moveTo(0, 0)
	} else if da > tau-epsilon {
		// Full circle or annulus
		path.moveTo(r1*math.Cos(a0), r1*math.Sin(a0))
		path.arc(0, 0, r1, a0, a1, !cw)
		if r0 > epsilon {
			path.moveTo(r0*math.Cos(a1), r0*math.Sin(a1))
			path.arc(0, 0, r0, a1, a0, cw)
		}
	} else {
		// Circular or annular sector
		a01 := a0
		a11 := a1
		a00 := a0
		a10 := a1
		da0 := da
		da1 := da
		ap := a.padAngle(d) / 2
		rp := 0.0
		if ap > epsilon {
			if a.padRadius != nil {
				rp = a.padRadius(d)
			} else {
				rp = math.Sqrt(r0*r0 + r1*r1)
			}
		}
		rc := math.Min(math.Abs(r1-r0)/2, a.cornerRadius(d))
		rc0 := rc
		rc1 := rc

		// Apply padding
		if rp > epsilon {
			p0 := asin(rp / r0 * math.Sin(ap))
			p1 := asin(rp / r1 * math.Sin(ap))
			da0 -= p0 * 2
			if da0 > epsilon {
				p0 *= sign2(cw)
				a00 += p0
				a10 -= p0
			} else {
				da0 = 0
				a00 = (a0 + a1) / 2
				a10 = a00
			}
			da1 -= p1 * 2
			if da1 > epsilon {
				p1 *= sign2(cw)
				a01 += p1
				a11 -= p1
			} else {
				da1 = 0
				a01 = (a0 + a1) / 2
				a11 = a01
			}
		}

		x01 := r1 * math.Cos(a01)
		y01 := r1 * math.Sin(a01)
		x10 := r0 * math.Cos(a10)
		y10 := r0 * math.Sin(a10)

		// Apply rounded corners?
		if rc > epsilon {
			x11 := r1 * math.Cos(a11)
			y11 := r1 * math.Sin(a11)
			x00 := r0 * math.Cos(a00)
			y00 := r0 * math.Sin(a00)
			var oc []float64

			// Restrict the corner radius according to the sector angle
			if da < pi {
				oc = intersect(x01, y01, x00, y00, x11, y11, x10, y10)
				if oc != nil {
					ax := x01 - oc[0]
					ay := y01 - oc[1]
					bx := x11 - oc[0]
					by := y11 - oc[1]
					kc := 1 / math.Sin(acos((ax*bx+ay*by)/(math.Sqrt(ax*ax+ay*ay)*math.Sqrt(bx*bx+by*by)))/2)
					lc := math.Sqrt(oc[0]*oc[0] + oc[1]*oc[1])
					rc0 = math.Min(rc, (r0-lc)/(kc-1))
					rc1 = math.Min(rc, (r1-lc)/(kc+1))
				} else {
					rc0 = 0
					rc1 = 0
				}
			}

			// Is the sector collapsed to a line?
			if !(da1 > epsilon) {
				path.moveTo(x01, y01)
			} else if rc1 > epsilon {
				// Outer ring has rounded corners
				t0 := cornerTangents(x00, y00, x01, y01, r1, rc1, cw)
				t1 := cornerTangents(x11, y11, x10, y10, r1, rc1, cw)

				path.moveTo(t0.cx+t0.x01, t0.cy+t0.y01)

				// Have the corners merged?
				if rc1 < rc {
					path.arc(t0.cx, t0.cy, rc1, math.Atan2(t0.y01, t0.x01), math.Atan2(t1.y01, t1.x01), !cw)
				} else {
					// Draw the two corners and the ring
					path.arc(t0.cx, t0.cy, rc1, math.Atan2(t0.y01, t0.x01), math.Atan2(t0.y11, t0.x11), !cw)
					path.arc(0, 0, r1, math.Atan2(t0.cy+t0.y11, t0.cx+t0.x11), math.Atan2(t1.cy+t1.y11, t1.cx+t1.x11), !cw)
					path.arc(t1.cx, t1.cy, rc1, math.Atan2(t1.y11, t1.x11), math.Atan2(t1.y01, t1.x01), !cw)
				}
			} else {
				// Outer ring is just a circular arc
				path.moveTo(x01, y01)
				path.arc(0, 0, r1, a01, a11, !cw)
			}

			// Is there no inner ring?
			if !(r0 > epsilon) || !(da0 > epsilon) {
				path.lineTo(x10, y10)
			} else if rc0 > epsilon {
				// Inner ring has rounded corners
				t0 := cornerTangents(x10, y10, x11, y11, r0, -rc0, cw)
				t1 := cornerTangents(x01, y01, x00, y00, r0, -rc0, cw)

				path.lineTo(t0.cx+t0.x01, t0.cy+t0.y01)

				// Have the corners merged?
				if rc0 < rc {
					path.arc(t0.cx, t0.cy, rc0, math.Atan2(t0.y01, t0.x01), math.Atan2(t1.y01, t1.x01), !cw)
				} else {
					// Draw the two corners and the ring
					path.arc(t0.cx, t0.cy, rc0, math.Atan2(t0.y01, t0.x01), math.Atan2(t0.y11, t0.x11), !cw)
					path.arc(0, 0, r0, math.Atan2(t0.cy+t0.y11, t0.cx+t0.x11), math.Atan2(t1.cy+t1.y11, t1.cx+t1.x11), cw)
					path.arc(t1.cx, t1.cy, rc0, math.Atan2(t1.y11, t1.x11), math.Atan2(t1.y01, t1.x01), !cw)
				}
			} else {
				// Inner ring is just a circular arc
				path.arc(0, 0, r0, a10, a00, cw)
			}
		} else {
			// No rounded corners
			if !(da1 > epsilon) {
				path.moveTo(x01, y01)
			} else {
				path.moveTo(x01, y01)
				path.arc(0, 0, r1, a01, a11, !cw)
			}

			if !(r0 > epsilon) || !(da0 > epsilon) {
				path.lineTo(x10, y10)
			} else {
				path.arc(0, 0, r0, a10, a00, cw)
			}
		}
	}

	path.closePath()
}

// sign2 returns +1 for true, -1 for false (d3 uses `cw ? 1 : -1`).
func sign2(cw bool) float64 {
	if cw {
		return 1
	}
	return -1
}

// intersect finds the intersection of two lines (d3-shape arc.js intersect).
func intersect(x0, y0, x1, y1, x2, y2, x3, y3 float64) []float64 {
	x10 := x1 - x0
	y10 := y1 - y0
	x32 := x3 - x2
	y32 := y3 - y2
	t := y32*x10 - x32*y10
	if t*t < epsilon {
		return nil
	}
	t = (x32*(y0-y2) - y32*(x0-x2)) / t
	return []float64{x0 + t*x10, y0 + t*y10}
}

// cornerTangent holds the results of cornerTangents.
type cornerTangent struct {
	cx, cy   float64
	x01, y01 float64
	x11, y11 float64
}

// cornerTangents computes the perpendicular offset line of length rc.
// Faithful transliteration of d3-shape's cornerTangents.
func cornerTangents(x0, y0, x1, y1, r1, rc float64, cw bool) cornerTangent {
	x01 := x0 - x1
	y01 := y0 - y1
	var lo float64
	if cw {
		lo = rc
	} else {
		lo = -rc
	}
	lo = lo / math.Sqrt(x01*x01+y01*y01)
	ox := lo * y01
	oy := -lo * x01
	x11 := x0 + ox
	y11 := y0 + oy
	x10 := x1 + ox
	y10 := y1 + oy
	x00 := (x11 + x10) / 2
	y00 := (y11 + y10) / 2
	dx := x10 - x11
	dy := y10 - y11
	d2 := dx*dx + dy*dy
	r := r1 - rc
	D := x11*y10 - x10*y11
	var d float64
	if dy < 0 {
		d = -1
	} else {
		d = 1
	}
	d = d * math.Sqrt(math.Max(0, r*r*d2-D*D))
	cx0 := (D*dy - dx*d) / d2
	cy0 := (-D*dx - dy*d) / d2
	cx1 := (D*dy + dx*d) / d2
	cy1 := (-D*dx + dy*d) / d2
	dx0 := cx0 - x00
	dy0 := cy0 - y00
	dx1 := cx1 - x00
	dy1 := cy1 - y00

	// Pick the closer of the two intersection points.
	if dx0*dx0+dy0*dy0 > dx1*dx1+dy1*dy1 {
		cx0 = cx1
		cy0 = cy1
	}

	return cornerTangent{
		cx:  cx0,
		cy:  cy0,
		x01: -ox,
		y01: -oy,
		x11: cx0 * (r1/r - 1),
		y11: cy0 * (r1/r - 1),
	}
}
