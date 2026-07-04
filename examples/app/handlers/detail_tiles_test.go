package handlers_test

import (
	"net/http"
	"strings"
	"testing"
)

// TestHeatmapColorSpaceToggle verifies the detail page's color-space switcher:
// the heatmap page offers a "color space" control, and switching to Lab
// actually changes the rendered SVG (perceptual interpolation differs from RGB).
func TestHeatmapColorSpaceToggle(t *testing.T) {
	h, _ := newServer(t)

	base := do(t, h, http.MethodGet, "/chart/heatmap")
	if base.Code != http.StatusOK {
		t.Fatalf("/chart/heatmap: status %d", base.Code)
	}
	if !strings.Contains(base.Body.String(), "color space") {
		t.Error("/chart/heatmap: missing 'color space' switcher")
	}

	rgb := do(t, h, http.MethodGet, "/chart/heatmap?space=rgb").Body.String()
	lab := do(t, h, http.MethodGet, "/chart/heatmap?space=lab").Body.String()
	lch := do(t, h, http.MethodGet, "/chart/heatmap?space=lch").Body.String()
	for name, body := range map[string]string{"rgb": rgb, "lab": lab, "lch": lch} {
		if !strings.Contains(body, "<svg") {
			t.Errorf("/chart/heatmap?space=%s: no <svg", name)
		}
	}
	// Lab and Lch must differ from RGB (the whole point of the toggle).
	if lab == rgb {
		t.Error("heatmap Lab render identical to RGB — perceptual space not applied")
	}
	if lch == rgb {
		t.Error("heatmap Lch render identical to RGB — perceptual space not applied")
	}
}

// TestGeoFitRenders verifies the choropleth detail page renders with the
// fitExtent auto-fit projection (no manual ProjectionScale).
func TestGeoFitRenders(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/chart/geo")
	if rec.Code != http.StatusOK {
		t.Fatalf("/chart/geo: status %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<svg") {
		t.Error("/chart/geo: missing <svg (fit projection failed?)")
	}
	// The feature paths must be present (a degenerate fit would collapse them).
	if !strings.Contains(body, "<path") {
		t.Error("/chart/geo: missing feature <path> elements")
	}
}

// TestNonSpaceChartHasNoSpaceSwitcher confirms the color-space control only
// appears for charts that support it (bar has no SpaceRender).
func TestNonSpaceChartHasNoSpaceSwitcher(t *testing.T) {
	h, _ := newServer(t)
	body := do(t, h, http.MethodGet, "/chart/bar").Body.String()
	if strings.Contains(body, "color space") {
		t.Error("/chart/bar: unexpected 'color space' switcher (bar has no SpaceRender)")
	}
}
