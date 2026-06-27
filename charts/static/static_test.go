package static_test

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/static"
)

func TestRenderChart_Bar(t *testing.T) {
	sample, ok := static.Samples[static.ChartTypeBar]
	if !ok {
		t.Fatal("missing bar sample")
	}
	out, err := static.RenderChart(static.ChartTypeBar, sample.Props, nil)
	if err != nil {
		t.Fatalf("RenderChart bar: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>, got: %q", out[:minLen(out)])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("expected closing </svg>")
	}
}

func TestRenderChart_Line(t *testing.T) {
	sample, ok := static.Samples[static.ChartTypeLine]
	if !ok {
		t.Fatal("missing line sample")
	}
	out, err := static.RenderChart(static.ChartTypeLine, sample.Props, nil)
	if err != nil {
		t.Fatalf("RenderChart line: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>, got: %q", out[:minLen(out)])
	}
}

func TestRenderChart_Pie(t *testing.T) {
	sample, ok := static.Samples[static.ChartTypePie]
	if !ok {
		t.Fatal("missing pie sample")
	}
	out, err := static.RenderChart(static.ChartTypePie, sample.Props, nil)
	if err != nil {
		t.Fatalf("RenderChart pie: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>, got: %q", out[:minLen(out)])
	}
}

func TestRenderChart_UnknownType(t *testing.T) {
	_, err := static.RenderChart("unknown", nil, nil)
	if err == nil {
		t.Error("expected error for unknown chart type")
	}
}

func TestRenderChart_StaticOverridesApplied(t *testing.T) {
	// Static render should have animate=false → no <animate> elements.
	sample := static.Samples[static.ChartTypeBar]
	out, _ := static.RenderChart(static.ChartTypeBar, sample.Props, nil)
	if strings.Contains(out, "<animate") {
		t.Error("expected no <animate> in static render (animate=false)")
	}
}

func TestRenderChart_OverrideWidth(t *testing.T) {
	sample := static.Samples[static.ChartTypeBar]
	override := map[string]any{"width": float64(300)}
	out, err := static.RenderChart(static.ChartTypeBar, sample.Props, override)
	if err != nil {
		t.Fatalf("RenderChart with override: %v", err)
	}
	// The SVG width should reflect the override.
	if !strings.Contains(out, `width="300"`) {
		t.Errorf("expected width=300 in output from override")
	}
}

func TestRenderChart_OverrideNotWhitelisted(t *testing.T) {
	sample := static.Samples[static.ChartTypeBar]
	// borderWidth is NOT in bar's runtimeProps, so this override should be ignored.
	override := map[string]any{"borderWidth": float64(999)}
	out, err := static.RenderChart(static.ChartTypeBar, sample.Props, override)
	if err != nil {
		t.Fatalf("RenderChart: %v", err)
	}
	if strings.Contains(out, `stroke-width="999"`) {
		t.Error("expected non-whitelisted borderWidth override to be ignored")
	}
}

func TestChartsMapping_AllPresent(t *testing.T) {
	for _, ct := range []static.ChartType{static.ChartTypeBar, static.ChartTypeLine, static.ChartTypePie} {
		if _, ok := static.ChartsMapping[ct]; !ok {
			t.Errorf("ChartsMapping missing %q", ct)
		}
	}
}

func TestSamples_AllPresent(t *testing.T) {
	for _, ct := range []static.ChartType{static.ChartTypeBar, static.ChartTypeLine, static.ChartTypePie} {
		if _, ok := static.Samples[ct]; !ok {
			t.Errorf("Samples missing %q", ct)
		}
	}
}

func TestRenderChart_BarPropsDirectly(t *testing.T) {
	props := bar.BarProps{
		Width: 400, Height: 200,
		Data: []bar.BarDatum{
			{"id": "a", "value": float64(10)},
			{"id": "b", "value": float64(20)},
		},
	}
	out, err := static.RenderChart(static.ChartTypeBar, props, nil)
	if err != nil {
		t.Fatalf("RenderChart: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>")
	}
}

func TestRenderChart_LinePropsDirectly(t *testing.T) {
	props := line.LineProps{
		Width: 400, Height: 200,
		Data: []line.LineSeries{
			{ID: "A", Data: []line.LinePointData{{X: "x", Y: float64(1)}, {X: "y", Y: float64(2)}}},
		},
	}
	out, err := static.RenderChart(static.ChartTypeLine, props, nil)
	if err != nil {
		t.Fatalf("RenderChart: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>")
	}
}

func TestRenderChart_PiePropsDirectly(t *testing.T) {
	props := pie.PieProps{
		Width: 400, Height: 400,
		Data: []any{
			map[string]any{"id": "A", "value": float64(10)},
			map[string]any{"id": "B", "value": float64(20)},
		},
	}
	out, err := static.RenderChart(static.ChartTypePie, props, nil)
	if err != nil {
		t.Fatalf("RenderChart: %v", err)
	}
	if !strings.HasPrefix(out, "<svg") {
		t.Errorf("expected <svg…>")
	}
}

func TestRenderChart_WrongPropsType(t *testing.T) {
	// Pass line props to bar — should error.
	_, err := static.RenderChart(static.ChartTypeBar, line.LineProps{}, nil)
	if err == nil {
		t.Error("expected error for wrong props type")
	}
}

func minLen(s string) int {
	if len(s) < 80 {
		return len(s)
	}
	return 80
}
