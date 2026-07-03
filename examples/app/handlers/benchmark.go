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

	a.renderPage(w,
		templates.LayoutProps{Title: "benchmark — templ-charts demo", Nav: "benchmark"},
		templ.Raw(b.String()))
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
