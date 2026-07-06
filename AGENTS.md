# AGENTS.md — templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart concepts as [templ](https://github.com/a-h/templ) components generating server-side SVG. It ships twenty-eight chart families (full nivo SVG parity) plus a hybrid interactivity layer (client-side hover + HTMX state changes), an opt-in Canvas backend for large-N scatterplot/heatmap, and a runnable demo app. See `docs/USAGE.md` for the consumer guide.

## Build & test commands

| Task | Command |
|---|---|
| Generate templ sources | `make templ` (or `go generate ./...` after directives are in place) |
| Build all packages | `make build` |
| Run unit tests | `make test` |
| Run go vet | `make vet` |
| Check gofmt | `make fmt` |
| Lint (vet + fmt) | `make lint` |
| Run tests with coverage | `make cover` |
| Regenerate golden snapshots | `make golden` |
| CI (lint + test) | `make ci` |
| Run the demo app | `make run-demo` |
| Tidy modules | `make tidy` |

**Always run `make lint` and `make test` after non-trivial Go/templ changes.**

### Golden snapshots

Golden SVG snapshots live under each chart package's `testdata/golden/`, arc
path strings under `charts/arcs/testdata/golden/`, and the ported-geometry path
strings under `internal/d3/*/testdata/golden/`. Tests compare rendered output
against these committed files via `internal/golden.Assert`.

When a render change is **intentional**, regenerate and commit the updated
snapshots:

```sh
make golden          # regenerates all golden files
```

The `-update` flag is only honored by packages that import `internal/golden`;
`make golden` scopes the run to those packages so unrelated test binaries
don't reject the flag.

## Layout

- `charts/` — library packages (mirror nivo package names): the 28 chart families plus core, theming, scales, colors, axes, arcs, text, tooltip, legends, annotations, interact, static, grid, polar-axes, htmx, canvas, samples, render
- `internal/d3/` — vendored pure-Go ports of d3-shape, d3-scale, d3-array, d3-format, d3-time-format, d3-color, d3-hierarchy, d3-delaunay, d3-force, d3-sankey, d3-chord, d3-geo, d3-quadtree
- `internal/golden/` — small snapshot-test helper (`Assert` + `-update` flag) used by the golden SVG/path tests
- `examples/app/` — runnable demo app (stdlib `net/http`, run via `make run-demo` → http://localhost:8000)
- `contrib/nivo/` — upstream nivo clone (gitignored, reference only; do NOT modify)

## Conventions

- Go 1.26.4 (matches `go.mod`).
- templ components live in `.templ` files; generated `templ_*.go` files are committed alongside sources. The `make templ` target pins the CLI to the version in `go.mod` (currently v0.3.1020).
- SVG is the default renderer; an opt-in Canvas backend (`charts/canvas`, `Render: theming.EngineCanvas`) is available for large-N scatterplot/heatmap.
- Interactivity via [htmx.org](https://htmx.org) (loaded via CDN in the demo); see `charts/htmx`.
- Animation via SMIL `<animate>` + CSS keyframes, gated by an `Animate bool` prop.
- Tests are standard `go test`; golden SVG snapshots regenerate via `make golden` (or `go test ./charts/{bar,line,pie,arcs} -update`).

## Dependencies

- `github.com/a-h/templ` — templ compiler/runtime
- htmx.org — CDN `<script>` in demo HTML (no Go dep)

## Reference

- Upstream nivo: `contrib/nivo/packages/*` (read-only; the design source of truth for types, defaults, compute logic)
- Consumer guide: `docs/USAGE.md`
