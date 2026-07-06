# Chart family reference

Exact constructors, data shapes, key props, and sample-data functions for all
28 chart families. Import path pattern:
`github.com/geoffjay/templ-charts/charts/<package>`.

Every family also has a runnable `ExampleXxx` in its package
(`charts/<package>/example_test.go`) — the authoritative minimal usage.

## bar

- `bar.Bar(bar.BarProps)` — `Data []bar.BarDatum` where `BarDatum = map[string]any`.
- Key props: `IndexBy` (row category accessor, e.g. `"country"`), `Keys []string`
  (columns to stack/group), `GroupMode` (`"stacked"`|`"grouped"`),
  `Layout` (`"vertical"`|`"horizontal"`).
- Sample: `samples.Bar() ([]bar.BarDatum, []string)`.

## line

- `line.Line(line.LineProps)` — `Data []line.LineSeries`:
  `LineSeries{ID string; Data []LinePointData}`, `LinePointData{X, Y any}`
  (numbers, strings, or time values depending on scale).
- Key props: `Curve core.CurveFactoryId` (e.g. `core.CurveMonotoneX`),
  `EnableArea bool`, `EnablePoints bool`.
- Sample: `samples.Line() []line.LineSeries`.

## pie

- `pie.Pie(pie.PieProps)` — `Data []any` with `ID` and `Value` accessors
  (string path or func).
- Key props: `InnerRadius float64` (0 = pie, >0 = donut), `PadAngle float64`.
- Sample: `samples.Pie() []any`.

## boxplot

- `boxplot.BoxPlot(boxplot.BoxPlotProps)` — `Data []boxplot.BoxPlotDatum`:
  `{Group, SubGroup string; Value float64}` (raw observations; quantiles are
  computed).
- Key props: `Layout`, `ColorBy` (`"subGroup"` default | `"group"`),
  `Quantiles []float64`.
- Sample: `samples.BoxPlot() []boxplot.BoxPlotDatum`.

## bullet

- `bullet.Bullet(bullet.BulletProps)` — `Data []bullet.BulletItemDatum`:
  `{ID, Title string; Ranges, Measures, Markers []float64}`.
- Key props: `Layout`, `RangeColors`/`MeasureColors`/`MarkerColors string`
  (sequential scheme ids, e.g. `"seq:cool"`).
- Sample: `samples.Bullet() []bullet.BulletItemDatum`.

## bump

- `bump.Bump(bump.BumpProps)` — `Data []bump.BumpSerie`:
  `{ID string; Data []BumpDatum}`, `BumpDatum{X, Y any}` (Y = rank).
- Key props: `Interpolation` (`"smooth"`|`"linear"`), `PointSize float64`.
- Sample: `samples.Bump() []bump.BumpSerie`.

## calendar

- `calendar.Calendar(calendar.CalendarProps)` — `Data []calendar.CalendarDatum`:
  `{Day string; Value float64}` with `Day` as `"YYYY-MM-DD"`.
- Key props: `From, To time.Time` (auto-derived if zero), `Direction`,
  `Colors []string` (quantize buckets) + `EmptyColor string`.
- Sample: `samples.Calendar() (time.Time, time.Time, []calendar.CalendarDatum)`.

## chord

- `chord.Chord(chord.ChordProps)` — `Data [][]float64` square flow matrix +
  `Keys []string` (Data[i][j] = flow from Keys[i] to Keys[j]).
- Key props: `PadAngle float64` (radians), `InnerRadiusRatio float64`.
- Sample: `samples.Chord() ([][]float64, []string)`.

## circlepacking

- `circlepacking.CirclePacking(circlepacking.CirclePackingProps)` —
  `Data circlepacking.CirclePackingNode` (recursive:
  `{ID string; Value float64; Children []CirclePackingNode}`).
- Key props: `Id`, `Value` accessors; colors default by depth.
- Sample: `samples.CirclePacking() circlepacking.CirclePackingNode`.

## funnel

- `funnel.Funnel(funnel.FunnelProps)` — `Data []funnel.FunnelDatum`:
  `{ID string; Value float64; Label string}` (ordered; width ∝ value).
- Key props: `Direction`, `Interpolation` (`"smooth"`|`"linear"`).
- Sample: `samples.Funnel() []funnel.FunnelDatum`.

## heatmap

