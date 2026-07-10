package demos

import (
	"fmt"
	"strings"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
)

// HTMLLegendDemoID is the tile on the /legends page whose legend is plain
// HTML rendered outside the SVG (the handler attaches the chips as the card's
// FooterHTML).
const HTMLLegendDemoID = "legends-html"

// LegendsDemos returns the /legends page tiles: symbol shapes, anchors and
// directions, symbol borders, the continuous color legend, and a chart whose
// legend is custom HTML outside the SVG.
func LegendsDemos() []Demo {
	return []Demo{
		{
			ID:          "legends-symbol-shapes",
			Title:       "Symbol shapes × anchors",
			Description: "One legend per corner, each with a different SymbolShape — circle (top-left), diamond (top-right), square (bottom-left), triangle (bottom-right). All four toggle the same series (HTMX).",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: core.Margin{Top: 40, Right: 130, Bottom: 60, Left: 60},
				Curve:  core.CurveMonotoneX,
				Data:   stylingLineData(),
				Legends: []legends.LegendProps{
					{
						Anchor: legends.LegendAnchorTopLeft, Direction: legends.LegendDirectionColumn,
						ItemWidth: 90, ItemHeight: 20, TranslateX: 25,
						SymbolShape: legends.SymbolShapeCircle,
					},
					{
						Anchor: legends.LegendAnchorTopRight, Direction: legends.LegendDirectionColumn,
						ItemWidth: 90, ItemHeight: 20, TranslateX: 100,
						SymbolShape: legends.SymbolShapeDiamond,
					},
					{
						Anchor: legends.LegendAnchorBottomLeft, Direction: legends.LegendDirectionColumn,
						ItemWidth: 90, ItemHeight: 20, TranslateX: 25,
						SymbolShape: legends.SymbolShapeSquare,
					},
					{
						Anchor: legends.LegendAnchorBottomRight, Direction: legends.LegendDirectionColumn,
						ItemWidth: 90, ItemHeight: 20, TranslateX: 100,
						SymbolShape: legends.SymbolShapeTriangle,
					},
				},
			},
		},
		{
			ID:          "legends-bottom-row",
			Title:       "Bottom row + symbol borders",
			Description: "A row legend anchored below the plot (Anchor: bottom, Direction: row, TranslateY into the bottom margin) with bordered diamond symbols — the most common real-world legend placement.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: core.Margin{Top: 40, Right: 50, Bottom: 95, Left: 60},
				Curve:  core.CurveMonotoneX,
				Data:   stylingLineData(),
				Legends: []legends.LegendProps{
					{
						Anchor: legends.LegendAnchorBottom, Direction: legends.LegendDirectionRow,
						ItemWidth: 100, ItemHeight: 20, TranslateY: 70,
						SymbolShape:       legends.SymbolShapeDiamond,
						SymbolSize:        14,
						SymbolBorderColor: "#333333",
						SymbolBorderWidth: 1.5,
					},
				},
			},
		},
		{
			ID:          "legends-continuous",
			Title:       "Continuous color legend",
			Description: "The value→color scale rendered as a titled gradient bar (legends.ContinuousColorsLegendSvg): Length, Thickness, Title, and anchor/translate are all configurable.",
			Kind:        htmx.KindHeatmap,
			Props: heatmap.HeatMapProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: core.Margin{Top: 60, Right: 90, Bottom: 90, Left: 90},
				Data:   heatmapData(),
				Colors: heatmap.HeatMapColorConfig{Type: "sequential", Scheme: "blues"},
				Legends: []heatmap.HeatMapLegend{
					{
						Anchor:     legends.LegendAnchorBottom,
						TranslateY: 66,
						Length:     320,
						Thickness:  14,
						Title:      "temperature (°C)",
					},
				},
			},
		},
		{
			ID:          HTMLLegendDemoID,
			Title:       "HTML legend outside the SVG",
			Description: "The chart renders with no SVG legend; the chips below are plain HTML built from the same series data, wired to the same HTMX toggle endpoint a built-in legend uses. Click a chip to hide/show its series.",
			Kind:        htmx.KindLine,
			Props: line.LineProps{
				Width:  commonChartWidth,
				Height: commonChartHeight,
				Margin: defaultMargin(),
				Curve:  core.CurveMonotoneX,
				Data:   stylingLineData(),

				EnableArea:  core.BoolPtr(true),
				AreaOpacity: 1,
				Defs: []core.Def{
					core.LinearGradientDef("htmlLegendGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 0.4},
						{Offset: 100, Color: "inherit", Opacity: 0.02},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "htmlLegendGrad", Match: "*"}},
			},
		},
	}
}

// HTMLLegendFooter builds the HTML legend chips for the HTMLLegendDemoID
// tile: one button per series, colored like the chart (same "nivo" scheme in
// series order) and posting to the chart's HTMX toggle endpoint. The onclick
// class toggle dims the chip locally; the series swap comes from the server.
func HTMLLegendFooter() string {
	scheme := colors.CategoricalColorSchemes["nivo"]
	var b strings.Builder
	b.WriteString(`<div class="html-legend">`)
	for i, s := range stylingLineData() {
		color := scheme[i%len(scheme)]
		last := s.Data[len(s.Data)-1].Y
		b.WriteString(fmt.Sprintf(
			`<button hx-post="/charts/%[1]s/toggle?series=%[2]s" hx-target="#chart-%[1]s" hx-swap="innerHTML" onclick="this.classList.toggle('off')">`+
				`<span class="dot" style="background:%[3]s"></span>%[2]s <span class="val">%.0[4]f</span></button>`,
			HTMLLegendDemoID, s.ID, color, last,
		))
	}
	b.WriteString(`</div>`)
	return b.String()
}
