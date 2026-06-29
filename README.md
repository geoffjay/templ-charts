# templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart
concepts as [templ](https://github.com/a-h/templ) components generating
**server-side SVG**. v1 ships **bar**, **line**, and **pie** charts plus an
HTMX-backed interactivity layer and a runnable demo app.

- **Render** charts as SVG strings from Go — no Canvas, no JS bundle.
- **Interact** via [HTMX](https://htmx.org): hover tooltips, series toggle,
  active-arc highlight, and slice crosshairs are server-rendered fragments.
- **Theme** with a nivo-faithful theme model (default, dark, custom).
- **Port** d3-shape / d3-scale / d3-array / d3-format / d3-time-format /
  d3-color to pure Go under `internal/d3/` (golden-tested against d3 output).

See [`docs/PLAN.md`](docs/PLAN.md) for the full design.

## Quickstart

```sh
make run-demo    # → http://localhost:8080
```

Browse `/bar`, `/line`, `/pie`, and `/themes`. Hover a bar/arc/line slice for
a tooltip; click a legend item to toggle a series.

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

For interactivity, register chart instances with `charts/htmx.Registry` and
mount `htmx.Handler` — see [`examples/app`](examples/app) for a complete
wiring.

## Repository layout

```
charts/         library packages (mirror nivo names): core, theming, scales,
                colors, axes, rects, arcs, text, tooltip, legends,
                annotations, static, bar, line, pie, htmx, grid, polar-axes
internal/d3/    pure-Go ports of d3-shape, d3-scale, d3-array, d3-format,
                d3-time-format, d3-color
internal/golden small snapshot-test helper
examples/app/   runnable demo app (stdlib net/http)
docs/PLAN.md    full implementation plan
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
- Golden **SVG snapshot** tests for `Bar`, `Line`, `Pie`
  (`charts/{bar,line,pie}/testdata/golden/`).
- Golden **path-string** tests for `GenerateSvgArc`, `BuildRoundedRectPath`,
  and `d3.Shape.Arc` (the arc-with-cornerRadius port is validated against
  real d3-shape output).
- `charts/htmx` and `examples/app/handlers` cover the HTMX endpoints and
  page rendering.

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

v1 complete. See `docs/PLAN.md` §9 for the scope summary (scaffold-only:
`grid`, `polar-axes`; skipped: Canvas).