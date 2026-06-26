package colors

import (
	"math"
	"strconv"

	d3color "github.com/geoffjay/templ-charts/internal/d3/color"
)

// ColorInterpolatorId identifies a named color interpolator (viridis, turbo,
// rainbow, etc.). Mirrors @nivo/colors colorInterpolatorIds.
type ColorInterpolatorId string

// ColorInterpolators maps interpolator id → func(t float64) string (CSS color).
// Sequential/diverging interpolators are derived from their scheme endpoints
// via interpolateRgbBasis (matching d3's ramp). Cyclical interpolators
// (rainbow, sinebow) and the named multi-hue interpolators (viridis, inferno,
// magma, plasma, warm, cool, cubehelixDefault, turbo, cividis) are implemented
// as functions.
var ColorInterpolators = map[string]func(float64) string{}

// ColorInterpolatorIds is the ordered list of interpolator ids.
var ColorInterpolatorIds []string

func init() {
	// Diverging interpolators: derived from the size-11 scheme arrays.
	for _, id := range divergingColorSchemeIds {
		scheme := DivergingColorSchemes[id][11]
		ColorInterpolators[id] = interpolateRgbBasis(scheme)
	}
	// Sequential interpolators: derived from the largest (size-9) scheme arrays.
	for _, id := range sequentialColorSchemeIds {
		if scheme, ok := SequentialColorSchemes[id]; ok {
			if arr, ok := scheme[9]; ok {
				ColorInterpolators[id] = interpolateRgbBasis(arr)
			}
		}
	}
	// Named multi-hue interpolators (no scheme array; use known gradients).
	ColorInterpolators["turbo"] = interpolateTurbo
	ColorInterpolators["viridis"] = interpolateViridis
	ColorInterpolators["inferno"] = interpolateInferno
	ColorInterpolators["magma"] = interpolateMagma
	ColorInterpolators["plasma"] = interpolatePlasma
	ColorInterpolators["warm"] = interpolateWarm
	ColorInterpolators["cool"] = interpolateCool
	ColorInterpolators["cubehelixDefault"] = interpolateCubehelixDefault
	ColorInterpolators["cividis"] = interpolateCividis
	// Cyclical.
	ColorInterpolators["rainbow"] = interpolateRainbow
	ColorInterpolators["sinebow"] = interpolateSinebow

	ColorInterpolatorIds = make([]string, 0, len(ColorInterpolators))
	for id := range ColorInterpolators {
		ColorInterpolatorIds = append(ColorInterpolatorIds, id)
	}
}

// interpolateRgbBasis returns a function that interpolates smoothly through
// the given color stops using a cubic B-spline in RGB space. Matches d3's
// interpolateRgbBasis (used by d3-scale-chromatic's ramp).
func interpolateRgbBasis(colors []string) func(float64) string {
	if len(colors) == 0 {
		return func(float64) string { return "#000000" }
	}
	if len(colors) == 1 {
		c := colors[0]
		return func(float64) string { return c }
	}
	pts := make([][3]float64, len(colors))
	for i, c := range colors {
		r, g, b := parseHexRGB(c)
		pts[i] = [3]float64{r, g, b}
	}
	return func(t float64) string {
		if t <= 0 {
			return colors[0]
		}
		if t >= 1 {
			return colors[len(colors)-1]
		}
		// Map t to segment.
		n := len(pts) - 1
		x := t * float64(n)
		i := int(math.Floor(x))
		if i >= n {
			i = n - 1
		}
		s := x - float64(i)
		// B-spline basis around pts[i], pts[i+1] with phantom endpoints.
		p0 := pts[i]
		p1 := pts[i+1]
		var c0, c2 [3]float64
		if i > 0 {
			c0 = pts[i-1]
		} else {
			c0 = [3]float64{2*p0[0] - p1[0], 2*p0[1] - p1[1], 2*p0[2] - p1[2]}
		}
		if i+2 < len(pts) {
			c2 = pts[i+2]
		} else {
			c2 = [3]float64{2*p1[0] - p0[0], 2*p1[1] - p0[1], 2*p1[2] - p0[2]}
		}
		r, g, b := bsplinePoint(c0, p0, p1, c2, s)
		return rgbToHex(r, g, b)
	}
}

// bsplinePoint evaluates the cubic B-spline at t in [0,1] over four controls.
func bsplinePoint(p0, p1, p2, p3 [3]float64, t float64) (r, g, b float64) {
	t2 := t * t
	t3 := t2 * t
	blend := func(a, b, c, d float64) float64 {
		return ((a + 4*b + c) / 6) + ((-a+c)/6)*t + ((a-2*b+c)/2)*t2 + ((-a+3*b-3*c+d)/6)*t3
	}
	return blend(p0[0], p1[0], p2[0], p3[0]), blend(p0[1], p1[1], p2[1], p3[1]), blend(p0[2], p1[2], p2[2], p3[2])
}

// parseHexRGB parses "#rrggbb" into 0..1 RGB channels.
func parseHexRGB(s string) (r, g, b float64) {
	if len(s) >= 7 && s[0] == '#' {
		ri, _ := strconv.ParseInt(s[1:7], 16, 64)
		r = float64((ri>>16)&0xFF) / 255
		g = float64((ri>>8)&0xFF) / 255
		b = float64(ri&0xFF) / 255
		return
	}
	return 0, 0, 0
}

// rgbToHex formats 0..1 RGB channels back to "#rrggbb".
func rgbToHex(r, g, b float64) string {
	clamp := func(v float64) float64 {
		if v < 0 {
			return 0
		}
		if v > 1 {
			return 1
		}
		return v
	}
	ri := int(math.Round(clamp(r) * 255))
	gi := int(math.Round(clamp(g) * 255))
	bi := int(math.Round(clamp(b) * 255))
	return "#" + hex2(ri) + hex2(gi) + hex2(bi)
}

func hex2(v int) string {
	const hexchars = "0123456789abcdef"
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	return string([]byte{hexchars[v>>4], hexchars[v&0xF]})
}

// ApplyColorModifiers applies a chain of brighter/darker/opacity modifiers to
// a base color string. Mirrors nivo's color modifier pipeline.
func ApplyColorModifiers(color string, modifiers [][2]any) string {
	if len(modifiers) == 0 {
		return color
	}
	c := d3color.RGBColor(color)
	for _, m := range modifiers {
		if len(m) < 2 {
			continue
		}
		kind, ok := m[0].(string)
		if !ok {
			continue
		}
		amt, _ := toFloat(m[1])
		switch kind {
		case "brighter":
			c = c.Brighter(amt)
		case "darker":
			c = c.Darker(amt)
		case "opacity":
			c.SetOpacity(amt)
		}
	}
	return c.String()
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}
