package arcs_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
	"github.com/geoffjay/templ-charts/charts/arcs"
)

// TestMain shrinks templ's render buffer to 1 byte so writer failures surface
// at the write site instead of being deferred to the final buffer flush. The
// rendered output is byte-identical either way; this only affects when errors
// from a broken io.Writer are observed, which TestComponents_WriteErrors
// relies on.
func TestMain(m *testing.M) {
	templruntime.DefaultBufferSize = 1
	os.Exit(m.Run())
}

var errWriterBroken = errors.New("writer broken")

// failingWriter accepts `remaining` bytes and then fails every write.
type failingWriter struct{ remaining int }

func (w *failingWriter) Write(p []byte) (int, error) {
	if len(p) > w.remaining {
		n := w.remaining
		w.remaining = 0
		return n, errWriterBroken
	}
	w.remaining -= len(p)
	return len(p), nil
}

// allComponents returns fully-populated instances of every templ component in
// the package, so every conditional attribute/branch is on the render path.
func allComponents() map[string]templ.Component {
	shape := arcs.ArcShapeProps{
		Path: "M0,0L10,0Z", Fill: "#f00",
		Stroke: "#000", StrokeWidth: 1, Opacity: 0.5,
		HxGet: "/h", HxTrigger: "mouseenter", HxSwap: "innerHTML", HxTarget: "#t",
		AriaLabel: "a", DataTooltip: "tip",
		Animate: true, AnimateFromPath: "M0,0",
	}
	label := arcs.ArcLabelProps{X: 1, Y: 2, Label: "l", Fill: "#333", FontSize: 10, FontFamily: "sans"}
	link := arcs.ArcLinkLabelProps{
		Link:  arcs.ArcLink{P0x: 1, P0y: 2, P1x: 3, P1y: 4, P2x: 5, P2y: 4, TextAnchor: "start", Side: "right"},
		Label: "l", Thickness: 1, Color: "#000", TextOffset: 2,
		Fill: "#333", FontSize: 10, FontFamily: "sans",
	}
	return map[string]templ.Component{
		"ArcShape": arcs.ArcShape(shape),
		"ArcsLayer": arcs.ArcsLayer(arcs.ArcsLayerProps{
			CenterX: 1, CenterY: 2,
			Items: []arcs.ArcLayerItem{{Props: shape}},
		}),
		"ArcLabel": arcs.ArcLabel(label),
		"ArcLabelsLayer": arcs.ArcLabelsLayer(arcs.ArcLabelsLayerProps{
			CenterX: 1, CenterY: 2,
			Items: []arcs.ArcLabelItem{{Props: label}},
		}),
		"ArcLinkLabel": arcs.ArcLinkLabel(link),
		"ArcLinkLabelsLayer": arcs.ArcLinkLabelsLayer(arcs.ArcLinkLabelsLayerProps{
			CenterX: 1, CenterY: 2,
			Items: []arcs.ArcLinkLabelItem{{Props: link}},
		}),
	}
}

// TestComponents_CancelledContext verifies rendering aborts up-front when the
// context is already cancelled.
func TestComponents_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for name, c := range allComponents() {
		t.Run(name, func(t *testing.T) {
			var b strings.Builder
			err := c.Render(ctx, &b)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("err = %v, want context.Canceled", err)
			}
			if b.Len() != 0 {
				t.Errorf("cancelled render must not write output, got %q", b.String())
			}
		})
	}
}

// TestComponents_WriteErrors verifies write errors from the destination
// propagate out of Render no matter where in the output they occur. Sweeping
// the failure point across the whole output exercises every write site.
func TestComponents_WriteErrors(t *testing.T) {
	for name, c := range allComponents() {
		t.Run(name, func(t *testing.T) {
			// Establish the full output length with a working writer.
			var b strings.Builder
			if err := c.Render(context.Background(), &b); err != nil {
				t.Fatalf("baseline render: %v", err)
			}
			total := b.Len()
			if total == 0 {
				t.Fatalf("baseline render produced no output")
			}
			for n := 0; n < total; n++ {
				err := c.Render(context.Background(), &failingWriter{remaining: n})
				if !errors.Is(err, errWriterBroken) {
					t.Fatalf("failure after %d bytes: err = %v, want errWriterBroken", n, err)
				}
			}
			// A writer with just enough capacity succeeds.
			if err := c.Render(context.Background(), &failingWriter{remaining: total}); err != nil {
				t.Errorf("render with exact capacity: %v", err)
			}
		})
	}
}
