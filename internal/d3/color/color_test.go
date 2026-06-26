package d3color

import (
	"math"
	"testing"
)

func TestParseHex(t *testing.T) {
	cases := []struct {
		in   string
		r    int
		g    int
		b    int
		opac float64
		hex  string
	}{
		{"#000000", 0, 0, 0, 1, "#000000"},
		{"#fff", 255, 255, 255, 1, "#ffffff"},
		{"#ff0000", 255, 0, 0, 1, "#ff0000"},
		{"#ff000080", 255, 0, 0, 0.5019607843137255, "rgba(255, 0, 0, 0.5019607843137255)"},
		{"#f0a", 255, 0, 170, 1, "#ff00aa"},
		{"#f0a8", 255, 0, 170, 0.5333333333333333, "rgba(255, 0, 170, 0.5333333333333333)"},
	}
	for _, c := range cases {
		got := RGBColor(c.in)
		if toByte(got.R) != c.r || toByte(got.G) != c.g || toByte(got.B) != c.b {
			t.Errorf("%s: got channels %d,%d,%d want %d,%d,%d",
				c.in, toByte(got.R), toByte(got.G), toByte(got.B), c.r, c.g, c.b)
		}
		if math.Abs(got.Opac-c.opac) > 1e-9 {
			t.Errorf("%s: got opacity %v want %v", c.in, got.Opac, c.opac)
		}
		if got.String() != c.hex {
			t.Errorf("%s: String() = %q want %q", c.in, got.String(), c.hex)
		}
	}
}

func TestNamedColors(t *testing.T) {
	c := RGBColor("red")
	if c.String() != "#ff0000" {
		t.Errorf("red = %q want #ff0000", c.String())
	}
	c = RGBColor("transparent")
	if c.Opac != 0 {
		t.Errorf("transparent opacity = %v want 0", c.Opac)
	}
}

func TestRGBFunc(t *testing.T) {
	c := RGBColor("rgb(100, 200, 50)")
	if toByte(c.R) != 100 || toByte(c.G) != 200 || toByte(c.B) != 50 {
		t.Errorf("rgb(100,200,50) channels = %d,%d,%d want 100,200,50",
			toByte(c.R), toByte(c.G), toByte(c.B))
	}
	c = RGBColor("rgba(0, 0, 0, 0.5)")
	if c.Opac != 0.5 {
		t.Errorf("rgba opacity = %v want 0.5", c.Opac)
	}
}

func TestBrighterDarker(t *testing.T) {
	c := RGBColor("#808080")
	bright := c.Brighter(1)
	if !(toByte(bright.R) > toByte(c.R)) {
		t.Errorf("Brighter(1) should be brighter: %d vs %d", toByte(bright.R), toByte(c.R))
	}
	dark := c.Darker(1)
	if !(toByte(dark.R) < toByte(c.R)) {
		t.Errorf("Darker(1) should be darker: %d vs %d", toByte(dark.R), toByte(c.R))
	}
	// Brighter then Darker by same amount should approximately return,
	// modulo clamping at the bright step (channels that brighten to 1.0
	// stay clamped). Use a small amount to avoid the clamp ceiling.
	orig := RGBColor("#808080")
	round := orig.Brighter(0.5).Darker(0.5)
	if math.Abs(round.R-orig.R) > 1e-9 {
		t.Errorf("Brighter(0.5).Darker(0.5) != identity: got %v want %v",
			round.R, orig.R)
	}
}

func TestSetOpacity(t *testing.T) {
	c := RGBColor("#ff0000")
	c.SetOpacity(0.5)
	if c.Opac != 0.5 {
		t.Errorf("opacity = %v want 0.5", c.Opac)
	}
	if c.String() != "rgba(255, 0, 0, 0.5)" {
		t.Errorf("String() = %q want rgba(255, 0, 0, 0.5)", c.String())
	}
	c.SetOpacity(1)
	if c.String() != "#ff0000" {
		t.Errorf("String() = %q want #ff0000", c.String())
	}
}

func TestApplyModifiers(t *testing.T) {
	// darker(1) should darken
	out, err := ApplyModifiers("#ff0000", [][2]any{{"darker", 1.0}})
	if err != nil {
		t.Fatalf("ApplyModifiers: %v", err)
	}
	// darker(1) on red = rgb(178,0,0) ≈ #b30000
	if out != "#b30000" {
		t.Errorf("darker(1) on red = %q want #b30000", out)
	}
	// opacity modifier changes alpha
	out, _ = ApplyModifiers("#ff0000", [][2]any{{"opacity", 0.5}})
	if out != "rgba(255, 0, 0, 0.5)" {
		t.Errorf("opacity modifier = %q want rgba(255, 0, 0, 0.5)", out)
	}
	// invalid modifier
	_, err = ApplyModifiers("#ff0000", [][2]any{{"garbage", 1.0}})
	if err == nil {
		t.Errorf("expected error for invalid modifier")
	}
}

func TestRoundtripOpaqueHex(t *testing.T) {
	cases := []string{"#000000", "#ffffff", "#ff0000", "#00ff00", "#0000ff", "#123456", "#abcdef"}
	for _, in := range cases {
		got := RGBColor(in).String()
		if got != in {
			t.Errorf("roundtrip %s -> %q want %s", in, got, in)
		}
	}
}

func TestUnrecognized(t *testing.T) {
	// unknown string falls back to opaque black
	c := RGBColor("not-a-color")
	if c.String() != "#000000" {
		t.Errorf("unrecognized color String() = %q want #000000", c.String())
	}
}
