package geo

// clipLiner is a per-line clipper (buffer or direct) exposing a Clean() marker
// used by the ring-stitching logic. Mirrors the object returned by a clipLine.
type clipLiner interface {
	Point(lambda, phi float64)
	LineStart()
	LineEnd()
	Clean() int
}

type clipLineFunc func(sink Sink) clipLiner

// newClip builds a preclip stream factory from a point-visibility test, a
// per-line clipper, a boundary interpolator, and a clip start point. Ported
// from d3-geo src/clip/index.js.
func newClip(pointVisible func(lambda, phi float64) bool, clipLine clipLineFunc, interpolate interpolateFunc, start [2]float64) func(Sink) Sink {
	return func(sink Sink) Sink {
		c := &clipStream{
			sink:         sink,
			pointVisible: pointVisible,
			interpolate:  interpolate,
			start:        start,
			ringBuf:      &clipBuffer{},
		}
		c.line = clipLine(sink)
		c.ringSink = clipLine(c.ringBuf)
		c.pointFn = c.point
		c.lineStartFn = c.lineStart
		c.lineEndFn = c.lineEnd
		return c
	}
}

type clipStream struct {
	sink         Sink
	pointVisible func(lambda, phi float64) bool
	interpolate  interpolateFunc
	start        [2]float64

	line     clipLiner
	ringBuf  *clipBuffer
	ringSink clipLiner

	polygonStarted bool
	polygon        [][][2]float64
	segments       [][][3]float64
	ring           [][2]float64

	pointFn     func(lambda, phi float64)
	lineStartFn func()
	lineEndFn   func()
}

func (c *clipStream) Point(lambda, phi float64) { c.pointFn(lambda, phi) }
func (c *clipStream) LineStart()                { c.lineStartFn() }
func (c *clipStream) LineEnd()                  { c.lineEndFn() }

func (c *clipStream) Sphere() {
	c.sink.PolygonStart()
	c.sink.LineStart()
	c.interpolate(nil, nil, 1, c.sink)
	c.sink.LineEnd()
	c.sink.PolygonEnd()
}

func (c *clipStream) PolygonStart() {
	c.pointFn = c.pointRing
	c.lineStartFn = c.ringStart
	c.lineEndFn = c.ringEnd
	c.segments = nil
	c.polygon = [][][2]float64{}
}

func (c *clipStream) PolygonEnd() {
	c.pointFn = c.point
	c.lineStartFn = c.lineStart
	c.lineEndFn = c.lineEnd
	startInside := polygonContains(c.polygon, c.start)
	if len(c.segments) > 0 {
		if !c.polygonStarted {
			c.sink.PolygonStart()
			c.polygonStarted = true
		}
		clipRejoin(c.segments, compareIntersection, startInside, c.interpolate, c.sink)
	} else if startInside {
		if !c.polygonStarted {
			c.sink.PolygonStart()
			c.polygonStarted = true
		}
		c.sink.LineStart()
		c.interpolate(nil, nil, 1, c.sink)
		c.sink.LineEnd()
	}
	if c.polygonStarted {
		c.sink.PolygonEnd()
		c.polygonStarted = false
	}
	c.segments = nil
	c.polygon = nil
}

func (c *clipStream) point(lambda, phi float64) {
	if c.pointVisible(lambda, phi) {
		c.sink.Point(lambda, phi)
	}
}

func (c *clipStream) pointLine(lambda, phi float64) { c.line.Point(lambda, phi) }

func (c *clipStream) lineStart() {
	c.pointFn = c.pointLine
	c.line.LineStart()
}

func (c *clipStream) lineEnd() {
	c.pointFn = c.point
	c.line.LineEnd()
}

func (c *clipStream) pointRing(lambda, phi float64) {
	c.ring = append(c.ring, [2]float64{lambda, phi})
	c.ringSink.Point(lambda, phi)
}

func (c *clipStream) ringStart() {
	c.ringSink.LineStart()
	c.ring = [][2]float64{}
}

func (c *clipStream) ringEnd() {
	c.pointRing(c.ring[0][0], c.ring[0][1])
	c.ringSink.LineEnd()

	clean := c.ringSink.Clean()
	ringSegments := c.ringBuf.result()
	n := len(ringSegments)

	c.ring = c.ring[:len(c.ring)-1]
	c.polygon = append(c.polygon, c.ring)
	c.ring = nil

	if n == 0 {
		return
	}

	// No intersections: the ring is entirely inside or outside.
	if clean&1 != 0 {
		segment := ringSegments[0]
		m := len(segment) - 1
		if m > 0 {
			if !c.polygonStarted {
				c.sink.PolygonStart()
				c.polygonStarted = true
			}
			c.sink.LineStart()
			for i := 0; i < m; i++ {
				c.sink.Point(segment[i][0], segment[i][1])
			}
			c.sink.LineEnd()
		}
		return
	}

	// Rejoin the first and last segments if the ring wrapped the boundary.
	if n > 1 && clean&2 != 0 {
		last := ringSegments[n-1]
		first := ringSegments[0]
		merged := append(append([][3]float64{}, last...), first...)
		ringSegments = append(ringSegments[1:n-1:n-1], merged)
	}

	for _, s := range ringSegments {
		if len(s) > 1 { // validSegment
			c.segments = append(c.segments, s)
		}
	}
}
