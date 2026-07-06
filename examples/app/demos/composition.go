package demos

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/axes"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// StaticDemo is one pre-rendered tile: no htmx registration, just the SVG.
// Used by the /composition page, whose charts are assembled by hand from the
// Use* hooks + sub-components instead of a top-level chart component.
type StaticDemo struct {
	ID          string
	Title       string
	Description string
	SVG         string
}

// CompositionDemos returns the /composition page tiles. Each chart here is
// built the way a custom d3 chart would be: call line.UseLine for the
// computed scales/series/generators, reuse the exported sub-components
// (axes.Grid, axes.Axes, line.Lines) for the standard layers, and write
// custom SVG marks for everything else.
func CompositionDemos() ([]StaticDemo, error) {
	custom, err := compCustomPoints()
	if err != nil {
		return nil, err
	}
	labels, err := compDirectLabels()
	if err != nil {
		return nil, err
	}
	sparks, err := compSparklines()
	if err != nil {
		return nil, err
	}
	return []StaticDemo{
		{
			ID:          "comp-custom-points",
			Title:       "Custom point symbols",
			Description: "line.UseLine computes scales, positioned series, and generators; grid/axes/lines render via the exported sub-components; the points layer is replaced with hand-written SVG — diamond markers, plus a ring + value label on each series' maximum.",
			SVG:         custom,
		},
		{
			ID:          "comp-direct-labels",
			Title:       "Direct labels instead of a legend",
			Description: "The classic d3 pattern: drop the legend and label each line at its end point, in the series color. Everything after UseLine is four Fprintf calls.",
			SVG:         labels,
		},
		{
			ID:          "comp-sparklines",
			Title:       "Sparkline KPI cards",
			Description: "Three KPI cards in one SVG: each runs its own UseLine at 190×48 with axes and grid omitted, draws a gradient area + line + end dot from the returned generators, and lays SVG text on top. The whole card row is a composition over the hooks — no chart component involved.",
			SVG:         sparks,
		},
	}, nil
}

