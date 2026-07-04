package handlers

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/treemap"
	"github.com/geoffjay/templ-charts/examples/app/demos"
	"github.com/geoffjay/templ-charts/examples/app/handlers/entries"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

// paletteChoice is one option in the palette switcher.
type paletteChoice struct {
	ID    colors.PaletteID
	Label string
}

// paletteChoices is the curated set of categorical palettes offered by the
// palette switcher (categorical because most charts color series ordinally).
var paletteChoices = []paletteChoice{
	{colors.PaletteNivo, "nivo"},
	{colors.PaletteCategory10, "category10"},
	{colors.PaletteTableau10, "tableau10"},
	{colors.PaletteSet2, "set2"},
	{colors.PalettePaired, "paired"},
	{colors.PaletteDark2, "dark2"},
	{colors.PaletteObservable10, "observable10"},
	{colors.PaletteOkabeIto, "okabe-ito"},
}

// resolveTheme maps a ?theme= name to its Theme, defaulting to "default".
func resolveTheme(name string) (theme *theming.Theme, resolved string) {
	for _, g := range demos.ThemeGroups() {
		if g.Name == name {
			return g.Theme, g.Name
		}
	}
	def := &theming.DefaultTheme
	return def, "default"
}

// Detail handles GET /chart/{slug}: the full-width detail page for one chart,
// with live theme + palette switchers and the Go snippet to reproduce it.
// Switching re-renders the whole page via ?theme=/?palette= query params (a
// plain, robust server render; htmx is loaded globally if a future pass wants
// fragment swaps).
func (a *App) Detail(w http.ResponseWriter, r *http.Request) {
	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/chart/"), "/")
	entry, ok := entries.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	theme, themeName := resolveTheme(r.URL.Query().Get("theme"))
	palette := colors.PaletteID(r.URL.Query().Get("palette"))
	animate := r.URL.Query().Get("animate") == "1"
	// Canvas engine toggle — honored only for Canvas-capable charts (those with
	// a CanvasRender); ignored otherwise.
	canvasEngine := r.URL.Query().Get("engine") == "canvas" && entry.CanvasRender != nil

	// The four hierarchy charts are click-to-zoom: register a zoomable instance
	// with the htmx registry and mount it (so hx-get="/charts/<id>/zoom" +
	// breadcrumb round-trip through the existing handler), instead of a static
	// SVG. Every other chart keeps its plain static detail render.
	var chartHTML string
	if id, kind, props, ok := a.zoomableChart(slug, theme, palette); ok {
		a.registry.Register(id, kind, props)
		svg, err := a.handler.RenderFull(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var b strings.Builder
		if err := htmx.Mount(htmx.MountProps{ID: id, SVG: svg, Interactive: true, Class: "tc-detail-chart chart"}).
			Render(context.Background(), &b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chartHTML = b.String()
	} else {
		renderFn := entry.Render
		if canvasEngine {
			renderFn = entry.CanvasRender
		}
		svg, err := renderFn(theme, palette, animate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chartHTML = `<div class="tc-detail-chart chart">` + svg + `</div>`
	}

	html := buildDetailHTML(entry, chartHTML, themeName, palette, animate, canvasEngine)
	a.renderPage(w,
		templates.LayoutProps{Title: entry.Title + " — templ-charts demo"},
		templ.Raw(html))
}

// zoomableChart returns a registered-instance definition for the four
// hierarchy charts (icicle/treemap/circle-packing/sunburst) with EnableZooming
// + a ChartID set, using a multi-level sample so the drill-down is meaningful.
// ok is false for every other slug.
func (a *App) zoomableChart(slug string, theme *theming.Theme, palette colors.PaletteID) (string, htmx.ChartKind, any, bool) {
	scheme := colors.OrdinalColorScaleConfig{}
	if palette != "" {
		scheme = colors.Scheme(palette)
	}
	id := "detail-" + slug
	switch slug {
	case "icicle":
		p := icicle.IcicleProps{Width: 720, Height: 440, Responsive: true, EnableZooming: true, Theme: theme, Data: icicleSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindIcicle, p, true
	case "treemap":
		p := treemap.TreemapProps{Width: 720, Height: 440, Responsive: true, EnableZooming: true, Theme: theme, Data: treemapSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindTreemap, p, true
	case "circle-packing":
		p := cp.CirclePackingProps{Width: 720, Height: 440, Responsive: true, EnableZooming: true, Theme: theme, Data: cpSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindCirclePack, p, true
	case "sunburst":
		p := sunburst.SunburstProps{Width: 720, Height: 440, Responsive: true, EnableZooming: true, Theme: theme, Data: sunburstSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindSunburst, p, true
	}
	return "", "", nil, false
}

func icicleSample() icicle.IcicleNode {
	return icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
		{ID: "analytics", Children: []icicle.IcicleNode{
			{ID: "charts", Children: []icicle.IcicleNode{{ID: "icicle", Value: 8}, {ID: "sunburst", Value: 6}}},
			{ID: "dashboards", Value: 12},
		}},
		{ID: "billing", Children: []icicle.IcicleNode{{ID: "invoices", Value: 10}, {ID: "reports", Value: 5}}},
		{ID: "support", Value: 9},
	}}
}

func treemapSample() treemap.TreemapNode {
	return treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
		{ID: "viz", Children: []treemap.TreemapNode{
			{ID: "charts", Children: []treemap.TreemapNode{{ID: "bar", Value: 14}, {ID: "line", Value: 9}}},
			{ID: "maps", Value: 11},
		}},
		{ID: "colors", Children: []treemap.TreemapNode{{ID: "categorical", Value: 12}, {ID: "sequential", Value: 7}}},
		{ID: "layout", Value: 10},
	}}
}

func cpSample() cp.CirclePackingNode {
	return cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
		{ID: "A", Children: []cp.CirclePackingNode{
			{ID: "A1", Children: []cp.CirclePackingNode{{ID: "A1a", Value: 8}, {ID: "A1b", Value: 5}}},
			{ID: "A2", Value: 10},
		}},
		{ID: "B", Children: []cp.CirclePackingNode{{ID: "B1", Value: 12}, {ID: "B2", Value: 6}}},
		{ID: "C", Value: 9},
	}}
}

func sunburstSample() sunburst.SunburstNode {
	return sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
		{ID: "fruit", Children: []sunburst.SunburstNode{
			{ID: "citrus", Children: []sunburst.SunburstNode{{ID: "orange", Value: 8}, {ID: "lemon", Value: 4}}},
			{ID: "berry", Value: 10},
		}},
		{ID: "veg", Children: []sunburst.SunburstNode{{ID: "root-veg", Value: 9}, {ID: "leafy", Value: 6}}},
		{ID: "grain", Value: 7},
	}}
}

