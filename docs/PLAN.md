# templ-charts — Implementation Plan (v1)

## 1. Vision

A Go library exposing nivo-equivalent chart types as **templ components** that generate server-side SVG. v1 delivers a complete framework plus three chart types (**bar**, **line**, **pie**), an HTMX-backed interactivity layer, and a runnable demo app.

## 2. Confirmed decisions

| Concern | Decision |
|---|---|
| Rendering | Server-rendered SVG via templ (no Canvas) |
| Interactivity | HTMX (htmx.org CDN) island pattern, stateful `htmx.Registry` |
| Animation | SMIL `<animate>` + scoped CSS keyframes, gated by `Animate bool` |
| Canvas | Skipped |
| Packages | Mirror nivo names under `charts/` |
| Demo | `examples/app`, stdlib `net/http`, minimal plain CSS, run via `go run ./examples/app` |
| D3 ports | Vendored minimal `internal/d3/` reimplementations; arc-with-cornerRadius ported faithfully as dense math |
| Sample data | Ported from `@nivo/static/samples` |
| Charts in v1 | bar, line, pie |

## 3. Repository layout

```
go.mod                                  github.com/geoffjay/templ-charts
AGENTS.md                               build/test/lint commands for opencode
docs/PLAN.md                            this plan
Makefile                                templ, build, test, lint, vet, run-demo
charts/
  core/                – Dimensions, Margin, Box, BoxAlign, PropertyAccessor,
                         ValueFormat, SvgDefsAndFill, CartesianMarkers,
                         SvgWrapper, Defs (gradients/patterns), DotsItem,
                         curve/stack/blend-mode enums, MotionProps
  theming/             – Theme (background, text, axis, grid, crosshair, legends,
                         labels, markers, dots, tooltip, annotations), DefaultTheme,
                         ExtendDefaultTheme, BorderRadius helpers, bridge attrs,
                         sanitize text styles
  scales/              – ScaleSpec types (linear/log/symlog/point/band/time),
                         ComputedSerieAxis, computeScale, computeXYScalesForSeries,
                         stack, getScaleTicks, centerScale, time helpers
  colors/              – InheritedColorConfig, OrdinalColorScale(+Config),
                         Sequential/Diverging/Quantize, full color scheme catalog,
                         color modifiers
  axes/                – AxisProps, CanvasAxisProps, GridValues, AxisLegendPosition,
                         AxisTickProps, Axes/Axis/AxisTick/Grid/GridLines/GridLine,
                         computeCartesianTicks, computeGridLines, getFormatter,
                         positions, defaultAxisProps
  grid/                – GridCell, generateGrid, computeCellDimensions,
                         BoundingBox/polygon helpers (scaffold-only in v1)
  polar-axes/          – CircularAxis, RadialAxis, PolarGrid (scaffold-only in v1)
  rects/               – Rect, RoundedRect (BuildRoundedRectPath), NodeA11yProps,
                         RectLabels anchor factories; SMIL enter animation support
  arcs/                – Arc, ArcGenerator (createArcGenerator),
                         generateSvgArc path builder, computeArcCenter,
                         computeArcBoundingBox, findArcUnderCursor, arc labels,
                         arc link labels (compute + components)
  text/                – text-anchor/baseline helpers, truncation helpers
  tooltip/             – BasicTooltip, Chip, TableTooltip, Crosshair types,
                         TooltipPosition, TooltipAnchor, CrosshairType
  legends/             – LegendProps, LegendAnchor, LegendDirection, LegendItem-
                         Direction, Datum, SymbolShape factories, BoxLegend,
                         ContinuousColorsLegend, compute helpers, legendDefaults
  annotations/         – AnnotationMatcher, AnnotationSpec (circle/dot/rect),
                         BoundAnnotation, AnnotationInstructions, bindAnnotations,
                         computeAnnotation, getLinkAngle
  static/              – ChartsMapping registry, RenderChart dispatcher,
                         defaults, samples (bar/line/pie)
  bar/                 – BarDatum, ComputedDatum, ComputedBarDatum, BarProps,
                         useBar equivalent, compute/{common,grouped,stacked,
                         legends,totals}, BarItem + BarTotals + BarAnnotations +
                         BarLegends templ comps; default layers: grid→axes→bars→
                         totals→markers→legends→annotations
  line/                – LineSeries, ComputedSeries, Point, SliceData,
                         useLine equivalent, line/area path generators, usePoints/
                         useSlices, Lines/LinesItem/Areas/Points/Slices/Mesh
                         templ comps; default layers: grid→markers→axes→areas→
                         crosshair→lines→points→slices→mesh→legends
  pie/                 – DefaultRawDatum, ComputedDatum, PieArc, PieProps,
                         normalizeData, pieArcs, pieFromBox (d3-pie/arc), Arcs,
                         ArcLinkLabels, ArcLabels, PieLegends, PieTooltip;
                         default layers: arcs→arcLinkLabels→arcLabels→legends
  htmx/                – ChartInstance, Registry, Handler (net/http): full render
                         + /hover + /click + /toggle + /slice endpoints
internal/
  d3/                  – Pure-Go ports of d3-shape (line, area, curve, stack,
                         pie, arc incl. cornerRadius), d3-scale (linear, log,
                         symlog, point, band, time), d3-array, d3-format,
                         d3-time-format, d3-color
examples/
  app/                 – demo app (see §6)
templates/             – shared templ fragments
```

