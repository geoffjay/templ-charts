package geo

// passThrough wraps a downstream Sink, forwarding every callback except Point,
// which it maps through pointFn. Used for the radians/rotate/scale-translate
// transform streams and resampleNone. Mirrors d3-geo's transformer().
type passThrough struct {
	stream  Sink
	pointFn func(x, y float64)
}

func (t *passThrough) Point(x, y float64) { t.pointFn(x, y) }
func (t *passThrough) LineStart()         { t.stream.LineStart() }
func (t *passThrough) LineEnd()           { t.stream.LineEnd() }
func (t *passThrough) PolygonStart()      { t.stream.PolygonStart() }
func (t *passThrough) PolygonEnd()        { t.stream.PolygonEnd() }
func (t *passThrough) Sphere()            { t.stream.Sphere() }

// transformRadians converts incoming (lon,lat) degrees to radians.
func transformRadians(sink Sink) Sink {
	return &passThrough{stream: sink, pointFn: func(x, y float64) {
		sink.Point(x*radians, y*radians)
	}}
}

// transformRotate applies a rotation before forwarding each point.
func transformRotate(rotate transform) func(Sink) Sink {
	return func(sink Sink) Sink {
		return &passThrough{stream: sink, pointFn: func(x, y float64) {
			lx, ly := rotate.forward(x, y)
			sink.Point(lx, ly)
		}}
	}
}
