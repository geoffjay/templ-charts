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
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/charts/treemap"
	"github.com/geoffjay/templ-charts/examples/app/demos"
	"github.com/geoffjay/templ-charts/examples/app/handlers/entries"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

// paletteKinds orders the optgroups of the palette dropdown. Every catalog
// palette is offered: sequential/diverging gradients are sampled into discrete
// steps by the ordinal scheme machinery, so they work on categorical charts too.
var paletteKinds = []colors.PaletteKind{colors.KindCategorical, colors.KindSequential, colors.KindDiverging}

// resolveSpace maps a ?space= name to a color-interpolation space, defaulting
// to RGB. Returns the space and its canonical name.
func resolveSpace(name string) (colors.Space, string) {
	switch name {
	case "lab":
		return colors.SpaceLab, "lab"
	case "lch":
		return colors.SpaceLch, "lch"
	default:
		return colors.SpaceRGB, "rgb"
	}
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
	// Color-space toggle — honored only for charts with a SpaceRender (a
	// sequential/diverging scale). Defaults to RGB, so other charts are unaffected.
	space, spaceName := resolveSpace(r.URL.Query().Get("space"))
	// Value-scale toggle — honored only for charts with a ScaleRender (a
	// continuous value axis). Defaults to linear, so other charts are unaffected.
	logScale := r.URL.Query().Get("scale") == "log"

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
		var svg string
		var err error
		switch {
		case canvasEngine:
			svg, err = entry.CanvasRender(theme, palette, animate)
		case entry.ScaleRender != nil:
			svg, err = entry.ScaleRender(theme, palette, animate, logScale)
		case entry.SpaceRender != nil:
			svg, err = entry.SpaceRender(theme, palette, animate, space)
		default:
			svg, err = entry.Render(theme, palette, animate)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chartHTML = `<div class="tc-detail-chart chart">` + svg + `</div>`
	}

	html := buildDetailHTML(entry, chartHTML, themeName, palette, animate, canvasEngine, spaceName, logScale)
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
		p := icicle.IcicleProps{Width: 720, Height: 440, Responsive: true, EnableZooming: core.BoolPtr(true), Theme: theme, Data: icicleSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindIcicle, p, true
	case "treemap":
		p := treemap.TreemapProps{Width: 720, Height: 440, Responsive: true, EnableZooming: core.BoolPtr(true), Theme: theme, Data: treemapSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindTreemap, p, true
	case "circle-packing":
		p := cp.CirclePackingProps{Width: 720, Height: 440, Responsive: true, EnableZooming: core.BoolPtr(true), Theme: theme, Data: cpSample()}
		if palette != "" {
			p.Colors = scheme
		}
		return id, htmx.KindCirclePack, p, true
	case "sunburst":
		p := sunburst.SunburstProps{Width: 720, Height: 440, Responsive: true, EnableZooming: core.BoolPtr(true), Theme: theme, Data: sunburstSample()}
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
			{ID: "charts", Children: []icicle.IcicleNode{{ID: "icicle", Value: 8}, {ID: "sunburst", Value: 6}, {ID: "treemap", Value: 7}, {ID: "bar", Value: 11}}},
			{ID: "dashboards", Children: []icicle.IcicleNode{{ID: "sales", Value: 9}, {ID: "ops", Value: 5}}},
			{ID: "exports", Value: 6},
		}},
		{ID: "billing", Children: []icicle.IcicleNode{
			{ID: "invoices", Children: []icicle.IcicleNode{{ID: "drafts", Value: 4}, {ID: "sent", Value: 7}}},
			{ID: "reports", Value: 5},
			{ID: "payments", Value: 8},
		}},
		{ID: "support", Children: []icicle.IcicleNode{{ID: "tickets", Value: 10}, {ID: "chat", Value: 6}, {ID: "kb", Value: 4}}},
		{ID: "auth", Children: []icicle.IcicleNode{{ID: "sso", Value: 5}, {ID: "sessions", Value: 7}}},
		{ID: "search", Value: 9},
		{ID: "notifications", Value: 6},
		{ID: "storage", Children: []icicle.IcicleNode{{ID: "blobs", Value: 8}, {ID: "cache", Value: 3}}},
		{ID: "admin", Value: 5},
	}}
}

func treemapSample() treemap.TreemapNode {
	return treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
		{ID: "viz", Children: []treemap.TreemapNode{
			{ID: "charts", Children: []treemap.TreemapNode{{ID: "bar", Value: 14}, {ID: "line", Value: 9}, {ID: "pie", Value: 6}, {ID: "radar", Value: 4}}},
			{ID: "maps", Children: []treemap.TreemapNode{{ID: "choropleth", Value: 7}, {ID: "tiles", Value: 4}}},
			{ID: "tables", Value: 6},
		}},
		{ID: "colors", Children: []treemap.TreemapNode{{ID: "categorical", Value: 12}, {ID: "sequential", Value: 7}, {ID: "diverging", Value: 5}}},
		{ID: "layout", Children: []treemap.TreemapNode{{ID: "grid", Value: 8}, {ID: "flex", Value: 5}}},
		{ID: "interact", Children: []treemap.TreemapNode{{ID: "tooltip", Value: 9}, {ID: "zoom", Value: 6}, {ID: "brush", Value: 3}}},
		{ID: "data", Children: []treemap.TreemapNode{{ID: "loaders", Value: 7}, {ID: "transforms", Value: 10}}},
		{ID: "themes", Value: 8},
		{ID: "legends", Value: 6},
		{ID: "axes", Value: 9},
	}}
}

