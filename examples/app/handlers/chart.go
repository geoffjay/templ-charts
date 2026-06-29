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
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
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
