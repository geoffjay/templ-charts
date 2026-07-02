package chord_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/chord"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleKeys() []string { return []string{"A", "B", "C", "D"} }

func sampleMatrix() [][]float64 {
	return [][]float64{
		{0, 5, 6, 4},
		{7, 0, 5, 4},
		{6, 8, 0, 3},
		{4, 5, 3, 0},
	}
}

func baseProps() chord.ChordProps {
	return chord.ChordProps{
		Width: 600, Height: 600,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 40},
		Data:   sampleMatrix(),
		Keys:   sampleKeys(),
	}
}

func render(t *testing.T, props chord.ChordProps) string {
	t.Helper()
	var b strings.Builder
	if err := chord.Chord(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Chord.Render: %v", err)
	}
	return b.String()
}

func TestChord_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestChord_ArcAndRibbonCount(t *testing.T) {
	out := render(t, baseProps())
	// 4 arcs + ribbons are all <path>; ribbons carry mix-blend-mode.
	if got := strings.Count(out, "mix-blend-mode"); got == 0 {
		t.Errorf("expected ribbon paths with mix-blend-mode, got %d", got)
	}
	// Every off-diagonal pair is non-zero → 6 ribbons; plus 4 arcs = 10 paths.
	if got := strings.Count(out, "<path"); got != 10 {
		t.Errorf("path count = %d, want 10 (4 arcs + 6 ribbons)", got)
	}
	if got := strings.Count(out, "mix-blend-mode"); got != 6 {
		t.Errorf("ribbon count = %d, want 6", got)
	}
	if got := strings.Count(out, "<text"); got != 4 {
		t.Errorf("text (label) count = %d, want 4", got)
	}
}

func TestChord_Golden(t *testing.T) {
	golden.Assert(t, "chord-basic", render(t, baseProps()))
}

func TestChord_PadAngleGolden(t *testing.T) {
	p := baseProps()
	p.PadAngle = 0.04
	p.InnerRadiusRatio = 0.86
	p.InnerRadiusOffset = 0.04
	golden.Assert(t, "chord-padded", render(t, p))
}

func TestChord_Deterministic(t *testing.T) {
	if render(t, baseProps()) != render(t, baseProps()) {
		t.Errorf("chord render is not deterministic")
	}
}

func TestChord_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Flows"
	p.Desc = "A small chord diagram."
	out := render(t, p)
	if !strings.Contains(out, "<title>Flows</title>") || !strings.Contains(out, "<desc>A small chord diagram.</desc>") {
		t.Errorf("expected title/desc threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestChord_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(render(t, p), "data-tc-tooltip") {
		t.Errorf("interactive chord should emit data-tc-tooltip")
	}
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive chord must not emit data-tc-tooltip")
	}
}

func TestChord_LabelsDisabled(t *testing.T) {
	p := baseProps()
	disabled := false
	p.EnableLabel = &disabled
	if strings.Count(render(t, p), "<text") != 0 {
		t.Errorf("disabled labels should emit no <text> elements")
	}
}

func TestChord_Legends(t *testing.T) {
	p := baseProps()
	p.Legends = []legends.LegendProps{{
		Anchor:    legends.LegendAnchorBottom,
		Direction: legends.LegendDirectionRow,
		ItemWidth: 80, ItemHeight: 20,
	}}
	out := render(t, p)
	for _, id := range sampleKeys() {
		if !strings.Contains(out, ">"+id+"<") {
			t.Errorf("legend missing entry for %q", id)
		}
	}
}
