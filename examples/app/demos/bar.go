// Package demos holds the demo chart props + data for the example app. Each
// file (bar/line/pie/themes) exports a set of Demo structs describing one
// chart instance: an ID, a title, a description, and the props to register
// with the htmx.Registry.
package demos

import (
	"github.com/geoffjay/templ-charts/charts/annotations"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/samples"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// Demo describes one chart instance on a page: the htmx instance id, a human
// title, a short description, and the props to register. The page handlers
// register the props with the registry and render the chart inline.
type Demo struct {
	ID          string
	Title       string
	Description string
	Kind        htmx.ChartKind
	Props       any
}

// commonChartSize is the default render size for demos (CSS scales the SVG to
// the card width; the viewBox keeps the aspect ratio).
const (
	commonChartWidth  = 700.0
	commonChartHeight = 400.0
)

func defaultMargin() core.Margin {
	return core.Margin{Top: 40, Right: 50, Bottom: 60, Left: 60}
}

// --- Bar demos ---

// BarDemos returns the set of bar chart demos for the /bar page.
func BarDemos() []Demo {
	return []Demo{
		{
			ID:          "bar-stacked",
			Title:       "Stacked vertical",
			Description: "Default stacked layout, multiple keys per index.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
				Data:    barData(),
				Margin:  defaultMargin(),
			},
		},
		{
			ID:          "bar-grouped",
			Title:       "Grouped horizontal",
			Description: "groupMode=grouped, layout=horizontal.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:     commonChartWidth,
				Height:    commonChartHeight,
				IndexBy:   "country",
				Keys:      []string{"hot dogs", "burgers", "sandwich", "kebab"},
				Data:      barData(),
				Margin:    defaultMargin(),
				GroupMode: bar.GroupModeGrouped,
				Layout:    bar.LayoutHorizontal,
			},
		},
		{
			ID:          "bar-markers",
			Title:       "Markers + annotations",
			Description: "A value marker line plus a circle annotation on the top bar.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich"},
				Data:    barData(),
				Margin:  defaultMargin(),
				Markers: []core.CartesianMarker{
					{Axis: "y", Value: float64(150), Legend: "target", LineColor: "#e25c3b"},
				},
				Annotations: []annotations.AnnotationSpec[bar.ComputedBarDatum]{
					{
						Type: annotations.AnnotationTypeCircle,
						Circle: &annotations.CircleAnnotationSpec[bar.ComputedBarDatum]{
							Match:       func(b bar.ComputedBarDatum) bool { return b.Data.ID == "burgers" && b.Data.IndexValue == "USA" },
							Radius:      6,
							Note:        "USA burgers",
							NoteOffsetX: 20, NoteOffsetY: -20,
						},
					},
				},
			},
		},
		{
			ID:          "bar-legend-toggle",
			Title:       "Legend + toggle (HTMX)",
			Description: "Click a legend item to toggle its series on/off via htmx.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:   commonChartWidth,
				Height:  commonChartHeight,
				IndexBy: "country",
				Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
				Data:    barData(),
				Margin:  defaultMargin(),
				Colors:  colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
				Legends: []bar.BarLegendProps{
					{
						LegendProps: legends.LegendProps{
							Anchor:    legends.LegendAnchorTopRight,
							Direction: legends.LegendDirectionColumn,
							ItemWidth: 100, ItemHeight: 20,
							SymbolShape: legends.SymbolShapeSquare,
							TranslateX:  10,
						},
						DataFrom: "keys",
					},
				},
			},
		},
		{
			ID:          "bar-totals",
			Title:       "Totals layer",
			Description: "Stacked bars with the totals label layer enabled.",
			Kind:        htmx.KindBar,
			Props: bar.BarProps{
				Width:        commonChartWidth,
				Height:       commonChartHeight,
				IndexBy:      "country",
				Keys:         []string{"hot dogs", "burgers", "sandwich", "kebab"},
				Data:         barData(),
				Margin:       defaultMargin(),
				EnableTotals: true,
				TotalsOffset: 12,
			},
		},
	}
}

// barData is the shared bar dataset, sourced from the public samples package
// (the single source of truth; samples.Bar also returns the key list).
func barData() []bar.BarDatum {
	d, _ := samples.Bar()
	return d
}

// BarScaleDemo returns the value-scale demo card (id "bar-scale"): one series
// spanning five orders of magnitude, rendered on a linear or log value axis.
// The page toggles logScale so the same data can be compared under both scales.
func BarScaleDemo(logScale bool) Demo {
	p := bar.BarProps{
		Width:     commonChartWidth,
		Height:    commonChartHeight,
		IndexBy:   "tier",
		Keys:      []string{"requests"},
		GroupMode: bar.GroupModeGrouped,
		Margin:    defaultMargin(),
		Data: []bar.BarDatum{
			{"tier": "cache", "requests": 1200000.0},
			{"tier": "cdn", "requests": 340000.0},
			{"tier": "app", "requests": 42000.0},
			{"tier": "db", "requests": 3800.0},
			{"tier": "queue", "requests": 210.0},
			{"tier": "audit", "requests": 12.0},
		},
	}
	desc := "Requests/day per tier, spanning ~12 to ~1.2M. On a linear axis the small tiers vanish; on a log axis every tier stays readable."
	if logScale {
		// Min pinned to 1 (not auto) so the smallest tier clears the baseline.
		p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()}
	}
	return Demo{ID: "bar-scale", Title: "Value scale (linear / log)", Description: desc, Kind: htmx.KindBar, Props: p}
}
