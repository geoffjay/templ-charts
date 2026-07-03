# templ-charts — Implementation Plan (v7)

## 1. Vision

By the end of v6 the nivo gap is, in substance, **closed**: v1–v3 reached full
nivo SVG chart parity (28 families), v4 made the library consumable, v5 was
fidelity & finish (correct geo, animation, zoom, hover), and v6 was scale (the
Canvas rendering path + the large-N performance work — a `d3-quadtree` port
driving Barnes–Hut, a Delaunator sweep-hull, and a client draw-list Canvas
backend).

**v7 is the "completeness & ergonomics" release — the last one.** It is a
deliberate mop-up of the three *opportunistic* items that outlived every themed
release because nothing forced them earlier: perceptual **color spaces**, a
first-class **consumption surface** (public sample data, a full static registry,
one color-setting API), and **geo measurement completeness** (`GeoPath`
bounds/centroid/area + `fitExtent`/`fitSize`, a fuller projection catalog, and
optional TopoJSON). None of these add a chart type or a rendering path; they
sand down the remaining rough edges so the library is complete *and* pleasant to
adopt end-to-end.

After v7 there is **no planned v8** — the remaining work is ongoing maintenance
(track upstream nivo/d3 changes, respond to consumer demand), not a themed
release. v7 follows the same per-package anatomy, verbatim-defaults, and
golden-test discipline as v1–v6.

See [`docs/PLAN.md`](PLAN.md) → [`docs/PLAN-v6.md`](PLAN-v6.md) for the prior
structure this mirrors, and [`docs/PLAN-deferred.md`](PLAN-deferred.md) §7–§8
(plus the geo-beyond-correctness remainder of §6) for the backlog it closes.

## 2. Confirmed decisions

| Concern | Decision |
|---|---|
| Theme | **Completeness & ergonomics.** Close the opportunistic remainder; no new chart types, no new rendering path. |
| Byte-stability | Unchanged discipline. Color-space additions must leave existing **RGB** interpolation byte-identical; the samples/static/geo-measurement work is **purely additive**. New behaviour is asserted by *new* goldens/tests only; `make ci` stays green. |
| Color scope | Add **HSL, Lab, Lch** to `internal/d3/color` behind the existing `Color` interface; wire them only where a chart/scale genuinely benefits from perceptually-uniform interpolation, opt-in, defaulting to today's RGB behaviour. |
| Consumption scope | A public typed **`samples`** package, extend the **`charts/static`** registry from 3 → all 28 charts, and add **one** unifying color-setting API over the five current config shapes (additive, not a breaking rename). |
| Geo scope | `GeoPath` **bounds/area/centroid**, **`fitExtent`/`fitSize`**, a **fuller projection catalog**, and **optional TopoJSON** decoding (behind a small decoder; callers still may supply GeoJSON). |
| Ordering | Independent workstreams; do color first (it unblocks nicer scales), then consumption surface, then geo — but they can land in any order and even in parallel. |

## 3. Color spaces

`internal/d3/color` is **RGB-only** today: `color.go:38` `type RGB`, and the
interface comment at `:15` literally says "the common interface implemented by
RGB (and, in future, HSL)." That "in future" is v7.

- Port **HSL, Lab, and Lch** as additional `Color` implementations, each with
  the standard conversions to/from RGB and each other, and `Interpolate` in its
  own space (matching d3-color's `interpolateHsl`/`interpolateLab`/
  `interpolateHcl`). Follow the existing RGB port's structure and its
  golden/unit tests validated against d3-color.
- **Wire opt-in, default unchanged.** Sequential/diverging color scales and the
  palette `Swatch`/gradient sampling currently interpolate in RGB; add a
  space selector (default `RGB`) so callers can request perceptually-uniform
  Lab/Lch interpolation without moving any existing golden. Lightness/darken
  modifiers (the `InheritedColorConfig.Modifiers` the color config already
  declares) gain a proper Lab-based implementation where they today approximate
  in RGB.
- **Risk**: the Lab/Lch white-point and gamma details are the fiddly part
  (same class as any color-space port); golden-test conversions and
  interpolation midpoints against d3-color values.

## 4. Consumption surface

v4 shipped the `charts/render` helpers and per-chart `ExampleXxx`; v7 finishes
the "easy to adopt" story. Three additive pieces:

- **Public typed sample data (`samples` package).** Real per-chart data
  generators exist but are **unexported** in `examples/app/demos`
  (`barData()`, `chordSample()`, `networkSample()`, … — confirmed across
  `examples/app/demos/*.go`), so a downstream user can't get ready-made data to
  try a chart. Promote them into a public `charts/samples` (or top-level
  `samples`) package returning typed, ready-to-render data per chart family,
  reused by the demo app so there's a single source of truth.
