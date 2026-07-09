// Package handlers implements the demo app's HTTP handlers. It wires the
// htmx.Handler for the interactive chart endpoints and renders the page
// handlers (/bar, /line, /pie, /themes) using the demo definitions in the
// demos package.
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/chord"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/funnel"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/marimekko"
	"github.com/geoffjay/templ-charts/charts/network"
	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/polarbar"
	"github.com/geoffjay/templ-charts/charts/radar"
	"github.com/geoffjay/templ-charts/charts/radialbar"
	"github.com/geoffjay/templ-charts/charts/sankey"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/stream"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
	"github.com/geoffjay/templ-charts/charts/tree"
	"github.com/geoffjay/templ-charts/charts/treemap"
	"github.com/geoffjay/templ-charts/charts/voronoi"
	"github.com/geoffjay/templ-charts/charts/waffle"
	"github.com/geoffjay/templ-charts/examples/app/demos"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

// App holds the shared state for the demo handlers: the htmx registry and
// handler. The page handlers register their demos with the registry on first
// request (idempotent) and render the initial SVG.
type App struct {
	registry *htmx.Registry
	handler  *htmx.Handler
}

// NewApp returns an App with a fresh registry + handler and pre-registers all
// interactive demos (bar/line/pie) so the htmx endpoints work even if a
// client hits /charts/<id>/... before visiting the page. The themes page is
// static and not registered.
func NewApp() *App {
	r := htmx.NewRegistry()
	app := &App{
		registry: r,
		handler:  htmx.NewHandler(r),
	}
	for _, d := range demos.BarDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.LineDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.PieDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.HeatmapDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.StylingDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.LegendsDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	for _, d := range demos.DashboardDemos() {
		r.Register(d.ID, d.Kind, d.Props)
	}
	return app
}

// Registry returns the underlying htmx registry (used by main.go to mount
// the htmx handler).
func (a *App) Registry() *htmx.Registry { return a.registry }

// Handler returns the htmx handler (mounted at /charts/ by main.go).
func (a *App) Handler() *htmx.Handler { return a.handler }

// --- page handlers ---

