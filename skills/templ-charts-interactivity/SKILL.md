---
name: templ-charts-interactivity
description: Add interactivity to templ-charts charts — hover tooltips, line crosshair/slices, HTMX series toggle and hierarchy zoom, hover-highlight for chord/sankey/network, SMIL animation, resize re-fetch, and the Canvas backend for large datasets. Use when making a templ-charts chart interactive, wiring charts/htmx endpoints, or rendering big scatterplots/heatmaps.
---

# templ-charts interactivity

templ-charts uses a **hybrid model**; everything is opt-in (defaults off, so
plain renders stay zero-JS and byte-stable):

| Want | Mechanism | What to do |
|---|---|---|
| Hover tooltips, crosshair | client-side, `data-tc-*` attrs | `Interactive: true` + `@interact.ScriptTag()` |
| Hover-highlight connected marks | pure scoped CSS `:has()` | `Interactive: true` on chord/sankey/network |
| Series toggle, active arc, hierarchy zoom | HTMX server fragments | `charts/htmx` Registry + Handler + Mount |
| Enter animation | SMIL, no JS | `Animate: true` (+ `MotionStagger` where offered) |
| Pixel-accurate re-render on resize | ResizeObserver re-fetch | `data-tc-observe="<url>"` on a container |
| 10k+ point scatterplot/heatmap | Canvas draw-list | `Render: theming.EngineCanvas` + `@canvas.CanvasScriptTag()` |

## Client-side hover (no server round-trip)

Set `Interactive: true` on the chart props and load the dependency-free
script **once per page** (layout `<head>` or before `</body>`):

```go
import "github.com/geoffjay/templ-charts/charts/interact"

// in a templ layout:
@interact.ScriptTag()
```

Marks then emit `data-tc-tooltip`; the script positions tooltips, tracks the
line crosshair, and does nearest-point hit-testing in the browser.

**Line hover modes** (`line.LineProps`, all gated by `Interactive: true`):

- `UseMesh: true` — nearest-point voronoi-mesh hover + crosshair, fully
  client-side (`data-tc-mesh` JSON). This is the modern default path.
- `EnableSlices: line.EnableSlicesX` (or `EnableSlicesY`) — slice hover via a
  server fragment (`GET /charts/{id}/slice?...`), needs the HTMX wiring below.
- `ServerHover: true` — legacy per-mousemove server round-trip (avoid).
  `ClientHover` is a retained no-op.

**Hover-highlight** (chord, sankey, network): with `Interactive: true`,
hovering a node/arc/link dims everything else and re-lights the hovered
element plus its connected elements — scoped CSS `:has()`, no JS. Tune with
`NodeHoverOpacity` / `NodeHoverOthersOpacity` / `LinkHoverOpacity` /
`LinkHoverOthersOpacity` (sankey; network has the two `*HoverOpacity` props).

## Server-side state changes (HTMX)

Anything that changes *what is rendered* — legend series toggle, bar/pie
activation, hierarchy zoom — goes through `charts/htmx`. Wiring:

```go
import "github.com/geoffjay/templ-charts/charts/htmx"

registry := htmx.NewRegistry()
handler  := htmx.NewHandler(registry)

// Register each interactive chart instance once, with a stable id.
registry.RegisterBar("revenue", bar.BarProps{ /* Interactive: true, … */ })
// Typed helpers: RegisterBar/Line/Pie/Heatmap/Icicle/Treemap/CirclePacking/Sunburst,
// or the generic Register(id, kind, props).

mux := http.NewServeMux()
mux.Handle("/charts/", handler)   // hover/slice/toggle/click/zoom endpoints
```

In the page, render the initial SVG **through the handler** (so registered
state like hidden series is applied) and mount it:

```go
svg, err := handler.RenderFull("revenue")
```

```templ
@htmx.Mount(htmx.MountProps{ID: "revenue", SVG: svg, Interactive: true})
@interact.ScriptTag()
```

`Mount` emits the `#chart-<id>` container and `#tooltip-<id>` target that the
chart's `hx-*` attributes point at. The page must load HTMX itself (CDN
`<script>`); it is not a Go dependency.

Handler routes (all under the mounted prefix): `GET /charts/{id}` (full SVG),
`/hover`, `/slice`, `/click`, `POST /charts/{id}/toggle?series=…`,
`GET /charts/{id}/zoom?node=…`.

**Hierarchy zoom** (icicle, treemap, circlepacking, sunburst): set
`EnableZooming: true` on the props and register the chart; clicking a node
re-renders focused on its subtree with a breadcrumb back to root.

A complete working wiring lives in the repo's `examples/app`
(`make run-demo` → http://localhost:8000).

## Animation (SMIL)

`Animate: true` (default off) emits native SMIL `<animate>` enter transitions
(600ms); charts with many marks (calendar, sankey, marimekko, tree,
circlepacking, …) also accept `MotionStagger float64` — seconds of delay per
successive item. No JS involved.

## Resize re-fetch

`Responsive: true` handles cosmetic fluid scaling. For axis-dense charts that
should re-render at the new size, add to a container:

```html
<div data-tc-observe="/fragments/revenue"></div>
```

The interact script debounces resize and re-fetches the URL with `?w=&h=`
appended (uses `window.htmx.ajax` when present, else `fetch`). Optional
`data-tc-observe-target="<selector>"` swaps into a different element.

## Canvas backend (large N)

Scatterplot and heatmap can render marks into a `<canvas>` draw-list instead
of per-mark SVG (~4.4× smaller payload at 20k points):

```go
import "github.com/geoffjay/templ-charts/charts/theming"

scatterplot.ScatterPlot(scatterplot.ScatterPlotProps{
    Width: 900, Height: 500,
    Data:    bigSeries,
    Render:  theming.EngineCanvas,   // default is EngineSVG
    ChartID: "scatter1",             // required stable id
})
```

Load the replay script once per page, alongside the interact script:

```go
import "github.com/geoffjay/templ-charts/charts/canvas"
// in the layout:
@canvas.CanvasScriptTag()
```

Grid renders behind and axes/legends in front as SVG panes; hover uses the
mesh overlay (there is no per-mark DOM). Output is deterministic JSON, so it
golden-tests like SVG.
