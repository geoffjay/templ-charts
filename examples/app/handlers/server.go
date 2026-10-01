package handlers

import (
	"net/http"
)

// pageRoutes is the demo's page route table: one entry per HTML page. Mux
// serves each; the static export (examples/app/cmd/export) prerenders each
// to <path>/index.html. The htmx endpoints (/charts/, /scale) and the
// parameterized detail route (/chart/{slug}) are wired separately in Mux —
// they are not pages, and the export renders the detail pages from the
// entries registry instead.
var pageRoutes = []struct {
	path string
	h    func(*App, http.ResponseWriter, *http.Request)
}{
	{"/", (*App).Index},
	{"/bar", (*App).Bar},
	{"/line", (*App).Line},
	{"/pie", (*App).Pie},
	{"/heatmap", (*App).Heatmap},
	{"/waffle", (*App).Waffle},
	{"/calendar", (*App).Calendar},
	{"/radar", (*App).Radar},
	{"/radial-bar", (*App).RadialBar},
	{"/scatterplot", (*App).ScatterPlot},
	{"/stream", (*App).Stream},
	{"/bullet", (*App).Bullet},
	{"/funnel", (*App).Funnel},
	{"/boxplot", (*App).BoxPlot},
	{"/bump", (*App).Bump},
	{"/marimekko", (*App).Marimekko},
	{"/parallel-coordinates", (*App).ParallelCoordinates},
	{"/polar-bar", (*App).PolarBar},
	{"/treemap", (*App).Treemap},
	{"/sunburst", (*App).Sunburst},
	{"/icicle", (*App).Icicle},
	{"/circle-packing", (*App).CirclePacking},
	{"/tree", (*App).Tree},
	{"/voronoi", (*App).Voronoi},
	{"/network", (*App).Network},
	{"/swarmplot", (*App).SwarmPlot},
	{"/sankey", (*App).Sankey},
	{"/chord", (*App).Chord},
	{"/geo", (*App).Geo},
	{"/scales", (*App).Scales},
	{"/styling", (*App).Styling},
	{"/legends", (*App).Legends},
	{"/composition", (*App).Composition},
	{"/dashboard", (*App).Dashboard},
	{"/palettes", (*App).Palettes},
	{"/themes", (*App).Themes},
	{"/benchmark", (*App).Benchmark},
}

// Pages returns the URL paths of the demo's HTML pages (every route in the
// table above). The static export prerenders one file per path.
func Pages() []string {
	out := make([]string, len(pageRoutes))
	for i, p := range pageRoutes {
		out[i] = p.path
	}
	return out
}

// Mux returns the demo server's full route table: the htmx chart endpoints
// under /charts/, the per-chart detail pages under /chart/, the /scale
// fragment endpoint, and one route per page in Pages. main.go passes the
// result to http.ListenAndServe; the export renders each page through it.
func (a *App) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/charts/", a.handler)
	mux.HandleFunc("/chart/", a.Detail)
	mux.HandleFunc("/scale", a.Scale)
	for _, p := range pageRoutes {
		mux.HandleFunc(p.path, func(w http.ResponseWriter, r *http.Request) {
			p.h(a, w, r)
		})
	}
	return mux
}

// NewStaticApp returns an App configured for the static export: pages render
// without the server-only interactivity — no registry mounts (so no hx-*
// hover/toggle/zoom wiring), no detail-page switchers, no htmx script — and
// the benchmark page notes that its timings were measured at export time.
// The client-side hover layer (charts/interact tooltips, mesh, canvas
// replay) keeps working, since it runs entirely in the browser.
// examples/app/cmd/export prerenders each page through Mux in this mode.
func NewStaticApp() *App {
	a := NewApp()
	a.static = true
	return a
}
