# AGENTS.md — templ-charts

A Go library that wraps [nivo](https://github.com/plouc/nivo)'s chart concepts as [templ](https://github.com/a-h/templ) components generating server-side SVG. v1 ships bar, line, and pie charts plus an HTMX-backed interactivity layer and a runnable demo app. See `docs/PLAN.md` for the full implementation plan.

## Build & test commands

| Task | Command |
|---|---|
| Generate templ sources | `make templ` (or `go generate ./...` after directives are in place) |
| Build all packages | `make build` |
| Run unit tests | `make test` |
| Run go vet | `make vet` |
| Check gofmt | `make fmt` |
| Lint (vet + fmt) | `make lint` |
| Run the demo app | `make run-demo` |
| Tidy modules | `make tidy` |

**Always run `make lint` and `make test` after non-trivial Go/templ changes.**

## Layout

- `charts/` — library packages (mirrors nivo package names; see `docs/PLAN.md` §3)
- `internal/d3/` — vendored pure-Go ports of d3-shape, d3-scale, d3-array, d3-format, d3-time-format, d3-color
- `examples/app/` — runnable demo app (stdlib `net/http`, run via `make run-demo` → http://localhost:8080)
- `contrib/nivo/` — upstream nivo clone (gitignored, reference only; do NOT modify)

## Conventions

- Go 1.26.4 (matches `go.mod`).
- templ components live in `.templ` files; generated `templ_*.go` files are committed alongside sources.
- No Canvas rendering in v1 — SVG only.
- Interactivity via [htmx.org](https://htmx.org) (loaded via CDN in the demo); see `charts/htmx`.
- Animation via SMIL `<animate>` + CSS keyframes, gated by an `Animate bool` prop.
- Tests are standard `go test`; golden SVG snapshots regenerate via `go test -update`.

## Dependencies

- `github.com/a-h/templ` — templ compiler/runtime
- htmx.org — CDN `<script>` in demo HTML (no Go dep)

## Reference

- Upstream nivo: `contrib/nivo/packages/*` (read-only; the design source of truth for types, defaults, compute logic)
- Full plan: `docs/PLAN.md`