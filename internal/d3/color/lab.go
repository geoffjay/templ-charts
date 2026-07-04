package d3color

import "math"

// CIELAB / CIELCh constants, matching d3-color's lab.js. d3 uses the D50 white
// point (Xn, Zn) rather than D65.
const (
	labKn = 18.0
	labXn = 0.96422
	labYn = 1.0
	labZn = 0.82521
	labT0 = 4.0 / 29.0
	labT1 = 6.0 / 29.0
	labT2 = 3 * labT1 * labT1
	labT3 = labT1 * labT1 * labT1
)

// Lab is the d3-color CIELAB color space. L is lightness (0..100), a and b are
// the green–red and blue–yellow axes (unbounded, typically ~[-128,127]).
type Lab struct {
	L, A, B, Opac float64
}

// NewLab constructs a Lab with the given channels and opacity (clamped 0..1).
func NewLab(l, a, b, opacity float64) *Lab {
	return &Lab{L: l, A: a, B: b, Opac: clampUnit(opacity)}
}

// LabColor parses a CSS color string and converts it to Lab (via RGB), matching
// d3-color's lab() factory on non-lab inputs.
func LabColor(s string) *Lab { return rgb(s).Lab() }

// Lab converts an RGB to CIELAB. Port of d3-color's rgb -> lab path.
func (c *RGB) Lab() *Lab {
	r := rgb2lrgb(c.R)
	g := rgb2lrgb(c.G)
	b := rgb2lrgb(c.B)
	y := xyz2lab((0.2225045*r + 0.7168786*g + 0.0606169*b) / labYn)
	var x, z float64
	if r == g && g == b {
		x, z = y, y
	} else {
		x = xyz2lab((0.4360747*r + 0.3850649*g + 0.1430804*b) / labXn)
		z = xyz2lab((0.0139322*r + 0.0971045*g + 0.7141733*b) / labZn)
	}
	return &Lab{L: 116*y - 16, A: 500 * (x - y), B: 200 * (y - z), Opac: c.Opac}
}

// RGB converts the Lab back to RGB. Port of d3-color's Lab.rgb. Channels may
// fall outside [0,1] for out-of-gamut colors (as in d3); toByte clamps on
// format so hex output matches d3 exactly.
func (c *Lab) RGB() *RGB {
	y := (c.L + 16) / 116
	x := y
	if !math.IsNaN(c.A) {
		x = y + c.A/500
	}
	z := y
	if !math.IsNaN(c.B) {
		z = y - c.B/200
	}
	x = labXn * lab2xyz(x)
	y = labYn * lab2xyz(y)
	z = labZn * lab2xyz(z)
	return &RGB{
		R:    lrgb2rgb(3.1338561*x - 1.6168667*y - 0.4906146*z),
		G:    lrgb2rgb(-0.9787684*x + 1.9161415*y + 0.0334540*z),
		B:    lrgb2rgb(0.0719453*x - 0.2289914*y + 1.4052427*z),
		Opac: c.Opac,
	}
}

// rgb2lrgb linearizes a gamma-encoded sRGB channel (0..1). Port of d3's
// rgb2lrgb without the /255 (channels are already 0..1 here).
func rgb2lrgb(x float64) float64 {
	if x <= 0.04045 {
		return x / 12.92
	}
	return math.Pow((x+0.055)/1.055, 2.4)
}

// lrgb2rgb gamma-encodes a linear-light channel back to sRGB (0..1), without
// the *255 factor.
func lrgb2rgb(x float64) float64 {
	if x <= 0.0031308 {
		return 12.92 * x
	}
	return 1.055*math.Pow(x, 1/2.4) - 0.055
}

func xyz2lab(t float64) float64 {
	if t > labT3 {
		return math.Cbrt(t)
	}
	return t/labT2 + labT0
}

func lab2xyz(t float64) float64 {
	if t > labT1 {
		return t * t * t
	}
	return labT2 * (t - labT0)
}

// Brighter lightens in Lab space by adding Kn*k to L. d3-color's Lab.brighter.
func (c *Lab) Brighter(k float64) Color {
	return &Lab{L: c.L + labKn*k, A: c.A, B: c.B, Opac: c.Opac}
}

// Darker darkens in Lab space by subtracting Kn*k from L. d3-color's Lab.darker.
func (c *Lab) Darker(k float64) Color {
	return &Lab{L: c.L - labKn*k, A: c.A, B: c.B, Opac: c.Opac}
}

// Interpolate blends from c to other in Lab space at t in [0,1] (linear L, a,
// b, and opacity). Matches d3's interpolateLab.
func (c *Lab) Interpolate(other Color, t float64) Color {
	o := other.RGB().Lab()
	return &Lab{
		L:    lerp(c.L, o.L, t),
		A:    lerp(c.A, o.A, t),
		B:    lerp(c.B, o.B, t),
		Opac: lerp(c.Opac, o.Opac, t),
	}
}

func (c *Lab) Displayable() bool    { return c.RGB().Displayable() }
func (c *Lab) FormatRGB() string    { return c.RGB().FormatRGB() }
func (c *Lab) FormatRGBA() string   { return c.RGB().FormatRGBA() }
func (c *Lab) FormatHex() string    { return c.RGB().FormatHex() }
func (c *Lab) FormatHex8() string   { return c.RGB().FormatHex8() }
func (c *Lab) String() string       { return c.RGB().String() }
func (c *Lab) Opacity() float64     { return c.Opac }
func (c *Lab) SetOpacity(o float64) { c.Opac = clampUnit(o) }
func (c *Lab) Copy() Color          { return &Lab{L: c.L, A: c.A, B: c.B, Opac: c.Opac} }
