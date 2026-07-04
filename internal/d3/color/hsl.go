package d3color

import "math"

// HSL is the d3-color hsl color space. H is the hue in degrees [0,360); S and L
// (saturation, lightness) and Opac are in [0,1]. A NaN hue (achromatic) or NaN
// saturation is preserved, matching d3-color's handling of gray colors.
//
// Channel math mirrors d3-color exactly; d3 stores rgb channels in 0..255 while
// this package uses 0..1, but the hsl<->rgb conversions are ratio-based so the
// two conventions agree bit-for-bit after the 0..1 scaling.
type HSL struct {
	H, S, L, Opac float64
}

// NewHSL constructs an HSL. H is wrapped/clamped only on conversion; S, L, and
// opacity are clamped to [0,1] here (NaN is preserved for H and S).
func NewHSL(h, s, l, opacity float64) *HSL {
	return &HSL{H: h, S: clampNaNUnit(s), L: clampNaNUnit(l), Opac: clampUnit(opacity)}
}

// HSLColor parses a CSS color string and converts it to HSL (via RGB), matching
// d3-color's hsl() factory on non-hsl inputs.
func HSLColor(s string) *HSL { return rgb(s).HSL() }

// HSL converts an RGB to HSL. Port of d3-color's rgbConvert -> hslConvert.
func (c *RGB) HSL() *HSL {
	r, g, b := c.R, c.G, c.B
	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	h := math.NaN()
	s := max - min
	l := (max + min) / 2
	if s != 0 {
		switch {
		case r == max:
			h = (g - b) / s
			if g < b {
				h += 6
			}
		case g == max:
			h = (b-r)/s + 2
		default:
			h = (r-g)/s + 4
		}
		h *= 60
		if l < 0.5 {
			s /= max + min
		} else {
			s /= 2 - max - min
		}
	} else {
		if l > 0 && l < 1 {
			s = 0
		} else {
			s = math.NaN()
		}
	}
	return &HSL{H: h, S: s, L: l, Opac: c.Opac}
}

// RGB converts the HSL back to RGB. Port of d3-color's Hsl.rgb / hsl2rgb.
func (c *HSL) RGB() *RGB {
	h := c.H
	if math.IsNaN(h) {
		h = 0
	}
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	s := c.S
	if math.IsNaN(c.H) || math.IsNaN(s) {
		s = 0
	}
	l := c.L
	var m2 float64
	if l < 0.5 {
		m2 = l * (1 + s)
	} else {
		m2 = l + s - l*s
	}
	m1 := 2*l - m2
	var h1, h3 float64
	if h >= 240 {
		h1 = h - 240
	} else {
		h1 = h + 120
	}
	if h < 120 {
		h3 = h + 240
	} else {
		h3 = h - 120
	}
	return &RGB{
		R:    hsl2rgb(h1, m1, m2),
		G:    hsl2rgb(h, m1, m2),
		B:    hsl2rgb(h3, m1, m2),
		Opac: c.Opac,
	}
}

// hsl2rgb resolves a single channel; h in [0,360). Matches d3-color's hsl2rgb
// (without the *255 factor, since channels are 0..1 here).
func hsl2rgb(h, m1, m2 float64) float64 {
	switch {
	case h < 60:
		return m1 + (m2-m1)*h/60
	case h < 180:
		return m2
	case h < 240:
		return m1 + (m2-m1)*(240-h)/60
	default:
		return m1
	}
}

// Brighter lightens in HSL space by scaling L by 0.7^-k. d3-color's
// Hsl.brighter. k==0 clones.
func (c *HSL) Brighter(k float64) Color {
	if k == 0 {
		return c.Copy()
	}
	return &HSL{H: c.H, S: c.S, L: c.L * math.Pow(0.7, -k), Opac: c.Opac}
}

// Darker darkens in HSL space by scaling L by 0.7^k. d3-color's Hsl.darker.
func (c *HSL) Darker(k float64) Color {
	if k == 0 {
		return c.Copy()
	}
	return &HSL{H: c.H, S: c.S, L: c.L * math.Pow(0.7, k), Opac: c.Opac}
}

// Interpolate blends from c to other in HSL space at t in [0,1]. Hue takes the
// shortest angular path (d3's interpolateHsl); S, L, and opacity are linear.
func (c *HSL) Interpolate(other Color, t float64) Color {
	o := other.RGB().HSL()
	return &HSL{
		H:    hueLerp(c.H, o.H, t),
		S:    lerp(c.S, o.S, t),
		L:    lerp(c.L, o.L, t),
		Opac: lerp(c.Opac, o.Opac, t),
	}
}

func (c *HSL) Displayable() bool  { return c.RGB().Displayable() }
func (c *HSL) FormatRGB() string  { return c.RGB().FormatRGB() }
func (c *HSL) FormatRGBA() string { return c.RGB().FormatRGBA() }
func (c *HSL) FormatHex() string  { return c.RGB().FormatHex() }
func (c *HSL) FormatHex8() string { return c.RGB().FormatHex8() }

// String formats via the RGB conversion (hex when opaque, rgba otherwise),
// consistent with the rest of the package's downstream CSS output.
func (c *HSL) String() string { return c.RGB().String() }

func (c *HSL) Opacity() float64     { return c.Opac }
func (c *HSL) SetOpacity(o float64) { c.Opac = clampUnit(o) }
func (c *HSL) Copy() Color          { return &HSL{H: c.H, S: c.S, L: c.L, Opac: c.Opac} }

// clampNaNUnit clamps to [0,1] but preserves NaN (used for hue/saturation).
func clampNaNUnit(v float64) float64 {
	if math.IsNaN(v) {
		return v
	}
	return clampFloat(v, 0, 1)
}

// hueLerp interpolates a hue angle taking the shortest path around the circle.
// Port of d3-interpolate's hue(). NaN endpoints hold the other endpoint.
func hueLerp(a, b, t float64) float64 {
	if math.IsNaN(a) {
		return b
	}
	if math.IsNaN(b) {
		return a
	}
	d := b - a
	if d > 180 || d < -180 {
		d -= 360 * math.Round(d/360)
	}
	return a + d*t
}
