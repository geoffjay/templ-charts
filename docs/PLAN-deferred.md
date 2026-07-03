# templ-charts — Deferred Backlog

A single living backlog of everything **not** yet built, consolidated from the
`§11 "Explicitly deferred"` sections of [`PLAN.md`](PLAN.md),
[`PLAN-v2.md`](PLAN-v2.md), [`PLAN-v3.md`](PLAN-v3.md), and the v4 audit (which
verified each item is genuinely absent in the code, with citations). It is
captured here so nothing discovered is lost and a future release can be scoped
from one place.

Items are grouped by theme and each carries a rough size, the code evidence that
it's still absent, and its likely future home.

## 0. Status — the full roadmap

Every item below is now scoped into a themed release. The sections stay
documented (with their code evidence) and are struck as they land.

**Done — v5 ("fidelity & finish", see [`PLAN-v5.md`](PLAN-v5.md)), which
finished the SVG story:**

- ~~**§3 Animation parity** → v5 §4 (wired `MotionProps` across the ~25 v2/v3 charts)~~ — **done**
- ~~**§4 Interactivity model gaps** → v5 §6 (unified hover-others, hierarchy zoom,
  retire line mousemove fallback, opt-in `ResizeObserver`)~~ — **done**
- ~~**§5 Feature completions** → v5 §5 (waffle `areas`, calendar month outline,
  sankey link gradients)~~ — **done**
