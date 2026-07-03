package bar_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/core"
)

// benchBarKeys is the set of series keys per index group (~6 series, matching
// the bar demo shape).
var benchBarKeys = []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"}

// benchBarData builds a deterministic bar dataset with n index groups, each
// carrying a value for every key in benchBarKeys. Values are index-derived so
// the data is stable across runs.
func benchBarData(n int) []bar.BarDatum {
	data := make([]bar.BarDatum, n)
	for i := 0; i < n; i++ {
		row := map[string]any{"country": fmt.Sprintf("c%d", i)}
		for j, k := range benchBarKeys {
			row[k] = float64((i*7+j*13)%100 + 10)
		}
		data[i] = row
	}
	return data
}

// benchBarProps assembles valid BarProps for a dataset of n index groups.
func benchBarProps(n int) bar.BarProps {
	return bar.BarProps{
		Width:   700,
		Height:  400,
		Margin:  core.Margin{Top: 40, Right: 50, Bottom: 60, Left: 60},
		IndexBy: "country",
		Keys:    benchBarKeys,
		Data:    benchBarData(n),
	}
}

// BenchmarkBar measures the SVG render cost of the bar chart as the number of
// index groups (n) grows.
func BenchmarkBar(b *testing.B) {
	for _, n := range []int{10, 100, 500} {
		props := benchBarProps(n)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var sb strings.Builder
				if err := bar.Bar(props).Render(context.Background(), &sb); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
