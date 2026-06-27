# templ-charts demo app

A runnable demo of the `templ-charts` library: a stdlib `net/http` server
serving bar, line, and pie chart demos with HTMX-backed interactivity.

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

| Route     | Contents                                                     |
|-----------|--------------------------------------------------------------|
| `/`       | Index listing all demos                                      |
| `/bar`    | Stacked, grouped, markers+annotations, legend toggle, totals |
| `/line`   | Single, multi+legend, area+points, slices, mesh             |
| `/pie`    | Plain, donut, half, sorted, active-arc hover, legend toggle  |
| `/themes` | bar / line / pie under default, dark, and custom themes      |

## Interactivity

Interactive demos register a chart instance with `charts/htmx.Registry` and
mount the `htmx.Handler` at `/charts/`. The chart components emit `hx-*`
attributes pointing at the handler's endpoints:

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
  main.go                  – http.ServeMux wiring
  demos/                   – bar/line/pie/themes demo props + data
    bar.go
    line.go
    pie.go
    themes.go
  handlers/chart.go        – page handlers + htmx wiring
  templates/               – layout + chart-card templ components
    layout.templ
    chart.templ
```

The CSS is hand-written and minimal (~80 lines, inline in the layout). HTMX
is loaded from the CDN (`htmx.org@1.9.12`); there are no other frontend deps.