## 4. Foundation packages — what each contains

### 4.1 `charts/core`
- Types: `Dimensions`, `Margin`, `Box`, `BoxAlign`, `Point`, `Padding`, `DatumValue`
- Accessors: `PropertyAccessor[D,V]` (string path or func), `GetPropertyAccessor`, `GetLabelGenerator`
- Formatters: `ValueFormat[V,C]`, `GetValueFormatter` (parses `time:`-prefixed, falls back to `fmt.Sprintf`)
- SVG defs: `Def`, `LinearGradientDef`, `PatternDotsDef`, `PatternLinesDef`, `PatternSquaresDef`, `SvgDefsAndFill`, `BindDefs` (with inherit color support)
- Templ components: `SvgWrapper` (svg + bg rect + `<g translate(margin)>`), `DotsItem`, `Defs`, `CartesianMarkers`, `CartesianMarkersItem`
- Enums: `CurveFactoryId`, `StackOrder`, `StackOffset`, `CssMixBlendMode` (consts + maps)
- `MotionProps` (Animate bool + config), `DefaultAnimate = true`

### 4.2 `charts/theming`
- `Theme` struct mirroring `@nivo/theming` defaults verbatim
- `TextStyle`, `PartialTheme`, `ThemeWithoutInheritance`
- `ExtendDefaultTheme(default, custom)` (deep merge + 9-path text inheritance)
- `ExtendAxisTheme`, `NormalizeBorderRadius`, `ConstrainBorderRadius`, `BorderRadiusToCss`
- `Engine` bridge tables (`svgStyleAttributesMapping`, `convertStyleAttribute`), `SanitizeSvgTextStyle`, `SanitizeHtmlTextStyle`

### 4.3 `charts/scales`
- `ScaleSpec` interface + concrete types: `ScaleLinearSpec`, `ScaleLogSpec`, `ScaleSymlogSpec`, `ScalePointSpec`, `ScaleBandSpec`, `ScaleTimeSpec`
- Concrete scales implementing d3 operations (`Domain`, `Range`, `RangeRound`, `Call`, `Bandwidth`, `Step`, `Nice`, `Clamp`, `Copy`)
- `ComputeScale`, `ComputeXYScalesForSeries`, `ComputedSerieAxis`, `StackAxis`/`StackX`/`StackY`
- `GetScaleTicks`, `CenterScale`, `TicksSpec[V]`, `TIME_PRECISION`, `CreateDateNormalizer`
- Backed by `internal/d3/scale` and `internal/d3/array`

### 4.4 `charts/colors`
- `InheritedColorConfig[D]` + `GetInheritedColor(config, theme, datum)`
- `ColorModifier` (brighter/darker/opacity) backed by `internal/d3/color`
- `OrdinalColorScaleConfig[D]` + `GetOrdinalColorScale`
- `Sequential/Diverging/QuantizeColorScale`
- Full color scheme catalog: categorical (nivo, category10, accent, dark2, paired, pastel1/2, set1/2/3), cyclical (rainbow, sinebow), diverging (BrBG, PiYG, PRGn, PuOr, RdBu, RdGy, RdYlBu, RdYlGn, spectral, 3–11 each), sequential single+multi-hue (Blues, Greens, Greys, Oranges, Purples, Reds, BuGn, BuPu, GnBu, OrRd, PuBu, PuBuGn, PuRd, RdPu, YlGn, YlGnBu, YlOrBr, YlOrRd)
- Interpolators (viridis, inferno, magma, plasma, warm, cool, cubehelixDefault)

