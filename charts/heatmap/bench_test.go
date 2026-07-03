package heatmap_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
)

// benchHeatMapData builds a deterministic n×n grid: n series (rows), each with
// n cells (columns). Cell values are index-derived, so the data is stable
// across runs.
func benchHeatMapData(n int) []heatmap.HeatMapSerie {
	out := make([]heatmap.HeatMapSerie, n)
	for i := 0; i < n; i++ {
		data := make([]heatmap.HeatMapDatum, n)
		for j := 0; j < n; j++ {
			v := float64((i*13+j*29)%100 + 5)
			data[j] = heatmap.HeatMapDatum{X: fmt.Sprintf("x%d", j), Y: &v}
		}
		out[i] = heatmap.HeatMapSerie{ID: fmt.Sprintf("s%d", i), Data: data}
	}
	return out
}

// benchHeatMapProps assembles valid HeatMapProps for an n×n grid.
func benchHeatMapProps(n int) heatmap.HeatMapProps {
	return heatmap.HeatMapProps{
		Width:  700,
		Height: 400,
		Margin: core.Margin{Top: 60, Right: 90, Bottom: 30, Left: 90},
		Data:   benchHeatMapData(n),
	}
}

// BenchmarkHeatMap measures the SVG render cost of the heatmap as the grid
// dimension (n) grows, yielding n×n cells.
func BenchmarkHeatMap(b *testing.B) {
	for _, n := range []int{10, 30, 60} {
		props := benchHeatMapProps(n)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var sb strings.Builder
				if err := heatmap.HeatMap(props).Render(context.Background(), &sb); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
