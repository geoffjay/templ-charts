// Command templ-charts-demo runs the examples/app demo server: a stdlib
// net/http app serving every chart family's demo page, the htmx chart
// endpoints under /charts/, and the per-chart detail pages under /chart/.
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

	addr := ":8000"
	log.Printf("templ-charts demo listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, app.Mux()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
