# templ-charts — Implementation Plan (v3)

## 1. Vision

v1 built the framework + three charts (**bar**, **line**, **pie**) on a deep,
reusable foundation. v2 added **ten** "quick-win" charts that reused that
foundation with *no* new d3 ports, plus hybrid interactivity, responsiveness,
and accessibility.

**v3 closes the gap to full nivo chart parity.** It ports the five remaining d3
modules that v2 explicitly deferred — **d3-hierarchy, d3-delaunay, d3-force,
d3-sankey, d3-chord, d3-geo** — and delivers the **15 remaining nivo chart
packages** that depend on them, plus four charts that need no new d3 at all.
After v3, the only material nivo feature not present in templ-charts is the
**Canvas** rendering path (still deferred — see §11); every SVG chart type nivo
ships has a templ-charts equivalent.

v3 stays **SVG-only** and follows the same per-package anatomy, verbatim
defaults, and golden-test discipline as v1/v2. The new work is concentrated in
`internal/d3/`: each new chart family is gated by exactly one new port, and the
ports are sequenced by leverage (charts unblocked per port) and risk.

See [`docs/PLAN.md`](PLAN.md) (v1) and [`docs/PLAN-v2.md`](PLAN-v2.md) (v2) for
the structure and conventions this document mirrors.

## 2. Confirmed decisions

| Concern | Decision |
|---|---|
| Chart scope | The **15 remaining nivo chart packages** (geo ships two components) + polar-bar |
| d3 strategy | Port the **five deferred d3 modules**: hierarchy, delaunay, force, sankey, chord, geo. Add `curveBumpX/Y` to `d3/shape`. |
| Rendering | Server-rendered SVG via templ (**Canvas still deferred** — §11) |
| Determinism | Every new port must be deterministic server-side. d3-force ports d3's built-in **LCG** + runs a **fixed tick count** (nivo default 120) — no wall-clock, no `Math.random`. |
| Interactivity | Reuse v2's `charts/interact` client layer. **d3-delaunay unblocks true voronoi-mesh hover**, retrofitted into line/scatterplot/bump/swarmplot/tree (replacing v2's direct-point-only hover). State changes stay on HTMX. |
| Accessibility / responsive | Reuse v2's `Responsive` prop, `<title>`/`<desc>`, `role`/ARIA — applied uniformly to all v3 charts. |
| geo data | Consume **GeoJSON** features supplied by the caller (as nivo does). **No TopoJSON** decoder in-tree; document converting to GeoJSON upstream. |
| Sequencing | Port-by-port phases, ordered by leverage/risk: hierarchy → delaunay → force → sankey → chord → geo, with the no-new-d3 charts as an opening fast-follow. |
| Demo | Extend `examples/app` with a page per new chart family, same Registry/handler pattern. |

## 3. `internal/d3` — the new ports

This is the heart of v3. Each port is validated by golden/path-string tests
against real d3 output, exactly like the v1 shape/scale ports.

### 3.1 `internal/d3/hierarchy` — **highest leverage (5 charts)**
Pure layout math, no randomness, fully deterministic. Public surface (union
across treemap/sunburst/icicle/circle-packing/tree):

- **Core**: `Hierarchy(data, children)` → root `*Node`; `Node.Sum(valueFn)`,
  `.Sort(cmp)`, `.Count()`, `.Each/.EachBefore/.EachAfter(fn)`, `.Descendants()`,
  `.Leaves()`, `.Ancestors()`, `.Links()`; fields `Parent`, `Children`, `Depth`,
  `Height`, `Value`.
