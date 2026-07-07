package axes_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"

	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// failAfterWriter fails once n bytes have been written.
type failAfterWriter struct{ n int }

var errWrite = errors.New("writer failed")

func (w *failAfterWriter) Write(p []byte) (int, error) {
	if w.n <= 0 {
		return 0, errWrite
	}
	if len(p) > w.n {
		n := w.n
		w.n = 0
		return n, errWrite
	}
	w.n -= len(p)
	return len(p), nil
}

// richAxisComponents returns components exercising most emit paths.
func richAxisComponents() map[string]templ.Component {
	legend := baseXAxisProps()
	legend.Legend = "L"
	legend.TickRotation = 45
	yAxis := baseXAxisProps()
	yAxis.Axis = "y"
	yAxis.Legend = "Y"
	x := baseXAxisProps()
	computed := axes.ComputeCartesianTicks(x, &theming.DefaultTheme)
	return map[string]templ.Component{
		"axis-x-legend-rotated": axes.Axis(legend, &theming.DefaultTheme),
		"axis-y-legend":         axes.Axis(yAxis, &theming.DefaultTheme),
		"axes":                  axes.Axes(axes.AxesProps{XAxis: &x, YAxis: &yAxis, Theme: &theming.DefaultTheme}),
		"axis-tick":             axes.AxisTick(computed, x, &theming.DefaultTheme),
		"axis-legend":           axes.AxisLegend(legend, &theming.DefaultTheme),
		"grid": axes.Grid(axes.GridProps{
			Axis: "x", Scale: linearScale(0, 100, 100), Width: 100, Height: 50,
			TickValues: []any{0.0, 50.0, 100.0}, Theme: &theming.DefaultTheme,
		}),
	}
}

func TestRender_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for name, c := range richAxisComponents() {
		t.Run(name, func(t *testing.T) {
			var b strings.Builder
			if err := c.Render(ctx, &b); !errors.Is(err, context.Canceled) {
				t.Errorf("err = %v, want context.Canceled", err)
			}
			if b.Len() != 0 {
				t.Errorf("cancelled render should write nothing, wrote %q", b.String())
			}
		})
	}
}

// TestRender_WriteErrorPropagation verifies a failing writer's error is
// surfaced no matter which write fails: it renders each component into
// writers that fail at every possible byte offset.
func TestRender_WriteErrorPropagation(t *testing.T) {
	old := templruntime.DefaultBufferSize
	templruntime.DefaultBufferSize = 1
	defer func() { templruntime.DefaultBufferSize = old }()

	for name, c := range richAxisComponents() {
		t.Run(name, func(t *testing.T) {
			var full strings.Builder
			if err := c.Render(context.Background(), &full); err != nil {
				t.Fatalf("baseline render: %v", err)
			}
			// Stop short of the end: trailing bytes may still sit in the
			// (unflushed) buffer when the caller owns it.
			for n := 0; n < full.Len()-4; n++ {
				buf := &templruntime.Buffer{}
				buf.Reset(&failAfterWriter{n: n})
				if err := c.Render(context.Background(), buf); !errors.Is(err, errWrite) {
					t.Fatalf("fail-after-%d: err = %v, want errWrite", n, err)
				}
			}
		})
	}
}

// TestRender_ExplicitNilChildren renders with an explicitly-nil children
// component in the context; output must match a plain render.
func TestRender_ExplicitNilChildren(t *testing.T) {
	ctx := templ.WithChildren(context.Background(), nil)
	for name, c := range richAxisComponents() {
		t.Run(name, func(t *testing.T) {
			var plain, withNil strings.Builder
			if err := c.Render(context.Background(), &plain); err != nil {
				t.Fatalf("plain render: %v", err)
			}
			if err := c.Render(ctx, &withNil); err != nil {
				t.Fatalf("nil-children render: %v", err)
			}
			if plain.String() != withNil.String() {
				t.Errorf("nil children changed output:\n%s\nvs\n%s", plain.String(), withNil.String())
			}
		})
	}
}

func TestAxisLegend_FontSizeFallsBackToRootText(t *testing.T) {
	theme := theming.Theme{}
	theme.Text.FontSize = 7
	props := baseXAxisProps()
	props.Legend = "L"
	var b strings.Builder
	if err := axes.AxisLegend(props, &theme).Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(b.String(), `font-size="7"`) {
		t.Errorf("legend should fall back to root text font-size: %s", b.String())
	}
}

func TestAxis_FullySpecifiedThemeText(t *testing.T) {
	// A theme whose axis text styles are fully populated: no fallback to the
	// root text style.
	theme := theming.Theme{}
	theme.Text = theming.TextStyle{Fill: "#000000", FontSize: 8, FontFamily: "mono"}
	theme.Axis.Ticks.Text = theming.TextStyle{Fill: "#111111", FontSize: 9, FontFamily: "TickFont"}
	theme.Axis.Legend.Text = theming.TextStyle{Fill: "#222222", FontSize: 10, FontFamily: "LegendFont"}
	props := baseXAxisProps()
	props.Legend = "L"
	out := render(t, axes.Axis(props, &theme))
	for _, want := range []string{
		`fill="#111111"`, `font-size="9"`, `font-family="TickFont"`,
		`fill="#222222"`, `font-size="10"`, `font-family="LegendFont"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q: %s", want, out)
		}
	}
	if strings.Contains(out, `fill="#000000"`) {
		t.Errorf("root text fill should not be used when axis text is set: %s", out)
	}
}
