# templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart concepts
as [templ](https://github.com/a-h/templ) components generating **server-side
SVG**. It ships **twenty-eight** chart families — **bar**, **line**, **pie**,
**heatmap**, **waffle**, **calendar**, **radar**, **radial-bar**,
**scatterplot**, **stream**, **bullet**, **funnel**, **box plot**, **bump**,
**marimekko**, **parallel-coordinates**, **polar-bar**, **treemap**,
**sunburst**, **icicle**, **circle-packing**, **tree**, **voronoi**,
**network**, **swarmplot**, **sankey**, **chord**, and **geo** (GeoMap +
Choropleth) — a hybrid interactivity layer, responsive + accessible output, and
a runnable demo app. This is **full nivo SVG chart parity**: every SVG chart
type nivo ships has a templ-charts equivalent. An opt-in **Canvas backend**
additionally renders scatterplot and heatmap into a `<canvas>` draw-list for
large-N datasets.

- **Render** charts as SVG strings from Go — no JS bundle required. An opt-in
  Canvas backend (`Render: theming.EngineCanvas`) is available for large-N
  scatterplot/heatmap.
- **Interact** with a hybrid model: ephemeral hover (tooltips, crosshair,
  nearest-point hit-testing) runs client-side from `data-*` attributes via a
  tiny dependency-free script (`charts/interact`); state changes (series
  toggle, active-arc) use [HTMX](https://htmx.org) server fragments.
- **Scale** fluidly with the `Responsive` prop (viewBox + `width:100%`, zero
  JS), and stay **accessible** with `role="img"`, `<title>`/`<desc>`, and
  ARIA labels on every chart.
- **Theme** with a nivo-faithful theme model (default, dark, custom).
- **Port** d3-shape / d3-scale / d3-array / d3-format / d3-time-format /
  d3-color / d3-hierarchy / d3-delaunay / d3-force / d3-sankey / d3-chord /
  d3-geo to pure Go under `internal/d3/` (golden-tested against d3 output).

See [`docs/USAGE.md`](docs/USAGE.md) for the full consumer guide.

## Quickstart

```sh
make run-demo    # → http://localhost:8000
```

Browse a page per chart family — `/bar`, `/line`, `/pie`, `/heatmap`,
`/waffle`, `/calendar`, `/radar`, `/radial-bar`, `/scatterplot`, `/stream`,
`/bullet`, `/funnel`, `/boxplot`, `/bump`, `/marimekko`,
`/parallel-coordinates`, `/polar-bar`, `/treemap`, `/sunburst`, `/icicle`,
`/circle-packing`, `/tree`, `/voronoi`, `/network`, `/swarmplot`, `/sankey`,
`/chord`, `/geo` — plus `/styling` (gradients, patterns, match rules),
`/legends`, `/composition` (build-your-own charts from the `Use*` hooks),
`/dashboard` (a composed dark-theme dashboard), `/palettes`, and `/themes`.
Hover any mark for a
tooltip (client-side); click a legend item to toggle a series (HTMX). The
`/palettes` page is the full color-palette catalog applied to bars.

## Usage

Every chart is a [`templ.Component`](https://templ.guide). The `charts/render`
helpers turn one into an SVG string (or write it to an `io.Writer`) in a single
call:

```go
import (
    "github.com/geoffjay/templ-charts/charts/bar"
    "github.com/geoffjay/templ-charts/charts/render"
)

func chart() (string, error) {
    return render.String(bar.Bar(bar.BarProps{
        Width: 700, Height: 400,
        IndexBy: "country",
        Keys:    []string{"hot dogs", "burgers"},
        Data:    data,
    }))
}
```

`render.To(w, component)` writes directly to an `io.Writer` (e.g. an
`http.ResponseWriter`); `render.StringCtx`/`render.ToCtx` take an explicit
`context.Context`. Since each chart is a plain `templ.Component`, you can also
render it yourself with `component.Render(ctx, w)` and embed it in a larger
templ page.

See [`docs/USAGE.md`](docs/USAGE.md) for the full consumer guide (props anatomy,
the color-config shapes, interactivity, responsive + a11y, theming), and the
runnable `ExampleXxx` in every `charts/<chart>` package.

### Interactivity

Hover interactions are **client-side**: set `Interactive: true` on a chart's
props (or `ClientHover: true` for `line` mesh/slice) so each mark emits a
`data-tc-tooltip`, and load the script once per page:

```go
import "github.com/geoffjay/templ-charts/charts/interact"
// in your layout <head> or before </body>:
@interact.ScriptTag()
```

The script (no dependencies) shows tooltips, tracks the line crosshair, and
does nearest-point hit-testing entirely in the browser. For `line`, the client
mesh/slice path is now the **default** (the legacy per-mousemove server
round-trip is opt-in via `ServerHover`). **State changes** (series toggle,
active-arc, hierarchy zoom) stay server-side: register chart instances with
`charts/htmx.Registry` and mount `htmx.Handler` — see
[`examples/app`](examples/app) for a complete wiring.

**Hover-highlight** (`chord`, `sankey`, `network`): with `Interactive: true`,
hovering a node/arc dims the rest and re-lights it plus its connected
elements — pure scoped CSS (`:has()`), no JS or server round-trip. The
`*HoverOpacity` / `*HoverOthersOpacity` props tune the lit/dimmed opacities.

**Hierarchy zoom** (`icicle`, `treemap`, `circle-packing`, `sunburst`): set
`EnableZooming: true` and register the chart with the `htmx.Registry`; clicking
a node re-renders focused on its subtree (with a breadcrumb back to the root),
via an HTMX full-SVG swap.

**Resize re-fetch** (opt-in): add `data-tc-observe="<url>"` to a container and
the script re-fetches it (with the new `?w=&h=`) on resize for a pixel-accurate
re-render of axis-dense charts. Cosmetic fluid scaling stays handled by
`Responsive`.

### Animation

Every chart accepts `Animate bool` (default off, so static output is
byte-stable) plus `MotionStagger float64`. When on, marks play a SMIL enter
transition — fade-in for rects/arcs/lines/cells, radius-scale for circles —
staggered by `MotionStagger` seconds. No JS: the animation is native SMIL
`<animate>` in the SVG.

### Canvas backend (large N)

Scatterplot and heatmap can render into a `<canvas>` instead of one SVG node per
mark — for datasets where thousands of DOM elements are too many. Set
`Render: theming.EngineCanvas` and give the chart a stable `ChartID`:

```go
import (
    "github.com/geoffjay/templ-charts/charts/scatterplot"
    "github.com/geoffjay/templ-charts/charts/theming"
)

scatterplot.ScatterPlot(scatterplot.ScatterPlotProps{
    Width: 900, Height: 500,
    Data:    bigSeries,             // thousands of points
    Render:  theming.EngineCanvas,  // draw marks into a <canvas>
    ChartID: "scatter1",            // unique per canvas on the page
})
```

The chart emits a wrapper `<div>` layering the marks (a `<canvas>` + a JSON
draw-list) between SVG panes for grid (below) and axes/legends (above), so it
lines up exactly with the SVG version. Include the replay script once per page —
`@canvas.CanvasScriptTag()` (a sibling of `interact.ScriptTag()`) — to paint the
draw-list; it HiDPI-scales and repaints on resize, with no server round-trip.
The default engine stays SVG, so existing output is unchanged. Hover on Canvas
goes through the `UseMesh` overlay (there is no per-mark DOM to hover).

### Accessibility & responsiveness

Every chart accepts `Responsive bool` (fluid `viewBox` scaling, zero JS) and
the a11y props `Role` (defaults to `img`), `Title`/`Desc` (rendered as
`<title>`/`<desc>` for the accessible name + description), `AriaLabel`,
`AriaLabelledBy`, `AriaDescribedBy`, and `IsFocusable`.

## Color palettes

Charts color series via the `Colors` field on each chart's props. The
`charts/colors` package ships a broad catalog of named palettes —
**categorical**, **sequential**, and **diverging** — including several
colorblind-safe options. The ergonomic way to pick one is `colors.Scheme`:

```go
import "github.com/geoffjay/templ-charts/charts/colors"

bar.BarProps{
    // ...
    Colors: colors.Scheme(colors.PaletteTableau10), // typed, autocomplete-friendly
}
```

Other ways to set colors:

```go
colors.Scheme(colors.PaletteOkabeIto)        // a named palette (colorblind-safe)
colors.PaletteColors("#4269d0", "#efb118")   // an explicit custom color list
```

Or use the unified `colors.Set` entry point, which infers intent from its
argument and resolves to whichever config shape a field expects:

```go
colors.Set(colors.PaletteTableau10).Ordinal()          // categorical `Colors`
colors.Set([]string{"#f00", "#0f0"}).Ordinal()          // explicit list
colors.Set(colors.PaletteViridis).InSpace(colors.SpaceLab).Sequential() // perceptual
colors.Set("#333").Inherited()                          // BorderColor/LabelTextColor
```

**Perceptual interpolation.** Sequential/diverging scales and palette gradient
sampling interpolate in RGB by default, but accept a `Space` selector for
perceptually-uniform ramps (`colors.SpaceLab` / `colors.SpaceLch`); the default
RGB behavior is unchanged. `internal/d3/color` implements HSL, Lab, and Lch
next to RGB.

Charts default to the `nivo` palette when `Colors` is left unset.

**Sample data.** The `charts/samples` package returns ready-to-render, typed,
deterministic data for every chart family (`samples.Bar()`, `samples.Chord()`,
…) — a fast way to try a chart before wiring your own data.

Enumerate the catalog at runtime (for pickers, docs, galleries):

```go
colors.Palettes()                          // full ordered catalog with metadata
colors.PalettesByKind(colors.KindSequential)
colors.ColorblindSafePalettes()
p, ok := colors.LookupPalette(colors.PaletteSunset)
swatch := p.Swatch(8)                      // preview colors (samples gradients)
```

See [`docs/PALETTES.md`](docs/PALETTES.md) for the full list of palette ids,
and the `/palettes` page in the demo app for a visual gallery.

## Agent skills

The repo ships [agent skills](https://skills.sh) under [`skills/`](skills)
that teach coding agents (Claude Code, Cursor, …) how to use the library.
Install them into your own project with:

```sh
npx skills add geoffjay/templ-charts
```

Four skills are included: `templ-charts` (core usage + a data-shape reference
for all 28 chart families), `templ-charts-interactivity` (hover, HTMX wiring,
zoom, Canvas backend), `templ-charts-theming` (colors, palettes, themes,
legends, gradients/patterns), and `templ-charts-composition` (build-your-own
charts from the `Use*` hooks). Every code snippet in the skills is
compile-checked against the library.

## Build and Test

| Task | Command |
|---|---|
| Generate templ sources | `make templ` |
| Build all packages | `make build` |
| Run tests | `make test` |
| Lint (vet + fmt) | `make lint` |
| Run tests with coverage | `make cover` |
| Regenerate golden snapshots | `make golden` |
| Run benchmarks (d3 ports + render paths) | `make bench` |
| CI (lint + test) | `make ci` |
| Run the demo app | `make run-demo` |
| Tidy modules | `make tidy` |

Golden SVG snapshots live under each chart package's `testdata/golden/`.
After an **intentional** render change, regenerate and commit the updated
snapshots:

```sh
make golden
```

## Testing

- Unit tests per foundational package (scales, colors, theming, axes, arcs)
  assert geometry and color output against known d3 values.
- Golden **SVG snapshot** tests for every chart type, under each package's
  `testdata/golden/`.
- Golden **layout / path-string** tests for the ported geometry: `GenerateSvgArc`,
  `BuildRoundedRectPath`, `d3.Shape.Arc`, the `curveBumpX/Y` bump curves, and
  the new d3 ports — hierarchy node coordinates (treemap tiling, pack circles,
  partition rects, tidy-tree positions), delaunay cells + `Find`, sankey node
  rects + link widths, chord group angles + ribbon paths, force positions after
  a fixed `Tick(120)`, and geo projected paths per projection — validated
  against real d3 output.
- `charts/interact` covers the client hover script; `charts/htmx` and
  `examples/app/handlers` cover the HTMX endpoints and page rendering.

Run everything:

```sh
make ci
```

## Dependencies

- [`github.com/a-h/templ`](https://github.com/a-h/templ) — templ
  compiler/runtime (the only Go dependency).
- [htmx.org](https://htmx.org) — loaded via CDN `<script>` in the demo (no
  Go dep).

## Status

templ-charts has **full nivo SVG chart parity** — all twenty-eight SVG chart
families are implemented, backed by pure-Go ports of the d3 modules they need.
Beyond the core charts:

- **Perceptual color spaces** — `internal/d3/color` implements **HSL, Lab, and
  Lch** alongside RGB (faithful d3-color conversions + in-space interpolation).
  Sequential/diverging scales and palette sampling take an opt-in `Space`
  selector (default RGB), and lightness modifiers can apply in Lab. See the
  **color space** switcher on the heatmap detail page.
- **Consumption surface** — a public typed **`charts/samples`** package (one
  ready-to-render dataset per chart family), a **`charts/static`** registry
  covering **all 28** chart families through one reflective adapter, and a
  unifying **`colors.Set(...)`** API over the color-config shapes.
- **Geo completeness** — `GeoPath` **bounds/area/centroid**, projection
  **`FitExtent`/`FitSize`/`FitWidth`/`FitHeight`**, the **conic** projection
  family (conformal/equal-area/equidistant with standard parallels), and an
  opt-in **TopoJSON** decoder (GeoJSON is still the default input). See the
  auto-fit choropleth on the geo detail page.
- **Canvas backend (large N)** — scatterplot and heatmap can render into a
  `<canvas>` draw-list (~4.4× smaller payload than SVG at n=20000), backed by
  the `d3-quadtree`/Barnes–Hut and Delaunator ports.

The Canvas engine and the color-space, TopoJSON, and conic additions are all
opt-in and default to prior behavior, so static SVG output is byte-stable across
releases. See [`docs/USAGE.md`](docs/USAGE.md) for the consumer guide.
