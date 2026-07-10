package heatmap

import (
	"context"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/theming"
)

func TestFmtH_Guards(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{math.NaN(), "0"},
		{math.Inf(1), "0"},
		{math.Inf(-1), "0"},
		{1.5, "1.5"},
		{2, "2"},
		{1.230, "1.23"},
		{-0.0001, "0"},
	}
	for _, c := range cases {
		if got := fmtH(c.in); got != c.want {
			t.Errorf("fmtH(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaxH(t *testing.T) {
	if got := maxH(3, 2); got != 3 {
		t.Errorf("maxH(3,2) = %v, want 3", got)
	}
	if got := maxH(-1, 0); got != 0 {
		t.Errorf("maxH(-1,0) = %v, want 0", got)
	}
}

func TestComputeLayout(t *testing.T) {
	// No forceSquare: identity.
	ox, oy, w, h := computeLayout(300, 150, 3, 3, false)
	if ox != 0 || oy != 0 || w != 300 || h != 150 {
		t.Errorf("plain layout = (%v,%v,%v,%v)", ox, oy, w, h)
	}
	// forceSquare with zero cols: identity guard.
	ox, oy, w, h = computeLayout(300, 150, 0, 3, true)
	if ox != 0 || oy != 0 || w != 300 || h != 150 {
		t.Errorf("zero-cols layout = (%v,%v,%v,%v)", ox, oy, w, h)
	}
	// forceSquare: squared + centered.
	ox, oy, w, h = computeLayout(300, 150, 3, 3, true)
	if ox != 75 || oy != 0 || w != 150 || h != 150 {
		t.Errorf("forceSquare layout = (%v,%v,%v,%v), want (75,0,150,150)", ox, oy, w, h)
	}
}

func TestResolveAxis_Nil(t *testing.T) {
	if got := resolveAxis(nil, "x", nil, 100, 0, 0, "before"); got != nil {
		t.Errorf("resolveAxis(nil, …) = %v, want nil", got)
	}
}

func TestValueFormatter(t *testing.T) {
	plain := valueFormatter("")
	if got := plain(1.5); got != "1.5" {
		t.Errorf("plain formatter = %q, want 1.5", got)
	}
	spec := valueFormatter(".2f")
	if got := spec(1.5); got != "1.50" {
		t.Errorf(".2f formatter = %q, want 1.50", got)
	}
}

func TestThemeHelpers(t *testing.T) {
	if resolveTheme(nil) != &theming.DefaultTheme {
		t.Errorf("resolveTheme(nil) must return the default theme")
	}
	th := theming.DefaultTheme
	th.Background = "#101010"
	if resolveTheme(&th) != &th {
		t.Errorf("resolveTheme must return the given theme")
	}
	if got := themeBackground(nil); got != "" {
		t.Errorf("themeBackground(nil) = %q, want empty", got)
	}
	if got := themeBackground(&th); got != "#101010" {
		t.Errorf("themeBackground = %q", got)
	}
}

// failingComponent always errors, to exercise renderComponent's error path.
var failingComponent = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
	return errors.New("boom")
})

func TestRenderComponent_ErrorReturnsEmpty(t *testing.T) {
	if got := renderComponent(failingComponent); got != "" {
		t.Errorf("renderComponent(failing) = %q, want empty string", got)
	}
}

// TestComponents_CancelledContext verifies the templ components propagate a
// cancelled context as a render error.
func TestComponents_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	components := map[string]templ.Component{
		"HeatMap":     HeatMap(canvasSampleProps()),
		"HeatMapCell": HeatMapCell(HeatMapCellProps{}),
	}
	for name, c := range components {
		var b strings.Builder
		if err := c.Render(ctx, &b); err == nil {
			t.Errorf("%s.Render with cancelled context: expected error, got nil", name)
		}
	}
}