// Index handles GET /: the home page listing all demo pages.
func (a *App) Index(w http.ResponseWriter, r *http.Request) {
	props := templates.IndexPageProps{
		Pages: []templates.IndexLink{
			{Href: "/chart/bar", Title: "Bar charts", Description: "Stacked, grouped, markers + annotations, legend toggle, totals."},
			{Href: "/chart/line", Title: "Line charts", Description: "Single & multi-series, area + points, slices, mesh hover."},
			{Href: "/chart/pie", Title: "Pie charts", Description: "Plain, donut, half, sorted, active-arc hover, legend toggle."},
			{Href: "/chart/heatmap", Title: "Heatmap", Description: "2D value grid: sequential/diverging color scales, labels, borders, continuous legend."},
			{Href: "/chart/waffle", Title: "Waffle", Description: "Part-of-whole cell grid: fill direction, borders, legend (built on charts/grid)."},
			{Href: "/chart/calendar", Title: "Calendar", Description: "Day-grid heatmap over a date range: quantized colors, month/year legends, horizontal/vertical."},
			{Href: "/chart/radar", Title: "Radar", Description: "Polar line/area chart: values per key around shared indices, circular/polygon grids, dots, legend."},
			{Href: "/chart/radial-bar", Title: "Radial bar", Description: "Stacked bars as polar arcs (charts/polar-axes + charts/arcs): tracks, radial/circular axes, labels."},
			{Href: "/chart/scatterplot", Title: "Scatterplot", Description: "{x,y} nodes on linear/time scales: grid, axes, per-series colors, legend."},
			{Href: "/chart/stream", Title: "Stream", Description: "Stacked areas with wiggle/silhouette/expand offsets and a smooth curve."},
			{Href: "/chart/bullet", Title: "Bullet", Description: "KPI ranges + measure bars + markers on a shared value scale, per-row axis."},
			{Href: "/chart/funnel", Title: "Funnel", Description: "Ordered parts as smooth/linear trapezoids with separators and labels."},
			{Href: "/chart/boxplot", Title: "Box plot", Description: "Quantile box + whisker glyphs from raw observations (d3/array.Quantile)."},
			{Href: "/chart/bump", Title: "Bump", Description: "Ranking over time: smooth (curveBumpX) or linear lines with end labels and point hover."},
			{Href: "/chart/marimekko", Title: "Marimekko", Description: "Variable-width stacked bars: width by value, segments stacked via d3.Stack."},
			{Href: "/chart/parallel-coordinates", Title: "Parallel coordinates", Description: "One axis per variable; each record a polyline across linear/point scales."},
			{Href: "/chart/polar-bar", Title: "Polar bar", Description: "Stacked bars wrapped into a full circle: angle band per index, radius-stacked keys."},
			{Href: "/chart/treemap", Title: "Treemap", Description: "Nested rectangles (d3-hierarchy): squarify/binary tiling, leaf + parent labels."},
			{Href: "/chart/sunburst", Title: "Sunburst", Description: "Radial partition (d3-hierarchy): arcs by value, colors inherited down the tree."},
			{Href: "/chart/icicle", Title: "Icicle", Description: "Depth-banded partition rectangles (d3-hierarchy), oriented four ways."},
			{Href: "/chart/circle-packing", Title: "Circle packing", Description: "Welzl enclosing-circle packing (d3-hierarchy), colored by depth."},
			{Href: "/chart/tree", Title: "Tree", Description: "Tidy-tree / dendrogram node-link diagrams (d3-hierarchy) with bump links."},
			{Href: "/chart/voronoi", Title: "Voronoi", Description: "Delaunay triangulation + Voronoi cells (d3-delaunay); links, cells, points, bounds."},
			{Href: "/chart/network", Title: "Network", Description: "Force-directed node/link graph (d3-force): link + many-body + centering forces, deterministic fixed-tick layout."},
			{Href: "/chart/swarmplot", Title: "Swarmplot", Description: "Grouped value distribution relaxed with d3-force (ForceX/Y + collide); voronoi-mesh hover."},
			{Href: "/chart/sankey", Title: "Sankey", Description: "Flow diagram (d3-sankey): node breadths + relaxation, variable-thickness monotone-curve ribbons."},
			{Href: "/chart/chord", Title: "Chord", Description: "Radial flow diagram (d3-chord): entity arcs sized by total flow, ribbons spanning each directed sub-flow; hover an entity to highlight it."},
			{Href: "/chart/geo", Title: "Geo", Description: "GeoJSON maps (d3-geo): GeoMap + value-bound Choropleth, ten projections, optional graticule and continuous legend."},
			{Href: "/scales", Title: "Scales", Description: "Linear vs log value axes on bar, line, scatterplot, and swarmplot — data spanning orders of magnitude, toggled per chart via HTMX."},
			{Href: "/styling", Title: "Styling", Description: "Gradients, pattern fills, conditional match rules, and blend modes via the Defs + Fill props."},
			{Href: "/legends", Title: "Legends", Description: "Symbol shapes, anchors/directions, symbol borders, the continuous color legend, and an HTML legend outside the SVG."},
			{Href: "/composition", Title: "Composition", Description: "Build-your-own charts from the Use* hooks + sub-components: custom point symbols, direct labels, sparkline KPI cards."},
			{Href: "/dashboard", Title: "Dashboard", Description: "A composed real-world dashboard: KPI sparklines, gradient trend chart, donut, stacked bars, and bullet goals under one dark theme."},
			{Href: "/palettes", Title: "Palettes", Description: "The full color-palette catalog (categorical, sequential, diverging) applied to bars, with swatches."},
			{Href: "/themes", Title: "Themes", Description: "Bar / line / pie under default, dark, and custom themes."},
			{Href: "/benchmark", Title: "Benchmark", Description: "Server-side render time + SVG size for a bar chart across dataset sizes (load/scaling showcase)."},
		},
	}
	a.renderPage(w, templates.LayoutProps{Title: "templ-charts demo", Nav: "home"}, templates.IndexPage(props))
}

// Bar handles GET /bar: the bar demos page.
func (a *App) Bar(w http.ResponseWriter, r *http.Request) {
	ds := demos.BarDemos()
	a.ensureRegistered(ds)
	cards := a.demoCards(ds)
	a.renderPage(w, templates.LayoutProps{Title: "Bar charts", Nav: "bar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Bar chart demos. Hover a bar for a tooltip; click a legend item to toggle a series.",
		Cards: cards,
	}))
}

// Line handles GET /line: the line demos page.
func (a *App) Line(w http.ResponseWriter, r *http.Request) {
	ds := demos.LineDemos()
	a.ensureRegistered(ds)
	cards := a.demoCards(ds)
	a.renderPage(w, templates.LayoutProps{Title: "Line charts", Nav: "line"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Line chart demos. The slices demo shows a per-series tooltip on vertical-slice hover.",
		Cards: cards,
	}))
}

// Pie handles GET /pie: the pie demos page.
func (a *App) Pie(w http.ResponseWriter, r *http.Request) {
	demos := demos.PieDemos()
	a.ensureRegistered(demos)
	cards := a.demoCards(demos)
	a.renderPage(w, templates.LayoutProps{Title: "Pie charts", Nav: "pie"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Pie chart demos. Hover an arc to pop it (active highlight); click a legend item to toggle a slice.",
		Cards: cards,
	}))
}

