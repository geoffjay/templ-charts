package geo

import "math"

// Spherical/cartesian conversions and vector helpers, ported from d3-geo's
// src/cartesian.js. Used by resampling and the clip machinery.

func spherical(c [3]float64) [2]float64 {
	return [2]float64{math.Atan2(c[1], c[0]), asin(c[2])}
}

func cartesian(s [2]float64) [3]float64 {
	lambda, phi := s[0], s[1]
	cosPhi := math.Cos(phi)
	return [3]float64{cosPhi * math.Cos(lambda), cosPhi * math.Sin(lambda), math.Sin(phi)}
}

func cartesianDot(a, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func cartesianCross(a, b [3]float64) [3]float64 {
	return [3]float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func cartesianAddInPlace(a *[3]float64, b [3]float64) {
	a[0] += b[0]
	a[1] += b[1]
	a[2] += b[2]
}

func cartesianScale(v [3]float64, k float64) [3]float64 {
	return [3]float64{v[0] * k, v[1] * k, v[2] * k}
}

func cartesianNormalizeInPlace(d *[3]float64) {
	l := math.Sqrt(d[0]*d[0] + d[1]*d[1] + d[2]*d[2])
	d[0] /= l
	d[1] /= l
	d[2] /= l
}
