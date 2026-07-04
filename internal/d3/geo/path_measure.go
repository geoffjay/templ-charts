package geo

import "math"

// This file ports d3-geo's planar path measurement sinks — path.area,
// path.bounds, path.centroid — as stream Sinks fed by the same projection
// stream that produces the SVG path string. Because they consume the projected
// (and clipped/resampled) point stream, they measure the shape exactly as it is
// drawn. Ported from d3-geo src/path/{area,bounds,centroid}.js.

// GeoObject is any GeoJSON value Path measurement accepts: a Geometry, a
// Feature, a FeatureCollection (value or pointer), or a slice of Features.
// Mirrors how d3-geo's geoPath accepts any GeoJSON object.
type GeoObject = any

// streamObject streams a supported GeoJSON object through sink.
func streamObject(obj GeoObject, sink Sink) {
	switch o := obj.(type) {
	case Geometry:
		StreamGeometry(o, sink)
	case *Geometry:
		StreamGeometry(*o, sink)
	case Feature:
		StreamGeometry(o.Geometry, sink)
	case *Feature:
		StreamGeometry(o.Geometry, sink)
	case FeatureCollection:
		for _, f := range o.Features {
			StreamGeometry(f.Geometry, sink)
		}
	case *FeatureCollection:
		for _, f := range o.Features {
			StreamGeometry(f.Geometry, sink)
		}
	case []Feature:
		for _, f := range o {
			StreamGeometry(f.Geometry, sink)
		}
	}
}

// Bounds returns the projected bounding box [[x0,y0],[x1,y1]] of obj, matching
// d3-geo path.bounds. An object that projects to nothing yields an inverted
// (infinite) box.
func (gp *Path) Bounds(obj GeoObject) [2][2]float64 {
	s := newBoundsSink()
	streamObject(obj, gp.projection.Stream(s))
	return s.result()
}

// Area returns the projected (planar) area of obj in square pixels, matching
// d3-geo path.area (the absolute area, summed per polygon so holes subtract).
func (gp *Path) Area(obj GeoObject) float64 {
	s := newAreaSink()
	streamObject(obj, gp.projection.Stream(s))
	return s.result()
}

// Centroid returns the projected centroid [x,y] of obj, matching d3-geo
// path.centroid: weighted by area for polygons, by length for lines, else the
// mean of points. Returns NaNs for empty input.
func (gp *Path) Centroid(obj GeoObject) [2]float64 {
	s := newCentroidSink()
	streamObject(obj, gp.projection.Stream(s))
	return s.result()
}

// --- bounds sink -----------------------------------------------------------

type boundsSink struct {
	x0, y0, x1, y1 float64
}

func newBoundsSink() *boundsSink {
	inf := math.Inf(1)
	return &boundsSink{x0: inf, y0: inf, x1: -inf, y1: -inf}
}

func (s *boundsSink) Point(x, y float64) {
	if x < s.x0 {
		s.x0 = x
	}
	if x > s.x1 {
		s.x1 = x
	}
	if y < s.y0 {
		s.y0 = y
	}
	if y > s.y1 {
		s.y1 = y
	}
}
func (s *boundsSink) LineStart()    {}
func (s *boundsSink) LineEnd()      {}
func (s *boundsSink) PolygonStart() {}
func (s *boundsSink) PolygonEnd()   {}
func (s *boundsSink) Sphere()       {}

func (s *boundsSink) result() [2][2]float64 {
	return [2][2]float64{{s.x0, s.y0}, {s.x1, s.y1}}
}

// --- area sink -------------------------------------------------------------

type areaSink struct {
	areaSum, areaRingSum float64
	x00, y00, x0, y0     float64

	point     func(x, y float64)
	lineStart func()
	lineEnd   func()
}

func newAreaSink() *areaSink {
	s := &areaSink{}
	s.point = s.noPoint
	s.lineStart = s.noLine
	s.lineEnd = s.noLine
	return s
}

func (s *areaSink) Point(x, y float64) { s.point(x, y) }
func (s *areaSink) LineStart()         { s.lineStart() }
func (s *areaSink) LineEnd()           { s.lineEnd() }
func (s *areaSink) Sphere()            {}

func (s *areaSink) PolygonStart() {
	s.lineStart = s.ringStart
	s.lineEnd = s.ringEnd
}

func (s *areaSink) PolygonEnd() {
	s.point = s.noPoint
	s.lineStart = s.noLine
	s.lineEnd = s.noLine
	s.areaSum += math.Abs(s.areaRingSum)
	s.areaRingSum = 0
}

