package samples

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/swarmplot"
)

// SwarmPlot returns 36 values spread across three groups (group A/B/C) using a
// fixed value sequence so the force layout renders deterministically. The
// returned data and groups populate a swarmplot.SwarmPlotProps' Data and Groups
// fields.
func SwarmPlot() ([]swarmplot.SwarmPlotDatum, []string) {
	groups := []string{"group A", "group B", "group C"}
	seq := []float64{
		12, 47, 23, 68, 34, 89, 5, 56, 78, 30, 41, 62,
		19, 73, 27, 51, 44, 8, 95, 60, 15, 38, 82, 49,
		22, 66, 11, 90, 33, 57, 71, 4, 84, 26, 53, 40,
	}
	data := make([]swarmplot.SwarmPlotDatum, 0, len(seq))
	for i, v := range seq {
		data = append(data, swarmplot.SwarmPlotDatum{
			ID: "n" + strconv.Itoa(i), Group: groups[i%len(groups)], Value: v,
		})
	}
	return data, groups
}
