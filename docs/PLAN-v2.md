# templ-charts — Implementation Plan (v2)

## 1. Vision

v1 delivered the framework + three chart types (**bar**, **line**, **pie**)
on a deep, reusable foundation. v2 broadens the chart catalog with **ten new
chart types** that reuse that foundation — adding *no* heavy new d3 ports —
and modernizes two cross-cutting concerns: **interactivity** (move ephemeral
hover off the server) and **responsiveness + accessibility**.

v2 deliberately stays **SVG-only** (Canvas remains deferred) and **reuse-only**
for d3: every new chart is buildable from the ports already present
(`internal/d3/{shape,scale,array,format,timeformat,color}`), with the single
exception of a small `Quantile` helper added to `d3/array` for boxplot. The
two scaffold-only packages from v1 — `grid` and `polar-axes` — are activated
and exercised by these charts.

See [`docs/PLAN.md`](PLAN.md) for the v1 plan whose structure and conventions
this document mirrors.

## 2. Confirmed decisions

| Concern | Decision |
|---|---|
| Chart scope | Ten "quick-win" Tier-2 charts that reuse the existing foundation |
| d3 strategy | **Reuse-only.** No new d3 modules. Only add `d3/array.Quantile` (boxplot) |
| Rendering | Server-rendered SVG via templ (Canvas still skipped) |
| Interactivity | **Hybrid**: thin client-side JS for hover/crosshair/tooltip; HTMX for state changes (toggle/activate) |
| Responsiveness | Fluid `viewBox` + `width:100%`; optional client resize-observer re-fetch |
| Accessibility | ARIA roles/labels/desc, titled `<svg>`, keyboard-focusable series; builds on `rects.NodeA11yProps` + pie `Role` |
| Scaffolds | Activate `charts/grid` (heatmap/waffle/calendar) and `charts/polar-axes` (radar/radial-bar) |
| Determinism | Fix nondeterministic CSS ordering in `styleFromMap` (prerequisite for stable goldens) |
| Demo | Extend `examples/app` with a page per new chart family, same Registry/handler pattern |

## 3. Chart roster

Each new chart follows the standard v1 package anatomy — `types.go` →
`defaults.go` → `hooks.go` (`Use{Chart}` orchestrator) → `render.go` →
`*.templ` (+ `compute/` where layout math warrants) — with golden SVG
snapshots under `testdata/golden/`. Defaults are ported **verbatim** from the
cited nivo package. All charts reuse the base packages (`core`, `theming`,
`scales`, `colors`, `axes`, `legends`, `tooltip`, `text`, `annotations`,
`htmx`); the table notes only what each adds on top.

| Chart | Reuses (beyond base) | New compute | d3 needs | nivo source |
|---|---|---|---|---|
| heatmap | `scales` (band x/y), `axes`, sequential/diverging colors, continuous legend | band-scale cell layout, color-by-value | existing | `@nivo/heatmap` |
| scatterplot | `axes`, `core.DotsItem`, annotations | node placement on x/y scales | existing | `@nivo/scatterplot` |
| radar | `polar-axes`, `arcs` angle helpers, line/area generator | polar point projection, grid rings | existing | `@nivo/radar` |
| radial-bar | `polar-axes`, `arcs` | angle/radius bands → arcs | existing | `@nivo/radial-bar` |
| stream | `axes`, area generator, `d3.Stack` | stacked-area baseline w/ wiggle/silhouette offsets | existing | `@nivo/stream` |
| waffle | `grid`, `rects` | fill cells by value share | existing | `@nivo/waffle` |
| calendar | `grid`-style layout, time scale, time-format | day/week/month cell geometry | existing | `@nivo/calendar` |
| bullet | `axes`, `rects` | range/measure/marker rects | existing | `@nivo/bullet` |
| funnel | `axes`, custom path | trapezoid band geometry | existing | `@nivo/funnel` |
| boxplot | `axes`, `rects`, line primitives | quartile/whisker aggregation | **`d3/array.Quantile`** | `@nivo/boxplot` |

### 3.1 heatmap
- **Data**: cells indexed by `(serieId, x)` → value; color mapped through a
  sequential/diverging color scale.
