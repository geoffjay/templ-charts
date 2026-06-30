# templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart
concepts as [templ](https://github.com/a-h/templ) components generating
**server-side SVG**. It ships **thirteen** chart families — **bar**, **line**,
**pie**, **heatmap**, **waffle**, **calendar**, **radar**, **radial-bar**,
**scatterplot**, **stream**, **bullet**, **funnel**, and **box plot** — a
hybrid interactivity layer, responsive + accessible output, and a runnable
demo app.

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
  d3-color to pure Go under `internal/d3/` (golden-tested against d3 output).

See [`docs/PLAN.md`](docs/PLAN.md) for the v1 design and
[`docs/PLAN-v2.md`](docs/PLAN-v2.md) for the v2 chart catalog + interactivity
and accessibility work.

## Quickstart

```sh
make run-demo    # → http://localhost:8080
```

Browse a page per chart family — `/bar`, `/line`, `/pie`, `/heatmap`,
`/waffle`, `/calendar`, `/radar`, `/radial-bar`, `/scatterplot`, `/stream`,
`/bullet`, `/funnel`, `/boxplot` — plus `/palettes` and `/themes`. Hover any
mark for a tooltip (client-side); click a legend item to toggle a series
(HTMX). The `/palettes` page is the full color-palette catalog applied to bars.

## Usage

```go
import (
    "context"
    "github.com/geoffjay/templ-charts/charts/bar"
)

func render() (string, error) {
    var b strings.Builder
    err := bar.Bar(bar.BarProps{
        Width: 700, Height: 400,
        IndexBy: "country",
        Keys:    []string{"hot dogs", "burgers"},
        Data:    data,
    }).Render(context.Background(), &b)
    return b.String(), err
}
```

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
                annotations, interact, static, grid, polar-axes, htmx, and the
                chart types: bar, line, pie, heatmap, waffle, calendar, radar,
                radialbar, scatterplot, stream, bullet, funnel, boxplot
internal/d3/    pure-Go ports of d3-shape, d3-scale, d3-array, d3-format,
                d3-time-format, d3-color
internal/golden small snapshot-test helper
examples/app/   runnable demo app (stdlib net/http)
docs/PLAN.md    v1 implementation plan
docs/PLAN-v2.md v2 plan (chart catalog, interactivity, a11y)
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
- Golden **SVG snapshot** tests for every chart type
  (`charts/{bar,line,pie,heatmap,waffle,calendar,radar,radialbar,scatterplot,stream,bullet,funnel,boxplot}/testdata/golden/`).
- Golden **path-string** tests for `GenerateSvgArc`, `BuildRoundedRectPath`,
  and `d3.Shape.Arc` (the arc-with-cornerRadius port is validated against
  real d3-shape output).
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

v2 complete: thirteen chart families, the hybrid client/server interactivity
layer, responsive scaling, and the accessibility pass. The v1 scaffolds
(`grid`, `polar-axes`) are now activated by the heatmap/waffle/calendar and
radar/radial-bar charts. Still SVG-only (Canvas remains out of scope). See
[`docs/PLAN-v2.md`](docs/PLAN-v2.md) §11 for explicitly deferred charts
(treemap, sunburst, sankey, chord, geo, voronoi-mesh hover, …) and the
optional `ResizeObserver` re-fetch.