package theming

import "testing"

func TestNewUniformBorderRadius(t *testing.T) {
	r := NewUniformBorderRadius(6)
	if r.Uniform == nil || *r.Uniform != 6 {
		t.Fatalf("Uniform = %v, want 6", r.Uniform)
	}
	c := NormalizeBorderRadius(r)
	if c != (BorderRadiusCorners{6, 6, 6, 6}) {
		t.Errorf("normalized = %+v, want all 6", c)
	}
}

func TestNormalizeBorderRadius_ExplicitCornerWinsOverGroup(t *testing.T) {
	top, tl := 2.0, 9.0
	c := NormalizeBorderRadius(BorderRadius{Object: BorderRadiusObject{Top: &top, TopLeft: &tl}})
	if c.TopLeft != 9 {
		t.Errorf("topLeft = %v, want explicit corner 9", c.TopLeft)
	}
	if c.TopRight != 2 {
		t.Errorf("topRight = %v, want group 2", c.TopRight)
	}
	if c.BottomLeft != 0 || c.BottomRight != 0 {
		t.Errorf("bottom corners = %v/%v, want 0", c.BottomLeft, c.BottomRight)
	}
}

func TestConstrainBorderRadius_ScalesDownAndClampsNegative(t *testing.T) {
	// Uniform 10 on a 10x10 box: each corner pair sums to 20 > 10, so each
	// corner is scaled by 0.5 twice? No — first the width pass halves them to 5,
	// then the height pass sees 5+5=10 which fits. All corners become 5.
	c := ConstrainBorderRadius(NewUniformBorderRadius(10), 10, 10)
	if c != (BorderRadiusCorners{5, 5, 5, 5}) {
		t.Errorf("constrained = %+v, want all 5", c)
	}
	// Fits: passes through untouched.
	c = ConstrainBorderRadius(NewUniformBorderRadius(3), 100, 100)
	if c != (BorderRadiusCorners{3, 3, 3, 3}) {
		t.Errorf("unconstrained = %+v, want all 3", c)
	}
	// Negative radii clamp to 0.
	neg := -4.0
	c = ConstrainBorderRadius(BorderRadius{Object: BorderRadiusObject{TopLeft: &neg}}, 100, 100)
	if c.TopLeft != 0 {
		t.Errorf("negative topLeft = %v, want 0", c.TopLeft)
	}
	// Bottom pair exceeding the width is scaled independently of the top pair.
	blv, brv := 8.0, 8.0
	c = ConstrainBorderRadius(BorderRadius{Object: BorderRadiusObject{BottomLeft: &blv, BottomRight: &brv}}, 8, 100)
	if c.BottomLeft != 4 || c.BottomRight != 4 || c.TopLeft != 0 {
		t.Errorf("bottom width-constrained = %+v, want bottoms 4", c)
	}
	// Right pair exceeding the height is scaled independently.
	trv, brv2 := 20.0, 20.0
	c = ConstrainBorderRadius(BorderRadius{Object: BorderRadiusObject{TopRight: &trv, BottomRight: &brv2}}, 100, 20)
	if c.TopRight != 10 || c.BottomRight != 10 {
		t.Errorf("right height-constrained = %+v, want rights 10", c)
	}
	// Height constraint kicks in independently: tall corners on a short box.
	tl, bl := 30.0, 30.0
	c = ConstrainBorderRadius(BorderRadius{Object: BorderRadiusObject{TopLeft: &tl, BottomLeft: &bl}}, 100, 30)
	if c.TopLeft != 15 || c.BottomLeft != 15 {
		t.Errorf("height-constrained left corners = %v/%v, want 15/15", c.TopLeft, c.BottomLeft)
	}
}

func TestBorderRadiusToCss(t *testing.T) {
	cases := []struct {
		in   BorderRadiusCorners
		want string
	}{
		{BorderRadiusCorners{4, 4, 0, 0}, "4px 4px 0 0"},
		{BorderRadiusCorners{0, 0, 0, 0}, "0 0 0 0"},
		{BorderRadiusCorners{2.5, 0, 1, 10}, "2.5px 0 1px 10px"},
		{BorderRadiusCorners{-3, 0, 0, 0}, "-3px 0 0 0"},
		{BorderRadiusCorners{123, 0, 0, 0}, "123px 0 0 0"},
	}
	for _, c := range cases {
		if got := BorderRadiusToCss(c.in); got != c.want {
			t.Errorf("BorderRadiusToCss(%+v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestItoa(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{0, "0"}, {7, "7"}, {42, "42"}, {-13, "-13"}, {1000, "1000"},
	}
	for _, c := range cases {
		if got := itoa(c.in); got != c.want {
			t.Errorf("itoa(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFtoa(t *testing.T) {
	if got := ftoa(4); got != "4" {
		t.Errorf("ftoa(4) = %q, want 4", got)
	}
	if got := ftoa(2.5); got != "2.5" {
		t.Errorf("ftoa(2.5) = %q, want 2.5", got)
	}
	if got := ftoaTrim(0.125); got != "0.125" {
		t.Errorf("ftoaTrim(0.125) = %q", got)
	}
}

func TestOrFallback(t *testing.T) {
	v := 3.0
	if got := orFallback(&v, 9); got != 3 {
		t.Errorf("orFallback(&3, 9) = %v", got)
	}
	if got := orFallback(nil, 9); got != 9 {
		t.Errorf("orFallback(nil, 9) = %v", got)
	}
}