### 4.5 `charts/axes`
- `AxisProps`, `CanvasAxisProps`, `AxisTickProps`, `AxisLegendPosition`, `GridValues[V]`, `TicksSpec`
- `ComputeCartesianTicks` → `{Ticks, TextAlign, TextBaseline}`
- `ComputeGridLines` → `[]Line`, `GetFormatter`, `Positions` const, `DefaultAxisProps`
- Templ components: `Axes`, `Axis`, `AxisTick`, `Grid`, `GridLines`, `GridLine`

### 4.6 `charts/grid` (scaffold-only in v1)
- `GridCell`, `BoundingBox`, `Vertex`, `GridFillDirection`
- `GenerateGrid`, `ComputeCellDimensions`, `AreBoundingBoxTouching`, `PerpendicularPolygon`, `GetCellsPolygons`
- Powers heatmap/waffle — not exercised in v1.

### 4.7 `charts/polar-axes` (scaffold-only in v1)
- Types: `CircularAxisConfig`, `RadialAxisConfig`, tick props
- Templ components: `CircularAxis`, `RadialAxis`, `PolarGrid`, `RadialGrid`, `CircularGrid`
- Not used by pie (pie does its own layout); grounds future radar/chord/sunburst.

### 4.8 `charts/rects`
- `Rect`, `NodeA11yProps`, `NodeWithRect`, `BorderRadiusCorners`, `BorderRadius`
- **`BuildRoundedRectPath(x,y,w,h,tl,tr,br,bl) string`** — primary rect path helper (used by bar)
- `RoundedRect` templ component (emits SMIL `<animate>` on width/height when entering + Animate)
- `RectLabels` anchor factories (9 anchors)

### 4.9 `charts/arcs` — fully implemented in v1
- `Arc{startAngle,endAngle,innerRadius,outerRadius}` (radians), `DatumWithArc`, `ArcGenerator`
- `CreateArcGenerator(cornerRadius, padAngle)` — wraps `internal/d3/shape.Arc()`
- `ComputeArcBoundingBox(centerX, centerY, radius, startAngleDeg, endAngleDeg, includeCenter)` — samples endpoints + axis crossings (multiples of 90°)
- `ComputeArcCenter(arc, offset) Point` — `angle = midAngle - π/2`
- `GetNormalizedAngle`, `FilterDataBySkipAngle`, `GenerateSvgArc`
- `ComputeArcLink(arc, offset, diagLength, straightLength) ArcLink` + `ComputeArcLinkTextAnchor(arc)`
- Templ components: `ArcShape`, `ArcsLayer`, `ArcLabel`, `ArcLabelsLayer`, `ArcLinkLabel`, `ArcLinkLabelsLayer`
- **Dropped (animation-only)**: useAnimatedArc, useArcsTransition, useArcTransitionMode, interpolateArc (replaced by direct `arcGenerator(arc)`), useArcLinkLabelsTransition, useArcCentersTransition, interpolateArcCenter, canvas files, ArcLine.

### 4.10 `charts/text`
- Shared text primitives: `SvgTextAttrs(textAlign, textBaseline) → {textAnchor, dominantBaseline}`
- Truncation helper (mirrors `truncateTickAt`)

### 4.11 `charts/tooltip`
- Types: `TooltipPosition` (cursor|fixed), `TooltipAnchor` (top|right|bottom|left|center), `CrosshairType` (12 variants)
- HTML templ components (float in container div): `BasicTooltip`, `Chip`, `TableTooltip`
- SVG `Crosshair` component (`theme.crosshair.line` styling)
- Plumbing routed through `charts/htmx` (HTML fragments swapped in by HTMX)

