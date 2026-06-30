package funnel_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/funnel"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleData() []funnel.FunnelDatum {
	return []funnel.FunnelDatum{
		{ID: "step_sent", Label: "Sent", Value: 60000},
		{ID: "step_viewed", Label: "Viewed", Value: 38000},
		{ID: "step_clicked", Label: "Clicked", Value: 22000},
		{ID: "step_add", Label: "Add to Cart", Value: 12000},
		{ID: "step_purchased", Label: "Purchased", Value: 7000},
	}
}

func baseProps() funnel.FunnelProps {
	return funnel.FunnelProps{
		Width: 500, Height: 500,
		Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
		Data:   sampleData(),
	}
}

func render(t *testing.T, props funnel.FunnelProps) string {
	t.Helper()
	var b strings.Builder
	if err := funnel.Funnel(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Funnel.Render: %v", err)
	}
	return b.String()
}

func TestFunnel_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestFunnel_PartCount(t *testing.T) {
	// 5 parts → 5 filled <path> + 2 border paths each (10) = 15 area/border
	// paths. Assert at least 5 fill-opacity paths (one per part).
	out := render(t, baseProps())
	if got := strings.Count(out, "fill-opacity"); got != 5 {
		t.Errorf("part fill count = %d, want 5", got)
	}
}

func TestFunnel_LabelsPresent(t *testing.T) {
	out := render(t, baseProps())
	for _, l := range []string{"Sent", "Viewed", "Purchased"} {
		if !strings.Contains(out, ">"+l+"<") {
			t.Errorf("missing label %q", l)
		}
	}
}

func TestFunnel_LinearInterpolation(t *testing.T) {
	p := baseProps()
	p.Interpolation = funnel.FunnelInterpolationLinear
	out := render(t, p)
	if got := strings.Count(out, "fill-opacity"); got != 5 {
		t.Errorf("linear: part fill count = %d, want 5", got)
	}
}

func TestFunnel_Golden(t *testing.T) {
	out := render(t, baseProps())
	golden.Assert(t, "funnel-smooth", out)
}

func TestFunnel_Golden_Horizontal(t *testing.T) {
	p := baseProps()
	p.Width, p.Height = 700, 300
	p.Direction = funnel.FunnelDirectionHorizontal
	out := render(t, p)
	golden.Assert(t, "funnel-horizontal", out)
}

func TestFunnel_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(render(t, p), "data-tc-tooltip") {
		t.Errorf("interactive funnel should emit data-tc-tooltip")
	}
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive funnel must not emit data-tc-tooltip")
	}
}
