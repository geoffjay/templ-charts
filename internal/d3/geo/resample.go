package geo

import "math"

const maxDepth = 16

var cosMinDistance = math.Cos(30 * radians)

// resample returns a stream factory that adaptively subdivides projected line
// segments so straight lines in (lon,lat) become smooth curves in the plane.
// delta2 is the squared pixel tolerance; delta2==0 disables resampling. Ported
// from d3-geo src/resample.js.
func resample(project transform, delta2 float64) func(Sink) Sink {
	if delta2 > 0 {
		return func(sink Sink) Sink {
			r := &resampleStream{sink: sink, project: project, delta2: delta2}
			r.pointFn = r.point
			r.lineStartFn = r.lineStart
			r.lineEndFn = r.lineEnd
			return r
		}
	}
	return func(sink Sink) Sink {
		return &passThrough{stream: sink, pointFn: func(x, y float64) {
			px, py := project.forward(x, y)
			sink.Point(px, py)
		}}
	}
}

type resampleStream struct {
	sink    Sink
	project transform
	delta2  float64

	pointFn     func(x, y float64)
	lineStartFn func()
	lineEndFn   func()

	// first point of the current ring
	lambda00, x00, y00, a00, b00, c00 float64
	// previous point
	lambda0, x0, y0, a0, b0, c0 float64
}

func (r *resampleStream) Point(x, y float64) { r.pointFn(x, y) }
func (r *resampleStream) LineStart()         { r.lineStartFn() }
func (r *resampleStream) LineEnd()           { r.lineEndFn() }
func (r *resampleStream) PolygonStart()      { r.sink.PolygonStart(); r.lineStartFn = r.ringStart }
func (r *resampleStream) PolygonEnd()        { r.sink.PolygonEnd(); r.lineStartFn = r.lineStart }
func (r *resampleStream) Sphere()            { r.sink.Sphere() }

func (r *resampleStream) resampleLineTo(x0, y0, lambda0, a0, b0, c0, x1, y1, lambda1, a1, b1, c1 float64, depth int) {
	dx := x1 - x0
	dy := y1 - y0
	d2 := dx*dx + dy*dy
	if d2 > 4*r.delta2 && depth > 0 {
		depth--
		a := a0 + a1
		b := b0 + b1
		c := c0 + c1
		m := math.Sqrt(a*a + b*b + c*c)
		c /= m
		phi2 := asin(c)
		var lambda2 float64
		if math.Abs(math.Abs(c)-1) < epsilon || math.Abs(lambda0-lambda1) < epsilon {
			lambda2 = (lambda0 + lambda1) / 2
		} else {
			lambda2 = math.Atan2(b, a)
		}
		x2, y2 := r.project.forward(lambda2, phi2)
		dx2 := x2 - x0
		dy2 := y2 - y0
		dz := dy*dx2 - dx*dy2
		if dz*dz/d2 > r.delta2 || math.Abs((dx*dx2+dy*dy2)/d2-0.5) > 0.3 || a0*a1+b0*b1+c0*c1 < cosMinDistance {
			a /= m
			b /= m
			r.resampleLineTo(x0, y0, lambda0, a0, b0, c0, x2, y2, lambda2, a, b, c, depth)
			r.sink.Point(x2, y2)
			r.resampleLineTo(x2, y2, lambda2, a, b, c, x1, y1, lambda1, a1, b1, c1, depth)
		}
	}
}

func (r *resampleStream) point(x, y float64) {
	px, py := r.project.forward(x, y)
	r.sink.Point(px, py)
}

func (r *resampleStream) lineStart() {
	r.x0 = math.NaN()
	r.pointFn = r.linePoint
	r.sink.LineStart()
}

func (r *resampleStream) linePoint(lambda, phi float64) {
	c := cartesian([2]float64{lambda, phi})
	px, py := r.project.forward(lambda, phi)
	r.resampleLineTo(r.x0, r.y0, r.lambda0, r.a0, r.b0, r.c0, px, py, lambda, c[0], c[1], c[2], maxDepth)
	r.x0, r.y0, r.lambda0, r.a0, r.b0, r.c0 = px, py, lambda, c[0], c[1], c[2]
	r.sink.Point(px, py)
}

func (r *resampleStream) lineEnd() {
	r.pointFn = r.point
	r.sink.LineEnd()
}

func (r *resampleStream) ringStart() {
	r.lineStart()
	r.pointFn = r.ringPoint
	r.lineEndFn = r.ringEnd
}

func (r *resampleStream) ringPoint(lambda, phi float64) {
	r.lambda00 = lambda
	r.linePoint(lambda, phi)
	r.x00, r.y00, r.a00, r.b00, r.c00 = r.x0, r.y0, r.a0, r.b0, r.c0
	r.pointFn = r.linePoint
}

func (r *resampleStream) ringEnd() {
	r.resampleLineTo(r.x0, r.y0, r.lambda0, r.a0, r.b0, r.c0, r.x00, r.y00, r.lambda00, r.a00, r.b00, r.c00, maxDepth)
	r.lineEndFn = r.lineEnd
	r.lineEnd()
}
