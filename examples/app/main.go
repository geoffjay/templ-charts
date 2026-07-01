// Command templ-charts-demo runs the examples/app demo server: a stdlib
// net/http app serving the bar/line/pie/themes pages plus the htmx chart
// endpoints under /charts/.
//
// Run via `make run-demo` (or `go run ./examples/app`) → http://localhost:8000
package main

import (
	"log"
	"net/http"

	"github.com/geoffjay/templ-charts/examples/app/handlers"
)

func main() {
	app := handlers.NewApp()

	mux := http.NewServeMux()

	// HTMX chart endpoints (full render, hover, slice, click, toggle).
	mux.Handle("/charts/", app.Handler())

	// Page routes. Using exact-match guards so /bar doesn't shadow /bar/foo.
	mux.HandleFunc("/", app.Index)
	mux.HandleFunc("/bar", app.Bar)
	mux.HandleFunc("/line", app.Line)
	mux.HandleFunc("/pie", app.Pie)
	mux.HandleFunc("/heatmap", app.Heatmap)
	mux.HandleFunc("/waffle", app.Waffle)
	mux.HandleFunc("/calendar", app.Calendar)
	mux.HandleFunc("/radar", app.Radar)
	mux.HandleFunc("/radial-bar", app.RadialBar)
	mux.HandleFunc("/scatterplot", app.ScatterPlot)
	mux.HandleFunc("/stream", app.Stream)
	mux.HandleFunc("/bullet", app.Bullet)
	mux.HandleFunc("/funnel", app.Funnel)
	mux.HandleFunc("/boxplot", app.BoxPlot)
	mux.HandleFunc("/bump", app.Bump)
	mux.HandleFunc("/marimekko", app.Marimekko)
	mux.HandleFunc("/parallel-coordinates", app.ParallelCoordinates)
	mux.HandleFunc("/polar-bar", app.PolarBar)
	mux.HandleFunc("/treemap", app.Treemap)
	mux.HandleFunc("/sunburst", app.Sunburst)
	mux.HandleFunc("/icicle", app.Icicle)
	mux.HandleFunc("/circle-packing", app.CirclePacking)
	mux.HandleFunc("/tree", app.Tree)
	mux.HandleFunc("/palettes", app.Palettes)
	mux.HandleFunc("/themes", app.Themes)

	addr := ":8000"
	log.Printf("templ-charts demo listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
