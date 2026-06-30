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
