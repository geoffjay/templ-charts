// Package geo is a pure-Go port of the subset of d3-geo needed by @nivo/geo
// (GeoMap + Choropleth): spherical projections, GeoJSON→SVG path generation,
// and a graticule generator. It mirrors d3-geo's module structure — a stream
// pipeline (radians → rotate → clip → resample → project) fed by a GeoJSON
// stream traversal — so output matches d3 for the common rendering path.
//
// Scope: all ten projection types nivo exposes are present, adaptive resampling
// and antimeridian clipping are faithful ports, and GeoPath emits SVG path
// strings. As of v5 the clip machinery also includes clipCircle (the azimuthal
// small-circle preclip that hides the far hemisphere) and clipExtent (the
// rectangular screen-space postclip), so every projection — cylindrical,
// pseudocylindrical, and azimuthal (orthographic/gnomonic/stereographic/
// azimuthal*) — renders correctly. Still deferred (GeoMap/Choropleth don't need
// them): GeoPath bounds/area/centroid and fitExtent/fitSize; projections beyond
// the ~10 nivo exposes; TopoJSON (callers supply GeoJSON). See docs/PLAN-v5.md
// §3 and NOTES.md.
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
