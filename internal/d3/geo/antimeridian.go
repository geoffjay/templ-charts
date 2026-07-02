package geo

import "math"

// clipAntimeridian is the default preclip: it splits geometry that crosses the
// ±180° meridian so it renders correctly on cylindrical/pseudocylindrical
// projections. Ported from d3-geo src/clip/antimeridian.js.
var clipAntimeridian = newClip(
	func(lambda, phi float64) bool { return true },
	func(sink Sink) clipLiner { return newAntimeridianLine(sink) },
	clipAntimeridianInterpolate,
	[2]float64{-pi, -halfPi},
)

// newAntimeridianLine seeds the previous-point state to NaN so the first
// Point of a line takes neither the pole- nor antimeridian-crossing branch
// (matching d3, where these vars start as NaN).
func newAntimeridianLine(sink Sink) *antimeridianLine {
	return &antimeridianLine{
		stream:  sink,
		lambda0: math.NaN(),
		phi0:    math.NaN(),
		sign0:   math.NaN(),
	}
}

type antimeridianLine struct {
	stream               Sink
	lambda0, phi0, sign0 float64
	clean                int
}

func (l *antimeridianLine) LineStart() {
	l.stream.LineStart()
	l.clean = 1
}

func (l *antimeridianLine) Point(lambda1, phi1 float64) {
	sign1 := pi
	if lambda1 <= 0 {
		sign1 = -pi
	}
	delta := math.Abs(lambda1 - l.lambda0)
	switch {
	case math.Abs(delta-pi) < epsilon: // line crosses a pole
		if (l.phi0+phi1)/2 > 0 {
			l.phi0 = halfPi
		} else {
			l.phi0 = -halfPi
		}
		l.stream.Point(l.lambda0, l.phi0)
		l.stream.Point(l.sign0, l.phi0)
		l.stream.LineEnd()
		l.stream.LineStart()
		l.stream.Point(sign1, l.phi0)
		l.stream.Point(lambda1, l.phi0)
		l.clean = 0
	case l.sign0 != sign1 && delta >= pi: // line crosses the antimeridian
		if math.Abs(l.lambda0-l.sign0) < epsilon {
			l.lambda0 -= l.sign0 * epsilon
		}
		if math.Abs(lambda1-sign1) < epsilon {
			lambda1 -= sign1 * epsilon
		}
		l.phi0 = clipAntimeridianIntersect(l.lambda0, l.phi0, lambda1, phi1)
		l.stream.Point(l.sign0, l.phi0)
		l.stream.LineEnd()
		l.stream.LineStart()
		l.stream.Point(sign1, l.phi0)
		l.clean = 0
	}
	l.lambda0 = lambda1
	l.phi0 = phi1
	l.stream.Point(l.lambda0, l.phi0)
	l.sign0 = sign1
}

func (l *antimeridianLine) LineEnd() {
	l.stream.LineEnd()
	l.lambda0 = math.NaN()
	l.phi0 = math.NaN()
}

func (l *antimeridianLine) Clean() int { return 2 - l.clean }

func clipAntimeridianIntersect(lambda0, phi0, lambda1, phi1 float64) float64 {
	sinLambda0Lambda1 := math.Sin(lambda0 - lambda1)
	if math.Abs(sinLambda0Lambda1) <= epsilon {
		return (phi0 + phi1) / 2
	}
	cosPhi0 := math.Cos(phi0)
	cosPhi1 := math.Cos(phi1)
	return math.Atan((math.Sin(phi0)*cosPhi1*math.Sin(lambda1) - math.Sin(phi1)*cosPhi0*math.Sin(lambda0)) /
		(cosPhi0 * cosPhi1 * sinLambda0Lambda1))
}

func clipAntimeridianInterpolate(from, to *[3]float64, direction float64, stream Sink) {
	var phi float64
	if from == nil {
		phi = direction * halfPi
		stream.Point(-pi, phi)
		stream.Point(0, phi)
		stream.Point(pi, phi)
		stream.Point(pi, 0)
		stream.Point(pi, -phi)
		stream.Point(0, -phi)
		stream.Point(-pi, -phi)
		stream.Point(-pi, 0)
		stream.Point(-pi, phi)
	} else if math.Abs(from[0]-to[0]) > epsilon {
		lambda := pi
		if from[0] >= to[0] {
			lambda = -pi
		}
		phi = direction * lambda / 2
		stream.Point(-lambda, phi)
		stream.Point(0, phi)
		stream.Point(lambda, phi)
	} else {
		stream.Point(to[0], to[1])
	}
}
