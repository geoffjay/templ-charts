package geo

import "math"

// Rotation transforms, ported from d3-geo src/rotation.js. rotateRadians builds
// the (lon,lat)→(lon',lat') rotation for the three-axis projection rotate
// parameter [λ,φ,γ] (already converted to radians).

func rotationIdentity() transform {
	f := func(lambda, phi float64) (float64, float64) {
		if math.Abs(lambda) > pi {
			lambda -= math.Round(lambda/tau) * tau
		}
		return lambda, phi
	}
	return transform{forward: f, invert: f}
}

func rotateRadians(deltaLambda, deltaPhi, deltaGamma float64) transform {
	deltaLambda = math.Mod(deltaLambda, tau)
	switch {
	case deltaLambda != 0:
		if deltaPhi != 0 || deltaGamma != 0 {
			return compose(rotationLambda(deltaLambda), rotationPhiGamma(deltaPhi, deltaGamma))
		}
		return rotationLambda(deltaLambda)
	case deltaPhi != 0 || deltaGamma != 0:
		return rotationPhiGamma(deltaPhi, deltaGamma)
	default:
		return rotationIdentity()
	}
}

func forwardRotationLambda(deltaLambda float64) func(float64, float64) (float64, float64) {
	return func(lambda, phi float64) (float64, float64) {
		lambda += deltaLambda
		switch {
		case lambda > pi:
			lambda -= tau
		case lambda < -pi:
			lambda += tau
		}
		return lambda, phi
	}
}

func rotationLambda(deltaLambda float64) transform {
	return transform{
		forward: forwardRotationLambda(deltaLambda),
		invert:  forwardRotationLambda(-deltaLambda),
	}
}

func rotationPhiGamma(deltaPhi, deltaGamma float64) transform {
	cosDeltaPhi := math.Cos(deltaPhi)
	sinDeltaPhi := math.Sin(deltaPhi)
	cosDeltaGamma := math.Cos(deltaGamma)
	sinDeltaGamma := math.Sin(deltaGamma)
	return transform{
		forward: func(lambda, phi float64) (float64, float64) {
			cosPhi := math.Cos(phi)
			x := math.Cos(lambda) * cosPhi
			y := math.Sin(lambda) * cosPhi
			z := math.Sin(phi)
			k := z*cosDeltaPhi + x*sinDeltaPhi
			return math.Atan2(y*cosDeltaGamma-k*sinDeltaGamma, x*cosDeltaPhi-z*sinDeltaPhi),
				asin(k*cosDeltaGamma + y*sinDeltaGamma)
		},
		invert: func(lambda, phi float64) (float64, float64) {
			cosPhi := math.Cos(phi)
			x := math.Cos(lambda) * cosPhi
			y := math.Sin(lambda) * cosPhi
			z := math.Sin(phi)
			k := z*cosDeltaGamma - y*sinDeltaGamma
			return math.Atan2(y*cosDeltaGamma+z*sinDeltaGamma, x*cosDeltaPhi+k*sinDeltaPhi),
				asin(k*cosDeltaPhi - x*sinDeltaPhi)
		},
	}
}
