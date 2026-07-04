# Using templ-charts

A practical guide to consuming the library. For the design/rationale see
[`PLAN.md`](PLAN.md) (v1), [`PLAN-v2.md`](PLAN-v2.md) (v2),
[`PLAN-v3.md`](PLAN-v3.md) (v3), and [`PLAN-v4.md`](PLAN-v4.md) (v4). Every chart
package also ships a runnable `ExampleXxx` (visible on pkg.go.dev and under
`go test`).

## Install

```sh
go get github.com/geoffjay/templ-charts
```

The only runtime dependency is [`github.com/a-h/templ`](https://templ.guide).
HTMX (used only for server-driven state changes, see [Interactivity](#interactivity))
is loaded from a CDN `<script>` in the browser — it is not a Go dependency.

## Rendering a chart

Every chart is a [`templ.Component`](https://templ.guide): you build a props
value, call the chart's constructor, and render the component. The
`charts/render` helpers collapse the render into one call:

```go
import (
    "github.com/geoffjay/templ-charts/charts/bar"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, err := render.String(bar.Bar(bar.BarProps{
    Width: 700, Height: 400,
    IndexBy: "country",
    Keys:    []string{"hot dogs", "burgers"},
    Data:    data,
}))
```

The four helpers:

| Helper | Signature | Use |
|---|---|---|
| `render.String` | `(templ.Component) (string, error)` | render to an SVG string |
| `render.To` | `(io.Writer, templ.Component) error` | render into a writer (e.g. `http.ResponseWriter`) |
| `render.StringCtx` | `(context.Context, templ.Component) (string, error)` | string, with cancellation |
| `render.ToCtx` | `(context.Context, io.Writer, templ.Component) error` | writer, with cancellation |

Because each chart is a plain `templ.Component`, you can also render it yourself
(`component.Render(ctx, w)`) or embed it inside a larger templ page — the
helpers are a convenience, not a requirement.

## Anatomy of a chart's props

Props follow a consistent v1/v2/v3 shape. Common fields present on (nearly)
every chart:

- **Geometry** — `Width`, `Height float64` and `Margin core.Margin`.
- **Data + mapping** — a `Data` field plus accessor fields (`IndexBy`, `Keys`,
  `Value`, `ID`, …) that vary per chart.
- **`Responsive bool`** — fluid scaling (see [Responsive](#responsive)).
- **A11y** — `Role`, `AriaLabel`, `AriaLabelledBy`, `AriaDescribedBy`, `Title`,
  `Desc string` and `IsFocusable bool` (see [Accessibility](#accessibility)).
- **`Interactive bool`** — client-side hover tooltips (see [Interactivity](#interactivity)).
- **`Theme *theming.Theme`** — see [Theming](#theming).
- **`Colors …`** — see [Colors](#colors).
- **`Layers []…LayerId`** — reorder/omit the chart's render layers.

Zero-valued fields fall back to each chart's exported `Defaults` (or `props.go`
for pie) via an internal `applyDefaults`. To see a chart's defaults, read its
`defaults.go` or `go doc github.com/geoffjay/templ-charts/charts/<chart>`.

## Colors

Charts color their marks through a `Colors` field (default: the `nivo` scheme).
There are **five** color-config shapes depending on what the chart needs — pick
by chart family:

1. **Ordinal / categorical** — most charts (bar, line, pie, radar, scatterplot,
   funnel, treemap, sankey, chord, sunburst, …). Field type
   `colors.OrdinalColorScaleConfig`; the ergonomic way to set it:

   ```go
   import "github.com/geoffjay/templ-charts/charts/colors"

   Colors: colors.Scheme(colors.PaletteTableau10),      // a named palette
   Colors: colors.PaletteColors("#4269d0", "#efb118"),  // explicit colors
   ```

2. **Continuous (sequential/diverging)** — heatmap. Field type
   `heatmap.HeatMapColorConfig` (a value→color scale) plus `EmptyColor string`.

3. **Quantize palette** — calendar. `Colors []string` (value bucketed into these
   colors) plus `EmptyColor string`.

4. **Quantize scheme id + steps** — geo Choropleth: `Colors string` (a scheme id,
   e.g. `"purple_blue_green"`) + `Steps int` + `UnknownColor string`. geo GeoMap
   uses a single `FillColor string`.

5. **Plain color strings** — network: `NodeColor`, `NodeBorderColor`,
   `LinkColor string` (per-node/link overrides fall back to these).

Enumerate the palette catalog at runtime with `colors.Palettes()`,
`colors.PalettesByKind`, `colors.ColorblindSafePalettes`, and
`colors.LookupPalette`. See [`PALETTES.md`](PALETTES.md) for the full id list and
the demo `/palettes` page for a gallery.

## Interactivity

templ-charts uses a **hybrid** model, and both halves are opt-in (defaults are
off, so static embeds are zero-JS and byte-stable):

- **Client-side hover (ephemeral)** — set `Interactive: true` on a chart's props
  and each mark emits `data-tc-*` attributes; a tiny dependency-free script shows
  tooltips and does nearest-point hit-testing entirely in the browser. For
  `line`, the client mesh/slice path is the **default** (the legacy per-mousemove
  server round-trip is opt-in via `ServerHover`; `ClientHover` is a retained
  no-op). Load the script once per page:

  ```go
  import "github.com/geoffjay/templ-charts/charts/interact"
  // in your <head> or before </body>:
  @interact.ScriptTag()
  ```

- **Client-side hover-highlight (CSS)** — `chord`, `sankey` and `network` (with
  `Interactive: true`) dim the other marks and re-light the hovered one plus its
  connected elements, using a scoped CSS `:has()` block — no JS, no round-trip.
  The `*HoverOpacity` / `*HoverOthersOpacity` props tune the lit/dimmed values.

- **Client-side animation (SMIL)** — set `Animate: true` (default off) plus an
  optional `MotionStagger` for a native SMIL enter transition; no JS.

- **Server-side state changes (HTMX)** — actions that change *what is rendered*
  (series toggle, active-arc, and hierarchy **zoom** for icicle/treemap/
  circle-packing/sunburst via `EnableZooming`). Register a chart instance with
  `charts/htmx.Registry` and mount `htmx.Handler`; the components emit `hx-*`
  attributes pointing at its endpoints. See [`examples/app`](../examples/app) for
  a complete wiring, and `htmx.Mount` for the container plumbing.

- **Resize re-fetch (opt-in)** — add `data-tc-observe="<url>"` to a container for
  a debounced re-fetch (with the new `?w=&h=`) on resize, for pixel-accurate
  re-render of axis-dense charts. Cosmetic scaling stays handled by `Responsive`.

For a fully static, zero-JS embed, leave `Interactive` off and simply place the
rendered SVG string in your page.

## Responsive

Set `Responsive: true` and the chart keeps its `viewBox` and intrinsic
dimensions while emitting `style="width:100%;height:auto;display:block"`, so it
scales fluidly to its container with **zero JS and zero consumer CSS**. Defaults
to `false` (fixed pixel size), so existing fixed-size renders are unchanged.

## Canvas backend (large N)

For large datasets, scatterplot and heatmap can render into a single `<canvas>`
instead of one SVG element per mark. Opt in with `Render: theming.EngineCanvas`
and a stable `ChartID`:

```go
scatterplot.ScatterPlot(scatterplot.ScatterPlotProps{
    Width: 900, Height: 500,
    Data:    bigSeries,
    Render:  theming.EngineCanvas,
    ChartID: "scatter1",
})
```

How it works: the same `Use{Chart}` layout hooks that drive the SVG render feed
a **draw-list** — an ordered sequence of primitive ops (`fillStyle`,
`fillCircle`, `fillRect`, `fillText`, `line`, `fillPath`, …) recorded by
`charts/canvas`. The chart emits a wrapper `<div>` with the marks as a
`<canvas>` + a JSON draw-list, layered between SVG panes for grid (behind) and
axes/legends (in front). A tiny dependency-free replay script paints the
draw-list into the 2D context, HiDPI-scaled and repainted on resize — include it
once per page:

```go
@canvas.CanvasScriptTag()   // alongside interact.ScriptTag()
```

Notes:

- The default engine is SVG; Canvas is purely opt-in, so existing renders are
  byte-identical.
- The draw-list is deterministic data — it golden-tests as JSON, and a future
  server-side PNG rasterizer could consume the same list.
- Hover uses the `UseMesh` overlay (there is no per-mark DOM to hover on Canvas).
- At 20k points a scatterplot's payload is ~4.4× smaller than SVG (2.3 MB →
  0.5 MB) and emits ~1.8× faster; see the demo app's `/benchmark` page.

## Accessibility

Every chart renders `role="img"` by default and accepts:

- `Title` / `Desc` — rendered as `<title>`/`<desc>` (the SVG-native accessible
  name + description).
- `AriaLabel` / `AriaLabelledBy` / `AriaDescribedBy` — ARIA attributes.
- `IsFocusable` — makes the `<svg>` keyboard-focusable (`tabindex`).

Defaults are empty/off, so static goldens stay byte-stable until you opt in.

## Theming

Charts accept a `Theme *theming.Theme` (nil → `theming.DefaultTheme`). Build a
custom theme by extending the default:

```go
import "github.com/geoffjay/templ-charts/charts/theming"

dark := theming.ExtendDefaultTheme(theming.DefaultTheme, &theming.PartialTheme{
    Background: "#1a1a1a",
    // Text, Grid, Axis, Tooltip, … overrides
})

bar.Bar(bar.BarProps{ /* … */ Theme: &dark })
```

`ExtendDefaultTheme` deep-merges and resolves text-style inheritance. See the
demo `/themes` page for default/dark/custom side by side.

## Where to look next

- **Runnable examples** — `ExampleXxx` in every `charts/<chart>` package
  (pkg.go.dev or `go test -run Example ./charts/...`).
- **The demo app** — [`examples/app`](../examples/app): a page per chart family
  plus `/palettes` and `/themes`. Run with `make run-demo`.
- **Package docs** — `go doc github.com/geoffjay/templ-charts/charts/<chart>`.