- `heatmap.HeatMap(heatmap.HeatMapProps)` — `Data []heatmap.HeatMapSerie`:
  `{ID string; Data []HeatMapDatum}`, `HeatMapDatum{X string; Y *float64}`
  (nil Y = empty cell).
- Key props: `Colors heatmap.HeatMapColorConfig` (continuous
  sequential/diverging scale — not the ordinal config), `ForceSquare bool`,
  `EmptyColor string`. Supports the Canvas backend.
- Sample: `samples.Heatmap() []heatmap.HeatMapSerie`.

## icicle

- `icicle.Icicle(icicle.IcicleProps)` — `Data icicle.IcicleNode` (recursive
  `{ID, Value, Children}`).
- Key props: `Orientation` (`"horizontal"`|`"vertical"`), `EnableZooming bool`
  (HTMX; see interactivity skill).
- Sample: `samples.Icicle() icicle.IcicleNode`.

## marimekko

- `marimekko.Marimekko(marimekko.MarimekkoProps)` — `Data []marimekko.MarimekkoDatum`:
  `{ID string; Value float64; Dimensions map[string]float64}` plus
  `Dimensions []marimekko.MarimekkoDimension` (`{ID, Key string}`).
- Key props: `Layout`, `Offset` (`"none"`|`"expand"`|`"diverging"`|`"silhouette"`|`"wiggle"`).
- Sample: `samples.Marimekko() ([]marimekko.MarimekkoDatum, []marimekko.MarimekkoDimension)`.

## network

- `network.Network(network.NetworkProps)` — `Nodes []network.NetworkInputNode`
  (`{ID string; Size float64; Color string}`) + `Links []network.NetworkInputLink`
  (`{Source, Target string; Distance float64}`). Force-directed layout.
- Key props: `LinkDistance`, `Repulsivity float64`, `Iterations int`.
- Sample: `samples.Network() ([]network.NetworkInputNode, []network.NetworkInputLink)`.

## parallelcoordinates

- `parallelcoordinates.ParallelCoordinates(parallelcoordinates.PCProps)` —
  `Data []PCDatum` (values per variable) + `Variables []PCVariable`
  (`{ID string; Type string; Min, Max float64}`).
- Sample: `samples.ParallelCoordinates() ([]pc.PCDatum, []pc.PCVariable)`.

## polarbar

- `polarbar.PolarBar(polarbar.PolarBarProps)` — `Data []polarbar.PolarBarDatum`
  (`{ID string; Value float64}`) + `Keys []string` (radial stack categories).
- Key props: `StartAngle`, `EndAngle float64` (degrees).
- Sample: `samples.PolarBar() ([]polarbar.PolarBarDatum, []string)`.

## radar

- `radar.Radar(radar.RadarProps)` — `Data []map[string]any` rows +
  `Keys []string` (one polygon layer per key).
- Key props: `IndexBy` accessor (axis label per row), `ValueFormat string`
  (d3-format spec).
- Sample: `samples.Radar() ([]map[string]any, []string)`.

## radialbar

- `radialbar.RadialBar(radialbar.RadialBarProps)` — `Data []radialbar.RadialBarSerie`:
  `{ID string; Data []RadialBarDatum}`, `RadialBarDatum{Category string; Value float64}`.
- Key props: `StartAngle`, `EndAngle float64`.
- Sample: `samples.RadialBar() []radialbar.RadialBarSerie`.

## sankey

- `sankey.Sankey(sankey.SankeyProps)` — `Nodes []sankey.SankeyInputNode`
  (`{ID string}`) + `Links []sankey.SankeyInputLink`
  (`{Source, Target string; Value float64}`).
- Key props: `Layout` (`"horizontal"`|`"vertical"`), `Align`
  (`"center"`|`"justify"`|`"start"`|`"end"`), `NodePadding float64`.
- Sample: `samples.Sankey() ([]sankey.SankeyInputNode, []sankey.SankeyInputLink)`.

## scatterplot

- `scatterplot.ScatterPlot(scatterplot.ScatterPlotProps)` —
  `Data []scatterplot.ScatterPlotSerie`: `{ID string; Data []ScatterPlotDatum}`,
  `ScatterPlotDatum{X, Y any}`.
