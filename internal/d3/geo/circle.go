package geo

import "math"

// clipCircle is the small-circle preclip the azimuthal family uses to hide the
// far hemisphere: geometry outside a circle of the given radius (centered on the
// projection's center after rotation) is clipped away, and rings that straddle
// the boundary are cut and rejoined along the circle. Ported from d3-geo
// src/clip/circle.js. Without it, orthographic/gnomonic/stereographic/azimuthal*
// render the whole sphere and can draw far-side geometry over the near side.
func clipCircle(radius float64) func(Sink) Sink {
	cr := math.Cos(radius)
	delta := 2 * radians
	smallRadius := cr > 0
	notHemisphere := math.Abs(cr) > epsilon // TODO optimise for the common hemisphere case

	visible := func(lambda, phi float64) bool {
		return math.Cos(lambda)*math.Cos(phi) > cr
	}

	interpolate := func(from, to *[3]float64, direction float64, stream Sink) {
		var t0, t1 *[2]float64
		if from != nil {
			t0 = &[2]float64{from[0], from[1]}
		}
		if to != nil {
			t1 = &[2]float64{to[0], to[1]}
		}
		circleStream(stream, radius, delta, direction, t0, t1)
	}

	clipLine := func(sink Sink) clipLiner {
		return &circleLine{
			stream:        sink,
			radius:        radius,
			cr:            cr,
			smallRadius:   smallRadius,
			notHemisphere: notHemisphere,
		}
	}

	var start [2]float64
	if smallRadius {
		start = [2]float64{0, -radius}
	} else {
		start = [2]float64{-pi, radius - pi}
	}
	return newClip(visible, clipLine, interpolate, start)
}

// circleLine cuts a line into visible segments against the clip circle. Clean()
// returns: bit0 = no intersections; bit1 = first & last points both visible (so
// the first/last segments should be rejoined). Ported from the clipLine closure
// in d3-geo src/clip/circle.js.
type circleLine struct {
	stream        Sink
	radius, cr    float64
	smallRadius   bool
	notHemisphere bool

	point0    [2]float64 // previous point
	hasPoint0 bool
	c0        int  // code for previous point
	v0        bool // visibility of previous point
	v00       bool // visibility of first point
	clean     int  // no intersections
}

func (l *circleLine) LineStart() {
	l.v00 = false
	l.v0 = false
	l.clean = 1
}

func (l *circleLine) Point(lambda, phi float64) {
	point1 := [2]float64{lambda, phi}
	v := l.cr < math.Cos(lambda)*math.Cos(phi)

	var c int
	if l.smallRadius {
		if !v {
			c = circleCode(lambda, phi, l.radius, l.smallRadius)
		}
	} else if v {
		adj := lambda - pi
		if lambda < 0 {
			adj = lambda + pi
		}
		c = circleCode(adj, phi, l.radius, l.smallRadius)
	}

	if !l.hasPoint0 {
		l.v00 = v
		l.v0 = v
		if v {
			l.stream.LineStart()
		}
	}

	if v != l.v0 {
		// The visibility changed between point0 and point1: cut at the circle.
		l.clean = 0
		if v {
			// Outside going in.
			l.stream.LineStart()
			p2, _ := circleIntersectOne(point1, l.point0, l.cr)
			l.stream.Point(p2[0], p2[1])
			l.point0 = p2
			l.hasPoint0 = true
		} else {
			// Inside going out.
			p2, _ := circleIntersectOne(l.point0, point1, l.cr)
			emitPoint(l.stream, p2[0], p2[1], 2)
			l.stream.LineEnd()
			l.point0 = p2
			l.hasPoint0 = true
		}
	} else if l.notHemisphere && l.hasPoint0 && (l.smallRadius != v) {
		// A segment between two points on the same side may still cross the
		// circle (both endpoints outside a small circle, or both inside a large
		// one). Only when their bounding-box codes don't overlap can it cross.
		if l.c0&c == 0 {
			if pair, ok := circleIntersectTwo(point1, l.point0, l.cr); ok {
				l.clean = 0
				if l.smallRadius {
					l.stream.LineStart()
					l.stream.Point(pair[0][0], pair[0][1])
					l.stream.Point(pair[1][0], pair[1][1])
					l.stream.LineEnd()
				} else {
					l.stream.Point(pair[1][0], pair[1][1])
					l.stream.LineEnd()
					l.stream.LineStart()
					emitPoint(l.stream, pair[0][0], pair[0][1], 3)
				}
			}
		}
	}

	if v && (!l.hasPoint0 || !pointEqual2(l.point0, point1)) {
		l.stream.Point(point1[0], point1[1])
	}
	l.point0 = point1
	l.hasPoint0 = true
	l.v0 = v
	l.c0 = c
}

func (l *circleLine) LineEnd() {
	if l.v0 {
		l.stream.LineEnd()
	}
	l.hasPoint0 = false
}

func (l *circleLine) Clean() int {
	c := l.clean
	if l.v00 && l.v0 {
		c |= 2
	}
	return c
}

// pointEqual2 reports whether two [lambda,phi] points coincide within epsilon.
func pointEqual2(a, b [2]float64) bool {
	return math.Abs(a[0]-b[0]) < epsilon && math.Abs(a[1]-b[1]) < epsilon
}

// circleCode returns a 4-bit box code for a point relative to the small circle's
// bounding box (left/right/below/above). Ported from d3-geo circle.js code().
func circleCode(lambda, phi, radius float64, smallRadius bool) int {
	r := radius
	if !smallRadius {
		r = pi - radius
	}
	code := 0
	if lambda < -r {
		code |= 1 // left
	} else if lambda > r {
		code |= 2 // right
	}
	if phi < -r {
		code |= 4 // below
	} else if phi > r {
		code |= 8 // above
	}
	return code
}