func cpSample() cp.CirclePackingNode {
	return cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
		{ID: "north", Children: []cp.CirclePackingNode{
			{ID: "n-web", Children: []cp.CirclePackingNode{{ID: "n-web-a", Value: 8}, {ID: "n-web-b", Value: 5}, {ID: "n-web-c", Value: 3}}},
			{ID: "n-mobile", Value: 10},
			{ID: "n-api", Value: 6},
		}},
		{ID: "south", Children: []cp.CirclePackingNode{{ID: "s-web", Value: 12}, {ID: "s-mobile", Value: 6}, {ID: "s-api", Value: 4}}},
		{ID: "east", Children: []cp.CirclePackingNode{
			{ID: "e-web", Children: []cp.CirclePackingNode{{ID: "e-web-a", Value: 7}, {ID: "e-web-b", Value: 4}}},
			{ID: "e-api", Value: 9},
		}},
		{ID: "west", Children: []cp.CirclePackingNode{{ID: "w-web", Value: 11}, {ID: "w-api", Value: 5}}},
		{ID: "central", Value: 9},
		{ID: "nordics", Children: []cp.CirclePackingNode{{ID: "no-web", Value: 6}, {ID: "no-mobile", Value: 4}}},
		{ID: "apac", Value: 7},
		{ID: "latam", Value: 5},
	}}
}

func sunburstSample() sunburst.SunburstNode {
	return sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
		{ID: "fruit", Children: []sunburst.SunburstNode{
			{ID: "citrus", Children: []sunburst.SunburstNode{{ID: "orange", Value: 8}, {ID: "lemon", Value: 4}, {ID: "lime", Value: 3}}},
			{ID: "berry", Children: []sunburst.SunburstNode{{ID: "strawberry", Value: 6}, {ID: "blueberry", Value: 4}}},
			{ID: "stone", Value: 5},
		}},
		{ID: "veg", Children: []sunburst.SunburstNode{{ID: "root-veg", Value: 9}, {ID: "leafy", Value: 6}, {ID: "squash", Value: 4}}},
		{ID: "grain", Children: []sunburst.SunburstNode{{ID: "wheat", Value: 7}, {ID: "rice", Value: 6}, {ID: "oats", Value: 3}}},
		{ID: "dairy", Children: []sunburst.SunburstNode{{ID: "milk", Value: 5}, {ID: "cheese", Value: 7}}},
		{ID: "protein", Children: []sunburst.SunburstNode{{ID: "fish", Value: 6}, {ID: "beans", Value: 4}, {ID: "nuts", Value: 3}}},
		{ID: "spice", Value: 4},
		{ID: "herbs", Value: 3},
		{ID: "oils", Value: 5},
	}}
}

