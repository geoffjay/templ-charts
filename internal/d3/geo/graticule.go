package geo

import "math"

// Graticule returns the default d3-geo graticule as a MultiLineString geometry:
// meridians and parallels every 10°, with major lines (every 90° of longitude,
// plus the equator) spanning the full ±90°/±180° extent and minor lines
// trimmed to ±80°/±180°. Densified with precision 2.5°. Ported from d3-geo
// src/graticule.js (geoGraticule with its default parameters).
func Graticule() Geometry {
	const (
		dx        = 10.0 // minor step (lon/lat)
		dy        = 10.0
		bigDX     = 90.0  // major meridian step
		bigDY     = 360.0 // major parallel step
		precision = 2.5
	)
	// extentMinor / extentMajor.
	x0, x1 := -180.0, 180.0
	y0, y1 := -80.0-epsilon, 80.0+epsilon
	bigX0, bigX1 := -180.0, 180.0
	bigY0, bigY1 := -90.0+epsilon, 90.0-epsilon

	minorMeridian := graticuleX(y0, y1, 90)
	minorParallel := graticuleY(x0, x1, 90)
	majorMeridian := graticuleX(bigY0, bigY1, precision)
	majorParallel := graticuleY(bigX0, bigX1, precision)

	var lines [][][2]float64

	for _, x := range rangeF(math.Ceil(bigX0/bigDX)*bigDX, bigX1, bigDX) {
		lines = append(lines, majorMeridian(x))
	}
	for _, y := range rangeF(math.Ceil(bigY0/bigDY)*bigDY, bigY1, bigDY) {
		lines = append(lines, majorParallel(y))
	}
	for _, x := range rangeF(math.Ceil(x0/dx)*dx, x1, dx) {
		if math.Abs(math.Mod(x, bigDX)) > epsilon {
			lines = append(lines, minorMeridian(x))
		}
	}
	for _, y := range rangeF(math.Ceil(y0/dy)*dy, y1, dy) {
		if math.Abs(math.Mod(y, bigDY)) > epsilon {
			lines = append(lines, minorParallel(y))
		}
	}

	return Geometry{Type: TypeMultiLineString, Coordinates: lines}
}

// graticuleX returns a meridian generator: given a longitude, a line of points
// from latitude y0 to y1 stepping dy (y1 appended).
func graticuleX(y0, y1, dy float64) func(x float64) [][2]float64 {
	ys := append(rangeF(y0, y1-epsilon, dy), y1)
	return func(x float64) [][2]float64 {
		out := make([][2]float64, len(ys))
		for i, y := range ys {
			out[i] = [2]float64{x, y}
		}
		return out
	}
}

// graticuleY returns a parallel generator: given a latitude, a line of points
// from longitude x0 to x1 stepping dx (x1 appended).
func graticuleY(x0, x1, dx float64) func(y float64) [][2]float64 {
	xs := append(rangeF(x0, x1-epsilon, dx), x1)
	return func(y float64) [][2]float64 {
		out := make([][2]float64, len(xs))
		for i, x := range xs {
			out[i] = [2]float64{x, y}
		}
		return out
	}
}

// rangeF mirrors d3.range(start, stop, step): [start, start+step, …) < stop.
func rangeF(start, stop, step float64) []float64 {
	n := int(math.Max(0, math.Ceil((stop-start)/step)))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = start + float64(i)*step
	}
	return out
}