- **Layouts** (each sets coordinates on nodes):
  - `Treemap()` — `.Size`, `.Round`, `.PaddingInner/Outer/Top/Right/Bottom/Left`,
    `.Tile(fn)`; tiling fns `TreemapSquarify` (default, needs golden ratio split),
    `TreemapBinary`, `TreemapDice`, `TreemapSlice`, `TreemapSliceDice`. Sets
    `x0,y0,x1,y1`.
  - `Partition()` — `.Size`, `.Round`, `.Padding`. Sets `x0,y0,x1,y1`
    (rectangular for icicle; sunburst maps to `[2π, r²]` then `√` for radii).
  - `Pack()` — `.Size`, `.Padding`; enclosing-circle algorithm (Welzl). Sets
    `x,y,r`.
  - `Tree()` / `Cluster()` — Reingold–Tilford tidy-tree (tree = mode
    `dendogram`/`tree`). Sets `x,y` in a unit space; charts rescale.
- **Risk**: `Pack` (smallest-enclosing-circle) and `TreemapSquarify` are the
  fiddly bits; golden-test node coordinates against d3-hierarchy fixtures.

### 3.2 `internal/d3/delaunay` — **medium; unblocks voronoi + 5 mesh retrofits**
Bowyer–Watson Delaunay triangulation + Voronoi dual + point location.
Deterministic. Public surface:
- `NewDelaunayFrom(points [][2]float64) *Delaunay`
- `(*Delaunay).Find(x, y float64) int` — nearest-site index (edge-walk); this is
  the **mesh-hover** primitive.
- `(*Delaunay).Voronoi(bounds [4]float64) *Voronoi`
- `(*Voronoi).CellPolygon(i int) [][2]float64`, `.Render() string`,
  `.RenderCell(i) string` (SVG path strings)
- `(*Delaunay).RenderPoints(radius) string`, `.Neighbors(i)` (voronoi "links")
- **Estimated** ~800–1200 Go LOC. Self-contained, no spherical math.

### 3.3 `internal/d3/force` — **iterative but deterministic**
Velocity-Verlet integrator with d3's built-in **LCG** for jiggle (so runs are
reproducible with no `Math.random`). Server-side model: attach forces, `Stop()`,
then `Tick(n)` a **fixed** count (nivo default 120) — no alpha-based auto-stop.
- `NewSimulation(nodes)` → `.Force(name, f)`, `.Tick(n)`, `.Stop()`, `.Nodes()`;
  node fields `X,Y,Vx,Vy,Index` (deterministic phyllotaxis init when unset).
- Forces: `ForceLink(links)` (`.Distance`, `.Strength`), `ForceManyBody`
  (`.Strength`, `.DistanceMin`, `.DistanceMax`), `ForceCenter(x,y)`,
  `ForceX/ForceY(accessor)` (`.Strength`), `ForceCollide(radius)`.
- `ForceManyBody` may use all-pairs (node counts here are small); Barnes–Hut
  quadtree is an optional optimization, not required for parity.
- **Risk**: matching d3's jiggle LCG exactly for byte-identical goldens; if that
  proves brittle, assert layout *stability/shape* (bounding box, non-overlap)
  rather than exact coordinates, and document the golden strategy.

### 3.4 `internal/d3/sankey` — **single-pass, deterministic, small**
- `Sankey()` → `.NodeAlign(fn)` (`SankeyLeft/Right/Center/Justify`),
  `.NodeSort(fn)`, `.LinkSort(fn)`, `.NodeWidth`, `.NodePadding`, `.Size`,
  `.NodeId(fn)`. `layout(graph)` mutates nodes (`x0,x1,y0,y1,value,depth,index,
  sourceLinks,targetLinks`) and links (`y0,y1,width`).
- Link paths use already-ported `curveMonotoneX/Y`.

### 3.5 `internal/d3/chord` — **single-pass, deterministic, small**
- `Chord()` → `.PadAngle`; `chord(matrix [][]float64)` → `{Groups, []Ribbon}`
  with `startAngle/endAngle/value/index` per group and `source/target` per
  ribbon. `Ribbon()` → `.Radius`; renders the ribbon path. Arc paths reuse the
  existing `charts/arcs` / `d3/shape.Arc`.