// Styling handles GET /styling: gradients, patterns, match rules, and blend
// modes via the Defs + Fill props.
func (a *App) Styling(w http.ResponseWriter, r *http.Request) {
	demos := demos.StylingDemos()
	a.ensureRegistered(demos)
	cards := a.demoCards(demos)
	a.renderPage(w, templates.LayoutProps{Title: "Styling", Nav: "styling"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Visual styling via SVG defs: linear gradients, dot/line patterns, and conditional fills bound with match rules (Defs []core.Def + Fill []core.DefRule, mirroring nivo's Patterns & Gradients guide), plus CSS mix-blend-mode on line areas. Colors marked \"inherit\" in a def resolve to each mark's own series color.",
		Cards: cards,
	}))
}

// Legends handles GET /legends: legend customization — symbol shapes,
// anchors/directions, symbol borders, the continuous color legend, and an
// HTML legend living outside the SVG.
func (a *App) Legends(w http.ResponseWriter, r *http.Request) {
	ds := demos.LegendsDemos()
	a.ensureRegistered(ds)
	cards := a.demoCards(ds)
	for i := range cards {
		if cards[i].ID == demos.HTMLLegendDemoID {
			cards[i].FooterHTML = demos.HTMLLegendFooter()
		}
	}
	a.renderPage(w, templates.LayoutProps{Title: "Legends", Nav: "legends"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Legend customization: the four symbol shapes across the four corner anchors, a bottom row legend with bordered symbols, the continuous color legend for value→color scales, and a legend built as plain HTML outside the SVG — same data, same HTMX toggle endpoint. Every legend here toggles its series on click.",
		Cards: cards,
	}))
}

// Composition handles GET /composition: charts assembled by hand from the
// Use* hooks + exported sub-components + custom SVG marks (static tiles).
func (a *App) Composition(w http.ResponseWriter, r *http.Request) {
	ds, err := demos.CompositionDemos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		cards = append(cards, templates.ChartCardProps{
			ID: d.ID, Title: d.Title, Description: d.Description,
			SVG: d.SVG, Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Composition", Nav: "composition"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Build-your-own charts: every chart package exports its Use* hook (the layout math) and its layer sub-components, so a custom chart is plain Go — call line.UseLine for scales/series/generators, reuse axes.Grid/axes.Axes/line.Lines for the standard layers, and write SVG for the marks you want to own. These tiles are the d3-style escape hatch: custom point symbols, direct labels, and sparkline KPI cards, all composed without a top-level chart component.",
		Cards: cards,
	}))
}

// Dashboard handles GET /dashboard: a composed real-world dashboard — KPI
// sparkline cards, an interactive trend chart, a donut + stacked-bar side
// column, and a bullet goal row, all sharing one dark theme + palette.
func (a *App) Dashboard(w http.ResponseWriter, r *http.Request) {
	ds := demos.DashboardDemos()
	a.ensureRegistered(ds)
	svgs := make(map[string]string, len(ds))
	for _, d := range ds {
		svg, err := a.handler.RenderFull(d.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		svgs[d.ID] = svg
	}

	kpis, err := demos.DashboardKPIs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cards := make([]templates.DashboardKPICard, 0, len(kpis))
	for _, k := range kpis {
		cards = append(cards, templates.DashboardKPICard{
			Label: k.Label, Value: k.Value, Delta: k.Delta, DeltaUp: k.DeltaUp, SVG: k.SVG,
		})
	}

	var bb strings.Builder
	if err := bullet.Bullet(demos.DashboardBulletProps()).Render(r.Context(), &bb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.renderPage(w, templates.LayoutProps{Title: "Dashboard", Nav: "dashboard"}, templates.DashboardPage(templates.DashboardPageProps{
		KPIs:      cards,
		MainID:    demos.DashboardMainID,
		MainSVG:   svgs[demos.DashboardMainID],
		DonutID:   demos.DashboardDonutID,
		DonutSVG:  svgs[demos.DashboardDonutID],
		BarID:     demos.DashboardBarID,
		BarSVG:    svgs[demos.DashboardBarID],
		BulletSVG: bb.String(),
	}))
}

// Palettes handles GET /palettes: the palette gallery page. Each tile renders
// a bar chart colored by one catalog palette plus a swatch strip. Static (no
// HTMX) — like the themes page.
func (a *App) Palettes(w http.ResponseWriter, r *http.Request) {
	pd := demos.PaletteDemos()
	cards := make([]templates.PaletteCardProps, 0, len(pd))
	for _, d := range pd {
		svg, err := renderStatic(d.Kind, d.Props)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.PaletteCardProps{
			ID:             string(d.Palette.ID),
			Name:           d.Palette.Name,
			Kind:           d.Palette.Kind.String(),
			Group:          d.Palette.Group,
			ColorblindSafe: d.Palette.ColorblindSafe,
			Swatch:         d.Swatch,
			SVG:            svg,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Palettes", Nav: "palettes"}, templates.PaletteGalleryPage(templates.PaletteGalleryPageProps{
		Intro: "Color palettes applied to a bar chart. Categorical palettes cycle discrete colors across series; sequential and diverging palettes are sampled into discrete steps. Set a palette via Colors: colors.Scheme(colors.PaletteTableau10).",
		Cards: cards,
	}))
}

// Heatmap handles GET /heatmap: the heatmap demos page. Heatmaps render
// statically (no HTMX) — interactivity arrives with the charts/interact client
// layer.
func (a *App) Heatmap(w http.ResponseWriter, r *http.Request) {
	demos := demos.HeatmapDemos()
	a.ensureRegistered(demos)
	cards := a.demoCards(demos)
	a.renderPage(w, templates.LayoutProps{Title: "Heatmap", Nav: "heatmap"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Heatmap demos: sequential and diverging color scales, cell borders, value labels, and a continuous legend. Hover a cell for a tooltip (the others dim), like the bar/pie demos.",
		Cards: cards,
	}))
}

