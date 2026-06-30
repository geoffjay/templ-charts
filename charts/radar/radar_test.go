package radar_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/radar"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// sampleData is a deterministic 5-index × 3-key dataset (taste profile).
func sampleData() []map[string]any {
	return []map[string]any{
		{"taste": "fruity", "chardonnay": 93.0, "carmenere": 61.0, "syrah": 114.0},
		{"taste": "bitter", "chardonnay": 91.0, "carmenere": 37.0, "syrah": 72.0},
		{"taste": "heavy", "chardonnay": 56.0, "carmenere": 95.0, "syrah": 99.0},
		{"taste": "strong", "chardonnay": 64.0, "carmenere": 90.0, "syrah": 30.0},
		{"taste": "sunny", "chardonnay": 119.0, "carmenere": 94.0, "syrah": 103.0},
	}
}

func baseProps() radar.RadarProps {
	return radar.RadarProps{
		Width: 500, Height: 460,
		Margin:  core.Margin{Top: 70, Right: 80, Bottom: 40, Left: 80},
		Data:    sampleData(),
		Keys:    []string{"chardonnay", "carmenere", "syrah"},
		IndexBy: "taste",
	}
}

func render(t *testing.T, props radar.RadarProps) string {
	t.Helper()
	var b strings.Builder
	if err := radar.Radar(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Radar.Render: %v", err)
	}
	return b.String()
}

func TestRadar_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestRadar_PolygonPerKey(t *testing.T) {
	// One closed <path> per key (3) in the layers group. Each path is a polygon
	// fill, so count <path with fill-opacity.
	out := render(t, baseProps())
	if got := strings.Count(out, "fill-opacity"); got != 3 {
		t.Errorf("polygon count via fill-opacity = %d, want 3", got)
	}
}

func TestRadar_DotCount(t *testing.T) {
	// 5 indices × 3 keys = 15 dots; each dot is a <circle>. The grid also emits
	// 5 concentric level circles (default circular shape). 15 + 5 = 20 circles.
	out := render(t, baseProps())
	if got := strings.Count(out, "<circle"); got != 20 {
		t.Errorf("circle count = %d, want 20 (15 dots + 5 grid levels)", got)
	}
}

func TestRadar_LinearGridUsesPolygons(t *testing.T) {
	p := baseProps()
	p.GridShape = radar.GridShapeLinear
	out := render(t, p)
	// No grid level circles now → only the 15 dots remain as circles.
	if got := strings.Count(out, "<circle"); got != 15 {
		t.Errorf("circle count = %d, want 15 (dots only; linear grid uses paths)", got)
	}
}

func TestRadar_DotsDisabled(t *testing.T) {
	p := baseProps()
	p.EnableDots = radar.BoolPtr(false)
	out := render(t, p)
	// Only the 5 grid level circles remain.
	if got := strings.Count(out, "<circle"); got != 5 {
		t.Errorf("circle count = %d, want 5 (grid levels only)", got)
	}
}

func TestRadar_IndexLabels(t *testing.T) {
	out := render(t, baseProps())
	for _, idx := range []string{"fruity", "bitter", "heavy", "strong", "sunny"} {
		if !strings.Contains(out, ">"+idx+"<") {
			t.Errorf("missing index label %q", idx)
		}
	}
}

func TestRadar_Golden(t *testing.T) {
	out := render(t, baseProps())
	golden.Assert(t, "radar-basic", out)
}

func TestRadar_Golden_Linear(t *testing.T) {
	p := baseProps()
	p.GridShape = radar.GridShapeLinear
	p.Rotation = 30
	out := render(t, p)
	golden.Assert(t, "radar-linear-rotated", out)
}
