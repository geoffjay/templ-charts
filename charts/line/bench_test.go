package line_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/line"
)

// benchLineData builds a deterministic line dataset of ~3 series, each with n
// points. X is the point index (numeric) and Y is index-derived, so the data
// is stable across runs.
func benchLineData(n int) []line.LineSeries {
	const seriesCount = 3
	series := make([]line.LineSeries, seriesCount)
	for s := 0; s < seriesCount; s++ {
		data := make([]line.LinePointData, n)
		for j := 0; j < n; j++ {
			data[j] = line.LinePointData{
				X: float64(j),
				Y: float64((s*3+j*11)%90 + 20),
			}
		}
		series[s] = line.LineSeries{ID: fmt.Sprintf("series-%d", s), Data: data}
	}
	return series
}

// benchLineProps assembles valid LineProps for ~3 series of n points each.
func benchLineProps(n int) line.LineProps {
	return line.LineProps{
		Width:  700,
		Height: 400,
		Margin: core.Margin{Top: 40, Right: 50, Bottom: 60, Left: 60},
		Curve:  core.CurveMonotoneX,
		Data:   benchLineData(n),
	}
}

// BenchmarkLine measures the SVG render cost of the line chart as the number
// of points per series (n) grows.
func BenchmarkLine(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		props := benchLineProps(n)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var sb strings.Builder
				if err := line.Line(props).Render(context.Background(), &sb); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