### 4.12 `charts/legends`
- Types: `LegendAnchor`, `LegendDirection`, `LegendItemDirection`, `Datum`, `LegendProps`, `BoxLegendSvgProps`, `LegendSvgItemProps`, `SymbolShape`, `SymbolProps`, `ContinuousColorsLegendProps`
- Compute: `ComputeDimensions`, `ComputePositionFromAnchor`, `ComputeItemLayout`, `ComputeContinuousColorsLegend`
- Templ components: `BoxLegendSvg`, `LegendSvg`, `LegendSvgItem`, `ContinuousColorsLegendSvg`, symbol factories (circle/diamond/square/triangle)
- `LegendDefaults`, `ContinuousColorsLegendDefaults`
- Legend toggle series routed through `charts/htmx`

### 4.13 `charts/annotations`
- Types: `RelativeOrAbsolutePosition`, `AnnotationMatcher`, `CircleAnnotationSpec`, `DotAnnotationSpec`, `RectAnnotationSpec`, `BoundAnnotation`, `AnnotationInstructions`
- Compute: `BindAnnotations(data, matchers, getPosition, getDimensions) []BoundAnnotation`, `ComputeAnnotation(bound) AnnotationInstructions`, `GetLinkAngle`
- Templ components: `Annotation` (link + outline + symbol + note), `CircleAnnotationOutline`, `DotAnnotationOutline`, `RectAnnotationOutline`
- Defaults: `noteWidth=120`, `noteTextOffset=8`, `dotSize=4`, rect `borderRadius=6`

### 4.14 `charts/static`
- `ChartType` enum ("bar", "line", "pie")
- `ChartComponent` interface `{ Render(props, override) (string, error) }`
- `ChartsMapping` map[ChartType]Mapping{`Component`, `RuntimeProps`, `Defaults`}
- `RenderChart(type, props, override) (string, error)` — applies static defaults `{animate:false, isInteractive:false, renderWrapper:false, theme:{}}` + chart defaults + user props + whitelisted overrides, renders to SVG string
- `Samples` map of demo data for bar/line/pie

## 5. Chart packages — bar, line, pie

### 5.1 `charts/bar`

**Types**: `BarDatum`, `ComputedDatum[D]`, `ComputedBarDatum[D]`, `BarLegendProps` (adds `DataFrom: indexes|keys`), `BarLayerId`, `BarProps[D]`/`BarSvgProps[D]`

**Compute** (`bar/compute/`):
- `common.go`: `GetIndexScale`, `NormalizeData`, `FilterNullValues`, `CoerceValue`, `ComputeLabelLayout`
- `grouped.go`: `GenerateGroupedBars` (vertical/horizontal)
- `stacked.go`: `GenerateStackedBars` using `d3-shape.stack.offset(diverging)`
- `legends.go`: `GetLegendData` with order-reversal rules
- `totals.go`: `ComputeBarTotals`

**`UseBar(props)` orchestrator**: color scales (`GetOrdinalColorScale` + `GetInheritedColor` for border/label), accessors, value formatter, generates bars, legendData, legendsWithData, barTotals, shouldRenderLabel.

**Templ components**: `Bar` (pipeline: dimensions → UseBar → label-layout → bound defs → layerById → SvgWrapper), `BarItem`, `BarTotals`, `BarAnnotations`, `BarLegends`.

**Defaults** (verbatim from `defaults.ts`): `groupMode:"stacked"`, `layout:"vertical"`, `valueScale:{type:linear,nice:true,round:false}`, `indexScale:{type:band,round:false}`, `padding:0.1`, `innerPadding:0`, `colors:{scheme:nivo}`, `colorBy:id`, `borderRadius:0`, `borderWidth:0`, `borderColor:{from:color}`, `enableLabel:true`, `labelPosition:middle`, `labelTextColor:{theme:labels.text.fill}`, layers `['grid','axes','bars','totals','markers','legends','annotations']`.

**Animation**: under `Animate:true`, each rect grows from 0 height (vertical) or 0 width (horizontal) via SMIL `<animate>` (`begin="0s" dur="0.6s" fill="freeze"`). Mirrors nivo's `enter` state.

### 5.2 `charts/line`

**Types**: `LineSeries`, `ComputedDatum[Series]`, `ComputedSeries[Series]`, `Point[Series]`, `SliceData[Series]`, `LineLayerId`, `CommonLineProps`, `LineSvgExtraProps`, `LineSvgProps`.