func (s *areaSink) ringStart() { s.point = s.pointFirst }

func (s *areaSink) pointFirst(x, y float64) {
	s.point = s.pointRing
	s.x00, s.x0 = x, x
	s.y00, s.y0 = y, y
}

func (s *areaSink) pointRing(x, y float64) {
	s.areaRingSum += s.y0*x - s.x0*y
	s.x0, s.y0 = x, y
}

func (s *areaSink) ringEnd()             { s.pointRing(s.x00, s.y00) }
func (s *areaSink) noPoint(x, y float64) {}
func (s *areaSink) noLine()              {}
func (s *areaSink) result() float64      { return s.areaSum / 2 }

// --- centroid sink ---------------------------------------------------------

type centroidSink struct {
	// dimension accumulators: 0=points, 1=lines (length-weighted), 2=areas.
	x0sum, y0sum, z0 float64
	x1sum, y1sum, z1 float64
	x2sum, y2sum, z2 float64
	x00, y00, x0, y0 float64

	point     func(x, y float64)
	lineStart func()
	lineEnd   func()
}

func newCentroidSink() *centroidSink {
	s := &centroidSink{}
	s.point = s.centroidPoint
	s.lineStart = s.lineStartFn
	s.lineEnd = s.lineEndFn
	return s
}

func (s *centroidSink) Point(x, y float64) { s.point(x, y) }
func (s *centroidSink) LineStart()         { s.lineStart() }
func (s *centroidSink) LineEnd()           { s.lineEnd() }
func (s *centroidSink) Sphere()            {}

func (s *centroidSink) PolygonStart() {
	s.lineStart = s.ringStart
	s.lineEnd = s.ringEnd
}

func (s *centroidSink) PolygonEnd() {
	s.point = s.centroidPoint
	s.lineStart = s.lineStartFn
	s.lineEnd = s.lineEndFn
}

// dimension 0: bare points
func (s *centroidSink) centroidPoint(x, y float64) {
	s.x0sum += x
	s.y0sum += y
	s.z0++
}

// dimension 1: line segments (length-weighted midpoints)
func (s *centroidSink) lineStartFn() { s.point = s.pointFirstLine }
func (s *centroidSink) lineEndFn()   { s.point = s.centroidPoint }

func (s *centroidSink) pointFirstLine(x, y float64) {
	s.point = s.pointLine
	s.x0, s.y0 = x, y
	s.centroidPoint(x, y)
}

func (s *centroidSink) pointLine(x, y float64) {
	dx, dy := x-s.x0, y-s.y0
	z := math.Sqrt(dx*dx + dy*dy)
	s.x1sum += z * (s.x0 + x) / 2
	s.y1sum += z * (s.y0 + y) / 2
	s.z1 += z
	s.x0, s.y0 = x, y
	s.centroidPoint(x, y)
}

// dimension 2: polygon rings (area-weighted)
func (s *centroidSink) ringStart() { s.point = s.pointFirstRing }
func (s *centroidSink) ringEnd()   { s.pointRing(s.x00, s.y00) }

func (s *centroidSink) pointFirstRing(x, y float64) {
	s.point = s.pointRing
	s.x00, s.x0 = x, x
	s.y00, s.y0 = y, y
	s.centroidPoint(x, y)
}

func (s *centroidSink) pointRing(x, y float64) {
	dx, dy := x-s.x0, y-s.y0
	z := math.Sqrt(dx*dx + dy*dy)
	s.x1sum += z * (s.x0 + x) / 2
	s.y1sum += z * (s.y0 + y) / 2
	s.z1 += z

	z = s.y0*x - s.x0*y
	s.x2sum += z * (s.x0 + x)
	s.y2sum += z * (s.y0 + y)
	s.z2 += z * 3
	s.x0, s.y0 = x, y
	s.centroidPoint(x, y)
}

func (s *centroidSink) result() [2]float64 {
	switch {
	case s.z2 != 0:
		return [2]float64{s.x2sum / s.z2, s.y2sum / s.z2}
	case s.z1 != 0:
		return [2]float64{s.x1sum / s.z1, s.y1sum / s.z1}
	case s.z0 != 0:
		return [2]float64{s.x0sum / s.z0, s.y0sum / s.z0}
	default:
		return [2]float64{math.NaN(), math.NaN()}
	}
}
