package demos

import (
	"github.com/geoffjay/templ-charts/charts/boxplot"
	"github.com/geoffjay/templ-charts/charts/core"
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

// boxplotData generates a deterministic spread per group (no randomness so the
// page render is stable).
func boxplotData() []boxplot.BoxPlotDatum {
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