- **`charts/static` registry → all 28.** The registry + `Samples` cover only
  **bar/line/pie** today (`charts/static/types.go`: `ChartTypeBar/Line/Pie`,
  `barComponent`/`lineComponent`/`pieComponent`; `Samples` "map of demo data for
  bar/line/pie"). Extend the `ChartComponent` adapter + `ChartType` enum + the
  reflective override path to every chart family, so a caller can render any
  chart by id + data + overrides through one uniform dispatcher.
- **Unified color-setting API.** Five distinct color-config shapes exist across
  families (v4 documented them; the ordinal `OrdinalColorScaleConfig`, the
  `InheritedColorConfig`, the per-chart `Colors`/`NodeColor`/etc.). Add **one**
  ergonomic entry point (e.g. a `colors.Set(...)` / typed option) that resolves
  to the right underlying shape per chart, **without** renaming or removing the
  existing fields — additive, so no consumer breaks and no golden moves.

## 5. Geo completeness (beyond correctness)

v5 made the projections *correct* (`clipCircle`/`clipExtent`); v7 makes
`GeoPath` *measure* and fits, and broadens the catalog. Today `internal/d3/geo`
ships path generation and ~10 projections (`GeoMercator`, `GeoEquirectangular`,
`GeoTransverseMercator`, `GeoEqualEarth`, `GeoAzimuthalEqualArea`,
`GeoAzimuthalEquidistant`, `GeoGnomonic`, `GeoOrthographic`, `GeoStereographic`,
`naturalEarth1`) but **no** `GeoPath` bounds/area/centroid and **no**
`fitExtent`/`fitSize`.

- **`GeoPath` measurement**: port the bounds, area, and centroid stream sinks
  (d3-geo's `path.bounds`/`path.area`/`path.centroid`), which also unblocks…
- **`fitExtent`/`fitSize`**: auto-scale + translate a projection so a GeoJSON
  object fills a given box — the single most-requested geo ergonomic (it removes
  the manual `scale`/`center` fiddling every geo demo needs today).
- **Fuller projection catalog**: add the commonly-used projections nivo/d3
  expose beyond the current ~10 (e.g. conic conformal/equal-area, azimuthal
  variants, additional pseudocylindricals), each with a path golden vs d3.
- **Optional TopoJSON**: a small decoder so callers may pass TopoJSON (still
  supporting GeoJSON directly); it was "out of scope by design" only because no
  consumer needed it — v7 adds it behind an opt-in decoder, not the default
  input path.
- **Risk**: the fit math and the measurement stream sinks are exacting; reuse
  the v3/v5 geo golden strategy (3-decimal rounding, Go-generated snapshots,
  coordinates unit-checked against known d3 output).

## 6. Build / test

- **Byte-stability**: color defaults stay RGB, so existing scale/palette goldens
  don't move; the samples/static/geo-measurement additions are new API surface
  with new tests. `make ci` (`make lint` + `make test`) stays green.
- **Ports vs d3**: HSL/Lab/Lch conversions + interpolation, `GeoPath`
  bounds/area/centroid, `fitExtent`/`fitSize`, and each new projection are
  golden/unit-tested against d3-color / d3-geo output.
- **Consumption**: the `charts/static` extension gets a render test per newly
  registered chart; the `samples` package gets a "every generator returns
  renderable data" test; the unified color API gets a test that it resolves to
  the same output as the underlying config shapes.
- **`make golden`** extended to any new golden-owning path; `.templ` (if any)
  regenerated **per-file**.

## 7. Demo app — `examples/app`

- Repoint the demo's data to the new public `samples` package (single source of
  truth), proving the export is genuinely usable.
- A **color-space** toggle on a sequential/diverging chart's detail page (RGB vs
  Lab/Lch interpolation) to show the perceptual difference.
- A **`fitExtent`** geo tile (auto-fit a feature collection to the frame),
  replacing manual scale/center in the geo demo.

## 8. Implementation order (topological, by leverage/risk)

1. **Color spaces** (§3) — self-contained port; unblocks perceptual scale
   interpolation and proper Lab modifiers.
2. **Consumption surface** (§4) — `samples` package first (single source of
   truth), then the `charts/static` → 28 extension, then the unified color API.
3. **Geo completeness** (§5) — `GeoPath` measurement → `fitExtent`/`fitSize` →
   catalog → optional TopoJSON.
4. **Demo + docs** (§6–§7): wire the demo to `samples`, add the color-space and
   `fitExtent` tiles, update `README`/`USAGE`/`NOTES`, `make golden` idempotent,
   `make ci` green.

## 9. Scope summary for v7

**Color**: HSL/Lab/Lch added to `internal/d3/color` behind the `Color`
interface, with in-space interpolation; opt-in perceptual interpolation for
sequential/diverging scales and palette sampling, and Lab-based lightness
modifiers — all defaulting to today's RGB behaviour (no golden drift).

**Consumption**: a public typed `samples` package, `charts/static` extended from
bar/line/pie to all 28 chart families, and one unifying (additive)
color-setting API over the five existing config shapes.

**Geo**: `GeoPath` bounds/area/centroid, `fitExtent`/`fitSize`, a fuller
projection catalog, and optional TopoJSON decoding.

**Tests/demo**: ports golden-tested vs d3; static/samples/color-API tests; the
demo re-based on `samples` with color-space + `fitExtent` tiles. All existing
goldens stay byte-stable.

## 10. After v7 — maintenance only

v7 closes the consolidated backlog. There is **no planned v8**: with SVG parity
(v1–v3), consumability (v4), fidelity (v5), scale/Canvas (v6), and completeness
(v7) all delivered, the remaining work is maintenance — tracking upstream
nivo/d3 changes, fixing bugs, and responding to consumer requests — rather than
a themed release. [`docs/PLAN-deferred.md`](PLAN-deferred.md) is expected to be
empty of committed work at this point; anything new lands there as it is
discovered.