// renderSVG renders a templ component to its SVG string.
func renderSVG(c templ.Component) (string, error) {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// rawY returns a series point's raw Y value as float64.
func rawY(d line.ComputedDatum) float64 {
	f, _ := d.Data.Y.(float64)
	return f
}

// defaultAxes builds bottom + left axis props for the computed scales.
func defaultAxes(result line.LineResult, innerW, innerH float64) (axes.AxisProps, axes.AxisProps) {
	xa := axes.DefaultAxisProps
	xa.Axis, xa.Scale, xa.Length = "x", result.XScale, innerW
	xa.Y = innerH
	xa.TicksPosition = "after"
	ya := axes.DefaultAxisProps
	ya.Axis, ya.Scale, ya.Length = "y", result.YScale, innerH
	ya.TicksPosition = "after"
	return xa, ya
}

// compCustomPoints builds the custom-point-symbols tile.
func compCustomPoints() (string, error) {
	dims := core.UseDimensions(commonChartWidth, commonChartHeight, defaultMargin())
	theme := &theming.DefaultTheme
	result := line.UseLine(line.LineProps{
		Width: dims.InnerWidth, Height: dims.InnerHeight,
		Curve: core.CurveMonotoneX,
		Data:  stylingLineData(),
	})

	var inner strings.Builder
	grid, err := renderSVG(axes.Grid(axes.GridProps{
		Axis: "y", Scale: result.YScale,
		Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
	}))
	if err != nil {
		return "", err
	}
	inner.WriteString(grid)

	xa, ya := defaultAxes(result, dims.InnerWidth, dims.InnerHeight)
	axesSVG, err := renderSVG(axes.Axes(axes.AxesProps{
		XAxis: &xa, YAxis: &ya,
		Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
	}))
	if err != nil {
		return "", err
	}
	inner.WriteString(axesSVG)

	lines, err := renderSVG(line.Lines(line.LinesProps{
		Series: result.Series, LineGenerator: result.LineGenerator, LineWidth: 2.5,
	}))
	if err != nil {
		return "", err
	}
	inner.WriteString(lines)

	// Custom marks: a rotated-square (diamond) at every point, and a ring +
	// value label on each series' maximum.
	for _, s := range result.Series {
		maxI := 0
		for i := range s.Data {
			if rawY(s.Data[i]) > rawY(s.Data[maxI]) {
				maxI = i
			}
		}
		for i, d := range s.Data {
			x, y := d.Position.X, d.Position.Y
			fmt.Fprintf(&inner,
				`<rect transform="translate(%.2f,%.2f) rotate(45)" x="-3.5" y="-3.5" width="7" height="7" fill="%s" stroke="#ffffff" stroke-width="1.5"></rect>`,
				x, y, s.Color)
			if i == maxI {
				fmt.Fprintf(&inner,
					`<circle cx="%.2f" cy="%.2f" r="9" fill="none" stroke="%s" stroke-width="1.5"></circle>`,
					x, y, s.Color)
				fmt.Fprintf(&inner,
					`<text x="%.2f" y="%.2f" text-anchor="middle" font-size="11" font-weight="600" font-family="sans-serif" fill="%s">%.0f</text>`,
					x, y-16, s.Color, rawY(d))
			}
		}
	}

	return renderSVG(core.SvgWrapper(core.SvgWrapperProps{
		Width: dims.OuterWidth, Height: dims.OuterHeight, Margin: dims.Margin,
		Role: "img", AriaLabel: "line chart with custom diamond point symbols",
	}, inner.String()))
}

// compDirectLabels builds the direct-labels tile.
func compDirectLabels() (string, error) {
	margin := core.Margin{Top: 40, Right: 110, Bottom: 60, Left: 60}
	dims := core.UseDimensions(commonChartWidth, commonChartHeight, margin)
	theme := &theming.DefaultTheme
	result := line.UseLine(line.LineProps{
		Width: dims.InnerWidth, Height: dims.InnerHeight,
		Curve: core.CurveMonotoneX,
		Data:  stylingLineData(),
	})

	var inner strings.Builder
	xa, _ := defaultAxes(result, dims.InnerWidth, dims.InnerHeight)
	axesSVG, err := renderSVG(axes.Axes(axes.AxesProps{
		XAxis: &xa,
		Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
	}))
	if err != nil {
		return "", err
	}
	inner.WriteString(axesSVG)

	lines, err := renderSVG(line.Lines(line.LinesProps{
		Series: result.Series, LineGenerator: result.LineGenerator, LineWidth: 3,
	}))
	if err != nil {
		return "", err
	}
	inner.WriteString(lines)

	// End-point dot + series label in the series color, replacing the legend.
	for _, s := range result.Series {
		if len(s.Data) == 0 {
			continue
		}
		end := s.Data[len(s.Data)-1]
		fmt.Fprintf(&inner,
			`<circle cx="%.2f" cy="%.2f" r="4.5" fill="%s" stroke="#ffffff" stroke-width="1.5"></circle>`,
			end.Position.X, end.Position.Y, s.Color)
		fmt.Fprintf(&inner,
			`<text x="%.2f" y="%.2f" font-size="13" font-weight="600" font-family="sans-serif" fill="%s">%s</text>`,
			end.Position.X+10, end.Position.Y+4, s.Color, s.ID)
	}

	return renderSVG(core.SvgWrapper(core.SvgWrapperProps{
		Width: dims.OuterWidth, Height: dims.OuterHeight, Margin: dims.Margin,
		Role: "img", AriaLabel: "line chart with direct series labels",
	}, inner.String()))
}

// compSparklines builds the sparkline-KPI-cards tile: three cards in one SVG.
func compSparklines() (string, error) {
	const (
		svgW, svgH   = 700.0, 170.0
		cardW, cardH = 212.0, 150.0
		gap          = 22.0
	)
	scheme := colors.CategoricalColorSchemes["nivo"]
	series := stylingLineData()
	labels := map[string]string{"mobile": "Mobile sessions", "desktop": "Desktop sessions", "tablet": "Tablet sessions"}

	var defs []core.Def
	var inner strings.Builder
	for i, s := range series {
		color := scheme[i%len(scheme)]
		gradID := fmt.Sprintf("sparkGrad-%d", i)
		defs = append(defs, core.LinearGradientDef(gradID, []core.GradientStop{
			{Offset: 0, Color: color, Opacity: 0.35},
			{Offset: 100, Color: color, Opacity: 0.03},
		}, nil))

		// Per-card mini line: its own UseLine with the plot sized to the card.
		plotW, plotH := cardW-32, 46.0
		result := line.UseLine(line.LineProps{
			Width: plotW, Height: plotH,
			Curve: core.CurveMonotoneX,
			Data:  []line.LineSeries{s},
		})
		cs := result.Series[0]
		pts := make([]line.PointXY, len(cs.Data))
		for j, d := range cs.Data {
			pts[j] = line.PointXY{X: d.Position.X, Y: d.Position.Y}
		}
		first, last := rawY(cs.Data[0]), rawY(cs.Data[len(cs.Data)-1])
		deltaPct := (last - first) / first * 100
		deltaColor, deltaSign := "#0a7d4b", "▲"
		if deltaPct < 0 {
			deltaColor, deltaSign = "#c2402a", "▼"
		}

		x0 := float64(i) * (cardW + gap)
		fmt.Fprintf(&inner, `<g transform="translate(%.0f,0)">`, x0)
		fmt.Fprintf(&inner, `<rect width="%.0f" height="%.0f" rx="10" fill="#fafbfc" stroke="#e2e6ec"></rect>`, cardW, cardH)
		fmt.Fprintf(&inner, `<text x="16" y="28" font-size="12" font-family="sans-serif" fill="#666">%s</text>`, labels[s.ID])
		fmt.Fprintf(&inner, `<text x="16" y="58" font-size="26" font-weight="700" font-family="sans-serif" fill="#222">%.0f</text>`, last)
		fmt.Fprintf(&inner, `<text x="%.0f" y="58" text-anchor="end" font-size="12" font-weight="600" font-family="sans-serif" fill="%s">%s %.1f%%</text>`,
			cardW-16, deltaColor, deltaSign, deltaPct)

		// The sparkline plot, inset in the lower half of the card.
		fmt.Fprintf(&inner, `<g transform="translate(16,%.0f)">`, cardH-plotH-16)
		fmt.Fprintf(&inner, `<path d="%s" fill="url(#%s)" stroke-width="0"></path>`, result.AreaGenerator(pts), gradID)
		fmt.Fprintf(&inner, `<path d="%s" fill="none" stroke="%s" stroke-width="2"></path>`, result.LineGenerator(pts), color)
		endPt := pts[len(pts)-1]
		fmt.Fprintf(&inner, `<circle cx="%.2f" cy="%.2f" r="3.5" fill="%s" stroke="#ffffff" stroke-width="1.5"></circle>`, endPt.X, endPt.Y, color)
		inner.WriteString(`</g></g>`)
	}

	return renderSVG(core.SvgWrapper(core.SvgWrapperProps{
		Width: svgW, Height: svgH,
		Margin: core.Margin{Top: 10, Left: 10},
		Defs:   defs,
		Role:   "img", AriaLabel: "three sparkline KPI cards",
	}, inner.String()))
}