### 3.6 `internal/d3/geo` — **heaviest; lowest chart-count**
Subset sufficient for nivo's geo package. Honest scope: ~5–8k Go LOC.
- **Projections** (the 10 nivo exposes): `GeoMercator`, `GeoEquirectangular`,
  `GeoTransverseMercator`, `GeoNaturalEarth1`, `GeoEqualEarth`,
  `GeoAzimuthalEqualArea`, `GeoAzimuthalEquidistant`, `GeoGnomonic`,
  `GeoOrthographic`, `GeoStereographic`. Each: `.Scale`, `.Translate`,
  `.Rotate(λ,φ,γ)`, `.Project(lon,lat)→(x,y)`.
- **Rendering**: `GeoPath(projection)` → GeoJSON geometry to SVG path string,
  incl. antimeridian/clip handling per projection family.
- **Graticule**: `GeoGraticule()` → lat/lon mesh feature.
- Spherical primitives: great-circle interpolation, rotation, clip.
- **Recommendation**: this is optional/last. Ship a **minimal viable** geo
  first (Mercator + equirectangular + `GeoPath` + graticule, ~500–800 LOC) to
  cover the common case, then add projections incrementally. Consider carving
  geo into its own **v3.1** release if the projection math slips.

### 3.7 `internal/d3/shape` additions
- **`curveBumpX` / `curveBumpY`** — cubic bump curves for `tree` links and
  `bump`'s smooth interpolation. Small. (monotone, basis, cardinal, catmullRom,
  step\* are already ported.)

## 4. Chart roster

Standard v1/v2 package anatomy per chart — `types.go` → `defaults.go` →
`hooks.go` (`Use{Chart}`) → `render.go` → `*.templ` (+ `compute/` where warranted)
— with golden SVG snapshots under `testdata/golden/`. Defaults ported
**verbatim** from the cited nivo package (captured below). All reuse the base
packages; the table notes only what each adds.

| Chart | Gated by | Reuses (beyond base) | nivo source |
|---|---|---|---|
| **bump** | — (uses existing scales/shape; mesh via §3.2) | point+linear scales, line gen, `interact` mesh | `@nivo/bump` |
| **marimekko** | — | `d3.Stack`, `axes`, `legends` | `@nivo/marimekko` |
| **parallel-coordinates** | — | linear+point scales, line gen, `axes` | `@nivo/parallel-coordinates` |
| **polar-bar** | — | `arcs`, `polar-axes`, `d3.Stack`, band+linear scales | `@nivo/polar-bar` |
| **treemap** | d3-hierarchy | `rects`, ordinal colors, labels | `@nivo/treemap` |
| **sunburst** | d3-hierarchy | `arcs`, arc labels, colors | `@nivo/sunburst` |
| **icicle** | d3-hierarchy | `rects`, labels, colors | `@nivo/icicle` |
| **circle-packing** | d3-hierarchy | colors, labels | `@nivo/circle-packing` |
| **tree** | d3-hierarchy + `curveBumpX/Y` | linear scale, link gen, `interact` mesh | `@nivo/tree` |
| **voronoi** | d3-delaunay | linear scales | `@nivo/voronoi` |
| **network** | d3-force | `annotations`, colors | `@nivo/network` |
| **swarmplot** | d3-force | `axes`, scales, `interact` mesh | `@nivo/swarmplot` |
| **sankey** | d3-sankey | `legends`, gradients (`core` defs), labels | `@nivo/sankey` |
| **chord** | d3-chord | `arcs`, `legends`, labels | `@nivo/chord` |
| **geo** (GeoMap + Choropleth) | d3-geo | sequential colors, continuous legend | `@nivo/geo` |

### 4.1 No-new-d3 fast-follows

**bump** — `layers:['grid','axes','labels','lines','points','mesh']`,
`interpolation:'smooth'`, `xPadding:0.6`, `xOuterPadding:0.5`, `yOuterPadding:0.5`,
`lineWidth:2`, `activeLineWidth:4`, `inactiveLineWidth:1`, `opacity:1`,
`activeOpacity:1`, `inactiveOpacity:0.3`, `startLabel:false`, `endLabel:true`,
`pointSize:6`, `activePointSize:8`, `inactivePointSize:4`, `useMesh:false`,
`debugMesh:false`. Point scale (X) + linear ranking scale (Y); lines via the area/
line generator. **Ships with direct hover; mesh layer lights up after §3.2.**

