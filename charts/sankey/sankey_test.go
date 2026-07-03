package sankey_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/sankey"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleNodes() []sankey.SankeyInputNode {
	return []sankey.SankeyInputNode{
		{ID: "A"}, {ID: "B"}, {ID: "C"}, {ID: "D"}, {ID: "E"},
	}
}

func sampleLinks() []sankey.SankeyInputLink {
	return []sankey.SankeyInputLink{
		{Source: "A", Target: "C", Value: 5},
		{Source: "B", Target: "C", Value: 3},
		{Source: "A", Target: "D", Value: 2},
		{Source: "C", Target: "E", Value: 6},
		{Source: "D", Target: "E", Value: 2},
	}
}

func baseProps() sankey.SankeyProps {
	return sankey.SankeyProps{
		Width: 700, Height: 400,
		Margin: core.Margin{Top: 20, Right: 80, Bottom: 20, Left: 80},
		Nodes:  sampleNodes(),
		Links:  sampleLinks(),
	}
}

func renderChart(t *testing.T, props sankey.SankeyProps) string {
	t.Helper()
	var b strings.Builder
	if err := sankey.Sankey(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Sankey.Render: %v", err)
	}
	return b.String()
}

func TestSankey_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestSankey_NodeAndLinkCount(t *testing.T) {
	out := renderChart(t, baseProps())
	// Node rects carry an x attribute; the SvgWrapper background rect does not.
	if got := strings.Count(out, "<rect x="); got != 5 {
		t.Errorf("rect (node) count = %d, want 5", got)
	}
	if got := strings.Count(out, "<path"); got != 5 {
		t.Errorf("path (link) count = %d, want 5", got)
	}
	// Labels on by default: one <text> per node (plus any title/desc — none here).
	if got := strings.Count(out, "<text"); got != 5 {
		t.Errorf("text (label) count = %d, want 5", got)
	}
}

func TestSankey_Golden(t *testing.T) {
	golden.Assert(t, "sankey-basic", renderChart(t, baseProps()))
}

func TestSankey_VerticalGolden(t *testing.T) {
	p := baseProps()
	p.Layout = sankey.SankeyLayoutVertical
	golden.Assert(t, "sankey-vertical", renderChart(t, p))
}

func TestSankey_Deterministic(t *testing.T) {
	if renderChart(t, baseProps()) != renderChart(t, baseProps()) {
		t.Errorf("sankey render is not deterministic")
	}
}

func TestSankey_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Flow"
	p.Desc = "A small sankey flow diagram."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Flow</title>") || !strings.Contains(out, "<desc>A small sankey flow diagram.</desc>") {
		t.Errorf("expected title/desc threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestSankey_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive sankey should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive sankey must not emit data-tc-tooltip")
	}
}

func TestSankey_LabelsDisabled(t *testing.T) {
	p := baseProps()
	disabled := false
	p.EnableLabels = &disabled
	if strings.Count(renderChart(t, p), "<text") != 0 {
		t.Errorf("disabled labels should emit no <text> elements")
	}
}

func TestSankey_Legends(t *testing.T) {
	p := baseProps()
	p.Legends = []legends.LegendProps{{
		Anchor:    legends.LegendAnchorBottom,
		Direction: legends.LegendDirectionRow,
		ItemWidth: 80, ItemHeight: 20,
	}}
	out := renderChart(t, p)
	// Legend renders one label per node id.
	for _, id := range []string{"A", "B", "C", "D", "E"} {
		if !strings.Contains(out, ">"+id+"<") {
			t.Errorf("legend missing entry for %q", id)
		}
	}
}
