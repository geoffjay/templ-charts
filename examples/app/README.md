# templ-charts demo app

A runnable demo of the `templ-charts` library: a stdlib `net/http` server
serving demos for all twenty-eight chart families with a hybrid interactivity
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

Then open <http://localhost:8000>.

## Pages

| Route            | Contents                                                     |
|------------------|--------------------------------------------------------------|
| `/`              | Index listing all charts; each card links to its detail page |
| `/chart/{slug}`  | Per-chart detail page: full-width chart + live theme & palette switchers + the Go snippet to reproduce it (e.g. `/chart/bar`, `/chart/sankey`) |
| `/benchmark`     | Live server-side render time + SVG size for a bar chart across dataset sizes (load/scaling showcase) |
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
| `/bump`        | Ranking over time: bump curves, end labels, point hover      |
| `/marimekko`   | Variable-width stacked bars (width by value, `d3.Stack`)     |
| `/parallel-coordinates` | One axis per variable; a polyline per record        |
| `/polar-bar`   | Stacked bars wrapped into a full circle                      |
| `/treemap`     | Nested rectangles (d3-hierarchy), squarify/binary tiling     |
| `/sunburst`    | Radial partition (d3-hierarchy), colors inherited down tree  |
| `/icicle`      | Depth-banded partition rectangles (d3-hierarchy)             |
| `/circle-packing` | Welzl enclosing-circle packing (d3-hierarchy)             |
| `/tree`        | Tidy-tree / dendrogram node-link diagrams, bump links        |
| `/voronoi`     | Delaunay triangulation + Voronoi cells (d3-delaunay)         |
| `/network`     | Force-directed graph (d3-force), deterministic fixed ticks   |
| `/swarmplot`   | Grouped distribution relaxed with d3-force, mesh hover       |
| `/sankey`      | Flow diagram (d3-sankey), monotone-curve ribbons             |
| `/chord`       | Radial flow diagram (d3-chord), entity arcs + ribbons        |
| `/geo`         | GeoMap + Choropleth (d3-geo), projections, graticule, legend |
| `/palettes`    | The full color-palette catalog applied to bars               |
| `/themes`      | bar / line / pie under default, dark, and custom themes      |

## Interactivity

The interactivity is **hybrid**:

- **Client-side** (`charts/interact`, loaded once in the layout): ephemeral
  hover — tooltips, the line crosshair, and nearest-point hit-testing — runs
  entirely in the browser off `data-tc-*` attributes, with no server
  round-trip. Most charts and `line`'s mesh/slice hover use this path; the
  `internal/d3/delaunay` port backs true voronoi-mesh hover on
  line/scatterplot/bump/swarmplot/tree.
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
documented trade-off is that this is fine for demos / small apps and would be
moved to a session/cookie store for horizontal scaling.

## Structure

```
examples/app/
  main.go                  – http.ServeMux wiring (one route per chart family)
  demos/                   – one *.go per family with demo props + data
                             (bar, line, pie, heatmap, waffle, calendar, radar,
                             radialbar, scatterplot, stream, bullet, funnel,
                             boxplot, bump, marimekko, parallelcoordinates,
                             polarbar, treemap, sunburst, icicle, circlepacking,
                             tree, voronoi, network, swarmplot, sankey, chord,
                             geo, palettes, themes; geo bundles a sample
                             world-countries GeoJSON)
  handlers/chart.go        – page handlers + htmx wiring
  templates/               – layout + chart-card templ components
    layout.templ           – loads htmx + the charts/interact client script
    chart.templ
```

The CSS is hand-written and minimal (inline in the layout). HTMX is loaded
from the CDN (`htmx.org@1.9.12`) and the `charts/interact` hover script is
inlined via `interact.ScriptTag()`; there are no other frontend deps.
