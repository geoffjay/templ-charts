package parallelcoordinates_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleVars() []pc.PCVariable {
	return []pc.PCVariable{
		{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
		{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
		{Key: "weight", Type: pc.PCScaleLinear, Label: "weight"},
		{Key: "volume", Type: pc.PCScaleLinear, Label: "volume"},
	}
}

func sampleData() []pc.PCDatum {
	return []pc.PCDatum{
		{ID: "A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "weight": 30.0, "volume": 8.0}},
		{ID: "B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "weight": 12.0, "volume": 15.0}},
		{ID: "C", Values: map[string]any{"temp": 12.0, "cost": 2.0, "weight": 25.0, "volume": 3.0}},
	}
}

func baseProps() pc.PCProps {
	return pc.PCProps{
		Width: 600, Height: 400,
		Margin:    core.Margin{Top: 50, Right: 60, Bottom: 50, Left: 60},
		Data:      sampleData(),
		Variables: sampleVars(),
	}
}

func renderChart(t *testing.T, props pc.PCProps) string {
	t.Helper()
	var b strings.Builder
	if err := pc.ParallelCoordinates(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ParallelCoordinates.Render: %v", err)
	}
	return b.String()
}

func TestPC_RendersSVG(t *testing.T) {
	out := renderChart(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestPC_LineCount(t *testing.T) {
	// One <path> polyline per datum (three).
	out := renderChart(t, baseProps())
	if got := strings.Count(out, `fill="none" stroke=`); got != 3 {
		t.Errorf("line path count = %d, want 3", got)
	}
}

func TestPC_AxisLegends(t *testing.T) {
	out := renderChart(t, baseProps())
	for _, want := range []string{"temperature", "cost", "weight", "volume"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected axis legend %q", want)
		}
	}
}

func TestPC_PointVariable(t *testing.T) {
	p := baseProps()
	p.Variables = append(p.Variables, pc.PCVariable{Key: "grade", Type: pc.PCScalePoint, Label: "grade"})
	for i := range p.Data {
		p.Data[i].Values["grade"] = []string{"low", "mid", "high"}[i]
	}
	out := renderChart(t, p)
	if !strings.Contains(out, "grade") {
		t.Errorf("expected the point variable axis")
	}
}

func TestPC_VerticalLayout(t *testing.T) {
	p := baseProps()
	p.Layout = pc.PCLayoutVertical
	out := renderChart(t, p)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("vertical layout did not render")
	}
}

func TestPC_Golden(t *testing.T) {
	out := renderChart(t, baseProps())
	golden.Assert(t, "parallelcoordinates-basic", out)
}

func TestPC_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Four variables"
	p.Desc = "Three records across four variables."
	out := renderChart(t, p)
	if !strings.Contains(out, "<title>Four variables</title>") {
		t.Errorf("expected <title>")
	}
	if !strings.Contains(out, "<desc>Three records across four variables.</desc>") {
		t.Errorf("expected <desc>")
	}
}

func TestPC_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(renderChart(t, p), "data-tc-tooltip") {
		t.Errorf("interactive PC should emit data-tc-tooltip")
	}
	if strings.Contains(renderChart(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive PC must not emit data-tc-tooltip")
	}
}