- ~~**§6 d3-geo `clipCircle`/`clipExtent`** → v5 §3 (the one *correctness* fix)~~ — **done**
- ~~**§9 Test depth** → v5 §7 (network's second golden + the new variant goldens)~~ — **done**

**Scoped — v6 ("scale", see [`PLAN-v6.md`](PLAN-v6.md)), the Canvas + large-N
theme:**

- **§1 Canvas rendering path** → v6 §4 (client draw-list backend)
- **§2 Large-N performance** → v6 §3 (`d3-quadtree` port + Barnes–Hut, Delaunator sweep-hull)

**Scoped — v7 ("completeness & ergonomics", see [`PLAN-v7.md`](PLAN-v7.md)), the
final opportunistic release:**

- **§7 Color spaces** → v7 §3 (HSL/Lab/Lch)
- **§8 Consumption surface** → v7 §4 (public `samples`, `charts/static` → 28, unified color API)
- **§6 geo remainder** (projection catalog, `GeoPath` bounds/centroid, `fitExtent`, TopoJSON) → v7 §5

After v7 the consolidated backlog is closed; remaining work is maintenance, not a
themed release. Nothing here is committed beyond the sections above.

## 1. Canvas rendering path — **scoped into v6 (§4)**

**The single largest remaining nivo capability.** nivo ships `*Canvas` variants
(bar, line, scatterplot, heatmap, network, swarmplot, geo, voronoi, …) for
large-N datasets; templ-charts is SVG-only. v6 §4 builds the **client draw-list**
backend (a deterministic server-emitted list of draw ops replayed into a
`<canvas>` by a dependency-free JS renderer, reusing every `Use{Chart}` hook).

- **Options**: a client `<canvas>` draw-list emitted from the same layout math
  the SVG path uses (v6's choice), or server-side raster (documented secondary).
- **Evidence absent**: `grep -ri canvas` finds only comments/aliases describing
  its absence (`charts/{bar,line,pie}/types.go:6` "no Canvas";
  `charts/axes/compute.go:76` `CanvasAxisProps = AxisProps` alias;
  `charts/theming/bridge.go` canvas style-mapping enum). No `<canvas>`,
  `getContext`, or `toDataURL` anywhere in code.
- **Size**: very large; the natural theme for its own release.
- **Depends on / pairs with**: the large-N performance work (§2) — Canvas only
  pays off once the layout math scales.

## 2. Large-N performance — **scoped into v6 (§3)**

The layout ports are correct and deterministic but use quadratic algorithms
(fine at chart scale, poor for large graphs/point clouds). v4 added the first
benchmarks to baseline these; v6 §3 does the optimizations (a new
`internal/d3/quadtree` port is the shared prerequisite).

- **Force: Barnes–Hut quadtree.** `internal/d3/force/manybody.go` is an exact
  all-pairs O(n²) charge sum (deliberate, "equivalent to θ=0"; the loop at `:65`);
  `collide.go` is O(n²) pairwise (`:44`, its header notes "d3 uses a quadtree
  purely to prune"). A quadtree brings both to O(n log n) with the θ parameter.
- **Delaunay: delaunator / sweep-hull.** `internal/d3/delaunay/delaunay.go` is
  classic Bowyer–Watson, O(n²) worst case (`:7-8`). A Delaunator sweep-hull port
  is O(n log n), behind the same public surface.
- **Render cost**: ~900 `WriteString` calls build SVG by string append; no
  streaming/large-N story. v6's Canvas draw-list emitter is written
  allocation-consciously; the SVG path itself is not rewritten.

## 3. Animation parity (v2/v3 charts) — **done in v5 (§4)**

**Done (v5 Phase 2):** a wired `Animate`/`MotionStagger` was threaded (default
off, byte-stable) through all ~25 v2/v3 charts via shared SMIL enter primitives
in `charts/core/animate.go` — opacity fade-in for rects/arcs/lines/cells,
radius-scale for circles. Six `*-animated` goldens lock the markup; `make ci`
green. See `docs/NOTES.md` (v5 Phase 2).

- **Remaining refinement (still deferred):** per-chart collapsed→final
  *geometry morph* (radius / `d` interpolation) beyond the v1 bar/line/pie
  primitives — the enter animation is currently fade-in / radius-scale. The
  shared `arcs.ArcShape` already supports opt-in d-morph via `AnimateFromPath`;
  extending morph geometry to the other families is a future pass.

## 4. Interactivity model gaps — **done in v5 (§6)**

All four sub-items shipped in v5 Phase 4 (default-off / opt-in where they add
markup, so existing goldens stayed byte-stable). See `docs/NOTES.md` (v5 Phase 4).

- ~~**Unified hover-others dimming.**~~ — **done in v5 (Phase 4 §6)**: chord's
  `:has()` hover-highlight pattern was factored into one shared
  `charts/interact` helper (`HoverHighlight`/`HoverGroup`/`HoverRule`); chord and
  sankey were migrated onto it **byte-identically**, and **network** was given
  the same treatment (dim others, re-light the hovered node + its links +
  neighbours, or a hovered link + its two endpoints), gated on `Interactive &&
  !UseMesh`. New `network-highlight` golden.
- ~~**Interactive zoom / drill-down** for icicle, treemap, circle-packing,
  sunburst.~~ — **done in v5 (Phase 4 §6)**: opt-in `EnableZooming` + `ChartID`
  turn each node into an htmx zoom target (`GET /charts/{id}/zoom?node=`); the
  `htmx` registry gained the four kinds + a `zoom` verb / `FocusID` state, and
  each `Use{Chart}` hook does the per-family focus recompute (partition rescale
  for icicle/sunburst, subtree re-layout for treemap, zoomable-pack transform
  for circle-packing) plus a clickable breadcrumb. Default-off → existing
  goldens byte-stable. See `docs/NOTES.md` (v5 Phase 4 — hierarchy zoom).
- ~~**Line per-mousemove server round-trip.**~~ — **done in v5 (Phase 4 §6)**:
  the client mesh/slice hover path is now the **default**; the per-mousemove
  htmx round-trip is only emitted when the new `ServerHover` opt-in is set (the
  htmx registry sets it automatically, since it *is* the server path).
  `ClientHover` retained as a compat no-op.
- ~~**`ResizeObserver` re-fetch.**~~ — **done in v5 (Phase 4 §6)**: an opt-in
  observer in `charts/interact/script.go`, gated by `data-tc-observe="<url>"`
  (+ optional `data-tc-observe-target`), debounce-re-fetches with the new
  `?w=&h=` on container resize (via htmx when present, else `fetch`); costs
  nothing unless opted in. Cosmetic fluid scaling stays covered by the
  `Responsive` viewBox.

## 5. Feature completions (partial charts) — **done in v5 (§5)**

All three shipped in v5 Phase 3 as opt-in / default-off features (existing
goldens byte-stable). See `docs/NOTES.md` (v5 Phase 3).

- ~~**waffle "areas" polygon layer**~~ — **done**: `WaffleLayerAreas` (opt-in via
  `Layers`) draws per-datum union outlines via `grid.GetCellsPolygons`.
- ~~**calendar month outline-path border**~~ — **done**: re-added
  `MonthBorderColor/Width` (default width 0 → off) backed by a real per-month
  outline path.
- ~~**sankey link gradients**~~ — **done**: re-added `EnableLinkGradient`
  (default false) backed by per-link `<linearGradient>` defs.

## 6. d3-geo completeness — **`clipCircle`/`clipExtent` done in v5 (§3)**

- ~~**`clipCircle`** — the small-circle preclip the azimuthal family needs to
  hide the far hemisphere.~~ **Done (v5 Phase 1):** ported into `internal/d3/geo`
  (`circle.go`), wired into the generic clip framework, and installed as the
  default preclip on all five azimuthal projection constructors — so
  orthographic/gnomonic/stereographic/azimuthal* now render only the visible
  hemisphere. This was the sole item that made output *wrong*; it is now closed.
- ~~**`clipExtent`** — rectangular clip.~~ **Done (v5 Phase 1):** ported the
  screen-space rectangle postclip (`clip_rectangle.go`) + `Projection.ClipExtent`.
- **Full projection catalog** beyond the ~10 nivo exposes; `GeoPath`
  bounds/area/centroid and `fitExtent`/`fitSize`. → **v7 §5.**
- **TopoJSON decoding** — out of scope by design (callers supply GeoJSON);
  → **v7 §5** adds it behind an opt-in decoder.

## 7. Color spaces — **scoped into v7 (§3)**

`internal/d3/color` is **RGB-only** (`color.go:38` `type RGB`; the interface
comment at `:15` says "and, in future, HSL"). HSL / Lab / Lch spaces are unported
— v7 §3 adds them behind the `Color` interface, with opt-in perceptual
interpolation for scales/palettes and Lab-based lightness modifiers (RGB stays
the default, so no golden moves).

## 8. Consumption surface (beyond v4's render helpers) — **scoped into v7 (§4)**

v4 ships `charts/render` helpers only. v7 §4 finishes the consumption story
(all three additive):

- **Public sample-data export.** Real per-chart data generators are unexported in
  `examples/app/demos` (`barData()`, `chordSample()`, …). A downstream user can't
  get ready-made data to try a chart. → **v7 §4** promotes them into a public
  typed `samples` package (reused by the demo as the single source of truth).
- **`charts/static` registry extension.** The registry + `Samples` cover only
  bar/line/pie (`charts/static/types.go`, `samples.go`) — 3 of 28. → **v7 §4**
  extends the dispatcher + samples to all 28 charts.
- **Unified color-setting API.** Five distinct color-config shapes exist across
  families; v4 documents them. → **v7 §4** adds one additive unifying entry
  point (no breaking rename).

## 9. Test depth — **done in v5 (§7)**

Every chart now has ≥2 goldens: v4 lifted 14 of the 15 single-golden charts to
two snapshots, and v5 added `network`'s second golden (`network-highlight`) —
the last one — plus the animated (6), on-variant completion, and azimuthal-geo
goldens from its other workstreams, and the `internal/d3/geo` `clipCircle` /
`clipExtent` port tests. `make golden` is idempotent.

## 10. Priority sketch (non-binding)

A rough ordering if these were to be scoped into future releases:

1. ~~**Correctness now-ish**: geo `clipCircle` (§6).~~ → **v5 §3.**
2. ~~**High-value features**: animation parity (§3), interactive zoom + unified
   hover-others (§4), the partial-chart completions (§5).~~ → **v5 §4–6.**
3. **The big theme (next)**: Canvas (§1) + large-N performance (§2), together —
   the **v6** theme ([`PLAN-v6.md`](PLAN-v6.md)). v5's animation/layout output is
   reused by the Canvas backend, and the quadtree/Delaunator perf ports land
   first so Canvas has scaled layout to render.
4. **Final completeness**: color spaces (§7), the consumption surface (§8 —
   sample-data export, `charts/static` → 28, unified color API), and the geo
   remainder (§6) — the **v7** theme ([`PLAN-v7.md`](PLAN-v7.md)), the last
   planned release. (Test depth §9 folded into v5 §7.)
