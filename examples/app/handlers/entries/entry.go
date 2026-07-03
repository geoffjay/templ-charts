// Package entries holds the per-chart detail-page definitions for the demo
// app. Each chart family gets its own file (bar.go, line.go, …) that registers
// one ChartEntry in init(); the handlers package looks them up by slug via Get
// and enumerates them via All.
package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ChartEntry is one chart family's detail-page definition: metadata, a
// copy-pasteable Go snippet, and a Render closure that owns the chart's
// concrete props and injects the selected theme + palette.
type ChartEntry struct {
	Slug        string
	Title       string
	Description string
	Snippet     string
	// Render builds the chart at the given theme (nil = default), palette
	// ("" = the chart's own default; ignored by charts whose colors are not
	// palette-driven), and animate flag (true enables the chart's SMIL enter
	// animation; default off) and returns the SVG string.
	Render func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error)
}

// registry maps slug → entry, populated by each chart file's init().
var registry = map[string]ChartEntry{}

// register adds an entry. Called from init() in the per-chart files; panics on
// a duplicate slug so collisions surface at startup.
func register(e ChartEntry) {
	if _, dup := registry[e.Slug]; dup {
		panic("entries: duplicate slug " + e.Slug)
	}
	registry[e.Slug] = e
}

// Get returns the entry for slug.
func Get(slug string) (ChartEntry, bool) {
	e, ok := registry[slug]
	return e, ok
}

// All returns all registered entries keyed by slug.
func All() map[string]ChartEntry { return registry }
