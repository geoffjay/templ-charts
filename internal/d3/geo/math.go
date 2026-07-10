// Package geo is a pure-Go port of the subset of d3-geo needed by @nivo/geo
// (GeoMap + Choropleth): spherical projections, GeoJSON→SVG path generation,
// and a graticule generator. It mirrors d3-geo's module structure — a stream
// pipeline (radians → rotate → clip → resample → project) fed by a GeoJSON
// stream traversal — so output matches d3 for the common rendering path.
//
// Scope: all ten projection types nivo exposes are present, adaptive resampling
// and antimeridian clipping are faithful ports, and GeoPath emits SVG path
// strings. The clip machinery also includes clipCircle (the azimuthal
// small-circle preclip that hides the far hemisphere) and clipExtent (the
// rectangular screen-space postclip), so every projection — cylindrical,
// pseudocylindrical, and azimuthal (orthographic/gnomonic/stereographic/
// azimuthal*) — renders correctly.
//
// The geo-measurement completeness pieces are also present:
// GeoPath.Bounds/Area/Centroid (planar stream sinks; see path_measure.go),
// Projection.FitExtent/FitSize/FitWidth/FitHeight (fit.go), the conic projection
// family — conicConformal/conicEqualArea/conicEquidistant with standard-parallels
// support (conic.go) — and an opt-in TopoJSON decoder (topojson.go; callers may
// still supply GeoJSON directly).
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
