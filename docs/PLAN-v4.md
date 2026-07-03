# templ-charts — Implementation Plan (v4)

## 1. Vision

v1–v3 built the framework and reached **full nivo SVG chart parity**: 28 chart
families across three releases, every SVG chart type nivo ships, on a deep
reusable foundation. The library is *complete* but not yet *comfortable to
consume* — and it doesn't yet **show itself off**.

**v4 is the "from complete to consumable" release.** It adds no new chart types
and no new d3 ports. Instead it closes the gaps that a real downstream user hits
on day one: the missing render convenience, the documentation and runnable
examples, the handful of dishonest/inconsistent props that erode trust in the
API, the absence of any performance evidence, and the lack of a proper showcase.
The centrepiece of the showcase is a **per-chart detail page**: a full-width
single-chart view with live theme and palette switchers and the exact Go code a
consumer would write to reproduce it.

Everything found during the v4 audit that is *larger than polish* — the Canvas
rendering path, full animation parity, the remaining d3-geo clipping, large-N
performance work, interactive zoom, TopoJSON, and the rest — is deliberately
**not** in v4. It is captured in [`docs/PLAN-deferred.md`](PLAN-deferred.md) as a
single living backlog so nothing discovered here is lost.

v4 stays **SVG-only** and follows the same conventions as v1/v2/v3. See
[`docs/PLAN.md`](PLAN.md), [`docs/PLAN-v2.md`](PLAN-v2.md), and
[`docs/PLAN-v3.md`](PLAN-v3.md) for the prior structure this document mirrors.

## 2. Confirmed decisions

| Concern | Decision |
|---|---|
| Theme | **Consumability & showcase.** No new charts, no new d3 ports. |
| Convenience API | **Render helpers only** — a small `charts/render` package (`String`/`To`). **No** public sample-data export and **no** extension of the `charts/static` registry beyond its v1 bar/line/pie scope (both deferred — see [`PLAN-deferred.md`](PLAN-deferred.md)). |
| API honesty | Fix props that advertise features the code doesn't deliver: remove/gate the dead `MotionProps{Animate}` defaults on non-animated charts, `sankey.EnableLinkGradient`, `calendar.MonthBorder*`. Full implementations are deferred; v4 only makes the surface honest. |
| API consistency | Standardize the interactivity flag on `Interactive` (rename the four v1 charts' `IsInteractive`; keep `line.ClientHover`), add the missing `pie.IsFocusable`, and document the five color-setting shapes. |
| Documentation | Runnable Go examples (`ExampleXxx`) per chart family, a consumer-facing usage/API doc, a package doc comment for `charts/colors`, and README field-name corrections. |
| Performance evidence | First benchmarks in the repo (d3 ports + render paths) plus a **load/showcase demo page** illustrating render cost and scaling under N. |
| Showcase | A **per-chart detail page** (full-width chart + theme switcher + palette switcher + copy-pasteable Go snippet); index cards link to it. |
| Deferred | Consolidate all v1/v2/v3 §11 deferrals + everything else the audit surfaced into [`docs/PLAN-deferred.md`](PLAN-deferred.md). |

## 3. Convenience render helpers — `charts/render`

Today every consumer hand-rolls the templ boilerplate for every chart:

```go
var b strings.Builder
if err := bar.Bar(props).Render(context.Background(), &b); err != nil { ... }
html := b.String()
```

The library itself repeats this exact triplet in an **unexported** `renderComponent`
in every chart package (e.g. `charts/bar/render.go:16`), so consumers can't reuse
it. v4 exports one:

```go
package render // charts/render

// String renders any chart component (or any templ.Component) to an SVG/HTML
// string using a background context.
func String(c templ.Component) (string, error)

// To renders a component into w. Thin wrapper over c.Render for symmetry.
func To(w io.Writer, c templ.Component) error

// StringCtx / ToCtx take an explicit context.Context for cancellation.
func StringCtx(ctx context.Context, c templ.Component) (string, error)
func ToCtx(ctx context.Context, w io.Writer, c templ.Component) error
```

Consumer usage becomes:

```go
svg, err := render.String(bar.Bar(props))
```

- Deliberately **generic over `templ.Component`** — no chart-type registry, no
  `any`-typed props. This keeps the surface tiny and matches the "render helpers
  only" scope.
- These four functions are what the detail-page **code snippets** show, so the
  documented snippet is literally copy-pasteable.
- Package location `charts/render` (importable as `render`) keeps the repo root
  clean and discoverable next to the other `charts/*` packages.

## 4. API honesty & consistency

The audit found props that either advertise a feature the code never delivers or
name the same concept two different ways. These are consumption papercuts, so
they land in v4 (the *implementations* they imply are deferred).

### 4.1 Honesty — remove/gate dead props
- **Animation is v1-only.** `MotionProps{Animate:true}` is defaulted on ≥8 v2/v3
  charts (heatmap, waffle, sankey, treemap, radar, stream, funnel, boxplot) but
  never read or rendered — only bar/line/pie emit SMIL `<animate>`. v4 removes the
  misleading defaults from the charts that don't animate (or documents them as
  no-ops). Full v2/v3 animation parity → [`PLAN-deferred.md`](PLAN-deferred.md).
