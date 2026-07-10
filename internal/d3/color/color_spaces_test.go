package d3color

import (
	"math"
	"testing"
)

// Reference values below were produced by d3-color@3 / d3-interpolate@3:
//
//	c.hsl(c.rgb(s)), c.lab(c.rgb(s)), c.hcl(c.rgb(s))  // conversions
//	it.interpolateHsl/Lab/Hcl(a, b)(t)                 // interpolation
//	c.lab(s).brighter(k).formatHex()                   // modifiers
//
// They pin the HSL/Lab/Lch ports to d3's exact math (D50 white point for Lab).

// approxNaN compares two floats, treating NaN==NaN as equal.
func approxNaN(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	return math.Abs(a-b) <= tol
}

func TestRGBToHSL(t *testing.T) {
	cases := []struct {
		in      string
		h, s, l float64
	}{
		{"#ff0000", 0, 1, 0.5},
		{"#00ff00", 120, 1, 0.5},
		{"#0000ff", 240, 1, 0.5},
		{"#808080", math.NaN(), 0, 0.501961},
		{"#123456", 210, 0.653846, 0.203922},
		{"#ffffff", math.NaN(), math.NaN(), 1},
		{"#000000", math.NaN(), math.NaN(), 0},
		{"#a52a2a", 0, 0.594203, 0.405882},
		{"#ffa500", 38.8235, 1, 0.5},
	}
	for _, c := range cases {
		got := RGBColor(c.in).HSL()
		if !approxNaN(got.H, c.h, 5e-3) || !approxNaN(got.S, c.s, 5e-3) || !approxNaN(got.L, c.l, 5e-3) {
			t.Errorf("%s HSL = (%.5f,%.5f,%.5f) want (%.5f,%.5f,%.5f)",
				c.in, got.H, got.S, got.L, c.h, c.s, c.l)
		}
	}
}

func TestRGBToLab(t *testing.T) {
	cases := []struct {
		in      string
		l, a, b float64
	}{
		{"#ff0000", 54.291734, 80.812455, 69.88504},
		{"#00ff00", 87.818128, -79.287281, 80.990256},
		{"#0000ff", 29.567573, 68.298653, -112.02943},
		{"#808080", 53.585013, 0, 0},
		{"#123456", 20.675279, -2.276767, -24.592974},
		{"#ffffff", 100, 0, 0},
		{"#000000", 0, 0, 0},
		{"#a52a2a", 38.149667, 50.388769, 31.834059},
		{"#ffa500", 75.590394, 27.519295, 79.116221},
	}
	for _, c := range cases {
		got := RGBColor(c.in).Lab()
		if !approxNaN(got.L, c.l, 5e-3) || !approxNaN(got.A, c.a, 5e-3) || !approxNaN(got.B, c.b, 5e-3) {
			t.Errorf("%s Lab = (%.5f,%.5f,%.5f) want (%.5f,%.5f,%.5f)",
				c.in, got.L, got.A, got.B, c.l, c.a, c.b)
		}
	}
}

func TestRGBToLch(t *testing.T) {
	cases := []struct {
		in      string
		l, c, h float64
	}{
		{"#ff0000", 54.291734, 106.838999, 40.8526},
		{"#00ff00", 87.818128, 113.339731, 134.3912},
		{"#0000ff", 29.567573, 131.207085, 301.3685},
		{"#808080", 53.585013, 0, math.NaN()},
		{"#123456", 20.675279, 24.698139, 264.7108},
		{"#ffffff", 100, math.NaN(), math.NaN()},
		{"#000000", 0, math.NaN(), math.NaN()},
		{"#a52a2a", 38.149667, 59.60231, 32.2834},
		{"#ffa500", 75.590394, 83.765673, 70.8206},
	}
	for _, c := range cases {
		got := RGBColor(c.in).Lch()
		if !approxNaN(got.L, c.l, 5e-3) || !approxNaN(got.C, c.c, 5e-3) || !approxNaN(got.H, c.h, 5e-3) {
			t.Errorf("%s Lch = (%.5f,%.5f,%.5f) want (%.5f,%.5f,%.5f)",
				c.in, got.L, got.C, got.H, c.l, c.c, c.h)
		}
	}
}

// TestSpaceRoundtripHex confirms every space converts back to the original RGB
// hex (d3: hsl/lab/hcl(s).formatHex() == s for these inputs).
func TestSpaceRoundtripHex(t *testing.T) {
	cases := []string{"#ff0000", "#00ff00", "#0000ff", "#808080", "#123456", "#ffffff", "#000000", "#a52a2a", "#ffa500"}
	for _, in := range cases {
		rgb := RGBColor(in)
		if got := rgb.HSL().FormatHex(); got != in {
			t.Errorf("HSL roundtrip %s -> %s", in, got)
		}
		if got := rgb.Lab().FormatHex(); got != in {
			t.Errorf("Lab roundtrip %s -> %s", in, got)
		}
		if got := rgb.Lch().FormatHex(); got != in {
			t.Errorf("Lch roundtrip %s -> %s", in, got)
		}
	}
}

