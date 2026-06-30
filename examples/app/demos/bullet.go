package demos

import (
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/core"
)

// BulletDemo is one bullet tile on the /bullet page.
type BulletDemo struct {
	ID          string
	Title       string
	Description string
	Props       bullet.BulletProps
}

// BulletDemos returns the bullet demos for the /bullet page.
func BulletDemos() []BulletDemo {
	data := []bullet.BulletItemDatum{
		{ID: "temperature", Ranges: []float64{0, 20, 60, 100}, Measures: []float64{55}, Markers: []float64{72}},
		{ID: "power", Ranges: []float64{0, 30, 70, 120}, Measures: []float64{48, 88}, Markers: []float64{95}},
		{ID: "volume", Ranges: []float64{0, 40, 80, 140}, Measures: []float64{110}, Markers: []float64{100}},
	}
	return []BulletDemo{
		{
			ID:          "bullet-horizontal",
			Title:       "Horizontal (default)",
			Description: "Three KPI rows: stacked qualitative ranges (seq:cool), measure bars (seq:red_purple), and comparative markers, with a per-row axis.",
			Props: bullet.BulletProps{
				Width: commonChartWidth, Height: 260,
				Margin:       core.Margin{Top: 20, Right: 30, Bottom: 40, Left: 100},
				Data:         data,
				TitleOffsetX: -80,
				Interactive:  true,
			},
		},
	}
}