**Hooks/compute** (`line/hooks.go`):
- `UseLineGenerator(curve)` → `func(points) string` (d3-shape line port with `.defined` honoring nulls)
- `UseAreaGenerator(curve, yScale, areaBaselineValue)` → d3-area port
- `UsePoints(...)` → `[]Point`
- `UseSlices(componentId, enableSlices, points, w, h)` → `[]SliceData`
- `UseLine(...)` — full orchestrator: `ComputeXYScalesForSeries`, color by series id, legendData, points, slices, generators

**Templ components**: `Line` (pipeline), `Lines`, `LinesItem` (SMIL `<animate>` on `d` when Animate), `Areas`, `Points`, `Slices`, `SlicesItem`, `Mesh`, `PointTooltip`, `SliceTooltip`.

**Defaults** (verbatim): `xScale:{type:point}`, `yScale:{type:linear,min:0,max:auto}`, `curve:linear`, `lineWidth:2`, `colors:{scheme:nivo}`, `pointSize:6`, `enableArea:false`, `areaOpacity:0.2`, `enableGridX/Y:true`, `crosshairType:bottom-left`, `useMesh:false`, `enableSlices:false`, `motionConfig:gentle`, layers `['grid','markers','axes','areas','crosshair','lines','points','slices','mesh','legends']`.

**Interaction**: three modes (direct point, slices, mesh) routed through `charts/htmx`; crosshair rendered server-side based on current state.

### 5.3 `charts/pie`

**Types** (`pie/types.go`):
- `DefaultRawDatum{ID, Value}`, `MayHaveLabel{Label}`
- `PieArc` (extends `arcs.Arc`, adds `Index`, `Angle`, `AngleDeg`, `Thickness`, `PadAngle`)
- `ComputedDatum[RawDatum]{ID, Label, Value, FormattedValue, Color, Fill?, Data, Arc, Hidden}`
- `LegendDatum[RawDatum]{ID, Label, Color, Hidden, Data}`
- `PieLayerId`: `arcs | arcLinkLabels | arcLabels | legends`
- `CommonPieProps[RawDatum]` + `PieSvgProps[RawDatum]`

**Defaults** (`pie/props.go`) — verbatim from `pie/src/props.ts`:
```
Layers: [arcs, arcLinkLabels, arcLabels, legends]
ID: "id", Value: "value", SortByValue: false
StartAngle: 0, EndAngle: 360, PadAngle: 0
Fit: true, InnerRadius: 0, CornerRadius: 0
ActiveInnerRadiusOffset: 0, ActiveOuterRadiusOffset: 0
Colors: {scheme:"nivo"}, BorderWidth: 0, BorderColor: {from:"color", modifiers:[["darker",1]]}
EnableArcLabels: true, ArcLabel: "formattedValue", ArcLabelsRadiusOffset: 0.5
ArcLabelsSkipAngle: 0, ArcLabelsSkipRadius: 0
ArcLabelsTextColor: {theme:"labels.text.fill"}
EnableArcLinkLabels: true, ArcLinkLabel: "id", ArcLinkLabelsSkipAngle: 0
ArcLinkLabelsOffset: 0, DiagonalLength: 16, StraightLength: 24, Thickness: 1, TextOffset: 6
ArcLinkLabelsTextColor: {theme:"labels.text.fill"}
ArcLinkLabelsColor: {theme:"axis.ticks.line.stroke"}
IsInteractive: true, Tooltip: PieTooltip, TransitionMode: "innerRadius"
Legends: [], Role: "img"
Animate: true, MotionConfig: "gentle"
```

**Hooks/compute** (`pie/hooks.go`) — pure functions:
- `NormalizeData(data, idAcc, valueAcc, valueFmt, getColor) []NormalizedDatum`
- `PieArcs(data, startAngleDeg, endAngleDeg, padAngleDeg, sortByValue, innerPx, outerPx, activeOffsets, activeID) ([]ComputedDatum, []LegendDatum)`
- `PieFromBox(data, width, height, fit, innerRadiusRatio, startAngleDeg, endAngleDeg, padAngleDeg, sortByValue, cornerRadius, activeOffsets, activeID, forwardLegendData) PieResult`:
  - `radius = min(w,h)/2`, `innerRadius = radius * ratio`, `centerX/Y = w/2, h/2`
  - If `fit`: `arcs.ComputeArcBoundingBox(centerX, centerY, radius, startAngle-90, endAngle-90)` → ratio → recenter/rescale
  - Run `d3.Pie().value().startAngle(rad).endAngle(rad).padAngle(rad)` (`.sortValues(null)` unless sortByValue) over non-hidden data
  - Build `arcGenerator = d3.Arc().innerRadius(fn).outerRadius(fn).cornerRadius().padAngle(rad)`