- **`sankey.EnableLinkGradient`** (`charts/sankey/types.go:160`) — declared,
  defaulted, never read. Gate/remove; gradient rendering → deferred.
- **`calendar.MonthBorderColor/MonthBorderWidth`** (`charts/calendar/types.go:94`)
  — declared, defaulted, never emitted. Gate/remove; the month outline path →
  deferred.

### 4.2 Consistency
- **Interactivity flag.** Standardize on **`Interactive`** (already used by 24 of
  28 families); rename `IsInteractive` on bar/line/pie/heatmap. Keep line's
  `ClientHover`. This is a breaking change for the four v1 charts, acceptable
  pre-1.0, and removes the README's documented-but-wrong field name.
- **`pie.IsFocusable`.** Add the field and thread it into `SvgWrapper`
  (`charts/pie/pie.templ:15`) — pie is the only entrypoint missing it.
- **Color-setting shapes.** Five distinct shapes exist (ordinal config; heatmap
  continuous; calendar `[]string`+`EmptyColor`; choropleth scheme-id+`Steps`;
  network plain strings). v4 does **not** unify them (out of scope) but
  **documents** each shape and when it applies (§5).

## 5. Documentation

- **Runnable Go examples.** Add `example_test.go` with an `ExampleXxx` per chart
  family (using `render.String`), so `go doc`/pkg.go.dev renders a working
  snippet for every chart. Zero exist today.
- **Consumer usage doc** (`docs/USAGE.md` or an expanded README section): the
  render-helper pattern, the `Interactive`/`ClientHover`/`Responsive`/a11y prop
  set, the five color-setting shapes and when each applies, and the
  static-embedding (zero-JS) vs interactive story.
- **`charts/colors` package doc comment** — the only package missing one.
- **README corrections** — fix the `Interactive` field-name reference and add the
  `render.String` quickstart.

## 6. Benchmarks & load/showcase demo

The repo has **zero benchmarks**. v4 adds the first, focused on the paths most
likely to regress or to scale non-linearly:

- **d3-port benchmarks:** `internal/d3/force` (all-pairs O(n²) charge/collide),
  `internal/d3/delaunay` (O(n²) Bowyer–Watson), and the hierarchy/sankey layouts,
  parameterized by N.
- **Render-path benchmarks:** the `~585 WriteString` SVG build for a few
  representative charts (bar, line, heatmap, network) at increasing N.
- **Load/showcase demo page** (`/benchmark` in the demo app): renders a chart at
  several dataset sizes and reports server-side render time, illustrating how the
  library behaves under load. This is documentation-by-demonstration, and it also
  gives the deferred large-N performance work (Barnes–Hut, delaunator — see
  [`PLAN-deferred.md`](PLAN-deferred.md)) a baseline to beat.

Benchmarks run via `go test -bench`; a `make bench` target is added. They are
**not** part of `make ci` (timing is environment-sensitive).

## 7. Per-chart detail page — the showcase

Today the demo index (`/`) is a text link-list to per-family pages that show
several static variations. v4 adds a **detail page per chart**: a full-width
single-chart view with live controls and the source to reproduce it.

### 7.1 What the page shows
- One representative chart rendered **full page width** (its own layout, not the
  card grid).
- A **theme switcher** — default / dark / custom, reusing the theme construction
  from `examples/app/demos/themes.go` (`theming.ExtendDefaultTheme`).
- A **palette switcher** — driven by `colors.Palettes()` and applied via
  `Palette.Ordinal()` / the `Colors` field (the mechanism from
  `examples/app/demos/palettes.go`).
- A **code panel** — the exact Go a consumer writes, ending in
  `render.String(...)` (§3). Hand-written per representative chart (~30 curated
  snippets); embedding demo source is brittle and rejected.
- Full interactivity for free: htmx and `interact.ScriptTag()` are already loaded
  globally in `layout.templ`.

