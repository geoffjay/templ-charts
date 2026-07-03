# templ-charts — Deferred Backlog

A single living backlog of everything **not** yet built, consolidated from the
`§11 "Explicitly deferred"` sections of [`PLAN.md`](PLAN.md),
[`PLAN-v2.md`](PLAN-v2.md), [`PLAN-v3.md`](PLAN-v3.md), and the v4 audit (which
verified each item is genuinely absent in the code, with citations). It is
captured here so nothing discovered is lost and a future release can be scoped
from one place.

Items are grouped by theme and each carries a rough size, the code evidence that
it's still absent, and its likely future home.

## 0. Status — what v5 now scopes

**v5 ("fidelity & finish" — see [`PLAN-v5.md`](PLAN-v5.md)) claims the following
sections**, finishing the SVG story. They stay documented below (with their code
evidence, still valid) but are **no longer part of the standing backlog** — track
them in `PLAN-v5.md` and strike them here as they land:

- **§3 Animation parity** → v5 §4 (wired `MotionProps` across the ~25 v2/v3 charts)
- **§4 Interactivity model gaps** → v5 §6 (unified hover-others, hierarchy zoom,
  retire line mousemove fallback, opt-in `ResizeObserver`)
- **§5 Feature completions** → v5 §5 (waffle `areas`, calendar month outline,
  sankey link gradients)