- `UsePieLayerContext(...)`

**Templ components**:
- `Pie` — full pipeline: `useDimensions` → `NormalizeData` → `PieFromBox` → `bindDefs` → `layerById` → `<SvgWrapper>` iterating layers
- `Arcs` (local wrapper) — `<g translate(center)>` over each `<path d=arcGenerator(arc) fill=color stroke=borderColor stroke-width=borderWidth>`; HTMX hover attrs when interactive
- `ArcLinkLabels` — `<g translate(center)>`: `arcs.ComputeArcLink` per datum → `<path d="M…L…L…">` + `<text>` at `p2 ± textOffset`
- `ArcLabels` — `<g translate(center)>`: `arcs.ComputeArcCenter(arc, radiusOffset)` → `<g translate(cx,cy)><text text-anchor=middle dominant-baseline=central>`
- `PieLegends` — `<BoxLegendSvg>` per legend config
- `PieTooltip` — `<BasicTooltip id=formattedValue enableChip color>`

**Default layers**: `['arcs', 'arcLinkLabels', 'arcLabels', 'legends']`

**Animation**: under `Animate`, arcs use `transitionMode: 'innerRadius'` → SMIL `<animate attributeName="d">` growing outer radius from `innerRadius` to `outerRadius` over 600ms ease. (Mirrors nivo's "enter" phase.)

**Angle convention**: User-facing `startAngle`/`endAngle`/`padAngle` in **degrees** (0 = top, clockwise). Converted to radians before d3-shape. nivo compensates with `-90` / `-π/2` offsets in `ComputeArcBoundingBox` and `ComputeArcCenter`. Default `0/360` = full pie; `0/180` = half-pie (rescaled by `fit`).

## 6. `charts/htmx` — interactivity layer

Two modes:

**Static (default)**: `bar.Render(props)` returns SVG string; no HTMX attrs emitted. `Props.Interactive=false`. Embed in any non-interactive page.

**Interactive**: `htmx.NewRegistry()` holds `map[InstanceID]*ChartInstance`. Each instance knows its kind (bar/line/pie), full `Props` snapshot, and mutable state (`HiddenIDs`, `HoveredKey`, `ActiveID`).

Endpoints:

| Route | Method | Action |
|---|---|---|
| `/charts/{id}` | GET | Full SVG render with `hx-*`-aware markup |
| `/charts/{id}/hover?bar=key` | GET | Returns tooltip HTML partial (`BasicTooltip`/`TableTooltip`) |
| `/charts/{id}/click?bar=key&verb=activate` | POST | Toggle activation server-side |
| `/charts/{id}/toggle?series=id` | POST | Re-render full SVG with series toggled |
| `/charts/{id}/slice?axis=x&x=value` | GET | Returns crosshair + slice tooltip (line) |
| `/charts/{id}/hover?arc=id` | GET | Returns pie tooltip + sets `ActiveID` (radius offset) |

`htmx.Handler` is a standard `http.Handler`; users register chart instances via `registry.Register("my-bar", barProps)`. Templates emit `hx-get`/`hx-trigger="mouseenter"`/`hx-target`/`hx-swap="innerHTML"`.

**Documented trade-off**: stateful server registry is fine for demos/small apps; horizontal-scale deployments would persist state in a session/cookie instead.

## 7. Demo app — `examples/app`

Stdlib `net/http` server, **minimal plain CSS** (~50–100 lines, hand-written; system fonts, subtle borders, card grid), zero external deps beyond htmx.org CDN `<script>`.

Pages:
- `/` — index listing all demos with thumbnails
- `/bar` — Bar demos: stacked vertical, grouped horizontal, markers + annotations, legend + toggle (HTMX), totals layer
- `/line` — Line demos: single series, multiple series + legend, area + points, slices + slice tooltip, mesh (voronoi hover)
- `/pie` — Pie demos: plain full pie, donut (`innerRadius`), half pie (`startAngle:0, endAngle:180`) demonstrating `fit`, sorted-by-value, arc labels + arc link labels, active-arc highlight via HTMX hover, legend + series toggle
- `/themes` — bar + line + pie under dark/light/custom themes

Structure:
```
examples/app/
  main.go                  – http.ServeMux wiring, registers instances
  demos/bar.go             – bar demo props + data
  demos/line.go            – line demo props + data
  demos/pie.go             – pie demo props + data
  handlers/chart.go        – /bar, /line, /pie, /themes handlers
  templates/layout.templ   – HTML shell + htmx.org CDN + CSS
  templates/chart.templ    – chart mount + tooltip container
  go.mod?                  – replace directive: github.com/geoffjay/templ-charts => ../../
  README.md                – how to run
```

Run: `go run ./examples/app` → http://localhost:8080

Added to root `go.mod`: `github.com/a-h/templ`. HTMX via CDN (no Go dep).

## 8. `internal/d3` — ports

- **`d3.Shape`**: `Pie()` (.value/.startAngle/.endAngle/.padAngle/.sortValues), `Arc()` (.innerRadius fn/.outerRadius fn/.cornerRadius/.padAngle → SVG path string — dense math port of d3-shape arc including cornerRadius lineTo/arc clamping; golden-tested against d3 outputs), `Line()` (.x/.y/.defined/.curve), `Area()` (.x/.y0/.y1/.defined/.curve), curves (linear, linearClosed, monotoneX/Y, basis, cardinal, catmullRom, step/Before/After), `Stack()` + diverging offset
- **`d3.Scale`**: linear, log, symlog, point, band, time (with `domain`, `range`, `rangeRound`, `nice`, `clamp`, `copy`, `bandwidth`, `step`, ticks)
- **`d3.Array`**: extent, ticks, bisect, range, etc.
- **`d3.Format`**: basic spec parsing
- **`d3.TimeFormat`**: basic spec parsing
- **`d3.Color`**: rgb/hsl conversions backing color modifiers

The arc-with-corner-radius port is the riskiest piece; plan is a near-faithful transliteration of d3-shape's `arc()` (~250 lines of dense math), validated by golden path-string tests against real d3-shape output.

## 9. Scope summary for v1

**Fully implemented** (exercised by bar/line/pie): `core`, `theming`, `scales`, `colors`, `axes`, `rects`, `text`, `tooltip`, `legends`, `annotations`, `static`, `arcs`, `bar`, `line`, `pie`, `htmx`, demo app

**Scaffold-only** (not exercised in v1): `grid`, `polar-axes`

**Skipped**: `canvas`

## 10. Build / test

- `Makefile` targets: `templ`, `build`, `test`, `lint`, `vet`, `run-demo`
- `AGENTS.md` documenting `make lint`, `make test`, `go generate ./...` (templ)
- `go vet` on all packages
- **Unit tests** per foundational package: scales (each scale type vs known d3 outputs), colors (schemes), theming (ExtendDefaultTheme), axes (ComputeCartesianTicks geometry), arcs (ComputeArcBoundingBox, ComputeArcCenter, ComputeArcLink)
- **Golden SVG tests** for `Bar`, `Line`, `Pie` (snapshot → on diff, regenerate with `-update`)
- **Path-string tests** for `BuildRoundedRectPath`, `GenerateSvgArc`, `d3.Shape.Arc` (golden vs d3-shape), line/area generators
- CI-ready: `go test ./...` from root

## 11. Implementation order (high-level, topological)

1. `go.mod` deps (`github.com/a-h/templ`), `Makefile`, `AGENTS.md`, `docs/PLAN.md`
2. `internal/d3/` — array, format, time-format, color; then scale (linear → band → time → log/symlog/point); then shape (line/area/curve → stack → pie → arc)
3. `charts/core` → `charts/theming` → `charts/colors` → `charts/scales`
4. `charts/text` → `charts/axes` → `charts/rects` → `charts/arcs` → `charts/tooltip` → `charts/legends` → `charts/annotations`
5. `charts/grid` + `charts/polar-axes` (scaffolds)
6. `charts/bar` (compute + components + tests)
7. `charts/line` (compute + components + tests)
8. `charts/pie` (compute + components + tests; depends on `arcs`)
9. `charts/static` (registry + samples + RenderChart)
10. `charts/htmx` (registry + handler + endpoints)
11. `examples/app` (handlers + templates + demos for bar/line/pie/themes)
12. Polish: golden tests, `make` targets, READMEs, AGENTS.md updates