**marimekko** — `layers:['grid','axes','bars','legends']`, `layout:'vertical'`,
`offset:'none'`, `outerPadding:0`, `innerPadding:3`, `enableGridX:false`,
`enableGridY:true`, `borderWidth:0`, `isInteractive:true`. Variable-width stacked
bars: primary dimension (thickness) is value-driven, secondary is the stacked
category via `d3.Stack`.

**parallel-coordinates** — `layers:['lines','axes','legends']`,
`layout:'horizontal'`, `curve:'linear'`, `lineWidth:2`, `lineOpacity:0.5`,
`colors:{scheme:'category10'}`, `axesTicksPosition:'after'`, `isInteractive:true`.
One linear/point scale per variable; polylines across parallel axes.

**polar-bar** (distinct from radial-bar — stacked, full-circle) —
`layers:['grid','arcs','axes','labels','legends']`, `indexBy:'id'`,
`keys:['value']`, `startAngle:0`, `endAngle:360`, `innerRadius:0`,
`cornerRadius:0`, `enableRadialGrid:true`, `enableCircularGrid:true`,
`enableArcLabels:false`, `arcLabel:'formattedValue'`, `arcLabelsRadiusOffset:0.5`,
`isInteractive:true`. `d3.Stack` layered radially; band scale (angle) + linear
scale (radius) → `arcs.Arc`; reuses `polar-axes` grids. (radial-bar = single
value + tracks, 0–270°, `innerRadius:0.3`; polar-bar = stacked, 0–360°,
`innerRadius:0`.)

### 4.2 d3-hierarchy charts

**treemap** — `layers:['nodes']`, `identity:'id'`, `value:'value'`,
`tile:'squarify'`, `leavesOnly:false`, `innerPadding:0`, `outerPadding:0`,
`colors:{scheme:'nivo'}`, `colorBy:'pathComponents.1'`, `nodeOpacity:0.33`,
`enableLabel:true`, `label:'formattedValue'`, `labelSkipSize:0`,
`enableParentLabel:true`, `parentLabel:'id'`, `parentLabelSize:20`,
`parentLabelPosition:'top'`, `parentLabelPadding:6`, `borderWidth:1`. Rects via
`rects`; `hierarchy().sum().round()` + `treemap()` tiling.

**sunburst** — `layers:['arcs','arcLabels']`, `id:'id'`, `value:'value'`,
`cornerRadius:0`, `colors:{scheme:'nivo'}`, `colorBy:'id'`,
`inheritColorFromParent:true`, `childColor:{from:'color'}`, `borderWidth:1`,
`borderColor:'white'`, `enableArcLabels:false`, `arcLabel:'formattedValue'`,
`arcLabelsRadiusOffset:0.5`, `arcLabelsSkipAngle:0`, `arcLabelsSkipRadius:0`,
`transitionMode:'innerRadius'`. `partition().size([2π, r²])`; arcs via existing
`charts/arcs` (`startAngle:x0, endAngle:x1, innerRadius:√y0, outerRadius:√y1`).

**icicle** — `layers:['rects','labels']`, `sort:'input'`, `identity:'id'`,
`value:'value'`, `orientation:'bottom'`, `gapX:1`, `gapY:1`,
`colors:{scheme:'nivo'}`, `colorBy:'id'`, `inheritColorFromParent:true`,
`childColor:{from:'color'}`, `borderRadius:0`, `borderWidth:0`,
`borderColor:{from:'color',modifiers:[['darker',0.6]]}`, `enableLabels:false`,
`label:'id'`, `labelBoxAnchor:'center'`, plus label align/skip/padding/offset/
rotation defaults, `enableZooming:true` (zoom deferred — static first).
`partition().size([w,h])` rectangular per orientation.

