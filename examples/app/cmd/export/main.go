// Command export prerenders the examples/app demo to static HTML for GitHub
// Pages: every page in handlers.Pages and every chart detail page lands as a
// directory index (index.html, bar/index.html, chart/bar/index.html, …), so
// the clean extensionless URLs the live app uses resolve on a static host.
// Pages render through the real handlers via the production Mux in static
// mode (handlers.NewStaticApp: no registry mounts, no detail-page switchers,
// no htmx script), so the export cannot drift from the live app.
//
//	export [-base /templ-charts] [outdir]
//
// -base is the URL path prefix the site is served under (GitHub Pages
// project sites live at /<repo>/); every internal link is rendered with it.
//
// Client-side interactions survive the static host (hover tooltips, the line
// mesh/crosshair, canvas replay) because charts/interact runs entirely in the
// browser. Server-side ones (series toggle, bar/pie hover, hierarchy zoom,
// theme switching) need the live server; the layout carries a note pointing
// at it. The export fails if any page would still reference an htmx endpoint
// or, with -base set, a root-absolute link that skips the base.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geoffjay/templ-charts/examples/app/handlers"
	"github.com/geoffjay/templ-charts/examples/app/handlers/entries"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

func main() {
	base := flag.String("base", "", `URL path prefix the site is served under, e.g. "/templ-charts"`)
	flag.Parse()
	outDir := "site"
	if flag.NArg() > 0 {
		outDir = flag.Arg(0)
	}

	pages, details, err := export(*base, outDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("export: wrote %d pages + %d chart details to %s\n", pages, details, outDir)
}

// export prerenders the whole demo site into outDir with every internal link
// prefixed by base, returning the page and chart-detail counts.
func export(base, outDir string) (pages, details int, err error) {
	// Normalize to "" or "/prefix" (no trailing slash) so URL("/x") joins
	// cleanly.
	b := strings.TrimSuffix(base, "/")
	if b != "" && !strings.HasPrefix(b, "/") {
		b = "/" + b
	}
	templates.Base = b

	mux := handlers.NewStaticApp().Mux()

	for _, route := range handlers.Pages() {
		out := filepath.Join(outDir, strings.Trim(route, "/"), "index.html")
		if err := renderPage(mux, route, out, b); err != nil {
			return 0, 0, err
		}
	}

	slugs := make([]string, 0, len(entries.All()))
	for slug := range entries.All() {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	for _, slug := range slugs {
		out := filepath.Join(outDir, "chart", slug, "index.html")
		if err := renderPage(mux, "/chart/"+slug, out, b); err != nil {
			return 0, 0, err
		}
	}

	return len(handlers.Pages()), len(slugs), nil
}

// renderPage GETs one route through the mux, checks the page against the
// static-host invariants, and writes it to out. Any non-200 or failed check
// aborts the export, so a broken page never gets published.
func renderPage(mux http.Handler, route, out, base string) error {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, route, nil))
	if rec.Code != http.StatusOK {
		return fmt.Errorf("%s: status %d", route, rec.Code)
	}
	if err := checkPage(rec.Body.String(), base); err != nil {
		return fmt.Errorf("%s: %w", route, err)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, rec.Body.Bytes(), 0o644)
}

// checkPage enforces what a static host needs from a rendered page: no htmx
// request attributes (their /charts/ endpoints exist only on the live
// server), and — when base is set — every root-absolute link prefixed with
// it. External links (https://…) never start with href="/", so they pass.
func checkPage(body, base string) error {
	for _, attr := range []string{`hx-get="`, `hx-post="`} {
		if strings.Contains(body, attr) {
			return fmt.Errorf("emits an htmx request attribute (%s…), which needs the live server", attr)
		}
	}
	if base != "" && strings.Count(body, `href="/`) != strings.Count(body, `href="`+base+`/`) {
		return fmt.Errorf("links to a root-absolute path not prefixed with %s", base)
	}
	return nil
}
