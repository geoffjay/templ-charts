package static

import (
	"fmt"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/chord"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/funnel"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/marimekko"
	"github.com/geoffjay/templ-charts/charts/network"
	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/polarbar"
	"github.com/geoffjay/templ-charts/charts/radar"
	"github.com/geoffjay/templ-charts/charts/radialbar"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/sankey"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/stream"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/swarmplot"
	"github.com/geoffjay/templ-charts/charts/tree"
	"github.com/geoffjay/templ-charts/charts/treemap"
	"github.com/geoffjay/templ-charts/charts/voronoi"
	"github.com/geoffjay/templ-charts/charts/waffle"
)

// ChartType constants for every supported family. Bar/Line/Pie keep their
// original string values (referenced by existing callers/tests).
const (
	ChartTypeBar                 ChartType = "bar"
	ChartTypeLine                ChartType = "line"
	ChartTypePie                 ChartType = "pie"
	ChartTypeBoxPlot             ChartType = "boxplot"
	ChartTypeBullet              ChartType = "bullet"
	ChartTypeBump                ChartType = "bump"
	ChartTypeCalendar            ChartType = "calendar"
	ChartTypeChord               ChartType = "chord"
	ChartTypeCirclePacking       ChartType = "circlepacking"
	ChartTypeFunnel              ChartType = "funnel"
	ChartTypeHeatMap             ChartType = "heatmap"
	ChartTypeIcicle              ChartType = "icicle"
	ChartTypeMarimekko           ChartType = "marimekko"
	ChartTypeNetwork             ChartType = "network"
	ChartTypeParallelCoordinates ChartType = "parallelcoordinates"
	ChartTypePolarBar            ChartType = "polarbar"
	ChartTypeRadar               ChartType = "radar"
	ChartTypeRadialBar           ChartType = "radialbar"
	ChartTypeSankey              ChartType = "sankey"
	ChartTypeScatterPlot         ChartType = "scatterplot"
	ChartTypeStream              ChartType = "stream"
	ChartTypeSunburst            ChartType = "sunburst"
	ChartTypeSwarmPlot           ChartType = "swarmplot"
	ChartTypeTree                ChartType = "tree"
	ChartTypeTreemap             ChartType = "treemap"
	ChartTypeVoronoi             ChartType = "voronoi"
	ChartTypeWaffle              ChartType = "waffle"
	ChartTypeChoropleth          ChartType = "choropleth"
)

// Sample is one demo chart: its type + ready-to-render props. Mirrors
// @nivo/static's { type, props } sample entry.
type Sample struct {
	Type  ChartType
	Props any
}

// comp builds a reflectComponent from a chart's templ constructor. The generic
// closure type-asserts the props (produced by the reflective adapter) to P and
// calls ctor; a wrong props type yields a descriptive error. This is the only
// per-family code the 28-way dispatch needs — everything else (static flags,
// margin default, whitelisted overrides) is handled generically in types.go.
func comp[P any](name ChartType, ctor func(P) templ.Component) reflectComponent {
	return reflectComponent{render: func(p any) (templ.Component, error) {
		pp, ok := p.(P)
		if !ok {
			return nil, fmt.Errorf("static: %s render expects %T, got %T", name, *new(P), p)
		}
		return ctor(pp), nil
	}}
}

// sampleW/sampleH are the default render geometry for the Samples entries.
const sampleW, sampleH = 700.0, 460.0

// colorRuntime is the whitelist for families whose Colors field is an ordinal
// color config; the rest expose only width/height (their color is bespoke,
// keyed, inherited, or a plain scheme string, which the generic override skips).
var (
	colorRuntime = []string{"width", "height", "colors"}
	sizeRuntime  = []string{"width", "height"}
)

