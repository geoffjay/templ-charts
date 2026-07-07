package line

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

func TestPointXYFromPosition(t *testing.T) {
	d := ComputedDatum{}
	d.Position.X = 3
	d.Position.Y = 4
	if got := pointXYFromPosition(d); got != (PointXY{X: 3, Y: 4}) {
		t.Errorf("pointXYFromPosition = %+v, want {3 4}", got)
	}
}

func TestFmtL_Guards(t *testing.T) {
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
		if got := fmtL(c.in); got != c.want {
			t.Errorf("fmtL(%v) = %q, want %q", c.in, got, c.want)
		}
		if got := fmtA(c.in); got != c.want {
			t.Errorf("fmtA(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNilFloatHelpers(t *testing.T) {
	if !isNilFloat(nilFloat()) {
		t.Errorf("nilFloat() must be NaN (isNilFloat true)")
	}
	if isNilFloat(1) {
		t.Errorf("isNilFloat(1) must be false")
	}
	if !isNilFloat(ptrToFloat(nil)) {
		t.Errorf("ptrToFloat(nil) must be the NaN sentinel")
	}
	v := 2.5
	if got := ptrToFloat(&v); got != 2.5 {
		t.Errorf("ptrToFloat(&2.5) = %v", got)
	}
	if !isNilAny(nil) || isNilAny(0) {
		t.Errorf("isNilAny misbehaves")
	}
}

func TestAreaFill(t *testing.T) {
	props := AreasProps{Fills: []string{"url(#g0)", ""}}
	s := ComputedSeries{Color: "#abc"}
	if got := areaFill(props, 0, s); got != "url(#g0)" {
		t.Errorf("areaFill override = %q, want url(#g0)", got)
	}
	if got := areaFill(props, 1, s); got != "#abc" {
		t.Errorf("areaFill empty override = %q, want series color", got)
	}
	if got := areaFill(props, 5, s); got != "#abc" {
		t.Errorf("areaFill out of range = %q, want series color", got)
	}
}

func TestResolveThemeAndBackground(t *testing.T) {
	if resolveTheme(nil) != &theming.DefaultTheme {
		t.Errorf("resolveTheme(nil) must return the default theme")
	}
	th := theming.DefaultTheme
	th.Background = "#222222"
	if resolveTheme(&th) != &th {
		t.Errorf("resolveTheme must return the given theme")
	}
	if got := themeBackground(nil); got != "" {
		t.Errorf("themeBackground(nil) = %q, want empty", got)
	}
	if got := themeBackground(&th); got != "#222222" {
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
	gen := func(points []PointXY) string { return "M0,0" }
	components := map[string]templ.Component{
		"Line": Line(LineProps{Width: 100, Height: 100, Data: []LineSeries{
			{ID: "A", Data: []LinePointData{{X: "a", Y: float64(1)}}},
		}}),
		"Lines":        Lines(LinesProps{LineGenerator: gen}),
		"LinesItem":    LinesItem(LinesItemProps{LineGenerator: gen}),
		"Areas":        Areas(AreasProps{AreaGenerator: gen}),
		"Points":       Points(PointsProps{}),
		"Slices":       Slices(SlicesProps{}),
		"Mesh":         Mesh(MeshProps{}),
		"PointTooltip": PointTooltip(PointTooltipProps{}),
		"SliceTooltip": SliceTooltip(SliceTooltipProps{}),
	}
	for name, c := range components {
		var b strings.Builder
		if err := c.Render(ctx, &b); err == nil {
			t.Errorf("%s.Render with cancelled context: expected error, got nil", name)
		}
	}
}
