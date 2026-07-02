package geo

// interpolateFunc walks the clip boundary from `from` to `to` in the given
// direction, emitting boundary points into stream. from/to are nil for a full
// sphere traversal. Mirrors the interpolate callback in d3-geo's clip modules.
type interpolateFunc func(from, to *[3]float64, direction float64, stream Sink)

// intersection is a node in the doubly-linked subject/clip rings used by the
// rejoin algorithm. Ported from d3-geo src/clip/rejoin.js Intersection.
type intersection struct {
	x    [3]float64    // the intersection point [x,y,m]
	z    [][3]float64  // the segment this node begins/ends (subject nodes only)
	o    *intersection // matching node on the other ring
	e    bool          // is this an entry point?
	v    bool          // visited?
	n, p *intersection // next / previous around the ring
}

func linkIntersections(arr []*intersection) {
	n := len(arr)
	if n == 0 {
		return
	}
	a := arr[0]
	for i := 1; i < n; i++ {
		b := arr[i]
		a.n = b
		b.p = a
		a = b
	}
	b := arr[0]
	a.n = b
	b.p = a
}

// clipRejoin stitches clipped line segments back into closed rings, walking
// between subject (geometry) segments and interpolated clip-boundary arcs.
// Ported from d3-geo src/clip/rejoin.js.
func clipRejoin(segments [][][3]float64, compare func(a, b *intersection) float64, startInside bool, interpolate interpolateFunc, stream Sink) {
	var subject, clip []*intersection

	for _, segment := range segments {
		n := len(segment) - 1
		if n <= 0 {
			continue
		}
		p0 := segment[0]
		p1 := segment[n]

		if pointEqual(p0, p1) {
			if p0[2] == 0 && p1[2] == 0 {
				stream.LineStart()
				for i := 0; i < n; i++ {
					stream.Point(segment[i][0], segment[i][1])
				}
				stream.LineEnd()
				continue
			}
			// Degenerate: nudge the last point so it differs from the first.
			segment[n][0] += 2 * epsilon
			p1 = segment[n]
		}

		a := &intersection{x: p0, z: segment, e: true}
		subject = append(subject, a)
		ao := &intersection{x: p0, o: a, e: false}
		a.o = ao
		clip = append(clip, ao)

		b := &intersection{x: p1, z: segment, e: false}
		subject = append(subject, b)
		bo := &intersection{x: p1, o: b, e: true}
		b.o = bo
		clip = append(clip, bo)
	}

	if len(subject) == 0 {
		return
	}

	insertionSortIntersections(clip, compare)
	linkIntersections(subject)
	linkIntersections(clip)

	for i := range clip {
		startInside = !startInside
		clip[i].e = startInside
	}

	start := subject[0]
	for {
		current := start
		isSubject := true
		for current.v {
			current = current.n
			if current == start {
				return
			}
		}
		points := current.z
		stream.LineStart()
		for {
			current.v = true
			current.o.v = true
			if current.e {
				if isSubject {
					for i := 0; i < len(points); i++ {
						stream.Point(points[i][0], points[i][1])
					}
				} else {
					interpolate(&current.x, &current.n.x, 1, stream)
				}
				current = current.n
			} else {
				if isSubject {
					points = current.p.z
					for i := len(points) - 1; i >= 0; i-- {
						stream.Point(points[i][0], points[i][1])
					}
				} else {
					interpolate(&current.x, &current.p.x, -1, stream)
				}
				current = current.p
			}
			current = current.o
			points = current.z
			isSubject = !isSubject
			if current.v {
				break
			}
		}
		stream.LineEnd()
	}
}

// insertionSortIntersections orders the clip ring by compare. A stable sort is
// required (d3 uses Array.sort, stable in modern engines); insertion sort keeps
// it dependency-free and stable for the small intersection counts here.
func insertionSortIntersections(a []*intersection, compare func(x, y *intersection) float64) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && compare(a[j-1], a[j]) > 0; j-- {
			a[j-1], a[j] = a[j], a[j-1]
		}
	}
}

// compareIntersection orders antimeridian intersections along the ±180° edge,
// treating the two hemispheres as one continuous boundary. Ported from d3-geo
// src/clip/index.js.
func compareIntersection(a, b *intersection) float64 {
	return intersectionOrder(a.x) - intersectionOrder(b.x)
}

func intersectionOrder(p [3]float64) float64 {
	if p[0] < 0 {
		return p[1] - halfPi - epsilon
	}
	return halfPi - p[1]
}