// circleIntersectOne intersects the great circle a→b with the clip circle,
// returning the single intersection point. Mirrors intersect(a,b) in d3-geo
// circle.js (the two==false case). ok is false only for the degenerate
// polar-point case, where d3 falls back to a itself.
func circleIntersectOne(a, b [2]float64, cr float64) (pt [2]float64, ok bool) {
	one, _, hasOne, _ := circleIntersect(a, b, cr, false)
	if hasOne {
		return one, true
	}
	return a, false // degenerate (determinant 0): d3 returns a
}

// circleIntersectTwo intersects the great circle a→b with the clip circle,
// returning both intersection points when the near one lies between a and b.
// Mirrors intersect(a,b,true) in d3-geo circle.js.
func circleIntersectTwo(a, b [2]float64, cr float64) (pair [2][2]float64, ok bool) {
	_, two, _, hasTwo := circleIntersect(a, b, cr, true)
	if hasTwo {
		return two, true
	}
	return pair, false
}

// circleIntersect is the shared core of intersect() in d3-geo circle.js. It
// finds where the great circle through a and b crosses the clip circle. For
// two==false it returns the near intersection (hasOne). For two==true it returns
// both intersections when the near one is between a and b (hasTwo).
func circleIntersect(a, b [2]float64, cr float64, two bool) (one [2]float64, pair [2][2]float64, hasOne bool, hasTwo bool) {
	pa := cartesian(a)
	pb := cartesian(b)

	// Two planes n1.p = d1 and n2.p = d2; their line of intersection passes
	// through the origin. n1 is the clip-circle plane normal [1,0,0].
	n1 := [3]float64{1, 0, 0}
	n2 := cartesianCross(pa, pb)
	n2n2 := cartesianDot(n2, n2)
	n1n2 := n2[0] // cartesianDot(n1, n2)
	determinant := n2n2 - n1n2*n1n2

	if determinant == 0 {
		// Two polar points: no proper intersection.
		return one, pair, false, false
	}

	c1 := cr * n2n2 / determinant
	c2 := -cr * n1n2 / determinant
	n1xn2 := cartesianCross(n1, n2)
	A := cartesianScale(n1, c1)
	B := cartesianScale(n2, c2)
	cartesianAddInPlace(&A, B)

	// Solve |p(t)|² = 1.
	u := n1xn2
	w := cartesianDot(A, u)
	uu := cartesianDot(u, u)
	t2 := w*w - uu*(cartesianDot(A, A)-1)
	if t2 < 0 {
		return one, pair, false, false
	}

	t := math.Sqrt(t2)
	q := cartesianScale(u, (-w-t)/uu)
	cartesianAddInPlace(&q, A)
	qs := spherical(q)

	if !two {
		return qs, pair, true, false
	}

	// Two intersection points: check the near one lies on the arc a→b.
	lambda0 := a[0]
	lambda1 := b[0]
	phi0 := a[1]
	phi1 := b[1]
	if lambda1 < lambda0 {
		lambda0, lambda1 = lambda1, lambda0
	}
	delta := lambda1 - lambda0
	polar := math.Abs(delta-pi) < epsilon
	meridian := polar || delta < epsilon
	if !polar && phi1 < phi0 {
		phi0, phi1 = phi1, phi0
	}

	var between bool
	if meridian {
		if polar {
			between = boolXor(phi0+phi1 > 0, qs[1] < ternaryPhi(math.Abs(qs[0]-lambda0) < epsilon, phi0, phi1))
		} else {
			between = phi0 <= qs[1] && qs[1] <= phi1
		}
	} else {
		between = boolXor(delta > pi, lambda0 <= qs[0] && qs[0] <= lambda1)
	}

	if between {
		q1 := cartesianScale(u, (-w+t)/uu)
		cartesianAddInPlace(&q1, A)
		pair[0] = qs
		pair[1] = spherical(q1)
		return one, pair, false, true
	}
	return one, pair, false, false
}

func ternaryPhi(cond bool, a, b float64) float64 {
	if cond {
		return a
	}
	return b
}

// circleStream emits boundary points along the clip circle from t0 to t1 in the
// given direction, or the whole circle when t0p is nil. Ported from d3-geo
// src/circle.js circleStream.
func circleStream(stream Sink, radius, delta, direction float64, t0p, t1p *[2]float64) {
	if delta == 0 {
		return
	}
	cosRadius := math.Cos(radius)
	sinRadius := math.Sin(radius)
	step := direction * delta

	var t0, t1 float64
	if t0p == nil {
		t0 = radius + direction*tau
		t1 = radius - step/2
	} else {
		t0 = circleRadiusAngle(cosRadius, *t0p)
		t1 = circleRadiusAngle(cosRadius, *t1p)
		if (direction > 0 && t0 < t1) || (direction < 0 && t0 > t1) {
			t0 += direction * tau
		}
	}

	for t := t0; (direction > 0 && t > t1) || (direction < 0 && t < t1); t -= step {
		point := spherical([3]float64{cosRadius, -sinRadius * math.Cos(t), -sinRadius * math.Sin(t)})
		stream.Point(point[0], point[1])
	}
}

// circleRadiusAngle returns the signed angle of a spherical point around the
// clip circle, used to order the interpolation endpoints. Ported from d3-geo
// src/circle.js circleRadius.
func circleRadiusAngle(cosRadius float64, point [2]float64) float64 {
	p := cartesian(point)
	p[0] -= cosRadius
	cartesianNormalizeInPlace(&p)
	radius := acos(-p[1])
	r := radius
	if -p[2] < 0 {
		r = -radius
	}
	return math.Mod(r+tau-epsilon, tau)
}
