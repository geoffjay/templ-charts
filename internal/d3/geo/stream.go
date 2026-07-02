package geo

// StreamGeometry walks a Geometry, emitting its primitives into sink. Mirrors
// d3-geo's src/stream.js streamGeometry.
func StreamGeometry(g Geometry, sink Sink) {
	switch g.Type {
	case TypeSphere:
		sink.Sphere()
	case TypePoint:
		if p, ok := g.Coordinates.([2]float64); ok {
			sink.Point(p[0], p[1])
		}
	case TypeMultiPoint:
		if pts, ok := g.Coordinates.([][2]float64); ok {
			for _, p := range pts {
				sink.Point(p[0], p[1])
			}
		}
	case TypeLineString:
		if pts, ok := g.Coordinates.([][2]float64); ok {
			streamLine(pts, sink, false)
		}
	case TypeMultiLineString:
		if lines, ok := g.Coordinates.([][][2]float64); ok {
			for _, l := range lines {
				streamLine(l, sink, false)
			}
		}
	case TypePolygon:
		if rings, ok := g.Coordinates.([][][2]float64); ok {
			streamPolygon(rings, sink)
		}
	case TypeMultiPolygon:
		if polys, ok := g.Coordinates.([][][][2]float64); ok {
			for _, poly := range polys {
				streamPolygon(poly, sink)
			}
		}
	case TypeGeometryColl:
		for _, sub := range g.Geometries {
			StreamGeometry(sub, sink)
		}
	}
}

// streamLine emits a line; closed=true drops the final (repeated) coordinate,
// as when streaming polygon rings.
func streamLine(coords [][2]float64, sink Sink, closed bool) {
	n := len(coords)
	if closed {
		n--
	}
	sink.LineStart()
	for i := 0; i < n; i++ {
		sink.Point(coords[i][0], coords[i][1])
	}
	sink.LineEnd()
}

func streamPolygon(rings [][][2]float64, sink Sink) {
	sink.PolygonStart()
	for _, ring := range rings {
		streamLine(ring, sink, true)
	}
	sink.PolygonEnd()
}

// StreamFeature streams a Feature's geometry.
func StreamFeature(f Feature, sink Sink) { StreamGeometry(f.Geometry, sink) }
