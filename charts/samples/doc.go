// Package samples provides public, typed, ready-to-render demo data for every
// templ-charts chart family. Each exported function returns the concrete data
// types the corresponding chart's Props expect — e.g. samples.Bar returns the
// []bar.BarDatum and the key list, samples.Chord returns the flow matrix and
// its labels — so a downstream user can try any chart without hand-authoring
// data first.
//
// Every generator is fully deterministic (no randomness, no wall-clock time),
// so renders built from these datasets are byte-stable and safe for golden
// tests. The datasets are intentionally small.
//
// Charts that need more than a single data slice to render (a key/dimension/
// variable list, a nodes+links pair, a features+values pair) return those
// extra inputs as additional results, in the order the chart consumes them.
// This package is the single source of truth for demo data: the example app
// and the charts/static registry both draw from it.
package samples
