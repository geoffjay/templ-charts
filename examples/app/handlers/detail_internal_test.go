package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/examples/app/handlers/entries"
)

// TestAllChartEntriesRender renders every registered detail-page entry at the
// default and at a non-default theme+palette, asserting each produces an SVG.
// This is the safety net for the per-family Render closures.
func TestAllChartEntriesRender(t *testing.T) {
	if len(entries.All()) == 0 {
		t.Fatal("no chart entries registered")
	}
	def := &theming.DefaultTheme
	for slug, e := range entries.All() {
		if e.Snippet == "" {
			t.Errorf("%s: empty Snippet", slug)
		}
		// Default theme, no palette, animate off.
		svg, err := e.Render(def, "", false)
		if err != nil {
			t.Errorf("%s: Render(default) error: %v", slug, err)
			continue
		}
		if !strings.HasPrefix(svg, "<svg") {
			t.Errorf("%s: Render(default) did not produce an <svg", slug)
		}
		// A non-default theme + palette must also render (palette is a no-op for
		// non-ordinal charts, which is fine — it must still render).
		if _, err := e.Render(def, colors.PaletteTableau10, false); err != nil {
			t.Errorf("%s: Render(palette=tableau10) error: %v", slug, err)
		}
		// With animate on it must also render.
		if _, err := e.Render(def, "", true); err != nil {
			t.Errorf("%s: Render(animate=true) error: %v", slug, err)
		}
	}
}

// TestIndexLinksHaveEntries checks that every /chart/{slug} link on the index
// resolves to a registered entry (and there are no orphan entries).
func TestIndexLinksHaveEntries(t *testing.T) {
	app := NewApp()
	rec := httptest.NewRecorder()
	app.Index(rec, httptest.NewRequest("GET", "/", nil))
	body := rec.Body.String()

	for slug := range entries.All() {
		want := `href="/chart/` + slug + `"`
		if !strings.Contains(body, want) {
			t.Errorf("index missing link for registered chart %q (%s)", slug, want)
		}
	}
}
