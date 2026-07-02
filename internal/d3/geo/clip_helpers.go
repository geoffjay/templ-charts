package geo

import "math"

// pointEqual reports whether two [x,y,...] points coincide within epsilon.
// Ported from d3-geo src/pointEqual.js.
func pointEqual(a, b [3]float64) bool {
	return math.Abs(a[0]-b[0]) < epsilon && math.Abs(a[1]-b[1]) < epsilon
}

func boolXor(a, b bool) bool { return a != b }

// polygonContains reports whether a spherical polygon (list of rings, in
// radians) contains a point. Ported from d3-geo src/polygonContains.js. Used by
// the clip framework to decide whether a fully-outside ring nonetheless
// encloses the clip start point. sum uses a plain accumulator rather than
// d3-array's compensated Adder; results agree for real-world inputs.
func polygonContains(polygon [][][2]float64, point [2]float64) bool {
	lambda := point[0]
	phi := point[1]
	sinPhi := math.Sin(phi)
	normal := [3]float64{math.Sin(lambda), -math.Cos(lambda), 0}
	angle := 0.0
	winding := 0
	sum := 0.0

	if sinPhi == 1 {
		phi = halfPi + epsilon
	} else if sinPhi == -1 {
		phi = -halfPi - epsilon
	}

	for _, ring := range polygon {
		m := len(ring)
		if m == 0 {
			continue
		}
		point0 := ring[m-1]
		lambda0 := point0[0]
		phi0 := point0[1]/2 + quarterPi
		sinPhi0 := math.Sin(phi0)
		cosPhi0 := math.Cos(phi0)

		for j := 0; j < m; j++ {
			point1 := ring[j]
			lambda1 := point1[0]
			phi1 := point1[1]/2 + quarterPi
			sinPhi1 := math.Sin(phi1)
			cosPhi1 := math.Cos(phi1)
			delta := lambda1 - lambda0
			sgn := 1.0
			if delta < 0 {
				sgn = -1
			}
			absDelta := sgn * delta
			antimeridian := absDelta > pi
			k := sinPhi0 * sinPhi1

			sum += math.Atan2(k*sgn*math.Sin(absDelta), cosPhi0*cosPhi1+k*math.Cos(absDelta))
			if antimeridian {
				angle += delta + sgn*tau
			} else {
				angle += delta
			}

			if boolXor(boolXor(antimeridian, lambda0 >= lambda), lambda1 >= lambda) {
				arc := cartesianCross(cartesian([2]float64{point0[0], point0[1]}), cartesian([2]float64{point1[0], point1[1]}))
				cartesianNormalizeInPlace(&arc)
				intersection := cartesianCross(normal, arc)
				cartesianNormalizeInPlace(&intersection)
				dir := boolXor(antimeridian, delta >= 0)
				phiArc := asin(intersection[2])
				if dir {
					phiArc = -phiArc
				}
				if phi > phiArc || (phi == phiArc && (arc[0] != 0 || arc[1] != 0)) {
					if dir {
						winding++
					} else {
						winding--
					}
				}
			}

			lambda0 = lambda1
			sinPhi0 = sinPhi1
			cosPhi0 = cosPhi1
			point0 = point1
		}
	}

	cond := angle < -epsilon || (angle < epsilon && sum < -epsilon2)
	return boolXor(cond, winding&1 == 1)
}

// clipBuffer collects clipped line segments as lists of [x,y,m] points, where m
// is an intersection marker (0 here — the antimeridian clipper adds no marks).
// Ported from d3-geo src/clip/buffer.js.
type clipBuffer struct {
	noopSink
	lines [][][3]float64
}

func (b *clipBuffer) Point(x, y float64) {
	n := len(b.lines)
	b.lines[n-1] = append(b.lines[n-1], [3]float64{x, y, 0})
}
func (b *clipBuffer) LineStart() { b.lines = append(b.lines, [][3]float64{}) }
func (b *clipBuffer) LineEnd()   {}

func (b *clipBuffer) result() [][][3]float64 {
	r := b.lines
	b.lines = nil
	return r
}