// Waffle handles GET /waffle: the waffle demos page (static SVG).
func (a *App) Waffle(w http.ResponseWriter, r *http.Request) {
	ds := demos.WaffleDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := waffle.Waffle(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{
			ID: d.ID, Title: d.Title, Description: d.Description,
			SVG: b.String(), Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Waffle", Nav: "waffle"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Waffle demos: part-of-whole grids built on the charts/grid layout, with fill direction, cell borders, and a legend. Static SVG.",
		Cards: cards,
	}))
}

// Calendar handles GET /calendar: the calendar demos page (static SVG).
func (a *App) Calendar(w http.ResponseWriter, r *http.Request) {
	ds := demos.CalendarDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := calendar.Calendar(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{
			ID: d.ID, Title: d.Title, Description: d.Description,
			SVG: b.String(), Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Calendar", Nav: "calendar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Calendar heatmap demos: a day grid over a date range with quantized colors, month/year legends, horizontal and vertical layouts. Static SVG.",
		Cards: cards,
	}))
}

// Radar handles GET /radar: the radar demos page (static SVG).
func (a *App) Radar(w http.ResponseWriter, r *http.Request) {
	ds := demos.RadarDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := radar.Radar(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{
			ID: d.ID, Title: d.Title, Description: d.Description,
			SVG: b.String(), Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Radar", Nav: "radar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Radar demos: values per key plotted around shared indices, with circular or polygon grid levels, filled areas, dots, and an optional legend. Static SVG.",
		Cards: cards,
	}))
}

// RadialBar handles GET /radial-bar: the radial-bar demos page (static SVG).
func (a *App) RadialBar(w http.ResponseWriter, r *http.Request) {
	ds := demos.RadialBarDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := radialbar.RadialBar(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{
			ID: d.ID, Title: d.Title, Description: d.Description,
			SVG: b.String(), Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Radial bar", Nav: "radial-bar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Radial-bar demos: stacked bars drawn as arcs in polar space (built on charts/polar-axes + charts/arcs), with background tracks, radial/circular axes, optional labels, and a legend. Static SVG.",
		Cards: cards,
	}))
}

// ScatterPlot handles GET /scatterplot: the scatterplot demos page (static SVG).
func (a *App) ScatterPlot(w http.ResponseWriter, r *http.Request) {
	ds := demos.ScatterPlotDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := scatterplot.ScatterPlot(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Scatterplot", Nav: "scatterplot"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Scatterplot demos: series of {x,y} nodes on linear scales with grid, axes, and a legend, plus a Voronoi-mesh hover tile (nearest-node detection via internal/d3/delaunay).",
		Cards: cards,
	}))
}

// Stream handles GET /stream: the stream demos page (static SVG).
func (a *App) Stream(w http.ResponseWriter, r *http.Request) {
	ds := demos.StreamDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := stream.Stream(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Stream", Nav: "stream"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Stream demos: stacked areas with configurable offset (wiggle/silhouette/expand) and a smooth curve. Static SVG.",
		Cards: cards,
	}))
}

// Bullet handles GET /bullet: the bullet demos page (static SVG).
func (a *App) Bullet(w http.ResponseWriter, r *http.Request) {
	ds := demos.BulletDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := bullet.Bullet(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Bullet", Nav: "bullet"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Bullet demos: stacked ranges, measure bars, and comparative markers on a shared value scale, with a per-row axis. Static SVG.",
		Cards: cards,
	}))
}