- **§6 d3-geo `clipCircle`/`clipExtent`** → v5 §3 (the one *correctness* fix)
- **§9 Test depth** → v5 §7 (network's second golden + the new variant goldens)

**The standing remainder after v5** — the active backlog a future release should
scope from — is:

- **§1 Canvas + §2 large-N performance** — paired, the v6 theme.
- **§7 Color spaces** — opportunistic.
- **§8 Consumption surface** (beyond v4's render helpers) — opportunistic.

Nothing below is committed to a version beyond the v5-scoped sections above.

## 1. Canvas rendering path — the headline deferral

**The single largest remaining nivo capability.** nivo ships `*Canvas` variants
(bar, line, scatterplot, heatmap, network, swarmplot, geo, voronoi, …) for
large-N datasets; templ-charts is SVG-only.

- **Options**: server-side raster, or a client `<canvas>` draw-list emitted from
  the same layout math the SVG path uses.
- **Evidence absent**: `grep -ri canvas` finds only comments/aliases describing
  its absence (`charts/{bar,line,pie}/types.go:6` "no Canvas";
  `charts/axes/compute.go:76` `CanvasAxisProps = AxisProps` alias;
  `charts/theming/bridge.go` canvas style-mapping enum). No `<canvas>`,
  `getContext`, or `toDataURL` anywhere in code.
- **Size**: very large; the natural theme for its own release.
- **Depends on / pairs with**: the large-N performance work (§2) — Canvas only
  pays off once the layout math scales.

## 2. Large-N performance

The layout ports are correct and deterministic but use quadratic algorithms
(fine at chart scale, poor for large graphs/point clouds). v4 adds the first
benchmarks to baseline these; the optimizations themselves are deferred.

- **Force: Barnes–Hut quadtree.** `internal/d3/force/manybody.go:73` is an exact
  all-pairs O(n²) charge sum (deliberate, "equivalent to θ=0"); `collide.go:53`
  is O(n²) pairwise. A quadtree brings charge to O(n log n).
- **Delaunay: delaunator / sweep-hull.** `internal/d3/delaunay/delaunay.go` is
  classic Bowyer–Watson, O(n²) worst case (each insertion scans all triangles,
  `:109`). A sweep-hull port is O(n log n).
- **Render cost**: ~585 `WriteString` calls build SVG by string append; no
  streaming/large-N story. Benchmark first (v4), optimize if needed.

## 3. Animation parity (v2/v3 charts) — **scoped into v5 (§4)**

Animation is **v1-only** but advertised repository-wide. `MotionProps{Animate}`
is defaulted on ≥8 v2/v3 charts (heatmap, waffle, sankey, treemap, radar, stream,
funnel, boxplot) yet **never read or rendered** — only bar/line/pie emit SMIL
`<animate>` (`charts/bar/bar_item.templ:151`, `charts/line/lines.templ:39`,
`charts/arcs/arc_shape.templ:66`).

- **v4** makes the surface honest (removes/gates the dead defaults).
- **Deferred**: actually wiring SMIL `<animate>` (or CSS keyframes) into the
  v2/v3 charts for real enter/update transitions.
- **Size**: medium-large (per-chart geometry interpolation, ~25 charts).

## 4. Interactivity model gaps — **scoped into v5 (§6)**

- **Unified hover-others dimming.** Flow/relationship charts handle "dim the
  others on hover" inconsistently: chord uses client CSS `:has()`
  (`charts/chord/render.go:99`), network has **no** hover-others at all, sankey's
  is deferred (props not even declared). A shared model (client CSS or HTMX
  state) across sankey/network/chord is unbuilt.
- **Interactive zoom / drill-down** for icicle, treemap, circle-packing,
  sunburst. Layout is present; the client-driven zoom transition is not
  (deferral comments in each `types.go`).
- **Line per-mousemove server round-trip.** The non-JS fallback
  (`charts/line/mesh.templ:64`, `mousemove throttle:40ms`) still round-trips to
  the server; the client mesh path (`charts/interact/script.go`) already replaces
  it when `DataMesh` is set. Fully retiring the fallback is deferred.
- **`ResizeObserver` re-fetch.** Pixel-accurate re-render on axes-heavy charts
  when the container resizes (cosmetic scaling is already covered by the
  `Responsive` viewBox prop). No `ResizeObserver` in `charts/interact/script.go`.

## 5. Feature completions (partial charts) — **scoped into v5 (§5)**

Small, well-scoped finishes to existing charts:

- **waffle "areas" polygon layer** — only `cells`/`legends` layers exist
  (`charts/waffle/types.go:6`, `render.go:75`); the polygon `areas` layer nivo
  ships is absent.
- **calendar month outline-path border** — `MonthBorderColor/Width` props exist
  and are defaulted (`charts/calendar/types.go:94`) but never emitted; the month
  outline path is unimplemented (`render.go:52`). (v4 gates the dead props; the
  path itself is here.)
- **sankey link gradients** — `EnableLinkGradient` declared/defaulted but never
  read (`charts/sankey/types.go:160`); gradient `<defs>` rendering deferred.

## 6. d3-geo completeness — **`clipCircle`/`clipExtent` scoped into v5 (§3)**

The `clipCircle`/`clipExtent` items below are v5's one correctness fix; the
remaining geo items (full projection catalog, `GeoPath` bounds/centroid/
`fitExtent`, TopoJSON) stay deferred.

- **`clipCircle`** — the small-circle preclip the azimuthal family needs to hide
  the far hemisphere. Absent (`internal/d3/geo` implements only
  `clipAntimeridian`). **Consequence today**: orthographic/gnomonic/stereographic/
  azimuthal* projections render the **whole sphere** and can show far-side
  geometry — a real correctness gap, not cosmetic.
- **`clipExtent`** — rectangular clip (`internal/d3/geo/projection.go:10`
  "clipExtent omitted").
- **Full projection catalog** beyond the ~10 nivo exposes; `GeoPath`
  bounds/area/centroid and `fitExtent`/`fitSize`.
- **TopoJSON decoding** — out of scope by design (callers supply GeoJSON); revisit
  only if a consumer needs it.

## 7. Color spaces

`internal/d3/color` is **RGB-only** (`color.go:38` `type RGB`; the interface
comment at `:15` says "and, in future, HSL"). HSL / Lab / Lch spaces are unported
— add only when a chart needs perceptually-uniform interpolation or lightness
modifiers beyond what RGB gives.

## 8. Consumption surface (beyond v4's render helpers)

v4 ships `charts/render` helpers only. Deferred by explicit v4 scope decision:

- **Public sample-data export.** Real per-chart data generators are unexported in
  `examples/app/demos` (`barData()`, `chordSample()`, …). A downstream user can't
  get ready-made data to try a chart. An exported `samples` package (typed data
  per chart) is deferred.
- **`charts/static` registry extension.** The registry + `Samples` cover only
  bar/line/pie (`charts/static/types.go:49`, `samples.go:20`) — 3 of 28. Extending
  the dispatcher and samples to all 28 charts is deferred.
- **Unified color-setting API.** Five distinct color-config shapes exist across
  families; v4 documents them, but a single unifying interface is deferred.

## 9. Test depth — **scoped into v5 (§7)**

Every chart has a test and at least one golden. v4 polish lifted 14 of the 15
single-golden charts to two snapshots; **`network` is the last chart still at a
single golden**. v5 adds network's second golden plus the new animated/on-
variant/azimuthal-geo goldens its other workstreams introduce.

## 10. Priority sketch (non-binding)

A rough ordering if these were to be scoped into future releases:

1. ~~**Correctness now-ish**: geo `clipCircle` (§6).~~ → **v5 §3.**
2. ~~**High-value features**: animation parity (§3), interactive zoom + unified
   hover-others (§4), the partial-chart completions (§5).~~ → **v5 §4–6.**
3. **The big theme (next)**: Canvas (§1) + large-N performance (§2), together —
   the natural **v6** theme. v5's animation/layout output is reused by a Canvas
   backend, so it lands first.
4. **Opportunistic**: color spaces (§7) and the consumption surface (§8 —
   sample-data export, `charts/static` registry extension, unified color API),
   as consumers demand them. (Test depth §9 folded into v5 §7.)
