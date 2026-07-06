package demos

import (
	"fmt"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// Dashboard chart instance ids (htmx registry keys).
const (
	DashboardMainID  = "dash-revenue"
	DashboardDonutID = "dash-share"
	DashboardBarID   = "dash-weekly"
)

// dashCardBg is the dashboard card background; chart themes use it as their
// Background so the SVGs blend into the cards seamlessly.
const dashCardBg = "#1e293b"

// dashPalette is the vibrant-on-dark series palette shared by every chart on
// the dashboard.
var dashPalette = []string{"#38bdf8", "#a78bfa", "#34d399"}

// dashTheme builds the shared dark theme for the dashboard charts.
func dashTheme() *theming.Theme {
	bg := dashCardBg
	gridLine := theming.GridLine{Extra: map[string]any{"stroke": "rgba(148,163,184,0.18)", "strokeWidth": float64(1)}}
	t := theming.ExtendDefaultTheme(theming.DefaultTheme, &theming.PartialTheme{
		Background: &bg,
		Text:       &theming.TextStyle{Fill: "#94a3b8", FontFamily: "sans-serif"},
		Grid:       &theming.PartialGridTheme{Line: &gridLine},
	})
	return &t
}

// dashRevenueSeries is the monthly recurring-revenue trend by segment ($k).
func dashRevenueSeries() []line.LineSeries {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	values := map[string][]float64{
		"subscriptions": {62, 68, 71, 79, 84, 92, 98, 104, 112, 118, 127, 138},
		"services":      {34, 30, 38, 35, 42, 39, 47, 44, 52, 49, 55, 58},
		"licenses":      {18, 22, 19, 25, 22, 28, 24, 30, 27, 33, 29, 35},
	}
	var out []line.LineSeries
	for _, id := range []string{"subscriptions", "services", "licenses"} {
		s := line.LineSeries{ID: id}
		for i, m := range months {
			s.Data = append(s.Data, line.LinePointData{X: m, Y: values[id][i]})
		}
		out = append(out, s)
	}
	return out
}

// DashboardDemos returns the htmx-registered dashboard charts (the KPI
// sparklines and bullets are static and rendered separately).
func DashboardDemos() []Demo {
	theme := dashTheme()
	return []Demo{
		{
			ID:   DashboardMainID,
			Kind: htmx.KindLine,
			Props: line.LineProps{
				Width: 820, Height: 430,
				Margin:     core.Margin{Top: 30, Right: 30, Bottom: 60, Left: 50},
				Curve:      core.CurveMonotoneX,
				Data:       dashRevenueSeries(),
				Theme:      theme,
				Colors:     colors.PaletteColors(dashPalette...),
				Responsive: true,

				LineWidth:   2.5,
				EnableArea:  true,
				AreaOpacity: 1,
				Defs: []core.Def{
					core.LinearGradientDef("dashAreaGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 0.35},
						{Offset: 100, Color: "inherit", Opacity: 0.02},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "dashAreaGrad", Match: "*"}},

				UseMesh:         true,
				EnableCrosshair: true,
				Legends: []legends.LegendProps{
					{
						Anchor: legends.LegendAnchorBottom, Direction: legends.LegendDirectionRow,
						ItemWidth: 120, ItemHeight: 20, TranslateY: 52,
						SymbolShape: legends.SymbolShapeCircle, SymbolSize: 10,
					},
				},
			},
		},
		{
			ID:   DashboardDonutID,
			Kind: htmx.KindPie,
			Props: pie.PieProps{
				Width: 380, Height: 240,
				Margin:       core.Margin{Top: 20, Right: 130, Bottom: 20, Left: 20},
				Data:         dashShareData(),
				Theme:        theme,
				Colors:       colors.PaletteColors(dashPalette...),
				Responsive:   true,
				InnerRadius:  0.65,
				PadAngle:     1.5,
				CornerRadius: 3,

				ActiveOuterRadiusOffset: 6,
				Layers:                  []pie.PieLayerId{pie.PieLayerArcs, pie.PieLayerLegends},
				Defs: []core.Def{
					core.LinearGradientDef("dashPieGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 1},
						{Offset: 100, Color: "inherit", Opacity: 0.55},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "dashPieGrad", Match: "*"}},
				Legends: []legends.LegendProps{
					{
						Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn,
						ItemWidth: 110, ItemHeight: 22, TranslateX: 115,
						SymbolShape: legends.SymbolShapeCircle, SymbolSize: 10,
					},
				},
			},
		},
		{
			ID:   DashboardBarID,
			Kind: htmx.KindBar,
			Props: bar.BarProps{
				Width: 380, Height: 240,
				Margin:     core.Margin{Top: 20, Right: 20, Bottom: 40, Left: 40},
				IndexBy:    "week",
				Keys:       []string{"subscriptions", "services", "licenses"},
				Data:       dashWeeklyData(),
				Theme:      theme,
				Colors:     colors.PaletteColors(dashPalette...),
				Responsive: true,

				BorderRadius: 2,
				// EnableLabel stays default-on; skip labels on segments too
				// small to fit them.
				LabelSkipHeight: 14,
				Defs: []core.Def{
					core.LinearGradientDef("dashBarGrad", []core.GradientStop{
						{Offset: 0, Color: "inherit", Opacity: 1},
						{Offset: 100, Color: "inherit", Opacity: 0.6},
					}, nil),
				},
				Fill: []core.DefRule{{ID: "dashBarGrad", Match: "*"}},
			},
		},
	}
}

// dashShareData is the revenue-share donut data (current quarter, $k).
func dashShareData() []any {
	return []any{
		map[string]any{"id": "subscriptions", "value": float64(383)},
		map[string]any{"id": "services", "value": float64(156)},
		map[string]any{"id": "licenses", "value": float64(97)},
	}
}

// dashWeeklyData is the last six weeks of revenue by segment ($k).
func dashWeeklyData() []bar.BarDatum {
	weeks := []string{"W27", "W28", "W29", "W30", "W31", "W32"}
	subs := []float64{28, 31, 29, 33, 34, 37}
	svcs := []float64{12, 10, 14, 11, 15, 13}
	lics := []float64{6, 8, 5, 9, 7, 10}
	out := make([]bar.BarDatum, len(weeks))
	for i, w := range weeks {
		out[i] = bar.BarDatum{
			"week": w, "subscriptions": subs[i], "services": svcs[i], "licenses": lics[i],
		}
	}
	return out
}

// DashboardBulletProps returns the dark-themed goal-tracking bullets.
func DashboardBulletProps() bullet.BulletProps {
	return bullet.BulletProps{
		Width: 1240, Height: 170,
		Margin:        core.Margin{Top: 10, Right: 30, Bottom: 40, Left: 170},
		TitleAlign:    "start",
		TitleOffsetX:  -150,
		RangeColors:   "seq:greys",
		MeasureColors: "seq:blues",
		MarkerColors:  "seq:oranges",
		Theme:         dashTheme(),
		Responsive:    true,
		Interactive:   true,
		Data: []bullet.BulletItemDatum{
			{ID: "revenue ($k)", Ranges: []float64{0, 400, 800, 1200}, Measures: []float64{1040}, Markers: []float64{1100}},
			{ID: "new customers", Ranges: []float64{0, 150, 300, 450}, Measures: []float64{310}, Markers: []float64{380}},
			{ID: "NPS", Ranges: []float64{0, 30, 60, 90}, Measures: []float64{58}, Markers: []float64{70}},
		},
	}
}

// DashboardKPI is one KPI card at the top of the dashboard.
type DashboardKPI struct {
	Label   string
	Value   string
	Delta   string
	DeltaUp bool
	SVG     string
}

// DashboardKPIs builds the KPI cards: label/value/delta plus a hand-composed
// sparkline (UseLine + area/line generators, no chart component).
func DashboardKPIs() ([]DashboardKPI, error) {
	defs := []struct {
		label, value, delta string
		up                  bool
		color               string
		values              []float64
	}{
		{"Monthly recurring revenue", "$138k", "+8.7%", true, "#38bdf8", []float64{62, 68, 71, 79, 84, 92, 98, 104, 112, 118, 127, 138}},
		{"Active customers", "3,204", "+3.1%", true, "#a78bfa", []float64{2610, 2640, 2705, 2680, 2760, 2850, 2910, 2955, 3040, 3080, 3150, 3204}},
		{"Net revenue retention", "112%", "+1.4%", true, "#34d399", []float64{104, 106, 105, 108, 107, 109, 108, 110, 109, 111, 110, 112}},
		{"Churn", "1.9%", "-0.3%", false, "#fbbf24", []float64{3.1, 2.9, 3.0, 2.7, 2.8, 2.5, 2.6, 2.3, 2.2, 2.1, 2.0, 1.9}},
	}
	out := make([]DashboardKPI, 0, len(defs))
	for i, d := range defs {
		svg, err := dashSparkline(i, d.color, d.values)
		if err != nil {
			return nil, err
		}
		out = append(out, DashboardKPI{
			Label: d.label, Value: d.value, Delta: d.delta, DeltaUp: d.up, SVG: svg,
		})
	}
	return out, nil
}

// dashSparkline renders one KPI sparkline: gradient area + line + end dot,
// composed from the UseLine generators into a bare responsive SVG.
func dashSparkline(idx int, color string, values []float64) (string, error) {
	const w, h = 240.0, 56.0
	s := line.LineSeries{ID: "kpi"}
	for i, v := range values {
		s.Data = append(s.Data, line.LinePointData{X: float64(i), Y: v})
	}
	result := line.UseLine(line.LineProps{
		Width: w - 8, Height: h - 8,
		Curve: core.CurveMonotoneX,
		// min:auto so narrow-range KPIs (e.g. 104→112) still show their shape.
		YScale: scales.ScaleLinearSpec{Min: scales.AutoFloat(), Max: scales.AutoFloat()},
		Data:   []line.LineSeries{s},
	})
	cs := result.Series[0]
	pts := make([]line.PointXY, len(cs.Data))
	for j, d := range cs.Data {
		pts[j] = line.PointXY{X: d.Position.X, Y: d.Position.Y}
	}
	gradID := fmt.Sprintf("dashKpiGrad-%d", idx)

	var inner strings.Builder
	fmt.Fprintf(&inner, `<path d="%s" fill="url(#%s)" stroke-width="0"></path>`, result.AreaGenerator(pts), gradID)
	fmt.Fprintf(&inner, `<path d="%s" fill="none" stroke="%s" stroke-width="2"></path>`, result.LineGenerator(pts), color)
	end := pts[len(pts)-1]
	fmt.Fprintf(&inner, `<circle cx="%.2f" cy="%.2f" r="3" fill="%s"></circle>`, end.X, end.Y, color)

	return renderSVG(core.SvgWrapper(core.SvgWrapperProps{
		Width: w, Height: h,
		Margin:     core.Margin{Top: 4, Left: 4},
		Responsive: true,
		Defs: []core.Def{
			core.LinearGradientDef(gradID, []core.GradientStop{
				{Offset: 0, Color: color, Opacity: 0.35},
				{Offset: 100, Color: color, Opacity: 0.03},
			}, nil),
		},
		Role: "img", AriaLabel: "sparkline",
	}, inner.String()))
}