// Funnel handles GET /funnel: the funnel demos page (static SVG).
func (a *App) Funnel(w http.ResponseWriter, r *http.Request) {
	ds := demos.FunnelDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := funnel.Funnel(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Funnel", Nav: "funnel"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Funnel demos: ordered parts with smooth or linear trapezoid bands, side borders, separators, and centered labels. Static SVG.",
		Cards: cards,
	}))
}

// BoxPlot handles GET /boxplot: the boxplot demos page (static SVG).
func (a *App) BoxPlot(w http.ResponseWriter, r *http.Request) {
	ds := demos.BoxPlotDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := boxplot.BoxPlot(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Box plot", Nav: "boxplot"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Box-plot demos: raw observations summarized to quantiles (q10/q25/median/q75/q90 via internal/d3/array.Quantile) and drawn as box + whisker glyphs. Static SVG.",
		Cards: cards,
	}))
}

// Bump handles GET /bump: the bump demos page (static SVG).
func (a *App) Bump(w http.ResponseWriter, r *http.Request) {
	ds := demos.BumpDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := bump.Bump(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Bump", Nav: "bump"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Bump demos: ranking over time as smooth (curveBumpX) or linear lines, with end labels and per-point hover, plus a Voronoi-mesh hover tile (nearest-point detection via internal/d3/delaunay). Static SVG.",
		Cards: cards,
	}))
}

// Marimekko handles GET /marimekko: the marimekko demos page (static SVG).
func (a *App) Marimekko(w http.ResponseWriter, r *http.Request) {
	ds := demos.MarimekkoDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := marimekko.Marimekko(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Marimekko", Nav: "marimekko"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Marimekko demos: variable-width stacked bars — bar width by value, segments stacked via d3.Stack (offset none/expand). Static SVG.",
		Cards: cards,
	}))
}

// ParallelCoordinates handles GET /parallel-coordinates (static SVG).
func (a *App) ParallelCoordinates(w http.ResponseWriter, r *http.Request) {
	ds := demos.ParallelCoordinatesDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := pc.ParallelCoordinates(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Parallel coordinates", Nav: "parallel-coordinates"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Parallel-coordinates demos: one axis per variable (linear or point), each record a polyline across them, horizontal and vertical layouts, plus a Voronoi-mesh hover tile that resolves to the nearest datum via its per-axis vertices (internal/d3/delaunay). Static SVG.",
		Cards: cards,
	}))
}

// PolarBar handles GET /polar-bar: the polar-bar demos page (static SVG).
func (a *App) PolarBar(w http.ResponseWriter, r *http.Request) {
	ds := demos.PolarBarDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := polarbar.PolarBar(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Polar bar", Nav: "polar-bar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Polar-bar demos: stacked bars wrapped into a full circle — angle band per index, keys stacked along the radius, with radial/circular grids and index labels. Static SVG.",
		Cards: cards,
	}))
}

// Treemap handles GET /treemap: the treemap demos page (static SVG).
func (a *App) Treemap(w http.ResponseWriter, r *http.Request) {
	ds := demos.TreemapDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := treemap.Treemap(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Treemap", Nav: "treemap"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Treemap demos: a hierarchy tiled as nested rectangles (squarify/binary), leaf + parent labels, colored by top-level ancestor. Static SVG.",
		Cards: cards,
	}))
}

// Sunburst handles GET /sunburst: the sunburst demos page (static SVG).
func (a *App) Sunburst(w http.ResponseWriter, r *http.Request) {
	ds := demos.SunburstDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := sunburst.Sunburst(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Sunburst", Nav: "sunburst"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Sunburst demos: the partition layout mapped onto polar arcs ([2π, r²]), colors inheriting from the top-level ancestor. Static SVG.",
		Cards: cards,
	}))
}

// Icicle handles GET /icicle: the icicle demos page (static SVG).
func (a *App) Icicle(w http.ResponseWriter, r *http.Request) {
	ds := demos.IcicleDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := icicle.Icicle(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Icicle", Nav: "icicle"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Icicle demos: the partition layout as depth-banded rectangles, oriented bottom/top/left/right. Static SVG.",
		Cards: cards,
	}))
}

// CirclePacking handles GET /circle-packing (static SVG).
func (a *App) CirclePacking(w http.ResponseWriter, r *http.Request) {
	ds := demos.CirclePackingDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := cp.CirclePacking(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Circle packing", Nav: "circle-packing"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Circle-packing demos: Welzl enclosing-circle packing (deterministic via the ported LCG), colored by depth. Static SVG.",
		Cards: cards,
	}))
}

