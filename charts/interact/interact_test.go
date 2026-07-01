package interact_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/interact"
)

func TestScript_NonEmpty(t *testing.T) {
	if len(interact.Script) < 100 {
		t.Fatalf("Script looks too short (%d bytes)", len(interact.Script))
	}
	for _, want := range []string{"data-tc-tooltip", "data-tc-mesh", "tc-chart", "tc-crosshair", "getScreenCTM"} {
		if !strings.Contains(interact.Script, want) {
			t.Errorf("Script missing %q", want)
		}
	}
}

// TestScript_LeavesHtmxChartsAlone guards the regression where the client
// script's fall-through hide() clobbered the htmx-populated tooltip on bar/pie
// charts: the script must gate on isClientChart so charts with no data-tc-*
// elements are never touched.
func TestScript_LeavesHtmxChartsAlone(t *testing.T) {
	if !strings.Contains(interact.Script, "isClientChart") {
		t.Errorf("Script must gate on isClientChart so it does not clobber htmx (bar/pie) tooltips")
	}
}

func TestScriptTag_WrapsScript(t *testing.T) {
	var b strings.Builder
	if err := interact.ScriptTag().Render(context.Background(), &b); err != nil {
		t.Fatalf("ScriptTag.Render: %v", err)
	}
	out := b.String()
	if !strings.HasPrefix(out, "<script>") || !strings.HasSuffix(out, "</script>") {
		t.Errorf("ScriptTag should wrap Script in <script>…</script>, got prefix/suffix %q…%q", out[:9], out[len(out)-10:])
	}
}

func TestTooltipHTML_LabelAndValue(t *testing.T) {
	out := interact.TooltipHTML("#ff0000", "France", "42")
	if !strings.Contains(out, "France: 42") {
		t.Errorf("expected 'France: 42' in %q", out)
	}
	if !strings.Contains(out, "background:#ff0000") {
		t.Errorf("expected chip color in %q", out)
	}
}

func TestTooltipHTML_EscapesText(t *testing.T) {
	out := interact.TooltipHTML("", "a<b>", "x&y")
	if strings.Contains(out, "a<b>") {
		t.Errorf("label should be HTML-escaped: %q", out)
	}
	if !strings.Contains(out, "a&lt;b&gt;") || !strings.Contains(out, "x&amp;y") {
		t.Errorf("expected escaped label/value in %q", out)
	}
}

func TestMeshData(t *testing.T) {
	out := interact.MeshData([]interact.MeshPoint{{X: 1, Y: 2, HTML: "a"}, {X: 3, Y: 4, HTML: "b"}})
	if !strings.Contains(out, `"x":1`) || !strings.Contains(out, `"html":"b"`) {
		t.Errorf("MeshData JSON unexpected: %q", out)
	}
	if got := interact.MeshData(nil); got != "[]" {
		t.Errorf("empty MeshData = %q, want []", got)
	}
}

func TestMeshCellsPath(t *testing.T) {
	pts := []interact.MeshPoint{{X: 20, Y: 20}, {X: 80, Y: 30}, {X: 50, Y: 70}, {X: 30, Y: 50}}
	got := interact.MeshCellsPath(pts, 100, 100)
	if !strings.HasPrefix(got, "M") || !strings.Contains(got, "Z") {
		t.Errorf("MeshCellsPath should be a closed path, got %q", got)
	}
	if interact.MeshCellsPath(nil, 100, 100) != "" {
		t.Errorf("empty points should yield empty path")
	}
}

func TestMeshOverlay(t *testing.T) {
	pts := []interact.MeshPoint{{X: 20, Y: 20, HTML: "a"}, {X: 80, Y: 30, HTML: "b"}, {X: 50, Y: 70, HTML: "c"}}
	out := interact.MeshOverlay(pts, 100, 100, false, 0)
	if !strings.Contains(out, "data-tc-mesh=") || !strings.Contains(out, `fill-opacity="0"`) {
		t.Errorf("overlay should carry a transparent capture rect with data-tc-mesh: %q", out)
	}
	if strings.Contains(out, "data-tc-mesh-radius") {
		t.Errorf("no radius attr expected when detectionRadius<=0")
	}
	// With debug + radius.
	dbg := interact.MeshOverlay(pts, 100, 100, true, 30)
	if !strings.Contains(dbg, "stroke=\"red\"") {
		t.Errorf("debug overlay should draw the voronoi cells")
	}
	if !strings.Contains(dbg, `data-tc-mesh-radius="30"`) {
		t.Errorf("expected detection-radius attribute")
	}
	if interact.MeshOverlay(nil, 100, 100, true, 10) != "" {
		t.Errorf("empty points should yield empty overlay")
	}
}

func TestScript_DetectionRadius(t *testing.T) {
	if !strings.Contains(interact.Script, "data-tc-mesh-radius") {
		t.Errorf("Script must honor the detection-radius attribute")
	}
}

func TestEscapeAttr(t *testing.T) {
	in := `<span style="x">a&b</span>`
	out := interact.EscapeAttr(in)
	for _, bad := range []string{`"`, "<", ">"} {
		if strings.Contains(out, bad) {
			t.Errorf("EscapeAttr left raw %q in %q", bad, out)
		}
	}
	if !strings.Contains(out, "&quot;") || !strings.Contains(out, "&lt;") {
		t.Errorf("EscapeAttr should encode quotes/brackets: %q", out)
	}
}
