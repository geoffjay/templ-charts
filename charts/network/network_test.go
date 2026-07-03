package network_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/network"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleNodes() []network.NetworkInputNode {
	return []network.NetworkInputNode{
		{ID: "0", Color: "#e8c1a0"},
		{ID: "1", Color: "#f47560"},
		{ID: "2", Color: "#f1e15b"},
		{ID: "3", Color: "#e8a838"},
		{ID: "4", Color: "#61cdbb"},
		{ID: "5", Color: "#97e3d5"},
	}
}

func sampleLinks() []network.NetworkInputLink {
	return []network.NetworkInputLink{
		{Source: "0", Target: "1"},
		{Source: "0", Target: "2"},
		{Source: "1", Target: "3"},
		{Source: "2", Target: "4"},
		{Source: "3", Target: "5"},
		{Source: "4", Target: "5"},
	}
}

func baseProps() network.NetworkProps {
	return network.NetworkProps{
		Width: 500, Height: 500,
		Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
		Nodes:  sampleNodes(),
		Links:  sampleLinks(),
	}
}

func renderChart(t *testing.T, props network.NetworkProps) string {
	t.Helper()
	var b strings.Builder
	if err := network.Network(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Network.Render: %v", err)
	}
	return b.String()
}

func TestNetwork_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
}

func TestNetwork_NodeAndLinkCount(t *testing.T) {
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 6 {
		t.Errorf("circle count = %d, want 6", got)
	}
	if got := strings.Count(out, "<line"); got != 6 {
		t.Errorf("line count = %d, want 6", got)
	}
}

func TestNetwork_Golden(t *testing.T) {
	golden.Assert(t, "network-basic", renderChart(t, baseProps()))
}

func TestNetwork_FitViewFillsArea(t *testing.T) {
	p := baseProps() // 500×500 outer, 20px margins → 460×460 inner
	p.FitView = true
	res := network.UseNetwork(mustInner(p))
	minX, maxX := res.Nodes[0].X, res.Nodes[0].X
	minY, maxY := res.Nodes[0].Y, res.Nodes[0].Y
	maxR := 0.0
	for _, n := range res.Nodes {
		minX, maxX = min(minX, n.X), max(maxX, n.X)
		minY, maxY = min(minY, n.Y), max(maxY, n.Y)
		maxR = max(maxR, n.Size/2)
	}
	const inner = 460.0
	// The fitted layout should span (almost) the full inner area on its larger
	// axis, inset by the node radius — i.e. cover well over half of it.
	spanX, spanY := maxX-minX, maxY-minY
	if spanX < inner*0.5 && spanY < inner*0.5 {
		t.Errorf("FitView layout too small: spanX=%.1f spanY=%.1f (inner=%.0f)", spanX, spanY, inner)
	}
	// And it must stay within the frame (allowing for the node radius).
	if minX < -1 || minY < -1 || maxX > inner+1 || maxY > inner+1 {
		t.Errorf("FitView layout out of bounds: x[%.1f,%.1f] y[%.1f,%.1f]", minX, maxX, minY, maxY)
	}
}

// mustInner mirrors what the templ entrypoint does before calling the hook:
// apply defaults and reduce Width/Height to the inner dimensions.
func mustInner(p network.NetworkProps) network.NetworkProps {
	p.Width = p.Width - p.Margin.Left - p.Margin.Right
	p.Height = p.Height - p.Margin.Top - p.Margin.Bottom
	if p.LinkDistance == 0 {
		p.LinkDistance = 30
	}
	if p.CenteringStrength == 0 {
		p.CenteringStrength = 1
	}
	if p.Repulsivity == 0 {
		p.Repulsivity = 10
	}
	if p.DistanceMin == 0 {
		p.DistanceMin = 1
	}
	if p.Iterations == 0 {
		p.Iterations = 120
	}
	if p.NodeSize == 0 {
		p.NodeSize = 12
	}
	return p
}

func TestNetwork_Deterministic(t *testing.T) {
	if renderChart(t, baseProps()) != renderChart(t, baseProps()) {
		t.Errorf("network render is not deterministic")
	}
}

func TestNetwork_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Graph"
	p.Desc = "A small force-directed graph."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Graph</title>") || !strings.Contains(out, "<desc>A small force-directed graph.</desc>") {
		t.Errorf("expected title/desc threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestNetwork_Interactive(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive network should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive network must not emit data-tc-tooltip")
	}
}

func TestNetwork_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.05
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("Animate=true should emit an <animate> enter transition")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("Animate=false should not emit any <animate>")
	}
}