// Tree handles GET /tree: the tree demos page (static SVG).
func (a *App) Tree(w http.ResponseWriter, r *http.Request) {
	ds := demos.TreeDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := tree.Tree(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Tree", Nav: "tree"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Tree demos: tidy-tree and dendrogram layouts with smooth bump links (curveBumpX/Y), in four orientations, plus a Voronoi-mesh hover tile (nearest-node detection via internal/d3/delaunay). Static SVG.",
		Cards: cards,
	}))
}

// Voronoi handles GET /voronoi: the voronoi demos page (static SVG).
func (a *App) Voronoi(w http.ResponseWriter, r *http.Request) {
	ds := demos.VoronoiDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := voronoi.Voronoi(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Voronoi", Nav: "voronoi"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Voronoi demos: Delaunay triangulation and its Voronoi dual (d3-delaunay), clipped to the chart — links, cells, points, and bounds, with per-cell hover and optional colored (ordinal-filled) cells. Static SVG.",
		Cards: cards,
	}))
}

// Network handles GET /network: the network demos page (static SVG; hover a
// node for a client-side tooltip).
func (a *App) Network(w http.ResponseWriter, r *http.Request) {
	ds := demos.NetworkDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := network.Network(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Network", Nav: "network"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Network demos: a force-directed graph laid out with internal/d3/force (link + many-body + centering forces, a fixed 120-iteration tick). The three tiles share one graph and vary only repulsivity — stronger charge spreads the nodes to fill more of the frame. The layout is deterministic; hover a node for its id. A fourth tile adds Voronoi-mesh hover (nearest-node detection via internal/d3/delaunay), so hovering anywhere resolves to the closest node.",
		Cards: cards,
	}))
}

// SwarmPlot handles GET /swarmplot: the swarmplot demos page (static SVG).
func (a *App) SwarmPlot(w http.ResponseWriter, r *http.Request) {
	ds := demos.SwarmPlotDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := swarmplot.SwarmPlot(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Swarmplot", Nav: "swarmplot"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Swarmplot demos: points grouped along one axis, positioned by value along the other, then relaxed with internal/d3/force (ForceX/ForceY + ForceCollide, deterministic fixed-tick). Vertical, horizontal, and voronoi-mesh hover.",
		Cards: cards,
	}))
}

// Sankey handles GET /sankey: the sankey demos page (static SVG; hover a node
// or ribbon for a client-side tooltip).
func (a *App) Sankey(w http.ResponseWriter, r *http.Request) {
	ds := demos.SankeyDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := sankey.Sankey(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Sankey", Nav: "sankey"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Sankey demos: flow diagrams laid out with internal/d3/sankey (a faithful d3-sankey port — node breadths, relaxation, collision resolution) and rendered as variable-thickness ribbons via internal/d3/shape's line generator + curveMonotoneX/Y, exactly as nivo does. Horizontal and vertical layouts, an alignment variant, and an interactive tile with a legend.",
		Cards: cards,
	}))
}

// Chord handles GET /chord: the chord demos page (static SVG; hover an arc or
// ribbon for a client-side tooltip).
func (a *App) Chord(w http.ResponseWriter, r *http.Request) {
	ds := demos.ChordDemos()
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		var b strings.Builder
		if err := chord.Chord(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Chord", Nav: "chord"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Chord demos: entity-to-entity flows laid out with internal/d3/chord (a faithful d3-chord port — group arcs sized by total flow, ribbons spanning each directed sub-flow). Arcs render via charts/arcs (d3-shape Arc) and ribbons via internal/d3/chord's Ribbon generator, exactly as nivo does. The default diagram, a padded/inset-ribbon variant, and an interactive tile where hovering an entity highlights it (and its ribbons) while the rest fade back.",
		Cards: cards,
	}))
}

// Geo handles GET /geo: the geo page with GeoMap and Choropleth demos (static
// SVG; hover a feature for a client-side tooltip on the interactive tiles).
func (a *App) Geo(w http.ResponseWriter, r *http.Request) {
	maps := demos.GeoMapDemos()
	choros := demos.ChoroplethDemos()
	cards := make([]templates.ChartCardProps, 0, len(maps)+len(choros))
	for _, d := range maps {
		var b strings.Builder
		if err := geo.GeoMap(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	for _, d := range choros {
		var b strings.Builder
		if err := geo.Choropleth(d.Props).Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{ID: d.ID, Title: d.Title, Description: d.Description, SVG: b.String()})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Geo", Nav: "geo"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Geo demos: GeoJSON features projected server-side with internal/d3/geo (a faithful d3-geo subset — projection machinery, three-axis rotation, adaptive resampling, antimeridian clipping, GeoPath→SVG, and a graticule generator). GeoMap fills every feature one color; Choropleth binds a {id,value} dataset onto the features by id and colors each by a quantize scale (nivo's 'PuBuGn' → purple_blue_green) with a continuous-color legend. The bundled world sample is coarsely simplified so the page stays self-contained.",
		Cards: cards,
	}))
}