// buildDetailHTML assembles the detail page body: switchers, the full-width
// chart, and the code snippet. chartHTML is the ready-to-inject chart block
// (a static SVG wrapper, or the htmx.Mount container for zoomable charts).
func buildDetailHTML(e entries.ChartEntry, chartHTML, themeName string, palette colors.PaletteID, animate, canvasEngine bool, spaceName string, logScale bool) string {
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
	// spaceSuffix carries a non-default color space (only for SpaceRender charts).
	spaceSuffix := ""
	if e.SpaceRender != nil && spaceName != "rgb" {
		spaceSuffix = "&space=" + spaceName
	}
	// scaleSuffix carries a non-default value scale (only for ScaleRender charts).
	scaleSuffix := ""
	if e.ScaleRender != nil && logScale {
		scaleSuffix = "&scale=log"
	}
	// sharedSuffix preserves animate + engine + space + scale on the theme/palette hrefs.
	sharedSuffix := animateSuffix + engineSuffix + spaceSuffix + scaleSuffix

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

	// Palette switcher: a dropdown over the full palette catalog, grouped by
	// kind, navigating on change (each option's value is the page URL for that
	// palette). A swatch strip previews the selected palette.
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">palette</span>`)
	b.WriteString(`<select class="tc-select" onchange="location.href=this.value">`)
	fmt.Fprintf(&b, `<option value="%s"%s>default (theme colors)</option>`,
		fmt.Sprintf("/chart/%s?theme=%s%s", e.Slug, themeName, sharedSuffix), selectedAttr(palette == ""))
	for _, kind := range paletteKinds {
		fmt.Fprintf(&b, `<optgroup label="%s">`, kind)
		for _, p := range colors.PalettesByKind(kind) {
			href := fmt.Sprintf("/chart/%s?theme=%s&palette=%s%s", e.Slug, themeName, p.ID, sharedSuffix)
			label := p.Name
			if p.ColorblindSafe {
				label += " · colorblind-safe"
			}
			fmt.Fprintf(&b, `<option value="%s"%s>%s</option>`, href, selectedAttr(palette == p.ID), html.EscapeString(label))
		}
		b.WriteString(`</optgroup>`)
	}
	b.WriteString(`</select>`)
	if p, ok := colors.LookupPalette(palette); ok {
		b.WriteString(`<span class="tc-swatch-strip">`)
		for _, c := range p.Swatch(8) {
			fmt.Fprintf(&b, `<span style="background:%s"></span>`, html.EscapeString(c))
		}
		b.WriteString(`</span>`)
	}
	b.WriteString(`</div>`)

	// baseTP is the theme(+palette) prefix shared by the animate and engine
	// switchers.
	baseTP := fmt.Sprintf("/chart/%s?theme=%s", e.Slug, themeName)
	if palette != "" {
		baseTP = fmt.Sprintf("/chart/%s?theme=%s&palette=%s", e.Slug, themeName, palette)
	}

	// Animate switcher: off (no animate param) / on (&animate=1). Hrefs preserve
	// theme + palette + engine + space + scale.
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">animate</span>`)
	b.WriteString(switchLink(baseTP+engineSuffix+spaceSuffix+scaleSuffix, "off", !animate))
	b.WriteString(switchLink(baseTP+"&animate=1"+engineSuffix+spaceSuffix+scaleSuffix, "on", animate))
	b.WriteString(`</div>`)

	// Engine switcher: only for Canvas-capable charts (scatterplot, heatmap).
	// svg (default, no engine param) / canvas. Hrefs preserve theme + palette +
	// animate + space + scale.
	if e.CanvasRender != nil {
		b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">engine</span>`)
		b.WriteString(switchLink(baseTP+animateSuffix+spaceSuffix+scaleSuffix, "svg", !canvasEngine))
		b.WriteString(switchLink(baseTP+animateSuffix+"&engine=canvas"+spaceSuffix+scaleSuffix, "canvas", canvasEngine))
		b.WriteString(`</div>`)
	}

	// Color-space switcher: only for charts with a SpaceRender (sequential/
	// diverging scale). rgb (default) / lab / lch. Hrefs preserve theme +
	// palette + animate + engine + scale.
	if e.SpaceRender != nil {
		aes := animateSuffix + engineSuffix + scaleSuffix
		b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">color space</span>`)
		b.WriteString(switchLink(baseTP+aes, "rgb", spaceName == "rgb"))
		b.WriteString(switchLink(baseTP+aes+"&space=lab", "lab", spaceName == "lab"))
		b.WriteString(switchLink(baseTP+aes+"&space=lch", "lch", spaceName == "lch"))
		b.WriteString(`</div>`)
	}

	// Value-scale switcher: only for charts with a ScaleRender (a continuous
	// value axis). linear (default) / log. Hrefs preserve theme + palette +
	// animate + engine + space.
	if e.ScaleRender != nil {
		aes := animateSuffix + engineSuffix + spaceSuffix
		b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">value scale</span>`)
		b.WriteString(switchLink(baseTP+aes, "linear", !logScale))
		b.WriteString(switchLink(baseTP+aes+"&scale=log", "log", logScale))
		b.WriteString(`</div>`)
	}

	// Full-width chart (static wrapper or zoomable htmx.Mount container).
	b.WriteString(chartHTML)

	// Code snippet, with a copy-to-clipboard button.
	b.WriteString(`<h3>Code</h3>`)
	b.WriteString(`<div class="tc-snippet-wrap">`)
	b.WriteString(`<button type="button" class="tc-copy" onclick="navigator.clipboard.writeText(this.nextElementSibling.textContent).then(()=>{this.textContent='copied';setTimeout(()=>{this.textContent='copy'},1200)})">copy</button>`)
	b.WriteString(`<pre class="tc-snippet"><code>` + html.EscapeString(e.Snippet) + `</code></pre>`)
	b.WriteString(`</div>`)

	return b.String()
}

// selectedAttr returns the selected attribute when active (for <option> tags).
func selectedAttr(active bool) string {
	if active {
		return ` selected`
	}
	return ""
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
.tc-select{font:inherit;font-size:.85rem;padding:.2rem .4rem;border:1px solid #ccc;border-radius:6px;background:#fff;color:#333;max-width:100%}
.tc-swatch-strip{display:inline-flex;height:18px;border:1px solid #ddd;border-radius:4px;overflow:hidden}
.tc-swatch-strip span{width:16px;height:100%}
.tc-detail-chart{max-width:760px;margin:1rem 0;border:1px solid #eee;border-radius:6px;padding:.5rem}
.tc-snippet-wrap{position:relative;max-width:760px}
.tc-copy{position:absolute;top:.5rem;right:.5rem;font:inherit;font-size:.75rem;padding:.15rem .55rem;border:1px solid #ccc;border-radius:6px;background:#fff;color:#333;cursor:pointer}
.tc-copy:hover{border-color:#888}
.tc-snippet{background:#f6f8fa;border:1px solid #e1e4e8;border-radius:6px;padding:1rem;overflow:auto;font-size:.85rem;line-height:1.4;margin:0}
</style>`
