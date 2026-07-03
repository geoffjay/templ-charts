package pie_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []any {
	return []any{
		map[string]any{"id": "A", "value": float64(10)},
		map[string]any{"id": "B", "value": float64(20)},
		map[string]any{"id": "C", "value": float64(30)},
	}
}

func renderChart(t *testing.T, props pie.PieProps) string {
	t.Helper()
	var b strings.Builder
	if err := pie.Pie(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Pie.Render: %v", err)
	}
	return b.String()
}

func TestPie_RendersSVG(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: sampleData(),
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got: %q", out[:minLen(out)])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("output missing closing </svg>")
	}
}

func TestPie_ArcCount(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: sampleData(),
	}
	out := renderChart(t, props)
	// 3 arcs → 3 <path> elements in the arcs layer.
	pathCount := strings.Count(out, "<path")
	if pathCount < 3 {
		t.Errorf("expected at least 3 <path> arcs, got %d", pathCount)
	}
}

func TestPie_Donut(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		InnerRadius: 0.5,
	}
	out := renderChart(t, props)
	// Donut should still produce 3 arc paths.
	pathCount := strings.Count(out, "<path")
	if pathCount < 3 {
		t.Errorf("expected at least 3 <path> arcs for donut, got %d", pathCount)
	}
}

func TestPie_HalfPie(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:       sampleData(),
		StartAngle: 0,
		EndAngle:   180,
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
	// Half-pie should still render 3 arcs.
	pathCount := strings.Count(out, "<path")
	if pathCount < 3 {
		t.Errorf("expected at least 3 <path> arcs for half-pie, got %d", pathCount)
	}
}

func TestPie_ArcLinkLabels(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:                sampleData(),
		EnableArcLinkLabels: true,
	}
	out := renderChart(t, props)
	// Arc link labels emit <path> (link) + <text> (label) per arc.
	// Check for link label text content ("A", "B", "C").
	if !strings.Contains(out, ">A<") {
		t.Errorf("expected arc link label 'A', not found")
	}
}

func TestPie_ArcLabels(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:            sampleData(),
		EnableArcLabels: true,
	}
	out := renderChart(t, props)
	// Arc labels emit <text> with the formatted value. Default arcLabel is
	// "formattedValue" which formats the numeric value.
	if !strings.Contains(out, "<text") {
		t.Errorf("expected arc label <text>, not found")
	}
}

func TestPie_DisableArcLinkLabels(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:                sampleData(),
		EnableArcLinkLabels: false,
	}
	out := renderChart(t, props)
	// Without arc link labels, the "A"/"B"/"C" text labels should not appear
	// (arc labels use formatted values, not ids). We check the link paths are
	// absent: link labels emit a second <path d="M…L…L…"> per arc.
	// The arc paths use d="M…" from the arc generator; link paths also use
	// d="M…". We count total paths — 3 arcs + 3 links = 6; without links = 3.
	pathCount := strings.Count(out, "<path")
	if pathCount > 5 {
		t.Errorf("expected no arc link label paths, but found %d total paths", pathCount)
	}
}

func TestPie_AnimationsEmitSMIL(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: sampleData(),
	}
	outOff := renderChart(t, props)
	if strings.Contains(outOff, "<animate") {
		t.Errorf("expected no <animate> when Animate=false, found one")
	}
	props.Animate = true
	outOn := renderChart(t, props)
	animateCount := strings.Count(outOn, "<animate")
	if animateCount == 0 {
		t.Errorf("expected <animate> elements when Animate=true, found 0")
	}
}

func TestPie_SortByValue(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: []any{
			map[string]any{"id": "small", "value": float64(5)},
			map[string]any{"id": "big", "value": float64(50)},
			map[string]any{"id": "mid", "value": float64(25)},
		},
		SortByValue: true,
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>")
	}
}

func TestPie_Legends(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: sampleData(),
		Legends: []legends.LegendProps{
			{
				Anchor: legends.LegendAnchorTopRight, Direction: legends.LegendDirectionColumn,
				ItemWidth: 80, ItemHeight: 20,
			},
		},
	}
	out := renderChart(t, props)
	// Legend emits a <g> with a symbol + <text>.
	if !strings.Contains(out, "<text") {
		t.Errorf("expected legend <text>, not found")
	}
}

func TestPie_PadAngle(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:     sampleData(),
		PadAngle: 2,
	}
	out := renderChart(t, props)
	pathCount := strings.Count(out, "<path")
	if pathCount < 3 {
		t.Errorf("expected at least 3 <path> arcs with padAngle, got %d", pathCount)
	}
}

func TestPie_CornerRadius(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:         sampleData(),
		CornerRadius: 5,
	}
	out := renderChart(t, props)
	pathCount := strings.Count(out, "<path")
	if pathCount < 3 {
		t.Errorf("expected at least 3 <path> arcs with cornerRadius, got %d", pathCount)
	}
}

func TestPie_CustomColors(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data: sampleData(),
		Colors: colors.OrdinalColorScaleConfig{
			Type:   colors.OrdinalTypeColors,
			Colors: []string{"#ff0000", "#00ff00", "#0000ff"},
		},
	}
	out := renderChart(t, props)
	if !strings.Contains(out, "#ff0000") {
		t.Errorf("expected custom color #ff0000 in output")
	}
}

func TestPie_BorderWidth(t *testing.T) {
	props := pie.PieProps{
		Width: 500, Height: 300,
		Data:        sampleData(),
		BorderWidth: 2,
	}
	out := renderChart(t, props)
	if !strings.Contains(out, `stroke-width="2"`) {
		t.Errorf("expected stroke-width=2 for border, not found")
	}
}

// TestPie_Golden renders a donut pie chart with arc + arc-link labels and
// compares the full SVG against a committed golden snapshot. Regenerate
// after an intentional render change with:
//
//	go test ./charts/pie -run TestPie_Golden -update
func TestPie_Golden(t *testing.T) {
	props := pie.PieProps{
		Width:        500,
		Height:       300,
		InnerRadius:  0.5,
		PadAngle:     0.5,
		CornerRadius: 3,
		Data:         sampleData(),
	}
	out := renderChart(t, props)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:minLen(out)])
	}
	golden.Assert(t, "pie-donut", out)
}

// TestPie_Golden_Half renders a half-pie (startAngle=0, endAngle=180, fit)
// and compares against a committed golden snapshot.
func TestPie_Golden_Half(t *testing.T) {
	props := pie.PieProps{
		Width:      500,
		Height:     300,
		StartAngle: 0,
		EndAngle:   180,
		Fit:        true,
		Data:       sampleData(),
	}
	out := renderChart(t, props)
	golden.Assert(t, "pie-half", out)
}

func minLen(s string) int {
	if len(s) < 50 {
		return len(s)
	}
	return 50
}