**circle-packing** — `layers:['circles','labels']`, `id:'id'`, `value:'value'`,
`padding:0`, `leavesOnly:false`, `colors:{scheme:'nivo'}`, `colorBy:'depth'`,
`inheritColorFromParent:false`,
`childColor:{from:'color',modifiers:[['darker',0.3]]}`, `borderWidth:0`,
`borderColor:{from:'color',modifiers:[['darker',0.3]]}`, `enableLabels:false`,
`label:'id'`, `labelsSkipRadius:8`. `pack().size().padding()` → `x,y,r` circles.

**tree** — `layers:['links','nodes','labels','mesh']`, `identity:'id'`,
`mode:'dendogram'`, `layout:'top-to-bottom'`, `nodeSize:12`,
`nodeColor:{scheme:'nivo'}`, `linkCurve:'bump'`, `linkThickness:1`,
`linkColor:{from:'source.color',modifiers:[['opacity',0.4]]}`,
`enableLabel:true`, `label:'id'`, `labelsPosition:'outward'`, `orientLabel:true`,
`labelOffset:6`, `useMesh:true`, `highlightAncestorNodes:true`,
`highlightAncestorLinks:true`, `nodeTooltipPosition:'fixed'`. `tree()`/`cluster()`
+ `curveBumpX/Y` links; **mesh via §3.2**.

### 4.3 d3-delaunay chart + mesh retrofit

**voronoi** — `xDomain:[0,1]`, `yDomain:[0,1]`,
`layers:['links','cells','points','bounds']`, `enableLinks:false`,
`linkLineWidth:1`, `linkLineColor:'#bbbbbb'`, `enableCells:true`,
`cellLineWidth:2`, `cellLineColor:'#000000'`, `enablePoints:true`, `pointSize:4`,
`pointColor:'#666666'`, `role:'img'`. Direct render of the delaunay/voronoi
structures.

**Mesh retrofit** (the high-leverage payoff): replace v2's direct-point-only
hover in **line, scatterplot, bump, swarmplot, tree** with true voronoi-mesh
hit-testing. The `charts/interact` client layer already does nearest-point
lookup off `data-tc-mesh`; v3 lets the *server* emit an accurate voronoi cell
partition (via §3.2 `Find`) so detection matches nivo exactly, including
`detectionRadius`. Resolves the v2 §11 "voronoi-mesh hover deferred" item.

### 4.4 d3-force charts

**network** — `linkDistance:30`, `centeringStrength:1`, `repulsivity:10`,
`distanceMin:1`, `distanceMax:Infinity`, `iterations:120`, `nodeSize:12`,
`activeNodeSize:18`, `inactiveNodeSize:8`, `nodeColor:'#000000'`,
`nodeBorderWidth:0`, `linkThickness:1`, `layers:['links','nodes','annotations']`.
`Stop()` then `Tick(iterations)`; render final node/link positions.

**swarmplot** — `id:'id'`, `value:'value'`,
`valueScale:{type:'linear',min:0,max:'auto'}`, `groupBy:'group'`, `size:6`,
`spacing:2`, `layout:'vertical'`, `gap:0`, `forceStrength:1`,
`simulationIterations:120`, `colors:{scheme:'nivo'}`, `colorBy:'group'`,
`borderWidth:0`, `layers:['grid','axes','circles','annotations','mesh']`,
`enableGridX:true`, `enableGridY:true`, `useMesh:false`. `ForceX/ForceY` +
`ForceCollide(size/2+spacing/2)`, fixed ticks; mesh via §3.2.

### 4.5 d3-sankey chart