// buildDetailHTML assembles the detail page body: switchers, the full-width
// chart, and the code snippet. chartHTML is the ready-to-inject chart block
// (a static SVG wrapper, or the htmx.Mount container for zoomable charts).
func buildDetailHTML(e entries.ChartEntry, chartHTML, themeName string, palette colors.PaletteID, animate, canvasEngine bool) string {
	var b strings.Builder

	// animateSuffix / engineSuffix carry the current animate + engine selections
	// through the theme and palette hrefs so toggling one control does not reset
	// the others. Both empty by default, so those hrefs (and the rendered page)
	// are unchanged for charts without the Canvas engine or animation.
	animateSuffix := ""
	if animate {
		animateSuffix = "&animate=1"
	}
	engineSuffix := ""
	if canvasEngine {
		engineSuffix = "&engine=canvas"
	}
	// sharedSuffix preserves both animate and engine on the theme/palette hrefs.
	sharedSuffix := animateSuffix + engineSuffix

	b.WriteString(detailCSS)
	fmt.Fprintf(&b, `<p><a href="/">← all charts</a></p>`)
	fmt.Fprintf(&b, `<h2>%s</h2><p>%s</p>`, html.EscapeString(e.Title), html.EscapeString(e.Description))

	// Theme switcher.
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">theme</span>`)
	for _, g := range demos.ThemeGroups() {
		href := fmt.Sprintf("/chart/%s?theme=%s&palette=%s%s", e.Slug, g.Name, palette, sharedSuffix)
		b.WriteString(switchLink(href, g.Name, g.Name == themeName))
	}
	b.WriteString(`</div>`)

	// Palette switcher (a "default" reset plus the curated categorical set).
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">palette</span>`)
	b.WriteString(switchLink(fmt.Sprintf("/chart/%s?theme=%s%s", e.Slug, themeName, sharedSuffix), "default", palette == ""))
	for _, pc := range paletteChoices {
		href := fmt.Sprintf("/chart/%s?theme=%s&palette=%s%s", e.Slug, themeName, pc.ID, sharedSuffix)
		b.WriteString(switchLink(href, pc.Label, palette == pc.ID))
	}
	b.WriteString(`</div>`)

	// baseTP is the theme(+palette) prefix shared by the animate and engine
	// switchers.
	baseTP := fmt.Sprintf("/chart/%s?theme=%s", e.Slug, themeName)
	if palette != "" {
		baseTP = fmt.Sprintf("/chart/%s?theme=%s&palette=%s", e.Slug, themeName, palette)
	}

	// Animate switcher: off (no animate param) / on (&animate=1). Hrefs preserve
	// theme + palette + engine.
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">animate</span>`)
	b.WriteString(switchLink(baseTP+engineSuffix, "off", !animate))
	b.WriteString(switchLink(baseTP+"&animate=1"+engineSuffix, "on", animate))
	b.WriteString(`</div>`)

	// Engine switcher: only for Canvas-capable charts (scatterplot, heatmap).
	// svg (default, no engine param) / canvas. Hrefs preserve theme + palette +
	// animate.
	if e.CanvasRender != nil {
		b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">engine</span>`)
		b.WriteString(switchLink(baseTP+animateSuffix, "svg", !canvasEngine))
		b.WriteString(switchLink(baseTP+animateSuffix+"&engine=canvas", "canvas", canvasEngine))
		b.WriteString(`</div>`)
	}

	// Full-width chart (static wrapper or zoomable htmx.Mount container).
	b.WriteString(chartHTML)

	// Code snippet.
	b.WriteString(`<h3>Code</h3>`)
	b.WriteString(`<pre class="tc-snippet"><code>` + html.EscapeString(e.Snippet) + `</code></pre>`)

	return b.String()
}

// switchLink renders one switcher option as an anchor, marked active when
// selected.
func switchLink(href, label string, active bool) string {
	cls := "tc-chip"
	if active {
		cls = "tc-chip tc-chip-active"
	}
	return fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, cls, href, html.EscapeString(label))
}

// detailCSS is the small page-scoped stylesheet for the detail page controls
// and snippet block (the layout ships no styles for these).
const detailCSS = `<style>
.tc-switch{margin:.5rem 0;display:flex;flex-wrap:wrap;gap:.4rem;align-items:center}
.tc-switch-label{font-size:.8rem;color:#666;width:4rem;text-transform:uppercase;letter-spacing:.04em}
.tc-chip{display:inline-block;padding:.15rem .55rem;border:1px solid #ccc;border-radius:999px;font-size:.85rem;text-decoration:none;color:#333}
.tc-chip:hover{border-color:#888}
.tc-chip-active{background:#333;color:#fff;border-color:#333}
.tc-detail-chart{max-width:100%;margin:1rem 0;border:1px solid #eee;border-radius:6px;padding:.5rem}
.tc-snippet{background:#f6f8fa;border:1px solid #e1e4e8;border-radius:6px;padding:1rem;overflow:auto;font-size:.85rem;line-height:1.4}
</style>`
