package static_test

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/static"
)

// TestRegistryCoversAllFamilies asserts the registry was extended from the
// original 3 families to the full catalog.
func TestRegistryCoversAllFamilies(t *testing.T) {
	if len(static.ChartsMapping) < 28 {
		t.Errorf("ChartsMapping has %d families, want >= 28", len(static.ChartsMapping))
	}
	if len(static.Samples) < 28 {
		t.Errorf("Samples has %d families, want >= 28", len(static.Samples))
	}
}

// TestRenderEveryRegisteredChart renders each registered family through the
// static dispatcher using its bundled sample, asserting a real SVG comes out —
// this is the per-family render test the plan calls for, and it doubles as the
// "every sample is renderable" guarantee.
func TestRenderEveryRegisteredChart(t *testing.T) {
	for ct, sample := range static.Samples {
		if _, ok := static.ChartsMapping[ct]; !ok {
			t.Errorf("%q: sample present but no ChartsMapping entry", ct)
			continue
		}
		out, err := static.RenderChart(ct, sample.Props, nil)
		if err != nil {
			t.Errorf("%q: RenderChart error: %v", ct, err)
			continue
		}
		if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
			preview := out
			if len(preview) > 80 {
				preview = preview[:80]
			}
			t.Errorf("%q: expected an <svg>…</svg>, got %q", ct, preview)
		}
	}
	// Every mapping should have a matching sample.
	for ct := range static.ChartsMapping {
		if _, ok := static.Samples[ct]; !ok {
			t.Errorf("%q: mapping present but no Sample", ct)
		}
	}
}
