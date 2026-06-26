// Package d3color provides a minimal port of d3-color used by templ-charts.
//
// Only the operations nivo exercises on inherited-color configs are supported:
// parsing a CSS color string into RGB, applying brighter/darker modifiers,
// overriding opacity, and formatting back to a CSS color string.
package d3color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color is the common interface implemented by RGB (and, in future, HSL).
// Mirrors the subset of d3-color's Color that nivo touches: brighter/darker
// modifiers, opacity mutation, and string formatting.
type Color interface {
	Brighter(k float64) Color
	Darker(k float64) Color
	Displayable() bool
	FormatRGB() string
	FormatRGBA() string
	FormatHex() string
	FormatHex8() string
	String() string
	Opacity() float64
	SetOpacity(o float64)
	// RGB returns the color as an *RGB; for an *RGB it returns itself.
	RGB() *RGB
	// Copy returns a deep clone of the color.
	Copy() Color
}

// RGB is the d3-color rgb color space. Channels are in [0,1] internally,
// matching d3-color's convention (the public FormatHex / FormatRGB convert
// to 0..255 on the way out).
type RGB struct {
	R, G, B, Opac float64
}

// NewRGB constructs an RGB with the given channel values (0..1 range) and
// opacity (0..1).
func NewRGB(r, g, b, opacity float64) *RGB {
	return &RGB{R: clampChannel(r), G: clampChannel(g), B: clampChannel(b), Opac: clampUnit(opacity)}
}

// rgb parses a CSS color string into an *RGB. Supported forms mirror the
// subset d3-color handles for typical chart usage:
//   - "#rgb", "#rrggbb", "#rgba", "#rrggbbaa"
//   - "rgb(r, g, b)" / "rgba(r, g, b, a)" (r,g,b in 0..255, a in 0..1)
//   - "transparent" (treated as rgb(0,0,0,0))
//   - named colors: black, white, red, green, blue, yellow, cyan, magenta,
//     gray/grey, orange, purple, brown, pink, navy, teal, silver, lime,
//     maroon, olive, aqua, fuchsia
//
// Unrecognized inputs fall back to opaque black (d3-color's behavior for
// unparseable strings is to return an RGBColor with NaN channels; we treat
// them as black instead since NaN channels produce no useful SVG output).
func rgb(s string) *RGB {
	c, ok := parseColor(s)
	if !ok {
		return NewRGB(0, 0, 0, 1)
	}
	return c
}

// RGBColor is the public entry point matching d3-color's `rgb()` factory.
// It returns an *RGB.
func RGBColor(s string) *RGB {
	return rgb(s)
}

// Brighter multiplies channel values by 0.7^-k (k=1 ≈ 1.43x brighter).
// Channels are clamped to [0,1] afterward. d3-color preserves opacity.
func (c *RGB) Brighter(k float64) *RGB {
	if k == 0 {
		return c.Copy().RGB()
	}
	k = math.Pow(0.7, -k)
	return NewRGB(c.R*k, c.G*k, c.B*k, c.Opac)
}

// Darker multiplies channel values by 0.7^k (k=1 ≈ 0.7x, i.e. ~30% darker).
func (c *RGB) Darker(k float64) *RGB {
	if k == 0 {
		return c.Copy().RGB()
	}
	k = math.Pow(0.7, k)
	return NewRGB(c.R*k, c.G*k, c.B*k, c.Opac)
}

// Displayable reports whether the RGB is displayable in CSS (true for all
// finite-channel RGBs; NaN is the only non-displayable case in d3-color).
func (c *RGB) Displayable() bool {
	return !math.IsNaN(c.R) && !math.IsNaN(c.G) && !math.IsNaN(c.B) && !math.IsNaN(c.Opac)
}

// FormatRGB returns "rgb(r, g, b)" with channels in 0..255.
func (c *RGB) FormatRGB() string {
	return fmt.Sprintf("rgb(%d, %d, %d)", toByte(c.R), toByte(c.G), toByte(c.B))
}

// FormatRGBA returns "rgba(r, g, b, a)" with channels in 0..255 and a in 0..1.
func (c *RGB) FormatRGBA() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %s)", toByte(c.R), toByte(c.G), toByte(c.B), formatFloat(c.Opac))
}

// FormatHex returns "#rrggbb".
func (c *RGB) FormatHex() string {
	return fmt.Sprintf("#%02x%02x%02x", toByte(c.R), toByte(c.G), toByte(c.B))
}

// FormatHex8 returns "#rrggbbaa" (8-digit hex with alpha).
func (c *RGB) FormatHex8() string {
	return fmt.Sprintf("#%02x%02x%02x%02x", toByte(c.R), toByte(c.G), toByte(c.B), toByte(c.Opac))
}

// String returns "rgba(...)" when opacity is < 1, otherwise the hex form —
// matching d3-color's toString behavior.
func (c *RGB) String() string {
	if math.IsNaN(c.Opac) || c.Opac >= 1 {
		return c.FormatHex()
	}
	return c.FormatRGBA()
}

// Opacity returns the alpha channel.
func (c *RGB) Opacity() float64 { return c.Opac }

