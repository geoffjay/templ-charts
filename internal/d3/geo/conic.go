package geo

import "math"

// Conic projection family, ported from d3-geo src/projection/{conic,
// conicConformal,conicEqualArea,conicEquidistant,cylindricalEqualArea}.js.
// Conic projections are parameterized by two standard parallels (see
// Projection.Parallels). Following d3-geo's conicProjection, the default
// parallels are [0°, 60°]; degenerate parallels fall back to the matching
// cylindrical projection.

// newConicProjection wraps a raw-conic factory with the standard machinery and
// the conic default parallels [0, π/3].
func newConicProjection(factory func(phi0, phi1 float64) transform) *Projection {
	phi0, phi1 := 0.0, pi/3
	p := newProjection(factory(phi0, phi1))
	p.rawFactory = factory
	p.phi0, p.phi1 = phi0, phi1
	return p
}

func tany(y float64) float64 { return math.Tan((halfPi + y) / 2) }

// conicConformalRaw is the Lambert conformal conic. Degenerates to Mercator
// when the parallels coincide at n=0.
func conicConformalRaw(y0, y1 float64) transform {
	cy0 := math.Cos(y0)
	var n float64
	if y0 == y1 {
		n = math.Sin(y0)
	} else {
		n = math.Log(cy0/math.Cos(y1)) / math.Log(tany(y1)/tany(y0))
	}
	if n == 0 {
		return mercatorRaw()
	}
	f := cy0 * math.Pow(tany(y0), n) / n
	return transform{
		forward: func(x, y float64) (float64, float64) {
			if f > 0 {
				if y < -halfPi+epsilon {
					y = -halfPi + epsilon
				}
			} else {
				if y > halfPi-epsilon {
					y = halfPi - epsilon
				}
			}
			r := f / math.Pow(tany(y), n)
			return r * math.Sin(n*x), f - r*math.Cos(n*x)
		},
		invert: func(x, y float64) (float64, float64) {
			fy := f - y
			r := sign(n) * math.Sqrt(x*x+fy*fy)
			l := math.Atan2(x, math.Abs(fy)) * sign(fy)
			if fy*n < 0 {
				l -= pi * sign(x) * sign(fy)
			}
			return l / n, 2*math.Atan(math.Pow(f/r, 1/n)) - halfPi
		},
	}
}

// cylindricalEqualAreaRaw is the fallback for conicEqualArea when the parallels
// coincide at n≈0.
func cylindricalEqualAreaRaw(phi0 float64) transform {
	cosPhi0 := math.Cos(phi0)
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			return lambda * cosPhi0, math.Sin(phi) / cosPhi0
		},
		invert: func(x, y float64) (float64, float64) {
			return x / cosPhi0, asin(y * cosPhi0)
		},
	}
}

// conicEqualAreaRaw is the Albers equal-area conic. Degenerates to a
// cylindrical equal-area projection when n≈0.
func conicEqualAreaRaw(y0, y1 float64) transform {
	sy0 := math.Sin(y0)
	n := (sy0 + math.Sin(y1)) / 2
	if math.Abs(n) < epsilon {
		return cylindricalEqualAreaRaw(y0)
	}
	c := 1 + sy0*(2*n-sy0)
	r0 := math.Sqrt(c) / n
	return transform{
		forward: func(x, y float64) (float64, float64) {
			r := math.Sqrt(c-2*n*math.Sin(y)) / n
			nx := n * x
			return r * math.Sin(nx), r0 - r*math.Cos(nx)
		},
		invert: func(x, y float64) (float64, float64) {
			r0y := r0 - y
			l := math.Atan2(x, math.Abs(r0y)) * sign(r0y)
			if r0y*n < 0 {
				l -= pi * sign(x) * sign(r0y)
			}
			return l / n, asin((c - (x*x+r0y*r0y)*n*n) / (2 * n))
		},
	}
}

// conicEquidistantRaw is the equidistant conic. Degenerates to equirectangular
// when n≈0.
func conicEquidistantRaw(y0, y1 float64) transform {
	cy0 := math.Cos(y0)
	var n float64
	if y0 == y1 {
		n = math.Sin(y0)
	} else {
		n = (cy0 - math.Cos(y1)) / (y1 - y0)
	}
	if math.Abs(n) < epsilon {
		return equirectangularRaw()
	}
	g := cy0/n + y0
	return transform{
		forward: func(x, y float64) (float64, float64) {
			gy := g - y
			nx := n * x
			return gy * math.Sin(nx), g - gy*math.Cos(nx)
		},
		invert: func(x, y float64) (float64, float64) {
			gy := g - y
			l := math.Atan2(x, math.Abs(gy)) * sign(gy)
			if gy*n < 0 {
				l -= pi * sign(x) * sign(gy)
			}
			return l / n, g - sign(n)*math.Sqrt(x*x+gy*gy)
		},
	}
}

// GeoConicConformal returns a Lambert conformal conic projection (default
// parallels 0°,60°; set your own with Parallels).
func GeoConicConformal() *Projection { return newConicProjection(conicConformalRaw) }

// GeoConicEqualArea returns an Albers equal-area conic projection.
func GeoConicEqualArea() *Projection { return newConicProjection(conicEqualAreaRaw) }

// GeoConicEquidistant returns an equidistant conic projection.
func GeoConicEquidistant() *Projection { return newConicProjection(conicEquidistantRaw) }