// Themes handles GET /themes: the themes page (no interactivity — static SVG).
func (a *App) Themes(w http.ResponseWriter, r *http.Request) {
	td := demos.ThemesDemos()
	cards := make([]templates.ChartCardProps, 0, len(td))
	for _, d := range td {
		svg, err := renderStatic(d.Kind, d.Props)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cards = append(cards, templates.ChartCardProps{
			ID:          d.ID,
			Title:       d.Title,
			Description: "Theme: " + d.Title,
			SVG:         svg,
			Interactive: false,
		})
	}
	a.renderPage(w, templates.LayoutProps{Title: "Themes", Nav: "themes"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "The same chart types under default, dark, and custom (warm) themes. Static (no HTMX).",
		Cards: cards,
	}))
}

// --- helpers ---

// ensureRegistered idempotently registers each demo with the htmx registry.
// Re-registering resets per-instance state (matches htmx.Registry.Register).
func (a *App) ensureRegistered(ds []demos.Demo) {
	for _, d := range ds {
		a.registry.Register(d.ID, d.Kind, d.Props)
	}
}

// demoCards renders the initial SVG for each demo (via the htmx full-render
// path, which applies Interactive=true + ChartID) and builds a ChartCard
// per demo.
func (a *App) demoCards(ds []demos.Demo) []templates.ChartCardProps {
	cards := make([]templates.ChartCardProps, 0, len(ds))
	for _, d := range ds {
		svg, err := a.handler.RenderFull(d.ID)
		if err != nil {
			cards = append(cards, templates.ChartCardProps{
				ID: d.ID, Title: d.Title, Description: d.Description,
				SVG: "<!-- render error: " + err.Error() + " -->",
			})
			continue
		}
		cards = append(cards, templates.ChartCardProps{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			SVG:         svg,
			Interactive: true,
		})
	}
	return cards
}

// Scales handles GET /scales: the value-scale showcase page. It collects the
// linear/log demo for each continuous-value chart family (bar, line,
// scatterplot, swarmplot); each card's toggle swaps just its chart via HTMX.
func (a *App) Scales(w http.ResponseWriter, r *http.Request) {
	cards := []templates.ChartCardProps{
		a.scaleCard("bar", "Bar"),
		a.scaleCard("line", "Line"),
		a.scaleCard("scatterplot", "Scatterplot"),
		a.scaleCard("swarmplot", "Swarmplot"),
	}
	a.renderPage(w, templates.LayoutProps{Title: "Scales", Nav: "scales"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Value scales: each chart plots data spanning several orders of magnitude. Toggle linear ↔ log to see how a log axis keeps small values readable where a linear axis flattens them. The switch swaps just that chart via HTMX — no full-page reload.",
		Cards: cards,
	}))
}

// scaleCard builds one linear/log value-scale demo card for the /scales page.
// bar/line go through the htmx registry (so hover + the /scale re-render stay
// consistent with the shown scale); scatterplot/swarmplot render static SVG.
// title is the card heading (the chart family); chartKey is what the toggle's
// htmx requests carry. Cards render linear initially — the toggle switches them.
func (a *App) scaleCard(chartKey, title string) templates.ChartCardProps {
	card := func(id, desc, svg string, interactive bool) templates.ChartCardProps {
		return templates.ChartCardProps{
			ID: id, Title: title, Description: desc,
			SVG: svg, Interactive: interactive,
			FooterHTML: scaleToggle(chartKey, id, false, false),
		}
	}
	registered := func(d demos.Demo) templates.ChartCardProps {
		a.registry.Register(d.ID, d.Kind, d.Props)
		svg, err := a.handler.RenderFull(d.ID)
		if err != nil {
			return card(d.ID, d.Description, "<!-- render error: "+err.Error()+" -->", false)
		}
		return card(d.ID, d.Description, svg, true)
	}
	static := func(id, desc, svg string, err error) templates.ChartCardProps {
		if err != nil {
			return card(id, desc, "<!-- render error: "+err.Error()+" -->", false)
		}
		return card(id, desc, svg, false)
	}
	switch chartKey {
	case "bar":
		return registered(demos.BarScaleDemo(false))
	case "line":
		return registered(demos.LineScaleDemo(false))
	case "scatterplot":
		d := demos.ScatterPlotScaleDemo(false)
		var b strings.Builder
		err := scatterplot.ScatterPlot(d.Props).Render(context.Background(), &b)
		return static(d.ID, d.Description, b.String(), err)
	case "swarmplot":
		d := demos.SwarmPlotScaleDemo(false)
		var b strings.Builder
		err := swarmplot.SwarmPlot(d.Props).Render(context.Background(), &b)
		return static(d.ID, d.Description, b.String(), err)
	}
	return templates.ChartCardProps{}
}

