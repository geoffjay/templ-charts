package bar

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

func TestFmtB_Guards(t *testing.T) {
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
		{-0.0001, "0"}, // rounds to -0.000 → normalized to 0
	}
	for _, c := range cases {
		if got := fmtB(c.in); got != c.want {
			t.Errorf("fmtB(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFmtT_Guards(t *testing.T) {
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
		if got := fmtT(c.in); got != c.want {
			t.Errorf("fmtT(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaxF(t *testing.T) {
	if got := maxF(3, 2); got != 3 {
		t.Errorf("maxF(3,2) = %v, want 3", got)
	}
	if got := maxF(-1, 0); got != 0 {
		t.Errorf("maxF(-1,0) = %v, want 0", got)
	}
}

func TestToFloatAny(t *testing.T) {
	cases := []struct {
		in   any
		want float64
	}{
		{float64(1.5), 1.5},
		{int(3), 3},
		{int64(4), 4},
		{"11px", 0},
		{nil, 0},
		{true, 0},
	}
	for _, c := range cases {
		if got := toFloatAny(c.in); got != c.want {
			t.Errorf("toFloatAny(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestThemeHelpers(t *testing.T) {
	if resolveTheme(nil) != &theming.DefaultTheme {
		t.Errorf("resolveTheme(nil) must return the default theme")
	}
	th := theming.DefaultTheme
	th.Background = "#123456"
	if resolveTheme(&th) != &th {
		t.Errorf("resolveTheme must return the given theme")
	}
	if got := themeBackground(nil); got != "" {
		t.Errorf("themeBackground(nil) = %q, want empty", got)
	}
	if got := themeBackground(&th); got != "#123456" {
		t.Errorf("themeBackground = %q, want #123456", got)
	}
	if got := themeLabelFontSize(nil); got != 0 {
		t.Errorf("themeLabelFontSize(nil) = %v, want 0", got)
	}
	th.Labels.Text.FontSize = float64(13)
	if got := themeLabelFontSize(&th); got != 13 {
		t.Errorf("themeLabelFontSize = %v, want 13", got)
	}
}

func TestTotalsAnchorBaseline(t *testing.T) {
	if got := totalsAnchor("horizontal"); got != "start" {
		t.Errorf("totalsAnchor(horizontal) = %q, want start", got)
	}
	if got := totalsAnchor("vertical"); got != "middle" {
		t.Errorf("totalsAnchor(vertical) = %q, want middle", got)
	}
	if got := totalsBaseline("horizontal"); got != "middle" {
		t.Errorf("totalsBaseline(horizontal) = %q, want middle", got)
	}
	if got := totalsBaseline("vertical"); got != "alphabetic" {
		t.Errorf("totalsBaseline(vertical) = %q, want alphabetic", got)
	}
}

func TestResolveAxis_Nil(t *testing.T) {
	if got := resolveAxis(nil, "x", nil, 100, 0, 0); got != nil {
		t.Errorf("resolveAxis(nil, …) = %v, want nil", got)
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

// TestBarEnterAnimate_Direct renders the enter-animation fragment standalone
// for both orientations.
func TestBarEnterAnimate_Direct(t *testing.T) {
	vertical := BarItemProps{Bar: ComputedBarDatum{Width: 30, Height: 50}}
	horizontal := vertical
	horizontal.Horizontal = true
	var b strings.Builder
	if err := barEnterAnimate(vertical).Render(context.Background(), &b); err != nil {
		t.Fatalf("render vertical: %v", err)
	}
	if !strings.Contains(b.String(), `attributeName="height" from="0" to="50"`) {
		t.Errorf("vertical animation must grow height to 50, got %s", b.String())
	}
	b.Reset()
	if err := barEnterAnimate(horizontal).Render(context.Background(), &b); err != nil {
		t.Fatalf("render horizontal: %v", err)
	}
	if !strings.Contains(b.String(), `attributeName="width" from="0" to="30"`) {
		t.Errorf("horizontal animation must grow width to 30, got %s", b.String())
	}
}

// TestComponents_CancelledContext verifies every exported templ component in
// the package propagates a cancelled context as a render error.
func TestComponents_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	components := map[string]templ.Component{
		"Bar": Bar(BarProps{Width: 100, Height: 100, Data: []BarDatum{
			{"id": "one", "value": float64(1)},
		}}),
		"BarItem":         BarItem(BarItemProps{}),
		"barEnterAnimate": barEnterAnimate(BarItemProps{}),
		"BarTotals":       BarTotals(BarTotalsProps{}),
		"BarLegends":      BarLegends(BarLegendsProps{}),
		"BarAnnotations":  BarAnnotations(BarAnnotationsProps{}),
	}
	for name, c := range components {
		var b strings.Builder
		if err := c.Render(ctx, &b); err == nil {
			t.Errorf("%s.Render with cancelled context: expected error, got nil", name)
		}
	}
}
