package geo

import "math"

// Raw projections and their constructors, ported from the corresponding
// d3-geo src/projection/*.js modules. Each constructor returns a *Projection
// wrapping the raw transform with the standard machinery. nivo always sets
// scale/translate/rotate explicitly, so the d3 per-projection default scales
// and clip angles are not replicated here (see NOTES.md for the deferred
// clipCircle that the azimuthal family needs to hide the far hemisphere).

// --- Cylindrical / pseudocylindrical (fully correct with antimeridian clip) ---

func mercatorRaw() transform {
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			return lambda, math.Log(math.Tan((halfPi + phi) / 2))
		},
		invert: func(x, y float64) (float64, float64) {
			return x, 2*math.Atan(math.Exp(y)) - halfPi
		},
	}
}

func equirectangularRaw() transform {
	id := func(a, b float64) (float64, float64) { return a, b }
	return transform{forward: id, invert: id}
}

func transverseMercatorRaw() transform {
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			return math.Log(math.Tan((halfPi + phi) / 2)), -lambda
		},
		invert: func(x, y float64) (float64, float64) {
			return -y, 2*math.Atan(math.Exp(x)) - halfPi
		},
	}
}

func naturalEarth1Raw() transform {
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			phi2 := phi * phi
			phi4 := phi2 * phi2
			return lambda * (0.8707 - 0.131979*phi2 + phi4*(-0.013791+phi4*(0.003971*phi2-0.001529*phi4))),
				phi * (1.007226 + phi2*(0.015085+phi4*(-0.044475+0.028874*phi2-0.005916*phi4)))
		},
	}
}

func equalEarthRaw() transform {
	const (
		a1 = 1.340264
		a2 = -0.081106
		a3 = 0.000893
		a4 = 0.003796
	)
	m := math.Sqrt(3) / 2
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			l := asin(m * math.Sin(phi))
			l2 := l * l
			l6 := l2 * l2 * l2
			return lambda * math.Cos(l) / (m * (a1 + 3*a2*l2 + l6*(7*a3+9*a4*l2))),
				l * (a1 + a2*l2 + l6*(a3+a4*l2))
		},
	}
}

// --- Azimuthal family (front hemisphere correct; far side needs clipCircle) ---

func azimuthalRaw(scale func(cxcy float64) float64) func(lambda, phi float64) (float64, float64) {
	return func(x, y float64) (float64, float64) {
		cx := math.Cos(x)
		cy := math.Cos(y)
		k := scale(cx * cy)
		if math.IsInf(k, 0) {
			return 2, 0
		}
		return k * cy * math.Sin(x), k * math.Sin(y)
	}
}

func azimuthalInvert(angle func(z float64) float64) func(x, y float64) (float64, float64) {
	return func(x, y float64) (float64, float64) {
		z := math.Sqrt(x*x + y*y)
		c := angle(z)
		sc := math.Sin(c)
		cc := math.Cos(c)
		yy := 0.0
		if z != 0 {
			yy = y * sc / z
		}
		return math.Atan2(x*sc, z*cc), asin(yy)
	}
}

func azimuthalEqualAreaRaw() transform {
	return transform{
		forward: azimuthalRaw(func(cxcy float64) float64 { return math.Sqrt(2 / (1 + cxcy)) }),
		invert:  azimuthalInvert(func(z float64) float64 { return 2 * asin(z/2) }),
	}
}

func azimuthalEquidistantRaw() transform {
	return transform{
		forward: azimuthalRaw(func(c float64) float64 {
			c = acos(c)
			if c == 0 {
				return 0
			}
			return c / math.Sin(c)
		}),
		invert: azimuthalInvert(func(z float64) float64 { return z }),
	}
}

func gnomonicRaw() transform {
	return transform{
		forward: func(x, y float64) (float64, float64) {
			cy := math.Cos(y)
			k := math.Cos(x) * cy
			return cy * math.Sin(x) / k, math.Sin(y) / k
		},
		invert: azimuthalInvert(math.Atan),
	}
}

func orthographicRaw() transform {
	return transform{
		forward: func(x, y float64) (float64, float64) {
			return math.Cos(y) * math.Sin(x), math.Sin(y)
		},
		invert: azimuthalInvert(asin),
	}
}

func stereographicRaw() transform {
	return transform{
		forward: func(x, y float64) (float64, float64) {
			cy := math.Cos(y)
			k := 1 + math.Cos(x)*cy
			return cy * math.Sin(x) / k, math.Sin(y) / k
		},
		invert: azimuthalInvert(func(z float64) float64 { return 2 * math.Atan(z) }),
	}
}

// Constructors mirror the geo* factory names.

func GeoMercator() *Projection             { return newProjection(mercatorRaw()) }
func GeoEquirectangular() *Projection      { return newProjection(equirectangularRaw()) }
func GeoTransverseMercator() *Projection   { return newProjection(transverseMercatorRaw()) }
func GeoNaturalEarth1() *Projection        { return newProjection(naturalEarth1Raw()) }
func GeoEqualEarth() *Projection           { return newProjection(equalEarthRaw()) }
func GeoAzimuthalEqualArea() *Projection   { return newProjection(azimuthalEqualAreaRaw()) }
func GeoAzimuthalEquidistant() *Projection { return newProjection(azimuthalEquidistantRaw()) }
func GeoGnomonic() *Projection             { return newProjection(gnomonicRaw()) }
func GeoOrthographic() *Projection         { return newProjection(orthographicRaw()) }
func GeoStereographic() *Projection        { return newProjection(stereographicRaw()) }

// ProjectionByType maps the nivo projectionType id to a fresh projection.
// Unknown ids fall back to mercator. Mirrors @nivo/geo's projectionById.
func ProjectionByType(name string) *Projection {
	switch name {
	case "azimuthalEqualArea":
		return GeoAzimuthalEqualArea()
	case "azimuthalEquidistant":
		return GeoAzimuthalEquidistant()
	case "gnomonic":
		return GeoGnomonic()
	case "orthographic":
		return GeoOrthographic()
	case "stereographic":
		return GeoStereographic()
	case "equalEarth":
		return GeoEqualEarth()
	case "equirectangular":
		return GeoEquirectangular()
	case "transverseMercator":
		return GeoTransverseMercator()
	case "naturalEarth1":
		return GeoNaturalEarth1()
	case "mercator":
		fallthrough
	default:
		return GeoMercator()
	}
}
