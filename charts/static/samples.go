package static

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
)

// Sample is one demo chart: its type + props. Mirrors @nivo/static samples
// entry shape { type, props }.
type Sample struct {
	Type  ChartType
	Props any
}

// Samples is the demo data registry for bar/line/pie. Mirrors @nivo/static
// samples (v1 only includes the three chart types implemented).
var Samples = map[ChartType]Sample{
	ChartTypeBar: {
		Type: ChartTypeBar,
		Props: bar.BarProps{
			Width:   1400,
			Height:  600,
			Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"},
			IndexBy: "country",
			Colors:  colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},
			Margin:  margin(40, 50, 40, 50),
			Data:    barSampleData(),
		},
	},
	ChartTypeLine: {
		Type: ChartTypeLine,
		Props: line.LineProps{
			Width:  800,
			Height: 500,
			Curve:  "monotoneX",
			Margin: margin(40, 50, 40, 50),
			Data:   lineSampleData(),
		},
	},
	ChartTypePie: {
		Type: ChartTypePie,
		Props: pie.PieProps{
			Width:        800,
			Height:       800,
			InnerRadius:  0.5,
			PadAngle:     0.5,
			CornerRadius: 5,
			Margin:       margin(100, 100, 100, 100),
			Data:         pieSampleData(),
		},
	},
}

// margin is a helper to build a core.Margin for the chart props structs.
func margin(top, right, bottom, left float64) core.Margin {
	return core.Margin{Top: top, Right: right, Bottom: bottom, Left: left}
}

// barSampleData generates 7 countries × 6 keys demo data. Mirrors nivo's
// generateCountriesData(keys, { size: 24 }) but with a fixed small set.
func barSampleData() []bar.BarDatum {
	countries := []string{"USA", "Germany", "France", "Japan", "Brazil", "India", "China"}
	keys := []string{"hot dogs", "burgers", "sandwich", "kebab", "fries", "donut"}
	// Deterministic pseudo-random values (nivo uses a seeded RNG; we use a
	// simple fixed pattern so golden tests are stable).
	data := make([]bar.BarDatum, len(countries))
	for i, c := range countries {
		row := map[string]any{"country": c}
		for j, k := range keys {
			row[k] = float64((i*7+j*13)%100 + 10)
		}
		data[i] = row
	}
	return data
}

// lineSampleData generates drink-stats-style demo data: 3 countries × 7 years.
func lineSampleData() []line.LineSeries {
	years := []string{"2018", "2019", "2020", "2021", "2022", "2023", "2024"}
	countries := []string{"USA", "Germany", "France"}
	series := make([]line.LineSeries, len(countries))
	for i, c := range countries {
		data := make([]line.LinePointData, len(years))
		for j, y := range years {
			data[j] = line.LinePointData{
				X: y,
				Y: float64((i*3+j*11)%90 + 20),
			}
		}
		series[i] = line.LineSeries{ID: c, Data: data}
	}
	return series
}

// pieSampleData generates programming-language-style demo data: 9 slices.
func pieSampleData() []any {
	labels := []string{"Go", "Rust", "TypeScript", "Python", "Ruby", "Elixir", "C", "Zig", "Kotlin"}
	data := make([]any, len(labels))
	for i, l := range labels {
		data[i] = map[string]any{
			"id":    l,
			"value": float64((i*17+3)%80 + 10),
		}
	}
	return data
}