// SetOpacity mutates the alpha channel in place. Matches d3-color's mutable
// `.opacity =` assignment used by nivo's "opacity" modifier.
func (c *RGB) SetOpacity(o float64) { c.Opac = clampUnit(o) }

// RGB satisfies the Color interface (returns itself).
func (c *RGB) RGB() *RGB { return c }

// Copy returns a deep clone.
func (c *RGB) Copy() *RGB { return &RGB{R: c.R, G: c.G, B: c.B, Opac: c.Opac} }

// --- parsing ---------------------------------------------------------------

// parseColor is the actual parser; returns (color, ok). `ok` is false for
// unrecognized forms.
func parseColor(s string) (*RGB, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil, false
	}
	if s == "transparent" {
		return NewRGB(0, 0, 0, 0), true
	}
	if c, ok := namedColor(s); ok {
		return c, true
	}
	if strings.HasPrefix(s, "#") {
		return parseHex(s)
	}
	if strings.HasPrefix(s, "rgb(") || strings.HasPrefix(s, "rgba(") {
		return parseRGBFunc(s)
	}
	return nil, false
}

// parseHex handles #rgb, #rgba, #rrggbb, #rrggbbaa.
func parseHex(s string) (*RGB, bool) {
	h := s[1:]
	switch len(h) {
	case 3: // #rgb
		r, g, b, ok := expandShort(h)
		if !ok {
			return nil, false
		}
		return NewRGB(float64(r)/255, float64(g)/255, float64(b)/255, 1), true
	case 4: // #rgba
		r, g, b, a, ok := expandShortAlpha(h)
		if !ok {
			return nil, false
		}
		return NewRGB(float64(r)/255, float64(g)/255, float64(b)/255, float64(a)/255), true
	case 6: // #rrggbb
		r, g, b, ok := splitHex6(h)
		if !ok {
			return nil, false
		}
		return NewRGB(float64(r)/255, float64(g)/255, float64(b)/255, 1), true
	case 8: // #rrggbbaa
		r, g, b, a, ok := splitHex8(h)
		if !ok {
			return nil, false
		}
		return NewRGB(float64(r)/255, float64(g)/255, float64(b)/255, float64(a)/255), true
	}
	return nil, false
}

func expandShort(h string) (r, g, b int, ok bool) {
	if len(h) != 3 {
		return 0, 0, 0, false
	}
	r, ok = hexDigit(h[0])
	if !ok {
		return 0, 0, 0, false
	}
	g, ok = hexDigit(h[1])
	if !ok {
		return 0, 0, 0, false
	}
	b, ok = hexDigit(h[2])
	if !ok {
		return 0, 0, 0, false
	}
	// expand by duplicating: 0xa -> 0xaa
	return r * 17, g * 17, b * 17, true
}

func expandShortAlpha(h string) (r, g, b, a int, ok bool) {
	if len(h) != 4 {
		return 0, 0, 0, 0, false
	}
	r, ok = hexDigit(h[0])
	if !ok {
		return 0, 0, 0, 0, false
	}
	g, ok = hexDigit(h[1])
	if !ok {
		return 0, 0, 0, 0, false
	}
	b, ok = hexDigit(h[2])
	if !ok {
		return 0, 0, 0, 0, false
	}
	a, ok = hexDigit(h[3])
	if !ok {
		return 0, 0, 0, 0, false
	}
	return r * 17, g * 17, b * 17, a * 17, true
}

func splitHex6(h string) (r, g, b int, ok bool) {
	rr, err := strconv.ParseInt(h[0:2], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	gg, err := strconv.ParseInt(h[2:4], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	bb, err := strconv.ParseInt(h[4:6], 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(rr), int(gg), int(bb), true
}

func splitHex8(h string) (r, g, b, a int, ok bool) {
	r, g, b, ok = splitHex6(h[:6])
	if !ok {
		return 0, 0, 0, 0, false
	}
	aa, err := strconv.ParseInt(h[6:8], 16, 32)
	if err != nil {
		return 0, 0, 0, 0, false
	}
	return r, g, b, int(aa), true
}

func hexDigit(b byte) (int, bool) {
	switch {
	case b >= '0' && b <= '9':
		return int(b - '0'), true
	case b >= 'a' && b <= 'f':
		return int(b - 'a' + 10), true
	}
	return 0, false
}

// parseRGBFunc handles "rgb(r,g,b)" / "rgba(r,g,b,a)" with r,g,b in 0..255
// (optionally as percentages) and a in 0..1.
func parseRGBFunc(s string) (*RGB, bool) {
	// strip "rgb(" / "rgba(" and trailing ")"
	prefix := "rgba("
	if strings.HasPrefix(s, "rgb(") {
		prefix = "rgb("
	}
	if !strings.HasSuffix(s, ")") {
		return nil, false
	}
	inner := s[len(prefix) : len(s)-1]
	parts := strings.Split(inner, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	switch len(parts) {
	case 3:
		r, ok := parseByteOrPct(parts[0], 255)
		if !ok {
			return nil, false
		}
		g, ok := parseByteOrPct(parts[1], 255)
		if !ok {
			return nil, false
		}
		b, ok := parseByteOrPct(parts[2], 255)
		if !ok {
			return nil, false
		}
		return NewRGB(r/255, g/255, b/255, 1), true
	case 4:
		r, ok := parseByteOrPct(parts[0], 255)
		if !ok {
			return nil, false
		}
		g, ok := parseByteOrPct(parts[1], 255)
		if !ok {
			return nil, false
		}
		b, ok := parseByteOrPct(parts[2], 255)
		if !ok {
			return nil, false
		}
		a, ok := parseAlpha(parts[3])
		if !ok {
			return nil, false
		}
		return NewRGB(r/255, g/255, b/255, a), true
	}
	return nil, false
}

// parseByteOrPct parses "0..255" or "0%..100%" into a 0..255 channel value.
// `max` is 255 for normal channels (255 if percent multiplier needed).
func parseByteOrPct(s string, max float64) (float64, bool) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, false
		}
		return clampFloat(v*max/100, 0, max), true
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return clampFloat(v, 0, max), true
}

func parseAlpha(s string) (float64, bool) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, false
		}
		return clampUnit(v / 100), true
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return clampUnit(v), true
}

