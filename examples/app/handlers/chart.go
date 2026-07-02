// Package handlers implements the demo app's HTTP handlers. It wires the
// htmx.Handler for the interactive chart endpoints and renders the page
// handlers (/bar, /line, /pie, /themes) using the demo definitions in the
// demos package.
package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/calendar"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/funnel"
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
			{Href: "/bar", Title: "Bar charts", Description: "Stacked, grouped, markers + annotations, legend toggle, totals."},
			{Href: "/line", Title: "Line charts", Description: "Single & multi-series, area + points, slices, mesh hover."},
			{Href: "/pie", Title: "Pie charts", Description: "Plain, donut, half, sorted, active-arc hover, legend toggle."},
			{Href: "/heatmap", Title: "Heatmap", Description: "2D value grid: sequential/diverging color scales, labels, borders, continuous legend."},
			{Href: "/waffle", Title: "Waffle", Description: "Part-of-whole cell grid: fill direction, borders, legend (built on charts/grid)."},
			{Href: "/calendar", Title: "Calendar", Description: "Day-grid heatmap over a date range: quantized colors, month/year legends, horizontal/vertical."},
			{Href: "/radar", Title: "Radar", Description: "Polar line/area chart: values per key around shared indices, circular/polygon grids, dots, legend."},
			{Href: "/radial-bar", Title: "Radial bar", Description: "Stacked bars as polar arcs (charts/polar-axes + charts/arcs): tracks, radial/circular axes, labels."},
			{Href: "/scatterplot", Title: "Scatterplot", Description: "{x,y} nodes on linear/time scales: grid, axes, per-series colors, legend."},
			{Href: "/stream", Title: "Stream", Description: "Stacked areas with wiggle/silhouette/expand offsets and a smooth curve."},
			{Href: "/bullet", Title: "Bullet", Description: "KPI ranges + measure bars + markers on a shared value scale, per-row axis."},
			{Href: "/funnel", Title: "Funnel", Description: "Ordered parts as smooth/linear trapezoids with separators and labels."},
			{Href: "/boxplot", Title: "Box plot", Description: "Quantile box + whisker glyphs from raw observations (d3/array.Quantile)."},
			{Href: "/bump", Title: "Bump", Description: "Ranking over time: smooth (curveBumpX) or linear lines with end labels and point hover."},
			{Href: "/marimekko", Title: "Marimekko", Description: "Variable-width stacked bars: width by value, segments stacked via d3.Stack."},
			{Href: "/parallel-coordinates", Title: "Parallel coordinates", Description: "One axis per variable; each record a polyline across linear/point scales."},
			{Href: "/polar-bar", Title: "Polar bar", Description: "Stacked bars wrapped into a full circle: angle band per index, radius-stacked keys."},
			{Href: "/treemap", Title: "Treemap", Description: "Nested rectangles (d3-hierarchy): squarify/binary tiling, leaf + parent labels."},
			{Href: "/sunburst", Title: "Sunburst", Description: "Radial partition (d3-hierarchy): arcs by value, colors inherited down the tree."},
			{Href: "/icicle", Title: "Icicle", Description: "Depth-banded partition rectangles (d3-hierarchy), oriented four ways."},
			{Href: "/circle-packing", Title: "Circle packing", Description: "Welzl enclosing-circle packing (d3-hierarchy), colored by depth."},
			{Href: "/tree", Title: "Tree", Description: "Tidy-tree / dendrogram node-link diagrams (d3-hierarchy) with bump links."},
			{Href: "/voronoi", Title: "Voronoi", Description: "Delaunay triangulation + Voronoi cells (d3-delaunay); links, cells, points, bounds."},
			{Href: "/network", Title: "Network", Description: "Force-directed node/link graph (d3-force): link + many-body + centering forces, deterministic fixed-tick layout."},
			{Href: "/swarmplot", Title: "Swarmplot", Description: "Grouped value distribution relaxed with d3-force (ForceX/Y + collide); voronoi-mesh hover."},
			{Href: "/sankey", Title: "Sankey", Description: "Flow diagram (d3-sankey): node breadths + relaxation, variable-thickness monotone-curve ribbons."},
			{Href: "/palettes", Title: "Palettes", Description: "The full color-palette catalog (categorical, sequential, diverging) applied to bars, with swatches."},
			{Href: "/themes", Title: "Themes", Description: "Bar / line / pie under default, dark, and custom themes."},
		},
	}
	a.renderPage(w, templates.LayoutProps{Title: "templ-charts demo", Nav: "home"}, templates.IndexPage(props))
}

// Bar handles GET /bar: the bar demos page.
func (a *App) Bar(w http.ResponseWriter, r *http.Request) {
	demos := demos.BarDemos()
	a.ensureRegistered(demos)
	cards := a.demoCards(demos)
	a.renderPage(w, templates.LayoutProps{Title: "Bar charts", Nav: "bar"}, templates.DemosPage(templates.DemosPageProps{
		Intro: "Bar chart demos. Hover a bar for a tooltip; click a legend item to toggle a series.",
		Cards: cards,
	}))
}

// Line handles GET /line: the line demos page.
func (a *App) Line(w http.ResponseWriter, r *http.Request) {
	demos := demos.LineDemos()
	a.ensureRegistered(demos)
	cards := a.demoCards(demos)
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
// statically (no HTMX) in v2 — interactivity arrives with the Phase 5 client
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
// path, which applies IsInteractive=true + ChartID) and builds a ChartCard
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
