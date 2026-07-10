package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/examples/app/templates"
)

// benchmarkSizes is the sweep of dataset sizes the /benchmark page renders.
var benchmarkSizes = []int{10, 50, 100, 500, 1000, 2000}

// Benchmark handles GET /benchmark: a load/showcase page that renders a bar
// chart at increasing dataset sizes and reports the server-side render time and
// SVG payload size for each, illustrating how the library scales under load.
// Timings are measured live per request (averaged over a few iterations), so
// they reflect the machine serving the page.
func (a *App) Benchmark(w http.ResponseWriter, r *http.Request) {
	const iters = 5

	var rows strings.Builder
	var showcase string
	for _, n := range benchmarkSizes {
		props := benchmarkBarProps(n)

		// Warm once (also gives us the payload size), then time the average.
		svg, err := render.String(bar.Bar(props))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		start := time.Now()
		for i := 0; i < iters; i++ {
			if _, err := render.String(bar.Bar(props)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		avgMs := float64(time.Since(start).Microseconds()) / float64(iters) / 1000.0

		fmt.Fprintf(&rows,
			`<tr><td style="text-align:right">%d</td><td style="text-align:right">%d</td><td style="text-align:right">%.3f</td></tr>`,
			n, len(svg), avgMs)
		if n == 500 {
			showcase = svg
		}
	}

	var b strings.Builder
	b.WriteString(`<h2>Render benchmark</h2>`)
	b.WriteString(`<p>Each row renders a stacked bar chart (4 keys per index) at the given number of index groups, ` +
		`averaged over ` + fmt.Sprint(iters) + ` server-side renders via <code>render.String</code>. ` +
		`Run <code>make bench</code> for the Go micro-benchmarks (d3 layout ports + chart render paths).</p>`)
	b.WriteString(`<table style="border-collapse:collapse" cellpadding="6" border="1">`)
	b.WriteString(`<thead><tr><th>index groups (N)</th><th>SVG bytes</th><th>render (ms, avg)</th></tr></thead><tbody>`)
	b.WriteString(rows.String())
	b.WriteString(`</tbody></table>`)
	b.WriteString(`<h3>Showcase — N = 500</h3>`)
	b.WriteString(`<div class="chart">` + showcase + `</div>`)

	if err := writeCanvasShowcase(&b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.renderPage(w,
		templates.LayoutProps{Title: "benchmark — templ-charts demo", Nav: "benchmark"},
		templ.Raw(b.String()))
}

// canvasShowcaseSizes is the large-N sweep for the SVG-vs-Canvas comparison.
var canvasShowcaseSizes = []int{1000, 5000, 20000}

// writeCanvasShowcase appends the large-N Canvas story: a table contrasting the
// SVG and Canvas backends' payload size and server render time for a
// scatterplot as the point count grows, plus one live Canvas scatterplot
// rendering thousands of points into a single <canvas>.
func writeCanvasShowcase(b *strings.Builder) error {
	const iters = 3

	b.WriteString(`<h2>Canvas backend — large N</h2>`)
	b.WriteString(`<p>The Canvas engine (<code>Render: theming.EngineCanvas</code>) draws the data marks into a ` +
		`single <code>&lt;canvas&gt;</code> draw-list instead of one SVG node per point, while grid and axes ` +
		`stay SVG. Below, the same scatterplot is rendered with each backend at increasing point counts — ` +
		`compare the payload size and server render time.</p>`)
	b.WriteString(`<table style="border-collapse:collapse" cellpadding="6" border="1">`)
	b.WriteString(`<thead><tr><th>points (N)</th><th>SVG bytes</th><th>Canvas bytes</th>` +
		`<th>SVG render (ms)</th><th>Canvas render (ms)</th></tr></thead><tbody>`)

	for _, n := range canvasShowcaseSizes {
		svgProps := showcaseScatterProps(n, theming.EngineSVG)
		canvasProps := showcaseScatterProps(n, theming.EngineCanvas)

		svg, err := render.String(scatterplot.ScatterPlot(svgProps))
		if err != nil {
			return err
		}
		cv, err := render.String(scatterplot.ScatterPlot(canvasProps))
		if err != nil {
			return err
		}
		svgMs := timeRender(iters, func() error {
			_, e := render.String(scatterplot.ScatterPlot(svgProps))
			return e
		})
		cvMs := timeRender(iters, func() error {
			_, e := render.String(scatterplot.ScatterPlot(canvasProps))
			return e
		})
		fmt.Fprintf(b,
			`<tr><td style="text-align:right">%d</td><td style="text-align:right">%d</td>`+
				`<td style="text-align:right">%d</td><td style="text-align:right">%.3f</td>`+
				`<td style="text-align:right">%.3f</td></tr>`,
			n, len(svg), len(cv), svgMs, cvMs)
	}
	b.WriteString(`</tbody></table>`)

	// Live Canvas showcase: 5000 points in one <canvas>.
	live, err := render.String(scatterplot.ScatterPlot(showcaseScatterProps(5000, theming.EngineCanvas)))
	if err != nil {
		return err
	}
	b.WriteString(`<h3>Live Canvas scatterplot — N = 5000</h3>`)
	b.WriteString(`<p>Five thousand points painted into a single canvas element (view source: one ` +
		`<code>&lt;canvas&gt;</code> + a JSON draw-list, no per-point DOM).</p>`)
	b.WriteString(`<div class="chart">` + live + `</div>`)
	return nil
}

// timeRender averages fn over iters runs and returns the mean in milliseconds.
func timeRender(iters int, fn func() error) float64 {
	start := time.Now()
	for i := 0; i < iters; i++ {
		if err := fn(); err != nil {
			return 0
		}
	}
	return float64(time.Since(start).Microseconds()) / float64(iters) / 1000.0
}

// showcaseScatterProps builds a deterministic single-series scatterplot of n
// points for the given backend.
func showcaseScatterProps(n int, engine theming.Engine) scatterplot.ScatterPlotProps {
	data := make([]scatterplot.ScatterPlotDatum, n)
	for i := 0; i < n; i++ {
		data[i] = scatterplot.ScatterPlotDatum{
			X: float64((i*37)%1000) + 1,
			Y: float64((i*53)%1000) + 1,
		}
	}
	return scatterplot.ScatterPlotProps{
		Width: 900, Height: 460,
		Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data:        []scatterplot.ScatterPlotSerie{{ID: "points", Data: data}},
		EnableGridX: core.BoolPtr(true),
		EnableGridY: core.BoolPtr(true),
		NodeSize:    4,
		Render:      engine,
		ChartID:     "benchmark-canvas-scatter",
		Responsive:  true,
	}
}

// benchmarkBarProps builds a deterministic stacked-bar props value with n index
// groups and four keys.
func benchmarkBarProps(n int) bar.BarProps {
	keys := []string{"alpha", "beta", "gamma", "delta"}
	data := make([]bar.BarDatum, n)
	for i := 0; i < n; i++ {
		row := bar.BarDatum{"idx": fmt.Sprintf("g%d", i)}
		for j, k := range keys {
			row[k] = float64((i*7+j*13)%100 + 10)
		}
		data[i] = row
	}
	return bar.BarProps{
		Width:      900,
		Height:     400,
		IndexBy:    "idx",
		Keys:       keys,
		Data:       data,
		Margin:     core.Margin{Top: 20, Right: 20, Bottom: 40, Left: 50},
		Responsive: true,
	}
}
