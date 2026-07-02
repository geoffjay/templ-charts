// Package geo is a pure-Go port of the subset of d3-geo needed by @nivo/geo
// (GeoMap + Choropleth): spherical projections, GeoJSON→SVG path generation,
// and a graticule generator. It mirrors d3-geo's module structure — a stream
// pipeline (radians → rotate → clip → resample → project) fed by a GeoJSON
// stream traversal — so output matches d3 for the common rendering path.
//
// Scope (v3 minimal-viable, see docs/PLAN-v3.md §3.6 and NOTES.md): all ten
// projection types nivo exposes are present, adaptive resampling and
// antimeridian clipping are faithful ports, and GeoPath emits SVG path
// strings. Deferred: clipCircle (azimuthal back-face hiding) and clipExtent
// (rectangular clip). Cylindrical/pseudocylindrical projections (mercator,
// equirectangular, transverseMercator, naturalEarth1, equalEarth) are fully
// correct; azimuthal-family projections render the whole sphere until
// clipCircle lands.
package geo

import "math"

// Constants mirror d3-geo's src/math.js.
const (
	epsilon   = 1e-6
	epsilon2  = 1e-12
	pi        = math.Pi
	halfPi    = pi / 2
	quarterPi = pi / 4
	tau       = pi * 2
	degrees   = 180 / pi
	radians   = pi / 180
)

// acos clamps to [-1,1] before Acos (d3-geo math.js).
func acos(x float64) float64 {
	switch {
	case x > 1:
		return 0
	case x < -1:
		return pi
	default:
		return math.Acos(x)
	}
}

// asin clamps to [-1,1] before Asin (d3-geo math.js).
func asin(x float64) float64 {
	switch {
	case x > 1:
		return halfPi
	case x < -1:
		return -halfPi
	default:
		return math.Asin(x)
	}
}

func haversin(x float64) float64 {
	x = math.Sin(x / 2)
	return x * x
}

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
