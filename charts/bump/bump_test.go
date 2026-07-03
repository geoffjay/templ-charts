package bump_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []bump.BumpSerie {
	return []bump.BumpSerie{
		{ID: "Serie 1", Data: []bump.BumpDatum{
			{X: "2018", Y: 1.0}, {X: "2019", Y: 2.0}, {X: "2020", Y: 1.0}, {X: "2021", Y: 3.0},
		}},
		{ID: "Serie 2", Data: []bump.BumpDatum{
			{X: "2018", Y: 2.0}, {X: "2019", Y: 1.0}, {X: "2020", Y: 3.0}, {X: "2021", Y: 1.0},
		}},
		{ID: "Serie 3", Data: []bump.BumpDatum{
			{X: "2018", Y: 3.0}, {X: "2019", Y: 3.0}, {X: "2020", Y: 2.0}, {X: "2021", Y: 2.0},
		}},
	}
}

func baseProps() bump.BumpProps {
	return bump.BumpProps{
		Width: 600, Height: 400,
		Margin: core.Margin{Top: 30, Right: 100, Bottom: 40, Left: 100},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props bump.BumpProps) string {
	t.Helper()
	var b strings.Builder
	if err := bump.Bump(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Bump.Render: %v", err)
	}
	return b.String()
}

func TestBump_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestBump_PointCount(t *testing.T) {
	// 3 series × 4 columns = 12 rank dots, each a <circle>.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 12 {
		t.Errorf("circle count = %d, want 12", got)
	}
}

func TestBump_LineCount(t *testing.T) {
	// One <path> line per serie (three) plus grid <line>s; assert the paths.
	out := renderChart(t, baseProps())
	if got := strings.Count(out, `fill="none" stroke=`); got != 3 {
		t.Errorf("serie line path count = %d, want 3", got)
	}
}

func TestBump_MissingRankBreaksLine(t *testing.T) {
	p := baseProps()
	p.Data[0].Data[1].Y = nil // drop a rank
	out := renderChart(t, p)
	// One fewer drawn point (11 instead of 12).
	if got := strings.Count(out, "<circle"); got != 11 {
		t.Errorf("circle count with a gap = %d, want 11", got)
	}
}

func TestBump_LinearInterpolation(t *testing.T) {
	p := baseProps()
	p.Interpolation = bump.InterpolationLinear
	out := renderChart(t, p)
	// Linear paths use L commands, not the C (cubic) of the smooth default.
	if !strings.Contains(out, `d="M`) {
		t.Errorf("expected a line path")
	}
}

func TestBump_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "bump-basic", out)
}

func TestBump_GoldenLinear(t *testing.T) {
	p := baseProps()
	p.Interpolation = bump.InterpolationLinear
	golden.Assert(t, "bump-linear", renderChart(t, p))
}

func TestBump_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Ranking over time"
	p.Desc = "Three series ranked across four years."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Ranking over time</title>") {
		t.Errorf("expected <title> threaded through to SvgWrapper")
	}
	if !strings.Contains(out, "<desc>Three series ranked across four years.</desc>") {
		t.Errorf("expected <desc> threaded through to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestBump_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive bump should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive bump must not emit data-tc-tooltip")
	}
}

func TestBump_Legend(t *testing.T) {
	p := baseProps()
	p.Legends = []legends.LegendProps{
		{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 90},
	}
	out := renderChart(t, p)
	if !strings.Contains(out, "Serie 1") {
		t.Errorf("expected legend to include serie ids")
	}
}

func TestBump_VoronoiMesh(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	p.UseMesh = true
	out := renderChart(t, p)
	if !strings.Contains(out, "data-tc-mesh") {
		t.Errorf("Interactive+UseMesh bump should emit a voronoi-mesh overlay")
	}
	// Default (no UseMesh) must not emit a mesh overlay.
	if strings.Contains(renderChart(t, baseProps()), "data-tc-mesh") {
		t.Errorf("bump without UseMesh must not emit data-tc-mesh")
	}
	// Debug draws the voronoi cells (faint red).
	p.DebugMesh = true
	if !strings.Contains(renderChart(t, p), `stroke="red"`) {
		t.Errorf("DebugMesh should draw the voronoi cells")
	}
}

func TestBump_AnimateFadeIn(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.05
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("bump with Animate=true should emit <animate> fade-in")
	}
	if strings.Contains(renderChart(t, baseProps()), "<animate") {
		t.Errorf("bump with Animate=false must not emit <animate>")
	}
}
