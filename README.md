# templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart
concepts as [templ](https://github.com/a-h/templ) components generating
**server-side SVG**. It ships **twenty-eight** chart families — **bar**,
**line**, **pie**, **heatmap**, **waffle**, **calendar**, **radar**,
**radial-bar**, **scatterplot**, **stream**, **bullet**, **funnel**, **box
plot**, **bump**, **marimekko**, **parallel-coordinates**, **polar-bar**,
**treemap**, **sunburst**, **icicle**, **circle-packing**, **tree**,
**voronoi**, **network**, **swarmplot**, **sankey**, **chord**, and **geo**
(GeoMap + Choropleth) — a hybrid interactivity layer, responsive + accessible
output, and a runnable demo app. This is **full nivo SVG chart parity**: every
SVG chart type nivo ships has a templ-charts equivalent (Canvas rendering
remains out of scope).

- **Render** charts as SVG strings from Go — no Canvas, no JS bundle.
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

See [`docs/PLAN.md`](docs/PLAN.md) for the v1 design,
[`docs/PLAN-v2.md`](docs/PLAN-v2.md) for the v2 chart catalog + interactivity
and accessibility work, and [`docs/PLAN-v3.md`](docs/PLAN-v3.md) for the v3
push to full nivo SVG parity (the five deferred d3 ports and the fifteen
charts they unblock).

## Quickstart

```sh
make run-demo    # → http://localhost:8000
```

Browse a page per chart family — `/bar`, `/line`, `/pie`, `/heatmap`,
`/waffle`, `/calendar`, `/radar`, `/radial-bar`, `/scatterplot`, `/stream`,
`/bullet`, `/funnel`, `/boxplot`, `/bump`, `/marimekko`,
`/parallel-coordinates`, `/polar-bar`, `/treemap`, `/sunburst`, `/icicle`,
`/circle-packing`, `/tree`, `/voronoi`, `/network`, `/swarmplot`, `/sankey`,
`/chord`, `/geo` — plus `/palettes` and `/themes`. Hover any mark for a
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
does nearest-point hit-testing entirely in the browser. **State changes**
(series toggle, active-arc) stay server-side: register chart instances with
`charts/htmx.Registry` and mount `htmx.Handler` — see
[`examples/app`](examples/app) for a complete wiring.

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

Charts default to the `nivo` palette when `Colors` is left unset.

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

## Repository layout

```
charts/         library packages (mirror nivo names): core, theming, scales,
                colors, axes, rects, arcs, text, tooltip, legends,
                annotations, interact, static, grid, polar-axes, htmx, render
                (convenience String/To helpers), and the
                chart types: bar, line, pie, heatmap, waffle, calendar, radar,
                radialbar, scatterplot, stream, bullet, funnel, boxplot, bump,
                marimekko, parallelcoordinates, polarbar, treemap, sunburst,
                icicle, circlepacking, tree, voronoi, network, swarmplot,
                sankey, chord, geo
internal/d3/    pure-Go ports of d3-shape, d3-scale, d3-array, d3-format,
                d3-time-format, d3-color, d3-hierarchy, d3-delaunay, d3-force,
                d3-sankey, d3-chord, d3-geo
internal/golden small snapshot-test helper
examples/app/   runnable demo app (stdlib net/http)
docs/PLAN.md    v1 implementation plan
docs/PLAN-v2.md v2 plan (chart catalog, interactivity, a11y)
docs/PLAN-v3.md v3 plan (five d3 ports, fifteen charts, full SVG parity)
docs/NOTES.md   port-by-port implementation notes + deferred items
contrib/nivo/   upstream nivo clone (gitignored, reference only)
```

## Build & test

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

v3 complete: **full nivo SVG chart parity**. On top of v1's three charts and
v2's ten, v3 ports the five deferred d3 modules (hierarchy, delaunay, force,
sankey, chord, geo) plus `curveBumpX/Y`, and delivers the fifteen charts they
unblock — bump, marimekko, parallel-coordinates, polar-bar, treemap, sunburst,
icicle, circle-packing, tree, voronoi, network, swarmplot, sankey, chord, and
geo (GeoMap + Choropleth). d3-delaunay also retrofits true voronoi-mesh hover
into line/scatterplot/bump/swarmplot/tree. Every SVG chart type nivo ships now
has a templ-charts equivalent.

Still SVG-only — **Canvas rendering remains the largest deferred capability**
(the natural v4 theme). Other documented deferrals: geo's `clipCircle` /
`clipExtent` (azimuthal-family projections render the whole sphere) and
interactive zoom for icicle/treemap/circle-packing/sunburst. See
[`docs/PLAN-v3.md`](docs/PLAN-v3.md) §11 and [`docs/NOTES.md`](docs/NOTES.md)
for the full deferred list and port-by-port implementation notes.
