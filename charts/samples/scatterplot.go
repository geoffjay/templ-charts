package samples

import "github.com/geoffjay/templ-charts/charts/scatterplot"

// ScatterPlot returns two series of six x/y observations each ("group A" and
// "group B"), ready to assign to a scatterplot.ScatterPlotProps.Data field.
func ScatterPlot() []scatterplot.ScatterPlotSerie {
	return []scatterplot.ScatterPlotSerie{
		{ID: "group A", Data: []scatterplot.ScatterPlotDatum{
			{X: 8.0, Y: 14.0}, {X: 22.0, Y: 40.0}, {X: 35.0, Y: 9.0}, {X: 51.0, Y: 60.0}, {X: 64.0, Y: 28.0}, {X: 78.0, Y: 73.0},
		}},
		{ID: "group B", Data: []scatterplot.ScatterPlotDatum{
			{X: 12.0, Y: 55.0}, {X: 27.0, Y: 22.0}, {X: 40.0, Y: 78.0}, {X: 58.0, Y: 33.0}, {X: 70.0, Y: 90.0}, {X: 88.0, Y: 47.0},
		}},
	}
}