// Scale handles GET /scale?chart=<key>&scale=<linear|log>: it re-renders a
// single value-scale demo chart under the requested scale and returns an htmx
// fragment — the bare chart SVG (swapped into #chart-<id>) plus an out-of-band
// copy of the toggle (updating the active chip). No full-page reload, so the
// scroll position is preserved. chart is one of bar/line/scatterplot/swarmplot.
func (a *App) Scale(w http.ResponseWriter, r *http.Request) {
	chart := r.URL.Query().Get("chart")
	logScale := r.URL.Query().Get("scale") == "log"
	var (
		svg, chartID string
		err          error
	)
	switch chart {
	case "bar":
		d := demos.BarScaleDemo(logScale)
		// Re-register so the interactive endpoints (hover) match the shown scale.
		a.registry.Register(d.ID, d.Kind, d.Props)
		chartID = d.ID
		svg, err = a.handler.RenderFull(d.ID)
	case "line":
		d := demos.LineScaleDemo(logScale)
		a.registry.Register(d.ID, d.Kind, d.Props)
		chartID = d.ID
		svg, err = a.handler.RenderFull(d.ID)
	case "scatterplot":
		d := demos.ScatterPlotScaleDemo(logScale)
		chartID = d.ID
		var b strings.Builder
		if err = scatterplot.ScatterPlot(d.Props).Render(r.Context(), &b); err == nil {
			svg = b.String()
		}
	case "swarmplot":
		d := demos.SwarmPlotScaleDemo(logScale)
		chartID = d.ID
		var b strings.Builder
		if err = swarmplot.SwarmPlot(d.Props).Render(r.Context(), &b); err == nil {
			svg = b.String()
		}
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, svg)
	fmt.Fprint(w, scaleToggle(chart, chartID, logScale, true))
}

// scaleToggle renders the linear/log switcher for a value-scale demo card as an
// htmx control. Each chip issues hx-get /scale, swapping the fresh SVG into
// #chart-<chartID> (innerHTML) so only the chart updates — no page reload. The
// response also carries this toggle again with hx-swap-oob so the active chip
// flips. Styles are inlined; the demo list pages don't ship the detail page's
// chip stylesheet. chartKey names the demo ("bar", "line", …).
func scaleToggle(chartKey, chartID string, logScale, oob bool) string {
	chip := func(scale, label string, active bool) string {
		style := "display:inline-block;padding:.15rem .6rem;margin-right:.4rem;border:1px solid #ccc;border-radius:999px;font-size:.85rem;cursor:pointer;color:#333"
		if active {
			style += ";background:#333;color:#fff;border-color:#333"
		}
		return fmt.Sprintf(
			`<a hx-get=%q hx-target=%q hx-swap="innerHTML" style=%q>%s</a>`,
			fmt.Sprintf("/scale?chart=%s&scale=%s", chartKey, scale),
			"#chart-"+chartID, style, label)
	}
	oobAttr := ""
	if oob {
		oobAttr = ` hx-swap-oob="true"`
	}
	return fmt.Sprintf(`<div id="scale-toggle-%s"%s style="margin-top:.6rem;display:flex;align-items:center;gap:.4rem">`, chartID, oobAttr) +
		`<span style="font-size:.72rem;color:#666;text-transform:uppercase;letter-spacing:.04em">value scale</span>` +
		chip("linear", "linear", !logScale) +
		chip("log", "log", logScale) +
		`</div>`
}

// renderPage wraps the page content in the Layout and writes the HTML.
func (a *App) renderPage(w http.ResponseWriter, props templates.LayoutProps, content templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.Layout(props, content).Render(context.Background(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderStatic renders a chart's SVG without htmx interactivity (used by the
// themes page).
func renderStatic(kind htmx.ChartKind, props any) (string, error) {
	var c templ.Component
	switch kind {
	case htmx.KindBar:
		c = bar.Bar(props.(bar.BarProps))
	case htmx.KindLine:
		c = line.Line(props.(line.LineProps))
	case htmx.KindPie:
		c = pie.Pie(props.(pie.PieProps))
	default:
		return "", errUnknownKind
	}
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

var errUnknownKind = stdErr("unknown chart kind")

type stdErr string

func (e stdErr) Error() string { return string(e) }