// TestInterpolateSpaces checks midpoint (and one quarter-point) interpolation
// against d3-interpolate's interpolateRgb/Hsl/Lab/Hcl, compared as final RGB
// hex (which is what downstream scales emit).
func TestInterpolateSpaces(t *testing.T) {
	mid := func(a, b string, sp Space) string {
		return RGBColor(a).In(sp).Interpolate(RGBColor(b), 0.5).RGB().FormatHex()
	}
	cases := []struct {
		a, b string
		sp   Space
		want string // d3 output as hex (rgb(...) values converted)
	}{
		// red -> blue at t=0.5
		{"#ff0000", "#0000ff", SpaceRGB, "#800080"}, // rgb(128,0,128)
		{"#ff0000", "#0000ff", SpaceHSL, "#ff00ff"}, // rgb(255,0,255)
		{"#ff0000", "#0000ff", SpaceLab, "#c10088"}, // rgb(193,0,136)
		{"#ff0000", "#0000ff", SpaceLch, "#f50086"}, // rgb(245,0,134)
		// red -> lime at t=0.5
		{"#ff0000", "#00ff00", SpaceRGB, "#808000"}, // rgb(128,128,0)
		{"#ff0000", "#00ff00", SpaceHSL, "#ffff00"}, // rgb(255,255,0)
		{"#ff0000", "#00ff00", SpaceLab, "#c8ac00"}, // rgb(200,172,0)
		{"#ff0000", "#00ff00", SpaceLch, "#d1a900"}, // rgb(209,169,0)
	}
	for _, c := range cases {
		if got := mid(c.a, c.b, c.sp); got != c.want {
			t.Errorf("%s->%s in %s @0.5 = %s want %s", c.a, c.b, c.sp, got, c.want)
		}
	}

	// Quarter-point green(#008000) -> blue in each space.
	quarter := func(sp Space) string {
		return RGBColor("#008000").In(sp).Interpolate(RGBColor("#0000ff"), 0.25).RGB().FormatHex()
	}
	q := []struct {
		sp   Space
		want string
	}{
		{SpaceLab, "#456c55"}, // rgb(69,108,85)
		{SpaceLch, "#007e57"}, // rgb(0,126,87)
		{SpaceHSL, "#00a050"}, // rgb(0,160,80)
	}
	for _, c := range q {
		if got := quarter(c.sp); got != c.want {
			t.Errorf("green->blue in %s @0.25 = %s want %s", c.sp, got, c.want)
		}
	}
}

// TestSpaceModifiers pins the per-space brighter/darker to d3 (Lab/Hcl add
// Kn=18 to L; HSL scales L by 0.7^∓k).
func TestSpaceModifiers(t *testing.T) {
	if got := RGBColor("#808080").Lab().Brighter(1).FormatHex(); got != "#afafaf" {
		t.Errorf("Lab brighter(1) #808080 = %s want #afafaf", got)
	}
	if got := RGBColor("#808080").Lab().Darker(1).FormatHex(); got != "#545454" {
		t.Errorf("Lab darker(1) #808080 = %s want #545454", got)
	}
	if got := RGBColor("#808080").Lch().Brighter(1).FormatHex(); got != "#afafaf" {
		t.Errorf("Lch brighter(1) #808080 = %s want #afafaf", got)
	}
	if got := RGBColor("#808080").HSL().Brighter(1).FormatHex(); got != "#b7b7b7" {
		t.Errorf("HSL brighter(1) #808080 = %s want #b7b7b7", got)
	}
}

// TestInterfaceSatisfied is a runtime companion to the compile-time asserts in
// spaces.go: every space is usable as a Color, and In(SpaceRGB) is identity.
func TestInterfaceSatisfied(t *testing.T) {
	cs := []Color{
		RGBColor("#123456"),
		RGBColor("#123456").HSL(),
		RGBColor("#123456").Lab(),
		RGBColor("#123456").Lch(),
	}
	for _, c := range cs {
		if c.RGB().FormatHex() != "#123456" {
			t.Errorf("%T.RGB() = %s want #123456", c, c.RGB().FormatHex())
		}
	}
	if _, ok := RGBColor("#123456").In(SpaceRGB).(*RGB); !ok {
		t.Errorf("In(SpaceRGB) should return *RGB unchanged")
	}
}
