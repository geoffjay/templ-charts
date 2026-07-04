package samples

import "github.com/geoffjay/templ-charts/charts/boxplot"

// BoxPlot returns ready-to-render boxplot data: a deterministic spread of raw
// observations across four groups (color-by "group"). Each group holds 30
// values, summarized to quantiles by the chart.
func BoxPlot() []boxplot.BoxPlotDatum {
	var out []boxplot.BoxPlotDatum
	groups := []string{"Alpha", "Beta", "Gamma", "Delta"}
	for gi, g := range groups {
		for i := 0; i < 30; i++ {
			v := float64((i*9+gi*13)%50) + float64(gi)*8 + 20
			out = append(out, boxplot.BoxPlotDatum{Group: g, Value: v})
		}
	}
	return out
}
