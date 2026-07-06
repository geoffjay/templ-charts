---
name: templ-charts
description: Render server-side SVG charts in Go with templ-charts — nivo-style templ components covering 28 chart families (bar, line, pie, heatmap, treemap, sankey, geo, …). Use when adding a chart to a Go templ/HTMX app, choosing a chart's data shape, rendering a chart to an SVG string, or wiring sample data. For hover/HTMX interactivity see templ-charts-interactivity; for colors/themes/legends see templ-charts-theming; for custom charts see templ-charts-composition.
---

# templ-charts

Go library (`github.com/geoffjay/templ-charts`) that renders nivo-style charts
as **server-side SVG** via [templ](https://templ.guide) components. No JS
bundle is required — a plain render produces a static, byte-stable `<svg>`
string. The only Go dependency is `github.com/a-h/templ`.

```sh
go get github.com/geoffjay/templ-charts
```

## Rendering a chart

Every chart is a `templ.Component`: build a props struct, call the chart's
constructor, render. The `charts/render` helpers collapse this to one call:

```go
import (
    "github.com/geoffjay/templ-charts/charts/bar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, err := render.String(bar.Bar(bar.BarProps{
    Width: 700, Height: 400,
    IndexBy: "country",
    Keys:    []string{"hot dogs", "burgers"},
    Data: []bar.BarDatum{
        {"country": "USA", "hot dogs": 42.0, "burgers": 33.0},
        {"country": "Japan", "hot dogs": 28.0, "burgers": 51.0},
    },
}))
```

| Helper | Signature | Use |
|---|---|---|
| `render.String` | `(templ.Component) (string, error)` | SVG string |
| `render.To` | `(io.Writer, templ.Component) error` | write to e.g. `http.ResponseWriter` |
| `render.StringCtx` | `(context.Context, templ.Component) (string, error)` | string, with cancellation |
| `render.ToCtx` | `(context.Context, io.Writer, templ.Component) error` | writer, with cancellation |

Since charts are plain `templ.Component`s you can also embed them directly in
a templ page (`@bar.Bar(props)`) or call `component.Render(ctx, w)` yourself.

## Props anatomy (consistent across all 28 families)

- **Geometry** — `Width`, `Height float64`, `Margin core.Margin`.
- **Data + accessors** — a `Data` field plus per-chart accessors (`IndexBy`,
  `Keys`, `Value`, `ID`, …). Accessor fields are
  `core.PropertyAccessor[D, V]` and accept **either a string key or a
  `func(D) V`** (e.g. `IndexBy: "country"` or `IndexBy: func(d bar.BarDatum) string {…}`).
- **`Responsive bool`** — `viewBox` + `width:100%;height:auto` fluid scaling,
  zero JS. Default `false` (fixed pixels).
- **A11y** — `Role` (default `"img"`), `Title`/`Desc` (rendered as SVG
  `<title>`/`<desc>`), `AriaLabel`/`AriaLabelledBy`/`AriaDescribedBy`,
  `IsFocusable` (emits `tabindex="0"`).
- **`Interactive bool`** — opt-in hover tooltips (see the
  templ-charts-interactivity skill).
- **`Theme *theming.Theme`**, **`Colors …`** — see the templ-charts-theming
  skill.
- **`Layers []<Chart>LayerId`** — reorder/omit render layers.

**Zero-valued fields fall back to the package's exported `Defaults`** via an
internal `applyDefaults` — a minimal props struct renders sensibly. Bool
props are the exception (Go can't distinguish unset from false): opt-in bools
like `Interactive`, `Responsive`, `Animate` default off. Read a chart's
`defaults.go` or `go doc github.com/geoffjay/templ-charts/charts/<chart>` for
its defaults.

**Gotcha:** map-shaped data (`bar.BarDatum = map[string]any`, radar rows)
needs numeric values as `float64` — write `42.0`, not `42`.

## Choosing a chart & its data shape

Quick map (full per-family catalog with exact types, key props, and sample
functions: read [reference.md](reference.md)):

| Data looks like | Charts |
|---|---|
| rows of `map[string]any` + keys | bar, radar |
| `[]Series{ID, Data []{X, Y}}` | line, scatterplot, bump, radialbar, heatmap |
| flat `[]{ID, Value}` items | pie, funnel, waffle, polarbar, swarmplot, boxplot (Group/SubGroup/Value), calendar (Day/Value), bullet (Ranges/Measures/Markers), voronoi (ID/X/Y) |
| recursive node tree `{ID, Value, Children}` | treemap, sunburst, icicle, circlepacking, tree (no Value) |
| nodes + links | sankey, network |
| square matrix + keys | chord |
| map rows + keys (stacked) | stream, marimekko (+ Dimensions), parallelcoordinates (+ Variables) |
| GeoJSON features | geo.GeoMap, geo.Choropleth (+ `[]{ID, Value}`) |

## Sample data & the static registry

`charts/samples` ships deterministic, typed demo data for every family — the
fastest way to try a chart. Multi-input charts return extras:

```go
import "github.com/geoffjay/templ-charts/charts/samples"

data, keys := samples.Bar()          // ([]bar.BarDatum, []string)
series     := samples.Line()         // []line.LineSeries
nodes, links := samples.Sankey()     // ([]sankey.SankeyInputNode, []sankey.SankeyInputLink)
```

`charts/static` renders any family by id through one dispatcher (useful for
chart-type-as-data scenarios like dashboards):

```go
import "github.com/geoffjay/templ-charts/charts/static"

sample := static.Samples[static.ChartTypeBar]        // bundled demo props
svg, err := static.RenderChart(static.ChartTypeBar, sample.Props, map[string]any{
    "width": 640, "height": 400,                     // whitelisted overrides
})
```

`static.ChartType*` constants exist for all 28 families
(`ChartTypeBar` … `ChartTypeChoropleth`).

## Verifying output

- Rendered output starts with `<svg` and is deterministic — snapshot/golden
  tests work well.
- Every chart package has a runnable `ExampleXxx` in `example_test.go`
  (`go test -run Example ./charts/...`) — copy these as starting points.
- Full consumer guide: `docs/USAGE.md` in the repo
  (github.com/geoffjay/templ-charts).
