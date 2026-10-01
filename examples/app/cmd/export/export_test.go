package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/examples/app/handlers"
	"github.com/geoffjay/templ-charts/examples/app/handlers/entries"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

// TestExport runs the full static export the way the Pages workflow does and
// checks what a static host needs. export itself fails if any page emits an
// htmx endpoint or an unprefixed root-absolute link, so a nil error already
// covers those invariants for every page.
func TestExport(t *testing.T) {
	t.Cleanup(func() { templates.Base = "" })
	dir := t.TempDir()

	pages, details, err := export("/templ-charts/", dir)
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if pages != len(handlers.Pages()) || details != len(entries.All()) {
		t.Errorf("wrote %d pages + %d details, want %d + %d", pages, details, len(handlers.Pages()), len(entries.All()))
	}

	read := func(rel string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		return string(b)
	}
	for _, c := range []struct{ file, want, why string }{
		// Each route lands in its own directory index (clean URLs), not one
		// file overwritten per route.
		{"index.html", `href="/templ-charts/chart/bar"`, "index links carry the normalized base"},
		{"bar/index.html", `id="chart-bar-stacked"`, "bar page has its own directory index"},
		{"chart/bar/index.html", `tc-detail-chart`, "detail page has its own directory index"},
		// Client-side hover survives the static host. Match the attribute
		// with its value: the inlined interact script names both attributes
		// as selectors on every page.
		{"line/index.html", `data-tc-mesh="`, "line mesh hover payload"},
		{"heatmap/index.html", `data-tc-tooltip="`, "heatmap cell tooltips"},
	} {
		if !strings.Contains(read(c.file), c.want) {
			t.Errorf("%s: missing %q (%s)", c.file, c.want, c.why)
		}
	}
	// Detail-page switchers navigate to query-param URLs only the live
	// server renders, so the export omits them.
	if strings.Contains(read("chart/heatmap/index.html"), `class="tc-switch"`) {
		t.Error("chart/heatmap/index.html: static export should omit the switchers")
	}
}

func TestCheckPage(t *testing.T) {
	for _, c := range []struct {
		name, body, base string
		wantErr          bool
	}{
		{"prefixed and external links", `<a href="/templ-charts/">h</a><a href="/templ-charts/bar">b</a><a href="https://example.com/">x</a>`, "/templ-charts", false},
		{"unprefixed root-absolute link", `<a href="/templ-charts/">h</a><a href="/bar">b</a>`, "/templ-charts", true},
		{"base prefix is not a path prefix", `<a href="/templ-chartsfoo">b</a>`, "/templ-charts", true},
		{"root-absolute links without a base", `<a href="/bar">b</a>`, "", false},
		{"htmx get endpoint", `<rect hx-get="/charts/x/hover"></rect>`, "", true},
		{"htmx post endpoint", `<g hx-post="/charts/x/toggle"></g>`, "/templ-charts", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			if err := checkPage(c.body, c.base); (err != nil) != c.wantErr {
				t.Errorf("checkPage error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
