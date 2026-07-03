package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/examples/app/handlers"
)

// do runs a request against the app's mux-less handler set. We build a tiny
// mux mirroring main.go so the page handlers can be exercised.
func newServer(t *testing.T) (http.Handler, *handlers.App) {
	t.Helper()
	app := handlers.NewApp()
	mux := http.NewServeMux()
	mux.Handle("/charts/", app.Handler())
	mux.HandleFunc("/", app.Index)
	mux.HandleFunc("/bar", app.Bar)
	mux.HandleFunc("/line", app.Line)
	mux.HandleFunc("/pie", app.Pie)
	mux.HandleFunc("/themes", app.Themes)
	mux.HandleFunc("/benchmark", app.Benchmark)
	mux.HandleFunc("/chart/", app.Detail)
	return mux, app
}

func TestDetailPage(t *testing.T) {
	h, _ := newServer(t)
	// Base render + theme + palette query params should all render an SVG.
	for _, target := range []string{
		"/chart/bar",
		"/chart/line?theme=dark",
		"/chart/pie?theme=custom&palette=tableau10",
	} {
		rec := do(t, h, http.MethodGet, target)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200", target, rec.Code)
			continue
		}
		body := rec.Body.String()
		if !strings.Contains(body, "<svg") {
			t.Errorf("%s: body missing <svg", target)
		}
		if !strings.Contains(body, "tc-snippet") {
			t.Errorf("%s: body missing code snippet panel", target)
		}
	}
	// Unknown slug → 404.
	if rec := do(t, h, http.MethodGet, "/chart/nope"); rec.Code != http.StatusNotFound {
		t.Errorf("/chart/nope: status = %d, want 404", rec.Code)
	}

	// The animate toggle must flip the chart's SMIL enter animation: with
	// ?animate=1 the SVG contains <animate elements; the default render (no
	// param) must not.
	animOn := do(t, h, http.MethodGet, "/chart/heatmap?animate=1").Body.String()
	if !strings.Contains(animOn, "<animate") {
		t.Error("/chart/heatmap?animate=1: body missing <animate element")
	}
	animOff := do(t, h, http.MethodGet, "/chart/heatmap").Body.String()
	if strings.Contains(animOff, "<animate") {
		t.Error("/chart/heatmap (default): body should not contain <animate")
	}
}

func TestBenchmarkPage(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/benchmark")
	if rec.Code != http.StatusOK {
		t.Fatalf("/benchmark: status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<table") {
		t.Error("/benchmark: body missing results <table")
	}
	if !strings.Contains(body, "<svg") {
		t.Error("/benchmark: body missing showcase <svg")
	}
}

func do(t *testing.T, h http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPages_ReturnHTML(t *testing.T) {
	h, _ := newServer(t)
	// "/" is the index (link list, no charts); the rest render SVG inline.
	svgRoutes := []string{"/bar", "/line", "/pie", "/themes"}
	for _, route := range []string{"/", "/bar", "/line", "/pie", "/themes"} {
		rec := do(t, h, http.MethodGet, route)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status = %d, want 200; body=%s", route, rec.Code, rec.Body.String()[:200])
			continue
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
			t.Errorf("%s: Content-Type = %q, want text/html", route, ct)
		}
	}
	for _, route := range svgRoutes {
		rec := do(t, h, http.MethodGet, route)
		if !strings.Contains(rec.Body.String(), "<svg") {
			t.Errorf("%s: body missing <svg (chart not rendered inline)", route)
		}
	}
}

func TestBarPage_HasFiveCards(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/bar")
	if got := strings.Count(rec.Body.String(), `class="card"`); got != 5 {
		t.Errorf("/bar card count = %d, want 5", got)
	}
}

func TestThemesPage_HasNineCards(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/themes")
	if got := strings.Count(rec.Body.String(), `class="card"`); got != 9 {
		t.Errorf("/themes card count = %d, want 9", got)
	}
}

func TestHtmxEndpoints_PreRegistered(t *testing.T) {
	h, _ := newServer(t)
	// Hover a bar without first visiting /bar — instances are pre-registered
	// in NewApp, so this should return a tooltip fragment, not 404.
	rec := do(t, h, http.MethodGet, "/charts/bar-legend-toggle/hover?bar=hot%20dogs.USA")
	if rec.Code != http.StatusOK {
		t.Fatalf("hover status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "nivo-tooltip-basic") {
		t.Errorf("hover body missing tooltip class; got %q", rec.Body.String())
	}
}

func TestHtmxEndpoints_ToggleReRendersSVG(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodPost, "/charts/bar-legend-toggle/toggle?series=hot%20dogs")
	if rec.Code != http.StatusOK {
		t.Fatalf("toggle status = %d, want 200", rec.Code)
	}
	if !strings.HasPrefix(rec.Body.String(), "<svg") {
		t.Errorf("toggle body should be full SVG, got %q", rec.Body.String()[:40])
	}
}

func TestHtmxEndpoints_PieHover(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/charts/pie-active/hover?arc=Go")
	if rec.Code != http.StatusOK {
		t.Fatalf("pie hover status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "nivo-tooltip-basic") {
		t.Errorf("pie hover body missing tooltip class; got %q", rec.Body.String())
	}
}

func TestHtmxEndpoints_UnknownInstance(t *testing.T) {
	h, _ := newServer(t)
	rec := do(t, h, http.MethodGet, "/charts/nope")
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown instance status = %d, want 404", rec.Code)
	}
}
