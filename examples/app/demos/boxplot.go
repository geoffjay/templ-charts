package demos

import (
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/samples"
)

// BoxPlotDemo is one boxplot tile on the /boxplot page.
type BoxPlotDemo struct {
	ID          string
	Title       string
	Description string
	Props       boxplot.BoxPlotProps
}

// BoxPlotDemos returns the boxplot demos for the /boxplot page.
func BoxPlotDemos() []BoxPlotDemo {
	return []BoxPlotDemo{
		{
			ID:          "boxplot-vertical",
			Title:       "Quantile boxes (vertical)",
			Description: "Raw observations per group summarized to the five-number quantile summary (q10/q25/median/q75/q90) and drawn as box + whisker glyphs.",
			Props: boxplot.BoxPlotProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 60},
				Data:        boxplotData(),
				ColorBy:     "group",
				Interactive: true,
			},
		},
	}
}

// boxplotData is the shared boxplot dataset, sourced from the public samples
// package.
func boxplotData() []boxplot.BoxPlotDatum {
	return samples.BoxPlot()
}
