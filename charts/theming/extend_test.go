package theming

import "testing"

func TestExtendDefaultTheme_NoOverride(t *testing.T) {
	got := ExtendDefaultTheme(DefaultTheme, nil)
	if got.Background != DefaultTheme.Background {
		t.Fatalf("background = %q, want %q", got.Background, DefaultTheme.Background)
	}
	// Text inheritance: axis.ticks.text should inherit root text fill.
	if got.Axis.Ticks.Text.Fill != DefaultTheme.Text.Fill {
		t.Fatalf("axis.ticks.text.fill = %q, want inherited %q", got.Axis.Ticks.Text.Fill, DefaultTheme.Text.Fill)
	}
}

func TestExtendDefaultTheme_PartialText(t *testing.T) {
	fill := "#ff0000"
	got := ExtendDefaultTheme(DefaultTheme, &PartialTheme{Text: &TextStyle{Fill: fill}})
	if got.Text.Fill != fill {
		t.Fatalf("text.fill = %q, want %q", got.Text.Fill, fill)
	}
	// Inherited nested text should pick up the new fill.
	if got.Axis.Ticks.Text.Fill != fill {
		t.Fatalf("axis.ticks.text.fill = %q, want inherited %q", got.Axis.Ticks.Text.Fill, fill)
	}
}

func TestNormalizeBorderRadius_Uniform(t *testing.T) {
	u := 4.0
	c := NormalizeBorderRadius(BorderRadius{Uniform: &u})
	if c.TopLeft != 4 || c.TopRight != 4 || c.BottomRight != 4 || c.BottomLeft != 4 {
		t.Fatalf("got %+v, want all 4", c)
	}
}

func TestNormalizeBorderRadius_Groups(t *testing.T) {
	top := 2.0
	left := 3.0
	c := NormalizeBorderRadius(BorderRadius{Object: BorderRadiusObject{Top: &top, Left: &left}})
	if c.TopLeft != 2 {
		t.Fatalf("topLeft = %v, want 2 (top wins over left)", c.TopLeft)
	}
	if c.TopRight != 2 {
		t.Fatalf("topRight = %v, want 2", c.TopRight)
	}
	if c.BottomLeft != 3 {
		t.Fatalf("bottomLeft = %v, want 3 (left)", c.BottomLeft)
	}
	if c.BottomRight != 0 {
		t.Fatalf("bottomRight = %v, want 0", c.BottomRight)
	}
}

func TestConstrainBorderRadius(t *testing.T) {
	u := 10.0
	c := ConstrainBorderRadius(BorderRadius{Uniform: &u}, 10, 10)
	// 10+10 = 20 > width 10 → scale to 5 each.
	if c.TopLeft != 5 || c.TopRight != 5 {
		t.Fatalf("got %+v, want 5,5", c)
	}
}

func TestConvertStyleAttribute_SVG(t *testing.T) {
	if got := ConvertStyleAttribute(EngineSVG, "textAlign", "center"); got != "middle" {
		t.Fatalf("textAlign center → %q, want middle", got)
	}
	if got := ConvertStyleAttribute(EngineSVG, "textBaseline", "top"); got != "text-before-edge" {
		t.Fatalf("textBaseline top → %q, want text-before-edge", got)
	}
}

func TestSanitizeSvgTextStyle(t *testing.T) {
	s := TextStyle{FontFamily: "sans", Fill: "#000", OutlineWidth: 2, OutlineColor: "#fff"}
	out := SanitizeSvgTextStyle(s)
	if _, ok := out["outlineWidth"]; ok {
		t.Fatal("outlineWidth should be stripped")
	}
	if out["fill"] != "#000" {
		t.Fatalf("fill = %v, want #000", out["fill"])
	}
}
