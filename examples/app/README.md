# templ-charts demo app

A runnable demo of the `templ-charts` library: a stdlib `net/http` server
serving demos for all thirteen chart families with a hybrid interactivity
model — client-side hover (tooltips + crosshair) plus HTMX-backed state
changes (series toggle, active arc).

## Run

From the repo root:

```
make run-demo
```

or

```
go run ./examples/app
```

Then open <http://localhost:8080>.

## Pages

| Route          | Contents                                                     |
|----------------|--------------------------------------------------------------|
| `/`            | Index listing all demos                                      |
| `/bar`         | Stacked, grouped, markers+annotations, legend toggle, totals |
| `/line`        | Single, multi+legend, area+points, slices, mesh+crosshair    |
| `/pie`         | Plain, donut, half, sorted, active-arc hover, legend toggle  |
| `/heatmap`     | 2D value grid: sequential/diverging colors, labels, legend   |
| `/waffle`      | Part-of-whole cell grid (built on `charts/grid`)             |
| `/calendar`    | Day-grid heatmap over a date range, month/year legends       |
| `/radar`       | Polar line/area: circular/polygon grids, dots, legend        |
| `/radial-bar`  | Stacked polar arcs (`charts/polar-axes` + `charts/arcs`)     |
| `/scatterplot` | {x,y} nodes on linear scales, grid, axes, legend             |
| `/stream`      | Stacked areas (wiggle/silhouette/expand offsets)             |
| `/bullet`      | KPI ranges + measures + markers on a shared scale            |
| `/funnel`      | Smooth/linear trapezoid parts, separators, labels            |
| `/boxplot`     | Quantile box + whisker glyphs from raw observations          |
| `/palettes`    | The full color-palette catalog applied to bars               |
| `/themes`      | bar / line / pie under default, dark, and custom themes      |

## Interactivity

The interactivity is **hybrid** (see `docs/PLAN-v2.md` §5):

- **Client-side** (`charts/interact`, loaded once in the layout): ephemeral
  hover — tooltips, the line crosshair, and nearest-point hit-testing — runs
  entirely in the browser off `data-tc-*` attributes, with no server
  round-trip. The ten v2 charts and `line`'s mesh/slice hover use this path.
- **Server-side** (HTMX): *state changes* that alter what is rendered. Bar and
  pie demos register a chart instance with `charts/htmx.Registry` and mount the
  `htmx.Handler` at `/charts/`; their components emit `hx-*` attributes
  pointing at the handler's endpoints:

| Route                          | Method | Action                                   |
|--------------------------------|--------|------------------------------------------|
| `/charts/{id}`                 | GET    | Full SVG render                          |
| `/charts/{id}/hover?bar=k`     | GET    | Bar hover tooltip HTML partial           |
| `/charts/{id}/hover?arc=id`    | GET    | Pie hover tooltip + active-arc highlight |
| `/charts/{id}/slice?axis=x&x=v`| GET    | Line slice tooltip HTML partial          |
| `/charts/{id}/click?bar=k`     | POST   | Bar activation toggle (re-render SVG)    |
| `/charts/{id}/toggle?series=id`| POST  | Series hide/show (re-render SVG)         |

State (hidden series, active arc) lives in the in-memory `Registry`; the
documented trade-off (see `docs/PLAN.md` §6) is that this is fine for demos /
small apps and would be moved to a session/cookie store for horizontal
scaling.

## Structure

```
examples/app/
  main.go                  – http.ServeMux wiring (one route per chart family)
  demos/                   – one *.go per family with demo props + data
                             (bar, line, pie, heatmap, waffle, calendar, radar,
                             radialbar, scatterplot, stream, bullet, funnel,
                             boxplot, palettes, themes)
  handlers/chart.go        – page handlers + htmx wiring
  templates/               – layout + chart-card templ components
    layout.templ           – loads htmx + the charts/interact client script
    chart.templ
```

The CSS is hand-written and minimal (inline in the layout). HTMX is loaded
from the CDN (`htmx.org@1.9.12`) and the `charts/interact` hover script is
inlined via `interact.ScriptTag()`; there are no other frontend deps.