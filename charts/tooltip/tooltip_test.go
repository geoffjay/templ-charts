package tooltip

import "testing"

func TestComputeCrosshairLines_Cross(t *testing.T) {
	lines := ComputeCrosshairLines(CrosshairProps{Type: CrosshairTypeCross, Width: 100, Height: 50, X: 30, Y: 20})
	if len(lines) != 2 {
		t.Fatalf("cross type got %d lines, want 2", len(lines))
	}
	// Horizontal line: full width at y.
	if lines[0].Y1 != 20 || lines[0].Y2 != 20 || lines[0].X1 != 0 || lines[0].X2 != 100 {
		t.Fatalf("horizontal line wrong: %+v", lines[0])
	}
	// Vertical line: full height at x.
	if lines[1].X1 != 30 || lines[1].X2 != 30 || lines[1].Y1 != 0 || lines[1].Y2 != 50 {
		t.Fatalf("vertical line wrong: %+v", lines[1])
	}
}

func TestComputeCrosshairLines_TopLeft(t *testing.T) {
	lines := ComputeCrosshairLines(CrosshairProps{Type: CrosshairTypeTopLeft, Width: 100, Height: 50, X: 30, Y: 20})
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	// Horizontal from 0 to x at y.
	if lines[0].X1 != 0 || lines[0].X2 != 30 || lines[0].Y1 != 20 {
		t.Fatalf("h-line wrong: %+v", lines[0])
	}
	// Vertical from 0 to y at x.
	if lines[1].Y1 != 0 || lines[1].Y2 != 20 || lines[1].X1 != 30 {
		t.Fatalf("v-line wrong: %+v", lines[1])
	}
}

func TestComputeCrosshairLines_BottomRight(t *testing.T) {
	lines := ComputeCrosshairLines(CrosshairProps{Type: CrosshairTypeBottomRight, Width: 100, Height: 50, X: 30, Y: 20})
	if len(lines) != 2 {
		t.Fatalf("got %d, want 2", len(lines))
	}
	if lines[0].X1 != 100 || lines[0].X2 != 30 {
		t.Fatalf("h-line from right wrong: %+v", lines[0])
	}
	if lines[1].Y1 != 50 || lines[1].Y2 != 20 {
		t.Fatalf("v-line from bottom wrong: %+v", lines[1])
	}
}

func TestComputeCrosshairLines_Single(t *testing.T) {
	for _, tc := range []struct {
		typ   CrosshairType
		count int
	}{
		{CrosshairTypeTop, 1},
		{CrosshairTypeBottom, 1},
		{CrosshairTypeLeft, 1},
		{CrosshairTypeRight, 1},
		{CrosshairTypeMiddle, 1},
	} {
		lines := ComputeCrosshairLines(CrosshairProps{Type: tc.typ, Width: 100, Height: 50, X: 30, Y: 20})
		if len(lines) != tc.count {
			t.Fatalf("%s got %d lines, want %d", tc.typ, len(lines), tc.count)
		}
	}
}

func TestLabelOrValue(t *testing.T) {
	if got := labelOrValue(BasicTooltipProps{ID: "foo", FormattedValue: "42"}); got != "foo: 42" {
		t.Fatalf("got %q, want 'foo: 42'", got)
	}
	if got := labelOrValue(BasicTooltipProps{ID: "foo"}); got != "foo" {
		t.Fatalf("got %q, want 'foo'", got)
	}
}

func TestStyleFromMap(t *testing.T) {
	s := styleFromMap(map[string]any{"color": "red", "fontSize": "11px"})
	if s == "" {
		t.Fatalf("expected non-empty style")
	}
}

func TestKebabCase(t *testing.T) {
	if got := kebabCase("fontSize"); got != "font-size" {
		t.Fatalf("got %q", got)
	}
	if got := kebabCase("borderRadius"); got != "border-radius" {
		t.Fatalf("got %q", got)
	}
}
