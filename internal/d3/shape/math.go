// math.go — math constants and helpers shared by d3-shape ports, matching
// d3-shape's src/math.js.
package d3shape

import "math"

const (
	epsilon    = 1e-12
	pi         = math.Pi
	halfPi     = pi / 2
	tau        = 2 * pi
	tauEpsilon = tau - epsilon
)

// acos clamps to [-1,1] (d3-shape math.js).
func acos(x float64) float64 {
	if x > 1 {
		return 0
	}
	if x < -1 {
		return pi
	}
	return math.Acos(x)
}

// asin clamps to [-1,1] (d3-shape math.js).
func asin(x float64) float64 {
	if x >= 1 {
		return halfPi
	}
	if x <= -1 {
		return -halfPi
	}
	return math.Asin(x)
}
