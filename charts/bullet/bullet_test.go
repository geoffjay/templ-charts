package bullet_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []bullet.BulletItemDatum {
	return []bullet.BulletItemDatum{
		{ID: "temp.", Ranges: []float64{0, 20, 80, 100}, Measures: []float64{60}, Markers: []float64{75}},
		{ID: "power", Ranges: []float64{0, 30, 70, 100}, Measures: []float64{45, 65}, Markers: []float64{80}},
	}
}

func baseProps() bullet.BulletProps {
	return bullet.BulletProps{
		Width: 700, Height: 200,
		Margin: core.Margin{Top: 20, Right: 30, Bottom: 30, Left: 80},
		Data:   sampleData(),
	}
}

func renderChart(t *testing.T, props bullet.BulletProps) string {
	t.Helper()
	var b strings.Builder
	if err := bullet.Bullet(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Bullet.Render: %v", err)
	}
	return b.String()
}

func TestBullet_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestBullet_TitlesPresent(t *testing.T) {
	out := renderChart(t, baseProps())
	for _, id := range []string{"temp.", "power"} {
		if !strings.Contains(out, ">"+id+"<") {
			t.Errorf("missing title %q", id)
		}
	}
}

func TestBullet_MarkerLines(t *testing.T) {
	// 1 marker per item × 2 items = 2 marker <line>s with stroke-width="2".
	out := renderChart(t, baseProps())
	if got := strings.Count(out, `stroke-width="2"`); got != 2 {
		t.Errorf("marker line count = %d, want 2", got)
	}
}

func TestBullet_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "bullet-basic", out)
}

func TestBullet_Golden_Vertical(t *testing.T) {
	p := baseProps()
	p.Width, p.Height = 300, 500
	p.Layout = bullet.BulletLayoutVertical
	out := renderChart(t, p)
	golden.Assert(t, "bullet-vertical", out)
}

func TestBullet_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive bullet should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive bullet must not emit data-tc-tooltip")
	}
}

func TestBullet_Animate(t *testing.T) {
	p := baseProps()
	p.Animate = true
	p.MotionStagger = 0.01
	if !strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("animated bullet should emit <animate>")
	}
	p.Animate = false
	if strings.Contains(renderChart(t, p), "<animate") {
		t.Errorf("non-animated bullet must not emit <animate>")
	}
}