**sankey** — `layout:'horizontal'`, `align:'center'`, `sort:'auto'`,
`colors:{scheme:'nivo'}`, `nodeOpacity:0.75`, `nodeHoverOpacity:1`,
`nodeHoverOthersOpacity:0.15`, `nodeThickness:12`, `nodeSpacing:12`,
`nodeInnerPadding:0`, `nodeBorderWidth:1`, `linkOpacity:0.25`,
`linkHoverOpacity:0.6`, `linkHoverOthersOpacity:0.15`, `linkContract:0`,
`linkBlendMode:'multiply'`, `enableLinkGradient:false`, `enableLabels:true`,
`label:'id'`, `labelPosition:'inside'`, `labelPadding:9`,
`labelOrientation:'horizontal'`, `layers:['links','nodes','labels','legends']`.
Link ribbons via `curveMonotoneX/Y`; optional gradients via `core` defs.

### 4.6 d3-chord chart

**chord** — `padAngle:0`, `innerRadiusRatio:0.9`, `innerRadiusOffset:0`,
`colors:{scheme:'nivo'}`, `arcOpacity:1`, `activeArcOpacity:1`,
`inactiveArcOpacity:0.15`, `arcBorderWidth:1`, `ribbonOpacity:0.5`,
`activeRibbonOpacity:0.85`, `inactiveRibbonOpacity:0.15`, `ribbonBorderWidth:1`,
`ribbonBlendMode:'normal'`, `enableLabel:true`, `label:'id'`, `labelOffset:12`,
`labelRotation:0`, `layers:['ribbons','arcs','labels','legends']`. Arcs reuse
`charts/arcs`; ribbons via §3.5.

### 4.7 d3-geo charts

**geo** (GeoMap + Choropleth) — common: `projectionType:'mercator'`,
`projectionScale:100`, `projectionTranslation:[0.5,0.5]`,
`projectionRotation:[0,0,0]`, `enableGraticule:false`, `graticuleLineWidth:0.5`,
`graticuleLineColor:'#999999'`, `borderWidth:0`, `borderColor:'#000000'`,
`isInteractive:true`, `layers:['graticule','features'(,'legends')]`, `legends:[]`.
GeoMap adds `fillColor:'#dddddd'`; Choropleth adds `match:'id'`, `label:'id'`,
`value:'value'`, `colors:'PuBuGn'`, `unknownColor:'#999'`, `fill:[]`, `defs:[]`,
and `'legends'` in layers. Caller supplies GeoJSON `features`.

## 5. Interactivity, responsive, accessibility

No new model — reuse v2. Specifics:
- **Mesh hover** becomes *accurate* once d3-delaunay lands (§4.3): line,
  scatterplot, bump, swarmplot, tree get voronoi-cell detection with
  `detectionRadius`, all resolved client-side off `data-tc-mesh` (no
  per-mousemove server round-trip).
- **State changes** (sankey node/link hover-others dimming, chord active-arc,
  network active-node, series toggles) stay on HTMX via the existing
  `htmx.Registry` + `Handler`.
- **Responsive/a11y**: thread the existing `Responsive`, `Title`/`Desc`,
  `Role`/`AriaLabel`/`IsFocusable` props through every new chart's props
  (defaults empty/off so static goldens are byte-stable).

## 6. Build / test

- **Golden SVG snapshots** per new chart under `charts/{chart}/testdata/golden/`,
  via `internal/golden.Assert`; `make golden` extended to every new package.
- **Path/layout golden tests** for the new geometry: hierarchy node coordinates
  (treemap tiling, pack circles, partition rects, tidy-tree positions) vs
  d3-hierarchy fixtures; delaunay cell polygons + `Find` results; sankey node
  rects + link widths; chord group angles + ribbon paths; geo projected paths
  for each projection; `curveBumpX/Y` path strings.
