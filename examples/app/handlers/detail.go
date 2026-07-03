package handlers

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
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

	svg, err := entry.Render(theme, palette, animate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buildDetailHTML(entry, svg, themeName, palette, animate)
	a.renderPage(w,
		templates.LayoutProps{Title: entry.Title + " — templ-charts demo"},
		templ.Raw(html))
}

// buildDetailHTML assembles the detail page body: switchers, the full-width
// chart, and the code snippet.
func buildDetailHTML(e entries.ChartEntry, svg, themeName string, palette colors.PaletteID, animate bool) string {
	var b strings.Builder

	// animateSuffix carries the current animate selection through the theme and
	// palette hrefs so toggling one control does not reset the others. Empty
	// when animate is off, so those hrefs (and the rendered page) are unchanged
	// from the pre-animate behavior.
	animateSuffix := ""
	if animate {
		animateSuffix = "&animate=1"
	}

	b.WriteString(detailCSS)
	fmt.Fprintf(&b, `<p><a href="/">← all charts</a></p>`)
	fmt.Fprintf(&b, `<h2>%s</h2><p>%s</p>`, html.EscapeString(e.Title), html.EscapeString(e.Description))

	// Theme switcher.
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">theme</span>`)
	for _, g := range demos.ThemeGroups() {
		href := fmt.Sprintf("/chart/%s?theme=%s&palette=%s%s", e.Slug, g.Name, palette, animateSuffix)
		b.WriteString(switchLink(href, g.Name, g.Name == themeName))
	}
	b.WriteString(`</div>`)

	// Palette switcher (a "default" reset plus the curated categorical set).
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">palette</span>`)
	b.WriteString(switchLink(fmt.Sprintf("/chart/%s?theme=%s%s", e.Slug, themeName, animateSuffix), "default", palette == ""))
	for _, pc := range paletteChoices {
		href := fmt.Sprintf("/chart/%s?theme=%s&palette=%s%s", e.Slug, themeName, pc.ID, animateSuffix)
		b.WriteString(switchLink(href, pc.Label, palette == pc.ID))
	}
	b.WriteString(`</div>`)

	// Animate switcher: off (no animate param) / on (&animate=1). Both hrefs
	// preserve the current theme + palette.
	base := fmt.Sprintf("/chart/%s?theme=%s", e.Slug, themeName)
	if palette != "" {
		base = fmt.Sprintf("/chart/%s?theme=%s&palette=%s", e.Slug, themeName, palette)
	}
	b.WriteString(`<div class="tc-switch"><span class="tc-switch-label">animate</span>`)
	b.WriteString(switchLink(base, "off", !animate))
	b.WriteString(switchLink(base+"&animate=1", "on", animate))
	b.WriteString(`</div>`)

	// Full-width chart.
	b.WriteString(`<div class="tc-detail-chart chart">` + svg + `</div>`)

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
