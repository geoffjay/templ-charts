package d3color

import "math"

const (
	degrees = 180 / math.Pi
	radians = math.Pi / 180
)

// Lch is the d3-color CIELCh(ab) color space: CIELAB expressed in cylindrical
// coordinates. L is lightness (0..100), C is chroma (>=0), H is hue in degrees
// [0,360). This is the same space d3 exposes as hcl (with reordered channels).
type Lch struct {
	L, C, H, Opac float64
}

// NewLch constructs an Lch with the given channels and opacity (clamped 0..1).
func NewLch(l, c, h, opacity float64) *Lch {
	return &Lch{L: l, C: c, H: h, Opac: clampUnit(opacity)}
}

// LchColor parses a CSS color string and converts it to Lch (via RGB/Lab),
// matching d3-color's lch()/hcl() factory on non-lch inputs.
func LchColor(s string) *Lch { return rgb(s).Lch() }

// Lch converts a Lab to Lch. Port of d3-color's hclConvert. For achromatic
// colors (a==b==0) the hue is undefined (NaN); chroma is 0 for mid lightnesses
// and NaN at the extremes (pure black/white), matching d3.
func (c *Lab) Lch() *Lch {
	if c.A == 0 && c.B == 0 {
		chroma := math.NaN()
		if c.L > 0 && c.L < 100 {
			chroma = 0
		}
		return &Lch{L: c.L, C: chroma, H: math.NaN(), Opac: c.Opac}
	}
	chroma := math.Sqrt(c.A*c.A + c.B*c.B)
	h := math.Atan2(c.B, c.A) * degrees
	if h < 0 {
		h += 360
	}
	return &Lch{L: c.L, C: chroma, H: h, Opac: c.Opac}
}

// Lch converts an RGB to Lch (via Lab).
func (c *RGB) Lch() *Lch { return c.Lab().Lch() }

// Lab converts the Lch back to Lab. Port of d3-color's hcl2lab.
func (c *Lch) Lab() *Lab {
	if math.IsNaN(c.H) {
		return &Lab{L: c.L, A: 0, B: 0, Opac: c.Opac}
	}
	h := c.H * radians
	return &Lab{L: c.L, A: math.Cos(h) * c.C, B: math.Sin(h) * c.C, Opac: c.Opac}
}

// RGB converts the Lch back to RGB (via Lab).
func (c *Lch) RGB() *RGB { return c.Lab().RGB() }

// Brighter lightens in Lch space by adding Kn*k to L. d3-color's Hcl.brighter.
func (c *Lch) Brighter(k float64) Color {
	return &Lch{L: c.L + labKn*k, C: c.C, H: c.H, Opac: c.Opac}
}

// Darker darkens in Lch space by subtracting Kn*k from L. d3-color's Hcl.darker.
func (c *Lch) Darker(k float64) Color {
	return &Lch{L: c.L - labKn*k, C: c.C, H: c.H, Opac: c.Opac}
}

// Interpolate blends from c to other in Lch space at t in [0,1]. Hue takes the
// shortest angular path (d3's interpolateHcl); L, C, and opacity are linear.
func (c *Lch) Interpolate(other Color, t float64) Color {
	o := other.RGB().Lch()
	return &Lch{
		L:    lerp(c.L, o.L, t),
		C:    lerp(c.C, o.C, t),
		H:    hueLerp(c.H, o.H, t),
		Opac: lerp(c.Opac, o.Opac, t),
	}
}

func (c *Lch) Displayable() bool    { return c.RGB().Displayable() }
func (c *Lch) FormatRGB() string    { return c.RGB().FormatRGB() }
func (c *Lch) FormatRGBA() string   { return c.RGB().FormatRGBA() }
func (c *Lch) FormatHex() string    { return c.RGB().FormatHex() }
func (c *Lch) FormatHex8() string   { return c.RGB().FormatHex8() }
func (c *Lch) String() string       { return c.RGB().String() }
func (c *Lch) Opacity() float64     { return c.Opac }
func (c *Lch) SetOpacity(o float64) { c.Opac = clampUnit(o) }
func (c *Lch) Copy() Color          { return &Lch{L: c.L, C: c.C, H: c.H, Opac: c.Opac} }
