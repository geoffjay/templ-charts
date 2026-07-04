package scatterplot_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scatterplot"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// benchScatterProps builds a single-series scatterplot of n points.
func benchScatterProps(n int, engine theming.Engine) scatterplot.ScatterPlotProps {
	data := make([]scatterplot.ScatterPlotDatum, n)
	for i := 0; i < n; i++ {
		data[i] = scatterplot.ScatterPlotDatum{
			X: float64((i*37)%1000) + 1,
			Y: float64((i*53)%1000) + 1,
		}
	}
	return scatterplot.ScatterPlotProps{
		Width: 900, Height: 500,
		Margin:      core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data:        []scatterplot.ScatterPlotSerie{{ID: "s", Data: data}},
		EnableGridX: true,
		EnableGridY: true,
		Render:      engine,
		ChartID:     "bench",
	}
}

func renderScatter(b *testing.B, props scatterplot.ScatterPlotProps) int {
	b.Helper()
	var sb strings.Builder
	if err := scatterplot.ScatterPlot(props).Render(context.Background(), &sb); err != nil {
		b.Fatal(err)
	}
	return sb.Len()
}

// BenchmarkScatterplotSVGvsCanvas contrasts the two backends' emit cost and
// payload size as the point count grows: the SVG path builds one <g><circle>
// per point (payload grows fast), while the Canvas path records one FillCircle
// op per point into a compact draw-list. Run with `make bench`.
func BenchmarkScatterplotSVGvsCanvas(b *testing.B) {
	for _, n := range []int{1000, 5000, 20000} {
		svgProps := benchScatterProps(n, theming.EngineSVG)
		b.Run(fmt.Sprintf("svg/n=%d", n), func(b *testing.B) {
			bytes := renderScatter(b, svgProps)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderScatter(b, svgProps)
			}
			b.ReportMetric(float64(bytes), "payload_B")
		})

		canvasProps := benchScatterProps(n, theming.EngineCanvas)
		b.Run(fmt.Sprintf("canvas/n=%d", n), func(b *testing.B) {
			bytes := renderScatter(b, canvasProps)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				renderScatter(b, canvasProps)
			}
			b.ReportMetric(float64(bytes), "payload_B")
		})
	}
}