// --- helpers ---------------------------------------------------------------

func clampChannel(v float64) float64 { return clampFloat(v, 0, 1) }
func clampUnit(v float64) float64    { return clampFloat(v, 0, 1) }

func clampFloat(v, lo, hi float64) float64 {
	if math.IsNaN(v) {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func toByte(v float64) int {
	if math.IsNaN(v) {
		return 0
	}
	x := int(math.Round(v * 255))
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return x
}

func formatFloat(v float64) string {
	if v == 1 {
		return "1"
	}
	// strip trailing zeros like d3-color does
	s := strconv.FormatFloat(v, 'f', -1, 64)
	return s
}

// --- named colors ----------------------------------------------------------
//
// Only the most common CSS named colors used in chart defaults are supported.
// `transparent` is handled separately.

func namedColor(name string) (*RGB, bool) {
	switch name {
	case "black":
		return NewRGB(0, 0, 0, 1), true
	case "white":
		return NewRGB(1, 1, 1, 1), true
	case "red":
		return NewRGB(1, 0, 0, 1), true
	case "green":
		return NewRGB(0, 0.5019607843137255, 0, 1), true // CSS green = #008000
	case "blue":
		return NewRGB(0, 0, 1, 1), true
	case "yellow":
		return NewRGB(1, 1, 0, 1), true
	case "cyan", "aqua":
		return NewRGB(0, 1, 1, 1), true
	case "magenta", "fuchsia":
		return NewRGB(1, 0, 1, 1), true
	case "gray", "grey":
		return NewRGB(0.5019607843137255, 0.5019607843137255, 0.5019607843137255, 1), true // #808080
	case "silver":
		return NewRGB(0.7529411764705882, 0.7529411764705882, 0.7529411764705882, 1), true // #c0c0c0
	case "maroon":
		return NewRGB(0.5019607843137255, 0, 0, 1), true // #800000
	case "olive":
		return NewRGB(0.5019607843137255, 0.5019607843137255, 0, 1), true // #808000
	case "lime":
		return NewRGB(0, 1, 0, 1), true
	case "purple":
		return NewRGB(0.5019607843137255, 0, 0.5019607843137255, 1), true // #800080
	case "teal":
		return NewRGB(0, 0.5019607843137255, 0.5019607843137255, 1), true // #008080
	case "navy":
		return NewRGB(0, 0, 0.5019607843137255, 1), true // #000080
	case "orange":
		return NewRGB(1, 0.6470588235294118, 0, 1), true // #ffa500
	case "brown":
		return NewRGB(0.6470588235294118, 0.16470588235294117, 0.16470588235294117, 1), true // #a52a2a
	case "pink":
		return NewRGB(1, 0.7529411764705882, 0.796078431372549, 1), true // #ffc0cb
	}
	return nil, false
}

// ApplyModifiers runs a chain of [type, amount] modifiers on the color,
// returning the resulting formatted string. Mirrors the reduce loop in
// nivo's inheritedColor.ts.
//
// Modifier types: "brighter", "darker", "opacity".
func ApplyModifiers(color string, modifiers [][2]any) (string, error) {
	c := RGBColor(color)
	for _, m := range modifiers {
		if len(m) != 2 {
			return "", fmt.Errorf("d3color: invalid modifier %#v", m)
		}
		kind, _ := m[0].(string)
		amount, ok := toFloat(m[1])
		if !ok {
			return "", fmt.Errorf("d3color: modifier amount is not a number: %#v", m[1])
		}
		switch kind {
		case "brighter":
			c = c.Brighter(amount)
		case "darker":
			c = c.Darker(amount)
		case "opacity":
			c.SetOpacity(amount)
		default:
			return "", fmt.Errorf("d3color: invalid color modifier %q (want brighter|darker|opacity)", kind)
		}
	}
	return c.String(), nil
}

// toFloat accepts numeric types and strings holding a number.
func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