- **Force determinism**: golden node positions after `Tick(120)` with the ported
  LCG; if exact-match proves brittle across platforms, fall back to asserting
  bounding box + pairwise non-overlap and **document** the chosen strategy in
  `NOTES.md` (don't silently weaken the golden).
- CI unchanged: `make ci` (`make lint` + `make test`) from root.

## 7. Demo app — `examples/app`

Extend the stdlib `net/http` app with a page per new family, reusing
`demos/` + `handlers/` + `htmx.Registry`:
- `/bump`, `/marimekko`, `/parallel-coordinates`, `/polar-bar`, `/treemap`,
  `/sunburst`, `/icicle`, `/circle-packing`, `/tree`, `/voronoi`, `/network`,
  `/swarmplot`, `/sankey`, `/chord`, `/geo` (GeoMap + Choropleth on one page).
- Nav entry per page; the client interactivity script is already loaded once in
  `layout.templ`. geo page ships a small bundled GeoJSON (world-countries)
  sample so the demo is self-contained.

## 8. Implementation order (topological, by leverage/risk)

1. [x] **Fast-follows (no new d3)** — `curveBumpX/Y` added to `d3/shape`, then
   **marimekko, parallel-coordinates, polar-bar, bump** (bump ships with direct
   hover; mesh layer added in Phase 3). Warms up the v3 chart pipeline with zero
   new-port risk.
2. [x] **`d3/hierarchy` port** → **treemap, sunburst, icicle, circle-packing, tree**
   (tree links use the Phase-1 bump curves; tree/icicle zoom deferred to a
   later pass — static first). Highest leverage: one port, five charts.
3. [x] **`d3/delaunay` port** → **voronoi** chart, then **retrofit voronoi-mesh
   hover** into line, scatterplot, bump, swarmplot(after Phase 4), tree.
   Resolves the v2 §11 mesh-hover deferral.
4. [x] **`d3/force` port** (LCG + fixed ticks) → **network, swarmplot** (swarmplot
   picks up the Phase-3 mesh).
5. [x] **`d3/sankey` port** → **sankey**.
6. [ ] **`d3/chord` port** → **chord**.
7. [ ] **`d3/geo` port** → **geo** (GeoMap + Choropleth). Heaviest, lowest
   chart-count — **may split to v3.1**. Ship minimal-viable (Mercator +
   equirectangular + `GeoPath` + graticule) first, add projections incrementally.
8. [ ] **Polish**: demo pages, `README`/`docs` updates for the full catalog, golden
   regeneration, `make ci` green.

## 9. Scope summary for v3

**New d3 ports**: `hierarchy`, `delaunay`, `force`, `sankey`, `chord`, `geo`;
`d3/shape` gains `curveBumpX/Y`.

**New charts (16 packages / 17 components)**: bump, marimekko,
parallel-coordinates, polar-bar, treemap, sunburst, icicle, circle-packing,
tree, voronoi, network, swarmplot, sankey, chord, geo (GeoMap + Choropleth).

**Retrofit**: true voronoi-mesh hover into line, scatterplot, bump, swarmplot,
tree.

**Deferred within v3** (feature-complete charts, cosmetic/interaction extras
left for a follow-up): icicle/treemap/circle-packing **zoom**, geo full
projection set beyond the minimal-viable subset.

## 10. Post-v3 nivo parity status

After v3, templ-charts covers **every SVG chart type nivo ships**. The
foundation packages (`core`, `theming`, `scales`, `colors`, `axes`, `grid`,
`polar-axes`, `arcs`, `rects`, `text`, `tooltip`, `legends`, `annotations`) are
all exercised. The remaining nivo surface not ported is intentional (§11).

## 11. Explicitly deferred (v4+)

- **Canvas rendering path** — nivo ships `*Canvas` variants (bar, line,
  scatterplot, heatmap, network, swarmplot, geo, voronoi, …) for large-N
  datasets. templ-charts remains **SVG-only**; a Canvas backend (server-side
  raster or a client `<canvas>` draw-list) is the single largest remaining nivo
  capability and the natural v4 theme.
- **TopoJSON decoding** — geo consumes GeoJSON; callers convert upstream.
- **d3-geo full projection catalog** — beyond the ~10 nivo exposes (only if a
  consumer needs more).
- **HSL/Lab color spaces** in `d3/color` — still RGB-only; add only on demand.
- **Interactive zoom** for icicle/treemap/circle-packing/sunburst — layout is
  present; the client-driven zoom transition is a later interaction pass.