// ChartsMapping maps each ChartType to its reflective component + whitelist.
var ChartsMapping = map[ChartType]Mapping{
	ChartTypeBar:                 {Component: comp(ChartTypeBar, bar.Bar), RuntimeProps: []string{"width", "height", "colors", "groupMode"}, Defaults: map[string]any{}},
	ChartTypeLine:                {Component: comp(ChartTypeLine, line.Line), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypePie:                 {Component: comp(ChartTypePie, pie.Pie), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeBoxPlot:             {Component: comp(ChartTypeBoxPlot, boxplot.BoxPlot), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeBullet:              {Component: comp(ChartTypeBullet, bullet.Bullet), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
	ChartTypeBump:                {Component: comp(ChartTypeBump, bump.Bump), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeCalendar:            {Component: comp(ChartTypeCalendar, calendar.Calendar), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
	ChartTypeChord:               {Component: comp(ChartTypeChord, chord.Chord), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeCirclePacking:       {Component: comp(ChartTypeCirclePacking, cp.CirclePacking), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeFunnel:              {Component: comp(ChartTypeFunnel, funnel.Funnel), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeHeatMap:             {Component: comp(ChartTypeHeatMap, heatmap.HeatMap), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
	ChartTypeIcicle:              {Component: comp(ChartTypeIcicle, icicle.Icicle), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeMarimekko:           {Component: comp(ChartTypeMarimekko, marimekko.Marimekko), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeNetwork:             {Component: comp(ChartTypeNetwork, network.Network), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
	ChartTypeParallelCoordinates: {Component: comp(ChartTypeParallelCoordinates, pc.ParallelCoordinates), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypePolarBar:            {Component: comp(ChartTypePolarBar, polarbar.PolarBar), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeRadar:               {Component: comp(ChartTypeRadar, radar.Radar), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeRadialBar:           {Component: comp(ChartTypeRadialBar, radialbar.RadialBar), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeSankey:              {Component: comp(ChartTypeSankey, sankey.Sankey), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeScatterPlot:         {Component: comp(ChartTypeScatterPlot, scatterplot.ScatterPlot), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeStream:              {Component: comp(ChartTypeStream, stream.Stream), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeSunburst:            {Component: comp(ChartTypeSunburst, sunburst.Sunburst), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeSwarmPlot:           {Component: comp(ChartTypeSwarmPlot, swarmplot.SwarmPlot), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeTree:                {Component: comp(ChartTypeTree, tree.Tree), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
	ChartTypeTreemap:             {Component: comp(ChartTypeTreemap, treemap.Treemap), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeVoronoi:             {Component: comp(ChartTypeVoronoi, voronoi.Voronoi), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeWaffle:              {Component: comp(ChartTypeWaffle, waffle.Waffle), RuntimeProps: colorRuntime, Defaults: map[string]any{}},
	ChartTypeChoropleth:          {Component: comp(ChartTypeChoropleth, geo.Choropleth), RuntimeProps: sizeRuntime, Defaults: map[string]any{}},
}

// Samples is the ready-to-render demo data registry for every family, sourced
// from the public charts/samples package (the single source of truth).
var Samples = map[ChartType]Sample{}

func init() {
	barData, barKeys := samples.Bar()
	Samples[ChartTypeBar] = Sample{Type: ChartTypeBar, Props: bar.BarProps{
		Width: sampleW, Height: sampleH, IndexBy: "country", Keys: barKeys, Data: barData,
	}}

	Samples[ChartTypeLine] = Sample{Type: ChartTypeLine, Props: line.LineProps{
		Width: sampleW, Height: sampleH, Curve: core.CurveMonotoneX, Data: samples.Line(),
	}}

	Samples[ChartTypePie] = Sample{Type: ChartTypePie, Props: pie.PieProps{
		Width: sampleW, Height: sampleH, InnerRadius: 0.5, PadAngle: 0.5, Data: samples.Pie(),
	}}

	Samples[ChartTypeBoxPlot] = Sample{Type: ChartTypeBoxPlot, Props: boxplot.BoxPlotProps{
		Width: sampleW, Height: sampleH, Data: samples.BoxPlot(), ColorBy: "group",
	}}

	Samples[ChartTypeBullet] = Sample{Type: ChartTypeBullet, Props: bullet.BulletProps{
		Width: sampleW, Height: 260, Margin: core.Margin{Top: 20, Right: 30, Bottom: 40, Left: 100}, Data: samples.Bullet(),
	}}

	Samples[ChartTypeBump] = Sample{Type: ChartTypeBump, Props: bump.BumpProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 30, Right: 100, Bottom: 40, Left: 100}, Data: samples.Bump(),
	}}

	calFrom, calTo, calData := samples.Calendar()
	Samples[ChartTypeCalendar] = Sample{Type: ChartTypeCalendar, Props: calendar.CalendarProps{
		Width: sampleW, Height: 220, Margin: core.Margin{Top: 40, Right: 20, Bottom: 10, Left: 40},
		From: calFrom, To: calTo, Data: calData,
	}}

	chordMatrix, chordKeys := samples.Chord()
	Samples[ChartTypeChord] = Sample{Type: ChartTypeChord, Props: chord.ChordProps{
		Width: 520, Height: 520, Margin: core.Margin{Top: 60, Right: 60, Bottom: 60, Left: 60},
		Data: chordMatrix, Keys: chordKeys,
	}}

	Samples[ChartTypeCirclePacking] = Sample{Type: ChartTypeCirclePacking, Props: cp.CirclePackingProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10}, Data: samples.CirclePacking(),
	}}

	Samples[ChartTypeFunnel] = Sample{Type: ChartTypeFunnel, Props: funnel.FunnelProps{
		Width: sampleW, Height: 460, Margin: core.Margin{Top: 20, Right: 30, Bottom: 20, Left: 30}, Data: samples.Funnel(),
	}}

	Samples[ChartTypeHeatMap] = Sample{Type: ChartTypeHeatMap, Props: heatmap.HeatMapProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 60, Right: 90, Bottom: 30, Left: 90}, Data: samples.Heatmap(),
	}}

	Samples[ChartTypeIcicle] = Sample{Type: ChartTypeIcicle, Props: icicle.IcicleProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10}, Data: samples.Icicle(),
	}}

	mkData, mkDims := samples.Marimekko()
	Samples[ChartTypeMarimekko] = Sample{Type: ChartTypeMarimekko, Props: marimekko.MarimekkoProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 20, Right: 160, Bottom: 50, Left: 60},
		Data: mkData, Dimensions: mkDims,
	}}

	netNodes, netLinks := samples.Network()
	Samples[ChartTypeNetwork] = Sample{Type: ChartTypeNetwork, Props: network.NetworkProps{
		Width: 460, Height: 460, Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
		Nodes: netNodes, Links: netLinks, LinkDistance: 90, Repulsivity: 120,
	}}

	pcData, pcVars := samples.ParallelCoordinates()
	Samples[ChartTypeParallelCoordinates] = Sample{Type: ChartTypeParallelCoordinates, Props: pc.PCProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 50, Right: 60, Bottom: 50, Left: 60},
		Data: pcData, Variables: pcVars,
	}}

	pbData, pbKeys := samples.PolarBar()
	Samples[ChartTypePolarBar] = Sample{Type: ChartTypePolarBar, Props: polarbar.PolarBarProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 40, Right: 60, Bottom: 40, Left: 60},
		Data: pbData, Keys: pbKeys,
	}}

	radarData, radarKeys := samples.Radar()
	Samples[ChartTypeRadar] = Sample{Type: ChartTypeRadar, Props: radar.RadarProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 70, Right: 90, Bottom: 50, Left: 90},
		Data: radarData, Keys: radarKeys, IndexBy: "taste",
	}}

	Samples[ChartTypeRadialBar] = Sample{Type: ChartTypeRadialBar, Props: radialbar.RadialBarProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 40, Right: 60, Bottom: 40, Left: 60}, Data: samples.RadialBar(),
	}}

	sankeyNodes, sankeyLinks := samples.Sankey()
	Samples[ChartTypeSankey] = Sample{Type: ChartTypeSankey, Props: sankey.SankeyProps{
		Width: sampleW, Height: 420, Margin: core.Margin{Top: 20, Right: 100, Bottom: 20, Left: 100},
		Nodes: sankeyNodes, Links: sankeyLinks,
	}}

	Samples[ChartTypeScatterPlot] = Sample{Type: ChartTypeScatterPlot, Props: scatterplot.ScatterPlotProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data: samples.ScatterPlot(), EnableGridX: core.BoolPtr(true), EnableGridY: core.BoolPtr(true),
	}}

	streamData, streamKeys := samples.Stream()
	Samples[ChartTypeStream] = Sample{Type: ChartTypeStream, Props: stream.StreamProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 50},
		Data: streamData, Keys: streamKeys,
	}}

	Samples[ChartTypeSunburst] = Sample{Type: ChartTypeSunburst, Props: sunburst.SunburstProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10}, Data: samples.Sunburst(),
	}}

	swarmData, swarmGroups := samples.SwarmPlot()
	Samples[ChartTypeSwarmPlot] = Sample{Type: ChartTypeSwarmPlot, Props: swarmplot.SwarmPlotProps{
		Width: 500, Height: 460, Margin: core.Margin{Top: 20, Right: 20, Bottom: 60, Left: 60},
		Data: swarmData, Groups: swarmGroups, Size: 8,
	}}

	Samples[ChartTypeTree] = Sample{Type: ChartTypeTree, Props: tree.TreeProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 30, Right: 40, Bottom: 40, Left: 40}, Data: samples.Tree(),
	}}

	Samples[ChartTypeTreemap] = Sample{Type: ChartTypeTreemap, Props: treemap.TreemapProps{
		Width: sampleW, Height: sampleH, Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10}, Data: samples.Treemap(),
	}}

	Samples[ChartTypeVoronoi] = Sample{Type: ChartTypeVoronoi, Props: voronoi.VoronoiProps{
		Width: 460, Height: 460, Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20}, Data: samples.Voronoi(),
	}}

	Samples[ChartTypeWaffle] = Sample{Type: ChartTypeWaffle, Props: waffle.WaffleProps{
		Width: sampleH, Height: sampleH, Margin: core.Margin{Top: 10, Right: 10, Bottom: 10, Left: 10},
		Total: 100, Rows: 10, Columns: 10, Data: samples.Waffle(),
	}}

	choroFeatures, choroData := samples.Choropleth()
	Samples[ChartTypeChoropleth] = Sample{Type: ChartTypeChoropleth, Props: geo.ChoroplethProps{
		Features: choroFeatures, Data: choroData,
		Colors: "purple_blue_green", Steps: 5, UnknownColor: "#dddddd", ValueFormat: ",.0f",
		GeoBase: geo.GeoBase{
			Width: 600, Height: 400, Margin: core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
			ProjectionType: geo.ProjectionEquirectangular, ProjectionScale: 150,
			BorderWidth: 0.5, BorderColor: "#152238",
		},
	}}
}
