package colors

import (
	"math"

	d3color "github.com/geoffjay/templ-charts/internal/d3/color"
)

// Space re-exports d3color.Space so callers select an interpolation space
// through the colors package (e.g. colors.SpaceLab). SpaceRGB is the zero
// value, so a zero-initialized Space preserves the default gamma-sRGB
// interpolation that every existing scale/palette golden was captured under.
type Space = d3color.Space

// Interpolation-space constants, re-exported from internal/d3/color.
const (
	SpaceRGB = d3color.SpaceRGB
	SpaceHSL = d3color.SpaceHSL
	SpaceLab = d3color.SpaceLab
	SpaceLch = d3color.SpaceLch
)

// interpolateInSpace returns a t→hex interpolator over the given color stops.
//
// For SpaceRGB it delegates to interpolateRgbBasis — byte-identical to the
// legacy path, so RGB goldens never move. For a perceptual space it blends
// piecewise-linearly between consecutive stops in that space via the validated
// d3color port (which matches d3's interpolateLab/interpolateHcl/interpolateHsl
// for the two-stop case, and carries an undefined hue/saturation through to the
// opposite endpoint, so gray stops don't produce NaN output).
func interpolateInSpace(cols []string, space Space) func(float64) string {
	if space == SpaceRGB {
		return interpolateRgbBasis(cols)
	}
	n := len(cols)
	if n == 0 {
		return func(float64) string { return "#000000" }
	}
	if n == 1 {
		c := cols[0]
		return func(float64) string { return c }
	}
	stops := make([]d3color.Color, n)
	for i, c := range cols {
		stops[i] = d3color.RGBColor(c).In(space)
	}
	return func(t float64) string {
		if t <= 0 {
			return stops[0].RGB().FormatHex()
		}
		if t >= 1 {
			return stops[n-1].RGB().FormatHex()
		}
		x := t * float64(n-1)
		i := int(math.Floor(x))
		if i >= n-1 {
			i = n - 2
		}
		return stops[i].Interpolate(stops[i+1], x-float64(i)).RGB().FormatHex()
	}
}

// schemeInterpolatorInSpace returns a t→hex interpolator for a named scheme,
// blended through the given space. For SpaceRGB it returns the pre-built RGB
// interpolator (unchanged). For a perceptual space it re-interpolates the
// scheme's color stops in that space; function-based interpolators (viridis,
// turbo, …) that have no stop array are first sampled densely in RGB, then
// re-blended in-space.
func schemeInterpolatorInSpace(scheme string, space Space) func(float64) string {
	if space == SpaceRGB {
		if f, ok := ColorInterpolators[scheme]; ok {
			return f
		}
		return ColorInterpolators["turbo"]
	}
	if stops := schemeStops(scheme); len(stops) > 0 {
		return interpolateInSpace(stops, space)
	}
	if base, ok := ColorInterpolators[scheme]; ok {
		const k = 16
		stops := make([]string, k)
		for i := 0; i < k; i++ {
			stops[i] = base(float64(i) / float64(k-1))
		}
		return interpolateInSpace(stops, space)
	}
	return ColorInterpolators["turbo"]
}

// schemeStops returns the discrete color stops backing a scheme id (the same
// arrays the RGB interpolators are built from), or nil for function-based
// interpolators with no array.
func schemeStops(scheme string) []string {
	if IsDivergingColorScheme(scheme) {
		return DivergingColorSchemes[scheme][11]
	}
	if IsSequentialColorScheme(scheme) {
		if arr, ok := SequentialColorSchemes[scheme][9]; ok {
			return arr
		}
	}
	return nil
}
