package canvas

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/internal/golden"
)

// sampleList records one op of every kind, in the order a chart layer might
// emit them: set up state, draw a heatmap cell, a scatter point, an axis rule,
// a label, and a shape path. It exercises every branch of the encoder.
func sampleList() []Op {
	r := NewRecorder()
	r.Save()
	r.Translate(60, 20)
	r.FillStyle("#e8c1a0")
	r.GlobalAlpha(0.85)
	r.FillRect(10, 20, 30.5, 15.25)
	r.StrokeStyle("#333333")
	r.LineWidth(1.5)
	r.StrokeRect(10, 20, 30.5, 15.25)
	r.FillCircle(120.1239, 80.0004, 3) // rounds to 120.124, 80
	r.StrokeCircle(120.124, 80, 3)
	r.Line(0, 0, 200, 0)
	r.Font("11px sans-serif")
	r.TextAlign("middle")
	r.TextBaseline("top")
	r.FillText("A & B <ok>", 5, 12)
	r.FillPath("M0,0L10,10Z")
	r.StrokePath("M0,0L10,10")
	r.Restore()
	return r.Ops()
}

// TestRecorderRecordsInOrder asserts the Recorder preserves op order and kinds.
func TestRecorderRecordsInOrder(t *testing.T) {
	ops := sampleList()
	want := []OpKind{
		OpSave, OpTranslate, OpFillStyle, OpGlobalAlpha, OpFillRect,
		OpStrokeStyle, OpLineWidth, OpStrokeRect, OpFillCircle, OpStrokeCircle,
		OpLine, OpFont, OpTextAlign, OpTextBaseline, OpFillText,
		OpFillPath, OpStrokePath, OpRestore,
	}
	if len(ops) != len(want) {
		t.Fatalf("op count = %d, want %d", len(ops), len(want))
	}
	for i, k := range want {
		if ops[i].Kind != k {
			t.Errorf("op %d kind = %d, want %d", i, ops[i].Kind, k)
		}
	}
}

// TestEncodeJSONGolden pins the canonical JSON form of the full op set.
func TestEncodeJSONGolden(t *testing.T) {
	golden.Assert(t, "canvas-drawlist", EncodeJSON(sampleList()))
}

// TestEncodeEmpty: an empty draw-list is the empty JSON array.
func TestEncodeEmpty(t *testing.T) {
	if got := EncodeJSON(nil); got != "[]" {
		t.Errorf("EncodeJSON(nil) = %q, want %q", got, "[]")
	}
}

// TestEncodeRounding: coordinates round to 3 decimals, shortest form, and a
// rounded-to-zero value emits "0" (no "-0").
func TestEncodeRounding(t *testing.T) {
	r := NewRecorder()
	r.FillCircle(1.23456, 2.0, -0.0001) // -0.0001 -> 0
	got := EncodeJSON(r.Ops())
	want := "[\n[\"fc\",1.235,2,0]\n]"
	if got != want {
		t.Errorf("EncodeJSON = %q, want %q", got, want)
	}
}

// TestEncodeEscapesForScriptEmbedding: <, >, & and quotes in string args are
// escaped so the JSON is safe inside an HTML <script> element.
func TestEncodeEscapesForScriptEmbedding(t *testing.T) {
	r := NewRecorder()
	r.FillText(`x "q" <b>&`, 0, 0)
	got := EncodeJSON(r.Ops())
	for _, bad := range []string{"<b>", " & "} {
		if strings.Contains(got, bad) {
			t.Errorf("encoded JSON contains unescaped %q: %s", bad, got)
		}
	}
	if !strings.Contains(got, "\\u003c") || !strings.Contains(got, "\\u0026") {
		t.Errorf("expected escaped < and &, got %s", got)
	}
	if !strings.Contains(got, `\"q\"`) {
		t.Errorf("expected escaped quotes, got %s", got)
	}
}

// TestEncodeDeterministic: same list encodes identically every time.
func TestEncodeDeterministic(t *testing.T) {
	if EncodeJSON(sampleList()) != EncodeJSON(sampleList()) {
		t.Error("EncodeJSON not deterministic")
	}
}

// TestMarkupShape checks the canvas element and its paired ops script.
func TestMarkupShape(t *testing.T) {
	m := Markup("chart1", 300, 200, sampleList())
	for _, want := range []string{
		`<canvas class="tc-canvas" id="chart1"`,
		`data-tc-canvas="chart1-ops"`,
		`data-tc-w="300" data-tc-h="200"`,
		`width="300" height="200"`,
		// The canvas fills its wrapper (which carries the intrinsic size +
		// aspect-ratio), so it scales down with the container.
		`style="display:block;width:100%;height:100%"`,
		`<script type="application/json" class="tc-canvas-ops" id="chart1-ops">`,
		`</script>`,
	} {
		if !strings.Contains(m, want) {
			t.Errorf("Markup missing %q\n---\n%s", want, m)
		}
	}
	// The ops script must carry the encoded draw-list verbatim.
	if !strings.Contains(m, EncodeJSON(sampleList())) {
		t.Error("Markup does not embed the encoded draw-list")
	}
}

// TestCanvasScriptTag renders the replay module inside a <script> tag and checks
// it carries the op-dispatch switch and HiDPI scaling.
func TestCanvasScriptTag(t *testing.T) {
	var b strings.Builder
	if err := CanvasScriptTag().Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}
	out := b.String()
	if !strings.HasPrefix(out, "<script>") || !strings.HasSuffix(out, "</script>") {
		t.Errorf("CanvasScriptTag not wrapped in <script>: %q", out[:min(40, len(out))])
	}
	for _, want := range []string{"devicePixelRatio", "new Path2D(", "getContext('2d')", "case 'fc':"} {
		if !strings.Contains(out, want) {
			t.Errorf("script missing %q", want)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
