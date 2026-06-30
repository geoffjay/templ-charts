package stream_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/stream"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []stream.StreamDatum {
	return []stream.StreamDatum{
		{"Raoul": 10, "Josiane": 20, "Marcel": 30},
		{"Raoul": 15, "Josiane": 18, "Marcel": 25},
		{"Raoul": 12, "Josiane": 22, "Marcel": 28},
		{"Raoul": 18, "Josiane": 16, "Marcel": 35},
		{"Raoul": 14, "Josiane": 24, "Marcel": 20},
	}
}

func baseProps() stream.StreamProps {
	return stream.StreamProps{
		Width: 600, Height: 400,
		Margin: core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 50},
		Data:   sampleData(),
		Keys:   []string{"Raoul", "Josiane", "Marcel"},
	}
}

func render(t *testing.T, props stream.StreamProps) string {
	t.Helper()
	var b strings.Builder
	if err := stream.Stream(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Stream.Render: %v", err)
	}
	return b.String()
}

func TestStream_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestStream_LayerPathPerKey(t *testing.T) {
	// One area <path d=…> per key (3). The y-grid also emits <line>s but no
	// <path d="M">, so count area paths.
	out := render(t, baseProps())
	if got := strings.Count(out, `<path d="M`); got != 3 {
		t.Errorf("area path count = %d, want 3", got)
	}
}

func TestStream_OffsetTypes(t *testing.T) {
	for _, off := range []core.StackOffset{
		core.StackOffsetWiggle, core.StackOffsetSilhouette, core.StackOffsetExpand, core.StackOffsetNone,
	} {
		p := baseProps()
		p.OffsetType = off
		out := render(t, p)
		if got := strings.Count(out, `<path d="M`); got != 3 {
			t.Errorf("offset %s: area path count = %d, want 3", off, got)
		}
	}
}

func TestStream_Golden(t *testing.T) {
	out := render(t, baseProps())
	golden.Assert(t, "stream-wiggle", out)
}

func TestStream_Golden_Expand(t *testing.T) {
	p := baseProps()
	p.OffsetType = core.StackOffsetExpand
	p.BorderWidth = 1
	out := render(t, p)
	golden.Assert(t, "stream-expand", out)
}

func TestStream_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(render(t, p), "data-tc-tooltip") {
		t.Errorf("interactive stream should emit data-tc-tooltip")
	}
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive stream must not emit data-tc-tooltip")
	}
}