- Key props: `XScale scales.ScaleSpec`, `YScale scales.ScaleLinearSpec`,
  `NodeSize float64`. Supports the Canvas backend for large N.
- Sample: `samples.ScatterPlot() []scatterplot.ScatterPlotSerie`.

## stream

- `stream.Stream(stream.StreamProps)` — `Data []stream.StreamDatum`
  (`StreamDatum = map[string]float64`) + `Keys []string` (layers).
- Key props: `Offset` (`"none"`|`"wiggle"`|`"silhouette"`|`"expand"`|`"diverging"`),
  `Curve core.CurveFactoryId`.
- Sample: `samples.Stream() ([]stream.StreamDatum, []string)`.

## sunburst

- `sunburst.Sunburst(sunburst.SunburstProps)` — `Data sunburst.SunburstNode`
  (recursive `{ID, Value, Children}`).
- Key props: `StartAngle`, `EndAngle float64`, `EnableZooming bool`.
- Sample: `samples.Sunburst() sunburst.SunburstNode`.

## swarmplot

- `swarmplot.SwarmPlot(swarmplot.SwarmPlotProps)` — `Data []swarmplot.SwarmPlotDatum`
  (`{ID, Group string; Value float64}`) + `Groups []string` (band order).
- Key props: `Size float64` (dot radius).
- Sample: `samples.SwarmPlot() ([]swarmplot.SwarmPlotDatum, []string)`.

## tree

- `tree.Tree(tree.TreeProps)` — `Data tree.TreeNode` (recursive
  `{ID string; Children []TreeNode}` — structure only, no values).
- Key props: `Mode` (`"tree"` tidy | `"dendrogram"` cluster),
  `NodeSize float64`, `LinkCurve core.CurveFactoryId`.
- Sample: `samples.Tree() tree.TreeNode`.

## treemap

- `treemap.Treemap(treemap.TreemapProps)` — `Data treemap.TreemapNode`
  (recursive `{ID, Value, Children}`).
- Key props: `Tile` (`"squarify"` default | `"slice"`|`"dice"`|`"sliceDice"`),
  `EnableZooming bool`.
- Sample: `samples.Treemap() treemap.TreemapNode`.

## voronoi

- `voronoi.Voronoi(voronoi.VoronoiProps)` — `Data []voronoi.VoronoiDatum`:
  `{ID string; X, Y float64}`.
- Key props: `XScale`/`YScale scales.ScaleLinearSpec`, `EnableLinks bool`
  (Delaunay edges), `EnableCells bool` (Voronoi polygons).
- Sample: `samples.Voronoi() []voronoi.VoronoiDatum`.

## waffle

- `waffle.Waffle(waffle.WaffleProps)` — `Data []waffle.WaffleDatum`:
  `{ID string; Value float64; Color string}`.
- Key props: `Total float64`, `Rows int`, `Columns int`.
- Sample: `samples.Waffle() []waffle.WaffleDatum`.

## geo (GeoMap + Choropleth)

- `geo.GeoMap(geo.GeoMapProps)` — `Features []geo.Feature` (GeoJSON),
  `FillColor string`, `ProjectionType string` (`"mercator"`, …),
  `ProjectionScale float64`, `Fit bool`.
- `geo.Choropleth(geo.ChoroplethProps)` — `Features []geo.Feature` +
  `Data []geo.ChoroplethDatum` (`{ID string; Value float64}`, matched by
  feature id). `Colors string` is a **quantize scheme id** (e.g.
  `"purple_blue_green"`) + `Steps int` + `Domain [2]float64` +
  `UnknownColor string`.
- Sample: `samples.Choropleth() ([]geo.Feature, []geo.ChoroplethDatum)`.

## static registry ids

`static.ChartTypeBar`, `Line`, `Pie`, `BoxPlot`, `Bullet`, `Bump`, `Calendar`,
`Chord`, `CirclePacking`, `Funnel`, `HeatMap`, `Icicle`, `Marimekko`,
`Network`, `ParallelCoordinates`, `PolarBar`, `Radar`, `RadialBar`, `Sankey`,
`ScatterPlot`, `Stream`, `Sunburst`, `SwarmPlot`, `Tree`, `Treemap`,
`Voronoi`, `Waffle`, `Choropleth` — with
`static.RenderChart(chartType, props, overrides) (string, error)` and bundled
demo props in `static.Samples[chartType]`.
