package polaraxes

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"

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

func polarComponents() map[string]templ.Component {
	return map[string]templ.Component{
		"circular-axis": CircularAxis(CircularAxisProps{
			Type: CircularAxisOuter, Radius: 80,
			StartAngle: 0, EndAngle: 360,
			Scale: angleScale(), Theme: &theming.DefaultTheme,
		}),
		"radial-axis": RadialAxis(RadialAxisProps{
			Angle: 0, Scale: radiusScale(),
			TicksPosition: TicksAfter, Theme: &theming.DefaultTheme,
		}),
		"polar-grid": PolarGrid(PolarGridProps{
			EnableRadialGrid: true, AngleScale: angleScale(),
			EnableCircularGrid: true, RadiusScale: radiusScale(),
			StartAngle: 0, EndAngle: 360, OuterRadius: 100,
			Theme: &theming.DefaultTheme,
		}),
		"radial-grid": RadialGrid(RadialGridProps{
			Scale: angleScale(), InnerRadius: 10, OuterRadius: 90,
			Theme: &theming.DefaultTheme,
		}),
		"circular-grid": CircularGrid(CircularGridProps{
			Scale: radiusScale(), StartAngle: 0, EndAngle: 360,
			Theme: &theming.DefaultTheme,
		}),
	}
}

func TestComponents_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for name, c := range polarComponents() {
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

// TestComponents_ExplicitNilChildren renders with an explicitly-nil children
// component in the context; output must match a plain render.
func TestComponents_ExplicitNilChildren(t *testing.T) {
	ctx := templ.WithChildren(context.Background(), nil)
	for name, c := range polarComponents() {
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

// TestComponents_WriteErrorPropagation verifies a failing writer's error is
// surfaced no matter which write fails.
func TestComponents_WriteErrorPropagation(t *testing.T) {
	old := templruntime.DefaultBufferSize
	templruntime.DefaultBufferSize = 1
	defer func() { templruntime.DefaultBufferSize = old }()

	for name, c := range polarComponents() {
		t.Run(name, func(t *testing.T) {
			var full strings.Builder
			if err := c.Render(context.Background(), &full); err != nil {
				t.Fatalf("baseline render: %v", err)
			}
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