- **Defaults** (verbatim `heatmap/src/defaults.ts`): `layers:['grid','axes','cells','legends','annotations']`,
  `forceSquare:false`, `xInnerPadding/xOuterPadding/yInnerPadding/yOuterPadding:0`,
  `opacity:1`, `activeOpacity:1`, `inactiveOpacity:0.15`, `borderWidth:0`,
  `borderColor:{from:color, modifiers:[[darker,0.8]]}`, `enableGridX/Y:false`,
  `enableLabels:true`, `label:'formattedValue'`, `labelTextColor:{from:color, modifiers:[[darker,2]]}`,
  `colors:{type:sequential, scheme:'brown_blueGreen'}`, `emptyColor:'#000000'`,
  `hoverTarget:'rowColumn'`, SVG extras: `borderRadius:0`, `cellComponent:'rect'`.
- **Compute**: band scales on both axes (the x band over x-values, the y band
  over serie ids — built with the non-reversed range so series run top→bottom
  like nivo); cell fill via the `colors` sequential/diverging scale; optional
  `forceSquare` layout centering. (`grid.GenerateGrid` is *not* used here — it's
  waffle's path; heatmap is band-scale driven like nivo.)
- **Layers**: grid → axes → cells → legends → annotations.

### 3.2 scatterplot
- **Data**: series of `{x, y}` nodes on two continuous (or time) scales.
- **Defaults** (verbatim `scatterplot/src/props.tsx`): `xScale:{type:linear, min:0, max:'auto'}`,
  `yScale:{type:linear, min:0, max:'auto'}`, `enableGridX/Y:true`,
  `axisTop/Right:null`, `axisBottom/Left:{}`, `nodeSize:9`, `colors:{scheme:'nivo'}`,
  `blendMode:'normal'`, `isInteractive:true`, `markers/legends/annotations:[]`.
- **Compute**: `ComputeXYScalesForSeries` (already exists) → node positions;
  render via `core.DotsItem`.
- **Caveat**: ship **direct point hover** only. nivo's voronoi-mesh hover
  needs d3-delaunay → **deferred to v3** (see §11).

### 3.3 radar — *activates `polar-axes`*
- **Data**: indices around the circle, one value per series per index.
- **Defaults** (verbatim `radar/src/defaults.ts`): `layers:['grid','layers','slices','dots','legends']`,
  `maxValue:'auto'`, `rotation:0`, `curve:'linearClosed'`, `borderWidth:2`,
  `borderColor:{from:color}`, `gridLevels:5`, `gridShape:'circular'`,
  `gridLabelOffset:16`, `enableDots:true`, `dotSize:6`, `colors:{scheme:'nivo'}`,
  `fillOpacity:0.25`, `blendMode:'normal'`.
- **Compute**: angle per index = `2π·i/n`; project `(value→radius)` via a
  linear radius scale; close the path with the `linearClosed` curve (ported).
  Grid rings/labels from `polar-axes` (`CircularAxis`, `PolarGrid`).

### 3.4 radial-bar — *activates `polar-axes`*
- **Data**: groups → stacked categories rendered as arcs in polar space.
- **Defaults** (verbatim `radial-bar/src/props.ts`): `maxValue:'auto'`,
  `layers:['grid','tracks','bars','labels','legends']`, `startAngle:0`,
  `endAngle:270`, `innerRadius:0.3`, `padding:0.2`, `padAngle:0`, `cornerRadius:0`,
  `enableTracks:true`, `tracksColor:'rgba(0,0,0,.15)'`, `enableRadialGrid:true`,
  `enableCircularGrid:true`, `colors:{scheme:'nivo'}`, `borderWidth:0`,
  `borderColor:{from:color, modifiers:[[darker,1]]}`, `enableLabels:false`,
  `labelsSkipAngle:10`, `labelsRadiusOffset:0.5`, `labelsTextColor:{theme:'labels.text.fill'}`.
- **Compute**: radius band scale + angle scale → `arcs.Arc`; reuse the v1 arc
  generator (cornerRadius/padAngle) and `polar-axes` radial/circular grids.

### 3.5 stream
- **Data**: stacked series over an index axis, smooth areas.
- **Defaults** (verbatim `stream/src/props.ts`): `label:'id'`, `order:'none'`,
  `offsetType:'wiggle'`, `curve:'catmullRom'`, `axisBottom/Left:{}`,
  `axisTop/Right:null`, `enableGridX:false`, `enableGridY:true`,
  `colors:{scheme:'nivo'}`, `fillOpacity:1`, `borderWidth:0`,
  `borderColor:{from:color, modifiers:[[darker,1]]}`, `enableDots:false`,
  `dotSize:6`, layers `['grid','axes','layers','dots','slices','legends']`,
  `motionConfig:'default'`.
- **Compute**: `d3.Stack` with **wiggle/silhouette/expand** offsets. v1 has the
  diverging offset; **add the remaining offsets** to `internal/d3/shape/stack.go`.
  Areas drawn with the existing area generator.

### 3.6 waffle — *activates `grid`*
- **Data**: parts of a whole laid into an N×N cell grid.
- **Defaults** (verbatim `waffle/src/defaults.ts`): `hiddenIds:[]`,
  `fillDirection:'top'`, `padding:1`, `colors:{scheme:'nivo'}`,
  `emptyColor:'#cccccc'`, `emptyOpacity:1`, `borderRadius:0`, `borderWidth:0`,
  `borderColor:{from:color, modifiers:[[darker,1]]}`, SVG layers
  `['cells','areas','legends']`, `motionStagger:0`.
- **Compute**: `grid.GenerateGrid` + `fillDirection` to assign cells to data
  by cumulative value share; cells via `rects.RoundedRect`.

### 3.7 calendar — *activates `grid`-style layout*
- **Data**: `{day:'YYYY-MM-DD', value}` over a date range.
- **Defaults** (verbatim `calendar/src/props.ts`): `colors:['#61cdbb','#97e3d5','#e8c1a0','#f47560']`,
  `align:'center'`, `direction:'horizontal'`, `emptyColor:'#fff'`, `minValue:0`,
  `maxValue:'auto'`, `yearSpacing:30`, `monthBorderWidth:2`, `monthBorderColor:'#000'`,
  `monthSpacing:0`, `daySpacing:0`, `dayBorderWidth:1`, `dayBorderColor:'#000'`,
  `legends:[]`, `role:'img'`.
- **Compute**: day-cell geometry from the date range using the ported **time
  scale** + **time-format**; quantize color scale over the value domain.

### 3.8 bullet
- **Data**: per-row `ranges[]`, `measures[]`, `markers[]` on a shared value scale.
- **Defaults** (verbatim `bullet/src/props.ts`): `layout:'horizontal'`,
  `reverse:false`, `spacing:30`, `minValue:0`, `maxValue:'auto'`,
  `axisPosition:'after'`, `titlePosition:'before'`, `titleAlign:'middle'`,
  `rangeColors:'seq:cool'`, `measureColors:'seq:red_purple'`,
  `markerColors:'seq:red_purple'`, `rangeBorderWidth:0`, `rangeBorderColor:{from:color}`.
- **Compute**: linear value scale → range/measure rects (`rects`) + marker lines;
  one axis (`axes`) per the layout.

### 3.9 funnel
- **Data**: ordered parts; widths interpolate between successive values.
- **Defaults** (verbatim `funnel/src/props.tsx`): `layers:['separators','parts','labels','annotations']`,
  `direction:'vertical'`, `interpolation:'smooth'`, `spacing:0`, `shapeBlending:0.66`,
  `colors:{scheme:'nivo'}`, `fillOpacity:1`, `borderWidth:6`, `borderColor:{from:color}`,
  `borderOpacity:0.66`, `enableLabel:true`, `labelColor:{theme:'background'}`,
  `enableBeforeSeparators/AfterSeparators:true`, `currentPartSizeExtension:0`.
- **Compute**: trapezoid band paths between successive part widths; `smooth`
  uses the ported bezier curve, `linear` uses straight edges.
- **Note**: the `d3.Stack` offsets stream needs (wiggle/silhouette/expand, plus
  diverging) are **already ported** in `internal/d3/shape/stack.go`; reuse directly.

### 3.10 boxplot — *adds `d3/array.Quantile`*
- **Data**: raw values grouped by `group`/`subGroup`; summarized to quantiles.
- **Defaults** (verbatim `boxplot/src/props.ts`): `value:'value'`, `groupBy:'group'`,
  `quantiles:[0.1,0.25,0.5,0.75,0.9]`, `layout:'vertical'`, `minValue/maxValue:'auto'`,
  `valueScale:{type:linear}`, `indexScale:{type:band, round:true}`, `padding:0.1`,
  `innerPadding:6`, `opacity:1`, `activeOpacity:1`, `inactiveOpacity:0.25`,
  `axisBottom/Left:{}`, `axisTop/Right:null`, `enableGridX:false`, `enableGridY:true`,
  `valueFormat:toPrecision(4)`, `colorBy:'subGroup'`, `colors:{scheme:'nivo'}`,
  layers `['grid','axes','boxPlots','markers','legends','annotations']`.
- **Compute**: per-group quantile aggregation via the new
  `internal/d3/array.Quantile` (linear interpolation, matching d3); box + whisker
  geometry from `rects` + line primitives.

## 4. Activating the scaffolds

Both scaffold packages already ship the hard geometry; v2 wires them into full
chart packages.

- **`charts/grid`** (heatmap, waffle, calendar): has `GenerateGrid[C]`,
  `ComputeCellDimensions`, `GridFillDirection`, polygon/bounding-box helpers.
  Charts consume these directly for cell layout; no new geometry needed.
- **`charts/polar-axes`** (radar, radial-bar): has `CircularAxis`, `RadialAxis`,
  `PolarGrid` components plus tick/position compute (`angleScaleToTicks`,
  `radialTickPositions`, `arcPath`). radar/radial-bar render their grids and
  axis labels through these; arcs come from the v1 `charts/arcs` package.

No new exported types are required in `grid`/`polar-axes` beyond what exists;
the work is the chart packages that *use* them (Props/Result/hooks/render/templ).

## 5. Interactivity v2 — hybrid client/server model

v1 routes **every** hover/mousemove through HTMX, re-rendering SVG server-side
(NOTES.md flags this for line crosshair). v2 splits responsibilities:

- **Client-side (new, thin JS — no framework, served like the existing inline
  demo script):** ephemeral interactions that must track the cursor —
  tooltip show/position, crosshair, and nearest-point/cell hit-testing.
  Charts emit plain `data-*` attributes (e.g. per-node `data-x/data-y/data-id`,
  per-cell bounds) that the script reads to find the hovered datum and render a
  tooltip locally. No server round-trip per mousemove.
- **Server-side (HTMX, unchanged):** *state-changing* actions that alter what
  is rendered — series **toggle**, **active** arc/segment, and any persisted
  selection. These keep using the `htmx.Registry` + `Handler` endpoints.

Concretely: the `line` mesh/slice hover migrates onto the client layer
(resolving the NOTES.md gap), and all ten new charts adopt the client layer for
hover tooltips; only toggle/activate hit the server. The boundary is documented
so non-interactive/static embedding still works with zero JS.

## 6. Responsive + accessibility

- **Responsive** ✅ *(done — CSS-native fluid scaling + a reusable mount)*:
  - `core.SvgWrapperProps.Responsive` (threaded through a `Responsive bool` on
    every chart's props) keeps the `viewBox` and intrinsic `width`/`height`
    while emitting an inline `style="width:100%;height:auto;display:block"`, so
    a chart scales fluidly to its container with **zero JS** and zero consumer
    CSS. Defaults to `false`, so existing fixed-size renders (and goldens) are
    unchanged.
  - `htmx.Mount(MountProps{...})` is the library-provided container that wires a
    rendered SVG into the interactivity layer: it emits the `#chart-<id>`
    container (the OOB swap target), the `#tooltip-<id>` swap target, and the
    `mouseleave` reset, matching the `/charts/{id}` routes and `chart-`/
    `tooltip-` id conventions. It ships **no styling** (visual chrome stays the
    consumer's concern via the `tc-chart`/`tc-chart-tooltip` class hooks or the
    `Class` field). The demo's `ChartCard` now consumes `htmx.Mount` instead of
    hand-rolling the markup, so the demo demonstrates the public API rather than
    hiding it.
  - *Still optional / future*: the client-side `ResizeObserver` re-fetch for a
    pixel-accurate re-render on axes-heavy charts where tick/label density
    matters. Purely cosmetic scaling is fully covered by `Responsive` above and
    needs no JS.
- **Accessibility**: build on `rects.NodeA11yProps` and pie's `Role` default.
  Add `<title>`/`<desc>` to the `SvgWrapper`, ARIA `role`/`aria-label`/
  `aria-describedby` on charts, and keyboard-focusable series/segments where
  meaningful (tabindex + focus styles). Apply uniformly to v1 and v2 charts.

## 7. `internal/d3` additions

- **`d3/array.Quantile`** (+ small helpers): linear-interpolation quantile over
  a sorted slice, matching d3-array semantics. Sole new d3 code in v2; gated by
  boxplot. Reuses existing `Ascending`/sort utilities.
- **`d3/shape/stack.go`**: the **wiggle / silhouette / expand / diverging**
  offsets stream needs are **already ported** — no work required (verified in
  Phase 1).
- `d3/color` remains RGB-only — sufficient for v2 (color modifiers/interpolators
  already work).

**Not ported (and the charts they gate, all deferred — see §11):** d3-hierarchy
(treemap, sunburst, icicle, circle-packing, tree), d3-force (network, swarmplot),
d3-sankey (sankey), d3-chord (chord), d3-geo (geo/choropleth), d3-delaunay
(voronoi diagram + voronoi-mesh hover for line/scatter).

## 8. Build / test

- **Golden SVG snapshots** per new chart (`charts/{chart}/testdata/golden/`),
  asserted via `internal/golden.Assert`; regenerate with `make golden`.
- **Path-string golden tests** for new geometry: funnel trapezoids, radar polar
  projection + `linearClosed` closure, calendar day cells, radial-bar arcs,
  and boxplot whisker/box paths.
- **Unit tests**: `d3/array.Quantile` vs known d3 outputs; new stack offsets vs
  d3-shape; `grid`/`polar-axes` compute (extend existing tests).
- **Determinism fix (prerequisite)**: sort keys in `styleFromMap` so emitted
  CSS/style ordering is stable (NOTES.md #8); without this the new goldens flake.
- CI unchanged: `make ci` (`make lint` + `make test`) from root.

## 9. Demo app — `examples/app`

Extend the existing stdlib `net/http` app with one page per new chart family,
reusing the `demos/` + `handlers/` + `htmx.Registry` pattern and the shared
template/CSS:
- `/heatmap`, `/scatterplot`, `/radar`, `/radial-bar`, `/stream`, `/waffle`,
  `/calendar`, `/bullet`, `/funnel`, `/boxplot` (group related ones if tidier).
- Add a nav entry per page (mirrors `/bar`,`/line`,`/pie`,`/palettes`,`/themes`).
- Load the new client-side interactivity script once in `layout.templ`.

## 10. Implementation order (topological)

1. **Foundations** ✅ *(done)*: fixed `styleFromMap` determinism (sorted keys);
   added `d3/array.Quantile`/`QuantileSorted`. (Stream stack offsets were already
   ported — no work needed.)
2. **Grid-based charts** ✅ *(done)*: heatmap (band-scale driven + continuous
   color/legend), waffle (built on `grid.GenerateGrid` — activates the scaffold),
   calendar (ported date math, quantized colors, month/year legends). Each has a
   full package, golden tests, and a demo page. *Deferred within this phase*:
   calendar's month outline-path border, and waffle's polygon "areas" layer.
3. **Activate `polar-axes`** ✅ *(done)*: radar (linear radius scale + projected
   closed polygons via the `linearClosed` curve, circular/polygon grid levels,
   index labels, dots, legend) and radial-bar (stacked arcs via `charts/arcs`,
   background tracks, and `polar-axes` PolarGrid/RadialAxis/CircularAxis). Added
   `scales.NewLinearScaleWithRange`/`NewBandScaleWithRange` so the polar charts
   can map onto angle/radius ranges (the gap noted in `PolarGridProps`). Each has
   a full package, golden tests, and a demo page. *Deferred within this phase*:
   the interactive `slices` (radar) / hover-tooltip (radial-bar) layers, which
   arrive with the Phase 5 client layer.
4. **Cartesian batch** ✅ *(done)*: scatterplot (XY scales + `core.DotsItem`),
   stream (stacked areas with wiggle/silhouette/expand offsets via `d3.Stack` +
   the area generator), bullet (range/measure rects + marker lines on a per-row
   value scale, sequential `seq:*` colors), funnel (smooth/linear trapezoid
   bands via the transposed area generator + side borders + separators), and
   boxplot (per-group quantile summaries via `d3/array.QuantileSorted` → box +
   whisker glyphs). Each has a full package, golden tests, and a demo page.
   *Deferred within this phase* (Phase 5 client layer): scatterplot point/mesh
   hover, stream slices/dots, and the bullet/funnel/boxplot hover tooltips.
5. **Interactivity layer** ✅ *(done)*: shipped `charts/interact` — a
   dependency-free client-side module (`Script` / `ScriptTag`) that handles
   ephemeral hover entirely in the browser off `data-tc-*` attributes:
   element-delegation tooltips (any element with `data-tc-tooltip`) and
   nearest-point **mesh** hover + crosshair (`data-tc-mesh`). It reuses the
   `.tc-chart` / `.tc-chart-tooltip` conventions and lazily creates the tooltip
   element for zero-markup static embeds. Shared components (`core.DotsItem`,
   `arcs.ArcShape`) gained tooltip fields; all ten v2 charts gained an
   `Interactive` flag that emits per-element tooltips (heatmap migrated off its
   htmx hover); `line` gained `ClientHover` which routes mesh/slice through the
   client layer (resolving the v1 NOTES.md per-mousemove round-trip) while
   keeping legend toggle on the server. The demo loads the script once in
   `layout.templ` and enables it across the new chart pages. Defaults are off,
   so static (zero-JS) embedding and all existing goldens are unchanged.
6. **Responsive + a11y** pass across v1 + v2 charts. *Responsive done* (the
   `Responsive` prop on `SvgWrapper`/all charts + the `htmx.Mount` container —
   see §6). *A11y done*: `core.SvgWrapper` gained `Title`/`Desc` that render
   `<title>`/`<desc>` as the first `<svg>` children (the SVG-native accessible
   name/description for `role="img"`), threaded uniformly through every v1 + v2
   chart's props alongside the existing `Role`/`AriaLabel`/`AriaLabelledBy`/
   `AriaDescribedBy`/`IsFocusable` (all charts default `role="img"`; the svg is
   keyboard-focusable via `IsFocusable`). Defaults stay empty so static goldens
   are unchanged. *Still deferred (optional/future, per §6)*: the client-side
   `ResizeObserver` re-fetch for pixel-accurate re-render on axes-heavy charts —
   cosmetic fluid scaling is already fully covered by `Responsive`, so this is
   left for a later pass.
7. **Demo pages**, `README`/`docs` updates, golden regeneration ✅ *(done)*:
   a demo page per chart family is wired (§9); the root `README`, the demo
   `examples/app/README`, and the demo index intro are updated for the v2
   catalog + hybrid interactivity + a11y; the `make golden` target now covers
   every golden-owning package and all snapshots are regenerated/current
   (including the previously stale `heatmap-basic` axis-tick offsets). `make
   ci` is fully green.

## 11. Explicitly deferred (v3+)

- **Charts** (need d3 modules not ported): treemap, sunburst, icicle,
  circle-packing, tree (d3-hierarchy); network, swarmplot (d3-force); sankey
  (d3-sankey); chord (d3-chord); geo/choropleth (d3-geo); voronoi diagram and
  voronoi-**mesh** hover for line/scatter (d3-delaunay); plus
  parallel-coordinates and marimekko (custom layout math).
- **Deferred by scope (no new d3 needed)**: bump (ranking-over-time line
  variant) — a natural fast-follow that reuses only existing ports; left out of
  the v2 batch to keep it at ten, and a strong first candidate for v2.1.
- **Rendering**: Canvas output path (SVG-only continues in v2).
- **Color spaces**: HSL/Lab in `d3/color` (only if a future chart needs them).
