package d3color

// Compile-time checks that every color space implements Color.
var (
	_ Color = (*RGB)(nil)
	_ Color = (*HSL)(nil)
	_ Color = (*Lab)(nil)
	_ Color = (*Lch)(nil)
)

// Space identifies a color space for interpolation and conversion. SpaceRGB is
// the zero value so a zero-initialized selector preserves today's RGB behavior.
type Space int

const (
	// SpaceRGB interpolates linearly in gamma-encoded sRGB (the default, and
	// the behavior every existing scale/palette golden was captured under).
	SpaceRGB Space = iota
	// SpaceHSL interpolates in HSL (hue takes the shortest angular path).
	SpaceHSL
	// SpaceLab interpolates in perceptually-uniform CIELAB.
	SpaceLab
	// SpaceLch interpolates in CIELCh (perceptual, hue shortest-path).
	SpaceLch
)

// String returns the lowercase space name (rgb/hsl/lab/lch).
func (s Space) String() string {
	switch s {
	case SpaceHSL:
		return "hsl"
	case SpaceLab:
		return "lab"
	case SpaceLch:
		return "lch"
	default:
		return "rgb"
	}
}

// In converts an RGB into the given space, returning a Color whose Interpolate
// blends in that space. Used by callers that want to opt into perceptual
// interpolation while keeping RGB as the default.
func (c *RGB) In(space Space) Color {
	switch space {
	case SpaceHSL:
		return c.HSL()
	case SpaceLab:
		return c.Lab()
	case SpaceLch:
		return c.Lch()
	default:
		return c
	}
}