### 7.2 The linchpin — a family registry (demo-app internal)
The 24 non-interactive families each declare their own `XDemo` struct with a
typed `Props` and **no shared `Kind`**, so there is no way to enumerate "one
chart per family" generically today. v4 adds a **demo-app-internal** registry
(not public API — consistent with the "render helpers only" decision):

```go
type ChartEntry struct {
    Slug, Title, Description string
    Snippet                  string // hand-written Go source
    Render func(theme *theming.Theme, palette colors.PaletteID) (string, error)
}
```

Each family contributes one entry whose `Render` closure owns its concrete props
and injects the chosen `Theme`/`Colors` (the per-family type assertion the
generic path can't do). This registry is the single new orchestration layer; it
also feeds the index card list.

### 7.3 Switching mechanics
Switchers use **HTMX fragment swaps** (htmx is already global): a new endpoint
`/chart/{slug}?theme=…&palette=…` returns just the re-rendered chart fragment
via the registry's `Render` closure; `hx-get`/`hx-target` on the switch controls.
A plain query-param full-page render is the no-JS fallback. No new dependency.

### 7.4 New presentation pieces
- A full-width detail `templ` layout (distinct from `.grid`).
- Theme-switcher and palette-switcher control components.
- A `<pre><code>` snippet panel + minimal CSS for it (no syntax highlighter).
- Index restyled from the text list to cards whose links point at `/chart/{slug}`.

## 8. Build / test

- **Unit/golden**: unchanged discipline. Existing goldens stay byte-stable; the
  §4 prop renames/removals regenerate the affected snapshots (via `make golden`).
- **Examples**: `ExampleXxx` functions double as compile-checked doc and run
  under `go test`.
- **Benchmarks**: `make bench` (not in CI).
- **Demo handler tests**: extend `examples/app/handlers` coverage to the new
  `/chart/{slug}` route and the theme/palette query params.
- CI unchanged: `make ci` (`make lint` + `make test`) stays green.

## 9. Implementation order (topological)

1. [x] **`charts/render` helpers** (§3) + README quickstart + `charts/colors`
   package doc. Small, foundational — everything else can use `render.String`.
2. [x] **API honesty & consistency** (§4): standardized `Interactive` (renamed
   `IsInteractive` on bar/line/pie/heatmap), added `pie.IsFocusable`, removed the
   dead `core.MotionProps` embed from the 25 non-animating charts and the dead
   `sankey.EnableLinkGradient` / `calendar.MonthBorder*` props. No golden changes
   (output byte-stable); README field name already correct after the rename.
3. **Documentation** (§5): `ExampleXxx` per family, `docs/USAGE.md`, color-shape
   docs.
4. **Benchmarks + load demo** (§6): bench suite, `make bench`, `/benchmark` page.
5. **Per-chart detail page** (§7): the family registry, the `/chart/{slug}`
   endpoint + switchers, the full-width layout + snippet panel, index → cards.
6. **Polish**: deepen the shallowest goldens where cheap, docs/README pass for the
   full v4 surface, `make ci` green.

## 10. Scope summary for v4

**New API**: `charts/render` (`String`/`To`/`StringCtx`/`ToCtx`).

**Fixes**: `Interactive` naming standardized; `pie.IsFocusable` added; dead props
(`MotionProps{Animate}` on non-animated charts, `sankey.EnableLinkGradient`,
`calendar.MonthBorder*`) removed/gated for an honest surface.

**Docs**: runnable examples per chart, `docs/USAGE.md`, `charts/colors` doc
comment, README corrections.

**Perf evidence**: first benchmarks (d3 ports + render paths), `make bench`, and a
`/benchmark` load-demo page.

**Showcase**: per-chart detail pages with theme + palette switchers and
copy-pasteable Go snippets, backed by a demo-app family registry.

**Explicitly out of scope** (→ [`PLAN-deferred.md`](PLAN-deferred.md)): Canvas
rendering, full animation parity, large-N perf (Barnes–Hut/delaunator), geo
`clipCircle`/`clipExtent` + full projection catalog + TopoJSON, interactive zoom,
unified hover-others dimming, waffle areas layer, calendar month outline path,
sankey link gradients, HSL/Lab/Lch color, `ResizeObserver` re-fetch, public
sample-data export, and `charts/static` registry extension.

## 11. Explicitly deferred

All larger work — everything beyond consumability polish and the showcase — lives
in [`docs/PLAN-deferred.md`](PLAN-deferred.md), the consolidated backlog that
folds together the §11 deferrals from v1/v2/v3 with everything the v4 audit
surfaced. v4 intentionally touches none of it.
