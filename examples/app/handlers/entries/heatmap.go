package entries

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// heatmapDemoData builds an 8 country x 10 transport-mode grid with a
// smooth-ish deterministic value surface (usage share, 0-100) so the
// sequential color scale shows a broad, continuous range.
func heatmapDemoData() []heatmap.HeatMapSerie {
	v := func(f float64) *float64 { return &f }
	rows := []string{"Japan", "France", "Germany", "Brazil", "Canada", "Norway", "Kenya", "India"}
	cols := []string{"Train", "Subway", "Bus", "Car", "Bike", "Scooter", "Taxi", "Ferry", "Tram", "Walk"}
	data := make([]heatmap.HeatMapSerie, len(rows))
	for ri, r := range rows {
		cells := make([]heatmap.HeatMapDatum, len(cols))
		for ci, c := range cols {
			val := 52 + 44*math.Sin(float64(ri)*0.55+float64(ci)*0.38)*math.Cos(float64(ci)*0.21-float64(ri)*0.33)
			cells[ci] = heatmap.HeatMapDatum{X: c, Y: v(math.Round(val))}
		}
		data[ri] = heatmap.HeatMapSerie{ID: r, Data: cells}
	}
	return data
}

func init() {
	register(ChartEntry{
		Slug:        "heatmap",
		Title:       "Heatmap",
		Description: "2D value grid colored through a sequential/diverging scale.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/heatmap"
    "github.com/geoffjay/templ-charts/charts/render"
)

v := func(f float64) *float64 { return &f }
svg, _ := render.String(heatmap.HeatMap(heatmap.HeatMapProps{
    Width: 720, Height: 440,
    Data: []heatmap.HeatMapSerie{
        {ID: "Japan", Data: []heatmap.HeatMapDatum{
            {X: "Train", Y: v(42)}, {X: "Bus", Y: v(18)},
        }},
        {ID: "France", Data: []heatmap.HeatMapDatum{
            {X: "Train", Y: v(33)}, {X: "Bus", Y: v(51)},
        }},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := heatmap.HeatMapProps{
				Width: 720, Height: 440, Responsive: true,
				Data:  heatmapDemoData(),
				Theme: theme,
			}
			p.Animate = animate
			return render.String(heatmap.HeatMap(p))
		},
		SpaceRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool, space colors.Space) (string, error) {
			// A wide value ramp under a multi-hue sequential scheme, so the
			// difference between RGB and Lab/Lch interpolation is visible.
			v := func(f float64) *float64 { return &f }
			cols := []string{"c0", "c1", "c2", "c3", "c4", "c5", "c6", "c7"}
			rows := []string{"low", "mid", "high"}
			data := make([]heatmap.HeatMapSerie, len(rows))
			for ri, r := range rows {
				cells := make([]heatmap.HeatMapDatum, len(cols))
				for ci, c := range cols {
					cells[ci] = heatmap.HeatMapDatum{X: c, Y: v(float64(ci*12 + ri*4))}
				}
				data[ri] = heatmap.HeatMapSerie{ID: r, Data: cells}
			}
			// The selected palette drives the sequential scale when it is a
			// gradient (sequential/diverging); categorical picks keep the default
			// multi-hue ramp, which shows the RGB vs Lab/Lch difference well.
			scheme := "yellow_orange_red"
			if pal, ok := colors.LookupPalette(palette); ok && pal.Kind != colors.KindCategorical {
				scheme = string(palette)
			}
			p := heatmap.HeatMapProps{
				Width: 720, Height: 320, Responsive: true,
				Data:   data,
				Colors: heatmap.HeatMapColorConfig{Type: "sequential", Scheme: scheme, Space: space},
				Theme:  theme,
			}
			p.Animate = animate
			return render.String(heatmap.HeatMap(p))
		},
		CanvasRender: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := heatmap.HeatMapProps{
				Width: 720, Height: 440, Responsive: true,
				Data:    heatmapDemoData(),
				Theme:   theme,
				Render:  theming.EngineCanvas,
				ChartID: "detail-heatmap-canvas",
			}
			return render.String(heatmap.HeatMap(p))
		},
	})
}
